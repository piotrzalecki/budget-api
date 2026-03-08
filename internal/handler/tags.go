package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/piotrzalecki/budget-api/internal/repo"
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
		h.logger.Error("failed to create tag", zap.Error(err), zap.String("name", request.Name))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create tag: " + err.Error(),
			"data":  nil,
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"data":  model.TagResponse{ID: tag.ID, Name: tag.Name},
		"error": nil,
	})
}

// UpdateTag handles PATCH /api/v1/tags/:id
// @Summary Update a tag
// @Description Update an existing tag's name
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param tag body model.UpdateTagRequest true "Tag update data"
// @Success 200 {object} map[string]interface{} "Tag updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 404 {object} map[string]interface{} "Tag not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security ApiKeyAuth
// @Router /tags/{id} [patch]
func (h *Handler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid tag ID",
			"data":  nil,
		})
		return
	}

	request, ok := GetValidatedRequest[model.UpdateTagRequest](c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get validated request",
			"data":  nil,
		})
		return
	}

	_, err = h.repo.GetTagByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tag not found",
			"data":  nil,
		})
		return
	}

	tag, err := h.repo.UpdateTag(c.Request.Context(), repo.UpdateTagParams{ID: id, Name: request.Name})
	if err != nil {
		h.logger.Error("failed to update tag", zap.Error(err), zap.Int64("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update tag: " + err.Error(),
			"data":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  model.TagResponse{ID: tag.ID, Name: tag.Name},
		"error": nil,
	})
}

// DeleteTag handles DELETE /api/v1/tags/:id
// @Summary Delete a tag
// @Description Delete an existing tag
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Success 204 "Tag deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid tag ID"
// @Failure 404 {object} map[string]interface{} "Tag not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security ApiKeyAuth
// @Router /tags/{id} [delete]
func (h *Handler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid tag ID",
			"data":  nil,
		})
		return
	}

	_, err = h.repo.GetTagByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tag not found",
			"data":  nil,
		})
		return
	}

	err = h.repo.DeleteTag(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to delete tag", zap.Error(err), zap.Int64("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete tag: " + err.Error(),
			"data":  nil,
		})
		return
	}

	c.Status(http.StatusNoContent)
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
		h.logger.Error("failed to list tags", zap.Error(err))
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