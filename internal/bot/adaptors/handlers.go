package adaptors

import (
	"context"

	"github.com/g4s8/openbots/internal/bot/handlers"
	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.uber.org/multierr"
)

func MessageRepply(sp types.StateProvider, s *spec.MessageReply, log zerolog.Logger) *handlers.MessageReply {
	var modifiers []handlers.MessageModifier
	if s.Markup != nil && len(s.Markup.Keyboard) > 0 {
		modifiers = append(modifiers, handlers.MessageWithKeyboard(s.Markup.Keyboard))
	}
	if s.Markup != nil && len(s.Markup.InlineKeyboard) > 0 {
		modifiers = append(modifiers, handlers.MessageWithInlineKeyboard(
			inlineKeyboardFromSpec(s.Markup.InlineKeyboard)))
	}
	if s.ParseMode != "" {
		modifiers = append(modifiers, handlers.MessageWithParseMode(string(s.ParseMode)))
	}
	return handlers.NewMessageReply(sp, s.Text, log, modifiers...)
}

func CallbackReply(s *spec.CallbackReply) *handlers.CallbackReply {
	return handlers.NewCallbackReply(s.Text, s.Alert)
}

func inlineKeyboardFromSpec(bts [][]spec.InlineButton) (res handlers.InlineKeyboard) {
	res = make(handlers.InlineKeyboard, len(bts))
	for i, row := range bts {
		res[i] = make([]handlers.InlineButton, len(row))
		for j, btn := range row {
			res[i][j].Text = btn.Text
			res[i][j].URL = btn.URL
			res[i][j].Callback = btn.Callback
		}
	}
	return
}

type multiHandler struct {
	handlers []types.Handler
}

func (h *multiHandler) Handle(ctx context.Context, upd *telegram.Update, bot *telegram.BotAPI) error {
	var merr error
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, upd, bot); err != nil {
			merr = multierr.Append(merr, err)
		}
	}
	if merr != nil {
		return errors.Wrap(merr, "one or more handler error")
	}
	return nil
}

func Replies(sp types.StateProvider, assets types.Assets,
	r []*spec.Reply, log zerolog.Logger) types.Handler {
	var handlers []types.Handler
	for _, reply := range r {
		if reply.Message != nil {
			handlers = append(handlers,
				MessageRepply(sp, reply.Message, log))
		}
		if reply.Callback != nil {
			handlers = append(handlers, CallbackReply(reply.Callback))
		}
		if reply.Edit != nil {
			handlers = append(handlers, newEdit(reply.Edit, sp, log))
		}
		if reply.Delete {
			handlers = append(handlers, newDelete(log))
		}
		if reply.Image != nil {
			handlers = append(handlers, newImageReply(reply.Image, assets, log))
		}
	}
	return &multiHandler{handlers}
}

func Webhook(s *spec.Webhook, sp types.StateProvider, log zerolog.Logger) types.Handler {
	return handlers.NewWebhook(s.URL, s.Method, s.Body, sp, log)
}

func newEdit(s *spec.Edit, sp types.StateProvider, log zerolog.Logger) types.Handler {
	if s.Message == nil {
		log.Fatal().Msg("Invalid edit spec: message is empty")
	}
	msg := s.Message
	return handlers.NewMessageEdit(msg.Caption, msg.Text, inlineKeyboardFromSpec(msg.InlineKeyboard),
		sp, log)
}

func newDelete(logger zerolog.Logger) types.Handler {
	return handlers.NewMessageDelete(logger)
}

func newImageReply(s *spec.ImageReply, assets types.Assets,
	log zerolog.Logger) types.Handler {
	if s.Key != "" {
		return handlers.NewReplyImageFile(s.Key, s.Name, assets, log)
	}
	return nil
}
