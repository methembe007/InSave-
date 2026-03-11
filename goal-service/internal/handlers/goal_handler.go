package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/insavein/goal-service/internal/goal"
)

type GoalHandler struct {
	service  goal.Service
	validate *validator.Validate
}

// NewGoalHandler creates a new goal handler
func NewGoalHandler(service goal.Service) *GoalHandler {
	return &GoalHandler{
		service:  service,
		validate: validator.New(),
	}
}

// CreateGoal handles POST /api/goals
func (h *GoalHandler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	// Parse request body
	var req goal.CreateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate request
	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	// Create goal
	createdGoal, err := h.service.CreateGoal(r.Context(), userID, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusCreated, createdGoal)
}

// GetActiveGoals handles GET /api/goals
func (h *GoalHandler) GetActiveGoals(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	// Get active goals
	goals, err := h.service.GetActiveGoals(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, goals)
}

// GetGoal handles GET /api/goals/:id
func (h *GoalHandler) GetGoal(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	// Get goal ID from URL
	vars := mux.Vars(r)
	goalID := vars["id"]
	
	// Get goal
	goalDetail, err := h.service.GetGoal(r.Context(), userID, goalID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, goalDetail)
}

// UpdateGoal handles PUT /api/goals/:id
func (h *GoalHandler) UpdateGoal(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	// Get goal ID from URL
	vars := mux.Vars(r)
	goalID := vars["id"]
	
	// Parse request body
	var req goal.UpdateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate request
	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	// Update goal
	updatedGoal, err := h.service.UpdateGoal(r.Context(), userID, goalID, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, updatedGoal)
}

// DeleteGoal handles DELETE /api/goals/:id
func (h *GoalHandler) DeleteGoal(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	// Get goal ID from URL
	vars := mux.Vars(r)
	goalID := vars["id"]
	
	// Delete goal
	if err := h.service.DeleteGoal(r.Context(), userID, goalID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Goal deleted successfully"})
}

// UpdateProgress handles POST /api/goals/:id/progress
func (h *GoalHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	// Get goal ID from URL
	vars := mux.Vars(r)
	goalID := vars["id"]
	
	// Parse request body
	var req goal.UpdateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate request
	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	// Update progress
	updatedGoal, err := h.service.UpdateProgress(r.Context(), goalID, req.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, updatedGoal)
}

// GetMilestones handles GET /api/goals/:id/milestones
func (h *GoalHandler) GetMilestones(w http.ResponseWriter, r *http.Request) {
	// Get goal ID from URL
	vars := mux.Vars(r)
	goalID := vars["id"]
	
	// Get milestones
	milestones, err := h.service.GetMilestones(r.Context(), goalID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, milestones)
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
