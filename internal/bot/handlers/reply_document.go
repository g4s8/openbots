package handlers

import (
	"bufio"
	"context"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ types.Handler = (*ReplyDocument)(nil)

// ReplyDocument sends document file to chat.
type ReplyDocument struct {
	key    string
	name   string
	assets types.Assets
	logger zerolog.Logger
}

// NewReplyDocument creates new ReplyDocument handler using file by key.
func NewReplyDocument(key, name string, assets types.Assets,
	logger zerolog.Logger,
) *ReplyDocument {
	return &ReplyDocument{
		key:    key,
		name:   name,
		assets: assets,
		logger: logger,
	}
}

func (h *ReplyDocument) Handle(ctx context.Context, upd *telegram.Update,
	api *telegram.BotAPI,
) error {
	asset, err := h.assets.LoadAsset(ctx, h.key)
	if err != nil {
		return errors.Wrap(err, "load asset")
	}
	defer func() {
		if err := asset.Close(); err != nil {
			h.logger.Error().Err(err).Msg("Failed to close asset file")
		}
	}()
	fr := telegram.FileReader{
		Name:   h.name,
		Reader: bufio.NewReader(asset),
	}
	chatID := ChatID(upd)
	msg := telegram.NewDocument(int64(chatID), fr)
	if _, err := api.Send(msg); err != nil {
		return errors.Wrap(err, "send document")
	}
	return nil
}
