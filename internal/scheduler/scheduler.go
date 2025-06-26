package scheduler

import (
	"context"
	"database/sql"
	"time"

	"github.com/piotrzalecki/budget-api/internal/repo"
	"go.uber.org/zap"
)

// RunScheduler implements the scheduler logic from the specification
// It materializes recurring rules, purges soft-deleted transactions, and optionally performs backup
func RunScheduler(ctx context.Context, db *sql.DB, today time.Time, logger *zap.Logger) (int, error) {
	// Create repository instance
	repository := repo.NewRepository(db)
	
	// Use transaction to ensure atomicity
	var processed int
	err := repository.WithTx(ctx, func(txRepo repo.Repository) error {
		// Get rules due on or before today
		rules, err := txRepo.GetRecurringDueOnDate(ctx, today)
		if err != nil {
			return err
		}
		
		processed = len(rules)
		
		// Process each due rule
		for _, rule := range rules {
			// Check if rule has ended
			if rule.EndDate.Valid && rule.EndDate.Time.Before(today) {
				// Rule has ended, deactivate it
				err := txRepo.ToggleRecurringActive(ctx, rule.ID)
				if err != nil {
					return err
				}
				continue
			}
			
			// Create transaction from recurring rule
			transactionParams := repo.CreateTransactionParams{
				UserID:          rule.UserID,
				AmountPence:     rule.AmountPence,
				TDate:           rule.NextDueDate,
				Note:            rule.Description,
				SourceRecurring: sql.NullInt64{Int64: rule.ID, Valid: true},
			}
			
			_, err := txRepo.CreateTransaction(ctx, transactionParams)
			if err != nil {
				return err
			}
			
			// Copy tags from recurring rule to transaction
			tags, err := txRepo.GetRecurringTags(ctx, rule.ID)
			if err != nil {
				return err
			}
			
			// Get the transaction we just created to get its ID
			// We'll need to get it by the recurring source and date
			transactions, err := txRepo.GetTransactionsByRecurringID(ctx, sql.NullInt64{Int64: rule.ID, Valid: true})
			if err != nil {
				return err
			}
			
			// Find the transaction we just created (should be the most recent one)
			var transactionID int64
			for _, tx := range transactions {
				if tx.TDate.Equal(rule.NextDueDate) {
					transactionID = tx.ID
					break
				}
			}
			
			// Add tags to the transaction
			for _, tag := range tags {
				tagParams := repo.CreateTransactionTagParams{
					TransactionID: transactionID,
					TagID:         tag.ID,
				}
				err := txRepo.CreateTransactionTag(ctx, tagParams)
				if err != nil {
					return err
				}
			}
			
			// Calculate next due date
			nextDueDate := calculateNextDueDate(rule, today)
			
			// Update recurring rule with new next due date
			updateParams := repo.UpdateRecurringNextDueParams{
				NextDueDate: nextDueDate,
				ID:          rule.ID,
			}
			err = txRepo.UpdateRecurringNextDue(ctx, updateParams)
			if err != nil {
				return err
			}
		}
		
		// Purge soft-deleted transactions older than 30 days
		cutoffDate := today.AddDate(0, 0, -30)
		purgeParams := sql.NullTime{Time: cutoffDate, Valid: true}
		err = txRepo.PurgeSoftDeletedTransactions(ctx, purgeParams)
		if err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		return 0, err
	}
	
	// Log the scheduler run
	logger.Info("scheduler", zap.Int("processed", processed))
	
	return processed, nil
}

// calculateNextDueDate calculates the next due date based on the recurring rule
func calculateNextDueDate(rule repo.Recurring, today time.Time) time.Time {
	nextDue := rule.NextDueDate
	
	switch rule.Frequency {
	case "daily":
		nextDue = nextDue.AddDate(0, 0, int(rule.IntervalN))
	case "weekly":
		nextDue = nextDue.AddDate(0, 0, 7*int(rule.IntervalN))
	case "monthly":
		nextDue = nextDue.AddDate(0, int(rule.IntervalN), 0)
	case "yearly":
		nextDue = nextDue.AddDate(int(rule.IntervalN), 0, 0)
	}
	
	return nextDue
} 