package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/tests/mocks"
)

func TestTransactionRepository_TransferCoins(t *testing.T) {
	poolMock := new(mocks.DBMock)
	txMock := new(mocks.TxMock)
	batchResultsMock := new(mocks.BatchResultsMock)
	repo := NewTransactionRepository(poolMock)
	commandTag := new(pgconn.CommandTag)
	rowMock := new(mocks.PgxRowMock)
	txOptions := pgx.TxOptions{}
	ctx := context.Background()

	txMock.On("Rollback", mock.Anything).Return(nil)

	t.Run("Successful transfer", func(t *testing.T) {
		poolMock.On("BeginTx", ctx, txOptions).Return(txMock, nil).Once()

		txMock.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
			*args[0].(*int) = 1000
		}).Return(nil).Once()

		txMock.On("SendBatch", ctx, mock.Anything).
			Return(batchResultsMock).Once()

		batchResultsMock.On("Exec").Return(*commandTag, nil).Times(3)
		batchResultsMock.On("Close").Return(nil).Once()

		txMock.On("Commit", ctx).Return(nil).Once()

		err := repo.TransferCoins(ctx, "user1", "user2", 500)
		assert.NoError(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("Transaction start error", func(t *testing.T) {
		poolMock.On("BeginTx", ctx, txOptions).
			Return(txMock, pgx.ErrTxClosed).Once()

		err := repo.TransferCoins(ctx, "user1", "user2", 500)
		assert.ErrorIs(t, err, pgx.ErrTxClosed)
	})

	t.Run("Balance retrieval error", func(t *testing.T) {
		poolMock.On("BeginTx", ctx, txOptions).Return(txMock, nil).Once()

		txMock.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Return(model.ErrInternalError).Once()

		err := repo.TransferCoins(ctx, "user1", "user2", 500)
		assert.ErrorIs(t, err, model.ErrInternalError)
	})

	t.Run("Insufficient funds", func(t *testing.T) {
		poolMock.On("BeginTx", ctx, txOptions).Return(txMock, nil).Once()

		txMock.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
			*args[0].(*int) = 200
		}).Return(nil).Once()

		err := repo.TransferCoins(ctx, "user1", "user2", 500)
		assert.ErrorIs(t, err, model.ErrInsufficientFunds)
	})

	t.Run("Request sending error", func(t *testing.T) {
		poolMock.On("BeginTx", ctx, txOptions).Return(txMock, nil).Once()

		txMock.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Run(func(args mock.Arguments) {
				*args[0].(*int) = 1000
			}).Return(nil).Once()

		txMock.On("SendBatch", ctx, mock.Anything).Return(batchResultsMock).Once()

		batchResultsMock.On("Exec").Return(*commandTag, nil).Twice()
		batchResultsMock.On("Exec").Return(*commandTag, pgx.ErrTxClosed).Once()

		err := repo.TransferCoins(ctx, "user1", "user2", 500)
		assert.ErrorIs(t, err, pgx.ErrTxClosed)
	})

	t.Run("Transaction commit error", func(t *testing.T) {
		poolMock.On("BeginTx", ctx, txOptions).Return(txMock, nil).Once()

		txMock.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Run(func(args mock.Arguments) {
				*args[0].(*int) = 1000
			}).Return(nil).Once()

		txMock.On("SendBatch", ctx, mock.Anything).Return(batchResultsMock).Once()

		batchResultsMock.On("Exec").Return(*commandTag, nil).Times(3)
		batchResultsMock.On("Close").Return(nil).Once()

		txMock.On("Commit", ctx).Return(pgx.ErrTxClosed).Once()

		err := repo.TransferCoins(ctx, "user1", "user2", 500)
		assert.ErrorIs(t, err, pgx.ErrTxClosed)
	})
}

func TestTransactionRepository_GetTransactionHistory(t *testing.T) {
	poolMock := new(mocks.DBMock)
	repo := NewTransactionRepository(poolMock)
	ctx := context.Background()

	recTrans := model.ReceivedTransaction{FromUser: "user2", Amount: 100, Timestamp: time.Time{}}

	senTrans := model.SentTransaction{ToUser: "user3", Amount: 200, Timestamp: time.Time{}}

	t.Run("Successful history retrieval", func(t *testing.T) {
		receivedRows := new(mocks.PgxRowsMock)

		poolMock.On("Query", ctx, mock.Anything, []interface{}{"user1"}).
			Return(receivedRows, nil).Once()

		receivedRows.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = recTrans.FromUser
				*args[1].(*int) = recTrans.Amount
				*args[2].(*time.Time) = recTrans.Timestamp
			}).Return(nil).Twice()

		receivedRows.On("Close").Return(nil).Once()
		receivedRows.On("Next").Return(true).Twice()
		receivedRows.On("Next").Return(false).Once()

		sentRows := new(mocks.PgxRowsMock)

		poolMock.On("Query", ctx, mock.Anything, []interface{}{"user1"}).Return(sentRows, nil).Once()

		sentRows.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = senTrans.ToUser
				*args[1].(*int) = senTrans.Amount
				*args[2].(*time.Time) = senTrans.Timestamp
			}).Return(nil).Once()

		sentRows.On("Close").Return(nil).Once()
		sentRows.On("Next").Return(true).Once()
		sentRows.On("Next").Return(false).Once()

		history, err := repo.GetTransactionHistory(ctx, "user1")
		assert.NoError(t, err)
		assert.Len(t, history.Received, 2)
		assert.Len(t, history.Sent, 1)
	})

	t.Run("First query execution error", func(t *testing.T) {
		receivedRows := new(mocks.PgxRowsMock)
		poolMock.On("Query", ctx, mock.Anything, []interface{}{"user1"}).
			Return(receivedRows, model.ErrInternalError).Once()

		_, err := repo.GetTransactionHistory(ctx, "user1")
		assert.ErrorIs(t, err, model.ErrInternalError)
	})

	t.Run("Error reading incoming transactions response", func(t *testing.T) {
		receivedRows := new(mocks.PgxRowsMock)

		poolMock.On("Query", ctx, mock.Anything, []interface{}{"user1"}).
			Return(receivedRows, nil).Once()

		receivedRows.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = recTrans.FromUser
				*args[1].(*int) = recTrans.Amount
				*args[2].(*time.Time) = recTrans.Timestamp
			}).Return(nil).Once()

		receivedRows.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time")).Run(func(args mock.Arguments) {
			*args[0].(*string) = recTrans.FromUser
			*args[1].(*int) = recTrans.Amount
			*args[2].(*time.Time) = recTrans.Timestamp
		}).Return(model.ErrInternalError).Once()

		receivedRows.On("Close").Return(nil).Once()
		receivedRows.On("Next").Return(true).Twice()
		receivedRows.On("Next").Return(false).Once()

		_, err := repo.GetTransactionHistory(ctx, "user1")
		assert.ErrorIs(t, err, model.ErrInternalError)
	})

	t.Run("Second query execution error", func(t *testing.T) {
		receivedRows := new(mocks.PgxRowsMock)

		poolMock.On("Query", ctx, mock.Anything, []interface{}{"user1"}).Return(receivedRows, nil).Once()

		receivedRows.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = recTrans.FromUser
				*args[1].(*int) = recTrans.Amount
				*args[2].(*time.Time) = recTrans.Timestamp
			}).Return(nil).Once()

		receivedRows.On("Close").Return(nil).Once()
		receivedRows.On("Next").Return(true).Once()
		receivedRows.On("Next").Return(false).Once()

		sentRows := new(mocks.PgxRowsMock)

		poolMock.On("Query", ctx, mock.Anything, []interface{}{"user1"}).
			Return(sentRows, model.ErrInternalError).Once()

		_, err := repo.GetTransactionHistory(ctx, "user1")
		assert.ErrorIs(t, err, model.ErrInternalError)
	})

	t.Run("Error reading outgoing transactions response", func(t *testing.T) {
		receivedRows := new(mocks.PgxRowsMock)

		poolMock.On("Query", ctx, mock.Anything, []interface{}{"user1"}).
			Return(receivedRows, nil).Once()

		receivedRows.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = recTrans.FromUser
				*args[1].(*int) = recTrans.Amount
				*args[2].(*time.Time) = recTrans.Timestamp
			}).Return(nil).Twice()

		receivedRows.On("Close").Return(nil).Once()
		receivedRows.On("Next").Return(true).Twice()
		receivedRows.On("Next").Return(false).Once()

		sentRows := new(mocks.PgxRowsMock)

		poolMock.On("Query", ctx, mock.Anything, []interface{}{"user1"}).
			Return(sentRows, nil).Once()

		sentRows.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = senTrans.ToUser
				*args[1].(*int) = senTrans.Amount
				*args[2].(*time.Time) = senTrans.Timestamp
			}).Return(model.ErrInternalError).Once()

		sentRows.On("Close").Return(nil).Once()
		sentRows.On("Next").Return(true).Once()

		_, err := repo.GetTransactionHistory(ctx, "user1")
		assert.ErrorIs(t, err, model.ErrInternalError)
	})
}
