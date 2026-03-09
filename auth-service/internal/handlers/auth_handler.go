package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/insavein/auth-service/internal/auth"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	service auth.Service
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Password == "" {
		respondError(w, http.StatusBadRequest, "password is required")
		return
	}
	if req.FirstName == "" {
		respondError(w, http.StatusBadRequest, "first_name is required")
		return
	}
	if req.LastName == "" {
		respondError(w, http.StatusBadRequest, "last_name is required")
		return
	}
	if req.DateOfBirth == "" {
		respondError(w, http.StatusBadRequest, "date_of_birth is required")
		return
	}

	response, err := h.service.Register(r.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "email already in use") {
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		if strings.Contains(err.Error(), "password must be at least") {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Password == "" {
		respondError(w, http.StatusBadRequest, "password is required")
		return
	}

	response, err := h.service.Login(r.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			respondError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if strings.Contains(err.Error(), "too many login attempts") {
			respondError(w, http.StatusTooManyRequests, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to login")
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// RefreshToken handles POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	response, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "revoked") {
			respondError(w, http.StatusUnauthorized, "invalid or expired refresh token")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to refresh token")
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	if err := h.service.Logout(r.Context(), userID, req.RefreshToken); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// ValidateToken handles GET /api/auth/validate (for internal use by other services)
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondError(w, http.StatusUnauthorized, "missing authorization header")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		respondError(w, http.StatusUnauthorized, "invalid authorization header format")
		return
	}

	token := parts[1]
	claims, err := h.service.ValidateToken(r.Context(), token)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	respondJSON(w, http.StatusOK, claims)
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
