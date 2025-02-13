package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *database.DB
}

func NewTransactionRepository(db *database.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, fromUserID, toUserID uuid.UUID, toUser, fromUser string, amount int) error {
	query := `
		INSERT INTO transactions (from_user_id, from_user, to_user_id, to_user, amount, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err := r.db.Exec(ctx, query, fromUserID, fromUser, toUserID, toUser, amount)
	if err != nil {
		return fmt.Errorf("ошибка при создании транзакции: %v", err)
	}
	return nil
}

func (r *TransactionRepository) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := `
		SELECT id, from_user, from_user_id, to_user, to_user_id, amount, created_at
		FROM transactions
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		log.Printf("ошибка при получении транзакций: %v", err)
		return nil, fmt.Errorf("ошибка при получении транзакций: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tr models.Transaction
		if err := rows.Scan(&tr.ID, &tr.FromUser, &tr.FromUserID, &tr.ToUser, &tr.ToUserID, &tr.Amount, &tr.CreatedAt); err != nil {
			log.Printf("ошибка при скане транзакции: %v", err)
			return nil, fmt.Errorf("ошибка при скане транзакции: %v", err)
		}
		transactions = append(transactions, tr)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ошибка при скане транзакций: %v", err)
		return nil, fmt.Errorf("ошибка при скане транзакций: %v", err)
	}

	return transactions, nil
}
