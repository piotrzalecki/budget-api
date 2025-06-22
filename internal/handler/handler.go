package handler

import (
	"github.com/piotrzalecki/budget-api/internal/repo"
)

// Handler holds all dependencies needed by HTTP handlers
type Handler struct {
	repo repo.Repository
}

// NewHandler creates a new Handler instance with the given dependencies
func NewHandler(repo repo.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// GetRepository returns the repository instance
func (h *Handler) GetRepository() repo.Repository {
	return h.repo
} 