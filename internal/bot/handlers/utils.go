package handlers

import (
	"os"
	"strings"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setStr(dst **string, src string) {
	*dst = new(string)
	**dst = src
}

func ChatID(upd *telegram.Update) int64 {
	if upd.Message != nil {
		return upd.Message.Chat.ID
	}
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Message.Chat.ID
	}
	return -1
}

type interpolator struct {
	state types.UserState
}

func newInterpolator(state types.UserState) *interpolator {
	return &interpolator{state: state}
}

func (i *interpolator) expander() func(string) string {
	return func(text string) string {
		if strings.HasPrefix(text, "state.") {
			return i.state[text[6:]]
		}
		return ""
	}
}

func (i *interpolator) interpolate(text string) string {
	return os.Expand(text, i.expander())
}
