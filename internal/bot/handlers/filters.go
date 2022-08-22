package handlers

import (
	"context"
	"errors"

	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/g4s8/openbots-go/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	_ types.EventFilter = (*MessageFilter)(nil)
	_ types.EventFilter = (*CallbackFilter)(nil)
)

// MessageFilter checks update by message criteria.
type MessageFilter struct {
	check messageCriteria
}

func NewMessageFilterFromSpec(s *spec.MessageTrigger) (types.EventFilter, error) {
	if s.Command != "" {
		return &MessageFilter{
			check: messageHasCommand(s.Command),
		}, nil
	}
	if len(s.Text) > 0 {
		return &MessageFilter{
			check: messageHasText(s.Text),
		}, nil
	}
	return nil, errors.New("unknown trigger")
}

func (h *MessageFilter) Check(ctx context.Context, update *telegram.Update) bool {
	return update.Message != nil && h.check(update.Message)
}

type messageCriteria func(*telegram.Message) bool

func messageHasCommand(cmd string) messageCriteria {
	return func(msg *telegram.Message) bool {
		return msg.Command() == cmd
	}
}

func messageHasText(texts []string) messageCriteria {
	return func(msg *telegram.Message) bool {
		for _, text := range texts {
			if msg.Text == text {
				return true
			}
		}
		return false
	}
}

// CallbackFilter check update callback data.
type CallbackFilter struct {
	callback string
}

func (h *CallbackFilter) Check(ctx context.Context, update *telegram.Update) bool {
	return update.CallbackQuery != nil && update.CallbackQuery.Data == h.callback
}

func NewCallbackFilterFromSpec(s *spec.CallbackTrigger) (types.EventFilter, error) {
	return &CallbackFilter{
		callback: s.Data,
	}, nil
}
