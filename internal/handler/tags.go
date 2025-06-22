package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piotrzalecki/budget-api/pkg/model"
)

// CreateTag handles POST /api/v1/tags
func (h *Handler) CreateTag(c *gin.Context) {
	// Get the validated request from context
	request, ok := GetValidatedRequest[model.CreateTagRequest](c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get validated request",
			"data":  nil,
		})
		return
	}

	// Create tag using the repository
	tag, err := h.repo.CreateTag(c.Request.Context(), request.Name)
	if err != nil {
		// Handle specific error cases
		// TODO: Add proper error handling for duplicate names, etc.
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create tag: " + err.Error(),
			"data":  nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id": tag.ID,
		},
		"error": nil,
	})
}

// GetTags handles GET /api/v1/tags
func (h *Handler) GetTags(c *gin.Context) {
	// Get tags from the repository
	tags, err := h.repo.ListTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get tags: " + err.Error(),
			"data":  nil,
		})
		return
	}

	// Convert database models to response DTOs
	tagResponses := make([]model.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = model.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tagResponses,
		"error": nil,
	})
} 