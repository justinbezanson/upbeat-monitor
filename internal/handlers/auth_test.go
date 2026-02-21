package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/justin/upbeat-monitor/internal/repository"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DDL for the full test schema
const createTestSchemaSQL = `
-- Create "accounts" table
CREATE TABLE "accounts" (
  "id" text NOT NULL,
  "name" text NULL,
  "tier" text NULL DEFAULT 'free',
  "created_at" datetime NULL,
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" text NOT NULL,
  "account_id" text NULL,
  "email" text NOT NULL,
  "password_hash" text NULL,
  "role" text NULL DEFAULT 'owner',
  "created_at" datetime NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_account" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id")
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create "oauth_clients" table
CREATE TABLE "oauth_clients" (
  "client_id" text NOT NULL,
  "account_id" text NULL,
  "client_secret_hash" text NULL,
  "name" text NULL,
  "created_at" datetime NULL,
  PRIMARY KEY ("client_id"),
  CONSTRAINT "fk_oauth_clients_account" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id")
);
-- Create "monitors" table
CREATE TABLE "monitors" (
  "id" text NOT NULL,
  "account_id" text NULL,
  "friendly_name" text NULL,
  "type" text NULL,
  "url" text NULL,
  "interval_seconds" integer NULL,
  "status" text NULL,
  "last_checked_at" datetime NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_monitors_account" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id")
);
-- Create "checks" table
CREATE TABLE "checks" (
  "id" text NOT NULL,
  "monitor_id" text NULL,
  "status_code" integer NULL,
  "latency_ms" integer NULL,
  "success" bool NULL,
  "created_at" datetime NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_checks_monitor" FOREIGN KEY ("monitor_id") REFERENCES "monitors" ("id")
);
`

// setupTestDB initializes an in-memory SQLite database for testing.
// It returns a *sql.DB instance and a cleanup function.
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}

	// Apply schema
	_, err = db.Exec(createTestSchemaSQL)
	if err != nil {
		db.Close()
		t.Fatalf("failed to create schema: %v", err)
	}

	return db, func() {
		db.Close()
	}
}

// newTestHandlers creates a new Handlers instance for testing.
func newTestHandlers(t *testing.T, db *sql.DB) *Handlers {
	userRepo := repository.NewSQLiteUserRepository(db)
	return &Handlers{
		DB:           db,
		UserRepository: userRepo,
		JWTSecret:    []byte("test_secret_key_12345678901234567890123456789012"), // Dummy key for tests
	}
}

// sendRequest is a helper to create and send an HTTP request during tests.
func sendRequest(method, path string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var reqBody bytes.Buffer
	if body != nil {
		json.NewEncoder(&reqBody).Encode(body)
	}

	req, _ := http.NewRequest(method, path, &reqBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	return req, rr
}

// TestMain sets up and tears down test environment for all auth tests.
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Clean up
	os.Exit(code)
}

func TestRegisterSuccess(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	h := newTestHandlers(t, db)

	reqBody := RegisterRequest{Email: "test@example.com", Password: "password123"}
	req, rr := sendRequest("POST", "/register", reqBody)

	handler := http.HandlerFunc(h.Register)
	handler.ServeHTTP(rr, req) 

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var res AuthResponse
	err := json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if res.Message != "User registered successfully" {
		t.Errorf("unexpected message: got %q want %q", res.Message, "User registered successfully")
	}
	if res.UserID == "" {
		t.Error("expected user ID in response, got empty")
	}

	// Verify user exists in DB
	user, err := h.UserRepository.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to get user from DB: %v", err)
	}
	if user == nil {
		t.Error("user not found in DB after registration")
	}
	if user.Email != "test@example.com" {
		t.Errorf("user email mismatch in DB: got %q want %q", user.Email, "test@example.com")
	}
}

func TestRegisterInvalidInput(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	h := newTestHandlers(t, db)

	testCases := []struct {
		name         string
		reqBody      RegisterRequest
		statusCode   int
		errorMessage string
	}{
		{
			name:         "Invalid Email Format",
			reqBody:      RegisterRequest{Email: "invalid-email", Password: "password123"},
			statusCode:   http.StatusBadRequest,
			errorMessage: "Invalid email format",
		},
		{
			name:         "Password Too Short",
			reqBody:      RegisterRequest{Email: "test@example.com", Password: "short"},
			statusCode:   http.StatusBadRequest,
			errorMessage: "Password must be at least 8 characters long",
		},
		{
			name:         "Empty Email",
			reqBody:      RegisterRequest{Email: "", Password: "password123"},
			statusCode:   http.StatusBadRequest,
			errorMessage: "Invalid email format", // This is because "" matches regex failure for isValidEmail
		},
		{
			name:         "Empty Password",
			reqBody:      RegisterRequest{Email: "test@example.com", Password: ""},
			statusCode:   http.StatusBadRequest,
			errorMessage: "Password must be at least 8 characters long",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, rr := sendRequest("POST", "/register", tc.reqBody)
			handler := http.HandlerFunc(h.Register)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.statusCode)
			}

			var res map[string]string
			err := json.NewDecoder(rr.Body).Decode(&res)
			if err != nil {
				t.Fatalf("could not decode error response: %v", err)
			}
			if res["error"] != tc.errorMessage {
				t.Errorf("unexpected error message: got %q want %q", res["error"], tc.errorMessage)
			}
		})
	}

	t.Run("Malformed JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(`{"email": "malformed",`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.Register)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for malformed JSON: got %v want %v", status, http.StatusBadRequest)
		}
		var res map[string]string
		err := json.NewDecoder(rr.Body).Decode(&res)
		if err != nil {
			t.Fatalf("could not decode error response: %v", err)
		}
		if res["error"] != "Invalid request payload" {
			t.Errorf("unexpected error message: got %q want %q", res["error"], "Invalid request payload")
		}
	})
}

func TestRegisterExistingUser(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	h := newTestHandlers(t, db)

	// Register user successfully first
	reqBody1 := RegisterRequest{Email: "existing@example.com", Password: "password123"}
	req1, rr1 := sendRequest("POST", "/register", reqBody1)
	handler := http.HandlerFunc(h.Register)
	handler.ServeHTTP(rr1, req1)

	if status := rr1.Code; status != http.StatusCreated {
		t.Fatalf("initial registration failed with status: %v", status)
	}

	// Attempt to register same user again
	reqBody2 := RegisterRequest{Email: "existing@example.com", Password: "anotherpassword"}
	req2, rr2 := sendRequest("POST", "/register", reqBody2)
	handler.ServeHTTP(rr2, req2)

	if status := rr2.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code for existing user: got %v want %v", status, http.StatusConflict)
	}

	var res map[string]string
	err := json.NewDecoder(rr2.Body).Decode(&res)
	if err != nil {
		t.Fatalf("could not decode error response: %v", err)
	}
	if res["error"] != "User with this email already exists" {
		t.Errorf("unexpected error message: got %q want %q", res["error"], "User with this email already exists")
	}
}

func TestLoginSuccess(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	h := newTestHandlers(t, db)

	// Prerequisite: Register a user
	registerReq := RegisterRequest{Email: "login@example.com", Password: "password123"}
	reqRegister, rrRegister := sendRequest("POST", "/register", registerReq)
	registerHandler := http.HandlerFunc(h.Register)
	registerHandler.ServeHTTP(rrRegister, reqRegister)
	if status := rrRegister.Code; status != http.StatusCreated {
		t.Fatalf("failed to setup user for login test: got status %v", status)
	}

	// Test Login
	loginReq := LoginRequest{Email: "login@example.com", Password: "password123"}
	reqLogin, rrLogin := sendRequest("POST", "/login", loginReq)
	loginHandler := http.HandlerFunc(h.Login)
	loginHandler.ServeHTTP(rrLogin, reqLogin)

	if status := rrLogin.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res AuthResponse
	err := json.NewDecoder(rrLogin.Body).Decode(&res)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if res.Message != "Logged in successfully" {
		t.Errorf("unexpected message: got %q want %q", res.Message, "Logged in successfully")
	}
	if res.Token == "" {
		t.Error("expected JWT token in response, got empty")
	}

	// Optional: Validate JWT token
	// This requires parsing and validating the token which can be complex.
	// For a basic test, checking if it's not empty might suffice.
	// A more robust test would involve parsing the token and checking its claims.
}

func TestLoginInvalidInput(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	h := newTestHandlers(t, db)

	// Test Malformed JSON
	t.Run("Malformed JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(`{"email": "malformed",`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.Login)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for malformed JSON: got %v want %v", status, http.StatusBadRequest)
		}
		var res map[string]string
		err := json.NewDecoder(rr.Body).Decode(&res)
		if err != nil {
			t.Fatalf("could not decode error response: %v", err)
		}
		if res["error"] != "Invalid request payload" {
			t.Errorf("unexpected error message: got %q want %q", res["error"], "Invalid request payload")
		}
	})
}

func TestLoginUserNotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	h := newTestHandlers(t, db)

	loginReq := LoginRequest{Email: "nonexistent@example.com", Password: "password123"}
	reqLogin, rrLogin := sendRequest("POST", "/login", loginReq)
	loginHandler := http.HandlerFunc(h.Login)
	loginHandler.ServeHTTP(rrLogin, reqLogin)

	if status := rrLogin.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var res map[string]string
	err := json.NewDecoder(rrLogin.Body).Decode(&res)
	if err != nil {
		t.Fatalf("could not decode error response: %v", err)
	}
	if res["error"] != "Invalid credentials" {
		t.Errorf("unexpected error message: got %q want %q", res["error"], "Invalid credentials")
	}
}

func TestLoginIncorrectPassword(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	h := newTestHandlers(t, db)

	// Prerequisite: Register a user
	registerReq := RegisterRequest{Email: "user@example.com", Password: "correctpassword"}
	reqRegister, rrRegister := sendRequest("POST", "/register", registerReq)
	registerHandler := http.HandlerFunc(h.Register)
	registerHandler.ServeHTTP(rrRegister, reqRegister)
	if status := rrRegister.Code; status != http.StatusCreated {
		t.Fatalf("failed to setup user for login test: got status %v", status)
	}

	// Test Login with incorrect password
	loginReq := LoginRequest{Email: "user@example.com", Password: "wrongpassword"}
	reqLogin, rrLogin := sendRequest("POST", "/login", loginReq)
	loginHandler := http.HandlerFunc(h.Login)
	loginHandler.ServeHTTP(rrLogin, reqLogin)

	if status := rrLogin.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var res map[string]string
	err := json.NewDecoder(rrLogin.Body).Decode(&res)
	if err != nil {
		t.Fatalf("could not decode error response: %v", err)
	}
	if res["error"] != "Invalid credentials" {
		t.Errorf("unexpected error message: got %q want %q", res["error"], "Invalid credentials")
	}
}
