package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"encoding/json"

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

func TestPurchaseIntegration(t *testing.T) {
	e := echo.New()

	db, err := database.ConnectDatabase("postgres://postgres:postgres@db:5432/avito?sslmode=disable")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	purchaseRepo := repositories.NewPurchaseRepository(db)
	merchandiseRepo := repositories.NewMerchRepository(db)

	authService := services.NewAuthService(userRepo, "secretsecret", utils.DefaultAuthUtils{})
	purchaseService := services.NewPurchaseService(purchaseRepo, merchandiseRepo, userRepo)
	userService := services.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(*authService, *userService)
	purchaseHandler := handlers.NewBuyHandler(*purchaseService)

	e.POST("/api/auth", authHandler.Auth)
	e.GET("/api/buy/:item", purchaseHandler.BuyItem, echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secretsecret"),
	}))

	// регистрация
	username := "purchaseuser"
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
		t.Fatalf("Ошибка декодирования токена: %v", err)
	}

	token := authResponse.Token
	if token == "" {
		t.Fatalf("токен не выдан")
	}

	// покупка мерча
	itemName := "t-shirt"
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/buy/"+itemName, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// проверка обновления баланса
	userID, err := getUserIDFromToken(token, "secretsecret")
	if err != nil {
		t.Fatalf("ошибка извлечения user_id из токена: %v", err)
	}

	user, err := userRepo.GetUserByID(context.Background(), userID)
	if err != nil {
		t.Fatalf("ошибка получения пользователя: %v", err)
	}

	expectedBalance := 1000 - 80
	assert.Equal(t, expectedBalance, user.Coins, "баланс пользователя не обновился корректно")

	// проверка на запись о покупке
	purchases, err := purchaseRepo.GetPurchasesByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("ошибка получения покупок: %v", err)
	}

	assert.Len(t, purchases, 1, "запись о покупке не была создана")
	assert.Equal(t, itemName, purchases[0].ItemName, "некорректное имя товара")
	assert.Equal(t, 1, purchases[0].Quantity, "некорректное количество товаров")
}
