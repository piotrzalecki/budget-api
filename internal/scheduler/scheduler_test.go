package scheduler

import (
	"testing"
	"time"

	"github.com/piotrzalecki/budget-api/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestCalculateNextDueDate(t *testing.T) {
	tests := []struct {
		name     string
		rule     repo.Recurring
		today    time.Time
		expected time.Time
	}{
		{
			name: "daily interval 1",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Frequency:   "daily",
				IntervalN:   1,
			},
			today:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "daily interval 2",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Frequency:   "daily",
				IntervalN:   2,
			},
			today:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "weekly interval 1",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Frequency:   "weekly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "weekly interval 2",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Frequency:   "weekly",
				IntervalN:   2,
			},
			today:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly interval 1",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly interval 3",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   3,
			},
			today:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "yearly interval 1",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				Frequency:   "yearly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "yearly interval 2",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
				Frequency:   "yearly",
				IntervalN:   2,
			},
			today:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2027, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateNextDueDate(tt.rule, tt.today)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCalculateNextDueDateFebruaryEdgeCases tests the advance date logic for February 28th and 29th edge cases
func TestCalculateNextDueDateFebruaryEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		rule     repo.Recurring
		today    time.Time
		expected time.Time
	}{
		// Monthly frequency tests for February 28th
		{
			name: "monthly from Jan 28 to Feb 28 (non-leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 28, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 1, 28, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Jan 29 to Feb 28 (non-leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 29, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 1, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Jan 30 to Feb 28 (non-leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 30, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 1, 30, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Jan 31 to Feb 28 (non-leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Feb 28 to Mar 28 (non-leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 28, 0, 0, 0, 0, time.UTC),
		},

		// Monthly frequency tests for February 29th (leap year)
		{
			name: "monthly from Jan 29 to Feb 29 (leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Jan 30 to Feb 29 (leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Jan 31 to Feb 29 (leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Feb 29 to Mar 29 (leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
		},

		// Monthly frequency tests for February 28th to March (leap year)
		{
			name: "monthly from Jan 28 to Feb 28 (leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Feb 28 to Mar 28 (leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 3, 28, 0, 0, 0, 0, time.UTC),
		},

		// Yearly frequency tests for February 28th and 29th
		{
			name: "yearly from Feb 28 2024 to Feb 28 2025 (leap to non-leap)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
				Frequency:   "yearly",
				IntervalN:   1,
			},
			today:    time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "yearly from Feb 29 2024 to Feb 28 2025 (leap to non-leap)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
				Frequency:   "yearly",
				IntervalN:   1,
			},
			today:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "yearly from Feb 28 2023 to Feb 28 2024 (non-leap to leap)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
				Frequency:   "yearly",
				IntervalN:   1,
			},
			today:    time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "yearly from Feb 29 2020 to Feb 28 2021 (leap to non-leap)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
				Frequency:   "yearly",
				IntervalN:   1,
			},
			today:    time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "yearly from Feb 29 2020 to Feb 29 2024 (leap to leap)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
				Frequency:   "yearly",
				IntervalN:   4,
			},
			today:    time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},

		// Edge cases for other months with 31 days
		{
			name: "monthly from Jan 31 to Feb 28 (non-leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Mar 31 to Apr 30",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 4, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from May 31 to Jun 30",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Jul 31 to Aug 31",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Aug 31 to Sep 30",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Oct 31 to Nov 30",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly from Dec 31 to Jan 31",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   1,
			},
			today:    time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateNextDueDate(tt.rule, tt.today)
			assert.Equal(t, tt.expected, result, 
				"Expected %v but got %v for test case: %s", 
				tt.expected.Format("2006-01-02"), 
				result.Format("2006-01-02"), 
				tt.name)
		})
	}
}

// TestCalculateNextDueDateFebruaryEdgeCasesWithIntervals tests edge cases with different intervals
func TestCalculateNextDueDateFebruaryEdgeCasesWithIntervals(t *testing.T) {
	tests := []struct {
		name     string
		rule     repo.Recurring
		today    time.Time
		expected time.Time
	}{
		// Monthly with interval 2
		{
			name: "monthly interval 2 from Jan 31 to Mar 31",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   2,
			},
			today:    time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly interval 2 from Feb 28 to Apr 28 (non-leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   2,
			},
			today:    time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 4, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monthly interval 2 from Feb 29 to Apr 29 (leap year)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
				Frequency:   "monthly",
				IntervalN:   2,
			},
			today:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 4, 29, 0, 0, 0, 0, time.UTC),
		},

		// Yearly with interval 2
		{
			name: "yearly interval 2 from Feb 29 2020 to Feb 28 2022 (leap to non-leap)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
				Frequency:   "yearly",
				IntervalN:   2,
			},
			today:    time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "yearly interval 2 from Feb 28 2023 to Feb 28 2025 (non-leap to non-leap)",
			rule: repo.Recurring{
				NextDueDate: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
				Frequency:   "yearly",
				IntervalN:   2,
			},
			today:    time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateNextDueDate(tt.rule, tt.today)
			assert.Equal(t, tt.expected, result, 
				"Expected %v but got %v for test case: %s", 
				tt.expected.Format("2006-01-02"), 
				result.Format("2006-01-02"), 
				tt.name)
		})
	}
} 