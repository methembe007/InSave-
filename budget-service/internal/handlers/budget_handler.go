package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/insavein/budget-service/internal/budget"
)

// BudgetHandler handles HTTP requests for budget operations
type BudgetHandler struct {
	service budget.Service
}

// NewBudgetHandler creates a new budget handler
func NewBudgetHandler(service budget.Service) *BudgetHandler {
	return &BudgetHandler{service: service}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// CreateBudget handles POST /api/budget
func (h *BudgetHandler) CreateBudget(w http.ResponseWriter, r *http.Request) {
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

	var req budget.CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Month.IsZero() {
		respondError(w, http.StatusBadRequest, "month is required")
		return
	}
	if len(req.Categories) == 0 {
		respondError(w, http.StatusBadRequest, "at least one category is required")
		return
	}

	budgetResult, err := h.service.CreateBudget(r.Context(), userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create budget")
		return
	}

	respondJSON(w, http.StatusCreated, budgetResult)
}

// GetCurrentBudget handles GET /api/budget/current
func (h *BudgetHandler) GetCurrentBudget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	budgetResult, err := h.service.GetCurrentBudget(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get current budget")
		return
	}

	if budgetResult == nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message": "no budget found for current month",
			"budget":  nil,
		})
		return
	}

	respondJSON(w, http.StatusOK, budgetResult)
}

// UpdateBudget handles PUT /api/budget/:id
func (h *BudgetHandler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Extract budget ID from URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		respondError(w, http.StatusBadRequest, "budget ID is required")
		return
	}
	budgetID := parts[len(parts)-1]

	var req budget.UpdateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if len(req.Categories) == 0 {
		respondError(w, http.StatusBadRequest, "at least one category is required")
		return
	}

	budgetResult, err := h.service.UpdateBudget(r.Context(), userID, budgetID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		if strings.Contains(err.Error(), "unauthorized") {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update budget")
		return
	}

	respondJSON(w, http.StatusOK, budgetResult)
}

// RecordSpending handles POST /api/budget/spending
func (h *BudgetHandler) RecordSpending(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req budget.SpendingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.CategoryID == "" {
		respondError(w, http.StatusBadRequest, "category_id is required")
		return
	}
	if req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}
	if req.Date.IsZero() {
		respondError(w, http.StatusBadRequest, "date is required")
		return
	}

	transaction, err := h.service.RecordSpending(r.Context(), userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "future") {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to record spending")
		return
	}

	respondJSON(w, http.StatusCreated, transaction)
}

// GetBudgetAlerts handles GET /api/budget/alerts
func (h *BudgetHandler) GetBudgetAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	alerts, err := h.service.CheckBudgetAlerts(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get budget alerts")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"alerts": alerts,
		"count":  len(alerts),
	})
}

// GetCategories handles GET /api/budget/categories
func (h *BudgetHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	categories, err := h.service.GetCategories(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get categories")
		return
	}

	respondJSON(w, http.StatusOK, categories)
}

// GetSpendingSummary handles GET /api/budget/summary
func (h *BudgetHandler) GetSpendingSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse month from query parameter (optional, defaults to current month)
	monthStr := r.URL.Query().Get("month")
	var month time.Time
	if monthStr != "" {
		var err error
		month, err = time.Parse("2006-01", monthStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid month format, use YYYY-MM")
			return
		}
	} else {
		month = time.Now().UTC()
	}

	summary, err := h.service.GetSpendingSummary(r.Context(), userID, month)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get spending summary")
		return
	}

	respondJSON(w, http.StatusOK, summary)
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
