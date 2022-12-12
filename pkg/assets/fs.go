package assets

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/g4s8/openbots/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _ types.Assets = (*FS)(nil)

// FS is assets provider from local filesystem.
type FS struct {
	root   string
	logger zerolog.Logger
}

// NewFS creates new filesystem assets provider.
func NewFS(root string, logger zerolog.Logger) *FS {
	return &FS{
		root: root,
		logger: logger.With().
			Str("component", "assets").
			Str("asset_provider", "fs").
			Str("root", root).
			Logger(),
	}
}

func (fs *FS) LoadAsset(ctx context.Context, key string) (io.ReadCloser, error) {
	f, err := os.Open(path.Join(fs.root, key))
	fs.logger.Info().Str("key", key).Msg("Load asset")
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	return f, nil
}
