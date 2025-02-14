package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/repository"
	"github.com/labstack/echo/v4"
)

type InfoHandler struct {
	userRepo        repository.UserRepositoryInt
	inventoryRepo   repository.InventoryRepositoryInt
	transactionRepo repository.TransactionRepositoryInt
}

func NewInfoHandler(
	uRepo repository.UserRepositoryInt,
	iRepo repository.InventoryRepositoryInt,
	tRepo repository.TransactionRepositoryInt,
) *InfoHandler {
	return &InfoHandler{
		userRepo:        uRepo,
		inventoryRepo:   iRepo,
		transactionRepo: tRepo,
	}
}

func (h *InfoHandler) GetUserInfo(c echo.Context) error {
	userID := c.Get("user_id").(string)

	user, err := h.userRepo.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{Errors: "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "database error"})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{Errors: "user not found"})
	}

	inventory, err := h.inventoryRepo.GetUserInventory(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "failed to get inventory"})
	}

	history, err := h.transactionRepo.GetTransactionHistory(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "failed to get history"})
	}

	return c.JSON(http.StatusOK, model.InfoResponse{
		Coins:       user.Coins,
		Inventory:   inventory,
		CoinHistory: *history,
	})
}
