package education

import "time"

// Lesson represents a financial education lesson
type Lesson struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	DurationMinutes int      `json:"duration_minutes"`
	Difficulty      string   `json:"difficulty"` // "beginner" | "intermediate" | "advanced"
	IsCompleted     bool     `json:"is_completed"`
	Order           int      `json:"order"`
}

// LessonDetail represents a lesson with full content
type LessonDetail struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Category        string     `json:"category"`
	DurationMinutes int        `json:"duration_minutes"`
	Difficulty      string     `json:"difficulty"`
	IsCompleted     bool       `json:"is_completed"`
	Order           int        `json:"order"`
	Content         string     `json:"content"`
	VideoURL        string     `json:"video_url,omitempty"`
	Resources       []Resource `json:"resources"`
}

// Resource represents an educational resource
type Resource struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"` // "article" | "video" | "pdf" | "link"
}

// EducationProgress represents a user's education progress
type EducationProgress struct {
	TotalLessons     int     `json:"total_lessons"`
	CompletedLessons int     `json:"completed_lessons"`
	ProgressPercent  float64 `json:"progress_percent"`
}

// LessonCompletion represents a lesson completion record
type LessonCompletion struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	LessonID    string     `json:"lesson_id"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
