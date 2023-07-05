// Package types provides base public types used by bot.
package types

import (
	"context"
	"strconv"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ChatID is a chat identifier.
type ChatID int64

func (cid ChatID) String() string {
	return strconv.FormatInt(int64(cid), 10)
}

func (cid ChatID) Int64() int64 {
	return int64(cid)
}

// Handler of telegram update message.
type Handler interface {
	Handle(context.Context, *telegram.Update, *telegram.BotAPI) error
}

// EventFilter checks that telegram update could be handlerd.
type EventFilter interface {
	Check(context.Context, *telegram.Update) (bool, error)
}

type DataContainer struct {
	data any
}

func (c *DataContainer) Set(data any) {
	c.data = data
}

func (c *DataContainer) Get() any {
	return c.data
}

type DataLoader interface {
	Load(ctx context.Context, c *DataContainer, upd *telegram.Update) error
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
