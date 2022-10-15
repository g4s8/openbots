package types

import "context"

type State interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
	Map() map[string]string
	Fill(map[string]string)
}

type StateProvider interface {
	Load(context.Context, ChatID, State) error
	Update(context.Context, ChatID, State) error
}

type StateOp interface {
	Apply(State)
}
