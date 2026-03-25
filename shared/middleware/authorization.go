package middleware

import (
	"context"
	"encoding/json"
	"net/http"
)

// AuthorizationMiddleware checks user ownership of resources
type AuthorizationMiddleware struct{}

// NewAuthorizationMiddleware creates a new authorization middleware
func NewAuthorizationMiddleware() *AuthorizationMiddleware {
	return &AuthorizationMiddleware{}
}

// ResourceOwnershipChecker is a function that checks if a user owns a resource
type ResourceOwnershipChecker func(ctx context.Context, userID string, resourceID string) (bool, error)

// RequireOwnership ensures the authenticated user owns the requested resource
func (m *AuthorizationMiddleware) RequireOwnership(checker ResourceOwnershipChecker, resourceIDParam string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context (set by auth middleware)
			userID, ok := r.Context().Value("user_id").(string)
			if !ok || userID == "" {
				respondAuthError(w, http.StatusUnauthorized, "Unauthorized: user not authenticated")
				return
			}

			// Get resource ID from request (URL param, query param, or body)
			resourceID := getResourceID(r, resourceIDParam)
			if resourceID == "" {
				respondAuthError(w, http.StatusBadRequest, "Resource ID not provided")
				return
			}

			// Check ownership
			isOwner, err := checker(r.Context(), userID, resourceID)
			if err != nil {
				respondAuthError(w, http.StatusInternalServerError, "Failed to verify resource ownership")
				return
			}

			if !isOwner {
				respondAuthError(w, http.StatusForbidden, "Forbidden: you do not have access to this resource")
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequireSelfOrAdmin ensures the user is accessing their own data or is an admin
func (m *AuthorizationMiddleware) RequireSelfOrAdmin(targetUserIDParam string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get authenticated user ID from context
			userID, ok := r.Context().Value("user_id").(string)
			if !ok || userID == "" {
				respondAuthError(w, http.StatusUnauthorized, "Unauthorized: user not authenticated")
				return
			}

			// Get target user ID from request
			targetUserID := getResourceID(r, targetUserIDParam)
			if targetUserID == "" {
				respondAuthError(w, http.StatusBadRequest, "Target user ID not provided")
				return
			}

			// Check if user is accessing their own data
			if userID != targetUserID {
				// Check if user has admin role
				roles, ok := r.Context().Value("roles").([]string)
				if !ok || !contains(roles, "admin") {
					respondAuthError(w, http.StatusForbidden, "Forbidden: you can only access your own data")
					return
				}
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequireRole ensures the user has a specific role
func (m *AuthorizationMiddleware) RequireRole(requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get roles from context
			roles, ok := r.Context().Value("roles").([]string)
			if !ok {
				respondAuthError(w, http.StatusForbidden, "Forbidden: no roles found")
				return
			}

			// Check if user has required role
			if !contains(roles, requiredRole) {
				respondAuthError(w, http.StatusForbidden, "Forbidden: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// getResourceID extracts resource ID from request
func getResourceID(r *http.Request, paramName string) string {
	// Try URL path parameter (mux.Vars)
	if vars, ok := r.Context().Value("vars").(map[string]string); ok {
		if id, exists := vars[paramName]; exists {
			return id
		}
	}

	// Try query parameter
	if id := r.URL.Query().Get(paramName); id != "" {
		return id
	}

	// Try form value
	if id := r.FormValue(paramName); id != "" {
		return id
	}

	return ""
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// respondAuthError sends an authorization error response
func respondAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]string{
		"error":   http.StatusText(status),
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}
