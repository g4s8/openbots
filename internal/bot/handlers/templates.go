package handlers

import (
	"bytes"
	"text/template"

	"github.com/g4s8/openbots/internal/bot/interpolator"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type templateContext struct {
	Update  *telegram.Update
	State   map[string]string
	Secrets map[string]string
	Data    any
}

func newTemplateContext(upd *telegram.Update, state types.State, secrets map[string]types.Secret, Data any) *templateContext {
	stateMap := state.Map()
	secretMap := make(map[string]string, len(secrets))
	for k, v := range secrets {
		secretMap[k] = string(v)
	}
	return &templateContext{
		Update:  upd,
		State:   stateMap,
		Secrets: secretMap,
		Data:    Data,
	}
}

type Templater func(string) (Template, error)

func NewDefaultTemplate(src string) (Template, error) {
	return &defaultTemplate{src: src}, nil
}

func NewGoTemplate(src string) (Template, error) {
	tpl, err := template.New("go").Parse(src)
	if err != nil {
		return nil, errors.Wrap(err, "parse template")
	}
	return &goTemplate{tpl: tpl}, nil
}

func NewNoTemplate(src string) (Template, error) {
	return &noTemplate{src: src}, nil
}

type Template interface {
	Format(ctx *templateContext) (string, error)
}

type defaultTemplate struct {
	src string
}

func (t *defaultTemplate) Format(ctx *templateContext) (string, error) {
	secrets := make(map[string]types.Secret, len(ctx.Secrets))
	for k, v := range ctx.Secrets {
		secrets[k] = types.Secret(v)
	}
	intp := interpolator.New(ctx.State, secrets, ctx.Update)
	processed := intp.Interpolate(t.src)
	return processed, nil
}

type goTemplate struct {
	tpl *template.Template
}

func (t *goTemplate) Format(ctx *templateContext) (string, error) {
	var buf bytes.Buffer
	if err := t.tpl.Execute(&buf, ctx); err != nil {
		return "", errors.Wrap(err, "execute template")
	}
	return buf.String(), nil
}

type noTemplate struct {
	src string
}

func (t *noTemplate) Format(ctx *templateContext) (string, error) {
	return t.src, nil
}
