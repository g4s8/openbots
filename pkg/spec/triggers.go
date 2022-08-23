package spec

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Trigger struct {
	Message  *MessageTrigger  `yaml:"message"`
	Callback *CallbackTrigger `yaml:"callback"`
}

func (t *Trigger) validate() (errs []error) {
	errs = make([]error, 0)
	if t.Message == nil && t.Callback == nil {
		errs = append(errs, errors.New("empty trigger"))
	}
	if t.Message != nil {
		errs = append(errs, t.Message.validate()...)
	}
	if t.Callback != nil {
		errs = append(errs, t.Callback.validate()...)
	}
	return
}

type MessageTrigger struct {
	Text    []string
	Command string
}

func (t *MessageTrigger) validate() []error {
	if len(t.Text) == 0 && t.Command == "" {
		return []error{errors.New("empty message trigger")}
	}
	return []error{}
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

func (t *CallbackTrigger) validate() []error {
	if t.Data == "" {
		return []error{errors.New("empty callback trigger")}
	}
	return []error{}
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
