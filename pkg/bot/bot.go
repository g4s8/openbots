package bot

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/g4s8/openbots/internal/bot/adaptors"
	"github.com/g4s8/openbots/internal/bot/handlers"
	"github.com/g4s8/openbots/pkg/api"
	ctx "github.com/g4s8/openbots/pkg/context"
	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	"github.com/pkg/errors"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ types.Bot = (*Bot)(nil)

type eventHandler struct {
	types.EventFilter
	types.Handler
}

type stateHandler struct {
	types.EventFilter
	types.StateHandler
}

type Bot struct {
	handlers      []*eventHandler
	stateHandlers []*stateHandler
	apiHandlers   map[string][]api.Handler

	apiAddr string

	context    types.ContextProvider
	state      types.StateProvider
	botAPI     *telegram.BotAPI
	apiService *api.Service
	stopOnce   sync.Once
	quitCh     chan struct{}
	doneCh     chan struct{}
}

func New(botAPI *telegram.BotAPI, state types.StateProvider, context types.ContextProvider,
	apiAddr string) *Bot {
	return &Bot{
		handlers:      make([]*eventHandler, 0),
		stateHandlers: make([]*stateHandler, 0),
		apiHandlers:   make(map[string][]api.Handler),
		apiAddr:       apiAddr,
		context:       context,
		state:         state,
		botAPI:        botAPI,
		quitCh:        make(chan struct{}, 1),
		doneCh:        make(chan struct{}, 1),
	}
}

func NewFromSpec(s *spec.Bot) (*Bot, error) {
	botAPI, err := telegram.NewBotAPI(s.Token)
	if err != nil {
		return nil, errors.Wrap(err, "create bot API")
	}
	botAPI.Debug = s.Debug

	botID := botAPI.Self.ID

	var (
		sp types.StateProvider
		cp types.ContextProvider
	)
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
		log.Println("Connecting database")
		db, err := sql.Open("postgres", conString)
		if err != nil {
			return nil, errors.Wrap(err, "open database")
		}
		if err := db.Ping(); err != nil {
			return nil, errors.Wrap(err, "ping database")
		}
		log.Println("Database connected")
		sp = state.NewDB(db, botID)
		cp = ctx.NewDBProvider(db, botID)
	}

	var apiAddr string
	if s.Config.Api != nil {
		apiAddr = s.Config.Api.Address
	}
	bot := New(botAPI, sp, cp, apiAddr)

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
			filter       types.EventFilter
			hs           []types.Handler
			stateHandler types.StateHandler
			err          error
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
		if h.Trigger.Context != "" && filter != nil {
			filter = handlers.NewContextFilter(filter, b.context, h.Trigger.Context)
		}
		if h.Replies != nil {
			hs = append(hs, adaptors.Replies(b.state, h.Replies))
		}
		if h.State != nil {
			stateHandler = handlers.NewStateHandlerFromSpec(b.state, h.State)
			if err != nil {
				return errors.Wrap(err, "create state handler")
			}
		}
		if len(hs) > 0 && h.Context != nil {
			for i, han := range hs {
				if h.Context.Set != "" {
					hs[i] = handlers.NewContextSetter(han, b.context, h.Context.Set)
				}
				if h.Context.Delete != "" {
					hs[i] = handlers.NewContextDeleter(han, b.context, h.Context.Delete)
				}
			}
		}
		if h.Webhook != nil {
			hs = append(hs, adaptors.Webhook(h.Webhook))
		}

		if filter == nil {
			return errors.New("no event filter")
		}
		if len(hs) == 0 && stateHandler == nil {
			return errors.New("no handler")
		}
		for _, h := range hs {
			b.Handle(filter, h)
		}
		if stateHandler != nil {
			b.HandleState(filter, stateHandler)
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

func (b *Bot) HandleState(filter types.EventFilter, h types.StateHandler) {
	b.stateHandlers = append(b.stateHandlers, &stateHandler{EventFilter: filter, StateHandler: h})
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
	for _, h := range b.handlers {
		if check, err := h.Check(ctx, upd); err != nil {
			log.Printf("Filter error: %v", err)
			continue
		} else if !check {
			continue
		}

		if err := h.Handle(ctx, upd, b.botAPI); err != nil {
			log.Printf("Handler error: %v", err)
			continue
		}
	}
	for _, sh := range b.stateHandlers {
		if check, err := sh.Check(ctx, upd); err != nil {
			log.Printf("State filter error: %v", err)
			continue
		} else if !check {
			continue
		}

		if err := sh.Handle(ctx, upd); err != nil {
			log.Printf("State handler error: %v", err)
			continue
		}
	}
}
