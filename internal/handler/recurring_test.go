package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piotrzalecki/budget-api/internal/repo"
	"github.com/piotrzalecki/budget-api/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

// Implement the required methods for the mock
func (m *MockRepository) WithTx(ctx context.Context, fn func(repo.Repository) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// Add other required methods as needed for testing
// For brevity, I'll just add the ones used in the recurring handlers

func (m *MockRepository) CreateRecurring(ctx context.Context, arg repo.CreateRecurringParams) (repo.Recurring, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Recurring), args.Error(1)
}

func (m *MockRepository) GetRecurringByID(ctx context.Context, id int64) (repo.Recurring, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repo.Recurring), args.Error(1)
}

func (m *MockRepository) ListRecurring(ctx context.Context, userID int64) ([]repo.Recurring, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]repo.Recurring), args.Error(1)
}

func (m *MockRepository) UpdateRecurring(ctx context.Context, arg repo.UpdateRecurringParams) (repo.Recurring, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Recurring), args.Error(1)
}

func (m *MockRepository) DeleteRecurring(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetRecurringByTag(ctx context.Context, tagID int64) ([]repo.Recurring, error) {
	args := m.Called(ctx, tagID)
	return args.Get(0).([]repo.Recurring), args.Error(1)
}

func (m *MockRepository) GetRecurringTags(ctx context.Context, recurringID int64) ([]repo.Tag, error) {
	args := m.Called(ctx, recurringID)
	return args.Get(0).([]repo.Tag), args.Error(1)
}

func (m *MockRepository) CreateRecurringTag(ctx context.Context, arg repo.CreateRecurringTagParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockRepository) DeleteAllRecurringTags(ctx context.Context, recurringID int64) error {
	args := m.Called(ctx, recurringID)
	return args.Error(0)
}

func (m *MockRepository) GetTagByID(ctx context.Context, id int64) (repo.Tag, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repo.Tag), args.Error(1)
}

// Add stub implementations for all other required methods
func (m *MockRepository) CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(repo.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id int64) (repo.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repo.User), args.Error(1)
}

func (m *MockRepository) ListUsers(ctx context.Context) ([]repo.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repo.User), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, arg repo.UpdateUserParams) (repo.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.User), args.Error(1)
}

func (m *MockRepository) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateTransaction(ctx context.Context, arg repo.CreateTransactionParams) (repo.Transaction, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Transaction), args.Error(1)
}

func (m *MockRepository) GetTransactionByID(ctx context.Context, id int64) (repo.Transaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repo.Transaction), args.Error(1)
}

func (m *MockRepository) ListTransactions(ctx context.Context, arg repo.ListTransactionsParams) ([]repo.Transaction, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repo.Transaction), args.Error(1)
}

func (m *MockRepository) ListTransactionsByDateRange(ctx context.Context, userID int64) ([]repo.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]repo.Transaction), args.Error(1)
}

func (m *MockRepository) GetTransactionsByRecurringID(ctx context.Context, sourceRecurring sql.NullInt64) ([]repo.Transaction, error) {
	args := m.Called(ctx, sourceRecurring)
	return args.Get(0).([]repo.Transaction), args.Error(1)
}

func (m *MockRepository) GetTransactionsByTag(ctx context.Context, tagID int64) ([]repo.Transaction, error) {
	args := m.Called(ctx, tagID)
	return args.Get(0).([]repo.Transaction), args.Error(1)
}

func (m *MockRepository) UpdateTransaction(ctx context.Context, arg repo.UpdateTransactionParams) (repo.Transaction, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Transaction), args.Error(1)
}

func (m *MockRepository) SoftDeleteTransaction(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) HardDeleteTransaction(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) PurgeSoftDeletedTransactions(ctx context.Context, deletedAt sql.NullTime) error {
	args := m.Called(ctx, deletedAt)
	return args.Error(0)
}

func (m *MockRepository) CreateTag(ctx context.Context, name string) (repo.Tag, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(repo.Tag), args.Error(1)
}

func (m *MockRepository) GetTagByName(ctx context.Context, name string) (repo.Tag, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(repo.Tag), args.Error(1)
}

func (m *MockRepository) ListTags(ctx context.Context) ([]repo.Tag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repo.Tag), args.Error(1)
}

func (m *MockRepository) UpdateTag(ctx context.Context, arg repo.UpdateTagParams) (repo.Tag, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Tag), args.Error(1)
}

func (m *MockRepository) DeleteTag(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateTransactionTag(ctx context.Context, arg repo.CreateTransactionTagParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockRepository) GetTransactionTags(ctx context.Context, transactionID int64) ([]repo.Tag, error) {
	args := m.Called(ctx, transactionID)
	return args.Get(0).([]repo.Tag), args.Error(1)
}

func (m *MockRepository) DeleteTransactionTag(ctx context.Context, arg repo.DeleteTransactionTagParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockRepository) DeleteAllTransactionTags(ctx context.Context, transactionID int64) error {
	args := m.Called(ctx, transactionID)
	return args.Error(0)
}

func (m *MockRepository) ListActiveRecurring(ctx context.Context, userID int64) ([]repo.Recurring, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]repo.Recurring), args.Error(1)
}

func (m *MockRepository) GetRecurringDueOnDate(ctx context.Context, nextDueDate time.Time) ([]repo.Recurring, error) {
	args := m.Called(ctx, nextDueDate)
	return args.Get(0).([]repo.Recurring), args.Error(1)
}

func (m *MockRepository) UpdateRecurringNextDue(ctx context.Context, arg repo.UpdateRecurringNextDueParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockRepository) ToggleRecurringActive(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) DeleteRecurringTag(ctx context.Context, arg repo.DeleteRecurringTagParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockRepository) CreateSetting(ctx context.Context, arg repo.CreateSettingParams) (repo.Setting, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Setting), args.Error(1)
}

func (m *MockRepository) GetSetting(ctx context.Context, key string) (repo.Setting, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(repo.Setting), args.Error(1)
}

func (m *MockRepository) ListSettings(ctx context.Context) ([]repo.Setting, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repo.Setting), args.Error(1)
}

func (m *MockRepository) UpdateSetting(ctx context.Context, arg repo.UpdateSettingParams) (repo.Setting, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.Setting), args.Error(1)
}

func (m *MockRepository) DeleteSetting(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRepository) GetMonthlyReport(ctx context.Context, arg repo.GetMonthlyReportParams) ([]repo.GetMonthlyReportRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]repo.GetMonthlyReportRow), args.Error(1)
}

func (m *MockRepository) GetMonthlyTotals(ctx context.Context, arg repo.GetMonthlyTotalsParams) (repo.GetMonthlyTotalsRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repo.GetMonthlyTotalsRow), args.Error(1)
}

// TestCreateRecurring tests the CreateRecurring handler
func TestCreateRecurring(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock repository
	mockRepo := new(MockRepository)
	
	// Create a handler with the mock repository
	handler := NewHandler(mockRepo)

	// Create a test request
	request := model.CreateRecurringRequest{
		Amount:        "-50.00",
		Description:   "Monthly subscription",
		Frequency:     "monthly",
		IntervalN:     1,
		FirstDueDate:  "2025-07-01",
		TagIDs:        []int64{1, 2},
	}

	// Convert request to JSON
	jsonData, _ := json.Marshal(request)

	// Create a test HTTP request
	req, _ := http.NewRequest("POST", "/api/v1/recurring", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "test-key")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Manually set the validated request in context (simulating the middleware)
	c.Set("validated_request", request)

	// Set up mock expectations
	mockRepo.On("CreateRecurring", mock.Anything, mock.AnythingOfType("repo.CreateRecurringParams")).Return(
		repo.Recurring{
			ID: 1,
			// Add other required fields as needed
		}, nil)
	
	mockRepo.On("GetTagByID", mock.Anything, int64(1)).Return(repo.Tag{ID: 1, Name: "Tag1"}, nil)
	mockRepo.On("GetTagByID", mock.Anything, int64(2)).Return(repo.Tag{ID: 2, Name: "Tag2"}, nil)
	mockRepo.On("CreateRecurringTag", mock.Anything, mock.AnythingOfType("repo.CreateRecurringTagParams")).Return(nil)

	// Call the handler
	handler.CreateRecurring(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

// TestGetRecurring tests the GetRecurring handler
func TestGetRecurring(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock repository
	mockRepo := new(MockRepository)
	
	// Create a handler with the mock repository
	handler := NewHandler(mockRepo)

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/v1/recurring", nil)
	req.Header.Set("X-API-Key", "test-key")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Set up mock expectations
	mockRepo.On("ListRecurring", mock.Anything, int64(1)).Return([]repo.Recurring{}, nil)

	// Call the handler
	handler.GetRecurring(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

// TestListActiveRecurring tests the ListActiveRecurring handler
func TestListActiveRecurring(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock repository
	mockRepo := new(MockRepository)
	
	// Create a handler with the mock repository
	handler := NewHandler(mockRepo)

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/v1/recurring/active", nil)
	req.Header.Set("X-API-Key", "test-key")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Set up mock expectations
	mockRepo.On("ListActiveRecurring", mock.Anything, int64(1)).Return([]repo.Recurring{}, nil)

	// Call the handler
	handler.ListActiveRecurring(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

// TestToggleRecurringActive tests the ToggleRecurringActive handler
func TestToggleRecurringActive(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock repository
	mockRepo := new(MockRepository)
	
	// Create a handler with the mock repository
	handler := NewHandler(mockRepo)

	// Create a test HTTP request
	req, _ := http.NewRequest("PATCH", "/api/v1/recurring/1/toggle", nil)
	req.Header.Set("X-API-Key", "test-key")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Set up mock expectations
	mockRepo.On("GetRecurringByID", mock.Anything, int64(1)).Return(repo.Recurring{ID: 1}, nil)
	mockRepo.On("ToggleRecurringActive", mock.Anything, int64(1)).Return(nil)

	// Call the handler
	handler.ToggleRecurringActive(c)

	// Assert the response
	assert.Equal(t, http.StatusNoContent, w.Code)
	
	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

// TestGetRecurringDueOnDate tests the GetRecurringDueOnDate handler
func TestGetRecurringDueOnDate(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock repository
	mockRepo := new(MockRepository)
	
	// Create a handler with the mock repository
	handler := NewHandler(mockRepo)

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/v1/recurring/due?date=2025-07-01", nil)
	req.Header.Set("X-API-Key", "test-key")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Set up mock expectations
	mockRepo.On("GetRecurringDueOnDate", mock.Anything, mock.AnythingOfType("time.Time")).Return([]repo.Recurring{}, nil)

	// Call the handler
	handler.GetRecurringDueOnDate(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

// TestGetRecurringDueOnDateWithoutDate tests the GetRecurringDueOnDate handler without date parameter
func TestGetRecurringDueOnDateWithoutDate(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock repository
	mockRepo := new(MockRepository)
	
	// Create a handler with the mock repository
	handler := NewHandler(mockRepo)

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/v1/recurring/due", nil)
	req.Header.Set("X-API-Key", "test-key")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Set up mock expectations
	mockRepo.On("GetRecurringDueOnDate", mock.Anything, mock.AnythingOfType("time.Time")).Return([]repo.Recurring{}, nil)

	// Call the handler
	handler.GetRecurringDueOnDate(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify mock expectations
	mockRepo.AssertExpectations(t)
} 