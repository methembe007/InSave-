package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/insavein/education-service/internal/education"
)

// EducationHandler handles HTTP requests for education operations
type EducationHandler struct {
	service education.Service
}

// NewEducationHandler creates a new education handler
func NewEducationHandler(service education.Service) *EducationHandler {
	return &EducationHandler{
		service: service,
	}
}

// GetLessons handles GET /api/education/lessons
func (h *EducationHandler) GetLessons(w http.ResponseWriter, r *http.Request) {
	// Get user_id from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	
	// Get lessons
	lessons, err := h.service.GetLessons(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lessons)
}

// GetLesson handles GET /api/education/lessons/:id
func (h *EducationHandler) GetLesson(w http.ResponseWriter, r *http.Request) {
	// Get user_id from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	
	// Get lesson ID from URL
	vars := mux.Vars(r)
	lessonID := vars["id"]
	if lessonID == "" {
		http.Error(w, "Lesson ID is required", http.StatusBadRequest)
		return
	}
	
	// Get lesson
	lesson, err := h.service.GetLesson(r.Context(), userID, lessonID)
	if err != nil {
		if err.Error() == "lesson not found" {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lesson)
}

// MarkLessonComplete handles POST /api/education/lessons/:id/complete
func (h *EducationHandler) MarkLessonComplete(w http.ResponseWriter, r *http.Request) {
	// Get user_id from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	
	// Get lesson ID from URL
	vars := mux.Vars(r)
	lessonID := vars["id"]
	if lessonID == "" {
		http.Error(w, "Lesson ID is required", http.StatusBadRequest)
		return
	}
	
	// Mark lesson complete
	if err := h.service.MarkLessonComplete(r.Context(), userID, lessonID); err != nil {
		if err.Error() == "lesson not found: lesson not found" {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Lesson marked as complete",
	})
}

// GetUserProgress handles GET /api/education/progress
func (h *EducationHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	// Get user_id from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	
	// Get progress
	progress, err := h.service.GetUserProgress(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}
