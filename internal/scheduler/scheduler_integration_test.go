package scheduler

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/piotrzalecki/budget-api/internal/repo"
)

// setupTestDB creates an in-memory SQLite database and runs migrations
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	
	// Create in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Run migrations
	err = goose.SetDialect("sqlite3")
	require.NoError(t, err)

	err = goose.Up(db, "../../migrations")
	require.NoError(t, err)

	return db
}

// createTestUser creates a test user and returns the user ID
func createTestUser(t *testing.T, repository repo.Repository) int64 {
	t.Helper()
	
	user, err := repository.CreateUser(context.Background(), repo.CreateUserParams{
		Email:  "test@example.com",
		PwHash: "hashedpassword",
	})
	require.NoError(t, err)
	
	return user.ID
}

// createTestTag creates a test tag and returns the tag ID
func createTestTag(t *testing.T, repository repo.Repository) int64 {
	t.Helper()
	
	tag, err := repository.CreateTag(context.Background(), "test-tag")
	require.NoError(t, err)
	
	return tag.ID
}

// createRecurringRule creates a recurring rule with the given parameters
func createRecurringRule(t *testing.T, repository repo.Repository, userID int64, firstDueDate time.Time, frequency string, intervalN int64, amountPence int64) repo.Recurring {
	t.Helper()
	
	rule, err := repository.CreateRecurring(context.Background(), repo.CreateRecurringParams{
		UserID:       userID,
		AmountPence:  amountPence,
		Description:  sql.NullString{String: "Test recurring rule", Valid: true},
		Frequency:    frequency,
		IntervalN:    intervalN,
		FirstDueDate: firstDueDate,
		NextDueDate:  firstDueDate,
		EndDate:      sql.NullTime{Valid: false},
		Active:       true,
	})
	require.NoError(t, err)
	
	return rule
}

// addTagToRecurring adds a tag to a recurring rule
func addTagToRecurring(t *testing.T, repository repo.Repository, recurringID, tagID int64) {
	t.Helper()
	
	err := repository.CreateRecurringTag(context.Background(), repo.CreateRecurringTagParams{
		RecurringID: recurringID,
		TagID:       tagID,
	})
	require.NoError(t, err)
}

// assertTransactionExists verifies that a transaction exists with the expected properties
func assertTransactionExists(t *testing.T, repository repo.Repository, userID int64, amountPence int64, tDate time.Time, sourceRecurringID int64) {
	t.Helper()
	
	// Get transactions for the user
	transactions, err := repository.ListTransactions(context.Background(), repo.ListTransactionsParams{
		UserID:  userID,
		TDate:   tDate,
		Column3: nil,
		TDate_2: tDate,
		Column5: nil,
	})
	require.NoError(t, err)
	
	// Find the transaction we expect
	var found bool
	for _, tx := range transactions {
		if tx.AmountPence == amountPence && 
		   tx.TDate.Equal(tDate) && 
		   tx.SourceRecurring.Valid && 
		   tx.SourceRecurring.Int64 == sourceRecurringID {
			found = true
			break
		}
	}
	
	assert.True(t, found, "Expected transaction not found: amount=%d, date=%v, source_recurring=%d", amountPence, tDate, sourceRecurringID)
}

// assertRecurringNextDueDate verifies that a recurring rule has the expected next due date
func assertRecurringNextDueDate(t *testing.T, repository repo.Repository, recurringID int64, expectedNextDue time.Time) {
	t.Helper()
	
	rule, err := repository.GetRecurringByID(context.Background(), recurringID)
	require.NoError(t, err)
	
	assert.True(t, rule.NextDueDate.Equal(expectedNextDue), 
		"Expected next due date %v, got %v", expectedNextDue, rule.NextDueDate)
}

// assertTransactionHasTags verifies that a transaction has the expected tags
func assertTransactionHasTags(t *testing.T, repository repo.Repository, userID int64, amountPence int64, tDate time.Time, expectedTagNames []string) {
	t.Helper()
	
	// Get transactions for the user
	transactions, err := repository.ListTransactions(context.Background(), repo.ListTransactionsParams{
		UserID:  userID,
		TDate:   tDate,
		Column3: nil,
		TDate_2: tDate,
		Column5: nil,
	})
	require.NoError(t, err)
	
	// Find the transaction
	var transactionID int64
	for _, tx := range transactions {
		if tx.AmountPence == amountPence && tx.TDate.Equal(tDate) {
			transactionID = tx.ID
			break
		}
	}
	require.NotZero(t, transactionID, "Transaction not found")
	
	// Get tags for the transaction
	tags, err := repository.GetTransactionTags(context.Background(), transactionID)
	require.NoError(t, err)
	
	// Verify tag count
	assert.Len(t, tags, len(expectedTagNames), "Expected %d tags, got %d", len(expectedTagNames), len(tags))
	
	// Verify tag names
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	
	for _, expectedName := range expectedTagNames {
		assert.Contains(t, tagNames, expectedName, "Expected tag '%s' not found", expectedName)
	}
}

func TestSchedulerIntegration_BasicFlow(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	
	repository := repo.NewRepository(db)
	userID := createTestUser(t, repository)
	tagID := createTestTag(t, repository)
	
	// Create a recurring rule that was due yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	rule := createRecurringRule(t, repository, userID, yesterday, "daily", 1, 1000)
	addTagToRecurring(t, repository, rule.ID, tagID)
	
	// Verify initial state
	initialRule, err := repository.GetRecurringByID(context.Background(), rule.ID)
	require.NoError(t, err)
	assert.True(t, initialRule.NextDueDate.Equal(yesterday))
	
	// Run scheduler with today's date
	today := time.Now().Truncate(24 * time.Hour)
	logger := zap.NewNop()
	
	processed, err := RunScheduler(context.Background(), db, today, logger)
	require.NoError(t, err)
	assert.Equal(t, 1, processed)
	
	// Verify transaction was created
	assertTransactionExists(t, repository, userID, 1000, yesterday, rule.ID)
	
	// Verify transaction has the expected tag
	assertTransactionHasTags(t, repository, userID, 1000, yesterday, []string{"test-tag"})
	
	// Verify next due date was advanced
	expectedNextDue := yesterday.AddDate(0, 0, 1)
	assertRecurringNextDueDate(t, repository, rule.ID, expectedNextDue)
}

func TestSchedulerIntegration_WeeklyFrequency(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	
	repository := repo.NewRepository(db)
	userID := createTestUser(t, repository)
	
	// Create a weekly recurring rule that was due yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	rule := createRecurringRule(t, repository, userID, yesterday, "weekly", 1, 2000)
	
	// Run scheduler
	today := time.Now().Truncate(24 * time.Hour)
	logger := zap.NewNop()
	
	processed, err := RunScheduler(context.Background(), db, today, logger)
	require.NoError(t, err)
	assert.Equal(t, 1, processed)
	
	// Verify transaction was created
	assertTransactionExists(t, repository, userID, 2000, yesterday, rule.ID)
	
	// Verify next due date was advanced by 7 days
	expectedNextDue := yesterday.AddDate(0, 0, 7)
	assertRecurringNextDueDate(t, repository, rule.ID, expectedNextDue)
}

func TestSchedulerIntegration_MonthlyFrequency(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	
	repository := repo.NewRepository(db)
	userID := createTestUser(t, repository)
	
	// Create a monthly recurring rule that was due yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	rule := createRecurringRule(t, repository, userID, yesterday, "monthly", 1, 3000)
	
	// Run scheduler
	today := time.Now().Truncate(24 * time.Hour)
	logger := zap.NewNop()
	
	processed, err := RunScheduler(context.Background(), db, today, logger)
	require.NoError(t, err)
	assert.Equal(t, 1, processed)
	
	// Verify transaction was created
	assertTransactionExists(t, repository, userID, 3000, yesterday, rule.ID)
	
	// Verify next due date was advanced by 1 month
	expectedNextDue := yesterday.AddDate(0, 1, 0)
	assertRecurringNextDueDate(t, repository, rule.ID, expectedNextDue)
}

func TestSchedulerIntegration_YearlyFrequency(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	
	repository := repo.NewRepository(db)
	userID := createTestUser(t, repository)
	
	// Create a yearly recurring rule that was due yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	rule := createRecurringRule(t, repository, userID, yesterday, "yearly", 1, 4000)
	
	// Run scheduler
	today := time.Now().Truncate(24 * time.Hour)
	logger := zap.NewNop()
	
	processed, err := RunScheduler(context.Background(), db, today, logger)
	require.NoError(t, err)
	assert.Equal(t, 1, processed)
	
	// Verify transaction was created
	assertTransactionExists(t, repository, userID, 4000, yesterday, rule.ID)
	
	// Verify next due date was advanced by 1 year
	expectedNextDue := yesterday.AddDate(1, 0, 0)
	assertRecurringNextDueDate(t, repository, rule.ID, expectedNextDue)
}

func TestSchedulerIntegration_MultipleRules(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	
	repository := repo.NewRepository(db)
	userID := createTestUser(t, repository)
	
	// Create multiple recurring rules that were due yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	rule1 := createRecurringRule(t, repository, userID, yesterday, "daily", 1, 1000)
	rule2 := createRecurringRule(t, repository, userID, yesterday, "weekly", 1, 2000)
	rule3 := createRecurringRule(t, repository, userID, yesterday, "monthly", 1, 3000)
	
	// Run scheduler
	today := time.Now().Truncate(24 * time.Hour)
	logger := zap.NewNop()
	
	processed, err := RunScheduler(context.Background(), db, today, logger)
	require.NoError(t, err)
	assert.Equal(t, 3, processed)
	
	// Verify all transactions were created
	assertTransactionExists(t, repository, userID, 1000, yesterday, rule1.ID)
	assertTransactionExists(t, repository, userID, 2000, yesterday, rule2.ID)
	assertTransactionExists(t, repository, userID, 3000, yesterday, rule3.ID)
	
	// Verify next due dates were advanced correctly
	assertRecurringNextDueDate(t, repository, rule1.ID, yesterday.AddDate(0, 0, 1))
	assertRecurringNextDueDate(t, repository, rule2.ID, yesterday.AddDate(0, 0, 7))
	assertRecurringNextDueDate(t, repository, rule3.ID, yesterday.AddDate(0, 1, 0))
}

func TestSchedulerIntegration_NoRulesDue(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	
	repository := repo.NewRepository(db)
	userID := createTestUser(t, repository)
	
	// Create a recurring rule that is due tomorrow (not today)
	tomorrow := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	rule := createRecurringRule(t, repository, userID, tomorrow, "daily", 1, 1000)
	
	// Run scheduler
	today := time.Now().Truncate(24 * time.Hour)
	logger := zap.NewNop()
	
	processed, err := RunScheduler(context.Background(), db, today, logger)
	require.NoError(t, err)
	assert.Equal(t, 0, processed)
	
	// Verify no transaction was created
	transactions, err := repository.ListTransactions(context.Background(), repo.ListTransactionsParams{
		UserID:  userID,
		TDate:   today,
		Column3: nil,
		TDate_2: today,
		Column5: nil,
	})
	require.NoError(t, err)
	assert.Len(t, transactions, 0)
	
	// Verify next due date was not changed
	assertRecurringNextDueDate(t, repository, rule.ID, tomorrow)
}

func TestSchedulerIntegration_RuleEndDate(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	
	repository := repo.NewRepository(db)
	userID := createTestUser(t, repository)
	
	// Create a recurring rule that was due yesterday but ended yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	rule, err := repository.CreateRecurring(context.Background(), repo.CreateRecurringParams{
		UserID:       userID,
		AmountPence:  1000,
		Description:  sql.NullString{String: "Test recurring rule", Valid: true},
		Frequency:    "daily",
		IntervalN:    1,
		FirstDueDate: yesterday,
		NextDueDate:  yesterday,
		EndDate:      sql.NullTime{Time: yesterday, Valid: true},
		Active:       true,
	})
	require.NoError(t, err)
	
	// Run scheduler
	today := time.Now().Truncate(24 * time.Hour)
	logger := zap.NewNop()
	
	processed, err := RunScheduler(context.Background(), db, today, logger)
	require.NoError(t, err)
	assert.Equal(t, 1, processed) // Rule was processed (found as due)
	
	// Verify no transaction was created (rule was ended, not materialized)
	transactions, err := repository.ListTransactions(context.Background(), repo.ListTransactionsParams{
		UserID:  userID,
		TDate:   yesterday,
		Column3: nil,
		TDate_2: yesterday,
		Column5: nil,
	})
	require.NoError(t, err)
	assert.Len(t, transactions, 0, "No transaction should be created for ended rules")
	
	// Verify rule was deactivated
	updatedRule, err := repository.GetRecurringByID(context.Background(), rule.ID)
	require.NoError(t, err)
	assert.False(t, updatedRule.Active, "Rule should have been deactivated")
}

func TestSchedulerIntegration_DuplicatePrevention(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	
	repository := repo.NewRepository(db)
	userID := createTestUser(t, repository)
	
	// Create a recurring rule that was due yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	rule := createRecurringRule(t, repository, userID, yesterday, "daily", 1, 1000)
	
	// Run scheduler with yesterday's date (not today)
	logger := zap.NewNop()
	
	// First run
	processed, err := RunScheduler(context.Background(), db, yesterday, logger)
	require.NoError(t, err)
	assert.Equal(t, 1, processed)
	
	// Second run with same date - should handle duplicate gracefully
	processed, err = RunScheduler(context.Background(), db, yesterday, logger)
	require.NoError(t, err)
	assert.Equal(t, 0, processed, "Second run should process 0 rules (transaction already exists)")
	
	// Verify only one transaction exists
	transactions, err := repository.ListTransactions(context.Background(), repo.ListTransactionsParams{
		UserID:  userID,
		TDate:   yesterday,
		Column3: nil,
		TDate_2: yesterday,
		Column5: nil,
	})
	require.NoError(t, err)
	assert.Len(t, transactions, 1, "Should only have one transaction, not duplicates")
	
	// Verify the transaction was created from the correct recurring rule
	assertTransactionExists(t, repository, userID, 1000, yesterday, rule.ID)
}

func TestSchedulerIntegration_PurgeSoftDeleted(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	
	repository := repo.NewRepository(db)
	userID := createTestUser(t, repository)
	
	// Create a transaction and soft delete it
	txn, err := repository.CreateTransaction(context.Background(), repo.CreateTransactionParams{
		UserID:      userID,
		AmountPence: 1000,
		TDate:       time.Now().AddDate(0, 0, -31), // 31 days ago
		Note:        sql.NullString{String: "Test transaction", Valid: true},
	})
	require.NoError(t, err)
	
	// Soft delete the transaction
	err = repository.SoftDeleteTransaction(context.Background(), txn.ID)
	require.NoError(t, err)
	
	// Verify transaction is soft deleted by querying directly
	var deletedAt sql.NullTime
	err = db.QueryRowContext(context.Background(), "SELECT deleted_at FROM transactions WHERE id = ?", txn.ID).Scan(&deletedAt)
	require.NoError(t, err)
	assert.True(t, deletedAt.Valid, "Transaction should be soft deleted")
	
	// Manually set deleted_at to 31 days ago so the scheduler will purge it
	thirtyOneDaysAgo := time.Now().AddDate(0, 0, -31)
	_, err = db.ExecContext(context.Background(), "UPDATE transactions SET deleted_at = ? WHERE id = ?", thirtyOneDaysAgo, txn.ID)
	require.NoError(t, err)
	
	// Run scheduler
	today := time.Now().Truncate(24 * time.Hour)
	logger := zap.NewNop()
	
	processed, err := RunScheduler(context.Background(), db, today, logger)
	require.NoError(t, err)
	assert.Equal(t, 0, processed) // No recurring rules to process
	
	// Verify soft deleted transaction was purged (should not be found)
	_, err = repository.GetTransactionByID(context.Background(), txn.ID)
	require.Error(t, err, "Transaction should have been purged")
	require.Equal(t, sql.ErrNoRows, err, "Expected no rows error after purging")
} 