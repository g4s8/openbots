package adaptors

import (
	"github.com/g4s8/openbots/internal/bot/filters"
	"github.com/g4s8/openbots/pkg/spec"
)

func NewPrecheckoutFilter(s *spec.PreCheckoutTrigger) *filters.PreCheckout {
	return filters.NewPreCheckout(s.InvoicePayload)
}

func NewPostcheckoutFilter(s *spec.PostCheckoutTrigger) *filters.PostCheckout {
	return filters.NewPostCheckout(s.InvoicePayload)
}
