package repo

import (
	"context"
	"database/sql"
)

// RepositoryImpl implements the Repository interface
type RepositoryImpl struct {
	*Queries
	db *sql.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) Repository {
	return &RepositoryImpl{
		Queries: New(db),
		db:      db,
	}
}

// WithTx executes a function within a database transaction
func (r *RepositoryImpl) WithTx(ctx context.Context, fn func(Repository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Create a new repository instance with the transaction
	txRepo := &RepositoryImpl{
		Queries: New(tx),
		db:      r.db, // Keep reference to original db for potential future use
	}

	// Execute the function with the transaction repository
	if err := fn(txRepo); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			// Return both errors, but prioritize the original error
			return err
		}
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// GetDB returns the underlying database connection
func (r *RepositoryImpl) GetDB() *sql.DB {
	return r.db
} 