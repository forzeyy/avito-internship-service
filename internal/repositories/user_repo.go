package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, coins FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	err := row.Scan(&user.ID, &user.Username, &user.Coins)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user. error: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, coins FROM users WHERE username = $1`
	row := r.db.QueryRow(ctx, query, username)

	err := row.Scan(&user.ID, &user.Username, &user.Coins)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user. error: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, username, password_hash, coins) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Username, user.PasswordHash, user.Coins)
	if err != nil {
		return fmt.Errorf("failed to create user. error: %v", err)
	}
	return nil
}

func (r *UserRepository) UpdateUserBalance(ctx context.Context, userID uuid.UUID, amount int) error {
	query := `UPDATE users SET coins = coins + $1 WHERE id = $2`
	tag, err := r.db.Exec(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update user balance. error: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
