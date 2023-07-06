package ctx

import (
	"context"
	"sync"

	"github.com/g4s8/openbots/pkg/types"
)

var _ types.Context = &Current{}

type pendingAction int

const (
	actionNone pendingAction = iota
	actionSet
	actionReset
)

// Current is a current context for a chat.
// It's immutable and changes are not applied antil `Save` is called.
type Current struct {
	current types.Context
	act     pendingAction
	newVal  string

	mx sync.Mutex
}

func (c *Current) Load(cp types.ContextProvider, chatID types.ChatID) {
	val := cp.UserContext(chatID)
	c.mx.Lock()
	c.current = val
	c.mx.Unlock()
}

func (c *Current) Save(ctx context.Context) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	switch c.act {
	case actionSet:
		return c.current.Set(ctx, c.newVal)
	case actionReset:
		return c.current.Reset(ctx)
	}
	return nil
}

func (c *Current) Set(ctx context.Context, val string) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.newVal = val
	c.act = actionSet

	return nil
}

func (c *Current) Reset(ctx context.Context) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.act = actionReset

	return nil
}

func (c *Current) Check(ctx context.Context, val string) (bool, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	return c.current.Check(ctx, val)
}
