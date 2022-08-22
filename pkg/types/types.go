// Package types provides base public types used by bot.
package types

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler of telegram update message.
type Handler interface {
	Handle(context.Context, *telegram.Update, *telegram.BotAPI) error
}

// EventFilter checks that telegram update could be handlerd.
type EventFilter interface {
	Check(context.Context, *telegram.Update) bool
}

// Bot instance.
type Bot interface {
	// Handle register bot's handler with filter.
	Handle(EventFilter, Handler)

	// Start bot.
	Start() error

	// Stop bot.
	Stop() error
}
