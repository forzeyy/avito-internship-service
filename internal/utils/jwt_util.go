package utils

import (
	"time"

	"github.com/forzeyy/avito-internship-service/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secret = middleware.JWTSecret

func GenerateAccessToken(userID uuid.UUID) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	return accessToken.SignedString(secret)
}

func VerifyAccessToken(tokenString string) (*uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(uuid.UUID)
		if !ok {
			return nil, jwt.ErrInvalidKey
		}
		return &userID, nil
	}

	return nil, jwt.ErrTokenMalformed
}
