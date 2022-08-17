package handlers

import (
	"context"
	"errors"

	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/g4s8/openbots-go/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	_ types.EventFilter = (*OnMessage)(nil)
	_ types.Handler     = Nop
)

var Nop = &nop{}

type nop struct{}

func (nop) Handle(context.Context, *telegram.Update, *telegram.BotAPI) error {
	return nil
}

type OnMessageCheck func(*telegram.Message) bool

func NewMessageCommandCheck(cmd string) OnMessageCheck {
	return func(msg *telegram.Message) bool {
		return msg.Command() == cmd
	}
}

func NewMessageTextCheck(text string) OnMessageCheck {
	return func(msg *telegram.Message) bool {
		return msg.Text == text
	}
}

type OnMessage struct {
	check OnMessageCheck
}

func (h *OnMessage) Check(ctx context.Context, update *telegram.Update) bool {
	if m := update.Message; m != nil && h.check(m) {
		return true
	}
	return false
}

type AnyOfCheck struct {
	filters []types.EventFilter
}

func (h *AnyOfCheck) Check(ctx context.Context, update *telegram.Update) bool {
	for _, filter := range h.filters {
		if filter.Check(ctx, update) {
			return true
		}
	}
	return false
}

func NewOnMessageFromSpec(s *spec.MessageTrigger) (types.EventFilter, error) {
	if s.Command != "" {
		return &OnMessage{
			check: NewMessageCommandCheck(s.Command),
		}, nil
	}
	if len(s.Text) > 0 {
		var anyOf AnyOfCheck
		for _, text := range s.Text {
			anyOf.filters = append(anyOf.filters, &OnMessage{check: NewMessageTextCheck(text)})
		}
		return &anyOf, nil
	}
	return nil, errors.New("unknown trigger")
}
