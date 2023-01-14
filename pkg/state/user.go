package state

import (
	"sort"
	"sync"

	"github.com/g4s8/openbots/pkg/types"
)

var _ types.State = (*UserState)(nil)

var statePool = sync.Pool{
	New: func() any {
		return makeUserState()
	},
}

type UserState struct {
	data map[string]string
	chng map[string]stateChange
}

func NewUserState() *UserState {
	return statePool.Get().(*UserState)
}

func makeUserState() *UserState {
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

func (s *UserState) Changes() (out types.StateChanges) {
	for k, v := range s.chng {
		switch v {
		case stateChangeWrite:
			out.Added = append(out.Added, k)
		case stateChangeDelete:
			out.Removed = append(out.Removed, k)
		}
	}
	sort.Strings(out.Added)
	sort.Strings(out.Removed)
	return
}

func (s *UserState) Close() {
	s.reset()
	statePool.Put(s)
}

func (s *UserState) reset() {
	for k := range s.data {
		delete(s.data, k)
	}
	for k := range s.chng {
		delete(s.chng, k)
	}
}
