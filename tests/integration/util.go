package integration

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func generateToken(userID uuid.UUID, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 7 * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func getUserIDFromToken(token, secret string) (uuid.UUID, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, errors.New("некорректный формат user_id")
	}

	return uuid.Parse(userIDStr)
}
