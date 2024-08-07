package bot

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"slices"
	"sync"
	"time"

	goerr "errors"

	"github.com/g4s8/openbots/internal/bot/adaptors"
	internal_api "github.com/g4s8/openbots/internal/bot/api"
	botctx "github.com/g4s8/openbots/internal/bot/ctx"
	"github.com/g4s8/openbots/internal/bot/data"
	"github.com/g4s8/openbots/internal/bot/filters"
	"github.com/g4s8/openbots/internal/bot/handlers"
	"github.com/g4s8/openbots/internal/bot/logger"
	"github.com/g4s8/openbots/pkg/api"
	"github.com/g4s8/openbots/pkg/assets"
	ctx "github.com/g4s8/openbots/pkg/context"
	logwrap "github.com/g4s8/openbots/pkg/log"
	"github.com/g4s8/openbots/pkg/payments"
	"github.com/g4s8/openbots/pkg/secrets"
	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ types.Bot = (*Bot)(nil)

type eventHandler struct {
	types.EventFilter
	types.Handler
	types.DataLoader
}

// Bot is a main bot instance.
type Bot struct {
	botAPI   *telegram.BotAPI
	apiAddr  string
	cp       *botctx.Provider
	state    types.StateProvider
	assets   types.Assets
	payments types.PaymentProviders
	secrets  types.Secrets
	httpCli  *http.Client
	ucp      *handlers.UpdateContextProvider
	log      zerolog.Logger

	handlers    []*eventHandler
	apiHandlers map[string][]api.Handler
	apiService  *api.Service

	stopOnce sync.Once
	quitCh   chan struct{}
	doneCh   chan struct{}
}

// NewWithOptions creates a new bot instance with options or default values for empty options.
func NewWithOptions(botAPI *telegram.BotAPI, opts ...Option) *Bot {
	b := &Bot{
		handlers:    make([]*eventHandler, 0),
		apiHandlers: make(map[string][]api.Handler),
		botAPI:      botAPI,
		quitCh:      make(chan struct{}, 1),
		doneCh:      make(chan struct{}, 1),
		log:         zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(b)
	}

	if b.httpCli == nil {
		b.httpCli = http.DefaultClient
	}
	if b.state == nil {
		b.state = state.NewMemory(nil)
	}
	if b.assets == nil {
		b.assets = assets.Dummy
	}
	if b.cp == nil {
		b.cp = botctx.NewProvider(ctx.NewMemoryProvider())
	}
	if b.payments == nil {
		b.payments = payments.EmptyProvider
	}
	if b.secrets == nil {
		b.secrets = secrets.Stub
	}
	b.ucp = handlers.NewUpdateContextProvider(b.secrets, b.state)

	return b
}

// New creates a new bot instance with default values for empty options.
// deprecated: use NewWithOptions instead
func New(botAPI *telegram.BotAPI, state types.StateProvider,
	context types.ContextProvider, assets types.Assets,
	paymentProviders types.PaymentProviders, secrets types.Secrets,
	apiAddr string, log zerolog.Logger,
) *Bot {
	return NewWithOptions(botAPI,
		WithStateProvider(state),
		WithContextProvider(context),
		WithAssets(assets),
		WithPaymentProviders(paymentProviders),
		WithSecrets(secrets),
		WithAPIAddr(apiAddr),
		WithLogger(log))
}

func NewFromSpec(s *spec.Bot) (*Bot, error) {
	botAPI, err := telegram.NewBotAPI(s.Token)
	if err != nil {
		return nil, errors.Wrap(err, "create bot API")
	}
	botAPI.Debug = s.Debug

	botID := botAPI.Self.ID

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	if s.Debug {
		log = log.Level(zerolog.DebugLevel)
	} else {
		log = log.Level(zerolog.InfoLevel)
	}

	var (
		sp types.StateProvider
		cp types.ContextProvider
		ap types.Assets
	)

	if s.Config == nil {
		s.Config = spec.DefaultConfig
	}
	if s.Config.Persistence == nil {
		s.Config.Persistence = spec.DefaultConfig.Persistence
	}
	if s.Config.Assets == nil {
		s.Config.Assets = spec.DefaultConfig.Assets
	}

	switch s.Config.Persistence.Type {
	case spec.MemoryPersistence:
		sp = state.NewMemory(s.State)
		cp = ctx.NewMemoryProvider()
	case spec.DatabasePersistence:
		conString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
			s.Config.Persistence.DBConfig.Host, s.Config.Persistence.DBConfig.Port,
			s.Config.Persistence.DBConfig.User, s.Config.Persistence.DBConfig.Password,
			s.Config.Persistence.DBConfig.Database)
		if s.Config.Persistence.DBConfig.NoSSL {
			conString += " sslmode=disable"
		}
		log.Debug().Str("host", s.Config.Persistence.DBConfig.Host).Msg("Connecting database")
		db, err := sql.Open("postgres", conString)
		if err != nil {
			return nil, errors.Wrap(err, "open database")
		}
		if err := db.Ping(); err != nil {
			return nil, errors.Wrap(err, "ping database")
		}
		log.Debug().Msg("Database connected")
		sp = state.NewDB(db, botID)
		cp = ctx.NewDBProvider(db, botID)
	}

	var apiAddr string
	if s.Config.Api != nil {
		apiAddr = s.Config.Api.Address
	}

	if s.Config.Assets.Provider == "fs" {
		var root string
		if r, ok := s.Config.Assets.Params["root"]; ok {
			root = r
		} else {
			wd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "get working directory")
			}
			root = wd
		}
		ap = assets.NewFS(root, log)
	}

	var paymentProviders types.PaymentProviders
	if len(s.Config.PaymentProviders) > 0 {
		m := make(map[string]string, len(s.Config.PaymentProviders))
		for _, p := range s.Config.PaymentProviders {
			m[p.Name] = p.Token
		}
		paymentProviders = payments.NewMapProvider(m)
	} else {
		paymentProviders = payments.EmptyProvider
	}

	sp = logwrap.WrapStateProvider(sp, log)
	cp = logwrap.WrapContextProvider(cp, log)

	bot := New(botAPI, sp, cp, ap, paymentProviders, secrets.Stub, apiAddr, log)

	if err := bot.SetupHandlersFromSpec(s.Handlers); err != nil {
		return nil, errors.Wrap(err, "setup handlers")
	}

	if s.Api != nil {
		if err := bot.SetupApiHandlersFromSpec(s.Api.Handlers); err != nil {
			return nil, errors.Wrap(err, "setup api handlers")
		}
	}

	return bot, nil
}

func (b *Bot) SetupHandlersFromSpec(src []*spec.Handler) error {
	for _, h := range src {
		var (
			filter types.EventFilter
			hs     []types.Handler
			dl     types.DataLoader
			err    error
		)

		if h.Trigger.Message != nil {
			filter, err = handlers.NewMessageFilterFromSpec(h.Trigger.Message)
			if err != nil {
				return errors.Wrap(err, "create message event filter")
			}
		}
		if h.Trigger.Callback != nil {
			filter, err = handlers.NewCallbackFilterFromSpec(h.Trigger.Callback)
			if err != nil {
				return errors.Wrap(err, "create callback event filter")
			}
		}
		if h.Trigger.Context != "" {
			filter = handlers.NewContextFilter(filter, b.cp, h.Trigger.Context)
		}
		// TODO: refactor all filters/triggers similat to handlers
		if h.Trigger.PreCheckout != nil {
			filter = adaptors.NewPrecheckoutFilter(h.Trigger.PreCheckout)
		}
		if h.Trigger.PostCheckout != nil {
			filter = adaptors.NewPostcheckoutFilter(h.Trigger.PostCheckout)
		}
		if len(h.Trigger.State) > 0 {
			f := adaptors.NewStateFilter(b.state, b.log, h.Trigger.State)
			filter = filters.Join(filter, f)
		}
		if filter == nil && h.Trigger.Fallback {
			filter = filters.Fallback
		}

		// validator should be the first handler
		if v := h.Validate; v != nil {
			h, err := adaptors.Validator(v, b.log)
			if err != nil {
				return errors.Wrap(err, "create validator handler")
			}
			hs = append(hs, h)
		}
		if h.Replies != nil {
			h, err := adaptors.Replies(b.botAPI, b.state, b.secrets, b.assets, b.payments, h.Replies, b.log)
			if err != nil {
				return errors.Wrap(err, "create replies handler")
			}
			hs = append(hs, h)
		}
		if h.State != nil {
			hs = append(hs, handlers.NewStateHandlerFromSpec(b.state, h.State, b.log))
		}
		if h.Context != nil {
			if h.Context.Set != "" {
				hs = append(hs, handlers.NewContextSetter(b.cp, h.Context.Set, b.log))
			}
			if h.Context.Delete != "" {
				hs = append(hs, handlers.NewContextDeleter(b.cp, h.Context.Delete, b.log))
			}
		}
		if h.Webhook != nil {
			hs = append(hs, adaptors.Webhook(h.Webhook, b.httpCli, b.state, b.secrets, b.log))
		}
		if h.Data != nil {
			d, err := adaptors.DataLoader(b.httpCli, b.state, b.secrets, h.Data, b.log)
			if err != nil {
				return errors.Wrap(err, "create data loader")
			}
			dl = d
		}

		if filter == nil {
			return errors.New("no event filter")
		}
		if len(hs) == 0 {
			return errors.New("no handler")
		}
		for _, h := range hs {
			b.HandleWithData(filter, h, dl)
		}
	}
	return nil
}

func (b *Bot) SetupApiHandlersFromSpec(src []*spec.ApiHandler) error {
	for _, h := range src {
		for _, act := range h.Actions {
			var hs []api.Handler
			if act.SendMessage != nil {
				reply, err := adaptors.MessageRepply(b.botAPI, b.state, b.secrets, act.SendMessage,
					b.log.With().Str("component", "api").Str("handler", h.ID).Logger())
				if err != nil {
					return errors.Wrap(err, "create api message reply handler")
				}
				hs = append(hs, reply)
			}

			if act.Context != nil {
				if act.Context.Set != "" {
					hs = append(hs, handlers.NewContextSetter(b.cp, act.Context.Set, b.log))
				}
				if act.Context.Delete != "" {
					hs = append(hs, handlers.NewContextDeleter(b.cp, act.Context.Delete, b.log))
				}
			}

			if act.State != nil {
				hs = append(hs, handlers.NewStateHandlerFromSpec(b.state, act.State, b.log))
			}

			for _, step := range hs {
				b.ApiHandler(h.ID, internal_api.NewSendMessage(act.ChatID, step))
			}
			b.log.Info().Int("handler", len(hs)).Str("id", h.ID).Msg("api handler registered")
		}
	}
	return nil
}

func (b *Bot) Handle(filter types.EventFilter, h types.Handler) {
	b.handlers = append(b.handlers, &eventHandler{EventFilter: filter, Handler: h})
	// fallback handler should be the last handler
	slices.SortFunc(b.handlers, func(left, right *eventHandler) int {
		leftFallback, rightFallback := left.EventFilter == filters.Fallback, right.EventFilter == filters.Fallback
		if leftFallback {
			return 1
		}
		if rightFallback {
			return -1
		}
		return 0
	})
}

func (b *Bot) HandleWithData(filter types.EventFilter, h types.Handler, dl types.DataLoader) {
	b.handlers = append(b.handlers, &eventHandler{EventFilter: filter, Handler: h, DataLoader: dl})
}

func (b *Bot) ApiHandler(id string, h api.Handler) {
	list, ok := b.apiHandlers[id]
	if !ok {
		list = make([]api.Handler, 0)
	}
	b.apiHandlers[id] = append(list, h)
}

func (b *Bot) Start() error {
	_ = telegram.SetLogger(logger.Wrap(b.log.With().Str("component", "telegram").Logger(), zerolog.InfoLevel))

	updCfg := telegram.NewUpdate(0)
	updCfg.Timeout = 30
	updCh := b.botAPI.GetUpdatesChan(updCfg)
	go func() {
		defer close(b.doneCh)
		for {
			select {
			case <-b.quitCh:
				return
			case upd := <-updCh:
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				b.HandleUpdate(ctx, &upd)
				cancel()
			}
		}
	}()
	if b.apiAddr != "" {
		b.apiService = b.HandlerAPI(api.Config{
			Addr:           b.apiAddr,
			ReadTimeout:    time.Second * 5,
			RequestTimeout: time.Second * 3,
		})
		if err := b.apiService.Start(context.TODO()); err != nil {
			return errors.Wrap(err, "start api service")
		}
		b.log.Info().Str("addr", b.apiAddr).Msg("API service started")
	}
	b.log.Info().Msg("Bot started")
	return nil
}

func (b *Bot) HandlerAPI(cfg api.Config) *api.Service {
	handlers := make(map[string]api.Handler, len(b.apiHandlers))
	for id, hs := range b.apiHandlers {
		handlers[id] = &apiHandlerGroup{handlers: hs}
	}
	return api.NewServiceWithLogger(cfg, handlers, b.log.With().Str("component", "api_svc").Logger())
}

func (b *Bot) Stop() error {
	b.stopOnce.Do(func() {
		b.log.Info().Msg("Stopping bot")
		b.botAPI.StopReceivingUpdates()
		close(b.quitCh)
		if b.apiService != nil {
			if err := b.apiService.Stop(context.TODO()); err != nil {
				b.log.Error().Err(err).Msg("Error stopping API service")
			}
		}
		<-b.doneCh
		b.log.Info().Msg("Bot stopped")
	})
	return nil
}

// HandleUpdate handles telegram update and log error if any.
// Deprecated: use HandleUpdateErr instead, in next major release HandlerUpdateErr will be renamed to HandleUpdate.
func (b *Bot) HandleUpdate(ctx context.Context, upd *telegram.Update) {
	if err := b.HandleUpdateErr(ctx, upd); err != nil {
		b.log.Error().Err(err).Msg("Handle update")
	}
}

type handlerWithData struct {
	handler types.Handler
	data    types.DataLoader
}

// HandleUpdateErr handles telegram update and returns error if any.
func (b *Bot) HandleUpdateErr(ctx context.Context, upd *telegram.Update) error {
	chatID := handlers.ChatID(upd)
	log := b.log.With().Str("chat_id", chatID.String()).Logger()
	log.Debug().Msg("Handling update")
	ctxCloser := b.cp.Begin(chatID)
	defer func() {
		if err := ctxCloser(ctx); err != nil {
			log.Error().Err(err).Msg("Close context")
		}
		log.Trace().Msg("Context closed")
	}()

	uctx, err := b.ucp.NewContext(ctx, upd)
	if err != nil {
		return errors.Wrap(err, "create update context")
	}
	ctx = uctx

	var errs []error
	hs := make([]handlerWithData, 0)
	var fallbackHandler handlerWithData
	for _, h := range b.handlers {
		if h.EventFilter == filters.Fallback {
			fallbackHandler = handlerWithData{handler: h.Handler, data: h.DataLoader}
			continue
		}
		if check, err := h.Check(ctx, upd); err != nil {
			errs = append(errs, errors.Wrap(err, "filter check"))
			continue
		} else if check {
			hs = append(hs, handlerWithData{handler: h.Handler, data: h.DataLoader})
		}
	}

	log.Debug().Int("handlers", len(hs)).Msg("Handlers found")
	var handled bool
	for i, h := range hs {
		log.Debug().Int("handler", i).Msg("Handling")

		ok, err := runHandler(ctx, b.botAPI, h, upd)
		if !handled && ok {
			handled = true
		}
		if err != nil {
			if errors.Is(err, handlers.ErrValidationFailed) {
				log.Info().Err(err).Msg("Validation failed")
				return nil
			}
			errs = append(errs, err)
		}
	}

	if !handled && fallbackHandler.handler != nil {
		log.Debug().Msg("Handling fallback")
		ok, err := runHandler(ctx, b.botAPI, fallbackHandler, upd)
		if err != nil {
			if errors.Is(err, handlers.ErrValidationFailed) {
				log.Info().Err(err).Msg("Validation failed")
				return nil
			}
			errs = append(errs, err)
		}
		if ok {
			log.Debug().Msg("Fallback handled")
		}
	}

	log.Debug().Msg("Update handled")
	return goerr.Join(errs...)
}

func runHandler(ctx context.Context, botAPI *telegram.BotAPI, h handlerWithData, upd *telegram.Update) (bool, error) {
	if h.data != nil {
		var c types.DataContainer
		if err := h.data.Load(ctx, &c, upd); err != nil {
			return false, errors.Wrap(err, "load data")
		}
		ctx = data.ContextWithContainer(ctx, &c)
	}

	if err := h.handler.Handle(ctx, upd, botAPI); err != nil {
		return true, errors.Wrap(err, "handler")
	}
	return true, nil
}
