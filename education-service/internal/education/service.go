package education

import (
	"context"
)

// Service defines the education service interface
type Service interface {
	// GetLessons retrieves all lessons with completion status for a user
	GetLessons(ctx context.Context, userID string) ([]Lesson, error)
	
	// GetLesson retrieves a specific lesson with detailed content
	GetLesson(ctx context.Context, userID string, lessonID string) (*LessonDetail, error)
	
	// MarkLessonComplete marks a lesson as completed for a user
	MarkLessonComplete(ctx context.Context, userID string, lessonID string) error
	
	// GetUserProgress calculates and returns user's education progress
	GetUserProgress(ctx context.Context, userID string) (*EducationProgress, error)
}

// Repository defines the data access interface for education operations
type Repository interface {
	// GetAllLessons retrieves all lessons from the database
	GetAllLessons(ctx context.Context) ([]Lesson, error)
	
	// GetLessonByID retrieves a specific lesson by ID
	GetLessonByID(ctx context.Context, lessonID string) (*LessonDetail, error)
	
	// GetUserCompletions retrieves all lesson completions for a user
	GetUserCompletions(ctx context.Context, userID string) (map[string]bool, error)
	
	// MarkLessonComplete creates or updates a lesson completion record
	MarkLessonComplete(ctx context.Context, userID string, lessonID string) error
	
	// GetTotalLessonsCount returns the total number of lessons
	GetTotalLessonsCount(ctx context.Context) (int, error)
	
	// GetCompletedLessonsCount returns the number of completed lessons for a user
	GetCompletedLessonsCount(ctx context.Context, userID string) (int, error)
}
