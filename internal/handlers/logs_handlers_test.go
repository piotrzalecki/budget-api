package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogs(t *testing.T) {
	testRepo := NewRepo(&testApp)

	// Mock the expected row data
	mockedRows := mockDB.NewRows([]string{"id", "log", "created_at", "updated_at"}).
		AddRow(1, "Test Log Entry", time.Now(), time.Now()).
		AddRow(2, "Another Log Entry", time.Now(), time.Now())

	// Set up the expected query
	mockDB.ExpectQuery("SELECT id, log, created_at, updated_at FROM logs").
		WillReturnRows(mockedRows)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/logs", nil)

	handler := http.HandlerFunc(testRepo.Logs)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}

	// Test error case
	mockDB.ExpectQuery("SELECT id, log, created_at, updated_at FROM logs").
		WillReturnError(sql.ErrNoRows)

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/dashboard/logs", nil)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		t.Error("expected error status code, but got 200 OK")
	}
}
