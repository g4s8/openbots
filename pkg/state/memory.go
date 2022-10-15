package state

import (
	"context"
	"sync"

	"github.com/g4s8/openbots/pkg/types"
)

var _ types.StateProvider = (*Memory)(nil)

type Memory struct {
	global map[string]string
	users  map[types.ChatID]map[string]string

	mux sync.RWMutex
}

func NewMemory(global map[string]string) *Memory {
	res := &Memory{
		global: make(map[string]string, len(global)),
		users:  make(map[types.ChatID]map[string]string),
	}
	for k, v := range global {
		res.global[k] = v
	}
	return res
}

// Get user state
func (m *Memory) Load(_ context.Context, user types.ChatID, state types.State) error {
	m.mux.RLock()
	defer m.mux.RUnlock()

	if user, ok := m.users[user]; ok {
		state.Fill(user)
	} else {
		state.Fill(m.global)
	}
	return nil
}

func (m *Memory) Update(_ context.Context, user types.ChatID, state types.State) error {
	var (
		isRep bool
		rep   changeReport
	)
	if r, ok := state.(reporter); ok {
		isRep = true
		rep = r.changes()
	}
	if isRep && len(rep.added) == 0 && len(rep.deleted) == 0 {
		return nil
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	var (
		data map[string]string
		ok   bool
	)
	cpy := state.Map()
	if data, ok = m.users[user]; !ok {
		m.users[user] = make(map[string]string, len(cpy))
		for k, v := range cpy {
			m.users[user][k] = v
		}
		return nil
	}

	if isRep {
		for _, key := range rep.added {
			data[key] = cpy[key]
		}
		for _, key := range rep.deleted {
			delete(data, key)
		}
	} else {
		m.users[user] = cpy
	}
	return nil
}
