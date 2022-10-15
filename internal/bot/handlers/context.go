package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type SetContextHandler struct {
	base  types.Handler
	cp    types.ContextProvider
	value string
}

func NewContextSetter(base types.Handler, cp types.ContextProvider, value string) *SetContextHandler {
	return &SetContextHandler{
		base:  base,
		cp:    cp,
		value: value,
	}
}

func (h *SetContextHandler) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	if err := h.base.Handle(ctx, upd, api); err != nil {
		return err
	}
	if err := h.cp.UserContext(ChatID(upd)).Set(ctx, h.value); err != nil {
		return errors.Wrap(err, "set context")
	}
	return nil
}

type DeleteContextHandler struct {
	base types.Handler
	cp   types.ContextProvider
	val  string
}

func NewContextDeleter(base types.Handler, cp types.ContextProvider, val string) *DeleteContextHandler {
	return &DeleteContextHandler{
		base: base,
		cp:   cp,
		val:  val,
	}
}

func (h *DeleteContextHandler) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	if err := h.base.Handle(ctx, upd, api); err != nil {
		return err
	}
	if err := h.cp.UserContext(ChatID(upd)).Reset(ctx); err != nil {
		return errors.Wrap(err, "delete context")
	}
	return nil
}
