package goal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type goalService struct {
	repo Repository
}

// NewGoalService creates a new goal service instance
func NewGoalService(repo Repository) Service {
	return &goalService{
		repo: repo,
	}
}

// GetActiveGoals retrieves all active goals for a user
func (s *goalService) GetActiveGoals(ctx context.Context, userID string) ([]Goal, error) {
	goals, err := s.repo.GetGoalsByUserAndStatus(ctx, userID, "active")
	if err != nil {
		return nil, fmt.Errorf("failed to get active goals: %w", err)
	}
	
	// Calculate progress percentage for each goal
	for i := range goals {
		goals[i].ProgressPercent = calculateProgressPercent(goals[i].CurrentAmount, goals[i].TargetAmount)
	}
	
	return goals, nil
}

// GetGoal retrieves a specific goal with its milestones
func (s *goalService) GetGoal(ctx context.Context, userID string, goalID string) (*GoalDetail, error) {
	goal, err := s.repo.GetGoalByID(ctx, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	// Verify the goal belongs to the user
	if goal.UserID != userID {
		return nil, fmt.Errorf("goal does not belong to user")
	}
	
	// Get milestones
	milestones, err := s.repo.GetMilestonesByGoalID(ctx, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones: %w", err)
	}
	
	// Calculate progress percentage
	goal.ProgressPercent = calculateProgressPercent(goal.CurrentAmount, goal.TargetAmount)
	
	return &GoalDetail{
		Goal:       *goal,
		Milestones: milestones,
	}, nil
}

// CreateGoal creates a new financial goal
func (s *goalService) CreateGoal(ctx context.Context, userID string, req CreateGoalRequest) (*Goal, error) {
	// Set default currency if not provided
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}
	
	// Create goal with initial values
	goal := &Goal{
		ID:            uuid.New().String(),
		UserID:        userID,
		Title:         req.Title,
		Description:   req.Description,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: 0, // Initialize to 0 as per requirement 9.2
		Currency:      currency,
		TargetDate:    req.TargetDate,
		Status:        "active", // Initialize to "active" as per requirement 9.2
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	
	// Calculate initial progress (will be 0)
	goal.ProgressPercent = calculateProgressPercent(goal.CurrentAmount, goal.TargetAmount)
	
	// Create goal in database
	if err := s.repo.CreateGoal(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}
	
	// Create milestones if provided
	for _, milestoneReq := range req.Milestones {
		milestone := &Milestone{
			ID:          uuid.New().String(),
			GoalID:      goal.ID,
			Title:       milestoneReq.Title,
			Amount:      milestoneReq.Amount,
			IsCompleted: false,
			Order:       milestoneReq.Order,
		}
		
		if err := s.repo.CreateMilestone(ctx, milestone); err != nil {
			return nil, fmt.Errorf("failed to create milestone: %w", err)
		}
	}
	
	return goal, nil
}

// UpdateGoal modifies an existing goal
func (s *goalService) UpdateGoal(ctx context.Context, userID string, goalID string, req UpdateGoalRequest) (*Goal, error) {
	// Get existing goal
	goal, err := s.repo.GetGoalByID(ctx, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	// Verify the goal belongs to the user
	if goal.UserID != userID {
		return nil, fmt.Errorf("goal does not belong to user")
	}
	
	// Update fields if provided
	if req.Title != "" {
		goal.Title = req.Title
	}
	if req.Description != "" {
		goal.Description = req.Description
	}
	if req.TargetAmount > 0 {
		goal.TargetAmount = req.TargetAmount
	}
	if !req.TargetDate.IsZero() {
		goal.TargetDate = req.TargetDate
	}
	if req.Status != "" {
		goal.Status = req.Status
	}
	
	// Update timestamp
	goal.UpdatedAt = time.Now()
	
	// Recalculate progress percentage
	goal.ProgressPercent = calculateProgressPercent(goal.CurrentAmount, goal.TargetAmount)
	
	// Update in database
	if err := s.repo.UpdateGoal(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}
	
	return goal, nil
}

// DeleteGoal removes a goal and all associated milestones
func (s *goalService) DeleteGoal(ctx context.Context, userID string, goalID string) error {
	// Get goal to verify ownership
	goal, err := s.repo.GetGoalByID(ctx, goalID)
	if err != nil {
		return fmt.Errorf("failed to get goal: %w", err)
	}
	
	// Verify the goal belongs to the user
	if goal.UserID != userID {
		return fmt.Errorf("goal does not belong to user")
	}
	
	// Delete goal (cascades to milestones via database constraint)
	if err := s.repo.DeleteGoal(ctx, goalID); err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}
	
	return nil
}

// GetMilestones retrieves all milestones for a goal
func (s *goalService) GetMilestones(ctx context.Context, goalID string) ([]Milestone, error) {
	milestones, err := s.repo.GetMilestonesByGoalID(ctx, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones: %w", err)
	}
	
	return milestones, nil
}

// UpdateProgress adds a contribution to a goal and updates milestones
// This implements requirements 10.1, 10.2, 10.3, 10.4, 10.5, 10.6, 16.4
func (s *goalService) UpdateProgress(ctx context.Context, goalID string, amount float64) (*Goal, error) {
	// Start database transaction for atomicity (requirement 16.4)
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed
	
	// Get goal with row-level lock to prevent race conditions (requirement 10.3)
	goal, err := tx.GetGoalByIDForUpdate(ctx, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal for update: %w", err)
	}
	
	// Only update active goals
	if goal.Status != "active" {
		return nil, fmt.Errorf("cannot update progress for non-active goal")
	}
	
	// Increase current amount by contribution (requirement 10.1)
	goal.CurrentAmount += amount
	
	// Check if goal is completed (requirement 10.2)
	if goal.CurrentAmount >= goal.TargetAmount {
		goal.Status = "completed"
	}
	
	// Calculate progress percentage (requirement 9.6)
	goal.ProgressPercent = calculateProgressPercent(goal.CurrentAmount, goal.TargetAmount)
	
	// Update timestamp
	goal.UpdatedAt = time.Now()
	
	// Update goal in database
	if err := tx.UpdateGoal(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}
	
	// Get uncompleted milestones ordered by amount (requirement 10.6)
	milestones, err := tx.GetUncompletedMilestones(ctx, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestones: %w", err)
	}
	
	// Update milestones that have been reached (requirements 10.4, 10.5, 10.6)
	// Process in ascending order by amount and stop at first unreached
	for _, milestone := range milestones {
		if goal.CurrentAmount >= milestone.Amount {
			// Mark milestone as completed with timestamp (requirement 10.5)
			now := time.Now()
			milestone.IsCompleted = true
			milestone.CompletedAt = &now
			
			if err := tx.UpdateMilestone(ctx, &milestone); err != nil {
				return nil, fmt.Errorf("failed to update milestone: %w", err)
			}
		} else {
			// Stop at first unreached milestone (requirement 10.6)
			break
		}
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return goal, nil
}

// calculateProgressPercent calculates progress as (current / target) × 100
func calculateProgressPercent(current, target float64) float64 {
	if target == 0 {
		return 0
	}
	return (current / target) * 100
}
