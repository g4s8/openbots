package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/api"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type SetContextHandler struct {
	cp    types.ContextProvider
	value string
	log   zerolog.Logger
}

func NewContextSetter(cp types.ContextProvider, value string, log zerolog.Logger) *SetContextHandler {
	return &SetContextHandler{
		cp:    cp,
		value: value,
		log:   log,
	}
}

func (h *SetContextHandler) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	if err := h.cp.UserContext(ChatID(upd)).Set(ctx, h.value); err != nil {
		return errors.Wrap(err, "set context")
	}
	return nil
}

func (h *SetContextHandler) Call(ctx context.Context, req api.Request) error {
	// TODO: refactor similar logic with Handle
	if err := h.cp.UserContext(req.ChatID).Set(ctx, h.value); err != nil {
		return errors.Wrap(err, "set context")
	}
	return nil
}

type DeleteContextHandler struct {
	cp  types.ContextProvider
	val string
	log zerolog.Logger
}

func NewContextDeleter(cp types.ContextProvider, val string, log zerolog.Logger) *DeleteContextHandler {
	return &DeleteContextHandler{
		cp:  cp,
		val: val,
		log: log,
	}
}

func (h *DeleteContextHandler) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	if err := h.cp.UserContext(ChatID(upd)).Reset(ctx); err != nil {
		return errors.Wrap(err, "delete context")
	}
	return nil
}

func (h *DeleteContextHandler) Call(ctx context.Context, req api.Request) error {
	// TODO: refactor similar logic with Handle
	if err := h.cp.UserContext(req.ChatID).Reset(ctx); err != nil {
		return errors.Wrap(err, "delete context")
	}
	return nil
}
