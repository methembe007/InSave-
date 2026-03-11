package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/insavein/notification-service/internal/notification"
)

// NotificationHandler handles HTTP requests for notification operations
type NotificationHandler struct {
	service notification.Service
}

// NewNotificationHandler creates a new notification handler instance
func NewNotificationHandler(service notification.Service) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

// GetUserNotifications handles GET /api/notifications
// Requirement 12.4: Retrieve user notifications
func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifications, err := h.service.GetUserNotifications(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get notifications for user %s: %v", userID, err)
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// MarkNotificationAsRead handles PUT /api/notifications/:id/read
// Requirement 12.5: Mark notification as read
func (h *NotificationHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get notification ID from URL parameters
	vars := mux.Vars(r)
	notificationID := vars["id"]
	if notificationID == "" {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	err := h.service.MarkAsRead(r.Context(), userID, notificationID)
	if err != nil {
		log.Printf("Failed to mark notification %s as read: %v", notificationID, err)
		http.Error(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification marked as read",
	})
}
