package handlers

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InlineButton struct {
	Text     string
	URL      string
	Callback string
}

type InlineKeyboard [][]InlineButton

func (k InlineKeyboard) telegramMarkup(ip Interpolator) telegram.InlineKeyboardMarkup {
	buttons := make([][]telegram.InlineKeyboardButton, len(k))
	for i, row := range k {
		buttonRow := make([]telegram.InlineKeyboardButton, len(row))
		for j, btn := range row {
			buttonRow[j].Text = ip.Interpolate(btn.Text)
			if btn.URL != "" {
				setStr(&buttonRow[j].URL, ip.Interpolate(btn.URL))
			} else if btn.Callback != "" {
				setStr(&buttonRow[j].CallbackData, btn.Callback)
			}
		}
		buttons[i] = buttonRow
	}
	return telegram.NewInlineKeyboardMarkup(buttons...)
}
