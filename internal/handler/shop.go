package handler

import (
	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/service"
	"github.com/labstack/echo/v4"
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
			return c.JSON(400, model.ErrorResponse{Errors: "item not found"})
		case model.ErrInsufficientFunds:
			return c.JSON(400, model.ErrorResponse{Errors: "insufficient coins"})
		default:
			return c.JSON(500, model.ErrorResponse{Errors: "internal error"})
		}
	}
	return c.JSON(200, map[string]interface{}{"status": "success"})
}
