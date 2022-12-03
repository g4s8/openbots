package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ types.Handler = (*MessageDelete)(nil)

type MessageDelete struct {
	logger zerolog.Logger
}

func NewMessageDelete(logger zerolog.Logger) *MessageDelete {
	return &MessageDelete{
		logger: logger,
	}
}

func (d *MessageDelete) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	if upd.CallbackQuery == nil || upd.CallbackQuery.Message == nil {
		return ErrNoCallbackMessage
	}

	msgID := upd.CallbackQuery.Message.MessageID
	chatID := ChatID(upd)
	msg := telegram.NewDeleteMessage(int64(chatID), msgID)

	d.logger.Debug().Int("message_id", msgID).Int64("chat_id", int64(chatID)).Msg("Deleting message")

	if _, err := api.Send(msg); err != nil {
		return errors.Wrap(err, "delete message")
	}
	return nil
}
