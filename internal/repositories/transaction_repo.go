package repositories

import (
	"context"
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/google/uuid"
)

type TransactionRepository struct {
	conn *database.DB
}

func NewTransactionRepository(conn *database.DB) *TransactionRepository {
	return &TransactionRepository{conn: conn}
}

func (db *TransactionRepository) CreateTransaction(ctx context.Context, fromUserID, toUserID uuid.UUID, amount int) error {
	query := `
		INSERT INTO transactions (from_user_id, to_user_id, amount, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := db.conn.Exec(ctx, query, fromUserID, toUserID, amount)
	if err != nil {
		return fmt.Errorf("failed to create transaction. error: %v", err)
	}
	return nil
}

func (db *TransactionRepository) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := `
		SELECT id, from_user_id, to_user_id, amount, created_at
		FROM transactions
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions. error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tr models.Transaction
		if err := rows.Scan(&tr.ID, &tr.FromUserID, &tr.ToUserID, &tr.Amount, &tr.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan transaction. error: %v", err)
		}
		transactions = append(transactions, tr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("transactions scanning error: %v", err)
	}

	return transactions, nil
}
