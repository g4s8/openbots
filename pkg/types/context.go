package types

import (
	ctx "context"
)

type ContextProvider interface {
	UserContext(ChatID) Context
}

type Context interface {

	// Set context value.
	Set(ctx.Context, string) error

	// Reset context.
	Reset(ctx.Context) error

	// Check context value.
	Check(ctx.Context, string) (bool, error)
}
