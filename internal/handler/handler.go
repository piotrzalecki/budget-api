package handler

import (
	"github.com/piotrzalecki/budget-api/internal/repo"
	"go.uber.org/zap"
)

// Handler holds all dependencies needed by HTTP handlers
type Handler struct {
	repo   repo.Repository
	logger *zap.Logger
}

// NewHandler creates a new Handler instance with the given dependencies
func NewHandler(repo repo.Repository, logger *zap.Logger) *Handler {
	return &Handler{
		repo:   repo,
		logger: logger,
	}
}

// GetRepository returns the repository instance
func (h *Handler) GetRepository() repo.Repository {
	return h.repo
} 