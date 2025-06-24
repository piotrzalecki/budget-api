package model

import (
	"time"
)

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	Amount  string  `json:"amount" validate:"required,currency"`
	TDate   string  `json:"t_date" validate:"required,date"`
	Note    *string `json:"note,omitempty"`
	TagIDs  []int64 `json:"tag_ids,omitempty"`
}

// UpdateTransactionRequest represents the request body for updating a transaction
type UpdateTransactionRequest struct {
	Deleted *bool   `json:"deleted,omitempty"`
	Note    *string `json:"note,omitempty"`
	TagIDs  []int64 `json:"tag_ids,omitempty"`
}

// CreateTagRequest represents the request body for creating a tag
type CreateTagRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// CreateRecurringRequest represents the request body for creating a recurring rule
type CreateRecurringRequest struct {
	Amount        string   `json:"amount" validate:"required,currency"`
	Description   string   `json:"description" validate:"required,min=1,max=255"`
	Frequency     string   `json:"frequency" validate:"required,oneof=daily weekly monthly yearly"`
	IntervalN     int      `json:"interval_n" validate:"required,min=1,max=365"`
	FirstDueDate  string   `json:"first_due_date" validate:"required,date"`
	EndDate       *string  `json:"end_date,omitempty" validate:"omitempty,date"`
	TagIDs        []int64  `json:"tag_ids,omitempty"`
}

// UpdateRecurringRequest represents the request body for updating a recurring rule
type UpdateRecurringRequest struct {
	Active        *bool    `json:"active,omitempty"`
	Amount        *string  `json:"amount,omitempty" validate:"omitempty,currency"`
	Description   *string  `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	Frequency     *string  `json:"frequency,omitempty" validate:"omitempty,oneof=daily weekly monthly yearly"`
	IntervalN     *int     `json:"interval_n,omitempty" validate:"omitempty,min=1,max=365"`
	FirstDueDate  *string  `json:"first_due_date,omitempty" validate:"omitempty,date"`
	EndDate       *string  `json:"end_date,omitempty" validate:"omitempty,date"`
	TagIDs        []int64  `json:"tag_ids,omitempty"`
}

// TransactionResponse represents a transaction in API responses
type TransactionResponse struct {
	ID             int64     `json:"id"`
	Amount         string    `json:"amount"`
	TDate          string    `json:"t_date"`
	Note           *string   `json:"note,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	SourceRecurring *int64   `json:"source_recurring,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	TagIDs         []int64   `json:"tag_ids,omitempty"`
}

// TagResponse represents a tag in API responses
type TagResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// RecurringResponse represents a recurring rule in API responses
type RecurringResponse struct {
	ID            int64     `json:"id"`
	Amount        string    `json:"amount"`
	Description   string    `json:"description"`
	Frequency     string    `json:"frequency"`
	IntervalN     int       `json:"interval_n"`
	FirstDueDate  string    `json:"first_due_date"`
	NextDueDate   string    `json:"next_due_date"`
	EndDate       *string   `json:"end_date,omitempty"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	TagIDs        []int64   `json:"tag_ids,omitempty"`
}

// MonthlyReportResponse represents the monthly report response
type MonthlyReportResponse struct {
	TotalIn  string                    `json:"total_in"`
	TotalOut string                    `json:"total_out"`
	ByTag    map[string]TagReportEntry `json:"by_tag"`
}

// TagReportEntry represents spending/income for a specific tag
type TagReportEntry struct {
	TotalIn  string `json:"total_in"`
	TotalOut string `json:"total_out"`
}

// SchedulerResponse represents the scheduler run response
type SchedulerResponse struct {
	Processed int `json:"processed"`
}

// APIResponse represents the standard API response envelope
type APIResponse struct {
	Data  interface{} `json:"data"`
	Error *string     `json:"error"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string                 `json:"error"`
	Data  map[string]interface{} `json:"data,omitempty"`
}

// PurgeTransactionsRequest represents the request body for purging soft deleted transactions
type PurgeTransactionsRequest struct {
	CutoffDate string `json:"cutoff_date" validate:"required,date"`
} 