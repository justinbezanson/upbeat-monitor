package handlers

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/justin/upbeat-monitor/internal/repository"
)

// Handlers struct holds dependencies for HTTP handlers.
type Handlers struct {
	DB           *sql.DB
	UserRepository repository.UserRepository
	// Add other dependencies here if needed, e.g., Logger
}

// corsMiddleware adds CORS headers to responses.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("FRONTEND_URL")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost" // Default to localhost if env var not set
		}
		// Allow requests from localhost (your frontend)
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RegisterRoutes registers the HTTP routes for the application
func RegisterRoutes(db *sql.DB) http.Handler {
	userRepo := repository.NewSQLiteUserRepository(db) // Initialize repository
	h := &Handlers{
		DB:           db,
		UserRepository: userRepo,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", h.Ping)
	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/login", h.Login)

	// Wrap the mux with the CORS middleware
	return corsMiddleware(mux)
}
