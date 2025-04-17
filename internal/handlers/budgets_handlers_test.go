package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestBudgets(t *testing.T) {
	testRepo := NewRepo(&testApp)

	// Mock the expected row data
	mockedRow := mockDB.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(1, "Test Budget", "Test Description", time.Now(), time.Now())

	// Set up the expected query
	mockDB.ExpectQuery("SELECT id, name, description, created_at, updated_at FROM budgets").
		WillReturnRows(mockedRow)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/budgets", nil)

	handler := http.HandlerFunc(testRepo.Budgets)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestBudgetsCreateUpdate(t *testing.T) {
	testRepo := NewRepo(&testApp)
	
	// Test Create (id = 0)
	mockDB.ExpectQuery(`INSERT INTO budgets \(name, description, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
		WithArgs("New Budget", "New Description", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(mockDB.NewRows([]string{"id"}).AddRow(1))

	rr := httptest.NewRecorder()
	reqBody := `{"name": "New Budget", "description": "New Description"}`
	req, _ := http.NewRequest("POST", "/dashboard/budgets/create", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.BudgetsCreateUpdate)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, but got %d", http.StatusAccepted, rr.Code)
	}

	// Clear expectations before next test
	mockDB.ExpectationsWereMet()

	// Test Update (id > 0)
	mockDB.ExpectExec(`UPDATE budgets SET name=\$1, description=\$2, updated_at=\$3 WHERE id=\$4`).
		WithArgs("Updated Budget", "Updated Description", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr = httptest.NewRecorder()
	reqBody = `{"id": 1, "name": "Updated Budget", "description": "Updated Description"}`
	req, _ = http.NewRequest("POST", "/dashboard/budgets/update", strings.NewReader(reqBody))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, but got %d", http.StatusAccepted, rr.Code)
	}
}

func TestBudgetsDelete(t *testing.T) {
	testRepo := NewRepo(&testApp)

	mockDB.ExpectExec("DELETE FROM budgets WHERE").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	reqBody := `{"id": 1}`
	req, _ := http.NewRequest("POST", "/dashboard/budgets/delete", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.BudgetsDelete)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestBudgetById(t *testing.T) {
	testRepo := NewRepo(&testApp)

	mockedRow := mockDB.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(1, "Test Budget", "Test Description", time.Now(), time.Now())

	mockDB.ExpectQuery("SELECT id, name, description, created_at, updated_at FROM budgets WHERE id=\\$1").
		WithArgs(1).
		WillReturnRows(mockedRow)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/budgets/1", nil)

	// Set up the chi context with URL parameters
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler := http.HandlerFunc(testRepo.BudgetsById)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}
