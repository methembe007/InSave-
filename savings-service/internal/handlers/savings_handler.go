package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/insavein/savings-service/internal/savings"
)

// SavingsHandler handles HTTP requests for savings operations
type SavingsHandler struct {
	service savings.Service
}

// NewSavingsHandler creates a new savings handler
func NewSavingsHandler(service savings.Service) *SavingsHandler {
	return &SavingsHandler{service: service}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// CreateTransaction handles POST /api/savings/transactions
func (h *SavingsHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
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

	var req savings.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}
	if req.Currency == "" {
		req.Currency = "USD" // Default currency
	}
	if len(req.Currency) != 3 {
		respondError(w, http.StatusBadRequest, "currency must be a 3-letter code")
		return
	}

	transaction, err := h.service.CreateTransaction(r.Context(), userID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create transaction")
		return
	}

	respondJSON(w, http.StatusCreated, transaction)
}

// GetHistory handles GET /api/savings/history
func (h *SavingsHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
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

	// Parse query parameters for pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // Default limit
	offset := 0 // Default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
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
func (h *SavingsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
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

	summary, err := h.service.GetSummary(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get summary")
		return
	}

	respondJSON(w, http.StatusOK, summary)
}

// GetStreak handles GET /api/savings/streak
func (h *SavingsHandler) GetStreak(w http.ResponseWriter, r *http.Request) {
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

	streak, err := h.service.GetStreak(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get streak")
		return
	}

	respondJSON(w, http.StatusOK, streak)
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
