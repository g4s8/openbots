// Package types provides base public types used by bot.
package types

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ChatID is a chat identifier.
type ChatID int64

// Handler of telegram update message.
type Handler interface {
	Handle(context.Context, *telegram.Update, *telegram.BotAPI) error
}

// StateHandler for state updates.
type StateHandler interface {
	Handle(context.Context, *telegram.Update) error
}

// EventFilter checks that telegram update could be handlerd.
type EventFilter interface {
	Check(context.Context, *telegram.Update) (bool, error)
}

// Bot instance.
type Bot interface {
	// Handle register bot's handler with filter.
	Handle(EventFilter, Handler)

	// HandleState register bot's state handler with filter.
	HandleState(EventFilter, StateHandler)

	// Start bot.
	Start() error

	// Stop bot.
	Stop() error
}
