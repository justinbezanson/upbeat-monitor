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
		var tokenString string

		// Try to get token from cookie first
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			tokenString = cookie.Value
		} else {
			// Fallback to Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			respondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return h.JWTSecret, nil
		})

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
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
