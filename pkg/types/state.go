package types

import "context"

type State interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
	Map() map[string]string
	Fill(map[string]string)
	Changes() StateChanges
}

type StateChanges struct {
	Added   []string
	Removed []string
}

func (s *StateChanges) IsEmpty() bool {
	return len(s.Added) == 0 && len(s.Removed) == 0
}

type StateProvider interface {
	Load(context.Context, ChatID, State) error
	Update(context.Context, ChatID, State) error
}

type StateOp interface {
	Apply(State, ...func(string) string) error
}
