package api

import (
	"github.com/gin-gonic/gin"
	"github.com/telemoz/backend/internal/handlers"
	"github.com/telemoz/backend/internal/middleware"
	"go.uber.org/zap"
)

func SetupRoutes(logger *zap.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (no authentication required)
		auth := api.Group("/auth")
		authHandler := handlers.NewAuthHandler()
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Profile routes
			profileHandler := handlers.NewProfileHandler()
			protected.GET("/profile", profileHandler.GetProfile)
			protected.PUT("/profile", profileHandler.UpdateProfile)

			// Trip routes (customer)
			tripHandler := handlers.NewTripHandler()
			trips := protected.Group("/trips")
			trips.Use(middleware.RequireUserType("customer"))
			{
				trips.POST("", tripHandler.CreateTrip)
				trips.GET("/active", tripHandler.GetActiveTrip)
				trips.GET("/history", tripHandler.GetTripHistory)
				trips.GET("/:id", tripHandler.GetTripByID)
				trips.PUT("/:id", tripHandler.UpdateTrip)
				trips.POST("/:id/cancel", tripHandler.CancelTrip)
			}

			// Job routes (driver)
			jobHandler := handlers.NewJobHandler()
			jobs := protected.Group("/jobs")
			jobs.Use(middleware.RequireUserType("driver"))
			{
				jobs.GET("/available", jobHandler.GetAvailableJobs)
				jobs.POST("/:id/accept", jobHandler.AcceptJob)
				jobs.POST("/:id/reject", jobHandler.RejectJob)
				jobs.GET("/active", jobHandler.GetActiveJob)
				jobs.GET("/history", jobHandler.GetJobHistory)
				jobs.PUT("/:id/status", jobHandler.UpdateJobStatus)
			}

			// Children routes (parent)
			childHandler := handlers.NewChildHandler()
			children := protected.Group("/children")
			children.Use(middleware.RequireUserType("parent"))
			{
				children.GET("", childHandler.ListChildren)
				children.POST("", childHandler.CreateChild)
				children.GET("/:id", childHandler.GetChildByID)
				children.PUT("/:id", childHandler.UpdateChild)
				children.DELETE("/:id", childHandler.DeleteChild)
			}

			// Bus routes
			busHandler := handlers.NewBusHandler()
			buses := protected.Group("/buses")
			{
				buses.GET("/child/:childId", busHandler.GetBusByChildID)
				buses.GET("/:id/track", busHandler.TrackBus)
			}

			// Notification routes
			notificationHandler := handlers.NewNotificationHandler()
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.ListNotifications)
				notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
				notifications.GET("/settings", notificationHandler.GetSettings)
				notifications.PUT("/settings", notificationHandler.UpdateSettings)
			}

			// Earnings routes (driver)
			earningsHandler := handlers.NewEarningsHandler()
			earnings := protected.Group("/earnings")
			earnings.Use(middleware.RequireUserType("driver"))
			{
				earnings.GET("/summary", earningsHandler.GetSummary)
				earnings.GET("/history", earningsHandler.GetHistory)
			}
		}
	}

	return router
}

