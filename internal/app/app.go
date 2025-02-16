package app

import (
	"fmt"

	"github.com/forzeyy/avito-internship-service/internal/config"
	"github.com/forzeyy/avito-internship-service/internal/database"
	"github.com/forzeyy/avito-internship-service/internal/handlers"
	"github.com/forzeyy/avito-internship-service/internal/middleware"
	"github.com/forzeyy/avito-internship-service/internal/repositories"
	"github.com/forzeyy/avito-internship-service/internal/services"
	"github.com/forzeyy/avito-internship-service/internal/utils"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config) error {
	e := echo.New()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	dbConn, err := database.ConnectDatabase(dsn)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbConn.Close()

	// repos
	userRepo := repositories.NewUserRepository(dbConn)
	transactionRepo := repositories.NewTransactionRepository(dbConn)
	purchaseRepo := repositories.NewPurchaseRepository(dbConn)
	merchRepo := repositories.NewMerchRepository(dbConn)

	// services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, utils.DefaultAuthUtils{})
	userService := services.NewUserService(userRepo)
	transactionService := services.NewTransactionService(transactionRepo, userRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo, merchRepo, userRepo)

	// handlers
	authHandler := handlers.NewAuthHandler(*authService, *userService)
	infoHandler := handlers.NewInfoHandler(*userService, *transactionService, *purchaseService)
	sendCoinHandler := handlers.NewSendCoinHandler(*transactionService)
	buyHandler := handlers.NewBuyHandler(*purchaseService)

	// open routes
	e.POST("/api/auth", authHandler.Auth)

	// protected routes
	protected := e.Group("")
	protected.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	protected.GET("/api/info", infoHandler.GetInfo)
	protected.POST("/api/sendCoin", sendCoinHandler.SendCoins)
	protected.GET("/api/buy/:item", buyHandler.BuyItem)

	e.Logger.Fatal(e.Start(":8080"))

	return nil
}
