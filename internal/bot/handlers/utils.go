package handlers

import (
	"os"
	"strconv"
	"strings"

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
	if upd.Message != nil {
		return upd.Message.Chat.ID
	}
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Message.Chat.ID
	}
	return -1
}

type interpolator struct {
	state   types.State
	secrets map[string]types.Secret

	message *telegram.Message
}

func newInterpolator(state types.State, secrets map[string]types.Secret, upd *telegram.Update) *interpolator {
	res := &interpolator{
		state:   state,
		secrets: secrets,
	}
	if upd.Message != nil {
		res.message = upd.Message
	}
	return res
}

func (i *interpolator) expander() func(string) string {
	message := make(map[string]string)
	if i.message != nil {
		message["id"] = strconv.Itoa(i.message.MessageID)
		message["text"] = i.message.Text
		message["from.id"] = strconv.FormatInt(i.message.From.ID, 10)
	}
	return func(text string) string {
		state := i.state.Map()
		if strings.HasPrefix(text, "state.") {
			return state[text[6:]]
		}
		if strings.HasPrefix(text, "message.") {
			return message[text[8:]]
		}
		if strings.HasPrefix(text, "secret.") {
			secret, ok := i.secrets[text[7:]]
			if ok {
				return secret.Value()
			}
		}
		return ""
	}
}

func (i *interpolator) interpolate(text string) string {
	return os.Expand(text, i.expander())
}
