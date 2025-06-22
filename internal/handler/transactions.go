package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	// TODO: Implement actual business logic here using request and h.repo
	// This would typically involve:
	// 1. Converting amount from string to pence
	// 2. Parsing the date
	// 3. Creating transaction with h.repo.CreateTransaction()
	// 4. Handling tag associations if provided
	
	_ = request // Use request in actual implementation
	_ = h.repo  // Use repository in actual implementation
	
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id": 123,
		},
		"error": nil,
	})
}

// GetTransactions handles GET /api/v1/transactions
func (h *Handler) GetTransactions(c *gin.Context) {
	// Get query parameters
	from := c.Query("from")
	to := c.Query("to")

	// TODO: Validate date parameters and implement business logic using h.repo
	// This would typically involve:
	// 1. Parsing from/to dates
	// 2. Getting user ID from context (when auth is implemented)
	// 3. Calling h.repo.ListTransactionsByDateRange() or similar
	// 4. Converting database models to response DTOs
	
	_ = from    // Use from in actual implementation
	_ = to      // Use to in actual implementation
	_ = h.repo  // Use repository in actual implementation

	// Mock response for now
	transactions := []model.TransactionResponse{
		{
			ID:     1,
			Amount: "-12.34",
			TDate:  "2025-06-17",
			Note:   nil,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  transactions,
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

	// TODO: Implement actual business logic here using id, request and h.repo
	// This would typically involve:
	// 1. Checking if transaction exists
	// 2. Updating fields based on request
	// 3. Handling soft delete if requested
	// 4. Updating tag associations if provided
	
	_ = id      // Use id in actual implementation
	_ = request // Use request in actual implementation
	_ = h.repo  // Use repository in actual implementation
	
	c.Status(http.StatusNoContent)
} 