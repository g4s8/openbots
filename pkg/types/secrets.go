package types

import (
	"context"
)

// Secrets is a collection of secrets.
type Secrets interface {
	// Get secret map, returns error if not found.
	Get(ctx context.Context) (map[string]Secret, error)
}

// Secret is a string that should not be printed
type Secret string

func (s Secret) String() string {
	return "***"
}

func (s Secret) Value() string {
	return string(s)
}
