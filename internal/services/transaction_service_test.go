package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepo struct {
	mock.Mock
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) CreateTransaction(ctx context.Context, fromUserID, toUserID uuid.UUID, toUsername, fromUsername string, amount int) error {
	args := m.Called(ctx, fromUserID, toUserID, toUsername, fromUsername, amount)
	return args.Error(0)
}

func (m *MockTransactionRepo) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateUserBalance(ctx context.Context, userID uuid.UUID, amount int) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func TestSendCoins(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepo)
	mockUserRepo := new(MockUserRepo)
	service := services.NewTransactionService(mockTransactionRepo, mockUserRepo)

	ctx := context.Background()
	fromUserID := uuid.New()
	toUserID := uuid.New()
	fromUser := &models.User{ID: fromUserID, Username: "avito_user111", Coins: 100}
	toUser := &models.User{ID: toUserID, Username: "avito_user222", Coins: 50}
	amount := 30

	mockUserRepo.On("GetUserByID", ctx, fromUserID).Return(fromUser, nil)
	mockUserRepo.On("GetUserByUsername", ctx, "avito_user222").Return(toUser, nil)
	mockTransactionRepo.On("CreateTransaction", ctx, fromUserID, toUserID, "avito_user222", "avito_user111", amount).Return(nil)
	mockUserRepo.On("UpdateUserBalance", ctx, fromUserID, -amount).Return(nil)
	mockUserRepo.On("UpdateUserBalance", ctx, toUserID, amount).Return(nil)

	err := service.SendCoins(ctx, fromUserID, "avito_user222", amount)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestSendCoins_InsufficientBalance(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepo)
	mockUserRepo := new(MockUserRepo)
	service := services.NewTransactionService(mockTransactionRepo, mockUserRepo)

	ctx := context.Background()
	fromUserID := uuid.New()
	toUserID := uuid.New()
	fromUser := &models.User{ID: fromUserID, Username: "avito_user111", Coins: 10}
	toUser := &models.User{ID: toUserID, Username: "avito_user222", Coins: 50}
	amount := 30

	mockUserRepo.On("GetUserByID", ctx, fromUserID).Return(fromUser, nil)
	mockUserRepo.On("GetUserByUsername", ctx, "avito_user222").Return(toUser, nil)

	err := service.SendCoins(ctx, fromUserID, "avito_user222", amount)

	assert.EqualError(t, err, "недостаточно монет")
	mockUserRepo.AssertExpectations(t)
}

func TestSendCoins_UserNotFound(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepo)
	mockUserRepo := new(MockUserRepo)
	service := services.NewTransactionService(mockTransactionRepo, mockUserRepo)

	ctx := context.Background()
	fromUserID := uuid.New()

	mockUserRepo.On("GetUserByID", ctx, fromUserID).Return((*models.User)(nil), errors.New("отправитель не найден"))

	err := service.SendCoins(ctx, fromUserID, "avito_user222", 30)

	assert.EqualError(t, err, "отправитель не найден")
	mockUserRepo.AssertExpectations(t)
}

func TestGetTransactionsByUserID(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepo)
	mockUserRepo := new(MockUserRepo)
	service := services.NewTransactionService(mockTransactionRepo, mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	transactions := []models.Transaction{
		{ID: uuid.New(), Amount: 50},
		{ID: uuid.New(), Amount: 30},
	}

	mockTransactionRepo.On("GetTransactionsByUserID", ctx, userID).Return(transactions, nil)

	result, err := service.GetTransactionsByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, transactions, result)
	mockTransactionRepo.AssertExpectations(t)
}
