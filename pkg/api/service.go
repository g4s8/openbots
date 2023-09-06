package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"

	"github.com/g4s8/openbots/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var rePath = regexp.MustCompile(`^/handlers/([a-zA-Z0-9-]+)$`)

type Service struct {
	cfg      Config
	handlers map[string]Handler
	logger   zerolog.Logger

	srv *http.Server
}

func NewService(cfg Config, handlers map[string]Handler) *Service {
	return NewServiceWithLogger(cfg, handlers, zerolog.Nop())
}

func NewServiceWithLogger(cfg Config, handlers map[string]Handler, logger zerolog.Logger) *Service {
	return &Service{
		cfg:      cfg,
		handlers: handlers,
		logger:   logger,
	}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		defer req.Body.Close()
	}

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !rePath.MatchString(req.URL.Path) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	matches := rePath.FindStringSubmatch(req.URL.Path)
	name := matches[1]
	handler, ok := s.handlers[name]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logger := s.logger.With().
		Str("method", req.Method).
		Str("path", req.URL.Path).
		Str("handler", name).
		Logger()

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.RequestTimeout)
	defer cancel()

	var jsonPayload struct {
		ChatID  int64             `json:"chat_id"`
		Payload map[string]string `json:"params"`
		Data    map[string]string `json:"data"`
	}
	if err := json.NewDecoder(req.Body).Decode(&jsonPayload); err != nil {
		logger.Err(err).Msg("Invalid payload")
		http.Error(w, fmt.Sprintf("invalid payload: %v", err), http.StatusBadRequest)
		return
	}
	if jsonPayload.Payload == nil && jsonPayload.Data != nil {
		jsonPayload.Payload = jsonPayload.Data
	}

	select {
	case <-ctx.Done():
		http.Error(w, "request timeout", http.StatusRequestTimeout)
		return
	default:
	}

	payload := Request{
		ChatID:  types.ChatID(jsonPayload.ChatID),
		Payload: jsonPayload.Payload,
	}
	err := handler.Call(ctx, payload)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var apiErr *Error
	if errors.As(err, &apiErr) {
		switch apiErr.kind {
		case InvalidRequestDataError:
			http.Error(w, apiErr.Error(), http.StatusBadRequest)
		default:
			http.Error(w, apiErr.Error(), http.StatusInternalServerError)
		}
		logger.Err(err).Object("details", apiErr).Msg("API error")
		return
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		logger.Err(err).Msg("Request timeout")
		http.Error(w, "request timeout", http.StatusRequestTimeout)
		return
	}
	logger.Err(err).Msg("Internal server error")
	http.Error(w, fmt.Sprintf("Internal server error: %v", err),
		http.StatusInternalServerError)
}

func (s *Service) Start(ctx context.Context) error {
	s.srv = &http.Server{
		Addr:        s.cfg.Addr,
		ReadTimeout: s.cfg.ReadTimeout,
	}
	mux := http.NewServeMux()
	mux.Handle("/handlers/", s)
	mux.Handle("/health", &health{}) // TODO: impl
	s.srv.Handler = mux
	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return errors.Wrap(err, "net listen")
	}
	go func() {
		if err := s.srv.Serve(ln); err != nil {
			s.logger.Err(err).Msg("HTTP server failed")
		}
	}()
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "server shutdown")
	}
	if err := s.srv.Close(); err != nil {
		return errors.Wrap(err, "server close")
	}
	return nil
}
