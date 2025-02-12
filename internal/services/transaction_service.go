package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/repositories"
	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo repositories.TransactionRepository
	userRepo        repositories.UserRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, userRepo repositories.UserRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (s *TransactionService) SendCoins(ctx context.Context, fromUserID, toUserID uuid.UUID, amount int) error {
	fromUser, err := s.userRepo.GetUserByID(ctx, fromUserID)
	if err != nil {
		return errors.New("отправитель не найден")
	}

	_, err = s.userRepo.GetUserByID(ctx, toUserID)
	if err != nil {
		return errors.New("получатель не найден")
	}

	if fromUser.Coins < amount {
		return errors.New("недостаточно монет")
	}

	if err := s.transactionRepo.CreateTransaction(ctx, fromUserID, toUserID, amount); err != nil {
		return fmt.Errorf("ошибка при создании транзакции: %v", err)
	}

	if err := s.userRepo.UpdateUserBalance(ctx, fromUserID, -amount); err != nil {
		return fmt.Errorf("ошибка при обновлении баланса отправителя: %v", err)
	}
	if err := s.userRepo.UpdateUserBalance(ctx, toUserID, amount); err != nil {
		return fmt.Errorf("ошибка при обновлении баланса получателя: %v", err)
	}

	return nil
}

func (s *TransactionService) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	return s.transactionRepo.GetTransactionsByUserID(ctx, userID)
}
