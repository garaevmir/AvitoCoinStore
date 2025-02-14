package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type TxMock struct {
	mock.Mock
}

func (m TxMock) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m TxMock) Commit(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m TxMock) Rollback(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m TxMock) CopyFrom(ctx context.Context, tableName pgx.Identifier, columns []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columns, rowSrc)
	return args.Get(0).(int64), args.Error(1)
}

func (m TxMock) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return m.Called(ctx, b).Get(0).(pgx.BatchResults)
}

func (m TxMock) LargeObjects() pgx.LargeObjects {
	return m.Called().Get(0).(pgx.LargeObjects)
}

func (m TxMock) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	return args.Get(0).(*pgconn.StatementDescription), args.Error(1)
}

func (m TxMock) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgconn.CommandTag), mockArgs.Error(1)
}

func (m TxMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgx.Rows), mockArgs.Error(1)
}

func (m TxMock) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return m.Called(ctx, sql, args).Get(0).(pgx.Row)
}

func (m TxMock) Conn() *pgx.Conn {
	return m.Called().Get(0).(*pgx.Conn)
}

type PgxRowMock struct {
	mock.Mock
}

func (r *PgxRowMock) Scan(dest ...any) error {
	return r.Called(dest...).Error(0)
}
