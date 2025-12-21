package utils

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// WithTx executes a function within a database transaction.
// It automatically handles commit on success and rollback on error or panic.
//
// Example usage:
//
//	err := utils.WithTx(ctx, db, func(tx *sqlx.Tx) error {
//	    // Your transactional operations here
//	    return repo.CreateUser(ctx, tx, user)
//	})
func WithTx(
	ctx context.Context,
	db *sqlx.DB,
	fn func(tx *sqlx.Tx) error,
) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// Panic occurred, rollback and re-panic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// Error occurred, rollback
			_ = tx.Rollback()
		} else {
			// Success, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
