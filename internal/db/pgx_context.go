package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type contextKey string

const txContextKey contextKey = "pgx_tx"

// InjectPgxTx stores a transaction in the context.
func InjectPgxTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

// ExtractPgxTx retrieves a transaction from the context.
func ExtractPgxTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey).(pgx.Tx)
	return tx, ok
}
