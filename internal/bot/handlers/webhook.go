package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

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
	data    map[string]string
	sp      types.StateProvider
	secrets types.Secrets
	log     zerolog.Logger

	cli *http.Client
}

func NewWebhook(url *url.URL, method string, data map[string]string, sp types.StateProvider,
	secrets types.Secrets, log zerolog.Logger,
) *Webhook {
	return &Webhook{
		url:     url,
		method:  method,
		data:    data,
		sp:      sp,
		secrets: secrets,
		cli:     http.DefaultClient,
		log:     log.With().Str("handler", "webhook").Logger(),
	}
}

type WebhookPayload struct {
	Data map[string]string `json:"data"`
	Meta struct {
		ChatID    int64     `json:"chat_id"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"meta"`
}

func (h *Webhook) Handle(ctx context.Context, upd *telegram.Update, _ *telegram.BotAPI) error {
	state := state.NewUserState()
	err := h.sp.Load(ctx, ChatID(upd), state)
	if err != nil {
		return errors.Wrap(err, "load state")
	}
	secretMap, err := h.secrets.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "get secrets")
	}
	interpolator := newInterpolator(state.Map(), secretMap, upd)
	values := make(map[string]string, len(h.data))
	for k, v := range h.data {
		values[k] = interpolator.interpolate(v)
	}
	payload := WebhookPayload{
		Data: values,
	}
	payload.Meta.ChatID = ChatID(upd).Int64()
	payload.Meta.Timestamp = time.Now().UTC()

	body, err := json.Marshal(&payload)
	if err != nil {
		return errors.Wrap(err, "marshal payload")
	}

	br := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, h.method, h.url.String(), br)
	if err != nil {
		return errors.Wrap(err, "make HTTP request")
	}
	req.Header.Add("Content-Type", "application/json")
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
