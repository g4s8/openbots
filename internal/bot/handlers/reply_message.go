package handlers

import (
	"github.com/g4s8/openbots-go/pkg/spec"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	messageModifier func(*telegram.MessageConfig)
	inlineButton    struct {
		Text     string
		URL      string
		Callback string
	}
)

func messageWithKeyboard(keyboard [][]string) messageModifier {
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

func messageWithInlinceKeyboard(keyboard [][]inlineButton) messageModifier {
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

func newMessageReplier(text string, modifiers ...messageModifier) Replier {
	return func(chatID int64, bot *telegram.BotAPI) error {
		msg := telegram.NewMessage(chatID, text)
		for _, modifier := range modifiers {
			modifier(&msg)
		}
		_, err := bot.Send(msg)
		return err
	}
}

func inlineButtonsFromSpec(bts [][]spec.InlineButton) (res [][]inlineButton) {
	res = make([][]inlineButton, len(bts))
	for i, row := range bts {
		res[i] = make([]inlineButton, len(row))
		for j, btn := range row {
			res[i][j].Text = btn.Text
			res[i][j].URL = btn.URL
			res[i][j].Callback = btn.Callback
		}
	}
	return
}
