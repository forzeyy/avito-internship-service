package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepositoryImpl struct {
	db *database.DB
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUserBalance(ctx context.Context, userID uuid.UUID, amount int) error
}

func NewUserRepository(db *database.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, coins FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&user.ID, &user.Username, &user.Coins)
	if err == pgx.ErrNoRows {
		return nil, errors.New("пользователь не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя: %v", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, password_hash, coins FROM users WHERE username = $1`
	row := r.db.QueryRow(ctx, query, username)

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins)
	if err == pgx.ErrNoRows {
		return nil, errors.New("пользователь не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя: %v", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *models.User) error {
	if user.ID == uuid.Nil {
		log.Printf("ID пользователя не должен быть нулевым")
		return errors.New("ID пользователя не должен быть нулевым")
	}

	if user.PasswordHash == "" {
		log.Printf("PasswordHash не может быть пустым")
		return errors.New("PasswordHash не может быть пустым")
	}

	query := `INSERT INTO users (id, username, password_hash, coins) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Username, user.PasswordHash, user.Coins)
	if err != nil {
		return fmt.Errorf("ошибка при создании пользователя: %v", err)
	}
	return nil
}

func (r *UserRepositoryImpl) UpdateUserBalance(ctx context.Context, userID uuid.UUID, amount int) error {
	query := `UPDATE users SET coins = coins + $1 WHERE id = $2`
	tag, err := r.db.Exec(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса пользователя: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return errors.New("пользователь не найден")
	}

	return nil
}
