package user

import "context"

// Service defines the user profile service interface
type Service interface {
	// GetProfile retrieves user profile by user ID
	GetProfile(ctx context.Context, userID string) (*UserProfile, error)

	// UpdateProfile updates user profile information
	UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error)

	// GetPreferences retrieves user preferences
	GetPreferences(ctx context.Context, userID string) (*UserPreferences, error)

	// UpdatePreferences updates user preferences
	UpdatePreferences(ctx context.Context, userID string, prefs UserPreferences) error

	// DeleteAccount deletes user account and all associated data
	DeleteAccount(ctx context.Context, userID string) error
}
