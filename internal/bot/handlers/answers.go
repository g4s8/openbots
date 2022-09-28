package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

var _ types.Handler = (*MessageReply)(nil)

// MessageReplier func reply to message in chat.
type MessageReplier func(ctx context.Context, chatID int64, bot *telegram.BotAPI) error

// MessageReply handler processes telegram updates and reply message to them.
type MessageReply struct {
	replier MessageReplier
}

// NewMessageReply from repliers funcs.
func NewMessageReply(replier MessageReplier) *MessageReply {
	return &MessageReply{replier: replier}
}

func (h *MessageReply) Handle(ctx context.Context, update *telegram.Update,
	bot *telegram.BotAPI) error {
	if err := h.replier(ctx, ChatID(update), bot); err != nil {
		return errors.Wrap(err, "reply message")
	}
	return nil
}

type CallbackReply struct {
	text  string
	alert bool
}

func NewCallbackReply(text string, alert bool) *CallbackReply {
	return &CallbackReply{text: text, alert: alert}
}

func (h *CallbackReply) Handle(ctx context.Context, update *telegram.Update,
	bot *telegram.BotAPI) error {
	if update.CallbackQuery == nil {
		return errors.New("update is not callback query")
	}

	resp := telegram.NewCallback(update.CallbackQuery.ID, h.text)
	resp.ShowAlert = h.alert
	_, err := bot.Send(resp)
	if err != nil {
		return errors.Wrap(err, "send callback response")
	}
	return nil
}
