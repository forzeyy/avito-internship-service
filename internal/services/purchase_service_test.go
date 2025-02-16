package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/repositories"
	"github.com/forzeyy/avito-internship-service/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPurchaseRepo struct {
	repositories.PurchaseRepository
	mock.Mock
}

type MockMerchRepo struct {
	repositories.MerchRepository
	mock.Mock
}

func (m *MockMerchRepo) GetItemPrice(ctx context.Context, itemName string) (int, error) {
	args := m.Called(ctx, itemName)
	return args.Int(0), args.Error(1)
}

func (m *MockPurchaseRepo) CreatePurchase(ctx context.Context, userID uuid.UUID, itemName string, quantity int) error {
	args := m.Called(ctx, userID, itemName, quantity)
	return args.Error(0)
}

func TestBuyItem_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	itemName := "Item"
	itemPrice := 100
	user := &models.User{ID: userID, Coins: 200}

	mockPurchaseRepo := new(MockPurchaseRepo)
	mockMerchRepo := new(MockMerchRepo)
	mockUserRepo := new(MockUserRepo)

	mockMerchRepo.On("GetItemPrice", ctx, itemName).Return(itemPrice, nil)
	mockUserRepo.On("GetUserByID", ctx, userID).Return(user, nil)
	mockPurchaseRepo.On("CreatePurchase", ctx, userID, itemName, 1).Return(nil)
	mockUserRepo.On("UpdateUserBalance", ctx, userID, -itemPrice).Return(nil)

	service := services.NewPurchaseService(mockPurchaseRepo, mockMerchRepo, mockUserRepo)
	err := service.BuyItem(ctx, userID, itemName)

	assert.NoError(t, err)
	mockMerchRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}

func TestBuyItem_InsufficientFunds(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	itemName := "Item"
	itemPrice := 100
	user := &models.User{ID: userID, Coins: 50}

	mockMerchRepo := new(MockMerchRepo)
	mockUserRepo := new(MockUserRepo)
	mockPurchaseRepo := new(MockPurchaseRepo)

	mockMerchRepo.On("GetItemPrice", ctx, itemName).Return(itemPrice, nil)
	mockUserRepo.On("GetUserByID", ctx, userID).Return(user, nil)

	service := services.NewPurchaseService(mockPurchaseRepo, mockMerchRepo, mockUserRepo)
	err := service.BuyItem(ctx, userID, itemName)

	assert.Error(t, err)
	assert.Equal(t, "недостаточно монет", err.Error())
}

func TestBuyItem_ItemNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	itemName := "Secret Item"

	mockMerchRepo := new(MockMerchRepo)
	mockUserRepo := new(MockUserRepo)
	mockPurchaseRepo := new(MockPurchaseRepo)

	mockMerchRepo.On("GetItemPrice", ctx, itemName).Return(0, errors.New("товар не найден"))

	service := services.NewPurchaseService(mockPurchaseRepo, mockMerchRepo, mockUserRepo)
	err := service.BuyItem(ctx, userID, itemName)

	assert.Error(t, err)
	assert.Equal(t, "товар не найден", err.Error())
}

func TestBuyItem_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	itemName := "Item"
	itemPrice := 100

	mockMerchRepo := new(MockMerchRepo)
	mockUserRepo := new(MockUserRepo)
	mockPurchaseRepo := new(MockPurchaseRepo)

	mockMerchRepo.On("GetItemPrice", ctx, itemName).Return(itemPrice, nil)
	mockUserRepo.On("GetUserByID", ctx, userID).Return(&models.User{}, errors.New("пользователь не найден"))

	service := services.NewPurchaseService(mockPurchaseRepo, mockMerchRepo, mockUserRepo)
	err := service.BuyItem(ctx, userID, itemName)

	assert.Error(t, err)
	assert.Equal(t, "пользователь не найден", err.Error())
}
