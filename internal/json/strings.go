package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var _ json.Unmarshaler = (*Strings)(nil)

// Strings decodes JSON stirng or array of strings or `null` into Go slice of string:
//   - If JSON value is an array of string, it's decoded directly into Go slice:
//     ["one", "two", "three"] -> ["one", "tho", "three"]
//   - JSON string value decoded into singleton Go slice:
//     "Helllo" -> ["Hello"]
//   - JSON `null` value decoded into an empty slice:
//     null -> []
type Strings []string

func (s *Strings) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	var inArray bool
	var values []string
	for {
		t, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("decode token: %w", err)
		}

		switch tt := t.(type) {
		case json.Delim:
			switch tt {
			case '[':
				if inArray {
					return errors.New("enter nested array")
				}
				if len(values) > 0 {
					return errors.New("starting array with not empty values")
				}
				inArray = true
			case ']':
				if !inArray {
					return errors.New("unexpected closing array bracet")
				}
				inArray = false
			default:
				return fmt.Errorf("unsupported delim: %s", t)
			}
		case string:
			values = append(values, tt)
		default:
			if tt == nil {
				continue
			}
			return fmt.Errorf("unsupported token type: %T (%v)", tt, tt)
		}
	}
	*s = values
	return nil
}
