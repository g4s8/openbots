package context

import (
	ctx "context"
	"database/sql"

	"github.com/g4s8/openbots/internal/db"
	"github.com/g4s8/openbots/pkg/types"
	"github.com/pkg/errors"
)

type dbProvider struct {
	db    *sql.DB
	botID int64
}

func NewDBProvider(db *sql.DB, botID int64) types.ContextProvider {
	return &dbProvider{db: db, botID: botID}
}

func (p *dbProvider) UserContext(chatID types.ChatID) types.Context {
	return &dbContext{db: p.db, botID: p.botID, chatID: chatID}
}

type dbContext struct {
	db     *sql.DB
	botID  int64
	chatID types.ChatID
}

func (c *dbContext) Set(ctx ctx.Context, value string) error {
	return db.Transactional(c.db, ctx, nil, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO bot_context (bot_id, chat_id, value) VALUES ($1, $2, $3)
			ON CONFLICT (bot_id, chat_id) DO UPDATE SET value = $3`,
			c.botID, c.chatID, value); err != nil {
			return errors.Wrap(err, "insert value")
		}
		return nil
	})
}

func (c *dbContext) Reset(ctx ctx.Context) error {
	return db.Transactional(c.db, ctx, nil, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM bot_context WHERE bot_id = $1 AND chat_id = $2`,
			c.botID, c.chatID); err != nil {
			return errors.Wrap(err, "delete value")
		}
		return nil
	})
}

func (c *dbContext) Check(ctx ctx.Context, value string) (bool, error) {
	var result string
	if err := db.Transactional(c.db, ctx, nil, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx,
			`SELECT value FROM bot_context WHERE bot_id = $1 AND chat_id = $2`,
			c.botID, c.chatID)
		if err != nil {
			return errors.Wrap(err, "select value")
		}
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(&result); err != nil {
				return errors.Wrap(err, "scan value")
			}
		}
		return nil
	}); err != nil {
		return false, err
	}
	return result == value, nil
}
