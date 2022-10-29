package bot

import (
	"context"

	"github.com/g4s8/openbots/pkg/api"
)

type apiHandlerGroup struct {
	handlers []api.Handler
}

func (g *apiHandlerGroup) Call(ctx context.Context, req api.Request) error {
	for _, h := range g.handlers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := h.Call(ctx, req); err != nil {
			return err
		}
	}
	return nil
}
