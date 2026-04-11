package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager defines transaction management operations.
type TxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type pgxTxManager struct {
	pool *pgxpool.Pool
}

// NewTxManager creates a new TxManager.
func NewTxManager(pool *pgxpool.Pool) TxManager {
	return &pgxTxManager{pool: pool}
}

// WithTransaction runs fn inside a transaction, reusing an existing one if present in ctx.
// Retries on deadlock up to 3 times with backoff.
func (m *pgxTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ExtractPgxTx(ctx); ok {
		return fn(ctx)
	}

	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := m.runInTx(ctx, fn)
		if err == nil {
			return nil
		}
		if isDeadlock(err) && attempt < maxRetries {
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}
		return err
	}
	return nil
}

func (m *pgxTxManager) runInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	txCtx := InjectPgxTx(ctx, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func isDeadlock(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgerrcode.DeadlockDetected
	}
	return false
}
