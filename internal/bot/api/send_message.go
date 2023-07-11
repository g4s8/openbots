package api

import (
	"context"

	"github.com/g4s8/openbots/pkg/api"
	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/types"
)

type SendMessage struct {
	chatID  spec.OptUint64
	handler api.Handler
}

func NewSendMessage(chatID spec.OptUint64, handler api.Handler) *SendMessage {
	return &SendMessage{
		chatID:  chatID,
		handler: handler,
	}
}

func (sm *SendMessage) Call(ctx context.Context, req api.Request) error {
	if req.ChatID == 0 && sm.chatID.Valid {
		req.ChatID = types.ChatID(sm.chatID.Value)
	}
	return sm.handler.Call(ctx, req)
}
