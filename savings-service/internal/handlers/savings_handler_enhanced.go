package handlers

import (
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/insavein/savings-service/internal/savings"
)

// SavingsHandlerEnhanced handles HTTP requests with validation and sanitization
type SavingsHandlerEnhanced struct {
	service  savings.Service
	validate *validator.Validate
}

// NewSavingsHandlerEnhanced creates a new enhanced savings handler
func NewSavingsHandlerEnhanced(service savings.Service) *SavingsHandlerEnhanced {
	return &SavingsHandlerEnhanced{
		service:  service,
		validate: validator.New(),
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents the error response
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// CreateTransaction handles POST /api/savings/transactions with validation
func (h *SavingsHandlerEnhanced) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req savings.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		h.respondValidationError(w, err)
		return
	}

	// Sanitize string inputs to prevent XSS and SQL injection
	req.Description = sanitizeString(req.Description)
	req.Category = sanitizeString(req.Category)

	// Create transaction
	transaction, err := h.service.CreateTransaction(r.Context(), userID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create transaction")
		return
	}

	respondJSON(w, http.StatusCreated, transaction)
}

// GetHistory handles GET /api/savings/history with validation
func (h *SavingsHandlerEnhanced) GetHistory(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse and validate query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // Default limit
	offset := 0 // Default offset

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 || l > 100 {
			respondError(w, http.StatusBadRequest, "limit must be between 1 and 100")
			return
		}
		limit = l
	}

	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil || o < 0 {
			respondError(w, http.StatusBadRequest, "offset must be non-negative")
			return
		}
		offset = o
	}

	params := savings.HistoryParams{
		Limit:  limit,
		Offset: offset,
	}

	history, err := h.service.GetHistory(r.Context(), userID, params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get history")
		return
	}

	respondJSON(w, http.StatusOK, history)
}

// GetSummary handles GET /api/savings/summary
func (h *SavingsHandlerEnhanced) GetSummary(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	summary, err := h.service.GetSummary(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get summary")
		return
	}

	respondJSON(w, http.StatusOK, summary)
}

// GetStreak handles GET /api/savings/streak
func (h *SavingsHandlerEnhanced) GetStreak(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	streak, err := h.service.GetStreak(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get streak")
		return
	}

	respondJSON(w, http.StatusOK, streak)
}

// sanitizeString removes potentially dangerous characters to prevent XSS and SQL injection
func sanitizeString(s string) string {
	// HTML escape to prevent XSS
	s = html.EscapeString(s)

	// Remove SQL injection patterns (basic protection - use parameterized queries as primary defense)
	s = strings.ReplaceAll(s, "--", "")
	s = strings.ReplaceAll(s, ";", "")
	s = strings.ReplaceAll(s, "/*", "")
	s = strings.ReplaceAll(s, "*/", "")
	s = strings.ReplaceAll(s, "xp_", "")
	s = strings.ReplaceAll(s, "sp_", "")

	return strings.TrimSpace(s)
}

// respondValidationError formats and sends validation errors
func (h *SavingsHandlerEnhanced) respondValidationError(w http.ResponseWriter, err error) {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Tag:     e.Tag(),
				Value:   e.Param(),
				Message: getErrorMessage(e),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ValidationErrorResponse{
		Error:   "Bad Request",
		Message: "Validation failed",
		Errors:  errors,
	}

	json.NewEncoder(w).Encode(response)
}

// getErrorMessage returns a user-friendly error message
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email address"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	case "gt":
		return e.Field() + " must be greater than " + e.Param()
	case "gte":
		return e.Field() + " must be greater than or equal to " + e.Param()
	case "lt":
		return e.Field() + " must be less than " + e.Param()
	case "lte":
		return e.Field() + " must be less than or equal to " + e.Param()
	case "len":
		return e.Field() + " must be exactly " + e.Param() + " characters"
	default:
		return e.Field() + " is invalid"
	}
}
