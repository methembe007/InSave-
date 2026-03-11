package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Missing authorization header"}`, http.StatusUnauthorized)
			return
		}
		
		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"Invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}
		
		tokenString := parts[1]
		
		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			
			// Get secret from environment
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return nil, fmt.Errorf("JWT_SECRET not configured")
			}
			
			return []byte(secret), nil
		})
		
		if err != nil {
			http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}
		
		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, `{"error":"Invalid token claims"}`, http.StatusUnauthorized)
			return
		}
		
		// Validate expiration
		exp, ok := claims["exp"].(float64)
		if !ok {
			http.Error(w, `{"error":"Missing expiration claim"}`, http.StatusUnauthorized)
			return
		}
		
		if time.Now().Unix() > int64(exp) {
			http.Error(w, `{"error":"Token expired"}`, http.StatusUnauthorized)
			return
		}
		
		// Extract user ID
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, `{"error":"Missing or invalid user_id"}`, http.StatusUnauthorized)
			return
		}
		
		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		
		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the user ID from the request context
func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
