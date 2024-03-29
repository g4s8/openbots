package spec

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
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

func (r *Reply) validate() (errs []error) {
	errs = make([]error, 0)
	if r.Message == nil && r.Callback == nil && r.Edit == nil && !r.Delete &&
		r.Image == nil && r.Document == nil && r.Invoice == nil && r.PreCheckout == nil {
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
	if r.Image != nil {
		errs = append(errs, r.Image.validate()...)
	}
	if r.Document != nil {
		errs = append(errs, r.Document.validate()...)
	}
	if r.Invoice != nil {
		errs = append(errs, r.Invoice.validate()...)
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

func (r *MessageReply) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		r.Text = node.Value
	case yaml.AliasNode:
		return r.UnmarshalYAML(node.Alias)
	case yaml.MappingNode:
		schema := &struct {
			Text      string        `yaml:"text"`
			ParseMode ParseMode     `yaml:"parseMode"`
			Markup    *ReplyMarkup  `yaml:"markup"`
			Template  TemplateStyle `yaml:"template"`
		}{}
		if err := node.Decode(schema); err != nil {
			return err
		}
		r.Text = schema.Text
		r.ParseMode = schema.ParseMode
		r.Markup = schema.Markup
		r.Template = schema.Template
	default:
		return fmt.Errorf("unexpected node kind: %v", node.Kind)
	}
	return nil
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
	if r.Template == "" {
		r.Template = TemplateDefault
	}
	errs = append(errs, r.Template.validate()...)

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
