package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	conn *database.DB
}

func NewUserRepository(conn *database.DB) *UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (db *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, coins FROM users WHERE id = $1`
	row := db.conn.QueryRow(ctx, query, id)

	err := row.Scan(&user.ID, &user.Username, &user.Coins)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user. error: %v", err)
	}

	return &user, nil
}

func (db *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, coins FROM users WHERE username = $1`
	row := db.conn.QueryRow(ctx, query, username)

	err := row.Scan(&user.ID, &user.Username, &user.Coins)
	if err == pgx.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user. error: %v", err)
	}

	return &user, nil
}

func (db *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, username, password_hash, coins) VALUES ($1, $2, $3, $4)`
	_, err := db.conn.Exec(ctx, query, user.ID, user.Username, user.PasswordHash, user.Coins)
	if err != nil {
		return fmt.Errorf("failed to create user. error: %v", err)
	}
	return nil
}

func (db *UserRepository) UpdateUserBalance(ctx context.Context, userID uint, amount int) error {
	query := `UPDATE users SET coins = coins + $1 WHERE id = $2`
	tag, err := db.conn.Exec(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update user balance. error: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
