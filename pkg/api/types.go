package api

import (
	"context"

	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/types"
)

type Request struct {
	ChatID  types.ChatID
	Payload map[string]string
}

type Handler interface {
	Call(ctx context.Context, req Request) error
}

type ReplyRequest struct {
	ChatID types.ChatID       `json:"chat_id"`
	Spec   *spec.MessageReply `json:"spec"`
}
