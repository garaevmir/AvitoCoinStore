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

func (m *TxMock) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *TxMock) Commit(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *TxMock) Rollback(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *TxMock) CopyFrom(ctx context.Context, tableName pgx.Identifier, columns []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columns, rowSrc)
	return args.Get(0).(int64), args.Error(1)
}

func (m *TxMock) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return m.Called(ctx, b).Get(0).(pgx.BatchResults)
}

func (m *TxMock) LargeObjects() pgx.LargeObjects {
	return m.Called().Get(0).(pgx.LargeObjects)
}

func (m *TxMock) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	return args.Get(0).(*pgconn.StatementDescription), args.Error(1)
}

func (m *TxMock) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgconn.CommandTag), mockArgs.Error(1)
}

func (m *TxMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgx.Rows), mockArgs.Error(1)
}

func (m *TxMock) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return m.Called(ctx, sql, args).Get(0).(pgx.Row)
}

func (m *TxMock) Conn() *pgx.Conn {
	return m.Called().Get(0).(*pgx.Conn)
}

type PgxRowMock struct {
	mock.Mock
}

func (r *PgxRowMock) Scan(dest ...any) error {
	return r.Called(dest...).Error(0)
}

type DBMock struct {
	mock.Mock
}

func (m *DBMock) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return m.Called(ctx, sql, args).Get(0).(pgx.Row)
}

func (m *DBMock) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	args := m.Called(ctx, txOptions)

	var tx pgx.Tx
	if args.Get(0) != nil {
		tx = args.Get(0).(pgx.Tx)
	}

	return tx, args.Error(1)
}

func (m *DBMock) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgx.Rows), mockArgs.Error(1)
}

type PgxRowsMock struct {
	mock.Mock
}

func (m *PgxRowsMock) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *PgxRowsMock) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *PgxRowsMock) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *PgxRowsMock) Close() {
	m.Called()
}

func (m *PgxRowsMock) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

func (m *PgxRowsMock) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	return args.Get(0).([]pgconn.FieldDescription)
}

func (m *PgxRowsMock) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *PgxRowsMock) Conn() *pgx.Conn {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgx.Conn)
}

func (m *PgxRowsMock) RawValues() [][]byte {
	args := m.Called()
	return args.Get(0).([][]byte)
}

type BatchResultsMock struct {
	mock.Mock
}

func (m *BatchResultsMock) Exec() (pgconn.CommandTag, error) {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *BatchResultsMock) Query() (pgx.Rows, error) {
	args := m.Called()
	return args.Get(0).(pgx.Rows), args.Error(1)
}

func (m *BatchResultsMock) QueryRow() pgx.Row {
	args := m.Called()
	return args.Get(0).(pgx.Row)
}

func (m *BatchResultsMock) Close() error {
	return m.Called().Error(0)
}
