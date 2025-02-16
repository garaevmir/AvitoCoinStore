package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPoolWrapper struct {
	pool *pgxpool.Pool
}

func (w *PgxPoolWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return w.pool.QueryRow(ctx, sql, args...)
}

func (w *PgxPoolWrapper) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return w.pool.BeginTx(ctx, txOptions)
}

func (w *PgxPoolWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return w.pool.Query(ctx, sql, args...)
}
