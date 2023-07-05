package adaptors

import (
	"net/http"

	"github.com/g4s8/openbots/internal/bot/data"
	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/types"
	"github.com/rs/zerolog"
)

func DataLoader(sp types.StateProvider, secrets types.Secrets, s *spec.Data, log zerolog.Logger) (*data.Loader, error) {
	if s.Fetch != nil {
		cfg := data.LoaderConfig{
			Method:  s.Fetch.Method,
			URL:     s.Fetch.URL,
			Headers: s.Fetch.Headers,
		}
		logger := log.With().Str("component", "data_loader").Logger()
		return data.NewLoader(http.DefaultClient, cfg, sp, secrets, logger), nil
	}
	return nil, nil
}
