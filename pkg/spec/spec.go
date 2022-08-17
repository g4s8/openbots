package spec

import (
	"io"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Spec struct {
	Bot *Bot `yaml:"bot"`
}

type Bot struct {
	Token    string     `yaml:"token" env:"BOT_TOKEN"`
	Handlers []*Handler `yaml:"handlers"`
}

type Handler struct {
	Trigger *Trigger `yaml:"on"`
	Replies []*Reply `yaml:"reply"`
}

type Trigger struct {
	Message *MessageTrigger `yaml:"message"`
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

type yamlScalarOrSeq struct {
	Value []string
}

func (s *yamlScalarOrSeq) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		s.Value = []string{node.Value}
	case yaml.SequenceNode:
		s.Value = make([]string, 0, len(node.Content))
		for i, node := range node.Content {
			if node.Kind == yaml.ScalarNode {
				s.Value = append(s.Value, node.Value)
			} else if node.Kind == yaml.AliasNode && node.Alias.Kind == yaml.ScalarNode {
				s.Value = append(s.Value, node.Alias.Value)
			} else {
				return errors.Errorf("%d: expected scalar node, got %v", i, node.Kind)
			}
		}
	case yaml.AliasNode:
		return s.UnmarshalYAML(node.Alias)
	default:
		return errors.Errorf("unexpected node kind: %v", node.Kind)
	}
	return nil
}

type MessageTrigger struct {
	Text    []string
	Command string
}

type Reply struct {
	Message *MessageReply `yaml:"message"`
}

type MessageReply struct {
	Text   string       `yaml:"text"`
	Markup *ReplyMarkup `yaml:"markup"`
}

type ReplyMarkup struct {
	Keyboard [][]string `yaml:"keyboard"`
}

func ParseYaml(r io.Reader) (*Spec, error) {
	var spec Spec
	if err := yaml.NewDecoder(r).Decode(&spec); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	if err := env.Parse(spec.Bot); err != nil {
		return nil, errors.Wrap(err, "env")
	}
	return &spec, nil
}
