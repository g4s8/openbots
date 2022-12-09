package handlers

import (
	"context"

	"github.com/g4s8/openbots/pkg/state"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type (
	// MessageModifier apply custom modifications to telegram message reply.
	MessageModifier func(*telegram.MessageConfig)
)

var _ types.Handler = (*MessageReply)(nil)

// MessageReply handler processes telegram updates and reply message to them.
type MessageReply struct {
	sp        types.StateProvider
	text      string
	modifiers []MessageModifier
	logger    zerolog.Logger
}

// NewMessageReply from repliers funcs.
func NewMessageReply(sp types.StateProvider, text string, logger zerolog.Logger, modifiers ...MessageModifier) *MessageReply {
	return &MessageReply{
		sp:        sp,
		text:      text,
		modifiers: modifiers,
		logger:    logger.With().Str("handler", "reply_message").Logger(),
	}
}

func (h *MessageReply) Handle(ctx context.Context, upd *telegram.Update,
	bot *telegram.BotAPI) error {
	state := state.NewUserState()
	defer state.Close()

	chatID := ChatID(upd)
	if err := h.sp.Load(ctx, chatID, state); err != nil {
		return errors.Wrap(err, "load state")
	}

	intp := newInterpolator(state, upd)
	processed := intp.interpolate(h.text)

	msg := telegram.NewMessage(int64(chatID), processed)
	for _, modifier := range h.modifiers {
		modifier(&msg)
	}
	if _, err := bot.Send(msg); err != nil {
		return errors.Wrap(err, "reply message")
	}
	return nil
}

// MessageWithKeyboard creates new message modifier to add
// custom keyboard to message.
func MessageWithKeyboard(keyboard [][]string) MessageModifier {
	return func(msg *telegram.MessageConfig) {
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
	return func(msg *telegram.MessageConfig) {
		if len(keyboard) == 0 {
			return
		}
		msg.ReplyMarkup = keyboard.telegramMarkup()
	}
}

// MessageWithParseMode creates new message modifier to set
// custom parse mode for message.
func MessageWithParseMode(mode string) MessageModifier {
	return func(msg *telegram.MessageConfig) {
		msg.ParseMode = mode
	}
}
