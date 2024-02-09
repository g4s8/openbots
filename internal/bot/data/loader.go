package data

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/g4s8/openbots/internal/bot/interpolator"
	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// LoaderConfig specifies data loader configuration options.
type LoaderConfig struct {
	Method  string
	URL     string
	Headers map[string]string
}

// Loader fetches data from the specified URL and stores it in the container.
type Loader struct {
	cli     *http.Client
	cfg     LoaderConfig
	sp      types.StateProvider
	secrets types.Secrets
	logger  zerolog.Logger
}

func NewLoader(cli *http.Client, cfg LoaderConfig, sp types.StateProvider, secrets types.Secrets, logger zerolog.Logger) *Loader {
	return &Loader{
		cli:     cli,
		cfg:     cfg,
		sp:      sp,
		secrets: secrets,
		logger:  logger,
	}
}

func (l *Loader) Load(ctx context.Context, c *types.DataContainer, upd *telegram.Update) error {
	state := state.NewUserState()
	defer state.Close()

	chat := upd.FromChat()
	var chatID int64
	if chat == nil {
		chatID = -1
	} else {
		chatID = chat.ID
	}

	if err := l.sp.Load(ctx, types.ChatID(chatID), state); err != nil {
		return errors.Wrap(err, "load state")
	}
	secretMap, err := l.secrets.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "get secrets")
	}
	ip := interpolator.New(state.Map(), secretMap, upd)
	req, err := http.NewRequestWithContext(ctx, l.cfg.Method, ip.Interpolate(l.cfg.URL), nil)
	if err != nil {
		return errors.Wrap(err, "create new request")
	}
	req.Header.Set("Accept", "application/json")
	for k, v := range l.cfg.Headers {
		req.Header.Set(k, ip.Interpolate(v))
	}
	resp, err := l.cli.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}

	l.logger.Info().
		Str("url", l.cfg.URL).
		Str("method", l.cfg.Method).
		Str("response_status", resp.Status).
		Msg("Data loaded")

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("bad status: %d", resp.StatusCode)
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.Wrap(err, "decode response")
	}
	c.Set(data)

	return nil
}
