package ctx

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
	"github.com/pkg/errors"
)

var _ types.ContextProvider = &Provider{}

// Provider - context provider for pending context state.
type Provider struct {
	origin   types.ContextProvider
	pendinng map[types.ChatID]*Current
}

func NewProvider(origin types.ContextProvider) *Provider {
	return &Provider{
		origin:   origin,
		pendinng: make(map[types.ChatID]*Current),
	}
}

// Closer - context closer. It saves pending context state on close.
type Closer func(context.Context) error

func newCloser(current *Current) Closer {
	return func(ctx context.Context) error {
		if err := current.Save(ctx); err != nil {
			return errors.Wrap(err, "save context")
		}
		return nil
	}
}

func (p *Provider) Begin(id types.ChatID) Closer {
	var pending Current
	pending.Load(p.origin, id)
	p.pendinng[id] = &pending
	return newCloser(&pending)
}

func (p *Provider) UserContext(id types.ChatID) types.Context {
	if ctx, ok := p.pendinng[id]; ok {
		return ctx
	}
	return nil
}
