package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type (
	MessageModifier func(*telegram.MessageConfig)
)

func MessageWithKeyboard(keyboard [][]string) MessageModifier {
	return func(msg *telegram.MessageConfig) {
		if len(keyboard) == 0 {
			return
		}
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

func MessageWithInlinceKeyboard(keyboard InlineKeyboard) MessageModifier {
	return func(msg *telegram.MessageConfig) {
		if len(keyboard) == 0 {
			return
		}
		msg.ReplyMarkup = keyboard.telegramMarkup()
	}
}

func MessageWithParseMode(mode string) MessageModifier {
	return func(msg *telegram.MessageConfig) {
		msg.ParseMode = mode
	}
}

func NewMessageReplier(sp types.StateProvider, text string, logger zerolog.Logger, modifiers ...MessageModifier) MessageReplier {
	return func(ctx context.Context, upd *telegram.Update, bot *telegram.BotAPI) error {
		logger = logger.With().Str("handler", "reply_message").Logger()

		state := state.NewUserState()
		defer state.Close()

		chatID := ChatID(upd)
		if err := sp.Load(ctx, chatID, state); err != nil {
			return errors.Wrap(err, "load state")
		}

		intp := newInterpolator(state, upd)
		processed := intp.interpolate(text)

		msg := telegram.NewMessage(int64(chatID), processed)
		for _, modifier := range modifiers {
			modifier(&msg)
		}
		_, err := bot.Send(msg)
		return err
	}
}
