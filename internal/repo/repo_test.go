package repo

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
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

func TestNewRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	assert.NotNil(t, repo)
	
	// Test that it implements the Repository interface
	var _ Repository = repo
}

func TestWithTx_Commit(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	// Create a test user
	user := CreateUserParams{
		Email:  "test@example.com",
		PwHash: "hashedpassword",
	}

	// Test successful transaction
	err := repo.WithTx(context.Background(), func(txRepo Repository) error {
		// Create user within transaction
		_, err := txRepo.CreateUser(context.Background(), user)
		return err
	})

	require.NoError(t, err)

	// Verify user was created (transaction committed)
	createdUser, err := repo.GetUserByEmail(context.Background(), "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.Email, createdUser.Email)
}

func TestWithTx_Rollback(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	// Create a test user
	user := CreateUserParams{
		Email:  "test@example.com",
		PwHash: "hashedpassword",
	}

	// Test transaction rollback
	err := repo.WithTx(context.Background(), func(txRepo Repository) error {
		// Create user within transaction
		_, err := txRepo.CreateUser(context.Background(), user)
		if err != nil {
			return err
		}

		// Try to create the same user again (should fail due to unique constraint)
		_, err = txRepo.CreateUser(context.Background(), user)
		return err // This should cause rollback
	})

	require.Error(t, err)

	// Verify user was not created (transaction rolled back)
	_, err = repo.GetUserByEmail(context.Background(), "test@example.com")
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestWithTx_NestedOperations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	// Test complex transaction with multiple operations
	err := repo.WithTx(context.Background(), func(txRepo Repository) error {
		// Create user
		user, err := txRepo.CreateUser(context.Background(), CreateUserParams{
			Email:  "test@example.com",
			PwHash: "hashedpassword",
		})
		if err != nil {
			return err
		}

		// Create tag
		tag, err := txRepo.CreateTag(context.Background(), "test-tag")
		if err != nil {
			return err
		}

		// Create transaction
		txn, err := txRepo.CreateTransaction(context.Background(), CreateTransactionParams{
			UserID:      user.ID,
			AmountPence: 1000,
			TDate:       time.Now(),
			Note:        sql.NullString{String: "Test transaction", Valid: true},
		})
		if err != nil {
			return err
		}

		// Link transaction to tag
		err = txRepo.CreateTransactionTag(context.Background(), CreateTransactionTagParams{
			TransactionID: txn.ID,
			TagID:         tag.ID,
		})
		return err
	})

	require.NoError(t, err)

	// Verify all operations were committed
	user, err := repo.GetUserByEmail(context.Background(), "test@example.com")
	require.NoError(t, err)

	tag, err := repo.GetTagByName(context.Background(), "test-tag")
	require.NoError(t, err)

	// Get transactions for the user
	txns, err := repo.ListTransactions(context.Background(), ListTransactionsParams{
		UserID:  user.ID,
		TDate:   time.Now(),
		Column3: nil,
		TDate_2: time.Now(),
		Column5: nil,
	})
	require.NoError(t, err)
	assert.Len(t, txns, 1)

	// Verify transaction tag relationship
	tags, err := repo.GetTransactionTags(context.Background(), txns[0].ID)
	require.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, tag.ID, tags[0].ID)
}

func TestWithTx_RollbackOnPanic(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	// Test that panic causes rollback
	var panicOccurred bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicOccurred = true
			}
		}()

		_ = repo.WithTx(context.Background(), func(txRepo Repository) error {
			// Create user within transaction
			_, err := txRepo.CreateUser(context.Background(), CreateUserParams{
				Email:  "test@example.com",
				PwHash: "hashedpassword",
			})
			if err != nil {
				return err
			}

			// Panic to test rollback
			panic("test panic")
		})
	}()

	assert.True(t, panicOccurred, "Expected panic to occur")

	// Verify user was not created (transaction rolled back)
	_, err := repo.GetUserByEmail(context.Background(), "test@example.com")
	assert.Error(t, err)
	// Note: The exact error type may vary depending on the database state after panic
	// We just check that the user doesn't exist
}

func TestRepository_ListTags(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	// Create some test tags
	_, err := repo.CreateTag(context.Background(), "tag1")
	require.NoError(t, err)

	_, err = repo.CreateTag(context.Background(), "tag2")
	require.NoError(t, err)

	// List all tags
	tags, err := repo.ListTags(context.Background())
	require.NoError(t, err)
	assert.Len(t, tags, 2)

	// Verify tags are returned
	tagNames := make(map[string]bool)
	for _, tag := range tags {
		tagNames[tag.Name] = true
	}
	assert.True(t, tagNames["tag1"])
	assert.True(t, tagNames["tag2"])
}

func TestRepository_CreateTransaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	// Create a user first
	user, err := repo.CreateUser(context.Background(), CreateUserParams{
		Email:  "test@example.com",
		PwHash: "hashedpassword",
	})
	require.NoError(t, err)

	// Create a transaction
	txn, err := repo.CreateTransaction(context.Background(), CreateTransactionParams{
		UserID:      user.ID,
		AmountPence: 1234, // £12.34
		TDate:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Note:        sql.NullString{String: "Test transaction", Valid: true},
	})
	require.NoError(t, err)

	// Verify transaction was created
	assert.Equal(t, user.ID, txn.UserID)
	assert.Equal(t, int64(1234), txn.AmountPence)
	assert.Equal(t, "Test transaction", txn.Note.String)
	assert.True(t, txn.Note.Valid)

	// Retrieve the transaction
	retrievedTxn, err := repo.GetTransactionByID(context.Background(), txn.ID)
	require.NoError(t, err)
	assert.Equal(t, txn.ID, retrievedTxn.ID)
	assert.Equal(t, txn.AmountPence, retrievedTxn.AmountPence)
} 