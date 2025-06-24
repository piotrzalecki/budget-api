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

// CreateTransaction handles POST /api/v1/transactions
func (h *Handler) CreateTransaction(c *gin.Context) {
	// Get the validated request from context
	request, ok := GetValidatedRequest[model.CreateTransactionRequest](c)
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

	// Parse the date
	tDate, err := model.ParseDate(request.TDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date format",
			"data":  nil,
		})
		return
	}

	// TODO: Get user ID from context when authentication is implemented
	// For now, use a default user ID of 1
	userID := int64(1)

	// Create transaction parameters
	params := repo.CreateTransactionParams{
		UserID:          userID,
		AmountPence:     amountPence,
		TDate:           tDate,
		Note:            model.StringToSQLNullString(request.Note),
		SourceRecurring: sql.NullInt64{Valid: false}, // Manual transaction
	}

	// Create transaction in database
	transaction, err := h.repo.CreateTransaction(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create transaction",
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

			// Create transaction-tag association
			tagParams := repo.CreateTransactionTagParams{
				TransactionID: transaction.ID,
				TagID:         tagID,
			}
			err = h.repo.CreateTransactionTag(c.Request.Context(), tagParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to associate tag with transaction",
					"data":  nil,
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id": transaction.ID,
		},
		"error": nil,
	})
}

// GetTransactions handles GET /api/v1/transactions
func (h *Handler) GetTransactions(c *gin.Context) {
	// Get query parameters
	from := c.Query("from")
	to := c.Query("to")

	// TODO: Get user ID from context when authentication is implemented
	// For now, use a default user ID of 1
	userID := int64(1)

	var transactions []repo.Transaction
	var err error

	if from != "" && to != "" {
		// Parse date range
		fromDate, err := model.ParseDate(from)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid from date format",
				"data":  nil,
			})
			return
		}

		toDate, err := model.ParseDate(to)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid to date format",
				"data":  nil,
			})
			return
		}

		// Use date range query with proper parameters
		params := repo.ListTransactionsParams{
			UserID:  userID,
			TDate:   fromDate,
			Column3: nil, // This represents the "OR ? IS NULL" condition
			TDate_2: toDate,
			Column5: nil, // This represents the "OR ? IS NULL" condition
		}
		transactions, err = h.repo.ListTransactions(c.Request.Context(), params)
	} else {
		// Get all transactions for user (no date filtering)
		// Use a very wide date range to get all transactions
		params := repo.ListTransactionsParams{
			UserID:  userID,
			TDate:   time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), // Very old date
			Column3: nil,
			TDate_2: time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC), // Very future date
			Column5: nil,
		}
		transactions, err = h.repo.ListTransactions(c.Request.Context(), params)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch transactions",
			"data":  nil,
		})
		return
	}

	// Convert to response DTOs
	response := make([]model.TransactionResponse, len(transactions))
	for i, txn := range transactions {
		// Get tags for this transaction
		tags, err := h.repo.GetTransactionTags(c.Request.Context(), txn.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch transaction tags",
				"data":  nil,
			})
			return
		}

		// Convert tag IDs
		tagIDs := make([]int64, len(tags))
		for j, tag := range tags {
			tagIDs[j] = tag.ID
		}

		response[i] = model.TransactionResponse{
			ID:             txn.ID,
			Amount:         model.PenceToCurrency(txn.AmountPence),
			TDate:          model.FormatDate(txn.TDate),
			Note:           model.SQLNullStringToString(txn.Note),
			CreatedAt:      txn.CreatedAt.Time,
			SourceRecurring: model.SQLNullInt64ToInt64(txn.SourceRecurring),
			DeletedAt:      model.SQLNullTimeToTimePtr(txn.DeletedAt),
			TagIDs:         tagIDs,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
}

// UpdateTransaction handles PATCH /api/v1/transactions/{id}
func (h *Handler) UpdateTransaction(c *gin.Context) {
	// Get transaction ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction ID",
			"data":  nil,
		})
		return
	}

	// Get the validated request from context
	request, ok := GetValidatedRequest[model.UpdateTransactionRequest](c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get validated request",
			"data":  nil,
		})
		return
	}

	// Check if transaction exists
	transaction, err := h.repo.GetTransactionByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "transaction not found",
				"data":  nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch transaction",
			"data":  nil,
		})
		return
	}

	// Handle soft delete if requested
	if request.Deleted != nil && *request.Deleted {
		err = h.repo.SoftDeleteTransaction(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete transaction",
				"data":  nil,
			})
			return
		}
		c.Status(http.StatusNoContent)
		return
	}

	// Update transaction fields if provided
	updateParams := repo.UpdateTransactionParams{
		ID:          id,
		AmountPence: transaction.AmountPence, // Keep existing amount
		TDate:       transaction.TDate,       // Keep existing date
		Note:        transaction.Note,        // Keep existing note
	}

	// Update note if provided
	if request.Note != nil {
		updateParams.Note = model.StringToSQLNullString(request.Note)
	}

	// Update transaction
	_, err = h.repo.UpdateTransaction(c.Request.Context(), updateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update transaction",
			"data":  nil,
		})
		return
	}

	// Handle tag associations if provided
	if request.TagIDs != nil {
		// Remove existing tags
		err = h.repo.DeleteAllTransactionTags(c.Request.Context(), id)
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

			// Create transaction-tag association
			tagParams := repo.CreateTransactionTagParams{
				TransactionID: id,
				TagID:         tagID,
			}
			err = h.repo.CreateTransactionTag(c.Request.Context(), tagParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to associate tag with transaction",
					"data":  nil,
				})
				return
			}
		}
	}

	c.Status(http.StatusNoContent)
}

// GetTransactionByID handles GET /api/v1/transactions/{id}
func (h *Handler) GetTransactionByID(c *gin.Context) {
	// Get transaction ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction ID",
			"data":  nil,
		})
		return
	}

	// Get transaction from database
	transaction, err := h.repo.GetTransactionByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "transaction not found",
				"data":  nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch transaction",
			"data":  nil,
		})
		return
	}

	// Get tags for this transaction
	tags, err := h.repo.GetTransactionTags(c.Request.Context(), transaction.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch transaction tags",
			"data":  nil,
		})
		return
	}

	// Convert tag IDs
	tagIDs := make([]int64, len(tags))
	for i, tag := range tags {
		tagIDs[i] = tag.ID
	}

	// Convert to response DTO
	response := model.TransactionResponse{
		ID:             transaction.ID,
		Amount:         model.PenceToCurrency(transaction.AmountPence),
		TDate:          model.FormatDate(transaction.TDate),
		Note:           model.SQLNullStringToString(transaction.Note),
		CreatedAt:      transaction.CreatedAt.Time,
		SourceRecurring: model.SQLNullInt64ToInt64(transaction.SourceRecurring),
		DeletedAt:      model.SQLNullTimeToTimePtr(transaction.DeletedAt),
		TagIDs:         tagIDs,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
}

// GetTransactionsByRecurringID handles GET /api/v1/transactions/by-recurring/{recurring_id}
func (h *Handler) GetTransactionsByRecurringID(c *gin.Context) {
	// Get recurring ID from URL
	recurringIDStr := c.Param("recurring_id")
	recurringID, err := strconv.ParseInt(recurringIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid recurring ID",
			"data":  nil,
		})
		return
	}

	// Get transactions by recurring ID
	sourceRecurring := sql.NullInt64{Int64: recurringID, Valid: true}
	transactions, err := h.repo.GetTransactionsByRecurringID(c.Request.Context(), sourceRecurring)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch transactions",
			"data":  nil,
		})
		return
	}

	// Convert to response DTOs
	response := make([]model.TransactionResponse, len(transactions))
	for i, txn := range transactions {
		// Get tags for this transaction
		tags, err := h.repo.GetTransactionTags(c.Request.Context(), txn.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch transaction tags",
				"data":  nil,
			})
			return
		}

		// Convert tag IDs
		tagIDs := make([]int64, len(tags))
		for j, tag := range tags {
			tagIDs[j] = tag.ID
		}

		response[i] = model.TransactionResponse{
			ID:             txn.ID,
			Amount:         model.PenceToCurrency(txn.AmountPence),
			TDate:          model.FormatDate(txn.TDate),
			Note:           model.SQLNullStringToString(txn.Note),
			CreatedAt:      txn.CreatedAt.Time,
			SourceRecurring: model.SQLNullInt64ToInt64(txn.SourceRecurring),
			DeletedAt:      model.SQLNullTimeToTimePtr(txn.DeletedAt),
			TagIDs:         tagIDs,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
}

// GetTransactionsByTag handles GET /api/v1/transactions/by-tag/{tag_id}
func (h *Handler) GetTransactionsByTag(c *gin.Context) {
	// Get tag ID from URL
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
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "tag not found",
				"data":  nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to verify tag",
			"data":  nil,
		})
		return
	}

	// Get transactions by tag
	transactions, err := h.repo.GetTransactionsByTag(c.Request.Context(), tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch transactions",
			"data":  nil,
		})
		return
	}

	// Convert to response DTOs
	response := make([]model.TransactionResponse, len(transactions))
	for i, txn := range transactions {
		// Get tags for this transaction
		tags, err := h.repo.GetTransactionTags(c.Request.Context(), txn.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch transaction tags",
				"data":  nil,
			})
			return
		}

		// Convert tag IDs
		tagIDs := make([]int64, len(tags))
		for j, tag := range tags {
			tagIDs[j] = tag.ID
		}

		response[i] = model.TransactionResponse{
			ID:             txn.ID,
			Amount:         model.PenceToCurrency(txn.AmountPence),
			TDate:          model.FormatDate(txn.TDate),
			Note:           model.SQLNullStringToString(txn.Note),
			CreatedAt:      txn.CreatedAt.Time,
			SourceRecurring: model.SQLNullInt64ToInt64(txn.SourceRecurring),
			DeletedAt:      model.SQLNullTimeToTimePtr(txn.DeletedAt),
			TagIDs:         tagIDs,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
}

// HardDeleteTransaction handles DELETE /api/v1/transactions/{id}
func (h *Handler) HardDeleteTransaction(c *gin.Context) {
	// Get transaction ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction ID",
			"data":  nil,
		})
		return
	}

	// Check if transaction exists
	_, err = h.repo.GetTransactionByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "transaction not found",
				"data":  nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch transaction",
			"data":  nil,
		})
		return
	}

	// Hard delete transaction
	err = h.repo.HardDeleteTransaction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete transaction",
			"data":  nil,
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// PurgeSoftDeletedTransactions handles POST /api/v1/transactions/purge
func (h *Handler) PurgeSoftDeletedTransactions(c *gin.Context) {
	// Get the validated request from context
	request, ok := GetValidatedRequest[model.PurgeTransactionsRequest](c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get validated request",
			"data":  nil,
		})
		return
	}

	// Parse the cutoff date
	cutoffDate, err := model.ParseDate(request.CutoffDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid cutoff date format",
			"data":  nil,
		})
		return
	}

	// Purge soft deleted transactions
	deletedAt := sql.NullTime{Time: cutoffDate, Valid: true}
	err = h.repo.PurgeSoftDeletedTransactions(c.Request.Context(), deletedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to purge transactions",
			"data":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "soft deleted transactions purged successfully",
		},
		"error": nil,
	})
} 