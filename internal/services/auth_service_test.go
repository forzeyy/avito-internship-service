package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/services"
	"github.com/forzeyy/avito-internship-service/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserBalance(ctx context.Context, userID uuid.UUID, amount int) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

type MockAuthUtils struct{}

func (m MockAuthUtils) CheckPassword(hashed, plain string) bool {
	return true
}

func (m MockAuthUtils) GenerateAccessToken(userID uuid.UUID, secret string) (string, error) {
	return "mocked_token", nil
}

func TestAuthenticate_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockAuthUtil := MockAuthUtils{}

	service := services.NewAuthService(mockRepo, "secret", mockAuthUtil)

	passwordHash, _ := utils.HashPassword("password")
	user := &models.User{ID: uuid.New(), Username: "test", PasswordHash: passwordHash}

	mockRepo.On("GetUserByUsername", mock.Anything, "test").Return(user, nil)

	token, err := service.Authenticate(context.Background(), "test", "password")

	assert.NoError(t, err)
	assert.Equal(t, "mocked_token", token)
}

func TestAuthenticate_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockAuthUtil := MockAuthUtils{}
	service := services.NewAuthService(mockRepo, "secret", mockAuthUtil)

	mockRepo.On("GetUserByUsername", mock.Anything, "test").Return((*models.User)(nil), errors.New("not found"))

	token, err := service.Authenticate(context.Background(), "test", "password")
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockAuthUtil := MockAuthUtils{}
	service := services.NewAuthService(mockRepo, "secret", mockAuthUtil)

	mockRepo.On("GetUserByUsername", mock.Anything, "test").Return((*models.User)(nil), errors.New("not found"))
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	err := service.Register(context.Background(), "test", "password")
	assert.NoError(t, err)
}

func TestRegister_UserExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockAuthUtil := MockAuthUtils{}
	service := services.NewAuthService(mockRepo, "secret", mockAuthUtil)

	user := &models.User{ID: uuid.New(), Username: "test", PasswordHash: "hash"}
	mockRepo.On("GetUserByUsername", mock.Anything, "test").Return(user, nil)

	err := service.Register(context.Background(), "test", "password")
	assert.Error(t, err)
	assert.Equal(t, "пользователь уже существует", err.Error())
}
