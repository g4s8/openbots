package data

import (
	"context"

	"github.com/g4s8/openbots/pkg/api"
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

func FromApiRequest(req *api.Request) *types.DataContainer {
	var data types.DataContainer
	values := make(map[string]string)
	for k, v := range req.Payload {
		values[k] = v
	}
	values["chat.id"] = req.ChatID.String()

	data.Set(values)
	return &data
}
