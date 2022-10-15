package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type TxOperation func(*sql.Tx) error

func Transactional(db *sql.DB, ctx context.Context, opts *sql.TxOptions, op TxOperation) (err error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		err = errors.Wrap(err, "begin transaction")
		return
	}
	defer func() {
		if err != nil {
			if terr := tx.Rollback(); err != nil {
				err = multierr.Append(err,
					errors.Wrap(terr, "rollback transaction"))
			}
		} else {
			if terr := tx.Commit(); terr != nil {
				err = errors.Wrap(terr, "commit transaction")
			}
		}
	}()
	if err = op(tx); err != nil {
		err = errors.Wrap(err, "operation failed")
		return
	}
	return
}
