package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/insavein/analytics-service/internal/analytics"
	"github.com/insavein/analytics-service/internal/middleware"
)

type AnalyticsHandler struct {
	service analytics.Service
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(service analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

// GetSpendingAnalysis handles GET /api/analytics/spending
func (h *AnalyticsHandler) GetSpendingAnalysis(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	
	// Parse period parameter (default to last 30 days)
	periodParam := r.URL.Query().Get("period")
	var period analytics.TimePeriod
	
	end := time.Now().UTC()
	
	switch periodParam {
	case "week":
		period = analytics.TimePeriod{
			Start: end.AddDate(0, 0, -7),
			End:   end,
		}
	case "month":
		period = analytics.TimePeriod{
			Start: end.AddDate(0, -1, 0),
			End:   end,
		}
	case "quarter":
		period = analytics.TimePeriod{
			Start: end.AddDate(0, -3, 0),
			End:   end,
		}
	case "year":
		period = analytics.TimePeriod{
			Start: end.AddDate(-1, 0, 0),
			End:   end,
		}
	default:
		// Default to last 30 days
		period = analytics.TimePeriod{
			Start: end.AddDate(0, 0, -30),
			End:   end,
		}
	}
	
	// Get spending analysis
	analysis, err := h.service.GetSpendingAnalysis(r.Context(), userID, period)
	if err != nil {
		http.Error(w, `{"error":"Failed to get spending analysis"}`, http.StatusInternalServerError)
		return
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}

// GetSavingsPatterns handles GET /api/analytics/patterns
func (h *AnalyticsHandler) GetSavingsPatterns(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	
	// Get savings patterns
	patterns, err := h.service.GetSavingsPatterns(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to get savings patterns"}`, http.StatusInternalServerError)
		return
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patterns)
}

// GetRecommendations handles GET /api/analytics/recommendations
func (h *AnalyticsHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	
	// Get recommendations
	recommendations, err := h.service.GetRecommendations(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to get recommendations"}`, http.StatusInternalServerError)
		return
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}

// GetFinancialHealth handles GET /api/analytics/health
func (h *AnalyticsHandler) GetFinancialHealth(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	
	// Get financial health score
	health, err := h.service.GetFinancialHealth(r.Context(), userID)
	if err != nil {
		// Check if it's an insufficient data error
		if err.Error()[:len("insufficient data")] == "insufficient data" {
			http.Error(w, `{"error":"Insufficient data: need at least 30 days of transaction history"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"error":"Failed to get financial health"}`, http.StatusInternalServerError)
		return
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}
