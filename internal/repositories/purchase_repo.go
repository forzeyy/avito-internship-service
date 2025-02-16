package repositories

import (
	"context"
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/google/uuid"
)

type PurchaseRepositoryImpl struct {
	db *database.DB
}

func NewPurchaseRepository(db *database.DB) *PurchaseRepositoryImpl {
	return &PurchaseRepositoryImpl{db: db}
}

type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, userID uuid.UUID, itemName string, quantity int) error
	GetPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error)
}

func (r *PurchaseRepositoryImpl) CreatePurchase(ctx context.Context, userID uuid.UUID, itemName string, quantity int) error {
	query := `
		INSERT INTO purchases (user_id, item_name, quantity, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := r.db.Exec(ctx, query, userID, itemName, quantity)
	if err != nil {
		return fmt.Errorf("ошибка при создании покупки: %v", err)
	}
	return nil
}

func (r *PurchaseRepositoryImpl) GetPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error) {
	var purchases []models.Purchase
	query := `
		SELECT id, user_id, item_name, quantity, created_at
		FROM purchases
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении покупок: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var p models.Purchase
		if err := rows.Scan(&p.ID, &p.UserID, &p.ItemName, &p.Quantity, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка при скане покупки: %v", err)
		}
		purchases = append(purchases, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при скане покупок: %v", err)
	}

	return purchases, nil
}
