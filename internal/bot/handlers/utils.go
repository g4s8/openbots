package handlers

import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func setStr(dst **string, src string) {
	*dst = new(string)
	**dst = src
}

func chatID(upd *telegram.Update) int64 {
	if upd.Message != nil {
		return upd.Message.Chat.ID
	}
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Message.Chat.ID
	}
	return -1
}
