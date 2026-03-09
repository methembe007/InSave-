package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/insavein/auth-service/internal/auth"
)

// AuthMiddleware validates JWT tokens for protected routes
type AuthMiddleware struct {
	service auth.Service
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(service auth.Service) *AuthMiddleware {
	return &AuthMiddleware{service: service}
}

// Authenticate validates the JWT token and adds user info to context
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		token := parts[1]
		claims, err := m.service.ValidateToken(r.Context(), token)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "roles", claims.Roles)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + http.StatusText(status) + `","message":"` + message + `"}`))
}
