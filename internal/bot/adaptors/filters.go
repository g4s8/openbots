package adaptors

import (
	"github.com/g4s8/openbots/internal/bot/filters"
	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/types"
	"github.com/rs/zerolog"
)

func NewPrecheckoutFilter(s *spec.PreCheckoutTrigger) *filters.PreCheckout {
	return filters.NewPreCheckout(s.InvoicePayload)
}

func NewPostcheckoutFilter(s *spec.PostCheckoutTrigger) *filters.PostCheckout {
	return filters.NewPostCheckout(s.InvoicePayload)
}

func NewStateFilter(sp types.StateProvider, logger zerolog.Logger, s []spec.StateCondition) filters.FilterChain {
	logger = logger.With().Str("component", "envet_filter").Str("filter", "state_filter").Logger()
	chain := make(filters.FilterChain, len(s))
	for i, c := range s {
		if c.Present.Valid {
			chain[i] = filters.NewStateFilterWithPresent(sp, logger, c.Key, c.Present.Value)
		} else if c.Eq != "" {
			chain[i] = filters.NewStateFilterEq(sp, logger, c.Key, c.Eq)
		} else if c.NEq != "" {
			chain[i] = filters.NewStateFilterNeq(sp, logger, c.Key, c.NEq)
		}
	}
	return chain
}
