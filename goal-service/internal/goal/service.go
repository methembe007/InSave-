package goal

import (
	"context"
)

// Service defines the goal service interface
type Service interface {
	// GetActiveGoals retrieves all active goals for a user
	GetActiveGoals(ctx context.Context, userID string) ([]Goal, error)
	
	// GetGoal retrieves a specific goal with its milestones
	GetGoal(ctx context.Context, userID string, goalID string) (*GoalDetail, error)
	
	// CreateGoal creates a new financial goal
	CreateGoal(ctx context.Context, userID string, req CreateGoalRequest) (*Goal, error)
	
	// UpdateGoal modifies an existing goal
	UpdateGoal(ctx context.Context, userID string, goalID string, req UpdateGoalRequest) (*Goal, error)
	
	// DeleteGoal removes a goal and all associated milestones
	DeleteGoal(ctx context.Context, userID string, goalID string) error
	
	// GetMilestones retrieves all milestones for a goal
	GetMilestones(ctx context.Context, goalID string) ([]Milestone, error)
	
	// UpdateProgress adds a contribution to a goal and updates milestones
	UpdateProgress(ctx context.Context, goalID string, amount float64) (*Goal, error)
}

// Repository defines the data access interface for goal operations
type Repository interface {
	// CreateGoal inserts a new goal record
	CreateGoal(ctx context.Context, goal *Goal) error
	
	// CreateMilestone inserts a new milestone
	CreateMilestone(ctx context.Context, milestone *Milestone) error
	
	// GetGoalByID retrieves a goal by its ID
	GetGoalByID(ctx context.Context, goalID string) (*Goal, error)
	
	// GetGoalsByUserAndStatus retrieves goals for a user filtered by status
	GetGoalsByUserAndStatus(ctx context.Context, userID string, status string) ([]Goal, error)
	
	// GetMilestonesByGoalID retrieves all milestones for a goal
	GetMilestonesByGoalID(ctx context.Context, goalID string) ([]Milestone, error)
	
	// UpdateGoal updates a goal record
	UpdateGoal(ctx context.Context, goal *Goal) error
	
	// DeleteGoal deletes a goal (cascades to milestones)
	DeleteGoal(ctx context.Context, goalID string) error
	
	// BeginTx starts a database transaction
	BeginTx(ctx context.Context) (Transaction, error)
}

// Transaction represents a database transaction for goal operations
type Transaction interface {
	// GetGoalByIDForUpdate retrieves a goal with row-level lock
	GetGoalByIDForUpdate(ctx context.Context, goalID string) (*Goal, error)
	
	// UpdateGoal updates a goal within a transaction
	UpdateGoal(ctx context.Context, goal *Goal) error
	
	// GetUncompletedMilestones retrieves uncompleted milestones ordered by amount
	GetUncompletedMilestones(ctx context.Context, goalID string) ([]Milestone, error)
	
	// UpdateMilestone updates a milestone within a transaction
	UpdateMilestone(ctx context.Context, milestone *Milestone) error
	
	// Commit commits the transaction
	Commit() error
	
	// Rollback rolls back the transaction
	Rollback() error
}
