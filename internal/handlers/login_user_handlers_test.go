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
	"golang.org/x/crypto/bcrypt"
)

func setupTest() {
	// Clear any existing expectations
	mockDB.ExpectationsWereMet()
}

func TestLogin(t *testing.T) {
	setupTest()
	testRepo := NewRepo(&testApp)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 12)

	// Fix: Match the exact query from user_model.go GetByEmail
	mockedUserRow := mockDB.NewRows([]string{
		"id", "email", "first_name", "last_name", "password", "user_active", "created_at", "updated_at",
	}).AddRow(
		1, "test@example.com", "Test", "User", hashedPassword, 1, time.Now(), time.Now(),
	)

	// Fix: Match exact query from GetByEmail
	mockDB.ExpectQuery(`select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = \$1`).
		WithArgs("test@example.com").
		WillReturnRows(mockedUserRow)

	// Fix: Match token insertion query from token_model.go Insert method
	mockDB.ExpectExec(`delete from tokens where user_id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mockDB.ExpectExec(`insert into tokens \(user_id, email, token, token_hash, created_at, updated_at, expiry\)`).
		WithArgs(1, "test@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	reqBody := `{"email": "test@example.com", "password": "password"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.Login)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestLogout(t *testing.T) {
	setupTest()
	testRepo := NewRepo(&testApp)

	mockDB.ExpectExec("delete from tokens").
		WithArgs("test-token").
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	reqBody := `{"token": "test-token"}`
	req, _ := http.NewRequest("POST", "/logout", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.Logout)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestAllUsers(t *testing.T) {
	setupTest()
	testRepo := NewRepo(&testApp)

	// Fix: Match the exact query and columns from user_model.go GetAll
	mockedRows := mockDB.NewRows([]string{
		"id", "email", "first_name", "last_name", "password", "user_active", "created_at", "updated_at", "has_token",
	}).AddRow(
		1, "test@example.com", "Test", "User", "password_hash", 1, time.Now(), time.Now(), 0,
	)

	// Fix: Match exact query from GetAll
	mockDB.ExpectQuery(`select id, email, first_name, last_name, password, user_active, created_at, updated_at,
		case
			when \(select count\(id\) from tokens t where user_id = users\.id and t\.expiry > NOW\(\)\) > 0 then 1
			else 0
		end as has_token
		from users order by last_name`).
		WillReturnRows(mockedRows)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)

	handler := http.HandlerFunc(testRepo.AllUsers)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestEditUser(t *testing.T) {
	setupTest()
	testRepo := NewRepo(&testApp)

	mockDB.ExpectQuery("insert into users").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnRows(mockDB.NewRows([]string{"id"}).AddRow(1))

	rr := httptest.NewRecorder()
	reqBody := `{
		"email": "new@example.com",
		"first_name": "New",
		"last_name": "User",
		"active": 1,
		"password": "newpassword"
	}`
	req, _ := http.NewRequest("POST", "/users/edit", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.EditUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, but got %d", http.StatusAccepted, rr.Code)
	}
}

func TestGetUser(t *testing.T) {
	setupTest()
	testRepo := NewRepo(&testApp)

	mockedRow := mockDB.NewRows([]string{
		"id", "email", "first_name", "last_name", "password", "user_active", "created_at", "updated_at",
	}).AddRow(
		1, "test@example.com", "Test", "User", "password_hash", 1, time.Now(), time.Now(),
	)

	mockDB.ExpectQuery("select (.+) from users where id = \\$1").
		WithArgs(1).
		WillReturnRows(mockedRow)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)

	// Set up the chi context with URL parameters
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler := http.HandlerFunc(testRepo.GetUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	setupTest()
	testRepo := NewRepo(&testApp)

	mockDB.ExpectExec("delete from users where").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	reqBody := `{"id": 1}`
	req, _ := http.NewRequest("POST", "/users/delete", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.DeleteUser)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestValidateToken(t *testing.T) {
	setupTest()
	testRepo := NewRepo(&testApp)

	// Fix the query to match the actual query in token_model.go ValidToken method
	mockDB.ExpectQuery(`select id, user_id, email, token, token_hash, created_at, updated_at, expiry 
		from tokens where token = \$1`).
		WithArgs("test-token").
		WillReturnRows(mockDB.NewRows([]string{
			"id", "user_id", "email", "token", "token_hash", "created_at", "updated_at", "expiry",
		}).AddRow(
			1, 1, "test@example.com", "test-token", []byte("hash"), time.Now(), time.Now(), time.Now().Add(time.Hour),
		))

	rr := httptest.NewRecorder()
	reqBody := `{"token": "test-token"}`
	req, _ := http.NewRequest("POST", "/validate-token", strings.NewReader(reqBody))

	handler := http.HandlerFunc(testRepo.ValidateToken)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}
}