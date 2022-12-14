package bot

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/g4s8/openbots/internal/bot/adaptors"
	"github.com/g4s8/openbots/internal/bot/handlers"
	"github.com/g4s8/openbots/pkg/api"
	"github.com/g4s8/openbots/pkg/assets"
	ctx "github.com/g4s8/openbots/pkg/context"
	logwrap "github.com/g4s8/openbots/pkg/log"
	"github.com/g4s8/openbots/pkg/payments"
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
}

type Bot struct {
	handlers    []*eventHandler
	apiHandlers map[string][]api.Handler

	apiAddr string

	context    types.ContextProvider
	state      types.StateProvider
	assets     types.Assets
	payments   types.PaymentProviders
	botAPI     *telegram.BotAPI
	apiService *api.Service
	stopOnce   sync.Once
	quitCh     chan struct{}
	doneCh     chan struct{}

	log zerolog.Logger
}

func New(botAPI *telegram.BotAPI, state types.StateProvider,
	context types.ContextProvider, assets types.Assets,
	paymentProviders types.PaymentProviders,
	apiAddr string, log zerolog.Logger) *Bot {
	return &Bot{
		handlers:    make([]*eventHandler, 0),
		apiHandlers: make(map[string][]api.Handler),
		apiAddr:     apiAddr,
		context:     context,
		state:       state,
		assets:      assets,
		payments:    paymentProviders,
		botAPI:      botAPI,
		quitCh:      make(chan struct{}, 1),
		doneCh:      make(chan struct{}, 1),
		log:         log,
	}
}

func NewFromSpec(s *spec.Bot) (*Bot, error) {
	botAPI, err := telegram.NewBotAPI(s.Token)
	if err != nil {
		return nil, errors.Wrap(err, "create bot API")
	}
	botAPI.Debug = s.Debug

	botID := botAPI.Self.ID

	log := zerolog.New(zerolog.ConsoleWriter{Out: log.Writer()}).With().Timestamp().Logger()
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
		log.Debug().Msg("Connecting database")
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

	bot := New(botAPI, sp, cp, ap, paymentProviders, apiAddr, log)

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
			filter = handlers.NewContextFilter(filter, b.context, h.Trigger.Context)
		}
		// TODO: refactor all filters/triggers similat to handlers
		if h.Trigger.PreCheckout != nil {
			filter = adaptors.NewPrecheckoutFilter(h.Trigger.PreCheckout)
		}
		if h.Trigger.PostCheckout != nil {
			filter = adaptors.NewPostcheckoutFilter(h.Trigger.PostCheckout)
		}

		if h.Replies != nil {
			hs = append(hs, adaptors.Replies(b.state, b.assets, b.payments, h.Replies, b.log))
		}
		if h.State != nil {
			hs = append(hs, handlers.NewStateHandlerFromSpec(b.state, h.State, b.log))
		}
		if h.Context != nil {
			if h.Context.Set != "" {
				hs = append(hs, handlers.NewContextSetter(b.context, h.Context.Set, b.log))
			}
			if h.Context.Delete != "" {
				hs = append(hs, handlers.NewContextDeleter(b.context, h.Context.Delete, b.log))
			}
		}
		if h.Webhook != nil {
			hs = append(hs, adaptors.Webhook(h.Webhook, b.state, b.log))
		}

		if filter == nil {
			return errors.New("no event filter")
		}
		if len(hs) == 0 {
			return errors.New("no handler")
		}
		for _, h := range hs {
			b.Handle(filter, h)
		}
	}
	return nil
}

func (b *Bot) SetupApiHandlersFromSpec(src []*spec.ApiHandler) error {
	for _, h := range src {
		for _, act := range h.Actions {
			if act.SendMessage != nil {
				b.ApiHandler(h.ID,
					adaptors.ApiSendMessage(b.botAPI, act.SendMessage))
			}
		}
	}
	return nil
}

func (b *Bot) Handle(filter types.EventFilter, h types.Handler) {
	b.handlers = append(b.handlers, &eventHandler{EventFilter: filter, Handler: h})
}

func (b *Bot) ApiHandler(id string, h api.Handler) {
	list, ok := b.apiHandlers[id]
	if !ok {
		list = make([]api.Handler, 0)
	}
	b.apiHandlers[id] = append(list, h)
}

func (b *Bot) Start() error {
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
		handlers := make(map[string]api.Handler, len(b.apiHandlers))
		for id, hs := range b.apiHandlers {
			handlers[id] = &apiHandlerGroup{handlers: hs}
		}
		b.apiService = api.NewService(api.Config{
			Addr:           b.apiAddr,
			ReadTimeout:    time.Second * 5,
			RequestTimeout: time.Second * 3,
		}, handlers)
		if err := b.apiService.Start(context.TODO()); err != nil {
			return errors.Wrap(err, "start api service")
		}
		log.Printf("API service started on `%s`", b.apiAddr)
	}
	log.Print("Bot started")
	return nil
}

func (b *Bot) Stop() error {
	b.stopOnce.Do(func() {
		log.Println("Stopping bot")
		b.botAPI.StopReceivingUpdates()
		close(b.quitCh)
		if b.apiService != nil {
			if err := b.apiService.Stop(context.TODO()); err != nil {
				log.Printf("Error stopping API service: %v", err)
			}
		}
		<-b.doneCh
		log.Print("Bot stopped")
	})
	return nil
}

func (b *Bot) HandleUpdate(ctx context.Context, upd *telegram.Update) {
	chatID := handlers.ChatID(upd)
	log := b.log.With().Str("chat-id", chatID.String()).Logger()
	log.Debug().Msg("Handling update")

	handlers := make([]types.Handler, 0)
	for _, h := range b.handlers {
		if check, err := h.Check(ctx, upd); err != nil {
			log.Printf("Filter error: %v", err)
			continue
		} else if check {
			handlers = append(handlers, h)
		}
	}
	log.Debug().Int("handlers", len(handlers)).Msg("Handlers found")
	for i, h := range handlers {
		log.Debug().Int("handler", i).Msg("Handling")
		if err := h.Handle(ctx, upd, b.botAPI); err != nil {
			log.Error().Err(err).Msg("Handler failed")
			continue
		}
	}
	log.Debug().Msg("Update handled")
}
