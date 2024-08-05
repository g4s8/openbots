package interpolator

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var emptyMap = make(map[string]string)

type Interpolator struct {
	state   map[string]string
	secrets map[string]types.Secret
	upd     *telegram.Update
	data    map[string]string
}

type InterpolatorOp func(*Interpolator)

func WithState(state map[string]string) InterpolatorOp {
	return func(i *Interpolator) {
		i.state = state
	}
}

func WithSecrets(secrets map[string]types.Secret) InterpolatorOp {
	return func(i *Interpolator) {
		i.secrets = secrets
	}
}

func WithUpdate(upd *telegram.Update) InterpolatorOp {
	return func(i *Interpolator) {
		i.upd = upd
	}
}

func WithData(data map[string]string) InterpolatorOp {
	return func(i *Interpolator) {
		i.data = data
	}
}

// NewWithOps interpolator with options.
func NewWithOps(ops ...InterpolatorOp) *Interpolator {
	i := &Interpolator{}
	for _, op := range ops {
		op(i)
	}
	return i
}

// New interpolator with state and secrets.
// deprecated: use NewWithOps
func New(state map[string]string, secrets map[string]types.Secret, upd *telegram.Update) *Interpolator {
	return NewWithOps(
		WithState(state),
		WithSecrets(secrets),
		WithUpdate(upd),
	)
}

func (i *Interpolator) expander() func(string) string {
	data := make(map[string]string)
	if upd := i.upd; upd != nil {
		if msg := upd.Message; msg != nil {
			data["message.id"] = strconv.Itoa(msg.MessageID)
			data["message.text"] = msg.Text
			data["message.from.id"] = strconv.FormatInt(msg.From.ID, 10)
		}
		if chat := upd.FromChat(); chat != nil {
			data["chat.id"] = strconv.FormatInt(chat.ID, 10)
			data["chat.type"] = chat.Type
			data["chat.title"] = chat.Title
			data["chat.first_name"] = chat.FirstName
			data["chat.last_name"] = chat.LastName
			data["chat.username"] = chat.UserName
		}
		if user := upd.SentFrom(); user != nil {
			data["user.id"] = strconv.FormatInt(user.ID, 10)
			data["user.is_bot"] = strconv.FormatBool(user.IsBot)
			data["user.first_name"] = user.FirstName
			data["user.last_name"] = user.LastName
			data["user.username"] = user.UserName
			data["user.language_code"] = user.LanguageCode
		}
	}

	for k, v := range i.data {
		data["data."+k] = v
	}

	return func(text string) string {
		if strings.HasPrefix(text, "state.") {
			return i.state[text[6:]]
		}
		if strings.HasPrefix(text, "secret.") {
			secret, ok := i.secrets[text[7:]]
			if ok {
				return secret.Value()
			}
		}
		if val, ok := data[text]; ok {
			return val
		}
		return ""
	}
}

func (i *Interpolator) Interpolate(text string) string {
	res := os.Expand(text, i.expander())
	fmt.Printf("Interpolated: %q -> %q\n", text, res)
	return res
}

func (i *Interpolator) InterpolateS(s fmt.Stringer) string {
	return i.Interpolate(s.String())
}
