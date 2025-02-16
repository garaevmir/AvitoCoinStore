package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/repository"
)

// A structure for info handler
type InfoHandler struct {
	userRepo        repository.UserRepositoryInt
	inventoryRepo   repository.InventoryRepositoryInt
	transactionRepo repository.TransactionRepositoryInt
}

// Constructor for info handler
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

// Function for /api/info request
func (h *InfoHandler) GetUserInfo(c echo.Context) error {
	userID := c.Get("user_id").(string)

	user, err := h.userRepo.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, model.ErrUserNotFound)
		}
		return c.JSON(http.StatusInternalServerError, model.ErrInternalError)
	}

	inventory, err := h.inventoryRepo.GetUserInventory(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrInventory)
	}

	history, err := h.transactionRepo.GetTransactionHistory(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrHistory)
	}

	return c.JSON(http.StatusOK, model.InfoResponse{
		Coins:       user.Coins,
		Inventory:   inventory,
		CoinHistory: *history,
	})
}
