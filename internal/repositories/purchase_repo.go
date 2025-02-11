package repositories

import (
	"context"
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/google/uuid"
)

type PurchaseRepository struct {
	db *database.DB
}

func NewPurchaseRepository(db *database.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) CreatePurchase(ctx context.Context, userID uuid.UUID, itemName string, quantity int) error {
	query := `
		INSERT INTO purchases (user_id, item_name, quantity, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := r.db.Exec(ctx, query, userID, itemName, quantity)
	if err != nil {
		return fmt.Errorf("failed to create purchase. error: %v", err)
	}
	return nil
}

func (r *PurchaseRepository) GetPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error) {
	var purchases []models.Purchase
	query := `
		SELECT id, user_id, item_name, quantity, created_at
		FROM purchases
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchases. error: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var p models.Purchase
		if err := rows.Scan(&p.ID, &p.UserID, &p.ItemName, &p.Quantity, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan purchase. error: %v", err)
		}
		purchases = append(purchases, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("purchase scanning failed. error: %v", err)
	}

	return purchases, nil
}
