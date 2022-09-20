package handlers

import (
	"context"

	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/g4s8/openbots-go/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateHandler struct {
	ops []types.UserStateOp
}

func (h *StateHandler) Handle(ctx context.Context, update *telegram.Update,
	state types.UserState) (types.UserState, error) {
	return state.Apply(h.ops...), nil
}

func NewStateHandlerFromSpec(spec *spec.State) *StateHandler {
	ops := make([]types.UserStateOp, 0, len(spec.Set)+len(spec.Delete))
	for setK, setV := range spec.Set {
		ops = append(ops, types.SetState(setK, setV))
	}
	for _, key := range spec.Delete {
		ops = append(ops, types.DeleteState(key))
	}
	return &StateHandler{ops: ops}
}
