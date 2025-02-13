package handlers

import (
	"net/http"

	"github.com/forzeyy/avito-internship-service/internal/services"
	"github.com/forzeyy/avito-internship-service/internal/utils"
	"github.com/labstack/echo/v4"
)

type BuyHandler struct {
	purchaseSvc services.PurchaseService
}

func NewBuyHandler(purchaseSvc services.PurchaseService) *BuyHandler {
	return &BuyHandler{
		purchaseSvc: purchaseSvc,
	}
}

func (h *BuyHandler) BuyItem(c echo.Context) error {
	itemName := c.Param("item")
	if itemName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Товар не указан.",
		})
	}

	userID, _ := utils.GetUserIDFromContext(c)

	err := h.purchaseSvc.BuyItem(c.Request().Context(), userID, itemName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
