package handlers

import (
	"context"

	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/g4s8/openbots-go/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/multierr"
)

var _ types.Handler = (*Reply)(nil)

type Replier func(chatID int64, bot *telegram.BotAPI) error

type MessageModifier func(*telegram.MessageConfig)

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

type InlineButton struct {
	Text     string
	URL      string
	Callback string
}

func MessageWithInlinceKeyboard(keyboard [][]InlineButton) MessageModifier {
	return func(msg *telegram.MessageConfig) {
		if len(keyboard) == 0 {
			return
		}
		buttons := make([][]telegram.InlineKeyboardButton, len(keyboard))
		for i, row := range keyboard {
			buttonRow := make([]telegram.InlineKeyboardButton, len(row))
			for j, btn := range row {
				buttonRow[j].Text = btn.Text
				if btn.URL != "" {
					setStr(&buttonRow[j].URL, btn.URL)
				} else if btn.Callback != "" {
					setStr(&buttonRow[j].CallbackData, btn.Callback)
				}
			}
			buttons[i] = buttonRow
		}
		msg.ReplyMarkup = telegram.NewInlineKeyboardMarkup(buttons...)
	}
}

func setStr(dst **string, src string) {
	*dst = new(string)
	**dst = src
}

func NewMessageReplier(text string, modifiers ...MessageModifier) Replier {
	return func(chatID int64, bot *telegram.BotAPI) error {
		msg := telegram.NewMessage(chatID, text)
		for _, modifier := range modifiers {
			modifier(&msg)
		}
		_, err := bot.Send(msg)
		return err
	}
}

type Reply struct {
	repliers []Replier
}

func NewReply(repliers ...Replier) *Reply {
	return &Reply{repliers: repliers}
}

func chatID(upd *telegram.Update) int64 {
	if upd.Message != nil {
		return upd.Message.Chat.ID
	}
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Message.Chat.ID
	}
	return -1
}

func (h *Reply) Handle(ctx context.Context, update *telegram.Update, bot *telegram.BotAPI) error {
	var err error
	for _, replier := range h.repliers {
		err = multierr.Append(err, replier(chatID(update), bot))
	}
	return err
}

func inlineButtonsFromSpec(bts [][]spec.InlineButton) (res [][]InlineButton) {
	res = make([][]InlineButton, len(bts))
	for i, row := range bts {
		res[i] = make([]InlineButton, len(row))
		for j, btn := range row {
			res[i][j].Text = btn.Text
			res[i][j].URL = btn.URL
			res[i][j].Callback = btn.Callback
		}
	}
	return
}

func NewReplyFromSpec(s []*spec.Reply) (*Reply, error) {
	var repliers []Replier
	for _, r := range s {
		var modifiers []MessageModifier
		if r.Message.Markup != nil && len(r.Message.Markup.Keyboard) > 0 {
			modifiers = append(modifiers, MessageWithKeyboard(r.Message.Markup.Keyboard))
		}
		if r.Message.Markup != nil && len(r.Message.Markup.InlineKeyboard) > 0 {
			modifiers = append(modifiers, MessageWithInlinceKeyboard(
				inlineButtonsFromSpec(r.Message.Markup.InlineKeyboard)))
		}
		repliers = append(repliers, NewMessageReplier(r.Message.Text, modifiers...))
	}
	return NewReply(repliers...), nil
}
