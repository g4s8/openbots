package spec

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Trigger struct {
	Message  *MessageTrigger  `yaml:"message"`
	Callback *CallbackTrigger `yaml:"callback"`
}

type MessageTrigger struct {
	Text    []string
	Command string
}

func (t *MessageTrigger) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode, yaml.SequenceNode, yaml.AliasNode:
		var s yamlScalarOrSeq
		if err := node.Decode(&s); err != nil {
			return err
		}
		t.Text = s.Value
	case yaml.MappingNode:
		schema := &struct {
			Text    yamlScalarOrSeq `yaml:"text"`
			Command string          `yaml:"command"`
		}{}
		if err := node.Decode(schema); err != nil {
			return err
		}
		t.Text = schema.Text.Value
		t.Command = schema.Command
	default:
		return errors.Errorf("unexpected node kind: %v", node.Kind)
	}
	return nil
}

type CallbackTrigger struct {
	Data string
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
		return errors.Errorf("unexpected node kind: %v", node.Kind)
	}
	return nil
}
