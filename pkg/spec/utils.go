package spec

import (
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
