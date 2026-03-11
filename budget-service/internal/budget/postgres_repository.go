package budget

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

// CreateBudget inserts a new budget record
func (r *postgresRepository) CreateBudget(ctx context.Context, budget *Budget) error {
	query := `
		INSERT INTO budgets (id, user_id, month, total_budget, total_spent, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		budget.ID,
		budget.UserID,
		budget.Month,
		budget.TotalBudget,
		budget.TotalSpent,
		budget.CreatedAt,
		budget.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create budget: %w", err)
	}
	
	return nil
}

// CreateBudgetCategory inserts a new budget category
func (r *postgresRepository) CreateBudgetCategory(ctx context.Context, category *BudgetCategory) error {
	query := `
		INSERT INTO budget_categories (id, budget_id, name, allocated_amount, spent_amount, color)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		category.ID,
		category.BudgetID,
		category.Name,
		category.AllocatedAmount,
		category.SpentAmount,
		category.Color,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create budget category: %w", err)
	}
	
	return nil
}

// GetBudgetByUserAndMonth retrieves a budget for a specific user and month
func (r *postgresRepository) GetBudgetByUserAndMonth(ctx context.Context, userID string, month time.Time) (*Budget, error) {
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
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}
	
	return &budget, nil
}

// GetBudgetByID retrieves a budget by its ID
func (r *postgresRepository) GetBudgetByID(ctx context.Context, budgetID string) (*Budget, error) {
	query := `
		SELECT id, user_id, month, total_budget, total_spent, created_at, updated_at
		FROM budgets
		WHERE id = $1
	`
	
	var budget Budget
	err := r.db.QueryRowContext(ctx, query, budgetID).Scan(
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
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}
	
	return &budget, nil
}

// GetCategoriesByBudgetID retrieves all categories for a budget
func (r *postgresRepository) GetCategoriesByBudgetID(ctx context.Context, budgetID string) ([]BudgetCategory, error) {
	query := `
		SELECT id, budget_id, name, allocated_amount, spent_amount, color
		FROM budget_categories
		WHERE budget_id = $1
		ORDER BY name
	`
	
	rows, err := r.db.QueryContext(ctx, query, budgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()
	
	var categories []BudgetCategory
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
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		cat.RemainingAmount = cat.AllocatedAmount - cat.SpentAmount
		categories = append(categories, cat)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}
	
	return categories, nil
}

// UpdateBudget updates a budget record
func (r *postgresRepository) UpdateBudget(ctx context.Context, budget *Budget) error {
	query := `
		UPDATE budgets
		SET total_budget = $1, total_spent = $2, updated_at = $3
		WHERE id = $4
	`
	
	_, err := r.db.ExecContext(ctx, query,
		budget.TotalBudget,
		budget.TotalSpent,
		budget.UpdatedAt,
		budget.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update budget: %w", err)
	}
	
	return nil
}

// UpdateBudgetCategory updates a budget category
func (r *postgresRepository) UpdateBudgetCategory(ctx context.Context, category *BudgetCategory) error {
	query := `
		UPDATE budget_categories
		SET name = $1, allocated_amount = $2, spent_amount = $3, color = $4
		WHERE id = $5
	`
	
	_, err := r.db.ExecContext(ctx, query,
		category.Name,
		category.AllocatedAmount,
		category.SpentAmount,
		category.Color,
		category.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	
	return nil
}

// CreateSpendingTransaction inserts a new spending transaction
func (r *postgresRepository) CreateSpendingTransaction(ctx context.Context, tx *SpendingTransaction) error {
	query := `
		INSERT INTO spending_transactions (id, user_id, budget_id, category_id, amount, description, merchant, transaction_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		tx.ID,
		tx.UserID,
		tx.BudgetID,
		tx.CategoryID,
		tx.Amount,
		tx.Description,
		tx.Merchant,
		tx.Date,
		tx.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create spending transaction: %w", err)
	}
	
	return nil
}

// UpdateCategorySpentAmount updates the spent amount for a category
func (r *postgresRepository) UpdateCategorySpentAmount(ctx context.Context, categoryID string, amount float64) error {
	query := `
		UPDATE budget_categories
		SET spent_amount = spent_amount + $1
		WHERE id = $2
	`
	
	_, err := r.db.ExecContext(ctx, query, amount, categoryID)
	if err != nil {
		return fmt.Errorf("failed to update category spent amount: %w", err)
	}
	
	return nil
}

// UpdateBudgetTotalSpent updates the total spent for a budget
func (r *postgresRepository) UpdateBudgetTotalSpent(ctx context.Context, budgetID string, amount float64) error {
	query := `
		UPDATE budgets
		SET total_spent = total_spent + $1, updated_at = $2
		WHERE id = $3
	`
	
	_, err := r.db.ExecContext(ctx, query, amount, time.Now().UTC(), budgetID)
	if err != nil {
		return fmt.Errorf("failed to update budget total spent: %w", err)
	}
	
	return nil
}

// GetCategoryByID retrieves a category by its ID
func (r *postgresRepository) GetCategoryByID(ctx context.Context, categoryID string) (*BudgetCategory, error) {
	query := `
		SELECT id, budget_id, name, allocated_amount, spent_amount, color
		FROM budget_categories
		WHERE id = $1
	`
	
	var cat BudgetCategory
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(
		&cat.ID,
		&cat.BudgetID,
		&cat.Name,
		&cat.AllocatedAmount,
		&cat.SpentAmount,
		&cat.Color,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	
	cat.RemainingAmount = cat.AllocatedAmount - cat.SpentAmount
	
	return &cat, nil
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

// CreateSpendingTransaction inserts a spending transaction within a transaction
func (t *postgresTransaction) CreateSpendingTransaction(ctx context.Context, spending *SpendingTransaction) error {
	query := `
		INSERT INTO spending_transactions (id, user_id, budget_id, category_id, amount, description, merchant, transaction_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := t.tx.ExecContext(ctx, query,
		spending.ID,
		spending.UserID,
		spending.BudgetID,
		spending.CategoryID,
		spending.Amount,
		spending.Description,
		spending.Merchant,
		spending.Date,
		spending.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create spending transaction: %w", err)
	}
	
	return nil
}

// UpdateCategorySpentAmount updates category spent amount within a transaction
func (t *postgresTransaction) UpdateCategorySpentAmount(ctx context.Context, categoryID string, amount float64) error {
	query := `
		UPDATE budget_categories
		SET spent_amount = spent_amount + $1
		WHERE id = $2
	`
	
	_, err := t.tx.ExecContext(ctx, query, amount, categoryID)
	if err != nil {
		return fmt.Errorf("failed to update category spent amount: %w", err)
	}
	
	return nil
}

// UpdateBudgetTotalSpent updates budget total spent within a transaction
func (t *postgresTransaction) UpdateBudgetTotalSpent(ctx context.Context, budgetID string, amount float64) error {
	query := `
		UPDATE budgets
		SET total_spent = total_spent + $1, updated_at = $2
		WHERE id = $3
	`
	
	_, err := t.tx.ExecContext(ctx, query, amount, time.Now().UTC(), budgetID)
	if err != nil {
		return fmt.Errorf("failed to update budget total spent: %w", err)
	}
	
	return nil
}

// GetCategoryByID retrieves a category within a transaction
func (t *postgresTransaction) GetCategoryByID(ctx context.Context, categoryID string) (*BudgetCategory, error) {
	query := `
		SELECT id, budget_id, name, allocated_amount, spent_amount, color
		FROM budget_categories
		WHERE id = $1
	`
	
	var cat BudgetCategory
	err := t.tx.QueryRowContext(ctx, query, categoryID).Scan(
		&cat.ID,
		&cat.BudgetID,
		&cat.Name,
		&cat.AllocatedAmount,
		&cat.SpentAmount,
		&cat.Color,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	
	cat.RemainingAmount = cat.AllocatedAmount - cat.SpentAmount
	
	return &cat, nil
}

// Commit commits the transaction
func (t *postgresTransaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *postgresTransaction) Rollback() error {
	return t.tx.Rollback()
}
