package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/piotrzalecki/budget-api/internal/docs" // This is the generated docs
	"github.com/piotrzalecki/budget-api/internal/handler"
	"github.com/piotrzalecki/budget-api/internal/repo"
)

// @title           Budget API
// @version         1.0
// @description     A REST API for personal budget management with transaction tracking, recurring payments, and reporting.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API key for authentication

// @tag.name transactions
// @tag.description Operations about transactions

// @tag.name tags
// @tag.description Operations about tags

// @tag.name recurring
// @tag.description Operations about recurring transactions

// @tag.name reports
// @tag.description Operations about reports

// @tag.name admin
// @tag.description Administrative operations

var version = "dev"

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

	// Seed service user if env vars are set
	seedServiceUser(context.Background(), repository, logger)

	// Initialize handlers with dependencies
	handlers := handler.NewHandler(repository, logger)

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add CORS middleware
	config := cors.DefaultConfig()
	// Allow requests from configured origins
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		// Default origins for development
		config.AllowOrigins = []string{
			"https://budget.lab.zalecki.uk",
			"http://budget.lab.zalecki.uk", // For development
		}
	} else {
		// Parse comma-separated origins from environment variable
		config.AllowOrigins = strings.Split(corsOrigins, ",")
		// Trim whitespace from each origin
		for i, origin := range config.AllowOrigins {
			config.AllowOrigins[i] = strings.TrimSpace(origin)
		}
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-API-Key", "Authorization"}
	config.AllowCredentials = false
	router.Use(cors.New(config))

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	// Setup routes
	setupRoutes(router, logger, handlers, repository, version)

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

// seedServiceUser creates (or updates) a permanent service user and session
// when SERVICE_USER_EMAIL and SERVICE_USER_TOKEN are set.
func seedServiceUser(ctx context.Context, r repo.Repository, logger *zap.Logger) {
	email := os.Getenv("SERVICE_USER_EMAIL")
	token := os.Getenv("SERVICE_USER_TOKEN")
	if email == "" || token == "" {
		return
	}

	user, err := r.GetUserByEmail(ctx, email)
	if err != nil {
		// User doesn't exist yet — create with a placeholder password hash
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
		if hashErr != nil {
			logger.Error("seedServiceUser: failed to hash token", zap.Error(hashErr))
			return
		}
		user, err = r.CreateUser(ctx, repo.CreateUserParams{
			Email:     email,
			PwHash:    string(hash),
			IsService: true,
		})
		if err != nil {
			logger.Error("seedServiceUser: failed to create user", zap.Error(err))
			return
		}
		logger.Info("seedServiceUser: created service user", zap.String("email", email))
	}

	// Upsert permanent session: delete old sessions, create a new permanent one
	if err := r.DeleteAllSessionsByUserID(ctx, user.ID); err != nil {
		logger.Error("seedServiceUser: failed to clear sessions", zap.Error(err))
		return
	}
	_, err = r.CreateSession(ctx, repo.CreateSessionParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: sql.NullTime{Valid: false}, // permanent
	})
	if err != nil {
		logger.Error("seedServiceUser: failed to create session", zap.Error(err))
		return
	}
	logger.Info("seedServiceUser: seeded permanent session", zap.String("email", email))
} 