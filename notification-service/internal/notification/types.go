package notification

import "time"

// EmailRequest represents a request to send an email notification
type EmailRequest struct {
	To           string                 `json:"to"`
	Subject      string                 `json:"subject"`
	TemplateID   string                 `json:"template_id"`
	TemplateData map[string]interface{} `json:"template_data"`
}

// PushNotificationRequest represents a request to send a push notification
type PushNotificationRequest struct {
	UserID string            `json:"user_id"`
	Title  string            `json:"title"`
	Body   string            `json:"body"`
	Data   map[string]string `json:"data"`
}

// ReminderRequest represents a request to schedule a reminder
type ReminderRequest struct {
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"` // "savings" | "budget" | "goal"
	ScheduledAt time.Time `json:"scheduled_at"`
	Message     string    `json:"message"`
}

// Notification represents a notification record
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
