package mock_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/garaevmir/avitocoinstore/internal/handler"
	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/tests/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInfoHandler_GetUserInfo(t *testing.T) {
	e := echo.New()
	userRepo := new(mocks.UserRepositoryMock)
	invRepo := new(mocks.InventoryRepositoryMock)
	txRepo := new(mocks.TransactionRepositoryMock)
	infoHandler := handler.NewInfoHandler(userRepo, invRepo, txRepo)

	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user_id", "user1")
			return next(c)
		}
	}

	t.Run("Successful information retrieval", func(t *testing.T) {
		mockUser := &model.User{
			ID:    "user1",
			Coins: 1500,
		}

		mockInventory := []model.InventoryItem{
			{Name: "t-shirt", Quantity: 2},
			{Name: "cup", Quantity: 1},
		}

		mockHistory := &model.TransactionHistory{
			Received: []model.ReceivedTransaction{
				{FromUser: "user2", Amount: 200, Timestamp: time.Now()},
			},
			Sent: []model.SentTransaction{
				{ToUser: "user3", Amount: 100, Timestamp: time.Now()},
			},
		}

		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(mockUser, nil).Once()

		invRepo.On("GetUserInventory", mock.Anything, "user1").
			Return(mockInventory, nil).Once()

		txRepo.On("GetTransactionHistory", mock.Anything, "user1").
			Return(mockHistory, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := middleware(infoHandler.GetUserInfo)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response model.InfoResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, 1500, response.Coins)
		assert.Len(t, response.Inventory, 2)
		assert.Len(t, response.CoinHistory.Received, 1)
		userRepo.AssertExpectations(t)
		invRepo.AssertExpectations(t)
		txRepo.AssertExpectations(t)
	})

	t.Run("Error: user not found", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(nil, sql.ErrNoRows).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := middleware(infoHandler.GetUserInfo)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var errorResp model.ErrorResponse
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "user not found", errorResp.Errors)
	})

	t.Run("Error: database error", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(nil, errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := middleware(infoHandler.GetUserInfo)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Empty transaction history", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1"}, nil).Once()

		invRepo.On("GetUserInventory", mock.Anything, "user1").
			Return([]model.InventoryItem{}, nil).Once()

		txRepo.On("GetTransactionHistory", mock.Anything, "user1").
			Return(&model.TransactionHistory{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := middleware(infoHandler.GetUserInfo)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response model.InfoResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Empty(t, response.Inventory)
		assert.Empty(t, response.CoinHistory.Received)
	})
}
