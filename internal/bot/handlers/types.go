package handlers

import (
	"context"
	"fmt"
	"strings"

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

type UpdateContext struct {
	upd     *telegram.Update
	state   map[string]string
	secrets map[string]types.Secret
	data    *types.DataContainer
}

func (c *UpdateContext) ChatID() types.ChatID {
	return types.ChatID(rawChatID(c.upd))
}

func (c *UpdateContext) MessageID() int {
	return c.upd.CallbackQuery.Message.MessageID
}

func (c *UpdateContext) templateContext() *templateContext {
	var data any
	if c.data != nil {
		data = c.data.Get()
	}
	return newTemplateContext(c.upd, c.state, c.secrets, data)
}

func (c *UpdateContext) Interpolator() Interpolator {
	opts := []interpolator.InterpolatorOp{
		interpolator.WithState(c.state),
		interpolator.WithSecrets(c.secrets),
		interpolator.WithUpdate(c.upd),
	}
	var data any
	if c.data != nil {
		data = c.data.Get()
	}
	if dataMap, ok := data.(map[string]any); ok {
		m := make(map[string]string, len(dataMap))
		for k, v := range dataMap {
			m[k] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, interpolator.WithData(m))
	}
	return interpolator.NewWithOps(opts...)
}

func (c *UpdateContext) String() string {
	var sb strings.Builder
	sb.WriteString("UpdateContext{")
	sb.WriteString(fmt.Sprintf("state: %v, ", c.state))
	sb.WriteString(fmt.Sprintf("secrets: %v, ", c.secrets))
	sb.WriteString(fmt.Sprintf("data: %v", c.data))
	sb.WriteString("}")
	return sb.String()
}

type updateContextKey struct{}

func UpdateContextFromCtx(ctx context.Context) *UpdateContext {
	if uctx, ok := ctx.Value(updateContextKey{}).(*UpdateContext); ok {
		res := uctx
		res.data = data.FromCtx(ctx)
		return res
	}
	return &UpdateContext{}
}

type UpdateContextProvider struct {
	secrets types.Secrets
	state   types.StateProvider
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
	if err := cp.state.Load(ctx, chatID, state); err != nil {
		return nil, errors.Wrap(err, "load state")
	}
	secretMap, err := cp.secrets.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get secrets")
	}
	c := &UpdateContext{
		upd:     upd,
		state:   state.Map(),
		secrets: secretMap,
	}
	return context.WithValue(ctx, updateContextKey{}, c), nil
}
