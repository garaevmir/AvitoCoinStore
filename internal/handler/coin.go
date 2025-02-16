package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/repository"
)

// A structure for a send coin handler
type CoinHandler struct {
	transactionRepo repository.TransactionRepositoryInt
	userRepo        repository.UserRepositoryInt
}

// Constructor for send coin handler
func NewCoinHandler(tRepo repository.TransactionRepositoryInt, uRepo repository.UserRepositoryInt) *CoinHandler {
	return &CoinHandler{transactionRepo: tRepo, userRepo: uRepo}
}

// Function for /api/sendCoin request
func (h *CoinHandler) SendCoins(c echo.Context) error {
	var req model.SendCoinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrInvalidRequest)
	}

	if req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, model.ErrNegAmount)
	}

	if req.ToUser == "" {
		return c.JSON(http.StatusBadRequest, model.ErrUserNotFound)
	}

	fromUserID := c.Get("user_id").(string)
	toUser, err := h.userRepo.GetUserByUsername(c.Request().Context(), req.ToUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrUserNotFound)
	}

	if toUser == nil {
		return c.JSON(http.StatusNotFound, model.ErrUserNotFound)
	}

	if err := h.transactionRepo.TransferCoins(c.Request().Context(), fromUserID, toUser.ID, req.Amount); err != nil {
		switch err {
		case model.ErrInsufficientFunds:
			return c.JSON(http.StatusBadRequest, model.ErrInsufficientFunds)
		default:
			return c.JSON(http.StatusInternalServerError, model.ErrInternalError)
		}
	}

	return c.JSON(200, map[string]interface{}{"status": "success"})
}
