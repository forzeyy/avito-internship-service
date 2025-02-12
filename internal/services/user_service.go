package services

import (
	"context"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/repositories"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetUserByUsername(ctx, username)
}
