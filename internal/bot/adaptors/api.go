package adaptors

import (
	"github.com/g4s8/openbots/internal/bot/api"
	"github.com/g4s8/openbots/pkg/spec"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ApiSendMessage(tg *telegram.BotAPI, s *spec.ApiSendMesageAction) *api.SendMessage {
	return api.NewSendMessage(tg, apiArg(s.Text))
}

func apiArg(s *spec.ApiArg) api.Argument {
	if s.Param != "" {
		return api.NewRefArg(s.Param)
	}
	return api.NewConstArg(s.Value)
}
