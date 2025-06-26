package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piotrzalecki/budget-api/internal/repo"
	"github.com/piotrzalecki/budget-api/pkg/model"
)

// CreateRecurring handles POST /api/v1/recurring
func (h *Handler) CreateRecurring(c *gin.Context) {
	// Get the validated request from context
	request, ok := GetValidatedRequest[model.CreateRecurringRequest](c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get validated request",
			"data":  nil,
		})
		return
	}

	// Convert amount from string to pence
	amountPence, err := model.CurrencyToPence(request.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid amount format",
			"data":  nil,
		})
		return
	}

	// Parse the first due date
	firstDueDate, err := model.ParseDate(request.FirstDueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid first_due_date format",
			"data":  nil,
		})
		return
	}

	// Parse the end date if provided
	var endDate sql.NullTime
	if request.EndDate != nil {
		parsedEndDate, err := model.ParseDate(*request.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid end_date format",
				"data":  nil,
			})
			return
		}
		endDate = sql.NullTime{Time: parsedEndDate, Valid: true}
	}

	// TODO: Get user ID from context when authentication is implemented
	// For now, use a default user ID of 1
	userID := int64(1)

	// Create recurring parameters
	params := repo.CreateRecurringParams{
		UserID:       userID,
		AmountPence:  amountPence,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Frequency:    request.Frequency,
		IntervalN:    int64(request.IntervalN),
		FirstDueDate: firstDueDate,
		NextDueDate:  firstDueDate, // Initially same as first due date
		EndDate:      endDate,
		Active:       true,
	}

	// Create recurring rule in database
	recurring, err := h.repo.CreateRecurring(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create recurring rule",
			"data":  nil,
		})
		return
	}

	// Handle tag associations if provided
	if len(request.TagIDs) > 0 {
		for _, tagID := range request.TagIDs {
			// Verify tag exists
			_, err := h.repo.GetTagByID(c.Request.Context(), tagID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid tag ID: " + strconv.FormatInt(tagID, 10),
					"data":  nil,
				})
				return
			}

			// Create recurring-tag association
			tagParams := repo.CreateRecurringTagParams{
				RecurringID: recurring.ID,
				TagID:       tagID,
			}
			err = h.repo.CreateRecurringTag(c.Request.Context(), tagParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to associate tag with recurring rule",
					"data":  nil,
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id": recurring.ID,
		},
		"error": nil,
	})
}

// GetRecurring handles GET /api/v1/recurring
func (h *Handler) GetRecurring(c *gin.Context) {
	// TODO: Get user ID from context when authentication is implemented
	// For now, use a default user ID of 1
	userID := int64(1)

	// Get all recurring rules for user
	recurringRules, err := h.repo.ListRecurring(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch recurring rules",
			"data":  nil,
		})
		return
	}

	// Convert to response DTOs
	response := make([]model.RecurringResponse, len(recurringRules))
	for i, rule := range recurringRules {
		// Get tags for this recurring rule
		tags, err := h.repo.GetRecurringTags(c.Request.Context(), rule.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch recurring rule tags",
				"data":  nil,
			})
			return
		}

		// Convert tag IDs
		tagIDs := make([]int64, len(tags))
		for j, tag := range tags {
			tagIDs[j] = tag.ID
		}

		// Convert end date
		var endDateStr *string
		if rule.EndDate.Valid {
			formatted := model.FormatDate(rule.EndDate.Time)
			endDateStr = &formatted
		}

		response[i] = model.RecurringResponse{
			ID:            rule.ID,
			Amount:        model.PenceToCurrency(rule.AmountPence),
			Description:   rule.Description.String,
			Frequency:     rule.Frequency,
			IntervalN:     int(rule.IntervalN),
			FirstDueDate:  model.FormatDate(rule.FirstDueDate),
			NextDueDate:   model.FormatDate(rule.NextDueDate),
			EndDate:       endDateStr,
			Active:        rule.Active,
			CreatedAt:     rule.CreatedAt.Time,
			TagIDs:        tagIDs,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
}

// GetRecurringByID handles GET /api/v1/recurring/:id
func (h *Handler) GetRecurringByID(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid recurring rule ID",
			"data":  nil,
		})
		return
	}

	// Get recurring rule by ID
	rule, err := h.repo.GetRecurringByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "recurring rule not found",
			"data":  nil,
		})
		return
	}

	// TODO: Check if user has access to this recurring rule when authentication is implemented
	// For now, allow access to any recurring rule

	// Get tags for this recurring rule
	tags, err := h.repo.GetRecurringTags(c.Request.Context(), rule.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch recurring rule tags",
			"data":  nil,
		})
		return
	}

	// Convert tag IDs
	tagIDs := make([]int64, len(tags))
	for i, tag := range tags {
		tagIDs[i] = tag.ID
	}

	// Convert end date
	var endDateStr *string
	if rule.EndDate.Valid {
		formatted := model.FormatDate(rule.EndDate.Time)
		endDateStr = &formatted
	}

	response := model.RecurringResponse{
		ID:            rule.ID,
		Amount:        model.PenceToCurrency(rule.AmountPence),
		Description:   rule.Description.String,
		Frequency:     rule.Frequency,
		IntervalN:     int(rule.IntervalN),
		FirstDueDate:  model.FormatDate(rule.FirstDueDate),
		NextDueDate:   model.FormatDate(rule.NextDueDate),
		EndDate:       endDateStr,
		Active:        rule.Active,
		CreatedAt:     rule.CreatedAt.Time,
		TagIDs:        tagIDs,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
}

// UpdateRecurring handles PATCH /api/v1/recurring/:id
func (h *Handler) UpdateRecurring(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid recurring rule ID",
			"data":  nil,
		})
		return
	}

	// Get the validated request from context
	request, ok := GetValidatedRequest[model.UpdateRecurringRequest](c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get validated request",
			"data":  nil,
		})
		return
	}

	// Get existing recurring rule
	existingRule, err := h.repo.GetRecurringByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "recurring rule not found",
			"data":  nil,
		})
		return
	}

	// TODO: Check if user has access to this recurring rule when authentication is implemented

	// Prepare update parameters
	updateParams := repo.UpdateRecurringParams{
		ID:           id,
		AmountPence:  existingRule.AmountPence,
		Description:  existingRule.Description,
		Frequency:    existingRule.Frequency,
		IntervalN:    existingRule.IntervalN,
		FirstDueDate: existingRule.FirstDueDate,
		NextDueDate:  existingRule.NextDueDate,
		EndDate:      existingRule.EndDate,
		Active:       existingRule.Active,
	}

	// Update fields if provided
	if request.Amount != nil {
		amountPence, err := model.CurrencyToPence(*request.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid amount format",
				"data":  nil,
			})
			return
		}
		updateParams.AmountPence = amountPence
	}

	if request.Description != nil {
		updateParams.Description = sql.NullString{String: *request.Description, Valid: true}
	}

	if request.Frequency != nil {
		updateParams.Frequency = *request.Frequency
	}

	if request.IntervalN != nil {
		updateParams.IntervalN = int64(*request.IntervalN)
	}

	if request.FirstDueDate != nil {
		firstDueDate, err := model.ParseDate(*request.FirstDueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid first_due_date format",
				"data":  nil,
			})
			return
		}
		updateParams.FirstDueDate = firstDueDate
	}

	if request.EndDate != nil {
		if *request.EndDate == "" {
			updateParams.EndDate = sql.NullTime{Valid: false}
		} else {
			endDate, err := model.ParseDate(*request.EndDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid end_date format",
					"data":  nil,
				})
				return
			}
			updateParams.EndDate = sql.NullTime{Time: endDate, Valid: true}
		}
	}

	if request.Active != nil {
		updateParams.Active = *request.Active
	}

	// Update recurring rule
	_, err = h.repo.UpdateRecurring(c.Request.Context(), updateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update recurring rule",
			"data":  nil,
		})
		return
	}

	// Handle tag associations if provided
	if request.TagIDs != nil {
		// Delete existing tags
		err = h.repo.DeleteAllRecurringTags(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to remove existing tags",
				"data":  nil,
			})
			return
		}

		// Add new tags
		for _, tagID := range request.TagIDs {
			// Verify tag exists
			_, err := h.repo.GetTagByID(c.Request.Context(), tagID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid tag ID: " + strconv.FormatInt(tagID, 10),
					"data":  nil,
				})
				return
			}

			// Create recurring-tag association
			tagParams := repo.CreateRecurringTagParams{
				RecurringID: id,
				TagID:       tagID,
			}
			err = h.repo.CreateRecurringTag(c.Request.Context(), tagParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to associate tag with recurring rule",
					"data":  nil,
				})
				return
			}
		}
	}

	c.JSON(http.StatusNoContent, gin.H{
		"data":  nil,
		"error": nil,
	})
}

// DeleteRecurring handles DELETE /api/v1/recurring/:id
func (h *Handler) DeleteRecurring(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid recurring rule ID",
			"data":  nil,
		})
		return
	}

	// Check if recurring rule exists
	_, err = h.repo.GetRecurringByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "recurring rule not found",
			"data":  nil,
		})
		return
	}

	// TODO: Check if user has access to this recurring rule when authentication is implemented

	// Delete all associated tags first
	err = h.repo.DeleteAllRecurringTags(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove associated tags",
			"data":  nil,
		})
		return
	}

	// Delete the recurring rule
	err = h.repo.DeleteRecurring(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete recurring rule",
			"data":  nil,
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"data":  nil,
		"error": nil,
	})
}

// GetRecurringByTag handles GET /api/v1/recurring/by-tag/:tag_id
func (h *Handler) GetRecurringByTag(c *gin.Context) {
	// Parse tag ID from URL
	tagIDStr := c.Param("tag_id")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid tag ID",
			"data":  nil,
		})
		return
	}

	// Verify tag exists
	_, err = h.repo.GetTagByID(c.Request.Context(), tagID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tag not found",
			"data":  nil,
		})
		return
	}

	// Get recurring rules by tag
	recurringRules, err := h.repo.GetRecurringByTag(c.Request.Context(), tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch recurring rules by tag",
			"data":  nil,
		})
		return
	}

	// Convert to response DTOs
	response := make([]model.RecurringResponse, len(recurringRules))
	for i, rule := range recurringRules {
		// Get tags for this recurring rule
		tags, err := h.repo.GetRecurringTags(c.Request.Context(), rule.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch recurring rule tags",
				"data":  nil,
			})
			return
		}

		// Convert tag IDs
		tagIDs := make([]int64, len(tags))
		for j, tag := range tags {
			tagIDs[j] = tag.ID
		}

		// Convert end date
		var endDateStr *string
		if rule.EndDate.Valid {
			formatted := model.FormatDate(rule.EndDate.Time)
			endDateStr = &formatted
		}

		response[i] = model.RecurringResponse{
			ID:            rule.ID,
			Amount:        model.PenceToCurrency(rule.AmountPence),
			Description:   rule.Description.String,
			Frequency:     rule.Frequency,
			IntervalN:     int(rule.IntervalN),
			FirstDueDate:  model.FormatDate(rule.FirstDueDate),
			NextDueDate:   model.FormatDate(rule.NextDueDate),
			EndDate:       endDateStr,
			Active:        rule.Active,
			CreatedAt:     rule.CreatedAt.Time,
			TagIDs:        tagIDs,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
}

// ListActiveRecurring handles GET /api/v1/recurring/active
func (h *Handler) ListActiveRecurring(c *gin.Context) {
	// TODO: Get user ID from context when authentication is implemented
	// For now, use a default user ID of 1
	userID := int64(1)

	// Get active recurring rules for user
	recurringRules, err := h.repo.ListActiveRecurring(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch active recurring rules",
			"data":  nil,
		})
		return
	}

	// Convert to response DTOs
	response := make([]model.RecurringResponse, len(recurringRules))
	for i, rule := range recurringRules {
		// Get tags for this recurring rule
		tags, err := h.repo.GetRecurringTags(c.Request.Context(), rule.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch recurring rule tags",
				"data":  nil,
			})
			return
		}

		// Convert tag IDs
		tagIDs := make([]int64, len(tags))
		for j, tag := range tags {
			tagIDs[j] = tag.ID
		}

		// Convert end date
		var endDateStr *string
		if rule.EndDate.Valid {
			formatted := model.FormatDate(rule.EndDate.Time)
			endDateStr = &formatted
		}

		response[i] = model.RecurringResponse{
			ID:            rule.ID,
			Amount:        model.PenceToCurrency(rule.AmountPence),
			Description:   rule.Description.String,
			Frequency:     rule.Frequency,
			IntervalN:     int(rule.IntervalN),
			FirstDueDate:  model.FormatDate(rule.FirstDueDate),
			NextDueDate:   model.FormatDate(rule.NextDueDate),
			EndDate:       endDateStr,
			Active:        rule.Active,
			CreatedAt:     rule.CreatedAt.Time,
			TagIDs:        tagIDs,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
}

// ToggleRecurringActive handles PATCH /api/v1/recurring/:id/toggle
func (h *Handler) ToggleRecurringActive(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid recurring rule ID",
			"data":  nil,
		})
		return
	}

	// Check if recurring rule exists
	_, err = h.repo.GetRecurringByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "recurring rule not found",
			"data":  nil,
		})
		return
	}

	// TODO: Check if user has access to this recurring rule when authentication is implemented

	// Toggle the active status
	err = h.repo.ToggleRecurringActive(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to toggle recurring rule status",
			"data":  nil,
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"data":  nil,
		"error": nil,
	})
}

// GetRecurringDueOnDate handles GET /api/v1/recurring/due?date=YYYY-MM-DD
func (h *Handler) GetRecurringDueOnDate(c *gin.Context) {
	// Get date from query parameter
	dateStr := c.Query("date")
	if dateStr == "" {
		// If no date provided, use today's date
		dateStr = model.FormatDate(time.Now())
	}

	// Parse the date
	dueDate, err := model.ParseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date format. Use YYYY-MM-DD",
			"data":  nil,
		})
		return
	}

	// Get recurring rules due on the specified date
	recurringRules, err := h.repo.GetRecurringDueOnDate(c.Request.Context(), dueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch recurring rules due on date",
			"data":  nil,
		})
		return
	}

	// Convert to response DTOs
	response := make([]model.RecurringResponse, len(recurringRules))
	for i, rule := range recurringRules {
		// Get tags for this recurring rule
		tags, err := h.repo.GetRecurringTags(c.Request.Context(), rule.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch recurring rule tags",
				"data":  nil,
			})
			return
		}

		// Convert tag IDs
		tagIDs := make([]int64, len(tags))
		for j, tag := range tags {
			tagIDs[j] = tag.ID
		}

		// Convert end date
		var endDateStr *string
		if rule.EndDate.Valid {
			formatted := model.FormatDate(rule.EndDate.Time)
			endDateStr = &formatted
		}

		response[i] = model.RecurringResponse{
			ID:            rule.ID,
			Amount:        model.PenceToCurrency(rule.AmountPence),
			Description:   rule.Description.String,
			Frequency:     rule.Frequency,
			IntervalN:     int(rule.IntervalN),
			FirstDueDate:  model.FormatDate(rule.FirstDueDate),
			NextDueDate:   model.FormatDate(rule.NextDueDate),
			EndDate:       endDateStr,
			Active:        rule.Active,
			CreatedAt:     rule.CreatedAt.Time,
			TagIDs:        tagIDs,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
} 