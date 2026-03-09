package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// UserService implements the Service interface
type UserService struct {
	repo Repository
}

// NewUserService creates a new user service
func NewUserService(repo Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetProfile retrieves user profile by user ID
func (s *UserService) GetProfile(ctx context.Context, userID string) (*UserProfile, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &UserProfile{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		DateOfBirth:     user.DateOfBirth.Format("2006-01-02"),
		ProfileImageURL: user.ProfileImageURL,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}, nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get existing user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format, use YYYY-MM-DD: %w", err)
		}
		user.DateOfBirth = dob
	}
	if req.ProfileImageURL != "" {
		user.ProfileImageURL = req.ProfileImageURL
	}

	// Update in database
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Return updated profile
	return s.GetProfile(ctx, userID)
}

// GetPreferences retrieves user preferences
func (s *UserService) GetPreferences(ctx context.Context, userID string) (*UserPreferences, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse preferences from JSONB
	prefs := &UserPreferences{
		Currency:             "USD",
		NotificationsEnabled: true,
		EmailNotifications:   true,
		PushNotifications:    true,
		SavingsReminders:     true,
		ReminderTime:         "09:00",
		Theme:                "light",
	}

	if user.Preferences != nil {
		// Convert map to JSON and back to struct for type safety
		prefsJSON, err := json.Marshal(user.Preferences)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal preferences: %w", err)
		}

		if err := json.Unmarshal(prefsJSON, prefs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal preferences: %w", err)
		}
	}

	return prefs, nil
}

// UpdatePreferences updates user preferences
func (s *UserService) UpdatePreferences(ctx context.Context, userID string, prefs UserPreferences) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Get existing user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Convert preferences to map for JSONB storage
	prefsMap := make(map[string]interface{})
	prefsJSON, err := json.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	if err := json.Unmarshal(prefsJSON, &prefsMap); err != nil {
		return fmt.Errorf("failed to unmarshal preferences: %w", err)
	}

	user.Preferences = prefsMap

	// Update in database
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	return nil
}

// DeleteAccount deletes user account and all associated data
func (s *UserService) DeleteAccount(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Delete user (cascade will handle all related data)
	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}
