package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/repositories"
	"github.com/forzeyy/avito-internship-service/internal/utils"
)

type AuthService struct {
	userRepo repositories.UserRepository
	jwtKey   []byte
}

func NewAuthService(userRepo repositories.UserRepository, jwtKey string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
	}
}

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("неверные учетные данные")
	}

	if !utils.CheckPassword(user.PasswordHash, password) {
		return "", errors.New("неверные учетные данные")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании токена: %v", err)
	}

	return accessToken, nil
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	_, err := s.userRepo.GetUserByUsername(ctx, username)
	if err == nil {
		return errors.New("пользователь уже существует")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("ошибка при хешировании пароля: %v", err)
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	return s.userRepo.CreateUser(ctx, user)
}
