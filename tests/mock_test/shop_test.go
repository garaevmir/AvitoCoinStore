package mock_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garaevmir/avitocoinstore/internal/handler"
	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/service"
	"github.com/garaevmir/avitocoinstore/tests/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShopHandler_BuyItem(t *testing.T) {
	e := echo.New()
	userRepo := new(mocks.UserRepositoryMock)
	txRepo := new(mocks.TransactionRepositoryMock)
	invRepo := new(mocks.InventoryRepositoryMock)
	txMock := new(mocks.TxMock)
	shopService := service.NewShopService(userRepo, txRepo, invRepo)
	shopHandler := handler.NewShopHandler(shopService)

	txMock.On("Commit", mock.Anything).Return(nil)
	txMock.On("Rollback", mock.Anything).Return(nil)
	txMock.On("Begin", mock.Anything).Return(txMock, nil)

	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user_id", "user1")
			return next(c)
		}
	}

	t.Run("Successful item purchase", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1", Coins: 500}, nil).Once()

		invRepo.On("AddToInventoryTx", mock.Anything, mock.Anything, "user1", "hoody", 1).
			Return(nil).Once()

		userRepo.On("UpdateUserCoinsTx", mock.Anything, mock.Anything, "user1", -300).
			Return(nil).Once()
		userRepo.On("BeginTx", mock.Anything).
			Return(txMock, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/buy/hoody", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/buy/:item")
		c.SetParamNames("item")
		c.SetParamValues("hoody")

		err := middleware(shopHandler.BuyItem)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		userRepo.AssertExpectations(t)
		invRepo.AssertExpectations(t)
	})

	t.Run("Error: insufficient funds", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1", Coins: 100}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/buy/hoody", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/buy/:item")
		c.SetParamNames("item")
		c.SetParamValues("hoody")

		err := middleware(shopHandler.BuyItem)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorResp model.ErrorResponse
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "insufficient coins", errorResp.Errors)
	})

	t.Run("Error: item not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/buy/unknown_item", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/buy/:item")
		c.SetParamNames("item")
		c.SetParamValues("unknown_item")

		err := middleware(shopHandler.BuyItem)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorResp model.ErrorResponse
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "item not found", errorResp.Errors)
	})

	t.Run("Error: database error during balance update", func(t *testing.T) {
		userRepo.On("GetUserByID", mock.Anything, "user1").
			Return(&model.User{ID: "user1", Coins: 500}, nil).Once()

		userRepo.On("UpdateUserCoinsTx", mock.Anything, mock.Anything, "user1", -300).
			Return(errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/buy/hoody", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/buy/:item")
		c.SetParamNames("item")
		c.SetParamValues("hoody")

		err := middleware(shopHandler.BuyItem)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
