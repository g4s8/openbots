package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gobwas/glob"
	"github.com/pkg/errors"
)

var (
	_ types.EventFilter = (*MessageFilter)(nil)
	_ types.EventFilter = (*CallbackFilter)(nil)
	_ types.EventFilter = (*ContextFilter)(nil)
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

func (h *MessageFilter) Check(ctx context.Context, update *telegram.Update) (bool, error) {
	return update.Message != nil && h.check(update.Message), nil
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
	callback glob.Glob
}

func (h *CallbackFilter) Check(ctx context.Context, update *telegram.Update) (bool, error) {
	return update.CallbackQuery != nil && h.callback.Match(update.CallbackQuery.Data), nil
}

func NewCallbackFilterFromSpec(s *spec.CallbackTrigger) (types.EventFilter, error) {
	g, err := glob.Compile(s.Data)
	if err != nil {
		return nil, errors.Wrap(err, "compile glob")
	}
	return &CallbackFilter{
		callback: g,
	}, nil
}

type ContextFilter struct {
	base types.EventFilter
	cp   types.ContextProvider
	val  string
}

func NewContextFilter(base types.EventFilter, cp types.ContextProvider, val string) types.EventFilter {
	return &ContextFilter{
		base: base,
		cp:   cp,
		val:  val,
	}
}

func (h *ContextFilter) Check(ctx context.Context, update *telegram.Update) (bool, error) {
	ctxCheck, err := h.cp.UserContext(ChatID(update)).Check(ctx, h.val)
	if err != nil {
		return false, errors.Wrap(err, "check context")
	}
	if !ctxCheck {
		return false, nil
	}
	if h.base != nil {
		check, err := h.base.Check(ctx, update)
		if err != nil {
			return false, err
		}
		if !check {
			return false, nil
		}
	}

	return true, nil
}
