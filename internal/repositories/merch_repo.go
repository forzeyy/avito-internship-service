package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/jackc/pgx/v5"
)

type MerchRepositoryImpl struct {
	db *database.DB
}

func NewMerchRepository(db *database.DB) *MerchRepositoryImpl {
	return &MerchRepositoryImpl{db: db}
}

type MerchRepository interface {
	GetItemPrice(ctx context.Context, itemName string) (int, error)
}

func (r *MerchRepositoryImpl) GetItemPrice(ctx context.Context, itemName string) (int, error) {
	var price int
	query := `SELECT price FROM merch WHERE name = $1`

	row := r.db.QueryRow(ctx, query, itemName)
	err := row.Scan(&price)
	if err == pgx.ErrNoRows {
		return 0, errors.New("мерч не найден")
	}
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении цены: %v", err)
	}

	return price, nil
}
