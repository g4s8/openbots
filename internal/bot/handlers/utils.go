package handlers

import (
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setStr(dst **string, src string) {
	*dst = new(string)
	**dst = src
}

func ChatID(upd *telegram.Update) types.ChatID {
	return types.ChatID(rawChatID(upd))
}

func rawChatID(upd *telegram.Update) int64 {
	if chat := upd.FromChat(); chat != nil {
		return chat.ID
	}
	return -1
}
