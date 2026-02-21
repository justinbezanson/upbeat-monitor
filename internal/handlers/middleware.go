package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ContextKey represents the type for context keys to avoid collisions.
type ContextKey string

const (
	// ContextKeyUserID is the key for the user ID in the request context.
	ContextKeyUserID ContextKey = "userID"
)

// AuthMiddleware creates a middleware that verifies JWT tokens.
func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return h.JWTSecret, nil
		})

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// Add UserID to context
			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		}
	})
}
