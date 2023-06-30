package secrets

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
)

type stub struct{}

func (s stub) Get(ctx context.Context) (map[string]types.Secret, error) {
	return map[string]types.Secret{}, nil
}

var _ types.Secrets = stub{}

// Stub is a stub implementation of secrets.
var Stub = stub{}
