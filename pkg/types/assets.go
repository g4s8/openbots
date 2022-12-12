package types

import (
	"context"
	"io"
)

// Assets provider interface.
type Assets interface {
	// LoadAsset loads the asset with the given key.
	LoadAsset(ctx context.Context, key string) (io.ReadCloser, error)
}
