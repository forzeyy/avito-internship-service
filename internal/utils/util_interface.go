package utils

import "github.com/google/uuid"

type AuthUtils interface {
	CheckPassword(hashed, plain string) bool
	GenerateAccessToken(userID uuid.UUID, secret string) (string, error)
}

type DefaultAuthUtils struct{}

func (d DefaultAuthUtils) CheckPassword(hashed, plain string) bool {
	return CheckPassword(hashed, plain)
}

func (d DefaultAuthUtils) GenerateAccessToken(userID uuid.UUID, secret string) (string, error) {
	return GenerateAccessToken(userID, secret)
}
