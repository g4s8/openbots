package handlers

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	"github.com/g4s8/openbots/pkg/types"
)

var _ types.Handler = (*Webhook)(nil)

type Webhook struct {
	url    *url.URL
	method string
	body   []byte

	cli *http.Client
}

func NewWebhook(url *url.URL, method string, body []byte) *Webhook {
	return &Webhook{
		url:    url,
		method: method,
		body:   body,
		cli:    http.DefaultClient,
	}
}

func (h *Webhook) Handle(ctx context.Context, _ *telegram.Update, _ *telegram.BotAPI) error {
	var br io.Reader
	if h.body != nil {
		br = bytes.NewReader(h.body)
	}
	req, err := http.NewRequestWithContext(ctx, h.method, h.url.String(), br)
	if err != nil {
		return errors.Wrap(err, "make HTTP request")
	}
	resp, err := h.cli.Do(req)
	if err != nil {
		return errors.Wrap(err, "call HTTP")
	}
	defer resp.Body.Close()
	var body string
	if resp.ContentLength > 0 {
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			return errors.Wrap(err, "read HTTP response body")
		}
		body = buf.String()
	}
	log.Printf("Call HTTP %s %s: %d - %s", h.method, h.url, resp.StatusCode, body)
	return nil
}
