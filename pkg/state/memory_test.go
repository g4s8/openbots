package state

import (
	"testing"

	m "github.com/g4s8/go-matchers"
	"github.com/g4s8/openbots/pkg/types"
)

func TestMemory(t *testing.T) {
	assert := m.Assert(t)

	mem := NewMemory(map[string]string{
		"foo": "bar",
	})
	state := NewUserState()
	_ = mem.Load(nil, types.ChatID(1), state)
	state.Set("foo", "baz")
	state.Set("one", "1")
	mem.Update(nil, types.ChatID(1), state)

	state.reset()
	_ = mem.Load(nil, types.ChatID(2), state)
	assert.That("second-get-foo", state, matchGet(assert, "foo", true, "bar"))
	state.Set("two", "2")
	mem.Update(nil, types.ChatID(2), state)

	state.reset()
	_ = mem.Load(nil, types.ChatID(1), state)
	assert.That("first-get-foo", state, matchGet(assert, "foo", true, "baz"))
	assert.That("first-get-one", state, matchGet(assert, "one", true, "1"))
	state.Set("one", "один")
	mem.Update(nil, types.ChatID(1), state)

	state.reset()
	_ = mem.Load(nil, types.ChatID(2), state)
	state.Delete("two")
	mem.Update(nil, types.ChatID(2), state)

	state.reset()
	_ = mem.Load(nil, types.ChatID(1), state)
	assert.That("first-get-foo", state, matchGet(assert, "foo", true, "baz"))
	assert.That("first-get-one", state, matchGet(assert, "one", true, "один"))
	assert.That("first-get-two", state, matchGet(assert, "two", false, ""))

	state.reset()
	_ = mem.Load(nil, types.ChatID(2), state)
	assert.That("second-get-foo", state, matchGet(assert, "foo", true, "bar"))
	assert.That("second-get-one", state, matchGet(assert, "one", false, ""))
	assert.That("second-get-two", state, matchGet(assert, "two", false, ""))
}
