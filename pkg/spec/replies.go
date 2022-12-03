package spec

import (
	"errors"
	"fmt"
)

type Reply struct {
	Message  *MessageReply  `yaml:"message"`
	Callback *CallbackReply `yaml:"callback"`
	Edit     *Edit          `yaml:"edit"`
	Delete   bool           `yaml:"delete"`
}

func (r *Reply) validate() (errs []error) {
	errs = make([]error, 0)
	if r.Message == nil && r.Callback == nil && r.Edit == nil && !r.Delete {
		errs = append(errs, errors.New("empty reply"))
	}
	if r.Message != nil {
		errs = append(errs, r.Message.validate()...)
	}
	if r.Callback != nil {
		errs = append(errs, r.Callback.validate()...)
	}
	if r.Edit != nil {
		errs = append(errs, r.Edit.validate()...)
	}
	return
}

// ParseMode of message reply
type ParseMode string

func (pm ParseMode) validate() []error {
	switch pm {
	case ModeMarkdown, ModeMarkdownV2, ModeHTML:
		return nil
	}
	return []error{fmt.Errorf("invalid parse mode %s", pm)}
}

const (
	ModeMarkdown   = ParseMode("Markdown")
	ModeMarkdownV2 = ParseMode("MarkdownV2")
	ModeHTML       = ParseMode("HTML")
)

type MessageReply struct {
	Text      string       `yaml:"text"`
	ParseMode ParseMode    `yaml:"parseMode"`
	Markup    *ReplyMarkup `yaml:"markup"`
}

func (r *MessageReply) validate() []error {
	var errs []error
	if r.Text == "" {
		errs = append(errs, errors.New("empty message reply"))
	}
	if r.Markup != nil {
		errs = append(errs, r.Markup.validate()...)
	}
	if r.ParseMode != "" {
		errs = append(errs, ParseMode(r.ParseMode).validate()...)
	}
	return errs
}

type InlineButton struct {
	Text     string `yaml:"text"`
	URL      string `yaml:"url"`
	Callback string `yaml:"callback"`
}

type ReplyMarkup struct {
	Keyboard       [][]string       `yaml:"keyboard"`
	InlineKeyboard [][]InlineButton `yaml:"inlineKeyboard"`
}

func (r *ReplyMarkup) validate() []error {
	if len(r.Keyboard) == 0 && len(r.InlineKeyboard) == 0 {
		return []error{errors.New("empty reply markup")}
	}
	errs := make([]error, 0)
	for i, row := range r.Keyboard {
		if len(row) == 0 {
			errs = append(errs, fmt.Errorf("empty keyboard row %d", i))
		}
		for j, button := range row {
			if button == "" {
				errs = append(errs, fmt.Errorf("empty keyboard button %d:%d", i, j))
			}
		}
	}

	for i, row := range r.InlineKeyboard {
		if len(row) == 0 {
			errs = append(errs, fmt.Errorf("empty inline keyboard row %d", i))
		}
		for j, button := range row {
			if button.Text == "" {
				errs = append(errs, fmt.Errorf("empty inline keyboard button %d:%d", i, j))
			}
			if button.URL == "" && button.Callback == "" {
				errs = append(errs, fmt.Errorf("empty inline keyboard button action %d:%d", i, j))
			}
		}
	}
	return errs
}

type CallbackReply struct {
	Text  string `yaml:"text"`
	Alert bool   `yaml:"alert"`
}

func (r *CallbackReply) validate() []error {
	if r.Text == "" {
		return []error{errors.New("empty callback reply")}
	}
	return []error{}
}
