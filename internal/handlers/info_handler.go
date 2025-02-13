package handlers

import (
	"net/http"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/services"
	"github.com/forzeyy/avito-internship-service/internal/utils"
	"github.com/labstack/echo/v4"
)

type InfoHandler struct {
	userSvc        services.UserService
	transactionSvc services.TransactionService
	purchaseSvc    services.PurchaseService
}

func NewInfoHandler(userSvc services.UserService, transactionSvc services.TransactionService, purchaseSvc services.PurchaseService) *InfoHandler {
	return &InfoHandler{
		userSvc:        userSvc,
		transactionSvc: transactionSvc,
		purchaseSvc:    purchaseSvc,
	}
}

func (h *InfoHandler) GetInfo(c echo.Context) error {
	userID, _ := utils.GetUserIDFromContext(c)
	user, err := h.userSvc.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Пользователь не найден.",
		})
	}

	transactions, err := h.transactionSvc.GetTransactionsByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Ошибка при получении истории транзакций.",
		})
	}

	received := make([]models.TransactionHistoryItem, 0)
	sent := make([]models.TransactionHistoryItem, 0)
	for _, tr := range transactions {
		if tr.FromUserID == userID {
			sent = append(sent, models.TransactionHistoryItem{
				ToUser: tr.ToUserID.String(),
				Amount: tr.Amount,
			})
		} else {
			received = append(received, models.TransactionHistoryItem{
				FromUser: tr.FromUserID.String(),
				Amount:   tr.Amount,
			})
		}
	}

	purchases, err := h.purchaseSvc.GetPurchasesByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Ошибка при получении списка покупок.",
		})
	}

	inventory := make(map[string]int)
	for _, p := range purchases {
		inventory[p.ItemName] += p.Quantity
	}

	response := models.InfoResponse{
		Coins:     user.Coins,
		Inventory: []models.InventoryItem{},
		CoinHistory: models.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}

	for itemName, quantity := range inventory {
		response.Inventory = append(response.Inventory, models.InventoryItem{
			Type:     itemName,
			Quantity: quantity,
		})
	}

	return c.JSON(http.StatusOK, response)
}
