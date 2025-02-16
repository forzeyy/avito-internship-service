package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestSendCoinsIntegration(t *testing.T) {
	e := echo.New()

	db, err := database.ConnectDatabase("postgres://postgres:postgres@db:5432/avito?sslmode=disable")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, "secretsecret", utils.DefaultAuthUtils{})
	transactionService := services.NewTransactionService(transactionRepo, userRepo)

	authHandler := handlers.NewAuthHandler(*authService, *userService)
	sendCoinHandler := handlers.NewSendCoinHandler(*transactionService)

	e.POST("/api/auth", authHandler.Auth)
	e.POST("/api/sendCoin", sendCoinHandler.SendCoins, echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secretsecret"),
	}))

	// рег отправителя и получателя
	var users []*models.User
	for _, username := range []string{"sender", "receiver"} {
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
			t.Fatalf("ошибка декодирования токена для пользователя %s: %v", username, err)
		}

		if authResponse.Token == "" {
			t.Fatalf("токен не был выдан для пользователя %s", username)
		}

		userID, err := getUserIDFromToken(authResponse.Token, "secretsecret")
		if err != nil {
			t.Fatalf("ошибка извлечения user_id из токена для пользователя %s: %v", username, err)
		}

		users = append(users, &models.User{ID: userID, Username: username})
	}

	sender, receiver := users[0], users[1]

	// отправка монет
	amount := 50
	reqBody := strings.NewReader(`{"toUser":"` + receiver.ID.String() + `","amount":` + strconv.Itoa(amount) + `}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", reqBody)
	req.Header.Set("Authorization", "Bearer "+generateToken(sender.ID, "secretsecret"))
	req.Header.Set("Content-Type", "application/json")

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// проверка обновления баланса отправителя
	senderAfter, err := userRepo.GetUserByID(context.Background(), sender.ID)
	if err != nil {
		t.Fatalf("ошибка получения информации о пользователе %s: %v", sender.Username, err)
	}

	expectedSenderBalance := 1000 - amount
	assert.Equal(t, expectedSenderBalance, senderAfter.Coins, "баланс отправителя не обновился корректно")

	// проверка обновления баланса получателя
	receiverAfter, err := userRepo.GetUserByID(context.Background(), receiver.ID)
	if err != nil {
		t.Fatalf("ошибка получения информации о пользователе %s: %v", receiver.Username, err)
	}

	expectedReceiverBalance := 1000 + amount
	assert.Equal(t, expectedReceiverBalance, receiverAfter.Coins, "баланс получателя не обновился корректно")

	// проверка записи
	transactions, err := transactionRepo.GetTransactionsByUserID(context.Background(), sender.ID)
	if err != nil {
		t.Fatalf("ошибка получения транзакций для пользователя %s: %v", sender.Username, err)
	}

	assert.Len(t, transactions, 1, "запись о транзакции не была создана")
	assert.Equal(t, amount, transactions[0].Amount, "некорректная сумма транзакции")
	assert.Equal(t, sender.ID, transactions[0].FromUserID, "некорректный отправитель транзакции")
	assert.Equal(t, receiver.ID, transactions[0].ToUserID, "некорректный получатель транзакции")
}
