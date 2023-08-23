package handlers

import (
	"context"

	"github.com/g4s8/openbots/internal/bot/interpolator"
	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

// CallbackReply send callback message reply.
type CallbackReply struct {
	sp      types.StateProvider
	secrets types.Secrets
	text    string
	alert   bool
}

// NewCallbackReply creates new callback reply handler.
func NewCallbackReply(sp types.StateProvider, secrets types.Secrets, text string, alert bool) *CallbackReply {
	return &CallbackReply{
		sp:      sp,
		secrets: secrets,
		text:    text,
		alert:   alert,
	}
}

var ErrUpdateNotSupported = errors.New("update is not valid for this handler")

func (h *CallbackReply) Handle(ctx context.Context, upd *telegram.Update,
	bot *telegram.BotAPI,
) error {
	if upd.CallbackQuery == nil {
		return errors.Wrap(ErrUpdateNotSupported, "not callback query")
	}

	chatID := ChatID(upd)
	state := state.NewUserState()
	defer state.Close()

	if err := h.sp.Load(ctx, chatID, state); err != nil {
		return errors.Wrap(err, "load state")
	}
	secretMap, err := h.secrets.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "get secrets")
	}
	interpolator := interpolator.New(state.Map(), secretMap, upd)
	text := interpolator.Interpolate(h.text)

	resp := telegram.NewCallback(upd.CallbackQuery.ID, text)
	resp.ShowAlert = h.alert
	if _, err := bot.Request(resp); err != nil {
		return errors.Wrap(err, "callback reply")
	}
	return nil
}
