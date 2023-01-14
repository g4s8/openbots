package state

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/g4s8/openbots/pkg/types"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

var _ types.StateProvider = (*DB)(nil)

type DB struct {
	con   *sql.DB
	botID int64
}

func NewDB(con *sql.DB, botID int64) *DB {
	return &DB{con: con, botID: botID}
}

func (db *DB) transactional(ctx context.Context, fn func(tx *sql.Tx) error) (errOut error) {
	tx, err := db.con.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	defer func() {
		if errOut != nil {
			if err := tx.Rollback(); err != nil {
				errOut = multierr.Append(errOut, errors.Wrap(err, "rollback tx"))
				return
			}
		} else {
			if err := tx.Commit(); err != nil {
				errOut = errors.Wrap(err, "commit tx")
				return
			}
		}
	}()
	errOut = fn(tx)
	return
}

func (db *DB) Load(ctx context.Context, uid types.ChatID, state types.State) error {
	return db.transactional(ctx, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx,
			`SELECT key, value FROM bot_state WHERE bot_id = $1 AND chat_id = $2`,
			db.botID, int64(uid))
		if err != nil {
			return errors.Wrap(err, "query state")
		}
		defer rows.Close()
		data := make(map[string]string)
		for rows.Next() {
			var key, value string
			if err = rows.Scan(&key, &value); err != nil {
				return errors.Wrap(err, "scan state")
			}
			data[key] = value
		}
		state.Fill(data)
		return nil
	})
}

func (db *DB) Update(ctx context.Context, uid types.ChatID, state types.State) error {
	return db.transactional(ctx, func(tx *sql.Tx) error {
		changes := state.Changes()
		if len(changes.Removed) > 0 {
			stmt, err := tx.PrepareContext(ctx,
				`DELETE FROM bot_state WHERE bot_id = $1 AND chat_id = $2 AND key = $3`)
			if err != nil {
				return errors.Wrap(err, "prepare delete")
			}
			defer stmt.Close()
			for _, key := range changes.Removed {
				fmt.Printf("DB(bot=%d chat=%d): delete %s\n", db.botID, uid, key)
				if _, err = stmt.ExecContext(ctx, db.botID, uid, key); err != nil {
					return errors.Wrap(err, "exec delete")
				}
			}
		}
		if len(changes.Added) > 0 {
			stmt, err := tx.PrepareContext(ctx,
				`INSERT INTO bot_state(bot_id, chat_id, key, value) VALUES($1, $2, $3, $4)
			ON CONFLICT(bot_id, chat_id, key) DO UPDATE SET value = $4`)
			if err != nil {
				return errors.Wrap(err, "prepare insert")
			}
			defer stmt.Close()
			for _, key := range changes.Added {
				val, _ := state.Get(key)
				fmt.Printf("DB(bot=%d chat=%d): add %s=%s\n", db.botID, uid, key, val)
				if _, err = stmt.ExecContext(ctx, db.botID, int64(uid), key, val); err != nil {
					return errors.Wrap(err, "exec insert")
				}
			}
		}
		return nil
	})
}
