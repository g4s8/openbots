package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
)

var _ types.Handler = (*Webhook)(nil)

type Webhook struct {
	url     *url.URL
	method  string
	payload map[string]string
	sp      types.StateProvider
	log     zerolog.Logger

	cli *http.Client
}

func NewWebhook(url *url.URL, method string, payload map[string]string, sp types.StateProvider,
	log zerolog.Logger) *Webhook {
	return &Webhook{
		url:     url,
		method:  method,
		payload: payload,
		sp:      sp,
		cli:     http.DefaultClient,
		log:     log.With().Str("handler", "webhook").Logger(),
	}
}

func (h *Webhook) Handle(ctx context.Context, upd *telegram.Update, _ *telegram.BotAPI) error {
	state := state.NewUserState()
	err := h.sp.Load(ctx, ChatID(upd), state)
	if err != nil {
		return errors.Wrap(err, "load state")
	}
	interpolator := newInterpolator(state, upd)
	values := make(map[string]string, len(h.payload))
	for k, v := range h.payload {
		values[k] = interpolator.interpolate(v)
	}
	body, err := json.Marshal(values)
	if err != nil {
		return errors.Wrap(err, "marshal payload")
	}

	br := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, h.method, h.url.String(), br)
	if err != nil {
		return errors.Wrap(err, "make HTTP request")
	}
	resp, err := h.cli.Do(req)
	if err != nil {
		return errors.Wrap(err, "call HTTP")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.ContentLength > 0 {
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "read response body")
		}
	}
	h.log.Printf("Call HTTP %s %s: %d", h.method, h.url, resp.StatusCode)
	return nil
}
