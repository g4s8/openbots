// Package types provides base public types used by bot.
package types

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	Handle(context.Context, *telegram.Update, *telegram.BotAPI) error
}

type EventFilter interface {
	Check(context.Context, *telegram.Update) bool
}

type Bot interface {
	Handle(EventFilter, Handler)
	Start() error
	Stop() error
}
