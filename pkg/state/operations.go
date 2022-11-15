package state

import "github.com/g4s8/openbots/pkg/types"

type StateSet struct {
	key, val string
}

func SetOp(key, val string) StateSet {
	return StateSet{key, val}
}

func (op StateSet) Apply(state types.State, modifiers ...func(string) string) {
	val := op.val
	for _, mod := range modifiers {
		val = mod(val)
	}

	state.Set(op.key, val)
}

type StateDel struct {
	key string
}

func DelOp(key string) StateDel {
	return StateDel{key}
}

func (op StateDel) Apply(state types.State, _ ...func(string) string) {
	state.Delete(op.key)
}
