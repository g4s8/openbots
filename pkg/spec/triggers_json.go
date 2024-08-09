package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	j "github.com/g4s8/openbots/internal/json"
)

var _ json.Unmarshaler = (*MessageTrigger)(nil)

func (t *MessageTrigger) UnmarshalJSON(data []byte) error {
	reader := bytes.NewReader(data)
	{
		var plain j.Strings
		err := json.NewDecoder(reader).Decode(&plain)
		if err == nil {
			t.Text = plain
			return nil
		}
		reader.Seek(0, io.SeekStart)
	}
	{
		schema := &struct {
			Text    j.Strings `json:"text"`
			Command string    `json:"command"`
		}{}
		if err := json.NewDecoder(reader).Decode(schema); err != nil {
			return fmt.Errorf("decode message trigger: %w", err)
		}
		t.Text = schema.Text
		t.Command = schema.Command
	}
	return nil
}
