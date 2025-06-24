package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/piotrzalecki/budget-api/internal/handler"
	"github.com/piotrzalecki/budget-api/internal/repo"
	"github.com/piotrzalecki/budget-api/pkg/model"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize database connection
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "dev.db" // Default to dev.db in current directory
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal("Failed to open database", zap.Error(err))
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	// Initialize repository
	repository := repo.NewRepository(db)

	// Initialize handlers with dependencies
	handlers := handler.NewHandler(repository)

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware(logger))

	// Setup routes
	setupRoutes(router, logger, handlers)

	// Create HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func setupRoutes(router *gin.Engine, logger *zap.Logger, handlers *handler.Handler) {
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
		
		// Recurring routes (TODO: Add validation when handlers are implemented)
		// v1.POST("/recurring", handler.ValidateRequest[model.CreateRecurringRequest](), handlers.CreateRecurring)
		// v1.PATCH("/recurring/:id", handler.ValidateRequest[model.UpdateRecurringRequest](), handlers.UpdateRecurring)
		// v1.GET("/recurring", handlers.GetRecurring)
		
		// Reports routes (TODO: Add when handlers are implemented)
		// v1.GET("/reports/monthly", handlers.GetMonthlyReport)
		
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
		// TODO: Add scheduler endpoint
		
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

func loggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	})
} 