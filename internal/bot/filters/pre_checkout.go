package filters

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ types.EventFilter = (*PreCheckout)(nil)

type PreCheckout struct {
	invoicePayload string
}

func NewPreCheckout(invoicePayload string) *PreCheckout {
	return &PreCheckout{
		invoicePayload: invoicePayload,
	}
}

func (f *PreCheckout) Check(ctx context.Context, upd *telegram.Update) (bool, error) {
	if upd.PreCheckoutQuery == nil {
		return false, nil
	}
	return upd.PreCheckoutQuery.InvoicePayload == f.invoicePayload, nil
}
