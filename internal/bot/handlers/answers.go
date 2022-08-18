package handlers

import (
	"context"

	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/g4s8/openbots-go/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/multierr"
)

var _ types.Handler = (*Reply)(nil)

type Replier func(chatID int64, bot *telegram.BotAPI) error

type MessageModifier func(*telegram.MessageConfig)

func MessageWithKeyboard(keyboard [][]string) MessageModifier {
	return func(msg *telegram.MessageConfig) {
		if len(keyboard) > 0 {
			buttons := make([][]telegram.KeyboardButton, len(keyboard))
			for i, row := range keyboard {
				buttonRow := make([]telegram.KeyboardButton, len(row))
				for j, btn := range row {
					buttonRow[j] = telegram.NewKeyboardButton(btn)
				}
				buttons[i] = buttonRow
			}
			msg.ReplyMarkup = telegram.NewReplyKeyboard(buttons...)
		}
	}
}

func NewMessageReplier(text string, modifiers ...MessageModifier) Replier {
	return func(chatID int64, bot *telegram.BotAPI) error {
		msg := telegram.NewMessage(chatID, text)
		for _, modifier := range modifiers {
			modifier(&msg)
		}
		_, err := bot.Send(msg)
		return err
	}
}

type Reply struct {
	repliers []Replier
}

func NewReply(repliers ...Replier) *Reply {
	return &Reply{repliers: repliers}
}

func (h *Reply) Handle(ctx context.Context, update *telegram.Update, bot *telegram.BotAPI) error {
	chatID := update.Message.Chat.ID
	var err error
	for _, replier := range h.repliers {
		err = multierr.Append(err, replier(chatID, bot))
	}
	return err
}

func NewReplyFromSpec(s []*spec.Reply) (*Reply, error) {
	var repliers []Replier
	for _, r := range s {
		var modifiers []MessageModifier
		if r.Message.Markup != nil && len(r.Message.Markup.Keyboard) > 0 {
			modifiers = append(modifiers, MessageWithKeyboard(r.Message.Markup.Keyboard))
		}
		repliers = append(repliers, NewMessageReplier(r.Message.Text, modifiers...))
	}
	return NewReply(repliers...), nil
}
