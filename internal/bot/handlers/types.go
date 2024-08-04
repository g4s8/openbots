package handlers

import (
	"context"

	"github.com/g4s8/openbots/internal/bot/data"
	"github.com/g4s8/openbots/internal/bot/interpolator"
	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Interpolator interface {
	Interpolate(string) string
}
type UpdateContextProvider struct {
	secrets types.Secrets
	state   types.StateProvider
}

type UpdateContext struct {
	upd     *telegram.Update
	state   map[string]string
	secrets map[string]types.Secret
	data    any
}

func (c *UpdateContext) ChatID() types.ChatID {
	return types.ChatID(rawChatID(c.upd))
}

func (c *UpdateContext) MessageID() int {
	return c.upd.CallbackQuery.Message.MessageID
}

func (c *UpdateContext) templateContext() *templateContext {
	return newTemplateContext(c.upd, c.state, c.secrets, c.data)
}

func (c *UpdateContext) Interpolator() Interpolator {
	opts := []interpolator.InterpolatorOp{
		interpolator.WithState(c.state),
		interpolator.WithSecrets(c.secrets),
		interpolator.WithUpdate(c.upd),
	}
	if dataMap, ok := c.data.(map[string]string); ok {
		opts = append(opts, interpolator.WithData(dataMap))
	}
	return interpolator.NewWithOps(opts...)
}

type updateContextKey struct{}

func UpdateContextFromCtx(ctx context.Context) *UpdateContext {
	if uctx, ok := ctx.Value(updateContextKey{}).(*UpdateContext); ok {
		return uctx
	}
	return &UpdateContext{}
}

func NewUpdateContextProvider(secrets types.Secrets, state types.StateProvider) *UpdateContextProvider {
	return &UpdateContextProvider{
		secrets: secrets,
		state:   state,
	}
}

func (cp *UpdateContextProvider) NewContext(ctx context.Context, upd *telegram.Update) (context.Context, error) {
	state := state.NewUserState()
	defer state.Close()

	chatID := ChatID(upd)
	if err := cp.state.Load(context.Background(), chatID, state); err != nil {
		return nil, errors.Wrap(err, "load state")
	}
	secretMap, err := cp.secrets.Get(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "get secrets")
	}
	dp := data.FromCtx(ctx)
	c := &UpdateContext{
		upd:     upd,
		state:   state.Map(),
		secrets: secretMap,
		data:    dp.Get(),
	}
	return context.WithValue(ctx, updateContextKey{}, c), nil
}
