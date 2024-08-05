package json

import (
	"encoding/json"
	"strings"
	"testing"
)

func stringsTester(data string, expect []string, wantErr bool) func(*testing.T) {
	return func(t *testing.T) {
		var actual Strings
		err := json.NewDecoder(strings.NewReader(data)).Decode(&actual)
		if wantErr && err == nil {
			t.Fatalf("Expect error but decoded successfully")
		}
		if err != nil && !wantErr {
			t.Fatalf("Failed to decode data: %v", err)
		}
		if err != nil && wantErr {
			t.Logf("Got expected error: %v", err)
			return
		}

		if len(expect) != len(actual) {
			t.Fatalf("Unexpected target length; want %d, got %d",
				len(expect), len(actual))
		}
		for i, itemExpect := range expect {
			itemActual := actual[i]
			if itemExpect != itemActual {
				t.Errorf("Unexpected item at position %d; want %q, got %q",
					i, itemExpect, itemActual)
			}

		}
	}
}

func TestStrings(t *testing.T) {
	t.Run("Array", stringsTester(`["one", "two", "three"]`, []string{"one", "two", "three"}, false))
	t.Run("Singleton", stringsTester(`"single"`, []string{"single"}, false))
	t.Run("Null", stringsTester(`null`, []string{}, false))
	t.Run("Error1", stringsTester(`{"foo": "bar"}`, nil, true))
	t.Run("Error2", stringsTester(`42`, nil, true))
}
