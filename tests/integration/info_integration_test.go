package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/forzeyy/avito-internship-service/internal/handlers"
	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/repositories"
	"github.com/forzeyy/avito-internship-service/internal/services"
	"github.com/forzeyy/avito-internship-service/internal/utils"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetInfoIntegration(t *testing.T) {
	e := echo.New()

	db, err := database.ConnectDatabase("postgres://postgres:postgres@db:5432/avito?sslmode=disable")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	purchaseRepo := repositories.NewPurchaseRepository(db)
	merchRepo := repositories.NewMerchRepository(db)

	authService := services.NewAuthService(userRepo, "secretsecret", utils.DefaultAuthUtils{})
	userService := services.NewUserService(userRepo)
	transactionService := services.NewTransactionService(transactionRepo, userRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo, merchRepo, userRepo)

	authHandler := handlers.NewAuthHandler(*authService, *userService)
	infoHandler := handlers.NewInfoHandler(*userService, *transactionService, *purchaseService)

	e.POST("/api/auth", authHandler.Auth)
	e.GET("/api/info", infoHandler.GetInfo, echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secretsecret"),
	}))

	// рег пользователя
	username := "infouser"
	password := "superpassword"
	reqBody := strings.NewReader(`{"username":"` + username + `","password":"` + password + `"}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth", reqBody)
	req.Header.Set("Content-Type", "application/json")

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var authResponse models.AuthResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &authResponse)
	if err != nil {
		t.Fatalf("ошибка декодирования токена: %v", err)
	}

	token := authResponse.Token
	if token == "" {
		t.Fatalf("токен не выдан")
	}

	// получение инфы о пользователе
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// декодирование ответа
	var infoResponse models.InfoResponse
	err = json.Unmarshal([]byte(rec.Body.String()), &infoResponse)
	if err != nil {
		t.Fatalf("ошибка декодирования ответа: %v", err)
	}

	// проверка данных
	assert.Equal(t, 1000, infoResponse.Coins, "начальный баланс должен быть равен 1000")
	assert.Empty(t, infoResponse.Inventory, "инвентарь должен быть пустым")
	assert.Empty(t, infoResponse.CoinHistory.Received, "история входящих транзакций должна быть пустой")
	assert.Empty(t, infoResponse.CoinHistory.Sent, "история исходящих транзакций должна быть пустой")
}
