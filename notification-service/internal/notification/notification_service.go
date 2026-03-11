package notification

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// notificationService implements the Service interface
type notificationService struct {
	repo          Repository
	emailProvider EmailProvider
	pushProvider  PushProvider
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(repo Repository, emailProvider EmailProvider, pushProvider PushProvider) Service {
	return &notificationService{
		repo:          repo,
		emailProvider: emailProvider,
		pushProvider:  pushProvider,
	}
}

// SendEmail sends an email notification using the configured email provider
// Requirement 12.1: Email notification delivery
func (s *notificationService) SendEmail(ctx context.Context, req EmailRequest) error {
	// Validate request
	if req.To == "" {
		return fmt.Errorf("recipient email is required")
	}
	if req.Subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if req.TemplateID == "" {
		return fmt.Errorf("template ID is required")
	}

	// Send email using provider
	if err := s.emailProvider.SendTemplateEmail(ctx, req); err != nil {
		log.Printf("Failed to send email to %s: %v", req.To, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", req.To)
	return nil
}

// SendPushNotification sends a push notification to user's registered devices
// Requirement 12.2: Push notification delivery
func (s *notificationService) SendPushNotification(ctx context.Context, req PushNotificationRequest) error {
	// Validate request
	if req.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if req.Title == "" {
		return fmt.Errorf("notification title is required")
	}
	if req.Body == "" {
		return fmt.Errorf("notification body is required")
	}

	// Check user preferences before sending
	// Requirement 12.6: Notification preference enforcement
	prefs, err := s.repo.GetUserPreferences(ctx, req.UserID)
	if err != nil {
		log.Printf("Failed to get user preferences for %s: %v", req.UserID, err)
		return fmt.Errorf("failed to get user preferences: %w", err)
	}

	// Skip sending if notifications are disabled
	if !prefs.NotificationsEnabled || !prefs.PushNotifications {
		log.Printf("Push notifications disabled for user %s, skipping", req.UserID)
		return nil
	}

	// Send push notification using provider
	if err := s.pushProvider.SendPush(ctx, req); err != nil {
		log.Printf("Failed to send push notification to user %s: %v", req.UserID, err)
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	// Create notification record
	notification := &Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Type:      "push",
		Title:     req.Title,
		Message:   req.Body,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		log.Printf("Failed to create notification record: %v", err)
		// Don't return error as notification was sent successfully
	}

	log.Printf("Push notification sent successfully to user %s", req.UserID)
	return nil
}

// ScheduleReminder schedules a reminder notification for future delivery
// Requirement 12.3: Reminder scheduling
func (s *notificationService) ScheduleReminder(ctx context.Context, req ReminderRequest) error {
	// Validate request
	if req.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if req.Type == "" {
		return fmt.Errorf("reminder type is required")
	}
	if req.Message == "" {
		return fmt.Errorf("reminder message is required")
	}
	if req.ScheduledAt.IsZero() {
		return fmt.Errorf("scheduled time is required")
	}

	// Validate reminder type
	validTypes := map[string]bool{"savings": true, "budget": true, "goal": true}
	if !validTypes[req.Type] {
		return fmt.Errorf("invalid reminder type: %s", req.Type)
	}

	// Create notification record with scheduled time
	notification := &Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Type:      req.Type + "_reminder",
		Title:     fmt.Sprintf("%s Reminder", req.Type),
		Message:   req.Message,
		IsRead:    false,
		CreatedAt: req.ScheduledAt,
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return fmt.Errorf("failed to schedule reminder: %w", err)
	}

	log.Printf("Reminder scheduled for user %s at %s", req.UserID, req.ScheduledAt)
	return nil
}

// GetUserNotifications retrieves all notifications for a user
// Requirement 12.4: Notification history retrieval
func (s *notificationService) GetUserNotifications(ctx context.Context, userID string) ([]Notification, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	notifications, err := s.repo.GetUserNotifications(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read
// Requirement 12.5: Notification read status update
func (s *notificationService) MarkAsRead(ctx context.Context, userID string, notificationID string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if notificationID == "" {
		return fmt.Errorf("notification ID is required")
	}

	if err := s.repo.MarkAsRead(ctx, userID, notificationID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}
