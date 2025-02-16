package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/service"
)

type ShopHandler struct {
	shopService *service.ShopService
}

func NewShopHandler(s *service.ShopService) *ShopHandler {
	return &ShopHandler{shopService: s}
}

func (h *ShopHandler) BuyItem(c echo.Context) error {
	itemName := c.Param("item")
	userID := c.Get("user_id").(string)

	if err := h.shopService.BuyItem(c.Request().Context(), userID, itemName); err != nil {
		switch err {
		case model.ErrItemNotFound:
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Errors: model.ErrItemNotFound.Error()})
		case model.ErrInsufficientFunds:
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Errors: model.ErrInsufficientFunds.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: model.ErrInternalError.Error()})
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}
