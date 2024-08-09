package spec

import (
	"errors"
	"fmt"
)

type Reply struct {
	Message     *MessageReply      `yaml:"message"`
	Callback    *CallbackReply     `yaml:"callback"`
	Edit        *Edit              `yaml:"edit"`
	Delete      bool               `yaml:"delete"`
	Image       *FileReply         `yaml:"image"`
	Document    *FileReply         `yaml:"document"`
	Invoice     *Invoice           `yaml:"invoice"`
	PreCheckout *PreCheckoutAnswer `yaml:"preCheckout"`
}

// ParseMode of message reply
type ParseMode string

const (
	ModeMarkdown   = ParseMode("Markdown")
	ModeMarkdownV2 = ParseMode("MarkdownV2")
	ModeHTML       = ParseMode("HTML")
)

// TemplateStyle of message text.
type TemplateStyle string

func (ts TemplateStyle) validate() []error {
	switch ts {
	case TemplateDefault, TemplateGo, TemplateNo:
		return nil
	}
	return []error{fmt.Errorf("invalid template style %s", ts)}
}

const (
	// TemplateDefault is default template style, uses interpolation of state variables.
	TemplateDefault = TemplateStyle("default")
	// TemplateGo uses go template engine.
	TemplateGo = TemplateStyle("go")
	// TemplateNo uses no template engine.
	TemplateNo = TemplateStyle("no")
)

type MessageReply struct {
	Text      string
	ParseMode ParseMode
	Markup    *ReplyMarkup
	Template  TemplateStyle
}

type InlineButton struct {
	Text     string `yaml:"text"`
	URL      string `yaml:"url"`
	Callback string `yaml:"callback"`
}

type ReplyMarkup struct {
	Keyboard       [][]string       `yaml:"keyboard" json:"keyboard"`
	InlineKeyboard [][]InlineButton `yaml:"inlineKeyboard" json:"inlineKeyboard"`
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

type FileReply struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

func (r *FileReply) validate() []error {
	var errs []error
	if r.Name == "" {
		errs = append(errs, errors.New("empty image name"))
	}
	if r.Key == "" {
		errs = append(errs, errors.New("empty image file"))
	}
	return errs
}
