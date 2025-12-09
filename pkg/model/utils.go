package model

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

// CurrencyToPence converts a currency string (e.g., "12.34" or "-12.34") to pence
func CurrencyToPence(amount string) (int64, error) {
	// Remove any leading/trailing whitespace
	amount = strings.TrimSpace(amount)
	
	// Parse as float first to handle the decimal point
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}
	
	// Convert to pence (multiply by 100 and round)
	pence := int64(amountFloat * 100)
	return pence, nil
}

// PenceToCurrency converts pence to a currency string (e.g., "12.34" or "-12.34")
func PenceToCurrency(pence int64) string {
	// Convert to float for proper decimal formatting
	amount := float64(pence) / 100.0
	
	// Format with exactly 2 decimal places
	return strconv.FormatFloat(amount, 'f', 2, 64)
}

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatDate formats a time.Time to YYYY-MM-DD string
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// StringToSQLNullString converts a string pointer to sql.NullString
func StringToSQLNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// SQLNullStringToString converts sql.NullString to a string pointer
func SQLNullStringToString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// Int64ToSQLNullInt64 converts an int64 pointer to sql.NullInt64
func Int64ToSQLNullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

// SQLNullInt64ToInt64 converts sql.NullInt64 to an int64 pointer
func SQLNullInt64ToInt64(ni sql.NullInt64) *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}

// SQLNullTimeToTimePtr converts sql.NullTime to a time.Time pointer
func SQLNullTimeToTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
} 