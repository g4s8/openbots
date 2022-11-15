package state

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/g4s8/openbots/pkg/types"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var _ types.StateProvider = (*DB)(nil)

type DB struct {
	con   *sql.DB
	botID int64
}

func NewDB(con *sql.DB, botID int64) *DB {
	return &DB{con: con, botID: botID}
}

func (db *DB) Load(ctx context.Context, uid types.ChatID, state types.State) error {
	var err error
	tx, err := db.con.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

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
}

func (db *DB) Update(ctx context.Context, uid types.ChatID, state types.State) error {
	var err error
	tx, err := db.con.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	reporter, ok := state.(reporter)
	if !ok {
		return errors.New("state doesn't support reporting")
	}

	changes := reporter.changes()
	if len(changes.deleted) > 0 {
		stmt, err := tx.PrepareContext(ctx,
			`DELETE FROM bot_state WHERE bot_id = $1 AND chat_id = $2 AND key = $3`)
		if err != nil {
			return errors.Wrap(err, "prepare delete")
		}
		defer stmt.Close()
		for _, key := range changes.deleted {
			fmt.Printf("DB(bot=%d chat=%d): delete %s\n", db.botID, uid, key)
			if _, err = stmt.ExecContext(ctx, db.botID, uid, key); err != nil {
				return errors.Wrap(err, "exec delete")
			}
		}
	}
	if len(changes.added) > 0 {
		stmt, err := tx.PrepareContext(ctx,
			`INSERT INTO bot_state(bot_id, chat_id, key, value) VALUES($1, $2, $3, $4)
			ON CONFLICT(bot_id, chat_id, key) DO UPDATE SET value = $4`)
		if err != nil {
			return errors.Wrap(err, "prepare insert")
		}
		defer stmt.Close()
		for _, key := range changes.added {
			val, _ := state.Get(key)
			fmt.Printf("DB(bot=%d chat=%d): add %s=%s\n", db.botID, uid, key, val)
			if _, err = stmt.ExecContext(ctx, db.botID, int64(uid), key, val); err != nil {
				return errors.Wrap(err, "exec insert")
			}
		}
	}
	return nil
}
