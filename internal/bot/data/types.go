package data

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
)

type ctxKey struct{}

func ContextWithContainer(ctx context.Context, c *types.DataContainer) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

func FromCtx(ctx context.Context) *types.DataContainer {
	val := ctx.Value(ctxKey{})
	if val == nil {
		return &types.DataContainer{}
	}
	return val.(*types.DataContainer)
}
