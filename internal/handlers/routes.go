package handlers

import (
	"net/http"
	"os"
)

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
func RegisterRoutes() http.Handler { // Changed return type to http.Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", Ping)

	// Wrap the mux with the CORS middleware
	return corsMiddleware(mux)
}
