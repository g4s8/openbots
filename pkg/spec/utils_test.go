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

func TestStrings(t *testing.T) {
	src := `
empty: []
foo: "foo"
bar: ["qwe", "asd", "zxc"]
baz:
- wer
- sdf
- xcv
`
	type spec struct {
		Empty Strings `yaml:"empty"`
		Foo   Strings `yaml:"foo"`
		Bar   Strings `yaml:"bar"`
		Baz   Strings `yaml:"baz"`
	}
	var s spec
	if err := yaml.Unmarshal([]byte(src), &s); err != nil {
		t.Fatal(err)
	}

	if len(s.Empty) != 0 {
		t.Errorf("empty: expected [], got %v", s.Empty)
	}
	if len(s.Foo) != 1 || s.Foo[0] != "foo" {
		t.Errorf("foo: expected [foo], got %v", s.Foo)
	}
	if len(s.Bar) != 3 || s.Bar[0] != "qwe" || s.Bar[1] != "asd" || s.Bar[2] != "zxc" {
		t.Errorf("bar: expected [qwe asd zxc], got %v", s.Bar)
	}
	if len(s.Baz) != 3 || s.Baz[0] != "wer" || s.Baz[1] != "sdf" || s.Baz[2] != "xcv" {
		t.Errorf("baz: expected [wer sdf xcv], got %v", s.Baz)
	}
}
