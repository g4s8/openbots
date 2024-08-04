package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/api"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type (
	// MessageModifier apply custom modifications to telegram message reply.
	MessageModifier func(context.Context, *telegram.MessageConfig)
)

var (
	_ types.Handler = (*MessageReply)(nil)
	_ api.Handler   = (*MessageReply)(nil)
)

// MessageReply handler processes telegram updates and reply message to them.
type MessageReply struct {
	bot       *telegram.BotAPI
	template  Template
	modifiers []MessageModifier
	logger    zerolog.Logger
}

// NewMessageReply from repliers funcs.
func NewMessageReply(
	bot *telegram.BotAPI,
	template Template, logger zerolog.Logger, modifiers ...MessageModifier,
) *MessageReply {
	return &MessageReply{
		bot:       bot,
		template:  template,
		modifiers: modifiers,
		logger:    logger.With().Str("handler", "reply_message").Logger(),
	}
}

func (h *MessageReply) Handle(ctx context.Context, upd *telegram.Update, _ *telegram.BotAPI) error {
	updCtx := UpdateContextFromCtx(ctx)

	chatID := updCtx.ChatID()

	response, err := h.template.Format(updCtx.templateContext())
	if err != nil {
		return errors.Wrap(err, "format template")
	}

	msg := telegram.NewMessage(int64(chatID), response)
	for _, modifier := range h.modifiers {
		modifier(ctx, &msg)
	}
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "reply message")
	}
	return nil
}

func (h *MessageReply) Call(ctx context.Context, req api.Request) error {
	updCtx := UpdateContextFromCtx(ctx)

	response, err := h.template.Format(updCtx.templateContext())
	if err != nil {
		return errors.Wrap(err, "format template")
	}

	msg := telegram.NewMessage(int64(req.ChatID), response)
	for _, modifier := range h.modifiers {
		modifier(ctx, &msg)
	}
	if _, err := h.bot.Send(msg); err != nil {
		return errors.Wrap(err, "send message")
	}
	return nil
}

// MessageWithKeyboard creates new message modifier to add
// custom keyboard to message.
func MessageWithKeyboard(keyboard [][]string) MessageModifier {
	return func(ctx context.Context, msg *telegram.MessageConfig) {
		if len(keyboard) == 0 {
			return
		}
		buttons := make([][]telegram.KeyboardButton, len(keyboard))
		for i, row := range keyboard {
			buttonRow := make([]telegram.KeyboardButton, len(row))
			for j, btn := range row {
				buttonRow[j] = telegram.NewKeyboardButton(btn)
			}
			buttons[i] = buttonRow
		}
		msg.ReplyMarkup = telegram.NewReplyKeyboard(buttons...)
	}
}

// MessageWithInlineKeyboard creates new message modifier to add
// custom inline keyboard to message.
func MessageWithInlineKeyboard(keyboard InlineKeyboard) MessageModifier {
	return func(ctx context.Context, msg *telegram.MessageConfig) {
		if len(keyboard) == 0 {
			return
		}
		u := UpdateContextFromCtx(ctx)
		msg.ReplyMarkup = keyboard.telegramMarkup(u.Interpolator())
	}
}

// MessageWithParseMode creates new message modifier to set
// custom parse mode for message.
func MessageWithParseMode(mode string) MessageModifier {
	return func(ctx context.Context, msg *telegram.MessageConfig) {
		msg.ParseMode = mode
	}
}
