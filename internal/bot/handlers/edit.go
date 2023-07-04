package handlers

import (
	"context"
	"fmt"

	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ types.Handler = (*MessageEdit)(nil)

type editMessageMode int

const (
	editMessageCaptionMode      editMessageMode = 0
	editMessageTextMode         editMessageMode = 1
	editMessageTextKeyboardMode editMessageMode = 2
)

var ErrNoCallbackMessage = errors.New("callback data doesn't have message id")

type MessageEdit struct {
	caption  string
	template Template
	keyboard InlineKeyboard
	sp       types.StateProvider
	secrets  types.Secrets
	logger   zerolog.Logger
}

func NewMessageEdit(caption string, template Template, keyboard InlineKeyboard,
	sp types.StateProvider, secrets types.Secrets, logger zerolog.Logger,
) *MessageEdit {
	return &MessageEdit{
		caption:  caption,
		template: template,
		keyboard: keyboard,
		sp:       sp,
		secrets:  secrets,
		logger:   logger.With().Str("handler", "message_edit").Logger(),
	}
}

func (h *MessageEdit) mode(text string) editMessageMode {
	if h.caption != "" {
		return editMessageCaptionMode
	}
	var res editMessageMode
	if text != "" {
		res |= editMessageTextMode
	}
	if len(h.keyboard) > 0 {
		res |= editMessageTextKeyboardMode
	}
	return res
}

func (h *MessageEdit) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	if upd.CallbackQuery == nil || upd.CallbackQuery.Message == nil {
		return ErrNoCallbackMessage
	}

	msgID := upd.CallbackQuery.Message.MessageID
	chatID := ChatID(upd)
	state := state.NewUserState()
	defer state.Close()
	if err := h.sp.Load(ctx, chatID, state); err != nil {
		return errors.Wrap(err, "load user state")
	}
	secretMap, err := h.secrets.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "get secrets")
	}
	templateCtx := newTemplateContext(upd, state, secretMap, nil)
	text, err := h.template.Format(templateCtx)
	if err != nil {
		return errors.Wrap(err, "format template")
	}

	h.logger.Debug().
		Int("message_id", msgID).
		Int64("chat_id", int64(chatID)).
		Str("text", text).
		Str("origin_text", text).
		Msg("Edit message")

	var msg telegram.Chattable
	switch mode := h.mode(text); mode {
	case editMessageCaptionMode:
		msg = telegram.NewEditMessageCaption(int64(chatID), msgID, h.caption)
	case editMessageTextMode:
		msg = telegram.NewEditMessageText(int64(chatID), msgID, text)
	case editMessageTextKeyboardMode:
		msg = telegram.NewEditMessageReplyMarkup(int64(chatID), msgID, h.keyboard.telegramMarkup())
	case editMessageTextKeyboardMode | editMessageTextMode:
		msg = telegram.NewEditMessageTextAndMarkup(int64(chatID), msgID, text, h.keyboard.telegramMarkup())
	default:
		return fmt.Errorf("unsupported edit message mode: %d", mode)
	}

	if _, err := api.Send(msg); err != nil {
		return errors.Wrap(err, "send message")
	}

	return nil
}
