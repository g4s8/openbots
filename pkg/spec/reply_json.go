package spec

import (
	"encoding/json"
	"errors"
	"fmt"
)

var _ (json.Unmarshaler) = (*MessageReply)(nil)

func (m *MessageReply) UnmarshalJSON(data []byte) error {
	// if it's  a string - parse it as text
	// if it's an object, parse it as a struct

	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		m.Text = text
		return nil
	} else if ute := new(json.UnmarshalTypeError); errors.As(err, &ute) && ute.Value == "object" {
		// it's not a string, continue
	} else {
		return fmt.Errorf("unmarshal text: %w", err)
	}

	obj := struct {
		Text      string        `json:"text"`
		ParseMode ParseMode     `json:"parseMode"`
		Markup    *ReplyMarkup  `json:"markup"`
		Template  TemplateStyle `json:"template"`
	}{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("unmarshal object: %w", err)
	}
	m.Text = obj.Text
	m.ParseMode = obj.ParseMode
	m.Markup = obj.Markup
	m.Template = obj.Template
	return nil
}
