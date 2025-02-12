package handlers

import (
	"net/http"

	"github.com/forzeyy/avito-internship-service/internal/models"
	"github.com/forzeyy/avito-internship-service/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authSvc services.AuthService
	userSvc services.UserService
}

func NewAuthHandler(authSvc services.AuthService, userSvc services.UserService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (h *AuthHandler) Auth(c echo.Context) error {
	var req models.AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Неверный запрос.",
		})
	}

	token, err := h.authSvc.Authenticate(c.Request().Context(), req.Username, req.Password)
	if err == nil {
		return c.JSON(http.StatusOK, models.AuthResponse{Token: token})
	}

	if err.Error() == "неверные учетные данные" {
		existingUser, _ := h.userSvc.GetUserByUsername(c.Request().Context(), req.Username)
		if existingUser == nil {
			err := h.authSvc.Register(c.Request().Context(), req.Username, req.Password)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"error": "Ошибка регистрации пользователя.",
				})
			}

			token, err := h.authSvc.Authenticate(c.Request().Context(), req.Username, req.Password)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"error": "Ошибка при выдаче токена.",
				})
			}

			return c.JSON(http.StatusOK, models.AuthResponse{Token: token})
		}
	}

	return c.JSON(http.StatusUnauthorized, echo.Map{
		"error": "Неверные учетные данные.",
	})
}
