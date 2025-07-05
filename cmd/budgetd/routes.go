package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/piotrzalecki/budget-api/internal/handler"
	"github.com/piotrzalecki/budget-api/pkg/model"
)

func setupRoutes(router *gin.Engine, logger *zap.Logger, handlers *handler.Handler) {
	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health endpoint (no auth required)
	router.GET("/health", healthHandler(logger))

	// API v1 routes (protected by API key)
	v1 := router.Group("/api/v1")
	v1.Use(handler.APIKeyAuth())
	{
		// Transaction routes with validation
		v1.POST("/transactions", handler.ValidateRequest[model.CreateTransactionRequest](), handlers.CreateTransaction)
		v1.GET("/transactions", handlers.GetTransactions)
		v1.GET("/transactions/:id", handlers.GetTransactionByID)
		v1.PATCH("/transactions/:id", handler.ValidateRequest[model.UpdateTransactionRequest](), handlers.UpdateTransaction)
		// v1.DELETE("/transactions/:id", handlers.HardDeleteTransaction) Commented out until Admin user will be implemented
		v1.GET("/transactions/by-recurring/:recurring_id", handlers.GetTransactionsByRecurringID)
		v1.GET("/transactions/by-tag/:tag_id", handlers.GetTransactionsByTag)
		v1.POST("/transactions/purge", handler.ValidateRequest[model.PurgeTransactionsRequest](), handlers.PurgeSoftDeletedTransactions)
		
		// Tag routes with validation
		v1.POST("/tags", handler.ValidateRequest[model.CreateTagRequest](), handlers.CreateTag)
		v1.GET("/tags", handlers.GetTags)
		
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
func healthHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Health check requested")
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"status":    "healthy",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "1.0.0",
			},
			"error": nil,
		})
	}
} 