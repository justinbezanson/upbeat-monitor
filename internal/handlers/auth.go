package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/justin/upbeat-monitor/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret key (should be loaded from environment variable)
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

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

	// Generate UUID for user ID
	userID := uuid.New().String()

	newUser := &repository.User{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Tier:         "free",      // Default tier
		CreatedAt:    time.Now(),
	}

	if err := h.UserRepository.CreateUser(newUser); err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

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
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not generate authentication token")
		return
	}

	respondWithJSON(w, http.StatusOK, AuthResponse{Message: "Logged in successfully", Token: tokenString, UserID: user.ID})
}
