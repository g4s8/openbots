package spec

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOptUint64(t *testing.T) {
	src := `
foo: 1
bar: 2235782385482385
`
	type spec struct {
		Foo OptUint64 `yaml:"foo"`
		Bar OptUint64 `yaml:"bar"`
		Baz OptUint64 `yaml:"baz"`
	}
	var s spec

	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatal(err)
	}

	if !s.Foo.Valid || s.Foo.Value != 1 {
		t.Errorf("foo: expected 1, got %v", s.Foo)
	}
	if !s.Bar.Valid || s.Bar.Value != 2235782385482385 {
		t.Errorf("bar: expected 2235782385482385, got %v", s.Bar)
	}
	if s.Baz.Valid {
		t.Errorf("baz: expected invalid, got %v", s.Baz)
	}
}
