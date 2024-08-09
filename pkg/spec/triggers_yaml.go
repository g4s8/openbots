package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func (t *Trigger) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.AliasNode {
		return t.UnmarshalYAML(node.Alias)
	}

	switch node.Kind {
	case yaml.ScalarNode:
		if node.Value == "*" {
			t.Fallback = true
			return nil
		}
		t.Message = &MessageTrigger{Text: []string{node.Value}}
	case yaml.SequenceNode:
		var s Strings
		if err := node.Decode(&s); err != nil {
			return err
		}
		t.Message = &MessageTrigger{Text: s}
	case yaml.MappingNode:
		var schema struct {
			Message      *MessageTrigger      `yaml:"message"`
			Callback     *CallbackTrigger     `yaml:"callback"`
			Context      string               `yaml:"context"`
			PreCheckout  *PreCheckoutTrigger  `yaml:"preCheckout"`
			PostCheckout *PostCheckoutTrigger `yaml:"postCheckout"`
			State        []StateCondition     `yaml:"state"`
			Fallback     bool                 `yaml:"fallback"`
		}
		if err := node.Decode(&schema); err != nil {
			return fmt.Errorf("decode trigger: %w", err)
		}
		t.Message = schema.Message
		t.Callback = schema.Callback
		t.Context = schema.Context
		t.PreCheckout = schema.PreCheckout
		t.PostCheckout = schema.PostCheckout
		t.State = schema.State
		t.Fallback = schema.Fallback
	default:
		return fmt.Errorf("unexpected node kind: %v", node.Kind)
	}
	return nil
}

func (t *MessageTrigger) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode, yaml.SequenceNode, yaml.AliasNode:
		var s Strings
		if err := node.Decode(&s); err != nil {
			return err
		}
		t.Text = s
	case yaml.MappingNode:
		schema := &struct {
			Text    Strings `yaml:"text"`
			Command string  `yaml:"command"`
		}{}
		if err := node.Decode(schema); err != nil {
			return err
		}
		t.Text = schema.Text
		t.Command = schema.Command
	default:
		return fmt.Errorf("unexpected node kind: %v", node.Kind)
	}
	return nil
}

func (ct *CallbackTrigger) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		ct.Data = node.Value
	case yaml.AliasNode:
		return ct.UnmarshalYAML(node.Alias)
	case yaml.MappingNode:
		var schema struct {
			Data string `yaml:"data"`
		}
		if err := node.Decode(&schema); err != nil {
			return err
		}
		ct.Data = schema.Data
	default:
		return fmt.Errorf("unexpected node kind: %v", node.Kind)
	}
	return nil
}
