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