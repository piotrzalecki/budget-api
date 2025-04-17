package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTransactionsAll(t *testing.T) {
	testRepo := NewRepo(&testApp)

	// Mock the expected row data
	mockedRows := mockDB.NewRows([]string{
		"id", "name", "description", "quote", "active", "starts", "ends",
		"created_at", "updated_at", "budget", "budget_name",
		"transaction_recurrence", "recurrence_name", "tag", "tag_name",
	}).AddRow(
		1, "Test Transaction", "Description", 100.50, true, time.Now(), time.Now(),
		time.Now(), time.Now(), 1, "Budget Name",
		1, "Monthly", 1, "Tag Name",
	)

	mockDB.ExpectQuery("SELECT tr.id, tr.name, tr.description, tr.quote").
		WillReturnRows(mockedRows)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/transactions", nil)

	handler := http.HandlerFunc(testRepo.TransactionsAll)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestTransactionsCreateUpdate(t *testing.T) {
	testRepo := NewRepo(&testApp)

	// Test Create
	mockDB.ExpectQuery("INSERT INTO transactions").
		WithArgs(
			"New Transaction", "Description", 1, 100.50, 1, true,
			AnyTime{}, AnyTime{}, 1, AnyTime{}, AnyTime{},
		).WillReturnRows(mockDB.NewRows([]string{"id"}).AddRow(1))

	rr := httptest.NewRecorder()
	reqBody := `{
		"name": "New Transaction",
		"description": "Description",
		"budget": {"id": 1},
		"quote": 100.50,
		"transaction_recurrence": {"id": 1},
		"active": true,
		"tag": {"id": 1}
	}`
	req, _ := http.NewRequest("POST", "/dashboard/transactions/create", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.TransactionsCreateUpdate)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, but got %d", http.StatusAccepted, rr.Code)
	}

	// Test Update
	mockDB.ExpectExec("UPDATE transactions SET").
		WithArgs(
			"Updated Transaction", "Updated Description", 1, 200.50, 1,
			true, AnyTime{}, AnyTime{}, 1, AnyTime{}, 1,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	rr = httptest.NewRecorder()
	reqBody = `{
		"id": 1,
		"name": "Updated Transaction",
		"description": "Updated Description",
		"budget": {"id": 1},
		"quote": 200.50,
		"transaction_recurrence": {"id": 1},
		"active": true,
		"tag": {"id": 1}
	}`
	req, _ = http.NewRequest("POST", "/dashboard/transactions/update", strings.NewReader(reqBody))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, but got %d", http.StatusAccepted, rr.Code)
	}
}

func TestTransactionsDelete(t *testing.T) {
	testRepo := NewRepo(&testApp)

	mockDB.ExpectExec("DELETE FROM transactions WHERE").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	reqBody := `{"id": 1}`
	req, _ := http.NewRequest("POST", "/dashboard/transactions/delete", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.TransactionsDelete)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestTransactionsSetStatus(t *testing.T) {
	testRepo := NewRepo(&testApp)

	mockDB.ExpectExec("UPDATE transactions SET active").
		WithArgs(true, AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	reqBody := `{"id": 1, "status": true}`
	req, _ := http.NewRequest("POST", "/dashboard/transactions/status", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.TransactionsSetStatus)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestTransactionsSetStatusAllActive(t *testing.T) {
	testRepo := NewRepo(&testApp)

	mockDB.ExpectExec("UPDATE transactions SET active").
		WithArgs(true, AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/dashboard/transactions/status/all", nil)

	handler := http.HandlerFunc(testRepo.TransactionsSetStatusAllActive)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}
