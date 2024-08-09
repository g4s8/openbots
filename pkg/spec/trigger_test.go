package spec

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type codec interface {
	fmt.Stringer
	decode(src string, out any) error
}

type (
	yamlCodec int
	jsonCodec int
)

var (
	codecYaml = yamlCodec(0)
	codecJson = jsonCodec(0)
)

func (c yamlCodec) decode(src string, out any) error {
	r := strings.NewReader(src)
	return yaml.NewDecoder(r).Decode(out)
}

func (c yamlCodec) String() string {
	return "yaml"
}

func (c jsonCodec) decode(src string, out any) error {
	r := strings.NewReader(src)
	return json.NewDecoder(r).Decode(out)
}

func (c jsonCodec) String() string {
	return "json"
}

func messageTriggerTester(data string, expect MessageTrigger, c codec) func(*testing.T) {
	return func(t *testing.T) {
		var actual MessageTrigger
		if err := c.decode(data, &actual); err != nil {
			t.Fatalf("Failed to decode data: %v", err)
		}
		if actual.Command != expect.Command {
			t.Errorf("Command didn't match: want %q, got %q", expect.Command, actual.Command)
		}
		if len(actual.Text) != len(expect.Text) {
			t.Errorf("Text length didn't match: want %d, got %d", len(expect.Text), len(actual.Text))
		} else {
			for i, expectText := range expect.Text {
				actualText := actual.Text[i]
				if expectText != actualText {
					t.Errorf("Text value at position %d didn't match: want %q, got %q",
						i, expectText, actualText)
				}
			}
		}
	}
}

func TestMessageTrigger(t *testing.T) {
	type testCase struct {
		name   string
		data   []string
		expect MessageTrigger
		codec  []codec
	}

	for _, tc := range []testCase{
		{
			name: "command",
			data: []string{
				`command: "start"`,
				`{"command": "start"}`,
			},
			expect: MessageTrigger{
				Command: "start",
			},
			codec: []codec{codecYaml, codecJson},
		},
		{
			name: "text-singleton",
			data: []string{
				`text: single`,
				`{"text": "single"}`,
			},
			expect: MessageTrigger{
				Text: []string{"single"},
			},
			codec: []codec{codecYaml, codecJson},
		},
		{
			name: "text-array",
			data: []string{
				`text: ["one", "two", "three"]`,
				`{"text": ["one", "two", "three"]}`,
			},
			expect: MessageTrigger{
				Text: []string{"one", "two", "three"},
			},
			codec: []codec{codecYaml, codecJson},
		},
		{
			name: "text-embed",
			data: []string{
				`"trigger"`,
				`"trigger"`,
			},
			expect: MessageTrigger{
				Text: []string{"trigger"},
			},
			codec: []codec{codecYaml, codecJson},
		},
	} {
		if len(tc.data) != len(tc.codec) {
			panic("codecs should match data sources")
		}
		t.Run(tc.name, func(t *testing.T) {
			for i, data := range tc.data {
				c := tc.codec[i]
				t.Run(c.String(), messageTriggerTester(data, tc.expect, c))
			}
		})
	}
}
