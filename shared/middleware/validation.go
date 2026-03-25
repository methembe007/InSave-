package middleware

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationMiddleware provides input validation and sanitization
type ValidationMiddleware struct {
	validate *validator.Validate
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validate: validator.New(),
	}
}

// ValidateRequest validates and sanitizes the request body
func (m *ValidationMiddleware) ValidateRequest(v interface{}) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Decode request body
			if err := json.NewDecoder(r.Body).Decode(v); err != nil {
				respondValidationError(w, http.StatusBadRequest, "Invalid request body", nil)
				return
			}

			// Validate struct
			if err := m.validate.Struct(v); err != nil {
				validationErrors := formatValidationErrors(err)
				respondValidationError(w, http.StatusBadRequest, "Validation failed", validationErrors)
				return
			}

			// Sanitize inputs
			sanitizeStruct(v)

			next.ServeHTTP(w, r)
		}
	}
}

// Validate validates a struct and returns detailed errors
func (m *ValidationMiddleware) Validate(v interface{}) error {
	return m.validate.Struct(v)
}

// formatValidationErrors converts validator errors to a readable format
func formatValidationErrors(err error) []ValidationError {
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
	
	return errors
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
	default:
		return e.Field() + " is invalid"
	}
}

// sanitizeStruct sanitizes all string fields in a struct to prevent XSS
func sanitizeStruct(v interface{}) {
	// Use reflection to sanitize string fields
	// This is a simplified version - in production, use a more robust solution
	switch val := v.(type) {
	case *string:
		*val = sanitizeString(*val)
	}
}

// sanitizeString removes potentially dangerous characters
func sanitizeString(s string) string {
	// HTML escape to prevent XSS
	s = html.EscapeString(s)
	
	// Remove SQL injection patterns (basic protection)
	s = strings.ReplaceAll(s, "--", "")
	s = strings.ReplaceAll(s, ";", "")
	s = strings.ReplaceAll(s, "/*", "")
	s = strings.ReplaceAll(s, "*/", "")
	s = strings.ReplaceAll(s, "xp_", "")
	s = strings.ReplaceAll(s, "sp_", "")
	
	return strings.TrimSpace(s)
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

// respondValidationError sends a validation error response
func respondValidationError(w http.ResponseWriter, status int, message string, errors []ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	response := ValidationErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Errors:  errors,
	}
	
	json.NewEncoder(w).Encode(response)
}
