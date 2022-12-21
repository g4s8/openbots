package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ types.Handler = (*SendInvoice)(nil)

type InvoicePrice struct {
	Label  string
	Amount int
}

type InvoiceConfig struct {
	Title       string
	Description string
	Payload     string
	Currency    string
	Prices      []InvoicePrice
}

type SendInvoice struct {
	providers    types.PaymentProviders
	providerName string
	config       InvoiceConfig
	logger       zerolog.Logger
}

func NewSendInvoice(providers types.PaymentProviders, providerName string, config InvoiceConfig,
	logger zerolog.Logger) *SendInvoice {
	return &SendInvoice{
		providers:    providers,
		providerName: providerName,
		config:       config,
		logger:       logger.With().Str("handler", "send_invoice").Logger(),
	}
}

func (h *SendInvoice) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	chatID := ChatID(upd)
	prices := make([]telegram.LabeledPrice, len(h.config.Prices))
	for i, p := range h.config.Prices {
		prices[i] = telegram.LabeledPrice{
			Label:  p.Label,
			Amount: p.Amount,
		}
	}
	h.logger.Debug().
		Int64("chat_id", chatID.Int64()).
		Str("provider", h.providerName).
		Str("payload", h.config.Payload).
		Msg("sending invoice")
	token := h.providers.PaymentToken(h.providerName)
	msg := telegram.NewInvoice(chatID.Int64(),
		h.config.Title, h.config.Description, h.config.Payload, token, "", h.config.Currency, prices)
	msg.MaxTipAmount = 10000
	msg.SuggestedTipAmounts = []int{100, 500, 1000, 5000}
	if _, err := api.Request(msg); err != nil {
		return errors.WithMessage(err, "send invoice")
	}
	return nil
}
