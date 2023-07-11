package interpolator

import (
	"os"
	"strconv"
	"strings"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Interpolator struct {
	state   map[string]string
	secrets map[string]types.Secret

	message *telegram.Message
}

func New(state map[string]string, secrets map[string]types.Secret, upd *telegram.Update) *Interpolator {
	res := &Interpolator{
		state:   state,
		secrets: secrets,
	}
	if upd != nil && upd.Message != nil {
		res.message = upd.Message
	}
	return res
}

func (i *Interpolator) expander() func(string) string {
	message := make(map[string]string)
	if i.message != nil {
		message["id"] = strconv.Itoa(i.message.MessageID)
		message["text"] = i.message.Text
		message["from.id"] = strconv.FormatInt(i.message.From.ID, 10)
	}
	return func(text string) string {
		if strings.HasPrefix(text, "state.") {
			return i.state[text[6:]]
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

func (i *Interpolator) Interpolate(text string) string {
	return os.Expand(text, i.expander())
}
