package notification

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// postgresRepository implements the Repository interface
type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository instance
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{
		db: db,
	}
}

// GetUserPreferences retrieves notification preferences for a user
func (r *postgresRepository) GetUserPreferences(ctx context.Context, userID string) (*UserPreferences, error) {
	query := `
		SELECT preferences
		FROM users
		WHERE id = $1
	`

	var prefsJSON []byte
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&prefsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user preferences: %w", err)
	}

	// Parse preferences JSON
	var prefs struct {
		NotificationsEnabled bool `json:"notifications_enabled"`
		EmailNotifications   bool `json:"email_notifications"`
		PushNotifications    bool `json:"push_notifications"`
	}

	if len(prefsJSON) > 0 {
		if err := json.Unmarshal(prefsJSON, &prefs); err != nil {
			return nil, fmt.Errorf("failed to parse preferences: %w", err)
		}
	} else {
		// Default to enabled if no preferences set
		prefs.NotificationsEnabled = true
		prefs.EmailNotifications = true
		prefs.PushNotifications = true
	}

	return &UserPreferences{
		NotificationsEnabled: prefs.NotificationsEnabled,
		EmailNotifications:   prefs.EmailNotifications,
		PushNotifications:    prefs.PushNotifications,
	}, nil
}

// CreateNotification creates a new notification record
func (r *postgresRepository) CreateNotification(ctx context.Context, notification *Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.IsRead,
		notification.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}

	return nil
}

// GetUserNotifications retrieves all notifications for a user ordered by date descending
// Requirement 12.4: Return notifications in date descending order
func (r *postgresRepository) GetUserNotifications(ctx context.Context, userID string) ([]Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifications: %w", err)
	}

	return notifications, nil
}

// MarkAsRead updates the is_read flag for a notification
// Requirement 12.5: Update is_read flag
func (r *postgresRepository) MarkAsRead(ctx context.Context, userID string, notificationID string) error {
	query := `
		UPDATE notifications
		SET is_read = true
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or does not belong to user")
	}

	return nil
}
