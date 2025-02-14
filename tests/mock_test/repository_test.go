package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/service"
	"github.com/garaevmir/avitocoinstore/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShopService_BuyItem(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	txRepo := new(mocks.TransactionRepositoryMock)
	invRepo := new(mocks.InventoryRepositoryMock)
	txMock := new(mocks.TxMock)

	txMock.On("Commit", mock.Anything).Return(nil)
	txMock.On("Rollback", mock.Anything).Return(nil)

	shopSvc := service.NewShopService(userRepo, txRepo, invRepo)

	t.Run("Successful purchase", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1", Coins: 500}, nil).Once()
		userRepo.On("BeginTx", mock.Anything).Return(txMock, nil).Once()
		userRepo.On("UpdateUserCoinsTx", mock.Anything, txMock, "user1", -300).
			Return(nil).Once()
		invRepo.On("AddToInventoryTx", mock.Anything, txMock, "user1", "hoody", 1).
			Return(nil).Once()

		err := shopSvc.BuyItem(context.Background(), "user1", "hoody")
		assert.NoError(t, err)
	})

	t.Run("Insufficient funds", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user2").
			Return(&model.User{ID: "user2", Coins: 100}, nil).Once()

		err := shopSvc.BuyItem(context.Background(), "user2", "hoody")
		assert.ErrorIs(t, err, model.ErrInsufficientFunds)
		userRepo.AssertExpectations(t)
	})

	t.Run("Item not found", func(t *testing.T) {
		err := shopSvc.BuyItem(context.Background(), "user1", "unknown_item")
		assert.ErrorIs(t, err, model.ErrItemNotFound)
	})

	t.Run("Database error during balance check", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user3").
			Return((*model.User)(nil), errors.New("database error")).Once()

		err := shopSvc.BuyItem(context.Background(), "user3", "hoody")
		assert.ErrorContains(t, err, "database error")
		userRepo.AssertExpectations(t)
	})
}
