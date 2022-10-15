package context

import (
	ctx "context"
	"sync"

	"github.com/g4s8/openbots/pkg/types"
)

type memoryProvider struct {
	values map[types.ChatID]string
	mux    sync.RWMutex
}

func NewMemoryProvider() types.ContextProvider {
	return &memoryProvider{
		values: make(map[types.ChatID]string),
	}
}

func (mp *memoryProvider) UserContext(uid types.ChatID) types.Context {
	return &memoryContext{
		provider: mp,
		uid:      uid,
	}
}

func (mp *memoryProvider) set(uid types.ChatID, value string) {
	mp.mux.Lock()
	defer mp.mux.Unlock()
	mp.values[uid] = value
}

func (mp *memoryProvider) get(uid types.ChatID) string {
	mp.mux.RLock()
	defer mp.mux.RUnlock()
	return mp.values[uid]
}

func (mp *memoryProvider) reset(uid types.ChatID) {
	mp.mux.Lock()
	defer mp.mux.Unlock()
	delete(mp.values, uid)
}

type memoryContext struct {
	provider *memoryProvider
	uid      types.ChatID
}

func (c *memoryContext) Set(_ ctx.Context, value string) error {
	c.provider.set(c.uid, value)
	return nil
}

func (c *memoryContext) Reset(_ ctx.Context) error {
	c.provider.reset(c.uid)
	return nil
}

func (c *memoryContext) Check(_ ctx.Context, val string) (bool, error) {
	return c.provider.get(c.uid) == val, nil
}
