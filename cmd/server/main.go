package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/telemoz/backend/api"
	"github.com/telemoz/backend/internal/config"
	"github.com/telemoz/backend/internal/database"
	"github.com/telemoz/backend/internal/models"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	var logger *zap.Logger
	var err error
	if config.AppConfig.Server.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	// Set Gin mode
	gin.SetMode(config.AppConfig.Server.GinMode)

	// Connect to database
	if err := database.Connect(); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations
	if err := database.Migrate(
		&models.User{},
		&models.Trip{},
		&models.Job{},
		&models.Child{},
		&models.Bus{},
		&models.BusLocation{},
		&models.Notification{},
		&models.NotificationSettings{},
		&models.DriverEarning{},
		&models.RefreshToken{},
	); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	logger.Info("Database migrations completed")

	// Setup routes
	router := api.SetupRoutes(logger)

	// Start server
	port := config.AppConfig.Server.Port
	logger.Info("Starting server", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

