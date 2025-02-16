package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/tests/mocks"
)

func TestInventoryRepository_GetUserInventory(t *testing.T) {
	dbMock := new(mocks.DBMock)
	repo := NewInventoryRepository(dbMock)
	ctx := context.Background()

	t.Run("Successful inventory retrieval", func(t *testing.T) {
		rowsMock := new(mocks.PgxRowsMock)
		expectedItems := []model.InventoryItem{
			{Name: "item1", Quantity: 5},
			{Name: "item2", Quantity: 3},
		}

		dbMock.On("Query", ctx, mock.Anything,
			[]interface{}{"user1"},
		).
			Return(rowsMock, nil).Once()

		rowsMock.On("Next").Return(true).Twice()
		rowsMock.On("Next").Return(false).Once()

		rowsMock.On("Scan", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				if len(expectedItems) > 0 {
					*args[0].(*string) = expectedItems[0].Name
					*args[1].(*int) = expectedItems[0].Quantity
					expectedItems = expectedItems[1:]
				}
			}).Return(nil).Twice()

		rowsMock.On("Close").Return(nil).Once()
		rowsMock.On("Err").Return(nil).Once()

		items, err := repo.GetUserInventory(ctx, "user1")
		assert.NoError(t, err)
		assert.Len(t, items, 2)
	})

	t.Run("Query execution error", func(t *testing.T) {
		expectedErr := errors.New("query error")
		rowsMock := new(mocks.PgxRowsMock)

		dbMock.On("Query", ctx, mock.Anything, mock.Anything).
			Return(rowsMock, expectedErr).Once()

		rowsMock.On("Close").Return().Once()
		rowsMock.On("Err").Return(expectedErr).Once()

		items, err := repo.GetUserInventory(ctx, "user2")
		assert.Nil(t, items)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("Row scanning error", func(t *testing.T) {
		rowsMock := new(mocks.PgxRowsMock)

		dbMock.On("Query", ctx, mock.Anything, mock.Anything).
			Return(rowsMock, nil).Once()

		rowsMock.On("Next").Return(true).Once()

		rowsMock.On("Scan", mock.Anything, mock.Anything).
			Return(errors.New("scan error")).Once()

		rowsMock.On("Close").Return(nil).Once()

		items, err := repo.GetUserInventory(ctx, "user3")
		assert.Nil(t, items)
		assert.ErrorContains(t, err, "scan error")
	})
}

func TestInventoryRepository_AddToInventoryTx(t *testing.T) {
	txMock := new(mocks.TxMock)
	repo := NewInventoryRepository(nil)
	commandTag := new(pgconn.CommandTag)
	ctx := context.Background()

	t.Run("Successful item addition", func(t *testing.T) {
		txMock.On("Exec", ctx, mock.Anything,
			[]interface{}{"user1", "sword", 1},
		).
			Return(*commandTag, nil).Once()

		err := repo.AddToInventoryTx(ctx, txMock, "user1", "sword", 1)
		assert.NoError(t, err)
		txMock.AssertExpectations(t)
	})

	t.Run("Query execution error", func(t *testing.T) {
		expectedErr := errors.New("exec error")

		txMock.On("Exec", ctx, mock.Anything, mock.Anything).
			Return(*commandTag, expectedErr).Once()

		err := repo.AddToInventoryTx(ctx, txMock, "user2", "shield", 1)
		assert.ErrorIs(t, err, expectedErr)
	})
}
