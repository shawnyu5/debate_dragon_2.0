package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute DB queries and transactions
type Store struct {
	*Queries // Embeds all sqlc generated queries automatically
	dbPool   *pgxpool.Pool
}

// NewStore creates a new store wrapper
func NewStore(dbPool *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(dbPool), // Initialize sqlc generated queries
		dbPool:  dbPool,
	}
}

// ExecTx executes a function within a database transaction
func (store *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Create an sqlc query instance bound to this specific transaction
	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
