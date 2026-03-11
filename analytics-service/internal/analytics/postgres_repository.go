package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{
		db: db,
	}
}

// GetSpendingTransactions retrieves spending transactions for a period
func (r *postgresRepository) GetSpendingTransactions(ctx context.Context, userID string, start, end time.Time) ([]SpendingTransaction, error) {
	query := `
		SELECT id, user_id, budget_id, category_id, amount, description, merchant, transaction_date, created_at
		FROM spending_transactions
		WHERE user_id = $1 AND transaction_date >= $2 AND transaction_date <= $3
		ORDER BY transaction_date DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query spending transactions: %w", err)
	}
	defer rows.Close()
	
	transactions := []SpendingTransaction{}
	for rows.Next() {
		var tx SpendingTransaction
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.BudgetID,
			&tx.CategoryID,
			&tx.Amount,
			&tx.Description,
			&tx.Merchant,
			&tx.Date,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan spending transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating spending transactions: %w", err)
	}
	
	return transactions, nil
}

// GetSavingsTransactions retrieves savings transactions for a period
func (r *postgresRepository) GetSavingsTransactions(ctx context.Context, userID string, start, end time.Time) ([]SavingsTransaction, error) {
	var query string
	var rows *sql.Rows
	var err error
	
	// If start is zero time, get all transactions
	if start.IsZero() {
		query = `
			SELECT id, user_id, amount, currency, description, category, created_at
			FROM savings_transactions
			WHERE user_id = $1
			ORDER BY created_at DESC
		`
		rows, err = r.db.QueryContext(ctx, query, userID)
	} else {
		query = `
			SELECT id, user_id, amount, currency, description, category, created_at
			FROM savings_transactions
			WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3
			ORDER BY created_at DESC
		`
		rows, err = r.db.QueryContext(ctx, query, userID, start, end)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to query savings transactions: %w", err)
	}
	defer rows.Close()
	
	transactions := []SavingsTransaction{}
	for rows.Next() {
		var tx SavingsTransaction
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Amount,
			&tx.Currency,
			&tx.Description,
			&tx.Category,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan savings transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating savings transactions: %w", err)
	}
	
	return transactions, nil
}

// GetBudgetForMonth retrieves the budget for a specific month
func (r *postgresRepository) GetBudgetForMonth(ctx context.Context, userID string, month time.Time) (*Budget, error) {
	// Normalize to first day of month
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	
	query := `
		SELECT id, user_id, month, total_budget, total_spent, created_at, updated_at
		FROM budgets
		WHERE user_id = $1 AND month = $2
	`
	
	var budget Budget
	err := r.db.QueryRowContext(ctx, query, userID, month).Scan(
		&budget.ID,
		&budget.UserID,
		&budget.Month,
		&budget.TotalBudget,
		&budget.TotalSpent,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query budget: %w", err)
	}
	
	// Get categories for the budget
	categoriesQuery := `
		SELECT id, budget_id, name, allocated_amount, spent_amount, color
		FROM budget_categories
		WHERE budget_id = $1
	`
	
	rows, err := r.db.QueryContext(ctx, categoriesQuery, budget.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query budget categories: %w", err)
	}
	defer rows.Close()
	
	categories := []BudgetCategory{}
	for rows.Next() {
		var cat BudgetCategory
		err := rows.Scan(
			&cat.ID,
			&cat.BudgetID,
			&cat.Name,
			&cat.AllocatedAmount,
			&cat.SpentAmount,
			&cat.Color,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget category: %w", err)
		}
		cat.RemainingAmount = cat.AllocatedAmount - cat.SpentAmount
		categories = append(categories, cat)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating budget categories: %w", err)
	}
	
	budget.Categories = categories
	budget.RemainingBudget = budget.TotalBudget - budget.TotalSpent
	
	return &budget, nil
}

// GetUserStreak retrieves the user's current and longest streak
func (r *postgresRepository) GetUserStreak(ctx context.Context, userID string) (currentStreak, longestStreak int, err error) {
	query := `
		SELECT 
			COALESCE((preferences->>'current_streak')::int, 0) as current_streak,
			COALESCE((preferences->>'longest_streak')::int, 0) as longest_streak
		FROM users
		WHERE id = $1
	`
	
	err = r.db.QueryRowContext(ctx, query, userID).Scan(&currentStreak, &longestStreak)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query user streak: %w", err)
	}
	
	return currentStreak, longestStreak, nil
}

// GetTotalSaved retrieves the total amount saved by the user
func (r *postgresRepository) GetTotalSaved(ctx context.Context, userID string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM savings_transactions
		WHERE user_id = $1
	`
	
	var total float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to query total saved: %w", err)
	}
	
	return total, nil
}

// GetMonthlyAverage retrieves the average monthly savings
func (r *postgresRepository) GetMonthlyAverage(ctx context.Context, userID string) (float64, error) {
	query := `
		SELECT COALESCE(AVG(monthly_total), 0)
		FROM (
			SELECT DATE_TRUNC('month', created_at) as month, SUM(amount) as monthly_total
			FROM savings_transactions
			WHERE user_id = $1
			GROUP BY DATE_TRUNC('month', created_at)
		) as monthly_totals
	`
	
	var average float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&average)
	if err != nil {
		return 0, fmt.Errorf("failed to query monthly average: %w", err)
	}
	
	return average, nil
}
