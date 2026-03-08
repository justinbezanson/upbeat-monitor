package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"time"

	"github.com/justin/upbeat-monitor/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Claims defines the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the response body for successful authentication
type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}

// ForgotPasswordRequest represents the request body for password reset request
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest represents the request body for resetting password
type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// respondWithJSON sends a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling JSON response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError sends an error JSON response.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// isValidEmail checks if the email format is valid.
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

// Register handles user registration.
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if !isValidEmail(req.Email) {
		respondWithError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	if len(req.Password) < 8 {
		respondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	// Check if user already exists
	existingUser, err := h.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Error checking for existing user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if existingUser != nil {
		respondWithError(w, http.StatusConflict, "User with this email already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Generate UUIDs for Account and User IDs
	accountID := uuid.New().String()
	userID := uuid.New().String()

	// Create new Account
	newAccount := &repository.Account{
		ID:        accountID,
		Name:      sql.NullString{String: "Default Account", Valid: true}, // Default name for now
		Tier:      "free",
		CreatedAt: time.Now(),
	}

	if err := h.UserRepository.CreateAccount(newAccount); err != nil {
		log.Printf("Error creating account: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create new User linked to the new Account
	newUser := &repository.User{
		ID:           userID,
		AccountID:    accountID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "owner", // Default role
		CreatedAt:    time.Now(),
	}

	if err := h.UserRepository.CreateUser(newUser); err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Generate JWT token for immediate login after registration
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: newUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JWTSecret)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		// We created the user but failed to log them in. 
		// They can still log in manually.
		respondWithJSON(w, http.StatusCreated, AuthResponse{Message: "User registered successfully", UserID: newUser.ID})
		return
	}

	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   true, // Always true for production/SSL
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	respondWithJSON(w, http.StatusCreated, AuthResponse{Message: "User registered successfully", UserID: newUser.ID})
}

// Login handles user login.
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if user == nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JWTSecret)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not generate authentication token")
		return
	}

	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   true, // Always true for production/SSL
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	respondWithJSON(w, http.StatusOK, AuthResponse{Message: "Logged in successfully", UserID: user.ID})
}

// Logout handles user logout by clearing the cookie.
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	respondWithJSON(w, http.StatusOK, AuthResponse{Message: "Logged out successfully"})
}

// ForgotPassword handles the request to send a password reset email.
func (h *Handlers) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// We don't want to leak if a user exists or not, so we always return OK.
	if user == nil {
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "If an account with that email exists, a password reset link has been sent."})
		return
	}

	// Generate a random token
	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// Create reset token in DB
	resetToken := &repository.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token expires in 1 hour
		CreatedAt: time.Now(),
	}

	// Delete any existing tokens for this user first
	if err := h.UserRepository.DeletePasswordResetTokensByUserID(user.ID); err != nil {
		log.Printf("Error deleting existing reset tokens: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if err := h.UserRepository.CreatePasswordResetToken(resetToken); err != nil {
		log.Printf("Error creating reset token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Send email
	if err := h.sendResetEmail(user.Email, token); err != nil {
		log.Printf("Error sending reset email: %v", err)
		// We still return OK to the user to avoid leaking existence, but this is a failure.
		// In a real system, you might want to handle this more robustly.
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "If an account with that email exists, a password reset link has been sent."})
}

// ResetPassword handles the password reset using the token.
func (h *Handlers) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Password) < 8 {
		respondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	hash := sha256.Sum256([]byte(req.Token))
	tokenHash := hex.EncodeToString(hash[:])

	token, err := h.UserRepository.GetPasswordResetToken(tokenHash)
	if err != nil {
		log.Printf("Error retrieving reset token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if token == nil || token.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusBadRequest, "Invalid or expired reset token")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Update user password
	if err := h.UserRepository.UpdateUserPassword(token.UserID, string(hashedPassword)); err != nil {
		log.Printf("Error updating password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Delete the used token
	if err := h.UserRepository.DeletePasswordResetToken(token.ID); err != nil {
		log.Printf("Error deleting used reset token: %v", err)
		// Not fatal, but good to clean up
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Password has been successfully reset."})
}

func (h *Handlers) sendResetEmail(to, token string) error {
	auth := smtp.PlainAuth("", h.SMTPUsername, h.SMTPPassword, h.SMTPAddress)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost"
	}
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	from := h.SMTPFrom
	if from == "" {
		from = "noreply@upbeat-monitor.com"
	}

	msg := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: Password Reset Request\r\n" +
		"\r\n" +
		"Please use the following link to reset your password:\r\n" +
		resetURL + "\r\n")

	addr := fmt.Sprintf("%s:%s", h.SMTPAddress, h.SMTPPort)
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}
