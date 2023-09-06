package api

import (
	"context"
	"errors"

	"github.com/g4s8/openbots/pkg/api"
	"github.com/g4s8/openbots/pkg/types"
)

type SendMessage struct {
	chats   []uint64
	handler api.Handler
}

func NewSendMessage(chats []uint64, handler api.Handler) *SendMessage {
	return &SendMessage{
		chats:   chats,
		handler: handler,
	}
}

func (sm *SendMessage) Call(ctx context.Context, req api.Request) error {
	if req.ChatID != 0 {
		return sm.handler.Call(ctx, req)
	}
	var errs []error
	for _, chat := range sm.chats {
		req.ChatID = types.ChatID(chat)
		if err := sm.handler.Call(ctx, req); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
