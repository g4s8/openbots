package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type StateHandler struct {
	provider types.StateProvider
	ops      []types.StateOp
}

func (h *StateHandler) Handle(ctx context.Context, update *telegram.Update) error {
	uid := ChatID(update)
	state := state.NewUserState()
	if err := h.provider.Load(ctx, uid, state); err != nil {
		return errors.Wrap(err, "load state")
	}
	for _, op := range h.ops {
		op.Apply(state)
	}
	if err := h.provider.Update(ctx, uid, state); err != nil {
		return errors.Wrap(err, "update state")
	}
	return nil
}

func NewStateHandlerFromSpec(provider types.StateProvider, spec *spec.State) *StateHandler {
	ops := make([]types.StateOp, 0, len(spec.Set)+len(spec.Delete))
	for setK, setV := range spec.Set {
		ops = append(ops, state.SetOp(setK, setV))
	}
	for _, key := range spec.Delete {
		ops = append(ops, state.DelOp(key))
	}
	return &StateHandler{provider: provider, ops: ops}
}
