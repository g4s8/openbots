package types

import "context"

var EmptyState = make(State, 0)

type State map[string]string

type StateOp func(State)

func (s State) Apply(ops ...StateOp) State {
	cpy := make(State, len(s))
	for k, v := range s {
		cpy[k] = v
	}
	for _, op := range ops {
		op(cpy)
	}
	return cpy
}

func SetState(key, value string) StateOp {
	return func(s State) {
		s[key] = value
	}
}

func DeleteState(keys ...string) StateOp {
	return func(s State) {
		for _, key := range keys {
			delete(s, key)
		}
	}
}

type StateCtxKey struct{}

func StateFromContext(ctx context.Context) State {
	s, ok := ctx.Value(StateCtxKey{}).(State)
	if !ok {
		return EmptyState
	}
	return s
}

func ContextWithState(ctx context.Context, s State) context.Context {
	return context.WithValue(ctx, StateCtxKey{}, s)
}
