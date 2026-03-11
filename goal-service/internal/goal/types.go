package goal

import "time"

// Goal represents a financial goal
type Goal struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	TargetAmount    float64   `json:"target_amount"`
	CurrentAmount   float64   `json:"current_amount"`
	Currency        string    `json:"currency"`
	TargetDate      time.Time `json:"target_date"`
	Status          string    `json:"status"` // "active" | "completed" | "paused"
	ProgressPercent float64   `json:"progress_percent"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// GoalDetail represents a goal with its milestones
type GoalDetail struct {
	Goal
	Milestones []Milestone `json:"milestones"`
}

// Milestone represents an intermediate checkpoint toward a goal
type Milestone struct {
	ID          string     `json:"id"`
	GoalID      string     `json:"goal_id"`
	Title       string     `json:"title"`
	Amount      float64    `json:"amount"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Order       int        `json:"order"`
}

// CreateGoalRequest represents the request to create a goal
type CreateGoalRequest struct {
	Title        string    `json:"title" validate:"required"`
	Description  string    `json:"description"`
	TargetAmount float64   `json:"target_amount" validate:"required,gt=0"`
	Currency     string    `json:"currency"`
	TargetDate   time.Time `json:"target_date" validate:"required"`
	Milestones   []CreateMilestoneRequest `json:"milestones"`
}

// CreateMilestoneRequest represents a milestone in goal creation
type CreateMilestoneRequest struct {
	Title  string  `json:"title" validate:"required"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Order  int     `json:"order" validate:"required,gte=0"`
}

// UpdateGoalRequest represents the request to update a goal
type UpdateGoalRequest struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	TargetAmount float64   `json:"target_amount" validate:"omitempty,gt=0"`
	TargetDate   time.Time `json:"target_date"`
	Status       string    `json:"status" validate:"omitempty,oneof=active completed paused"`
}

// UpdateProgressRequest represents the request to update goal progress
type UpdateProgressRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}
