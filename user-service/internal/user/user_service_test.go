package user

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	users map[string]*User
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users: make(map[string]*User),
	}
}

func (m *MockRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockRepository) UpdateUser(ctx context.Context, user *User) error {
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockRepository) DeleteUser(ctx context.Context, userID string) error {
	if _, exists := m.users[userID]; !exists {
		return errors.New("user not found")
	}
	delete(m.users, userID)
	return nil
}

func TestGetProfile(t *testing.T) {
	repo := NewMockRepository()
	service := NewUserService(repo)

	// Setup test user
	userID := "test-user-id"
	repo.users[userID] = &User{
		ID:              userID,
		Email:           "test@example.com",
		FirstName:       "John",
		LastName:        "Doe",
		DateOfBirth:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		ProfileImageURL: "https://example.com/image.jpg",
		Preferences:     make(map[string]interface{}),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Test successful profile retrieval
	profile, err := service.GetProfile(context.Background(), userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if profile.ID != userID {
		t.Errorf("Expected ID %s, got %s", userID, profile.ID)
	}
	if profile.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", profile.Email)
	}
	if profile.FirstName != "John" {
		t.Errorf("Expected first name John, got %s", profile.FirstName)
	}

	// Test with non-existent user
	_, err = service.GetProfile(context.Background(), "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	// Test with empty user ID
	_, err = service.GetProfile(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty user ID, got nil")
	}
}

func TestUpdateProfile(t *testing.T) {
	repo := NewMockRepository()
	service := NewUserService(repo)

	// Setup test user
	userID := "test-user-id"
	repo.users[userID] = &User{
		ID:              userID,
		Email:           "test@example.com",
		FirstName:       "John",
		LastName:        "Doe",
		DateOfBirth:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		ProfileImageURL: "",
		Preferences:     make(map[string]interface{}),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Test successful profile update
	req := UpdateProfileRequest{
		FirstName:       "Jane",
		LastName:        "Smith",
		ProfileImageURL: "https://example.com/new-image.jpg",
	}

	profile, err := service.UpdateProfile(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if profile.FirstName != "Jane" {
		t.Errorf("Expected first name Jane, got %s", profile.FirstName)
	}
	if profile.LastName != "Smith" {
		t.Errorf("Expected last name Smith, got %s", profile.LastName)
	}
	if profile.ProfileImageURL != "https://example.com/new-image.jpg" {
		t.Errorf("Expected profile image URL to be updated")
	}

	// Test with invalid date format
	req = UpdateProfileRequest{
		DateOfBirth: "invalid-date",
	}
	_, err = service.UpdateProfile(context.Background(), userID, req)
	if err == nil {
		t.Error("Expected error for invalid date format, got nil")
	}
}

func TestGetPreferences(t *testing.T) {
	repo := NewMockRepository()
	service := NewUserService(repo)

	// Setup test user with preferences
	userID := "test-user-id"
	repo.users[userID] = &User{
		ID:          userID,
		Email:       "test@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Preferences: map[string]interface{}{
			"currency":              "EUR",
			"notifications_enabled": true,
			"theme":                 "dark",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test successful preferences retrieval
	prefs, err := service.GetPreferences(context.Background(), userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if prefs.Currency != "EUR" {
		t.Errorf("Expected currency EUR, got %s", prefs.Currency)
	}
	if prefs.Theme != "dark" {
		t.Errorf("Expected theme dark, got %s", prefs.Theme)
	}
}

func TestUpdatePreferences(t *testing.T) {
	repo := NewMockRepository()
	service := NewUserService(repo)

	// Setup test user
	userID := "test-user-id"
	repo.users[userID] = &User{
		ID:          userID,
		Email:       "test@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Preferences: make(map[string]interface{}),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test successful preferences update
	prefs := UserPreferences{
		Currency:             "GBP",
		NotificationsEnabled: false,
		EmailNotifications:   false,
		PushNotifications:    true,
		SavingsReminders:     true,
		ReminderTime:         "10:00",
		Theme:                "dark",
	}

	err := service.UpdatePreferences(context.Background(), userID, prefs)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify preferences were updated
	updatedPrefs, err := service.GetPreferences(context.Background(), userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedPrefs.Currency != "GBP" {
		t.Errorf("Expected currency GBP, got %s", updatedPrefs.Currency)
	}
	if updatedPrefs.Theme != "dark" {
		t.Errorf("Expected theme dark, got %s", updatedPrefs.Theme)
	}
	if updatedPrefs.NotificationsEnabled != false {
		t.Error("Expected notifications_enabled to be false")
	}
}

func TestDeleteAccount(t *testing.T) {
	repo := NewMockRepository()
	service := NewUserService(repo)

	// Setup test user
	userID := "test-user-id"
	repo.users[userID] = &User{
		ID:          userID,
		Email:       "test@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Preferences: make(map[string]interface{}),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test successful account deletion
	err := service.DeleteAccount(context.Background(), userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify user was deleted
	_, err = service.GetProfile(context.Background(), userID)
	if err == nil {
		t.Error("Expected error for deleted user, got nil")
	}

	// Test deleting non-existent user
	err = service.DeleteAccount(context.Background(), "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	// Test with empty user ID
	err = service.DeleteAccount(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty user ID, got nil")
	}
}
