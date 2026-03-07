package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/piotrzalecki/budget-api/internal/handler"
	"github.com/piotrzalecki/budget-api/internal/repo"
	"github.com/piotrzalecki/budget-api/pkg/model"
)

func setupRoutes(router *gin.Engine, logger *zap.Logger, handlers *handler.Handler, repository repo.Repository, version string) {
	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health endpoint (no auth required)
	router.GET("/health", healthHandler(logger, version))

	// Public auth routes (no middleware)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login", handler.ValidateRequest[model.LoginRequest](), handlers.Login)
	}

	// API v1 routes (protected by session token)
	v1 := router.Group("/api/v1")
	v1.Use(handler.SessionAuth(repository))
	{
		// Auth
		v1.POST("/auth/logout", handlers.Logout)

		// User routes
		v1.GET("/users", handlers.ListUsers)
		v1.POST("/users", handler.ValidateRequest[model.CreateUserRequest](), handlers.CreateUser)
		v1.GET("/users/:id", handlers.GetUserByID)
		v1.PATCH("/users/:id", handler.ValidateRequest[model.UpdateUserRequest](), handlers.UpdateUser)
		v1.DELETE("/users/:id", handlers.DeleteUser)

		// Transaction routes with validation
		v1.POST("/transactions", handler.ValidateRequest[model.CreateTransactionRequest](), handlers.CreateTransaction)
		v1.GET("/transactions", handlers.GetTransactions)
		v1.GET("/transactions/:id", handlers.GetTransactionByID)
		v1.PATCH("/transactions/:id", handler.ValidateRequest[model.UpdateTransactionRequest](), handlers.UpdateTransaction)
		v1.DELETE("/transactions/:id", handlers.HardDeleteTransaction) //Commented out until Admin user will be implemented
		v1.GET("/transactions/by-recurring/:recurring_id", handlers.GetTransactionsByRecurringID)
		v1.GET("/transactions/by-tag/:tag_id", handlers.GetTransactionsByTag)
		v1.POST("/transactions/purge", handler.ValidateRequest[model.PurgeTransactionsRequest](), handlers.PurgeSoftDeletedTransactions)
		
		// Tag routes with validation
		v1.POST("/tags", handler.ValidateRequest[model.CreateTagRequest](), handlers.CreateTag)
		v1.GET("/tags", handlers.GetTags)
		v1.PATCH("/tags/:id", handler.ValidateRequest[model.UpdateTagRequest](), handlers.UpdateTag)
		v1.DELETE("/tags/:id", handlers.DeleteTag)
		
		// Recurring routes with validation
		v1.POST("/recurring", handler.ValidateRequest[model.CreateRecurringRequest](), handlers.CreateRecurring)
		v1.GET("/recurring", handlers.GetRecurring)
		v1.GET("/recurring/:id", handlers.GetRecurringByID)
		v1.PATCH("/recurring/:id", handler.ValidateRequest[model.UpdateRecurringRequest](), handlers.UpdateRecurring)
		v1.DELETE("/recurring/:id", handlers.DeleteRecurring)
		v1.GET("/recurring/by-tag/:tag_id", handlers.GetRecurringByTag)
		v1.GET("/recurring/active", handlers.ListActiveRecurring)
		v1.PATCH("/recurring/:id/toggle", handlers.ToggleRecurringActive)
		v1.GET("/recurring/due", handlers.GetRecurringDueOnDate)
		
		// Reports routes
		v1.GET("/reports/monthly", handlers.GetMonthlyReport)
		v1.GET("/reports/monthly/totals", handlers.GetMonthlyTotals)
		
		// Placeholder route to use v1 variable
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data":  "Budget API v1",
				"error": nil,
			})
		})
	}

	// Admin routes (protected by API key)
	admin := router.Group("/admin")
	admin.Use(handler.APIKeyAuth())
	{
		// Scheduler endpoint
		admin.POST("/run-scheduler", handlers.RunScheduler)
		
		// Placeholder route to use admin variable
		admin.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data":  "Admin endpoint",
				"error": nil,
			})
		})
	}

	// Add a catch-all route for undefined endpoints
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Endpoint not found",
			"data":  nil,
		})
	})
}

// @Summary Health check
// @Description Check if the API is healthy and running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Health status"
// @Router /health [get]
func healthHandler(logger *zap.Logger, version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Health check requested")
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"status":    "healthy",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   version,
			},
			"error": nil,
		})
	}
} 