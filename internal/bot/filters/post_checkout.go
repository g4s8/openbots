package filters

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ types.EventFilter = (*PreCheckout)(nil)

type PostCheckout struct {
	invoicePayload string
}

func NewPostCheckout(invoicePayload string) *PostCheckout {
	return &PostCheckout{
		invoicePayload: invoicePayload,
	}
}

func (f *PostCheckout) Check(ctx context.Context, upd *telegram.Update) (bool, error) {
	if upd.Message == nil || upd.Message.SuccessfulPayment == nil {
		return false, nil
	}
	return upd.Message.SuccessfulPayment.InvoicePayload == f.invoicePayload, nil
}
