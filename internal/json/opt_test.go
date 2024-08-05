package json

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func optTesterCmp[T any](data string, want T, wantOK bool, cmp func(T, T) bool) func(*testing.T) {
	return func(t *testing.T) {
		var opt Opt[T]
		if err := json.NewDecoder(strings.NewReader(data)).Decode(&opt); err != nil {
			t.Fatalf("Failed to decode: %v", err)
		}
		got, ok := opt.Value()
		if ok != wantOK {
			if wantOK {
				t.Fatal("Wanted value bug got nothing")
			} else {
				t.Fatalf("Got value but not wanted")
			}
		}
		if !cmp(want, got) {
			t.Fatalf("Unexpected value; want %v but got %v", want, got)
		}
	}
}

func optTester[T comparable](data string, want T, wantOK bool) func(*testing.T) {
	return optTesterCmp(data, want, wantOK, func(a, b T) bool {
		return a == b
	})
}

func TestOpt(t *testing.T) {
	t.Run("JustInt", optTester(`42`, 42, true))
	t.Run("JustString", optTester(`"hello"`, "hello", true))
	t.Run("JustBool", optTester(`true`, true, true))
	t.Run("EmptyInt", optTester(`null`, 0, false))
	type testStruct struct {
		Foo   string      `json:"foo"`
		Num   int         `json:"num"`
		Empty interface{} `json:"empty"`
	}
	t.Run("Object", optTesterCmp(`{
		"foo": "bar",
		"num": 4,
		"empty": null
	}`, &testStruct{
		Foo:   "bar",
		Num:   4,
		Empty: nil,
	}, true, func(a, b *testStruct) bool {
		return *a == *b
	}))
	t.Run("Array", optTesterCmp(`[1,2,3]`,
		[]int{1, 2, 3}, true, func(a, b []int) bool {
			return slices.Compare(a, b) == 0
		},
	))

	t.Run("Methods", func(t *testing.T) {
		var opt Opt[int]
		opt.value = 8
		if x := opt.OrDefault(42); x != 42 {
			t.Errorf("Expected default 42 but was %d", x)
		}
		if _, ok := opt.Value(); ok {
			t.Error("Expected opt has no value")
		}
		opt.decoded = true
		if x := opt.OrDefault(42); x != 8 {
			t.Errorf("Expected default 8 but was %d", x)
		}
		x, ok := opt.Value()
		if !ok {
			t.Error("Expected opt has value")
		}
		if x != 8 {
			t.Errorf("Expected opt has value 8 but was %d", x)
		}
	})
}
