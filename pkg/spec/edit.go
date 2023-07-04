package spec

import "github.com/pkg/errors"

type Edit struct {
	Message *EditMessage `yaml:"message"`
}

func (r *Edit) validate() []error {
	var errs []error
	if r.Message == nil {
		errs = append(errs, errors.New("empty edit"))
	}
	if r.Message != nil {
		errs = append(errs, r.Message.validate()...)
	}
	return errs
}

type EditMessage struct {
	Caption        string           `yaml:"caption"`
	Text           string           `yaml:"text"`
	InlineKeyboard [][]InlineButton `yaml:"inlineKeyboard"`
	Template       TemplateStyle    `yaml:"template"`
}

func (r *EditMessage) validate() []error {
	var errs []error

	if r.Text == "" && r.Caption == "" && r.InlineKeyboard == nil {
		errs = append(errs, errors.New("empty edit message"))
	}

	if r.Text != "" && r.Caption != "" {
		errs = append(errs, errors.New("both text and caption are set"))
	}

	if r.Caption != "" && r.InlineKeyboard != nil {
		errs = append(errs, errors.New("caption and inline keyboard are set"))
	}

	if r.Template == "" {
		r.Template = TemplateDefault
	}
	errs = append(errs, r.Template.validate()...)

	return errs
}
