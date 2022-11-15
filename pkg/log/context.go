package log

import (
	ctx "context"

	"github.com/g4s8/openbots/pkg/types"
	"github.com/rs/zerolog"
)

type ContextProvider struct {
	base types.ContextProvider
	log  zerolog.Logger
}

func WrapContextProvider(base types.ContextProvider, log zerolog.Logger) *ContextProvider {
	return &ContextProvider{
		base: base,
		log:  log.With().Str("component", "context-provider").Logger(),
	}
}

func (c *ContextProvider) UserContext(chatID types.ChatID) types.Context {
	return WrapContext(c.base.UserContext(chatID), c.log.With().Str("chat", chatID.String()).Logger())
}

type Context struct {
	base types.Context
	log  zerolog.Logger
}

func WrapContext(base types.Context, log zerolog.Logger) *Context {
	return &Context{
		base: base,
		log:  log.With().Str("component", "context").Logger(),
	}
}

func (c *Context) Set(ctx ctx.Context, val string) error {
	c.log.Debug().Str("val", val).Msg("Set context")
	return c.base.Set(ctx, val)
}

func (c *Context) Reset(ctx ctx.Context) error {
	c.log.Debug().Msg("Reset context")
	return c.base.Reset(ctx)
}

func (c *Context) Check(ctx ctx.Context, val string) (bool, error) {
	res, err := c.base.Check(ctx, val)
	c.log.Debug().Str("val", val).Bool("check", res).Msg("Check context")
	return res, err
}
