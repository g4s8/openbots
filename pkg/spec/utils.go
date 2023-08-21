package spec

import (
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

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

// OptUint64 is a uint64 that can be omitted in YAML.
type OptUint64 struct {
	Value uint64
	Valid bool
}

func (o *OptUint64) UnmarshalYAML(node *yaml.Node) error {
	res, err := unmarshalOptYaml(node, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 64)
	})
	if err != nil {
		return err
	}
	o.Value = res
	o.Valid = true
	return nil
}

type OptBool struct {
	Value bool
	Valid bool
}

func (o *OptBool) UnmarshalYAML(node *yaml.Node) error {
	res, err := unmarshalOptYaml(node, strconv.ParseBool)
	if err != nil {
		return err
	}
	o.Value = res
	o.Valid = true
	return nil
}

func unmarshalOptYaml[T any](node *yaml.Node, parser func(string) (T, error)) (T, error) {
	switch node.Kind {
	case yaml.ScalarNode:
		return parser(node.Value)
	case yaml.AliasNode:
		return unmarshalOptYaml(node.Alias, parser)
	default:
		var zero T
		return zero, errors.Errorf("expected scalar or alias node, got %v", node.Kind)
	}
}
