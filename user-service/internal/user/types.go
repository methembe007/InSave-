package user

import "time"

// UserProfile represents complete user profile information
type UserProfile struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	DateOfBirth     string    `json:"date_of_birth"`
	ProfileImageURL string    `json:"profile_image_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UpdateProfileRequest represents profile update data
type UpdateProfileRequest struct {
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	DateOfBirth     string `json:"date_of_birth,omitempty"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
}

// UserPreferences represents user settings and preferences
type UserPreferences struct {
	Currency             string `json:"currency"`
	NotificationsEnabled bool   `json:"notifications_enabled"`
	EmailNotifications   bool   `json:"email_notifications"`
	PushNotifications    bool   `json:"push_notifications"`
	SavingsReminders     bool   `json:"savings_reminders"`
	ReminderTime         string `json:"reminder_time"`
	Theme                string `json:"theme"`
}
