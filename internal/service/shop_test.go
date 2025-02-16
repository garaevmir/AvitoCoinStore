package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/tests/mock_test/mocks"
)

func TestShopService_BuyItem(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	txRepo := new(mocks.TransactionRepositoryMock)
	invRepo := new(mocks.InventoryRepositoryMock)
	txMock := new(mocks.TxMock)

	txMock.On("Rollback", mock.Anything).Return(nil)

	shopSvc := NewShopService(userRepo, txRepo, invRepo)

	t.Run("Successful purchase", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1", Coins: 500}, nil).Once()
		userRepo.On("BeginTx", mock.Anything).
			Return(txMock, nil).Once()
		userRepo.On("UpdateUserCoinsTx", mock.Anything, txMock, "user1", -300).
			Return(nil).Once()
		invRepo.On("AddToInventoryTx", mock.Anything, txMock, "user1", "hoody", 1).
			Return(nil).Once()
		txMock.On("Commit", mock.Anything).
			Return(nil).Once()

		err := shopSvc.BuyItem(context.Background(), "user1", "hoody")
		assert.NoError(t, err)
	})

	t.Run("Item not found", func(t *testing.T) {
		err := shopSvc.BuyItem(context.Background(), "user1", "unknown_item")
		assert.ErrorIs(t, err, model.ErrItemNotFound)
	})

	t.Run("Database error during balance check", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user3").
			Return((*model.User)(nil), model.ErrInternalError).Once()

		err := shopSvc.BuyItem(context.Background(), "user3", "hoody")
		assert.ErrorIs(t, err, model.ErrInternalError)
		userRepo.AssertExpectations(t)
	})

	t.Run("Insufficient funds", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user2").
			Return(&model.User{ID: "user2", Coins: 100}, nil).Once()

		err := shopSvc.BuyItem(context.Background(), "user2", "hoody")
		assert.ErrorIs(t, err, model.ErrInsufficientFunds)
		userRepo.AssertExpectations(t)
	})

	t.Run("Begin transaction error", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1", Coins: 500}, nil).Once()
		userRepo.On("BeginTx", mock.Anything).
			Return(txMock, model.ErrInternalError).Once()

		err := shopSvc.BuyItem(context.Background(), "user1", "hoody")
		assert.ErrorIs(t, err, model.ErrInternalError)
	})

	t.Run("Update user coins error", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1", Coins: 500}, nil).Once()
		userRepo.On("BeginTx", mock.Anything).
			Return(txMock, nil).Once()
		userRepo.On("UpdateUserCoinsTx", mock.Anything, txMock, "user1", -300).
			Return(model.ErrInternalError).Once()

		err := shopSvc.BuyItem(context.Background(), "user1", "hoody")
		assert.ErrorIs(t, err, model.ErrInternalError)
	})

	t.Run("Adding to inventory error", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1", Coins: 500}, nil).Once()
		userRepo.On("BeginTx", mock.Anything).
			Return(txMock, nil).Once()
		userRepo.On("UpdateUserCoinsTx", mock.Anything, txMock, "user1", -300).
			Return(nil).Once()
		invRepo.On("AddToInventoryTx", mock.Anything, txMock, "user1", "hoody", 1).
			Return(model.ErrInternalError).Once()

		err := shopSvc.BuyItem(context.Background(), "user1", "hoody")
		assert.ErrorIs(t, err, model.ErrInternalError)
	})

	t.Run("Transaction commit error", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1", Coins: 500}, nil).Once()
		userRepo.On("BeginTx", mock.Anything).
			Return(txMock, nil).Once()
		userRepo.On("UpdateUserCoinsTx", mock.Anything, txMock, "user1", -300).
			Return(nil).Once()
		invRepo.On("AddToInventoryTx", mock.Anything, txMock, "user1", "hoody", 1).
			Return(nil).Once()
		txMock.On("Commit", mock.Anything).
			Return(model.ErrInternalError).Once()

		err := shopSvc.BuyItem(context.Background(), "user1", "hoody")
		assert.ErrorIs(t, err, model.ErrInternalError)
	})
}
