package types

import "context"

var EmptyUserState = make(UserState, 0)

type State struct {
	users  map[int64]UserState
	global map[string]string
}

func (s *State) User(id int64) UserState {
	if s.users == nil {
		s.users = make(map[int64]UserState)
	}
	if _, ok := s.users[id]; !ok {
		s.users[id] = make(UserState)
	}
	res := make(UserState)
	for k, v := range s.users[id] {
		res[k] = v
	}
	for k, v := range s.global {
		if _, ok := res[k]; !ok {
			res[k] = v
		}
	}
	return res
}

func (s *State) Save(id int64, state UserState) {
	s.users[id] = state
}

type UserState map[string]string

type UserStateOp func(UserState)

func NewState(global map[string]string) State {
	var state State
	state.global = make(map[string]string)
	if global != nil {
		for k, v := range global {
			state.global[k] = v
		}
	}
	state.users = make(map[int64]UserState)
	return state
}

func (s UserState) Apply(ops ...UserStateOp) UserState {
	cpy := make(UserState, len(s))
	for k, v := range s {
		cpy[k] = v
	}
	for _, op := range ops {
		op(cpy)
	}
	return cpy
}

func SetState(key, value string) UserStateOp {
	return func(s UserState) {
		s[key] = value
	}
}

func DeleteState(keys ...string) UserStateOp {
	return func(s UserState) {
		for _, key := range keys {
			delete(s, key)
		}
	}
}

type StateCtxKey struct{}

func StateFromContext(ctx context.Context, user int64) UserState {
	s, ok := ctx.Value(StateCtxKey{}).(State)
	if !ok {
		return EmptyUserState
	}
	us, ok := s.users[user]
	if !ok {
		return EmptyUserState
	}
	for k, v := range s.global {
		if _, ok := us[k]; !ok {
			us[k] = v
		}
	}
	return us
}

func ContextWithState(ctx context.Context, user int64, s State) context.Context {
	return context.WithValue(ctx, StateCtxKey{}, s)
}
