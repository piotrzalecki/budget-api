package repo

import (
	"context"
	"database/sql"
	"time"
)

// Repository defines the interface for all database operations
type Repository interface {
	// Transaction management
	WithTx(ctx context.Context, fn func(Repository) error) error

	// User operations
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	DeleteUser(ctx context.Context, id int64) error

	// Transaction operations
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	GetTransactionByID(ctx context.Context, id int64) (Transaction, error)
	ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]Transaction, error)
	ListTransactionsByDateRange(ctx context.Context, userID int64) ([]Transaction, error)
	GetTransactionsByRecurringID(ctx context.Context, sourceRecurring sql.NullInt64) ([]Transaction, error)
	GetTransactionsByTag(ctx context.Context, tagID int64) ([]Transaction, error)
	UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) (Transaction, error)
	SoftDeleteTransaction(ctx context.Context, id int64) error
	HardDeleteTransaction(ctx context.Context, id int64) error
	PurgeSoftDeletedTransactions(ctx context.Context, deletedAt sql.NullTime) error

	// Tag operations
	CreateTag(ctx context.Context, name string) (Tag, error)
	GetTagByID(ctx context.Context, id int64) (Tag, error)
	GetTagByName(ctx context.Context, name string) (Tag, error)
	ListTags(ctx context.Context) ([]Tag, error)
	UpdateTag(ctx context.Context, arg UpdateTagParams) (Tag, error)
	DeleteTag(ctx context.Context, id int64) error

	// Transaction tag operations
	CreateTransactionTag(ctx context.Context, arg CreateTransactionTagParams) error
	GetTransactionTags(ctx context.Context, transactionID int64) ([]Tag, error)
	DeleteTransactionTag(ctx context.Context, arg DeleteTransactionTagParams) error
	DeleteAllTransactionTags(ctx context.Context, transactionID int64) error

	// Recurring operations
	CreateRecurring(ctx context.Context, arg CreateRecurringParams) (Recurring, error)
	GetRecurringByID(ctx context.Context, id int64) (Recurring, error)
	ListRecurring(ctx context.Context, userID int64) ([]Recurring, error)
	ListActiveRecurring(ctx context.Context, userID int64) ([]Recurring, error)
	GetRecurringByTag(ctx context.Context, tagID int64) ([]Recurring, error)
	GetRecurringDueOnDate(ctx context.Context, nextDueDate time.Time) ([]Recurring, error)
	UpdateRecurring(ctx context.Context, arg UpdateRecurringParams) (Recurring, error)
	UpdateRecurringNextDue(ctx context.Context, arg UpdateRecurringNextDueParams) error
	ToggleRecurringActive(ctx context.Context, id int64) error
	DeleteRecurring(ctx context.Context, id int64) error

	// Recurring tag operations
	CreateRecurringTag(ctx context.Context, arg CreateRecurringTagParams) error
	GetRecurringTags(ctx context.Context, recurringID int64) ([]Tag, error)
	DeleteRecurringTag(ctx context.Context, arg DeleteRecurringTagParams) error
	DeleteAllRecurringTags(ctx context.Context, recurringID int64) error

	// Settings operations
	CreateSetting(ctx context.Context, arg CreateSettingParams) (Setting, error)
	GetSetting(ctx context.Context, key string) (Setting, error)
	ListSettings(ctx context.Context) ([]Setting, error)
	UpdateSetting(ctx context.Context, arg UpdateSettingParams) (Setting, error)
	DeleteSetting(ctx context.Context, key string) error

	// Report operations
	GetMonthlyReport(ctx context.Context, arg GetMonthlyReportParams) ([]GetMonthlyReportRow, error)
	GetMonthlyTotals(ctx context.Context, arg GetMonthlyTotalsParams) (GetMonthlyTotalsRow, error)
} 