package spec

import (
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Strings is a slice of strings that can be unmarshalled from YAML scalars or
// sequences.
type Strings []string

func (s *Strings) UnmarshalYAML(node *yaml.Node) error {
	return unmarhalYAMLSeqOrScaler(node, (*[]string)(s), func(s string) (string, error) {
		return s, nil
	})
}

// Uints64 is a slice of uint64s that can be unmarshalled from YAML scalars or
// sequences.
type Uints64 []uint64

func (u *Uints64) UnmarshalYAML(node *yaml.Node) error {
	return unmarhalYAMLSeqOrScaler(node, (*[]uint64)(u), func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 64)
	})
}

func unmarhalYAMLSeqOrScaler[T any](node *yaml.Node, target *[]T, parser func(string) (T, error)) error {
	switch node.Kind {
	case yaml.ScalarNode:
		val, err := parser(node.Value)
		if err != nil {
			return errors.Wrap(err, "parsing target")
		}
		*target = []T{val}
	case yaml.SequenceNode:
		*target = make([]T, 0, len(node.Content))
		for i, node := range node.Content {
			if node.Kind == yaml.ScalarNode {
				val, err := parser(node.Value)
				if err != nil {
					return errors.Wrapf(err, "%d: parsing target", i)
				}
				*target = append(*target, val)
			} else if node.Kind == yaml.AliasNode && node.Alias.Kind == yaml.ScalarNode {
				val, err := parser(node.Alias.Value)
				if err != nil {
					return errors.Wrapf(err, "%d: parsing target", i)
				}
				*target = append(*target, val)
			} else {
				return errors.Errorf("%d: expected scalar node, got %v", i, node.Kind)
			}
		}
	case yaml.AliasNode:
		return unmarhalYAMLSeqOrScaler(node.Alias, target, parser)
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
