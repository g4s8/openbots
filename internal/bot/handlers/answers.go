package handlers

import (
	"context"

	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/g4s8/openbots-go/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/multierr"
)

var _ types.Handler = (*Reply)(nil)

// Replier func reply to message in chat.
type Replier func(ctx context.Context, chatID int64, bot *telegram.BotAPI) error

// Reply handler processes telegram updates and reply to them.
type Reply struct {
	repliers []Replier
}

// NewReplyFromSpec creates reply object from reply specification.
func NewReplyFromSpec(s []*spec.Reply) (*Reply, error) {
	var repliers []Replier
	for _, r := range s {
		var modifiers []messageModifier
		if r.Message.Markup != nil && len(r.Message.Markup.Keyboard) > 0 {
			modifiers = append(modifiers, messageWithKeyboard(r.Message.Markup.Keyboard))
		}
		if r.Message.Markup != nil && len(r.Message.Markup.InlineKeyboard) > 0 {
			modifiers = append(modifiers, messageWithInlinceKeyboard(
				inlineButtonsFromSpec(r.Message.Markup.InlineKeyboard)))
		}
		repliers = append(repliers, newMessageReplier(r.Message.Text, modifiers...))
	}
	return NewReply(repliers...), nil
}

// NewReply from repliers funcs.
func NewReply(repliers ...Replier) *Reply {
	return &Reply{repliers: repliers}
}

func (h *Reply) Handle(ctx context.Context, update *telegram.Update,
	bot *telegram.BotAPI) error {
	var err error
	for _, replier := range h.repliers {
		err = multierr.Append(err, replier(ctx, chatID(update), bot))
	}
	return err
}
