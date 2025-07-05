package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piotrzalecki/budget-api/pkg/model"
)

// CreateTag handles POST /api/v1/tags
// @Summary Create a new tag
// @Description Create a new tag for categorizing transactions
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body model.CreateTagRequest true "Tag data"
// @Success 200 {object} map[string]interface{} "Tag created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security ApiKeyAuth
// @Router /tags [post]
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
// @Summary Get all tags
// @Description Get all available tags for the authenticated user
// @Tags tags
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of tags"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security ApiKeyAuth
// @Router /tags [get]
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