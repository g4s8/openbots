package bot

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/g4s8/openbots-go/internal/bot/adaptors"
	"github.com/g4s8/openbots-go/internal/bot/handlers"
	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/g4s8/openbots-go/pkg/types"
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

	context  *types.Context
	state    types.State
	botAPI   *telegram.BotAPI
	stopOnce sync.Once
	quitCh   chan struct{}
	doneCh   chan struct{}
}

func New(botAPI *telegram.BotAPI) *Bot {
	return &Bot{
		handlers:      make([]*eventHandler, 0),
		stateHandlers: make([]*stateHandler, 0),
		context:       new(types.Context),
		state:         types.NewState(nil),
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
	botAPI.Debug = true
	bot := New(botAPI)
	bot.state = types.NewState(s.State)
	for _, h := range s.Handlers {
		var filter types.EventFilter
		var hs []types.Handler
		var stateHandler types.StateHandler

		if h.Trigger.Message != nil {
			filter, err = handlers.NewMessageFilterFromSpec(h.Trigger.Message)
			if err != nil {
				return nil, errors.Wrap(err, "create message event filter")
			}
		}
		if h.Trigger.Callback != nil {
			filter, err = handlers.NewCallbackFilterFromSpec(h.Trigger.Callback)
			if err != nil {
				return nil, errors.Wrap(err, "create callback event filter")
			}
		}
		if h.Trigger.Context != "" && filter != nil {
			filter = handlers.NewContextFilter(filter, bot.context, h.Trigger.Context)
		}
		if h.Replies != nil {
			hs = append(hs, adaptors.Replies(h.Replies))
		}
		if h.State != nil {
			stateHandler = handlers.NewStateHandlerFromSpec(h.State)
			if err != nil {
				return nil, errors.Wrap(err, "create state handler")
			}
		}
		if len(hs) > 0 && h.Context != nil {
			for i, han := range hs {
				if h.Context.Set != "" {
					hs[i] = handlers.NewContextSetter(han, bot.context, h.Context.Set)
				}
				if h.Context.Delete != "" {
					hs[i] = handlers.NewContextDeleter(han, bot.context, h.Context.Delete)
				}
			}
		}
		if h.Webhook != nil {
			hs = append(hs, adaptors.Webhook(h.Webhook))
		}

		if filter == nil {
			return nil, errors.New("no event filter")
		}
		if len(hs) == 0 && stateHandler == nil {
			return nil, errors.New("no handler")
		}
		for _, h := range hs {
			bot.Handle(filter, h)
		}
		if stateHandler != nil {
			bot.HandleState(filter, stateHandler)
		}
	}
	return bot, nil
}

func (b *Bot) Handle(filter types.EventFilter, h types.Handler) {
	b.handlers = append(b.handlers, &eventHandler{EventFilter: filter, Handler: h})
}

func (b *Bot) HandleState(filter types.EventFilter, h types.StateHandler) {
	b.stateHandlers = append(b.stateHandlers, &stateHandler{EventFilter: filter, StateHandler: h})
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
				b.handleUpdate(&upd)
			}
		}
	}()
	log.Print("Bot started")
	return nil
}

func (b *Bot) Stop() error {
	b.stopOnce.Do(func() {
		log.Println("Stopping bot")
		b.botAPI.StopReceivingUpdates()
		close(b.quitCh)
		<-b.doneCh
		log.Print("Bot stopped")
	})
	return nil
}

func (b *Bot) handleUpdate(upd *telegram.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userID := handlers.ChatID(upd)
	ctx = types.ContextWithState(ctx, userID, b.state)
	for _, h := range b.handlers {
		if !h.Check(ctx, upd) {
			continue
		}
		if err := h.Handle(ctx, upd, b.botAPI); err != nil {
			log.Printf("Handler error: %v\n", err)
		}
	}
	state := b.state.User(userID)
	for _, h := range b.stateHandlers {
		var err error
		if !h.Check(ctx, upd) {
			continue
		}
		state, err = h.Handle(ctx, upd, state)
		if err != nil {
			log.Printf("State handler error: %v\n", err)
		}
	}
	b.state.Save(userID, state)
}
