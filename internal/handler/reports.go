package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piotrzalecki/budget-api/internal/repo"
	"github.com/piotrzalecki/budget-api/pkg/model"
)

// GetMonthlyReport handles GET /api/v1/reports/monthly
func (h *Handler) GetMonthlyReport(c *gin.Context) {
	// Get query parameter for year-month (YYYY-MM format)
	ym := c.Query("ym")
	if ym == "" {
		// Default to current month if not provided
		now := time.Now()
		ym = now.Format("2006-01")
	}

	// Parse the year-month parameter
	yearMonth, err := time.Parse("2006-01", ym)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid year-month format. Use YYYY-MM (e.g., 2025-06)",
			"data":  nil,
		})
		return
	}

	// TODO: Get user ID from context when authentication is implemented
	// For now, use a default user ID of 1
	userID := int64(1)

	// Get monthly totals
	totalsParams := repo.GetMonthlyTotalsParams{
		UserID: userID,
		TDate:  yearMonth,
	}
	totals, err := h.repo.GetMonthlyTotals(c.Request.Context(), totalsParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch monthly totals",
			"data":  nil,
		})
		return
	}

	// Get monthly report by tag
	reportParams := repo.GetMonthlyReportParams{
		UserID: userID,
		TDate:  yearMonth,
	}
	reportRows, err := h.repo.GetMonthlyReport(c.Request.Context(), reportParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch monthly report",
			"data":  nil,
		})
		return
	}

	// Build response
	byTag := make(map[string]model.TagReportEntry)
	for _, row := range reportRows {
		tagName := "Untagged"
		if row.TagName.Valid {
			tagName = row.TagName.String
		}

		// Convert pence to currency strings
		totalIn := "0.00"
		if row.TotalInPence.Valid {
			totalIn = model.PenceToCurrency(int64(row.TotalInPence.Float64))
		}

		totalOut := "0.00"
		if row.TotalOutPence.Valid {
			totalOut = model.PenceToCurrency(int64(row.TotalOutPence.Float64))
		}

		byTag[tagName] = model.TagReportEntry{
			TotalIn:  totalIn,
			TotalOut: totalOut,
		}
	}

	// Convert totals to currency strings
	totalIn := "0.00"
	if totals.TotalInPence.Valid {
		totalIn = model.PenceToCurrency(int64(totals.TotalInPence.Float64))
	}

	totalOut := "0.00"
	if totals.TotalOutPence.Valid {
		totalOut = model.PenceToCurrency(int64(totals.TotalOutPence.Float64))
	}

	response := model.MonthlyReportResponse{
		TotalIn:  totalIn,
		TotalOut: totalOut,
		ByTag:    byTag,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
}

// GetMonthlyTotals handles GET /api/v1/reports/monthly/totals
func (h *Handler) GetMonthlyTotals(c *gin.Context) {
	// Get query parameter for year-month (YYYY-MM format)
	ym := c.Query("ym")
	if ym == "" {
		// Default to current month if not provided
		now := time.Now()
		ym = now.Format("2006-01")
	}

	// Parse the year-month parameter
	yearMonth, err := time.Parse("2006-01", ym)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid year-month format. Use YYYY-MM (e.g., 2025-06)",
			"data":  nil,
		})
		return
	}

	// TODO: Get user ID from context when authentication is implemented
	// For now, use a default user ID of 1
	userID := int64(1)

	// Get monthly totals
	params := repo.GetMonthlyTotalsParams{
		UserID: userID,
		TDate:  yearMonth,
	}
	totals, err := h.repo.GetMonthlyTotals(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch monthly totals",
			"data":  nil,
		})
		return
	}

	// Convert to currency strings
	totalIn := "0.00"
	if totals.TotalInPence.Valid {
		totalIn = model.PenceToCurrency(int64(totals.TotalInPence.Float64))
	}

	totalOut := "0.00"
	if totals.TotalOutPence.Valid {
		totalOut = model.PenceToCurrency(int64(totals.TotalOutPence.Float64))
	}

	response := gin.H{
		"total_in":         totalIn,
		"total_out":        totalOut,
		"transaction_count": totals.TransactionCount,
		"year_month":       ym,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"error": nil,
	})
} 