package handlers

import (
	"net/http"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/services"
	"github.com/forzeyy/avito-internship-service/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SendCoinHandler struct {
	transactionSvc services.TransactionService
}

func NewSendCoinHandler(transactionSvc services.TransactionService) *SendCoinHandler {
	return &SendCoinHandler{
		transactionSvc: transactionSvc,
	}
}

func (h *SendCoinHandler) SendCoins(c echo.Context) error {
	var req models.SendCoinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Неверный запрос.",
		})
	}

	fromUserID, _ := utils.GetUserIDFromContext(c)
	toUserID, err := uuid.Parse(req.ToUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Неверный ID получателя.",
		})
	}

	err = h.transactionSvc.SendCoins(c.Request().Context(), fromUserID, toUserID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
