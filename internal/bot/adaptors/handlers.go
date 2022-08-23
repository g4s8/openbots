package adaptors

import (
	"context"

	"github.com/g4s8/openbots-go/internal/bot/handlers"
	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/g4s8/openbots-go/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

func MessageRepply(s *spec.MessageReply) *handlers.MessageReply {
	var modifiers []handlers.MessageModifier
	if s.Markup != nil && len(s.Markup.Keyboard) > 0 {
		modifiers = append(modifiers, handlers.MessageWithKeyboard(s.Markup.Keyboard))
	}
	if s.Markup != nil && len(s.Markup.InlineKeyboard) > 0 {
		modifiers = append(modifiers, handlers.MessageWithInlinceKeyboard(
			inlineButtonsFromSpec(s.Markup.InlineKeyboard)))
	}
	replier := handlers.NewMessageReplier(s.Text, modifiers...)
	return handlers.NewMessageReply(replier)
}

func CallbackReply(s *spec.CallbackReply) *handlers.CallbackReply {
	return handlers.NewCallbackReply(s.Text, s.Alert)
}

func inlineButtonsFromSpec(bts [][]spec.InlineButton) (res [][]handlers.InlineButton) {
	res = make([][]handlers.InlineButton, len(bts))
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

func Replies(r []*spec.Reply) types.Handler {
	var handlers []types.Handler
	for _, reply := range r {
		var handler types.Handler
		if reply.Message != nil {
			handler = MessageRepply(reply.Message)
		} else if reply.Callback != nil {
			handler = CallbackReply(reply.Callback)
		}
		if handler != nil {
			handlers = append(handlers, handler)
		}
	}
	return &multiHandler{handlers}
}
