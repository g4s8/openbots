package state

import "github.com/g4s8/openbots/pkg/types"

type StateSet struct {
	key, val string
}

func SetOp(key, val string) StateSet {
	return StateSet{key, val}
}

func (op StateSet) Apply(state types.State) {
	state.Set(op.key, op.val)
}

type StateDel struct {
	key string
}

func DelOp(key string) StateDel {
	return StateDel{key}
}

func (op StateDel) Apply(state types.State) {
	state.Delete(op.key)
}
