package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/secrets"
	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type StateHandler struct {
	provider types.StateProvider
	secrets  types.Secrets
	ops      []types.StateOp
	log      zerolog.Logger
}

func (h *StateHandler) Handle(ctx context.Context, update *telegram.Update, _ *telegram.BotAPI) error {
	uid := ChatID(update)
	state := state.NewUserState()
	if err := h.provider.Load(ctx, uid, state); err != nil {
		return errors.Wrap(err, "load state")
	}
	secretMap, err := h.secrets.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "get secrets")
	}
	interpolator := newInterpolator(state, secretMap, update)
	for _, op := range h.ops {
		op.Apply(state, interpolator.interpolate)
	}
	if err := h.provider.Update(ctx, uid, state); err != nil {
		return errors.Wrap(err, "update state")
	}
	return nil
}

type logOp struct {
	op  types.StateOp
	log zerolog.Logger
}

func opWithLog(op types.StateOp, log zerolog.Logger) types.StateOp {
	return &logOp{op: op, log: log}
}

func (o *logOp) Apply(target types.State, modifiers ...func(string) string) {
	o.op.Apply(target, modifiers...)
	o.log.Debug().Msg("apply")
}

func NewStateHandlerFromSpec(provider types.StateProvider, spec *spec.State, log zerolog.Logger) *StateHandler {
	ops := make([]types.StateOp, 0, len(spec.Set)+len(spec.Delete))
	for setK, setV := range spec.Set {
		ops = append(ops, opWithLog(state.SetOp(setK, setV),
			log.With().
				Str("op", "set").
				Str("key", setK).
				Str("value", setV).
				Logger()))
	}
	for _, key := range spec.Delete {
		ops = append(ops, opWithLog(state.DelOp(key),
			log.With().
				Str("op", "del").
				Str("key", key).
				Logger()))
	}
	return &StateHandler{
		provider: provider,
		secrets:  secrets.Stub,
		ops:      ops,
		log:      log.With().Str("handler", "state").Logger(),
	}
}
