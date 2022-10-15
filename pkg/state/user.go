package state

import "github.com/g4s8/openbots/pkg/types"

var _ types.State = (*UserState)(nil)

type UserState struct {
	data map[string]string
	chng map[string]stateChange
}

func NewUserState() *UserState {
	res := &UserState{
		data: make(map[string]string),
		chng: make(map[string]stateChange),
	}
	return res
}

func (s *UserState) Get(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *UserState) Set(key, value string) {
	s.data[key] = value
	s.chng[key] = stateChangeWrite
}

func (s *UserState) Delete(key string) {
	delete(s.data, key)
	s.chng[key] = stateChangeDelete
}

func (s *UserState) Map() (out map[string]string) {
	out = make(map[string]string, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return
}

func (s *UserState) Fill(data map[string]string) {
	for k, v := range data {
		s.data[k] = v
	}
}

func (s *UserState) changes() (out changeReport) {
	for k, v := range s.chng {
		switch v {
		case stateChangeWrite:
			out.added = append(out.added, k)
		case stateChangeDelete:
			out.deleted = append(out.deleted, k)
		}
	}
	return
}

func (s *UserState) reset() {
	// for tests
	s.data = make(map[string]string)
	s.chng = make(map[string]stateChange)
}
