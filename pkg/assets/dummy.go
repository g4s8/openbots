package assets

import (
	"context"
	"errors"
	"io"

	"github.com/g4s8/openbots/pkg/types"
)

var _ types.Assets = (*dummy)(nil)

type dummy int

func (d dummy) LoadAsset(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, errors.New("not found")
}

// Dummy assets provider.
var Dummy types.Assets = dummy(0)
