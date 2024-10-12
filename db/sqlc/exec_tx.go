package db

import (
	"context"
	"fmt"
)

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// starting transaction
	tx, err := store.connPool.Begin(ctx)

	// if starting is unsuccessful close transaction
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
