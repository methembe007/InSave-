package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens for protected routes
type AuthMiddleware struct {
	jwtSecret []byte
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(jwtSecret),
	}
}

// TokenClaims represents JWT token payload
type TokenClaims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	ExpiresAt int64    `json:"exp"`
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

		tokenString := parts[1]
		claims, err := m.validateToken(tokenString)
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

// validateToken verifies JWT token signature and expiration
func (m *AuthMiddleware) validateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Validate expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("missing expiration claim")
	}

	if time.Now().Unix() > int64(exp) {
		return nil, errors.New("token expired")
	}

	// Extract required fields
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return nil, errors.New("missing or invalid user_id")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("missing or invalid email")
	}

	// Extract roles (optional)
	var roles []string
	if rolesInterface, ok := claims["roles"].([]interface{}); ok {
		for _, role := range rolesInterface {
			if roleStr, ok := role.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	return &TokenClaims{
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		ExpiresAt: int64(exp),
	}, nil
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + http.StatusText(status) + `","message":"` + message + `"}`))
}
