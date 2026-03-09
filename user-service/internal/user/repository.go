package user

import (
	"context"
	"time"
)

// User represents a user record in the database
type User struct {
	ID              string
	Email           string
	FirstName       string
	LastName        string
	DateOfBirth     time.Time
	ProfileImageURL string
	Preferences     map[string]interface{}
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Repository defines database operations for user profile
type Repository interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, userID string) (*User, error)

	// UpdateUser updates user profile fields
	UpdateUser(ctx context.Context, user *User) error

	// DeleteUser deletes a user and all associated data (cascade)
	DeleteUser(ctx context.Context, userID string) error
}
