package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Interface over pgxpool.Pool, needed for tests
type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// The wrapper over pgxpool.Pool, needed for tests
type PgxPoolWrapper struct {
	pool *pgxpool.Pool
}

// The wrapper over pgxpool.Pool.QueryRow
func (w *PgxPoolWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return w.pool.QueryRow(ctx, sql, args...)
}

// The wrapper over pgxpool.Pool.BeginTx
func (w *PgxPoolWrapper) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return w.pool.BeginTx(ctx, txOptions)
}

// The wrapper over pgxpool.Pool.Query
func (w *PgxPoolWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return w.pool.Query(ctx, sql, args...)
}
