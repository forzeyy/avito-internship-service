package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GenerateAccessToken(userID uuid.UUID, secret string) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("секрет JWT не может быть пустым")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 7 * 24).Unix(),
	})

	return accessToken.SignedString([]byte(secret))
}

func VerifyAccessToken(tokenString string, secret string) (*uuid.UUID, error) {
	if len(secret) == 0 {
		return nil, errors.New("секрет JWT не может быть пустым")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok || userIDStr == "" {
			return nil, jwt.ErrInvalidKey
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.New("некорректный формат user_id")
		}

		return &userID, nil
	}

	return nil, jwt.ErrTokenMalformed
}

func GetUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	token := c.Get("user").(*jwt.Token)
	if token == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Неавторизован.")
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(string)
	userUUID, _ := uuid.Parse(userID)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Неавторизован.")
	}

	return userUUID, nil
}
