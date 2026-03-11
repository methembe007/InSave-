package goal

import (
	"context"
	"database/sql"
	"fmt"
)

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

// CreateGoal inserts a new goal record
func (r *postgresRepository) CreateGoal(ctx context.Context, goal *Goal) error {
	query := `
		INSERT INTO goals (id, user_id, title, description, target_amount, current_amount, 
			currency, target_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		goal.ID,
		goal.UserID,
		goal.Title,
		goal.Description,
		goal.TargetAmount,
		goal.CurrentAmount,
		goal.Currency,
		goal.TargetDate,
		goal.Status,
		goal.CreatedAt,
		goal.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create goal: %w", err)
	}
	
	return nil
}

// CreateMilestone inserts a new milestone
func (r *postgresRepository) CreateMilestone(ctx context.Context, milestone *Milestone) error {
	query := `
		INSERT INTO goal_milestones (id, goal_id, title, amount, is_completed, completed_at, "order")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		milestone.ID,
		milestone.GoalID,
		milestone.Title,
		milestone.Amount,
		milestone.IsCompleted,
		milestone.CompletedAt,
		milestone.Order,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}
	
	return nil
}

// GetGoalByID retrieves a goal by its ID
func (r *postgresRepository) GetGoalByID(ctx context.Context, goalID string) (*Goal, error) {
	query := `
		SELECT id, user_id, title, description, target_amount, current_amount,
			currency, target_date, status, created_at, updated_at
		FROM goals
		WHERE id = $1
	`
	
	var goal Goal
	err := r.db.QueryRowContext(ctx, query, goalID).Scan(
		&goal.ID,
		&goal.UserID,
		&goal.Title,
		&goal.Description,
		&goal.TargetAmount,
		&goal.CurrentAmount,
		&goal.Currency,
		&goal.TargetDate,
		&goal.Status,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("goal not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	
	return &goal, nil
}

// GetGoalsByUserAndStatus retrieves goals for a user filtered by status
func (r *postgresRepository) GetGoalsByUserAndStatus(ctx context.Context, userID string, status string) ([]Goal, error) {
	query := `
		SELECT id, user_id, title, description, target_amount, current_amount,
			currency, target_date, status, created_at, updated_at
		FROM goals
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query goals: %w", err)
	}
	defer rows.Close()
	
	var goals []Goal
	for rows.Next() {
		var goal Goal
		err := rows.Scan(
			&goal.ID,
			&goal.UserID,
			&goal.Title,
			&goal.Description,
			&goal.TargetAmount,
			&goal.CurrentAmount,
			&goal.Currency,
			&goal.TargetDate,
			&goal.Status,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, goal)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating goals: %w", err)
	}
	
	return goals, nil
}

// GetMilestonesByGoalID retrieves all milestones for a goal
func (r *postgresRepository) GetMilestonesByGoalID(ctx context.Context, goalID string) ([]Milestone, error) {
	query := `
		SELECT id, goal_id, title, amount, is_completed, completed_at, "order"
		FROM goal_milestones
		WHERE goal_id = $1
		ORDER BY "order" ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to query milestones: %w", err)
	}
	defer rows.Close()
	
	var milestones []Milestone
	for rows.Next() {
		var milestone Milestone
		err := rows.Scan(
			&milestone.ID,
			&milestone.GoalID,
			&milestone.Title,
			&milestone.Amount,
			&milestone.IsCompleted,
			&milestone.CompletedAt,
			&milestone.Order,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}
		milestones = append(milestones, milestone)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating milestones: %w", err)
	}
	
	return milestones, nil
}

// UpdateGoal updates a goal record
func (r *postgresRepository) UpdateGoal(ctx context.Context, goal *Goal) error {
	query := `
		UPDATE goals
		SET title = $1, description = $2, target_amount = $3, current_amount = $4,
			currency = $5, target_date = $6, status = $7, updated_at = $8
		WHERE id = $9
	`
	
	_, err := r.db.ExecContext(ctx, query,
		goal.Title,
		goal.Description,
		goal.TargetAmount,
		goal.CurrentAmount,
		goal.Currency,
		goal.TargetDate,
		goal.Status,
		goal.UpdatedAt,
		goal.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}
	
	return nil
}

// DeleteGoal deletes a goal (cascades to milestones)
func (r *postgresRepository) DeleteGoal(ctx context.Context, goalID string) error {
	query := `DELETE FROM goals WHERE id = $1`
	
	_, err := r.db.ExecContext(ctx, query, goalID)
	if err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}
	
	return nil
}

// BeginTx starts a database transaction
func (r *postgresRepository) BeginTx(ctx context.Context) (Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	return &postgresTransaction{tx: tx}, nil
}

// postgresTransaction implements the Transaction interface
type postgresTransaction struct {
	tx *sql.Tx
}

// GetGoalByIDForUpdate retrieves a goal with row-level lock
func (t *postgresTransaction) GetGoalByIDForUpdate(ctx context.Context, goalID string) (*Goal, error) {
	query := `
		SELECT id, user_id, title, description, target_amount, current_amount,
			currency, target_date, status, created_at, updated_at
		FROM goals
		WHERE id = $1
		FOR UPDATE
	`
	
	var goal Goal
	err := t.tx.QueryRowContext(ctx, query, goalID).Scan(
		&goal.ID,
		&goal.UserID,
		&goal.Title,
		&goal.Description,
		&goal.TargetAmount,
		&goal.CurrentAmount,
		&goal.Currency,
		&goal.TargetDate,
		&goal.Status,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("goal not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get goal for update: %w", err)
	}
	
	return &goal, nil
}

// UpdateGoal updates a goal within a transaction
func (t *postgresTransaction) UpdateGoal(ctx context.Context, goal *Goal) error {
	query := `
		UPDATE goals
		SET title = $1, description = $2, target_amount = $3, current_amount = $4,
			currency = $5, target_date = $6, status = $7, updated_at = $8
		WHERE id = $9
	`
	
	_, err := t.tx.ExecContext(ctx, query,
		goal.Title,
		goal.Description,
		goal.TargetAmount,
		goal.CurrentAmount,
		goal.Currency,
		goal.TargetDate,
		goal.Status,
		goal.UpdatedAt,
		goal.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}
	
	return nil
}

// GetUncompletedMilestones retrieves uncompleted milestones ordered by amount
func (t *postgresTransaction) GetUncompletedMilestones(ctx context.Context, goalID string) ([]Milestone, error) {
	query := `
		SELECT id, goal_id, title, amount, is_completed, completed_at, "order"
		FROM goal_milestones
		WHERE goal_id = $1 AND is_completed = false
		ORDER BY amount ASC
	`
	
	rows, err := t.tx.QueryContext(ctx, query, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to query milestones: %w", err)
	}
	defer rows.Close()
	
	var milestones []Milestone
	for rows.Next() {
		var milestone Milestone
		err := rows.Scan(
			&milestone.ID,
			&milestone.GoalID,
			&milestone.Title,
			&milestone.Amount,
			&milestone.IsCompleted,
			&milestone.CompletedAt,
			&milestone.Order,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone: %w", err)
		}
		milestones = append(milestones, milestone)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating milestones: %w", err)
	}
	
	return milestones, nil
}

// UpdateMilestone updates a milestone within a transaction
func (t *postgresTransaction) UpdateMilestone(ctx context.Context, milestone *Milestone) error {
	query := `
		UPDATE goal_milestones
		SET title = $1, amount = $2, is_completed = $3, completed_at = $4, "order" = $5
		WHERE id = $6
	`
	
	_, err := t.tx.ExecContext(ctx, query,
		milestone.Title,
		milestone.Amount,
		milestone.IsCompleted,
		milestone.CompletedAt,
		milestone.Order,
		milestone.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}
	
	return nil
}

// Commit commits the transaction
func (t *postgresTransaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *postgresTransaction) Rollback() error {
	return t.tx.Rollback()
}
