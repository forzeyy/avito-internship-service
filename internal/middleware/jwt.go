package middleware

import (
	"os"

	jwtMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

func JWTMiddleware() echo.MiddlewareFunc {
	return jwtMiddleware.WithConfig(jwtMiddleware.Config{
		SigningKey:    []byte(JWTSecret),
		TokenLookup:   "header:Authorization:Bearer ",
		SigningMethod: "HS256",
	})
}
