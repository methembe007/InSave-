package auth

import (
	"context"
	"time"
)

// User represents a user in the database
type User struct {
	ID           string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	DateOfBirth  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Repository defines database operations for auth
type Repository interface {
	// CreateUser inserts a new user into the database
	CreateUser(ctx context.Context, user *User) error

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, userID string) (*User, error)

	// EmailExists checks if an email is already registered
	EmailExists(ctx context.Context, email string) (bool, error)
}
