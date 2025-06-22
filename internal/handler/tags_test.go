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

// mockRepo implements repo.Repository with only the tag methods needed for tests
// All other methods panic if called

type mockRepo struct {
	tags []repo.Tag
}

func (m *mockRepo) CreateTag(ctx context.Context, name string) (repo.Tag, error) {
	if name == "" {
		return repo.Tag{}, errors.New("name required")
	}
	if len(name) > 100 {
		return repo.Tag{}, errors.New("name too long")
	}
	for _, t := range m.tags {
		if t.Name == name {
			return repo.Tag{}, errors.New("duplicate name")
		}
	}
	tag := repo.Tag{ID: int64(len(m.tags) + 1), Name: name}
	m.tags = append(m.tags, tag)
	return tag, nil
}

func (m *mockRepo) ListTags(ctx context.Context) ([]repo.Tag, error) {
	if len(m.tags) == 0 {
		return []repo.Tag{
			{ID: 1, Name: "groceries"},
			{ID: 2, Name: "entertainment"},
			{ID: 3, Name: "transport"},
		}, nil
	}
	return m.tags, nil
}

// All other methods panic if called
func (m *mockRepo) WithTx(ctx context.Context, fn func(repo.Repository) error) error { panic("not implemented") }
func (m *mockRepo) CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) { panic("not implemented") }
func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (repo.User, error) { panic("not implemented") }
func (m *mockRepo) GetUserByID(ctx context.Context, id int64) (repo.User, error) { panic("not implemented") }
func (m *mockRepo) ListUsers(ctx context.Context) ([]repo.User, error) { panic("not implemented") }
func (m *mockRepo) UpdateUser(ctx context.Context, arg repo.UpdateUserParams) (repo.User, error) { panic("not implemented") }
func (m *mockRepo) DeleteUser(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockRepo) CreateTransaction(ctx context.Context, arg repo.CreateTransactionParams) (repo.Transaction, error) { panic("not implemented") }
func (m *mockRepo) GetTransactionByID(ctx context.Context, id int64) (repo.Transaction, error) { panic("not implemented") }
func (m *mockRepo) ListTransactions(ctx context.Context, arg repo.ListTransactionsParams) ([]repo.Transaction, error) { panic("not implemented") }
func (m *mockRepo) ListTransactionsByDateRange(ctx context.Context, userID int64) ([]repo.Transaction, error) { panic("not implemented") }
func (m *mockRepo) GetTransactionsByRecurringID(ctx context.Context, sourceRecurring sql.NullInt64) ([]repo.Transaction, error) { panic("not implemented") }
func (m *mockRepo) GetTransactionsByTag(ctx context.Context, tagID int64) ([]repo.Transaction, error) { panic("not implemented") }
func (m *mockRepo) UpdateTransaction(ctx context.Context, arg repo.UpdateTransactionParams) (repo.Transaction, error) { panic("not implemented") }
func (m *mockRepo) SoftDeleteTransaction(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockRepo) HardDeleteTransaction(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockRepo) PurgeSoftDeletedTransactions(ctx context.Context, deletedAt sql.NullTime) error { panic("not implemented") }
func (m *mockRepo) GetTagByID(ctx context.Context, id int64) (repo.Tag, error) { panic("not implemented") }
func (m *mockRepo) GetTagByName(ctx context.Context, name string) (repo.Tag, error) { panic("not implemented") }
func (m *mockRepo) UpdateTag(ctx context.Context, arg repo.UpdateTagParams) (repo.Tag, error) { panic("not implemented") }
func (m *mockRepo) DeleteTag(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockRepo) CreateTransactionTag(ctx context.Context, arg repo.CreateTransactionTagParams) error { panic("not implemented") }
func (m *mockRepo) GetTransactionTags(ctx context.Context, transactionID int64) ([]repo.Tag, error) { panic("not implemented") }
func (m *mockRepo) DeleteTransactionTag(ctx context.Context, arg repo.DeleteTransactionTagParams) error { panic("not implemented") }
func (m *mockRepo) DeleteAllTransactionTags(ctx context.Context, transactionID int64) error { panic("not implemented") }
func (m *mockRepo) CreateRecurring(ctx context.Context, arg repo.CreateRecurringParams) (repo.Recurring, error) { panic("not implemented") }
func (m *mockRepo) GetRecurringByID(ctx context.Context, id int64) (repo.Recurring, error) { panic("not implemented") }
func (m *mockRepo) ListRecurring(ctx context.Context, userID int64) ([]repo.Recurring, error) { panic("not implemented") }
func (m *mockRepo) ListActiveRecurring(ctx context.Context, userID int64) ([]repo.Recurring, error) { panic("not implemented") }
func (m *mockRepo) GetRecurringByTag(ctx context.Context, tagID int64) ([]repo.Recurring, error) { panic("not implemented") }
func (m *mockRepo) GetRecurringDueOnDate(ctx context.Context, nextDueDate time.Time) ([]repo.Recurring, error) { panic("not implemented") }
func (m *mockRepo) UpdateRecurring(ctx context.Context, arg repo.UpdateRecurringParams) (repo.Recurring, error) { panic("not implemented") }
func (m *mockRepo) UpdateRecurringNextDue(ctx context.Context, arg repo.UpdateRecurringNextDueParams) error { panic("not implemented") }
func (m *mockRepo) ToggleRecurringActive(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockRepo) DeleteRecurring(ctx context.Context, id int64) error { panic("not implemented") }
func (m *mockRepo) CreateRecurringTag(ctx context.Context, arg repo.CreateRecurringTagParams) error { panic("not implemented") }
func (m *mockRepo) GetRecurringTags(ctx context.Context, recurringID int64) ([]repo.Tag, error) { panic("not implemented") }
func (m *mockRepo) DeleteRecurringTag(ctx context.Context, arg repo.DeleteRecurringTagParams) error { panic("not implemented") }
func (m *mockRepo) DeleteAllRecurringTags(ctx context.Context, recurringID int64) error { panic("not implemented") }
func (m *mockRepo) CreateSetting(ctx context.Context, arg repo.CreateSettingParams) (repo.Setting, error) { panic("not implemented") }
func (m *mockRepo) GetSetting(ctx context.Context, key string) (repo.Setting, error) { panic("not implemented") }
func (m *mockRepo) ListSettings(ctx context.Context) ([]repo.Setting, error) { panic("not implemented") }
func (m *mockRepo) UpdateSetting(ctx context.Context, arg repo.UpdateSettingParams) (repo.Setting, error) { panic("not implemented") }
func (m *mockRepo) DeleteSetting(ctx context.Context, key string) error { panic("not implemented") }
func (m *mockRepo) GetMonthlyReport(ctx context.Context, arg repo.GetMonthlyReportParams) ([]repo.GetMonthlyReportRow, error) { panic("not implemented") }
func (m *mockRepo) GetMonthlyTotals(ctx context.Context, arg repo.GetMonthlyTotalsParams) (repo.GetMonthlyTotalsRow, error) { panic("not implemented") }

func TestCreateTag(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockRepo{}
	h := NewHandler(mock)
	router := gin.New()
	router.POST("/tags", ValidateRequest[model.CreateTagRequest](), h.CreateTag)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid tag creation",
			requestBody: map[string]interface{}{
				"name": "groceries",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "empty name",
			requestBody: map[string]interface{}{
				"name": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "missing name",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "name too long",
			requestBody: map[string]interface{}{
				"name": "this_is_a_very_long_tag_name_that_exceeds_the_maximum_length_allowed_by_the_validation_rules_and_should_fail_the_test",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/tags", bytes.NewBuffer(body))
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

func TestGetTags(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockRepo{}
	h := NewHandler(mock)
	router := gin.New()
	router.GET("/tags", h.GetTags)

	req := httptest.NewRequest("GET", "/tags", nil)
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
	firstTag, ok := data[0].(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, firstTag, "id")
	assert.Contains(t, firstTag, "name")
} 