package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"

	"github.com/garaevmir/avitocoinstore/internal/model"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)

	var user *model.User
	if args.Get(0) != nil {
		user = args.Get(0).(*model.User)
	}

	return user, args.Error(1)
}

func (m *UserRepositoryMock) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepositoryMock) UpdateUserCoinsTx(ctx context.Context, tx pgx.Tx, userID string, delta int) error {
	args := m.Called(ctx, tx, userID, delta)
	return args.Error(0)
}

func (m *UserRepositoryMock) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

type TransactionRepositoryMock struct {
	mock.Mock
}

func (m *TransactionRepositoryMock) TransferCoins(ctx context.Context, fromUserID, toUserID string, amount int) error {
	args := m.Called(ctx, fromUserID, toUserID, amount)
	return args.Error(0)
}

func (m *TransactionRepositoryMock) GetTransactionHistory(ctx context.Context, userID string) (*model.TransactionHistory, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*model.TransactionHistory), args.Error(1)
}

type InventoryRepositoryMock struct {
	mock.Mock
}

func (m *InventoryRepositoryMock) GetUserInventory(ctx context.Context, userID string) ([]model.InventoryItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.InventoryItem), args.Error(1)
}

func (m *InventoryRepositoryMock) AddToInventoryTx(ctx context.Context, tx pgx.Tx, userID, item string, quantity int) error {
	args := m.Called(ctx, tx, userID, item, quantity)
	return args.Error(0)
}
