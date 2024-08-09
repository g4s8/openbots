package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

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
