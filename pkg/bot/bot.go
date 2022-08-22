package bot

import (
	"context"
	"log"
	"time"

	"github.com/g4s8/openbots-go/internal/bot/handlers"
	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/g4s8/openbots-go/pkg/types"
	"github.com/pkg/errors"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ types.Bot = (*Bot)(nil)

type State map[string]string

type eventHandler struct {
	types.EventFilter
	types.Handler
}

type Bot struct {
	handlers []*eventHandler

	// state  State
	botAPI *telegram.BotAPI
	doneCh chan struct{}
}

func New(botAPI *telegram.BotAPI) *Bot {
	return &Bot{
		handlers: make([]*eventHandler, 0),
		botAPI:   botAPI,
		doneCh:   make(chan struct{}),
	}
}

func NewFromSpec(s *spec.Bot) (*Bot, error) {
	botAPI, err := telegram.NewBotAPI(s.Token)
	if err != nil {
		return nil, errors.Wrap(err, "create bot API")
	}
	botAPI.Debug = true
	bot := New(botAPI)
	for _, h := range s.Handlers {
		var filter types.EventFilter
		var handler types.Handler

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
		if h.Replies != nil {
			handler, err = handlers.NewReplyFromSpec(h.Replies)
			if err != nil {
				return nil, errors.Wrap(err, "create handler")
			}
		}
		if filter == nil {
			return nil, errors.New("no event filter")
		}
		if handler == nil {
			return nil, errors.New("no handler")
		}
		bot.Handle(filter, handler)
	}
	return bot, nil
}

func (b *Bot) Handle(filter types.EventFilter, h types.Handler) {
	b.handlers = append(b.handlers, &eventHandler{EventFilter: filter, Handler: h})
}

func (b *Bot) Start() error {
	updCfg := telegram.NewUpdate(0)
	updCfg.Timeout = 30
	updCh := b.botAPI.GetUpdatesChan(updCfg)
	go func() {
		for {
			select {
			case <-b.doneCh:
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
	b.botAPI.StopReceivingUpdates()
	close(b.doneCh)
	log.Print("Bot stopped")
	return nil
}

func (b *Bot) handleUpdate(upd *telegram.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, h := range b.handlers {
		if !h.Check(ctx, upd) {
			continue
		}
		if err := h.Handle(ctx, upd, b.botAPI); err != nil {
			log.Printf("Handler error: %v\n", err)
		}
	}
}
