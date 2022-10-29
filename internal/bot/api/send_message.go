package api

import (
	"context"
	"log"

	"github.com/g4s8/openbots/pkg/api"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type SendMessage struct {
	tg   *telegram.BotAPI
	text Argument
}

func NewSendMessage(tg *telegram.BotAPI, text Argument) *SendMessage {
	return &SendMessage{
		tg:   tg,
		text: text,
	}
}

func (sm *SendMessage) Call(ctx context.Context, req api.Request) error {
	log.Printf("send message: req=%+v", req)
	text, err := sm.text.Get(req)
	if err != nil {
		return errors.Wrap(err, "get 'text' argument")
	}

	msg := telegram.NewMessage(int64(req.ChatID), text)
	if _, err := sm.tg.Send(msg); err != nil {
		return errors.Wrap(err, "send message")
	}

	return nil
}
