package handlers

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type PreCheckout struct {
	ok     bool
	err    string
	logger zerolog.Logger
}

func NewPreCheckout(ok bool, err string, logger zerolog.Logger) *PreCheckout {
	return &PreCheckout{
		ok:     ok,
		err:    err,
		logger: logger.With().Str("handler", "pre_checkout").Logger(),
	}
}

var ErrPrecheckoutQueryEmpty = errors.New("precheckout query is empty")

func (h *PreCheckout) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	if upd.PreCheckoutQuery == nil {
		return ErrPrecheckoutQueryEmpty
	}

	logger := h.logger.With().Str("query", upd.PreCheckoutQuery.ID).Bool("ok", h.ok).Logger()
	if h.err != "" {
		logger = logger.With().Str("error", h.err).Logger()
	}
	logger.Debug().Msg("Sending precheckout query answer")

	msg := telegram.PreCheckoutConfig{
		PreCheckoutQueryID: upd.PreCheckoutQuery.ID,
		OK:                 h.ok,
		ErrorMessage:       h.err,
	}
	if _, err := api.Request(msg); err != nil {
		return errors.WithMessage(err, "send precheckout answer")
	}
	return nil
}
