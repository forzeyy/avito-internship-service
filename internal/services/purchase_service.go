package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/repositories"
	"github.com/google/uuid"
)

type PurchaseService struct {
	purchaseRepo repositories.PurchaseRepository
	merchRepo    repositories.MerchRepository
	userRepo     repositories.UserRepository
}

func NewPurchaseService(purchaseRepo repositories.PurchaseRepository, merchRepo repositories.MerchRepository, userRepo repositories.UserRepository) *PurchaseService {
	return &PurchaseService{
		purchaseRepo: purchaseRepo,
		merchRepo:    merchRepo,
		userRepo:     userRepo,
	}
}

func (s *PurchaseService) BuyItem(ctx context.Context, userID uuid.UUID, itemName string) error {
	price, err := s.merchRepo.GetItemPrice(ctx, itemName)
	if err != nil {
		return errors.New("товар не найден")
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	if user.Coins < price {
		return errors.New("недостаточно монет")
	}

	if err := s.purchaseRepo.CreatePurchase(ctx, userID, itemName, 1); err != nil {
		return fmt.Errorf("ошибка при создании записи о покупке: %v", err)
	}

	if err := s.userRepo.UpdateUserBalance(ctx, userID, -price); err != nil {
		return fmt.Errorf("ошибка при обновлении баланса: %v", err)
	}

	return nil
}

func (s *PurchaseService) GetPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error) {
	return s.purchaseRepo.GetPurchasesByUserID(ctx, userID)
}
