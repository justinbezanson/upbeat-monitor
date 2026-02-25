package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/justin/upbeat-monitor/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// DDL for the full test schema, including password_reset_tokens
const createPasswordResetTestSchemaSQL = `
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
-- Create "password_reset_tokens" table
CREATE TABLE "password_reset_tokens" (
  "id" text NOT NULL,
  "user_id" text NOT NULL,
  "token_hash" text NOT NULL,
  "expires_at" datetime NOT NULL,
  "created_at" datetime NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_password_reset_tokens_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
-- Create index "idx_password_reset_tokens_token_hash" to table: "password_reset_tokens"
CREATE UNIQUE INDEX "idx_password_reset_tokens_token_hash" ON "password_reset_tokens" ("token_hash");
`

func TestForgotPassword(t *testing.T) {
	db, cleanup := setupTestDBWithReset(t)
	defer cleanup()
	h := newTestHandlers(t, db)

	// Register a user
	registerUser(t, h, "reset@example.com", "password123")

	t.Run("Successful Forgot Password Request", func(t *testing.T) {
		reqBody := ForgotPasswordRequest{Email: "reset@example.com"}
		req, rr := sendRequest("POST", "/forgot-password", reqBody)
		
		h.ForgotPassword(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var res map[string]string
		json.NewDecoder(rr.Body).Decode(&res)
		expectedMsg := "If an account with that email exists, a password reset link has been sent."
		if res["message"] != expectedMsg {
			t.Errorf("unexpected message: got %q want %q", res["message"], expectedMsg)
		}

		// Verify token was created in DB
		user, _ := h.UserRepository.GetUserByEmail("reset@example.com")
        var count int
        db.QueryRow("SELECT COUNT(*) FROM password_reset_tokens WHERE user_id = ?", user.ID).Scan(&count)
        if count != 1 {
            t.Errorf("expected 1 reset token in DB, got %d", count)
        }
	})

	t.Run("Non-existent Email", func(t *testing.T) {
		reqBody := ForgotPasswordRequest{Email: "nonexistent@example.com"}
		req, rr := sendRequest("POST", "/forgot-password", reqBody)
		h.ForgotPassword(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var res map[string]string
		json.NewDecoder(rr.Body).Decode(&res)
		expectedMsg := "If an account with that email exists, a password reset link has been sent."
		if res["message"] != expectedMsg {
			t.Errorf("unexpected message: got %q want %q", res["message"], expectedMsg)
		}
	})
}

func TestResetPassword(t *testing.T) {
	db, cleanup := setupTestDBWithReset(t)
	defer cleanup()
	h := newTestHandlers(t, db)

	// Register a user
	user := registerUser(t, h, "reset_target@example.com", "oldpassword")

	t.Run("Successful Password Reset", func(t *testing.T) {
		// Create a manual token
		token := "test-token"
		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])
		
		resetToken := &repository.PasswordResetToken{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		}
		h.UserRepository.CreatePasswordResetToken(resetToken)

		reqBody := ResetPasswordRequest{
			Token:    token,
			Password: "newpassword123",
		}
		req, rr := sendRequest("POST", "/reset-password", reqBody)
		h.ResetPassword(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify password was updated
		updatedUser, _ := h.UserRepository.GetUserByEmail("reset_target@example.com")
		err := bcrypt.CompareHashAndPassword([]byte(updatedUser.PasswordHash), []byte("newpassword123"))
		if err != nil {
			t.Errorf("password was not correctly updated: %v", err)
		}

		// Verify token was deleted
		dbToken, _ := h.UserRepository.GetPasswordResetToken(tokenHash)
		if dbToken != nil {
			t.Error("expected reset token to be deleted after use")
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		token := "expired-token"
		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])
		
		resetToken := &repository.PasswordResetToken{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			CreatedAt: time.Now().Add(-2 * time.Hour),
		}
		h.UserRepository.CreatePasswordResetToken(resetToken)

		reqBody := ResetPasswordRequest{
			Token:    token,
			Password: "anothernewpassword",
		}
		req, rr := sendRequest("POST", "/reset-password", reqBody)
		h.ResetPassword(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		var res map[string]string
		json.NewDecoder(rr.Body).Decode(&res)
		if res["error"] != "Invalid or expired reset token" {
			t.Errorf("unexpected error: got %q want %q", res["error"], "Invalid or expired reset token")
		}
	})
}

// Helper functions for testing
func setupTestDBWithReset(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}

	_, err = db.Exec(createPasswordResetTestSchemaSQL)
	if err != nil {
		db.Close()
		t.Fatalf("failed to create schema: %v", err)
	}

	return db, func() {
		db.Close()
	}
}

func registerUser(t *testing.T, h *Handlers, email, password string) *repository.User {
	reqBody := RegisterRequest{Email: email, Password: password}
	req, rr := sendRequest("POST", "/register", reqBody)
	h.Register(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("failed to register user for test: %v", rr.Code)
	}
	user, _ := h.UserRepository.GetUserByEmail(email)
	return user
}
