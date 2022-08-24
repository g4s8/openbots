package handlers

import (
	"context"

	"github.com/g4s8/openbots-go/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SetContextHandler struct {
	base    types.Handler
	context *types.Context
	value   string
}

func NewContextSetter(base types.Handler, context *types.Context, value string) *SetContextHandler {
	return &SetContextHandler{
		base:    base,
		context: context,
		value:   value,
	}
}

func (h *SetContextHandler) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	if err := h.base.Handle(ctx, upd, api); err != nil {
		return err
	}
	h.context.Set(h.value)
	return nil
}

type DeleteContextHandler struct {
	base    types.Handler
	context *types.Context
	val     string
}

func NewContextDeleter(base types.Handler, context *types.Context, val string) *DeleteContextHandler {
	return &DeleteContextHandler{
		base:    base,
		context: context,
		val:     val,
	}
}

func (h *DeleteContextHandler) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	if err := h.base.Handle(ctx, upd, api); err != nil {
		return err
	}
	h.context.Delete(h.val)
	return nil
}
