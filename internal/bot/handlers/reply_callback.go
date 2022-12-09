package handlers

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

// CallbackReply send callback message reply.
type CallbackReply struct {
	text  string
	alert bool
}

// NewCallbackReply creates new callback reply handler.
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
	if _, err := bot.Request(resp); err != nil {
		return errors.Wrap(err, "callback reply")
	}
	return nil
}
