package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/repositories"
	"github.com/forzeyy/avito-internship-service/internal/utils"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo repositories.UserRepository
	jwtKey   []byte
	authUtil utils.AuthUtils
}

func NewAuthService(userRepo repositories.UserRepository, jwtKey string, authUtil utils.AuthUtils) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
		authUtil: authUtil,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	log.Printf("Аутентификация пользователя: %s", username)
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		log.Printf("Пользователь %s не найден: %v", username, err)
		return "", errors.New("пользователь не найден")
	}

	if user.PasswordHash == "" {
		log.Printf("PasswordHash пользователя %s пустой", username)
		return "", errors.New("password hash отсутствует")
	}

	if !s.authUtil.CheckPassword(user.PasswordHash, password) {
		log.Printf("Неверный пароль для пользователя %s: %v", username, err)
		return "", errors.New("неверные учетные данные")
	}

	if len(s.jwtKey) == 0 {
		log.Printf("Секрет JWT отсутствует")
		return "", errors.New("секрет JWT отсутствует")
	}

	accessToken, err := s.authUtil.GenerateAccessToken(user.ID, string(s.jwtKey))
	if err != nil {
		log.Printf("Ошибка создания токена для пользователя %s: %v", username, err)
		return "", fmt.Errorf("ошибка при создании токена: %v", err)
	}

	if accessToken == "" {
		log.Printf("Токен пустой для пользователя %s", username)
		return "", errors.New("не удалось создать токен")
	}

	return accessToken, nil
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	log.Printf("Регистрация пользователя: %s", username)
	_, err := s.userRepo.GetUserByUsername(ctx, username)
	if err == nil {
		log.Printf("Пользователь %s уже существует", username)
		return errors.New("пользователь уже существует")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Ошибка хеширования пароля для пользователя %s: %v", username, err)
		return fmt.Errorf("ошибка при хешировании пароля: %v", err)
	}

	log.Printf("Хеш пароля для пользователя %s: %s", username, string(hashedPassword))

	user := &models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: hashedPassword,
		Coins:        1000,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		log.Printf("Ошибка создания пользователя %s: %v", username, err)
		return fmt.Errorf("ошибка при создании пользователя: %v", err)
	}

	log.Printf("Пользователь %s успешно зарегистрирован с ID: %s", username, user.ID.String())
	return nil
}
