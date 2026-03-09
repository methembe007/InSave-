package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/insavein/user-service/internal/user"
)

// UserHandler handles HTTP requests for user profile operations
type UserHandler struct {
	service user.Service
}

// NewUserHandler creates a new user handler
func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service: service}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// GetProfile handles GET /api/user/profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get profile")
		return
	}

	respondJSON(w, http.StatusOK, profile)
}

// UpdateProfile handles PUT /api/user/profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req user.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.service.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		if strings.Contains(err.Error(), "invalid date") {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	respondJSON(w, http.StatusOK, profile)
}

// GetPreferences handles GET /api/user/preferences
func (h *UserHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	preferences, err := h.service.GetPreferences(r.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get preferences")
		return
	}

	respondJSON(w, http.StatusOK, preferences)
}

// UpdatePreferences handles PUT /api/user/preferences
func (h *UserHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var prefs user.UserPreferences
	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.UpdatePreferences(r.Context(), userID, prefs); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update preferences")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "preferences updated successfully"})
}

// DeleteAccount handles DELETE /api/user/account
func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.DeleteAccount(r.Context(), userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete account")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "account deleted successfully"})
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
