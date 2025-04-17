package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTransactionsRecurrences(t *testing.T) {
	testRepo := NewRepo(&testApp)

	// Mock the expected row data
	mockedRows := mockDB.NewRows([]string{
		"id", "name", "description", "add_time", "created_at", "updated_at",
	}).AddRow(
		1, "Monthly", "Monthly recurrence", "1 month", time.Now(), time.Now(),
	).AddRow(
		2, "Weekly", "Weekly recurrence", "1 week", time.Now(), time.Now(),
	)

	mockDB.ExpectQuery("SELECT id, name, description, add_time, created_at, updated_at FROM transactions_recurrences").
		WillReturnRows(mockedRows)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/transactions/recurrences", nil)

	handler := http.HandlerFunc(testRepo.TransactionsRecurrences)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestTransactionsRecurrencesCreate(t *testing.T) {
	testRepo := NewRepo(&testApp)

	mockDB.ExpectQuery("INSERT INTO transactions_recurrences").
		WithArgs(
			"New Recurrence",
			"Test Description",
			"2 weeks",
			AnyTime{},
			AnyTime{},
		).WillReturnRows(mockDB.NewRows([]string{"id"}).AddRow(1))

	rr := httptest.NewRecorder()
	reqBody := `{
		"name": "New Recurrence",
		"description": "Test Description",
		"add_time": "2 weeks"
	}`
	req, _ := http.NewRequest("POST", "/dashboard/transactions/recurrences/create", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.TransactionsRecurrencesCreate)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, but got %d", http.StatusAccepted, rr.Code)
	}
}
