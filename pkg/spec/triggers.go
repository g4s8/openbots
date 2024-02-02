package spec

import (
	"errors"
	"fmt"
	"slices"

	"gopkg.in/yaml.v3"
)

// Trigger is a handler trigger whcich configures when the handler should be
// executed.
type Trigger struct {
	Message      *MessageTrigger
	Callback     *CallbackTrigger
	Context      string
	PreCheckout  *PreCheckoutTrigger
	PostCheckout *PostCheckoutTrigger
	State        []StateCondition
	Fallback     bool
}

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

// TriggerType is a type of trigger.
//
//go:generate stringer -type=TriggerType
type TriggerType int

// Trigger types.
const (
	TriggerTypeMessage TriggerType = 1 + iota
	TriggerTypeCallback
	TriggerTypeContext
	TriggerTypePreCheckout
	TriggerTypePostCheckout
	TriggerTypeState
	TriggerTypeFallback
)

// Types returns a list of trigger types.
func (t *Trigger) Types() []TriggerType {
	var typ []TriggerType
	if t.Message != nil {
		typ = append(typ, TriggerTypeMessage)
	}
	if t.Callback != nil {
		typ = append(typ, TriggerTypeCallback)
	}
	if t.Context != "" {
		typ = append(typ, TriggerTypeContext)
	}
	if t.PreCheckout != nil {
		typ = append(typ, TriggerTypePreCheckout)
	}
	if t.PostCheckout != nil {
		typ = append(typ, TriggerTypePostCheckout)
	}
	if len(t.State) > 0 {
		typ = append(typ, TriggerTypeState)
	}
	if t.Fallback {
		typ = append(typ, TriggerTypeFallback)
	}
	return typ
}

// ErrEmptyTrigger is returned when no triggers are specified.
var ErrEmptyTrigger = errors.New("empty trigger")

// ErrInvalidTriggerCombination is returned when trigger combination is invalid.
var ErrInvalidTriggerCombination = errors.New("invalid trigger combination")

func (t *Trigger) validate() error {
	types := t.Types()
	if len(types) == 0 {
		return ErrEmptyTrigger
	}
	// message, callback, preCheckout, postCheckout, fallback could be combined in any combination
	// any type except fallback could be combined with context and state types
	// fallback could not be combined with any other type
	if len(types) > 1 && slices.Contains(types, TriggerTypeFallback) {
		return fmt.Errorf("fallback with other triggers: %w", ErrInvalidTriggerCombination)
	}
	unmixable := []TriggerType{TriggerTypeMessage, TriggerTypeCallback, TriggerTypePreCheckout, TriggerTypePostCheckout}
	var unmixableCnt int
	for _, u := range unmixable {
		if slices.Contains(types, u) {
			unmixableCnt++
		}
	}
	if unmixableCnt > 1 {
		return fmt.Errorf("unmixable triggers (%s): %w", types, ErrInvalidTriggerCombination)
	}

	var errs []error
	if t.Message != nil {
		errs = append(errs, t.Message.validate()...)
	}
	if t.Callback != nil {
		errs = append(errs, t.Callback.validate()...)
	}
	if len(t.State) > 0 {
		for _, sc := range t.State {
			errs = append(errs, sc.validate()...)
		}
	}
	return errors.Join(errs...)
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
		return fmt.Errorf("unexpected node kind: %v", node.Kind)
	}
	return nil
}

type StateCondition struct {
	Key     string  `yaml:"key"`
	Present OptBool `yaml:"present"`
	Eq      string  `yaml:"eq"`
	NEq     string  `yaml:"neq"`
}

func (sc *StateCondition) validate() []error {
	if sc.Key == "" {
		return []error{errors.New("empty state condition key")}
	}
	if !sc.Present.Valid && sc.Eq == "" && sc.NEq == "" {
		return []error{errors.New("empty state condition value")}
	}
	return []error{}
}
