package handlers

import (
	"context"

	"github.com/g4s8/openbots/internal/bot/interpolator"
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
	for _, op := range h.ops {
		interpolator := interpolator.New(state.Map(), secretMap, update)
		if err := op.Apply(state, interpolator.Interpolate); err != nil {
			return errors.Wrap(err, "apply state op")
		}
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

func (o *logOp) Apply(target types.State, modifiers ...func(string) string) error {
	err := o.op.Apply(target, modifiers...)
	log := o.log.Debug()
	if err != nil {
		log = log.Err(err)
	}
	log.Msg("apply")
	return err
}

func NewStateHandlerFromSpec(provider types.StateProvider, s *spec.State, log zerolog.Logger) *StateHandler {
	ops := make([]types.StateOp, 0, len(s.Set)+len(s.Delete))
	for setK, setV := range s.Set {
		ops = append(ops, opWithLog(state.SetOp(setK, setV),
			log.With().
				Str("op", "set").
				Str("key", setK).
				Str("value", setV).
				Logger()))
	}
	for _, key := range s.Delete {
		ops = append(ops, opWithLog(state.DelOp(key),
			log.With().
				Str("op", "del").
				Str("key", key).
				Logger()))
	}
	for _, op := range s.Ops {
		switch op.Kind {
		case spec.StateUpdateOpKindSet:
			ops = append(ops, opWithLog(state.SetOp(op.Key, op.Value),
				log.With().
					Str("op", "set").
					Str("key", op.Key).
					Str("value", op.Value).
					Logger()))
		case spec.StateUpdateOpKindDelete:
			ops = append(ops, opWithLog(state.DelOp(op.Key),
				log.With().
					Str("op", "del").
					Str("key", op.Key).
					Logger()))

		case spec.StateUpdateOpKindAdd:
			ops = append(ops, opWithLog(state.AddOp(op.Key, op.Value),
				log.With().
					Str("op", "add").
					Str("key", op.Key).
					Str("value", op.Value).
					Logger()))
		case spec.StateUpdateOpKindSub:
			ops = append(ops, opWithLog(state.SubOp(op.Key, op.Value),
				log.With().
					Str("op", "sub").
					Str("key", op.Key).
					Str("value", op.Value).
					Logger()))
		case spec.StateUpdateOpKindMul:
			ops = append(ops, opWithLog(state.MulOp(op.Key, op.Value),
				log.With().
					Str("op", "mul").
					Str("key", op.Key).
					Str("value", op.Value).
					Logger()))
		case spec.StateUpdateOpKindDiv:
			ops = append(ops, opWithLog(state.DivOp(op.Key, op.Value),
				log.With().
					Str("op", "div").
					Str("key", op.Key).
					Str("value", op.Value).
					Logger()))

		}
	}
	return &StateHandler{
		provider: provider,
		secrets:  secrets.Stub,
		ops:      ops,
		log:      log.With().Str("handler", "state").Logger(),
	}
}
