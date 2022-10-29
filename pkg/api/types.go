package api

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
)

type Request struct {
	ChatID  types.ChatID
	Payload map[string]string
}

type Handler interface {
	Call(ctx context.Context, req Request) error
}
