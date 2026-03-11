package education

import (
	"context"
	"fmt"
)

// educationService implements the Service interface
type educationService struct {
	repo Repository
}

// NewEducationService creates a new education service instance
func NewEducationService(repo Repository) Service {
	return &educationService{
		repo: repo,
	}
}

// GetLessons retrieves all lessons with completion status for a user
func (s *educationService) GetLessons(ctx context.Context, userID string) ([]Lesson, error) {
	// Get all lessons
	lessons, err := s.repo.GetAllLessons(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lessons: %w", err)
	}
	
	// Get user's completions
	completions, err := s.repo.GetUserCompletions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user completions: %w", err)
	}
	
	// Mark lessons as completed based on user's progress
	for i := range lessons {
		if completions[lessons[i].ID] {
			lessons[i].IsCompleted = true
		}
	}
	
	return lessons, nil
}

// GetLesson retrieves a specific lesson with detailed content
func (s *educationService) GetLesson(ctx context.Context, userID string, lessonID string) (*LessonDetail, error) {
	// Get lesson details
	lesson, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}
	
	// Get user's completions
	completions, err := s.repo.GetUserCompletions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user completions: %w", err)
	}
	
	// Set completion status
	lesson.IsCompleted = completions[lessonID]
	
	return lesson, nil
}

// MarkLessonComplete marks a lesson as completed for a user
func (s *educationService) MarkLessonComplete(ctx context.Context, userID string, lessonID string) error {
	// Verify lesson exists
	_, err := s.repo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return fmt.Errorf("lesson not found: %w", err)
	}
	
	// Mark as complete
	if err := s.repo.MarkLessonComplete(ctx, userID, lessonID); err != nil {
		return fmt.Errorf("failed to mark lesson complete: %w", err)
	}
	
	return nil
}

// GetUserProgress calculates and returns user's education progress
func (s *educationService) GetUserProgress(ctx context.Context, userID string) (*EducationProgress, error) {
	// Get total lessons count
	totalLessons, err := s.repo.GetTotalLessonsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total lessons count: %w", err)
	}
	
	// Get completed lessons count
	completedLessons, err := s.repo.GetCompletedLessonsCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed lessons count: %w", err)
	}
	
	// Calculate progress percentage
	var progressPercent float64
	if totalLessons > 0 {
		progressPercent = (float64(completedLessons) / float64(totalLessons)) * 100
	}
	
	return &EducationProgress{
		TotalLessons:     totalLessons,
		CompletedLessons: completedLessons,
		ProgressPercent:  progressPercent,
	}, nil
}
