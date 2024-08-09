package spec

import (
	"errors"
	"fmt"
	"slices"
)

var (
	// ErrEmptyTrigger is returned when no triggers are specified.
	ErrEmptyTrigger = errors.New("empty trigger")

	// ErrInvalidTriggerCombination is returned when trigger combination is invalid.
	ErrInvalidTriggerCombination = errors.New("invalid trigger combination")

	// ErrEmptyMessageTrigger is returned when message trigger data is empty.
	ErrEmptyMessageTrigger = errors.New("empty message trigger")

	// ErrEmptyCallbackTrigger is returned when callback trigger data is empty.
	ErrEmptyCallbackTrigger = errors.New("empty callback trigger")

	// ErrEmptyStateConditionKey is returned when state condition key is empty.
	ErrEmptyStateConditionKey = errors.New("empty state condition key")

	// ErrEmptyStateConditionValue is returned when state condition value is empty.
	ErrEmptyStateConditionValue = errors.New("empty state condition value")
)

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

func (t *MessageTrigger) validate() []error {
	if len(t.Text) == 0 && t.Command == "" {
		return []error{ErrEmptyTrigger}
	}
	return []error{}
}

func (t *CallbackTrigger) validate() []error {
	if t.Data == "" {
		return []error{ErrEmptyCallbackTrigger}
	}
	return []error{}
}

func (sc *StateCondition) validate() []error {
	if sc.Key == "" {
		return []error{ErrEmptyStateConditionKey}
	}
	if !sc.Present.Valid && sc.Eq == "" && sc.NEq == "" {
		return []error{ErrEmptyStateConditionValue}
	}
	return []error{}
}
