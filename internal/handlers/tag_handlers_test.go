package handlers

import (
	"context"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestTags(t *testing.T) {

	testRepo := NewRepo(&testApp)

	var mockedRow = mockDB.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).AddRow(1, "Test Tag", "Test Description", time.Now(), time.Now())

	mockDB.ExpectQuery("SELECT id, name, description, created_at, updated_at FROM tags ORDER BY name ASC").WillReturnRows(mockedRow)

	rr:= httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/dashboard/tags", nil)

	handler := http.HandlerFunc(testRepo.Tags)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}

}

func TestTagsCreateUpdate(t *testing.T) {
	testRepo := NewRepo(&testApp)

	// Test Create (id = 0)
	mockDB.ExpectQuery("INSERT INTO tags").
		WithArgs("New Tag", "New Description", AnyTime{}, AnyTime{}).
		WillReturnRows(mockDB.NewRows([]string{"id"}).AddRow(1))

	rr := httptest.NewRecorder()
	reqBody := `{"name": "New Tag", "description": "New Description"}`
	req, _ := http.NewRequest("POST", "/dashboard/tags/create", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.TagsCreateUpdate)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, but got %d", http.StatusAccepted, rr.Code)
	}

	// Test Update (id > 0)
	mockDB.ExpectExec("UPDATE tags SET").
		WithArgs("Updated Tag", "Updated Description", AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr = httptest.NewRecorder()
	reqBody = `{"id": 1, "name": "Updated Tag", "description": "Updated Description"}`
	req, _ = http.NewRequest("POST", "/dashboard/tags/update", strings.NewReader(reqBody))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, but got %d", http.StatusAccepted, rr.Code)
	}
}

func TestTagsDelete(t *testing.T) {
	testRepo := NewRepo(&testApp)

	mockDB.ExpectExec("DELETE FROM tags WHERE").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	reqBody := `{"id": 1}`
	req, _ := http.NewRequest("POST", "/dashboard/tags/delete", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.TagsDelete)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestTagById(t *testing.T) {
	testRepo := NewRepo(&testApp)

	mockedRow := mockDB.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
		AddRow(1, "Test Tag", "Test Description", time.Now(), time.Now())

	mockDB.ExpectQuery("SELECT (.+) FROM tags WHERE").
		WithArgs(1).
		WillReturnRows(mockedRow)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/tags/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
	chiCtx := chi.RouteContext(req.Context())
	chiCtx.URLParams.Add("id", "1")

	handler := http.HandlerFunc(testRepo.TagById)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

// Helper struct for matching time.Time arguments
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}