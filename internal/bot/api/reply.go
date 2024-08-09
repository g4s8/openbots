package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/g4s8/openbots/internal/bot/adaptors"
	"github.com/g4s8/openbots/internal/bot/handlers"
	"github.com/g4s8/openbots/pkg/api"
	"github.com/g4s8/openbots/pkg/types"
	"github.com/rs/zerolog"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ReplyService struct {
	tg      *telegram.BotAPI
	secrets types.Secrets
	sp      types.StateProvider
	logger  zerolog.Logger
}

func NewReplyService(tg *telegram.BotAPI, secrets types.Secrets, sp types.StateProvider,
	logger zerolog.Logger,
) *ReplyService {
	return &ReplyService{
		tg:      tg,
		secrets: secrets,
		sp:      sp,
		logger:  logger,
	}
}

func (s *ReplyService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		defer req.Body.Close()
	}
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data api.ReplyRequest
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		s.logger.Debug().Err(err).Msg("decode request")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	handler, err := adaptors.MessageRepply(s.tg, s.sp, s.secrets, data.Spec, s.logger)
	if err != nil {
		s.logger.Err(err).Msg("create handler")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	uctxp := handlers.NewUpdateContextProvider(s.secrets, s.sp)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	upd := &telegram.Update{
		Message: &telegram.Message{
			Chat: &telegram.Chat{
				ID: data.ChatID.Int64(),
			},
		},
	}
	uctx, err := uctxp.NewContext(ctx, upd)
	if err != nil {
		s.logger.Err(err).Msg("create update context")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := handler.Handle(uctx, nil, s.tg); err != nil {
		s.logger.Err(err).Msg("handle message")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
