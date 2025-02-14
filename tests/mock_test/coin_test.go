package mock_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garaevmir/avitocoinstore/internal/handler"
	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/tests/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCoinHandler_SendCoins(t *testing.T) {
	e := echo.New()
	userRepo := new(mocks.UserRepositoryMock)
	txRepo := new(mocks.TransactionRepositoryMock)
	coinHandler := handler.NewCoinHandler(txRepo, userRepo)

	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user_id", "user1")
			return next(c)
		}
	}

	t.Run("Successful coin transfer", func(t *testing.T) {
		reqBody := model.SendCoinRequest{
			ToUser: "user2",
			Amount: 100,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		userRepo.On("GetUserByUsername", mock.Anything, "user2").
			Return(&model.User{ID: "user2"}, nil).Once()

		txRepo.On("TransferCoins", mock.Anything, "user1", "user2", 100).
			Return(nil).Once()

		err := middleware(coinHandler.SendCoins)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		txRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("Error: insufficient funds", func(t *testing.T) {
		reqBody := model.SendCoinRequest{
			ToUser: "user2",
			Amount: 1000,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		userRepo.On("GetUserByUsername", mock.Anything, "user2").
			Return(&model.User{ID: "user2"}, nil).Once()

		txRepo.On("TransferCoins", mock.Anything, "user1", "user2", 1000).
			Return(model.ErrInsufficientFunds).Once()

		err := middleware(coinHandler.SendCoins)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorResp model.ErrorResponse
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "insufficient funds", errorResp.Errors)
	})

	t.Run("Error: recipient not found", func(t *testing.T) {
		reqBody := model.SendCoinRequest{
			ToUser: "unknown_user",
			Amount: 100,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		userRepo.On("GetUserByUsername", mock.Anything, "unknown_user").
			Return((*model.User)(nil), nil).Once()

		err := middleware(coinHandler.SendCoins)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var errorResp model.ErrorResponse
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "user not found", errorResp.Errors)
	})

	t.Run("Error: negative coin amount", func(t *testing.T) {
		reqBody := model.SendCoinRequest{
			ToUser: "user2",
			Amount: -100,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set("user_id", "user1")
				return next(c)
			}
		}

		err := middleware(coinHandler.SendCoins)(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorResp model.ErrorResponse
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "amount must be positive", errorResp.Errors)
	})
}
