package spec

import (
	"errors"
	"fmt"
)

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

func (pm ParseMode) validate() []error {
	switch pm {
	case ModeMarkdown, ModeMarkdownV2, ModeHTML:
		return nil
	}
	return []error{fmt.Errorf("invalid parse mode %s", pm)}
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
