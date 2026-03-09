package auth

import "time"

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
}

// LoginRequest represents user login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse contains authentication tokens and user info
type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	User         UserSummary `json:"user"`
}

// TokenResponse contains new tokens after refresh
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// TokenClaims represents JWT token payload
type TokenClaims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	ExpiresAt int64    `json:"exp"`
}

// UserSummary contains basic user information
type UserSummary struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}
