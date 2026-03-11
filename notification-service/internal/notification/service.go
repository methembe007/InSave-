package notification

import (
	"context"
)

// Service defines the notification service interface
type Service interface {
	// SendEmail sends an email notification using the configured email provider
	SendEmail(ctx context.Context, req EmailRequest) error
	
	// SendPushNotification sends a push notification to user's registered devices
	SendPushNotification(ctx context.Context, req PushNotificationRequest) error
	
	// ScheduleReminder schedules a reminder notification for future delivery
	ScheduleReminder(ctx context.Context, req ReminderRequest) error
	
	// GetUserNotifications retrieves all notifications for a user
	GetUserNotifications(ctx context.Context, userID string) ([]Notification, error)
	
	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, userID string, notificationID string) error
}

// Repository defines the data access interface for notification operations
type Repository interface {
	// GetUserPreferences retrieves notification preferences for a user
	GetUserPreferences(ctx context.Context, userID string) (*UserPreferences, error)
	
	// CreateNotification creates a new notification record
	CreateNotification(ctx context.Context, notification *Notification) error
	
	// GetUserNotifications retrieves all notifications for a user ordered by date descending
	GetUserNotifications(ctx context.Context, userID string) ([]Notification, error)
	
	// MarkAsRead updates the is_read flag for a notification
	MarkAsRead(ctx context.Context, userID string, notificationID string) error
}

// UserPreferences represents user notification preferences
type UserPreferences struct {
	NotificationsEnabled bool `json:"notifications_enabled"`
	EmailNotifications   bool `json:"email_notifications"`
	PushNotifications    bool `json:"push_notifications"`
}

// EmailProvider defines the interface for email delivery providers
type EmailProvider interface {
	// SendTemplateEmail sends an email using a template
	SendTemplateEmail(ctx context.Context, req EmailRequest) error
}

// PushProvider defines the interface for push notification providers
type PushProvider interface {
	// SendPush sends a push notification to user's devices
	SendPush(ctx context.Context, req PushNotificationRequest) error
}
