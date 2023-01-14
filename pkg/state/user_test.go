package state

import (
	"fmt"
	"reflect"
	"testing"

	m "github.com/g4s8/go-matchers"
)

type getMatcher struct {
	assert m.Assertion

	key string
	has bool
	val string
}

func (mc *getMatcher) Check(x interface{}) bool {
	s, ok := x.(*UserState)
	if !ok {
		return false
	}
	val, ok := s.Get(mc.key)
	if mc.has != ok {
		return false
	}
	if ok {
		if !reflect.DeepEqual(val, mc.val) {
			return false
		}
	}
	return true
}

func (mc *getMatcher) String() string {
	return fmt.Sprintf("expect has=%v, val=%v", mc.has, mc.val)
}

func matchGet(assert m.Assertion, key string, has bool, val string) m.Matcher {
	return &getMatcher{assert: assert, key: key, has: has, val: val}
}

func newState() *UserState {
	s := NewUserState()
	s.Set("one", "1")
	s.Set("two", "2")
	return s
}

func TestUserState(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		s := newState()
		assert := m.Assert(t)
		assert.That("one", s, matchGet(assert, "one", true, "1"))
		assert.That("two", s, matchGet(assert, "two", true, "2"))
		assert.That("three", s, matchGet(assert, "three", false, ""))
	})
	t.Run("set", func(t *testing.T) {
		s := newState()
		assert := m.Assert(t)
		s.Set("one", "11")
		s.Set("three", "3")
		assert.That("one", s, matchGet(assert, "one", true, "11"))
		assert.That("two", s, matchGet(assert, "two", true, "2"))
		assert.That("three", s, matchGet(assert, "three", true, "3"))
	})
	t.Run("delete", func(t *testing.T) {
		s := newState()
		assert := m.Assert(t)
		s.Delete("one")
		assert.That("one", s, matchGet(assert, "one", false, ""))
		assert.That("two", s, matchGet(assert, "two", true, "2"))
	})
	t.Run("map", func(t *testing.T) {
		s := newState()
		assert := m.Assert(t)
		assert.That("map", s.Map(), m.Eq(map[string]string{
			"one": "1",
			"two": "2",
		}))
	})
	t.Run("changes", func(t *testing.T) {
		s := newState()
		assert := m.Assert(t)
		s.Set("one", "11")
		s.Set("three", "3")
		s.Delete("two")
		changes := s.Changes()
		assert.That("added", changes.Added, m.Eq([]string{"one", "three"}))
		assert.That("deleted", changes.Removed, m.Eq([]string{"two"}))
	})
}
