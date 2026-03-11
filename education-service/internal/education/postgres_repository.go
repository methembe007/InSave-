package education

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// postgresRepository implements the Repository interface
type postgresRepository struct {
	db         *sql.DB // Primary database for writes
	replicaDB  *sql.DB // Read replica for reads
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB, replicaDB *sql.DB) Repository {
	return &postgresRepository{
		db:        db,
		replicaDB: replicaDB,
	}
}

// GetAllLessons retrieves all lessons from the read replica
func (r *postgresRepository) GetAllLessons(ctx context.Context) ([]Lesson, error) {
	query := `
		SELECT id, title, description, category, duration_minutes, difficulty, "order"
		FROM lessons
		ORDER BY "order" ASC
	`
	
	// Use read replica for lesson retrieval
	rows, err := r.replicaDB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query lessons: %w", err)
	}
	defer rows.Close()
	
	var lessons []Lesson
	for rows.Next() {
		var lesson Lesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.Description,
			&lesson.Category,
			&lesson.DurationMinutes,
			&lesson.Difficulty,
			&lesson.Order,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lesson: %w", err)
		}
		lessons = append(lessons, lesson)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lessons: %w", err)
	}
	
	return lessons, nil
}

// GetLessonByID retrieves a specific lesson by ID from the read replica
func (r *postgresRepository) GetLessonByID(ctx context.Context, lessonID string) (*LessonDetail, error) {
	query := `
		SELECT id, title, description, category, duration_minutes, difficulty, 
		       "order", content, video_url, resources
		FROM lessons
		WHERE id = $1
	`
	
	// Use read replica for lesson retrieval
	var lesson LessonDetail
	var videoURL sql.NullString
	var resourcesJSON []byte
	
	err := r.replicaDB.QueryRowContext(ctx, query, lessonID).Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.Description,
		&lesson.Category,
		&lesson.DurationMinutes,
		&lesson.Difficulty,
		&lesson.Order,
		&lesson.Content,
		&videoURL,
		&resourcesJSON,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("lesson not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query lesson: %w", err)
	}
	
	// Set video URL if present
	if videoURL.Valid {
		lesson.VideoURL = videoURL.String
	}
	
	// Parse resources JSON
	if len(resourcesJSON) > 0 {
		if err := json.Unmarshal(resourcesJSON, &lesson.Resources); err != nil {
			return nil, fmt.Errorf("failed to parse resources: %w", err)
		}
	}
	
	// Initialize empty array if nil
	if lesson.Resources == nil {
		lesson.Resources = []Resource{}
	}
	
	return &lesson, nil
}

// GetUserCompletions retrieves all lesson completions for a user
func (r *postgresRepository) GetUserCompletions(ctx context.Context, userID string) (map[string]bool, error) {
	query := `
		SELECT lesson_id, is_completed
		FROM education_progress
		WHERE user_id = $1 AND is_completed = true
	`
	
	// Use read replica for reading progress
	rows, err := r.replicaDB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query completions: %w", err)
	}
	defer rows.Close()
	
	completions := make(map[string]bool)
	for rows.Next() {
		var lessonID string
		var isCompleted bool
		if err := rows.Scan(&lessonID, &isCompleted); err != nil {
			return nil, fmt.Errorf("failed to scan completion: %w", err)
		}
		completions[lessonID] = isCompleted
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating completions: %w", err)
	}
	
	return completions, nil
}

// MarkLessonComplete creates or updates a lesson completion record
func (r *postgresRepository) MarkLessonComplete(ctx context.Context, userID string, lessonID string) error {
	query := `
		INSERT INTO education_progress (user_id, lesson_id, is_completed, completed_at)
		VALUES ($1, $2, true, $3)
		ON CONFLICT (user_id, lesson_id)
		DO UPDATE SET is_completed = true, completed_at = $3
	`
	
	// Use primary database for writes
	_, err := r.db.ExecContext(ctx, query, userID, lessonID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to mark lesson complete: %w", err)
	}
	
	return nil
}

// GetTotalLessonsCount returns the total number of lessons
func (r *postgresRepository) GetTotalLessonsCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM lessons`
	
	// Use read replica for counting
	var count int
	err := r.replicaDB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count lessons: %w", err)
	}
	
	return count, nil
}

// GetCompletedLessonsCount returns the number of completed lessons for a user
func (r *postgresRepository) GetCompletedLessonsCount(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM education_progress
		WHERE user_id = $1 AND is_completed = true
	`
	
	// Use read replica for counting
	var count int
	err := r.replicaDB.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count completed lessons: %w", err)
	}
	
	return count, nil
}
