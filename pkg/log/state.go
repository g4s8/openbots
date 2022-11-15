package log

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
	"github.com/rs/zerolog"
)

type StateProvider struct {
	base types.StateProvider
	log  zerolog.Logger
}

func WrapStateProvider(base types.StateProvider, log zerolog.Logger) *StateProvider {
	return &StateProvider{
		base: base,
		log:  log.With().Str("component", "state-provider").Logger(),
	}
}

func (s *StateProvider) Load(ctx context.Context, chatID types.ChatID, state types.State) error {
	s.log.Debug().Str("chat", chatID.String()).Msg("Load state")
	return s.base.Load(ctx, chatID, state)
}

func (s *StateProvider) Update(ctx context.Context, chatID types.ChatID, state types.State) error {
	s.log.Debug().Str("chat", chatID.String()).Msg("Update state")
	return s.base.Update(ctx, chatID, state)
}
