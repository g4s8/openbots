package spec

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestReplyJson(t *testing.T) {
	t.Run("parse text", func(t *testing.T) {
		data := []byte(`"hello"`)
		var r MessageReply
		if err := json.Unmarshal(data, &r); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.Text != "hello" {
			t.Fatalf("unexpected text: %s", r.Text)
		}
	})
	t.Run("parse object", func(t *testing.T) {
		data := []byte(`{"text":"hello","parseMode":"markdown"}`)
		var r MessageReply
		if err := json.Unmarshal(data, &r); err != nil {
			t.Fatalf("unexpected error: %v %T", err, errors.Unwrap(err))
		}
		if r.Text != "hello" {
			t.Fatalf("unexpected text: %s", r.Text)
		}
		if r.ParseMode != "markdown" {
			t.Fatalf("unexpected parse mode: %s", r.ParseMode)
		}
	})
	t.Run("parse error", func(t *testing.T) {
		data := []byte(`{"text":123}`)
		var r MessageReply
		if err := json.Unmarshal(data, &r); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
	t.Run("parse error type", func(t *testing.T) {
		data := []byte(`[1,2,3]`)
		var r MessageReply
		if err := json.Unmarshal(data, &r); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
