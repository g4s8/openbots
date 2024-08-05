package json

import (
	"encoding/json"
	"fmt"
)

type Opt[T any] struct {
	value   T
	decoded bool
}

func (o *Opt[T]) init() {
	// `Opt` could be reused potentially,
	// we have to be ensure it's in initial state before unmarshaling.
	o.decoded = false
	var def T
	o.value = def
}

func (o *Opt[T]) UnmarshalJSON(data []byte) error {
	o.init()

	var val *T
	if err := json.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	if val != nil {
		o.value = *val
		o.decoded = true
	}
	return nil
}

func (o *Opt[T]) Value() (T, bool) {
	return o.value, o.decoded
}

func (o *Opt[T]) OrDefault(def T) T {
	if o.decoded {
		return o.value
	}
	return def
}

func (o *Opt[T]) Decoded() bool {
	return o.decoded
}
