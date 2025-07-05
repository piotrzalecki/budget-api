package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/piotrzalecki/budget-api/internal/scheduler"
	"github.com/piotrzalecki/budget-api/pkg/model"
)

// RunScheduler handles POST /admin/run-scheduler
// @Summary Run the scheduler
// @Description Manually trigger the scheduler to process recurring transactions due today
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Scheduler execution result"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security ApiKeyAuth
// @Router /admin/run-scheduler [post]
func (h *Handler) RunScheduler(c *gin.Context) {
	// Get logger from context or create a new one
	logger := zap.L()
	if log, exists := c.Get("logger"); exists {
		if l, ok := log.(*zap.Logger); ok {
			logger = l
		}
	}

	// Get database connection from repository
	db := h.repo.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database connection not available",
			"data":  nil,
		})
		return
	}

	// Run the scheduler with today's date
	today := time.Now().UTC().Truncate(24 * time.Hour)
	processed, err := scheduler.RunScheduler(c.Request.Context(), db, today, logger)
	if err != nil {
		logger.Error("scheduler failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "scheduler execution failed",
			"data":  nil,
		})
		return
	}

	// Return success response with processed count
	response := model.SchedulerResponse{
		Processed: processed,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
} 