package handlers

import (
	"bufio"
	"context"
	"os"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ types.Handler = (*ReplyImage)(nil)

// ReplyImage sends image to chat.
type ReplyImage struct {
	file   string
	name   string
	logger zerolog.Logger
}

// NewReplyImageFile creates new ReplyImage handler using image from file path.
func NewReplyImageFile(file, name string, logger zerolog.Logger) *ReplyImage {
	return &ReplyImage{
		file:   file,
		name:   name,
		logger: logger,
	}
}

func (h *ReplyImage) Handle(ctx context.Context, upd *telegram.Update,
	api *telegram.BotAPI) error {
	f, err := os.Open(h.file)
	if err != nil {
		return errors.Wrap(err, "open image file")
	}
	buff := bufio.NewReader(f)
	defer func() {
		if err := f.Close(); err != nil {
			h.logger.Error().Err(err).Msg("Close image file")
		}
	}()
	fr := telegram.FileReader{
		Name:   h.name,
		Reader: buff,
	}
	chatID := ChatID(upd)
	msg := telegram.NewPhoto(int64(chatID), fr)
	if _, err := api.Send(msg); err != nil {
		return errors.Wrap(err, "send photo")
	}
	return nil
}
