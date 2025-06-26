package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piotrzalecki/budget-api/internal/repo"
	"github.com/piotrzalecki/budget-api/pkg/model"
	"github.com/stretchr/testify/assert"
)

// mockTransactionRepo implements repo.Repository with transaction methods for tests
type mockTransactionRepo struct {
	transactions []repo.Transaction
	tags         []repo.Tag
	transactionTags map[int64][]repo.Tag // transactionID -> tags
}

func (m *mockTransactionRepo) GetDB() *sql.DB {
	return nil
}

func (m *mockTransactionRepo) WithTx(ctx context.Context, fn func(repo.Repository) error) error {
	return fn(m)
}

func (m *mockTransactionRepo) CreateTransaction(ctx context.Context, arg repo.CreateTransactionParams) (repo.Transaction, error) {
	transaction := repo.Transaction{
		ID:              int64(len(m.transactions) + 1),
		UserID:          arg.UserID,
		AmountPence:     arg.AmountPence,
		TDate:           arg.TDate,
		Note:            arg.Note,
		CreatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		SourceRecurring: arg.SourceRecurring,
		DeletedAt:       sql.NullTime{Valid: false},
	}
	m.transactions = append(m.transactions, transaction)
	return transaction, nil
}

func (m *mockTransactionRepo) GetTransactionByID(ctx context.Context, id int64) (repo.Transaction, error) {
	for _, t := range m.transactions {
		if t.ID == id && !t.DeletedAt.Valid {
			return t, nil
		}
	}
	return repo.Transaction{}, sql.ErrNoRows
}

func (m *mockTransactionRepo) ListTransactions(ctx context.Context, arg repo.ListTransactionsParams) ([]repo.Transaction, error) {
	var result []repo.Transaction
	for _, t := range m.transactions {
		if t.UserID == arg.UserID && !t.DeletedAt.Valid {
			if t.TDate.After(arg.TDate) || t.TDate.Equal(arg.TDate) {
				if t.TDate.Before(arg.TDate_2) || t.TDate.Equal(arg.TDate_2) {
					result = append(result, t)
				}
			}
		}
	}
	return result, nil
}

func (m *mockTransactionRepo) ListTransactionsByDateRange(ctx context.Context, userID int64) ([]repo.Transaction, error) {
	var result []repo.Transaction
	for _, t := range m.transactions {
		if t.UserID == userID && !t.DeletedAt.Valid {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *mockTransactionRepo) UpdateTransaction(ctx context.Context, arg repo.UpdateTransactionParams) (repo.Transaction, error) {
	for i, t := range m.transactions {
		if t.ID == arg.ID && !t.DeletedAt.Valid {
			m.transactions[i].AmountPence = arg.AmountPence
			m.transactions[i].TDate = arg.TDate
			m.transactions[i].Note = arg.Note
			return m.transactions[i], nil
		}
	}
	return repo.Transaction{}, sql.ErrNoRows
}

func (m *mockTransactionRepo) SoftDeleteTransaction(ctx context.Context, id int64) error {
	for i, t := range m.transactions {
		if t.ID == id && !t.DeletedAt.Valid {
			m.transactions[i].DeletedAt = sql.NullTime{Time: time.Now(), Valid: true}
			return nil
		}
	}
	return sql.ErrNoRows
}

func (m *mockTransactionRepo) GetTagByID(ctx context.Context, id int64) (repo.Tag, error) {
	for _, tag := range m.tags {
		if tag.ID == id {
			return tag, nil
		}
	}
	return repo.Tag{}, sql.ErrNoRows
}

func (m *mockTransactionRepo) CreateTransactionTag(ctx context.Context, arg repo.CreateTransactionTagParams) error {
	// Verify tag exists
	found := false
	for _, tag := range m.tags {
		if tag.ID == arg.TagID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("tag not found")
	}

	// Add to transaction tags
	m.transactionTags[arg.TransactionID] = append(m.transactionTags[arg.TransactionID], repo.Tag{ID: arg.TagID})
	return nil
}

func (m *mockTransactionRepo) GetTransactionTags(ctx context.Context, transactionID int64) ([]repo.Tag, error) {
	tags, exists := m.transactionTags[transactionID]
	if !exists {
		return []repo.Tag{}, nil
	}
	return tags, nil
}

func (m *mockTransactionRepo) DeleteAllTransactionTags(ctx context.Context, transactionID int64) error {
	delete(m.transactionTags, transactionID)
	return nil
}

// All other methods panic if called
func (m *mockTransactionRepo) CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) { panic("not implemented") }
func (m *mockTransactionRepo) GetUserByEmail(ctx context.Context, email string) (repo.User, error) { panic("not implemented") }
func (m *mockTransactionRepo) GetUserByID(ctx context.Context, id int64) (repo.User, error) { panic("not implemented") }
func (m *mockTransactionRepo) ListUsers(ctx context.Context) ([]repo.User, error) { panic("not implemented") }
func (m *mockTransactionRepo) UpdateUser(ctx context.Context, arg repo.UpdateUserParams) (repo.User, error) { panic("not implemented") }
func (m *mockTransactionRepo) DeleteUser(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockTransactionRepo) GetTransactionsByRecurringID(ctx context.Context, sourceRecurring sql.NullInt64) ([]repo.Transaction, error) { panic("not implemented") }
func (m *mockTransactionRepo) GetTransactionsByTag(ctx context.Context, tagID int64) ([]repo.Transaction, error) { panic("not implemented") }
func (m *mockTransactionRepo) HardDeleteTransaction(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockTransactionRepo) PurgeSoftDeletedTransactions(ctx context.Context, deletedAt sql.NullTime) error { panic("not implemented") }
func (m *mockTransactionRepo) CreateTag(ctx context.Context, name string) (repo.Tag, error) { panic("not implemented") }
func (m *mockTransactionRepo) GetTagByName(ctx context.Context, name string) (repo.Tag, error) { panic("not implemented") }
func (m *mockTransactionRepo) ListTags(ctx context.Context) ([]repo.Tag, error) { panic("not implemented") }
func (m *mockTransactionRepo) UpdateTag(ctx context.Context, arg repo.UpdateTagParams) (repo.Tag, error) { panic("not implemented") }
func (m *mockTransactionRepo) DeleteTag(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockTransactionRepo) DeleteTransactionTag(ctx context.Context, arg repo.DeleteTransactionTagParams) error { panic("not implemented") }
func (m *mockTransactionRepo) CreateRecurring(ctx context.Context, arg repo.CreateRecurringParams) (repo.Recurring, error) { panic("not implemented") }
func (m *mockTransactionRepo) GetRecurringByID(ctx context.Context, id int64) (repo.Recurring, error) { panic("not implemented") }
func (m *mockTransactionRepo) ListRecurring(ctx context.Context, userID int64) ([]repo.Recurring, error) { panic("not implemented") }
func (m *mockTransactionRepo) ListActiveRecurring(ctx context.Context, userID int64) ([]repo.Recurring, error) { panic("not implemented") }
func (m *mockTransactionRepo) GetRecurringByTag(ctx context.Context, tagID int64) ([]repo.Recurring, error) { panic("not implemented") }
func (m *mockTransactionRepo) GetRecurringDueOnDate(ctx context.Context, nextDueDate time.Time) ([]repo.Recurring, error) { panic("not implemented") }
func (m *mockTransactionRepo) UpdateRecurring(ctx context.Context, arg repo.UpdateRecurringParams) (repo.Recurring, error) { panic("not implemented") }
func (m *mockTransactionRepo) UpdateRecurringNextDue(ctx context.Context, arg repo.UpdateRecurringNextDueParams) error { panic("not implemented") }
func (m *mockTransactionRepo) ToggleRecurringActive(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockTransactionRepo) DeleteRecurring(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockTransactionRepo) CreateRecurringTag(ctx context.Context, arg repo.CreateRecurringTagParams) error { panic("not implemented") }
func (m *mockTransactionRepo) GetRecurringTags(ctx context.Context, recurringID int64) ([]repo.Tag, error) { panic("not implemented") }
func (m *mockTransactionRepo) DeleteRecurringTag(ctx context.Context, arg repo.DeleteRecurringTagParams) error { panic("not implemented") }
func (m *mockTransactionRepo) DeleteAllRecurringTags(ctx context.Context, recurringID int64) error { panic("not implemented") }
func (m *mockTransactionRepo) CreateSetting(ctx context.Context, arg repo.CreateSettingParams) (repo.Setting, error) { panic("not implemented") }
func (m *mockTransactionRepo) GetSetting(ctx context.Context, key string) (repo.Setting, error) { panic("not implemented") }
func (m *mockTransactionRepo) ListSettings(ctx context.Context) ([]repo.Setting, error) { panic("not implemented") }
func (m *mockTransactionRepo) UpdateSetting(ctx context.Context, arg repo.UpdateSettingParams) (repo.Setting, error) { panic("not implemented") }
func (m *mockTransactionRepo) DeleteSetting(ctx context.Context, key string) error { panic("not implemented") }
func (m *mockTransactionRepo) GetMonthlyReport(ctx context.Context, arg repo.GetMonthlyReportParams) ([]repo.GetMonthlyReportRow, error) { panic("not implemented") }
func (m *mockTransactionRepo) GetMonthlyTotals(ctx context.Context, arg repo.GetMonthlyTotalsParams) (repo.GetMonthlyTotalsRow, error) { panic("not implemented") }

func TestCreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockTransactionRepo{
		tags: []repo.Tag{
			{ID: 1, Name: "groceries"},
			{ID: 2, Name: "entertainment"},
		},
		transactionTags: make(map[int64][]repo.Tag),
	}
	h := NewHandler(mock)
	router := gin.New()
	router.POST("/transactions", ValidateRequest[model.CreateTransactionRequest](), h.CreateTransaction)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid transaction creation",
			requestBody: map[string]interface{}{
				"amount":  "-12.34",
				"t_date":  "2025-06-17",
				"note":    "Test transaction",
				"tag_ids": []int64{1, 2},
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "valid transaction without tags",
			requestBody: map[string]interface{}{
				"amount": "123.45",
				"t_date": "2025-06-17",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "invalid amount format",
			requestBody: map[string]interface{}{
				"amount": "invalid",
				"t_date": "2025-06-17",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid date format",
			requestBody: map[string]interface{}{
				"amount": "12.34",
				"t_date": "invalid-date",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid tag ID",
			requestBody: map[string]interface{}{
				"amount":  "12.34",
				"t_date":  "2025-06-17",
				"tag_ids": []int64{999}, // Non-existent tag
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if tt.expectedError {
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.NotNil(t, response["error"])
			} else {
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
				assert.Nil(t, response["error"])
				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.Contains(t, data, "id")
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockTransactionRepo{
		transactions: []repo.Transaction{
			{
				ID:          1,
				UserID:      1,
				AmountPence: -1234, // -12.34
				TDate:       time.Date(2025, 6, 17, 0, 0, 0, 0, time.UTC),
				Note:        sql.NullString{String: "Test transaction", Valid: true},
				CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
				DeletedAt:   sql.NullTime{Valid: false},
			},
		},
		transactionTags: make(map[int64][]repo.Tag),
	}
	h := NewHandler(mock)
	router := gin.New()
	router.GET("/transactions", h.GetTransactions)

	req := httptest.NewRequest("GET", "/transactions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")
	assert.Nil(t, response["error"])
	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Greater(t, len(data), 0)
	firstTransaction, ok := data[0].(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, firstTransaction, "id")
	assert.Contains(t, firstTransaction, "amount")
	assert.Contains(t, firstTransaction, "t_date")
}

func TestUpdateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		transactionID  string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name:          "update note",
			transactionID: "1",
			requestBody: map[string]interface{}{
				"note": "Updated note",
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:          "soft delete transaction",
			transactionID: "1",
			requestBody: map[string]interface{}{
				"deleted": true,
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:          "update tags",
			transactionID: "1",
			requestBody: map[string]interface{}{
				"tag_ids": []int64{1},
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:          "transaction not found",
			transactionID: "999",
			requestBody: map[string]interface{}{
				"note": "Updated note",
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:          "invalid transaction ID",
			transactionID: "invalid",
			requestBody: map[string]interface{}{
				"note": "Updated note",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh mock for each test to avoid state pollution
			mock := &mockTransactionRepo{
				transactions: []repo.Transaction{
					{
						ID:          1,
						UserID:      1,
						AmountPence: -1234,
						TDate:       time.Date(2025, 6, 17, 0, 0, 0, 0, time.UTC),
						Note:        sql.NullString{String: "Original note", Valid: true},
						CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
						DeletedAt:   sql.NullTime{Valid: false},
					},
				},
				tags: []repo.Tag{
					{ID: 1, Name: "groceries"},
				},
				transactionTags: make(map[int64][]repo.Tag),
			}
			h := NewHandler(mock)
			router := gin.New()
			router.PATCH("/transactions/:id", ValidateRequest[model.UpdateTransactionRequest](), h.UpdateTransaction)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PATCH", "/transactions/"+tt.transactionID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.NotNil(t, response["error"])
			}
		})
	}
} 