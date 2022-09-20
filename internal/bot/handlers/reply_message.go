package handlers

import (
	"context"

	"github.com/g4s8/openbots-go/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	MessageModifier func(*telegram.MessageConfig)
	InlineButton    struct {
		Text     string
		URL      string
		Callback string
	}
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

func MessageWithInlinceKeyboard(keyboard [][]InlineButton) MessageModifier {
	return func(msg *telegram.MessageConfig) {
		if len(keyboard) == 0 {
			return
		}
		buttons := make([][]telegram.InlineKeyboardButton, len(keyboard))
		for i, row := range keyboard {
			buttonRow := make([]telegram.InlineKeyboardButton, len(row))
			for j, btn := range row {
				buttonRow[j].Text = btn.Text
				if btn.URL != "" {
					setStr(&buttonRow[j].URL, btn.URL)
				} else if btn.Callback != "" {
					setStr(&buttonRow[j].CallbackData, btn.Callback)
				}
			}
			buttons[i] = buttonRow
		}
		msg.ReplyMarkup = telegram.NewInlineKeyboardMarkup(buttons...)
	}
}

func MessageWithParseMode(mode string) MessageModifier {
	return func(msg *telegram.MessageConfig) {
		msg.ParseMode = mode
	}
}

func NewMessageReplier(text string, modifiers ...MessageModifier) MessageReplier {
	return func(ctx context.Context, chatID int64, bot *telegram.BotAPI) error {
		state := types.StateFromContext(ctx, chatID)
		intp := newInterpolator(state)
		processed := intp.interpolate(text)

		msg := telegram.NewMessage(chatID, processed)
		for _, modifier := range modifiers {
			modifier(&msg)
		}
		_, err := bot.Send(msg)
		return err
	}
}
