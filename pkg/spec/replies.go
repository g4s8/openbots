package spec

import (
	"errors"
	"fmt"
)

type Reply struct {
	Message  *MessageReply  `yaml:"message"`
	Callback *CallbackReply `yaml:"callback"`
}

func (r *Reply) validate() (errs []error) {
	errs = make([]error, 0)
	if r.Message == nil && r.Callback == nil {
		errs = append(errs, errors.New("empty reply"))
	}
	if r.Message != nil {
		errs = append(errs, r.Message.validate()...)
	}
	if r.Callback != nil {
		errs = append(errs, r.Callback.validate()...)
	}
	return
}

type MessageReply struct {
	Text   string       `yaml:"text"`
	Markup *ReplyMarkup `yaml:"markup"`
}

func (r *MessageReply) validate() []error {
	if r.Text == "" {
		return []error{errors.New("empty message reply")}
	}
	if r.Markup != nil {
		return r.Markup.validate()
	}
	return []error{}
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
