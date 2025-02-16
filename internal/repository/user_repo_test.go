package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/tests/mock_test/mocks"
)

func TestUserRepository_CreateUser(t *testing.T) {
	dbMock := new(mocks.DBMock)
	userRepo := NewUserRepository(dbMock)
	rowMock := new(mocks.PgxRowMock)
	ctx := context.Background()

	t.Run("Successful user creation", func(t *testing.T) {
		testUser := &model.User{
			Username:     "test_user",
			PasswordHash: "hash",
			Coins:        100,
		}

		dbMock.On("QueryRow", ctx, mock.Anything, []interface{}{"test_user", "hash", 100}).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = "generated-id-123"
			}).Return(nil).Once()

		err := userRepo.CreateUser(ctx, testUser)
		assert.NoError(t, err)
		assert.Equal(t, "generated-id-123", testUser.ID)
		dbMock.AssertExpectations(t)
		rowMock.AssertExpectations(t)
	})

	t.Run("Error inserting into database", func(t *testing.T) {
		testUser := &model.User{Username: "error_user"}

		dbMock.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Return(pgx.ErrNoRows).Once()

		err := userRepo.CreateUser(ctx, testUser)
		assert.ErrorIs(t, err, pgx.ErrNoRows)
		assert.Empty(t, testUser.ID)
	})
}

func TestUserRepository_GetUserByUsername(t *testing.T) {
	dbMock := new(mocks.DBMock)
	userRepo := NewUserRepository(dbMock)
	rowMock := new(mocks.PgxRowMock)
	ctx := context.Background()

	testUser := &model.User{
		ID:           "123",
		Username:     "test_user",
		PasswordHash: "hash",
		Coins:        100,
	}

	t.Run("Successful user retrieval", func(t *testing.T) {
		dbMock.On("QueryRow", ctx, mock.Anything, []interface{}{"test_user"}).
			Return(rowMock).Once()

		rowMock.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
		).Run(func(args mock.Arguments) {
			*args[0].(*string) = testUser.ID
			*args[1].(*string) = testUser.Username
			*args[2].(*string) = testUser.PasswordHash
			*args[3].(*int) = testUser.Coins
		}).Return(nil).Once()

		user, err := userRepo.GetUserByUsername(ctx, "test_user")
		assert.NoError(t, err)
		assert.Equal(t, testUser, user)
		dbMock.AssertExpectations(t)
		rowMock.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		dbMock.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
		).Run(func(args mock.Arguments) {
			*args[0].(*string) = testUser.ID
			*args[1].(*string) = testUser.Username
			*args[2].(*string) = testUser.PasswordHash
			*args[3].(*int) = testUser.Coins
		}).Return(pgx.ErrNoRows).Once()

		user, err := userRepo.GetUserByUsername(ctx, "unknown_user")
		assert.Nil(t, user)
		assert.NoError(t, err)
	})

	t.Run("Database error", func(t *testing.T) {
		expectedErr := pgx.ErrTooManyRows

		dbMock.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan",
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
		).Run(func(args mock.Arguments) {
			*args[0].(*string) = testUser.ID
			*args[1].(*string) = testUser.Username
			*args[2].(*string) = testUser.PasswordHash
			*args[3].(*int) = testUser.Coins
		}).Return(expectedErr).Once()

		user, err := userRepo.GetUserByUsername(ctx, "error_user")
		assert.Nil(t, user)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestUserRepository_GetUserByID(t *testing.T) {
	dbMock := new(mocks.DBMock)
	userRepo := NewUserRepository(dbMock)
	rowMock := new(mocks.PgxRowMock)
	ctx := context.Background()

	testUser := &model.User{
		ID:           "123",
		Username:     "test_user",
		PasswordHash: "hash",
		Coins:        100,
	}

	t.Run("User found by ID", func(t *testing.T) {
		dbMock.On("QueryRow", ctx, mock.Anything, []interface{}{"123"}).
			Return(rowMock).Once()

		rowMock.On("Scan",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Run(func(args mock.Arguments) {
			*args[0].(*string) = testUser.ID
			*args[1].(*string) = testUser.Username
			*args[2].(*string) = testUser.PasswordHash
			*args[3].(*int) = testUser.Coins
		}).Return(nil).Once()

		user, err := userRepo.GetUserByID(ctx, "123")
		assert.NoError(t, err)
		assert.Equal(t, testUser, user)
		dbMock.AssertExpectations(t)
		rowMock.AssertExpectations(t)
	})

	t.Run("User not found by ID", func(t *testing.T) {
		dbMock.On("QueryRow", ctx, mock.Anything, []interface{}{"123"}).
			Return(rowMock).Once()

		rowMock.On("Scan",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Run(func(args mock.Arguments) {
			*args[0].(*string) = testUser.ID
			*args[1].(*string) = testUser.Username
			*args[2].(*string) = testUser.PasswordHash
			*args[3].(*int) = testUser.Coins
		}).Return(model.ErrInternalError).Once()

		_, err := userRepo.GetUserByID(ctx, "123")
		assert.Error(t, err)
	})
}

func TestUserRepository_BeginTx(t *testing.T) {
	dbMock := new(mocks.DBMock)
	userRepo := NewUserRepository(dbMock)
	txMock := new(mocks.TxMock)
	txOptions := pgx.TxOptions{}
	ctx := context.Background()

	t.Run("Successful transaction start", func(t *testing.T) {
		dbMock.On("BeginTx", ctx, txOptions).
			Return(txMock, nil).Once()

		tx, err := userRepo.BeginTx(ctx)
		assert.NoError(t, err)
		assert.Equal(t, txMock, tx)
		dbMock.AssertExpectations(t)
	})

	t.Run("Error starting transaction", func(t *testing.T) {
		expectedErr := pgx.ErrTxClosed

		dbMock.On("BeginTx", ctx, txOptions).
			Return(nil, expectedErr).Once()

		tx, err := userRepo.BeginTx(ctx)
		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, tx)
		dbMock.AssertExpectations(t)
	})
}

func TestUserRepository_UpdateUserCoinsTx(t *testing.T) {
	poolMock := &pgxpool.Pool{}
	userRepoMock := NewUserRepository(poolMock)
	txMock := new(mocks.TxMock)
	commandTag := new(pgconn.CommandTag)
	ctx := context.Background()

	t.Run("Success coins update", func(t *testing.T) {
		txMock.On("Exec", ctx, mock.Anything, []interface{}{50, "user1"}).
			Return(*commandTag, nil).Once()

		err := userRepoMock.UpdateUserCoinsTx(ctx, txMock, "user1", 50)
		assert.NoError(t, err)
	})

	t.Run("Database error on update", func(t *testing.T) {
		txMock.On("Exec", ctx, mock.Anything, []interface{}{50, "user1"}).
			Return(*commandTag, model.ErrInternalError).Once()

		err := userRepoMock.UpdateUserCoinsTx(ctx, txMock, "user1", 50)
		assert.Error(t, err)
	})
}
