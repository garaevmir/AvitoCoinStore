package handler

import (
	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/repository"
	"github.com/labstack/echo/v4"
)

type CoinHandler struct {
	transactionRepo repository.TransactionRepositoryInt
	userRepo        repository.UserRepositoryInt
}

func NewCoinHandler(tRepo repository.TransactionRepositoryInt, uRepo repository.UserRepositoryInt) *CoinHandler {
	return &CoinHandler{transactionRepo: tRepo, userRepo: uRepo}
}

func (h *CoinHandler) SendCoins(c echo.Context) error {
	var req model.SendCoinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, model.ErrorResponse{Errors: "invalid request"})
	}

	if req.Amount <= 0 {
		return c.JSON(400, model.ErrorResponse{Errors: "amount must be positive"})
	}

	fromUserID := c.Get("user_id").(string)
	toUser, err := h.userRepo.GetUserByUsername(c.Request().Context(), req.ToUser)
	if err != nil {
		return c.JSON(404, model.ErrorResponse{Errors: "user not found"})
	}

	if toUser == nil {
		return c.JSON(404, model.ErrorResponse{Errors: "user not found"})
	}

	if err := h.transactionRepo.TransferCoins(c.Request().Context(), fromUserID, toUser.ID, req.Amount); err != nil {
		return c.JSON(400, model.ErrorResponse{Errors: err.Error()})
	}

	return c.JSON(200, map[string]interface{}{"status": "success"})
}
