package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

var _ types.Handler = (*Validator)(nil)

// Check is a type for message validator checks.
type Check string

func ParseCheckString(str string) (Check, error) {
	switch c := Check(str); c {
	case CheckNotEmpty, CheckIsInt, CheckIsFloat, CheckIsBool:
		return c, nil
	}
	return Check(""), fmt.Errorf("unknown validator: %q", str)
}

func (c Check) perform(upd *telegram.Update) error {
	if upd.Message == nil {
		return fmt.Errorf("update is not a message")
	}
	text := upd.Message.Text
	switch c {
	case CheckNotEmpty:
		if text == "" {
			return fmt.Errorf("message is empty")
		}
	case CheckIsInt:
		if _, err := strconv.ParseInt(text, 10, 64); err != nil {
			return fmt.Errorf("message is not an integer: %w", err)
		}
	case CheckIsFloat:
		if _, err := strconv.ParseFloat(text, 64); err != nil {
			return fmt.Errorf("message is not a float: %w", err)
		}
	case CheckIsBool:
		if _, err := strconv.ParseBool(text); err != nil {
			return fmt.Errorf("message is not a bool: %w", err)
		}
	default:
		return fmt.Errorf("unknown validator: %q", c)
	}
	return nil
}

const (
	CheckNotEmpty Check = "not_empty"
	CheckIsInt    Check = "is_int"
	CheckIsFloat  Check = "is_float"
	CheckIsBool   Check = "is_bool"
)

// Validator is a struct for message validation.
type Validator struct {
	errMessage string
	checks     []Check
	logger     zerolog.Logger
}

// NewValidator is a constructor for Validator.
// It accepts a list of checks to perform and an error message to send if validation fails.
func NewValidator(logger zerolog.Logger, errMessage string, checks ...Check) *Validator {
	return &Validator{
		errMessage: errMessage,
		checks:     checks,
		logger:     logger,
	}
}

var ErrValidationFailed = fmt.Errorf("validation failed")

func (v *Validator) Handle(ctx context.Context, upd *telegram.Update, api *telegram.BotAPI) error {
	for _, check := range v.checks {
		if err := check.perform(upd); err != nil {
			v.logger.Debug().Err(err).Msg("Validation failed")
			if upd.Message != nil && upd.Message.Chat != nil {
				msg := telegram.NewMessage(upd.Message.Chat.ID, v.errMessage)
				if _, err := api.Send(msg); err != nil {
					return fmt.Errorf("failed to send validation error message: %w", err)
				}
			}
			return ErrValidationFailed
		}
	}
	return nil
}
