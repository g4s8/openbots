package filters

import (
	"context"

	"github.com/g4s8/openbots/internal/bot/handlers"
	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ types.EventFilter = (*StateFilter)(nil)

type StateFilter struct {
	sp     types.StateProvider
	logger zerolog.Logger

	key     string
	eq      string
	neq     string
	present *bool
}

func NewStateFilterWithPresent(sp types.StateProvider, logger zerolog.Logger, key string, present bool) *StateFilter {
	return &StateFilter{
		sp:      sp,
		logger:  logger,
		key:     key,
		present: &present,
	}
}

func NewStateFilterEq(sp types.StateProvider, logger zerolog.Logger, key, eq string) *StateFilter {
	return &StateFilter{
		sp:     sp,
		logger: logger,
		key:    key,
		eq:     eq,
	}
}

func NewStateFilterNeq(sp types.StateProvider, logger zerolog.Logger, key, neq string) *StateFilter {
	return &StateFilter{
		sp:     sp,
		logger: logger,
		key:    key,
		neq:    neq,
	}
}

func (f *StateFilter) Check(ctx context.Context, upd *telegram.Update) (bool, error) {
	uid := handlers.ChatID(upd)
	state := state.NewUserState()
	if err := f.sp.Load(ctx, uid, state); err != nil {
		return false, errors.Wrap(err, "load state")
	}
	val, ok := state.Get(f.key)
	if f.present != nil {
		return ok == *f.present, nil
	}
	if f.eq != "" {
		return ok && val == f.eq, nil
	}
	if f.neq != "" {
		return !ok || val != f.neq, nil
	}
	return false, errors.New("invalid state filter")
}
