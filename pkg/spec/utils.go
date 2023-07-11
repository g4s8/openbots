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
	if node.Kind == yaml.ScalarNode {
		val, err := strconv.ParseUint(node.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "parse %q as uint64", node.Value)
		}
		o.Value = val
		o.Valid = true
	} else {
		return errors.Errorf("expected scalar node, got %v", node.Kind)
	}
	return nil
}
