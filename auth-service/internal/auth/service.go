package auth

import "context"

// Service defines the authentication service interface
type Service interface {
	// Register creates a new user account with hashed password
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)

	// Login authenticates user and returns JWT tokens
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)

	// RefreshToken issues new access and refresh tokens
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)

	// ValidateToken verifies JWT token signature and expiration
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)

	// Logout invalidates the refresh token
	Logout(ctx context.Context, userID string, refreshToken string) error
}
