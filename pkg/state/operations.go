package state

import (
	"fmt"
	"strconv"

	"github.com/g4s8/openbots/pkg/types"
)

type StateSet struct {
	key, val string
}

func SetOp(key, val string) StateSet {
	return StateSet{key, val}
}

func (op StateSet) Apply(state types.State, modifiers ...func(string) string) error {
	val := op.val
	for _, mod := range modifiers {
		val = mod(val)
	}

	state.Set(op.key, val)
	return nil
}

type StateDel struct {
	key string
}

func DelOp(key string) StateDel {
	return StateDel{key}
}

func (op StateDel) Apply(state types.State, _ ...func(string) string) error {
	state.Delete(op.key)
	return nil
}

type StateAdd struct {
	key, val string
}

func AddOp(key, val string) StateAdd {
	return StateAdd{key, val}
}

func (op StateAdd) Apply(state types.State, modifiers ...func(string) string) error {
	return modifyIntStateWithVal(state, op.key, op.val, func(x, y int) int { return x + y }, modifiers...)
}

type StateSub struct {
	key, val string
}

func SubOp(key, val string) StateSub {
	return StateSub{key, val}
}

func (op StateSub) Apply(state types.State, modifiers ...func(string) string) error {
	return modifyIntStateWithVal(state, op.key, op.val, func(x, y int) int { return x - y }, modifiers...)
}

type StateMul struct {
	key, val string
}

func MulOp(key, val string) StateMul {
	return StateMul{key, val}
}

func (op StateMul) Apply(state types.State, modifiers ...func(string) string) error {
	return modifyIntStateWithVal(state, op.key, op.val, func(x, y int) int { return x * y }, modifiers...)
}

type StateDiv struct {
	key, val string
}

func DivOp(key, val string) StateDiv {
	return StateDiv{key, val}
}

func (op StateDiv) Apply(state types.State, modifiers ...func(string) string) error {
	return modifyIntStateWithVal(state, op.key, op.val, func(x, y int) int { return x / y }, modifiers...)
}

func modifyIntStateWithVal(state types.State, key, val string, op func(x, y int) int, modifiers ...func(string) string) error {
	for _, mod := range modifiers {
		val = mod(val)
	}

	var x int
	current, ok := state.Get(key)
	if ok {
		v, err := strconv.Atoi(current)
		if err != nil {
			return fmt.Errorf("failed to parse value `x` %q as int: %w", current, err)
		}
		x = v
	}

	y, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("failed to parse value `y` %q as int: %w", val, err)
	}

	res := op(x, y)
	state.Set(key, strconv.Itoa(res))
	return nil
}
