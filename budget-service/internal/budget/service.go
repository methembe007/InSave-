package budget

import (
	"context"
	"time"
)

// Service defines the budget service interface
type Service interface {
	// GetCurrentBudget retrieves the budget for the current month
	GetCurrentBudget(ctx context.Context, userID string) (*Budget, error)
	
	// CreateBudget creates a new monthly budget with category allocations
	CreateBudget(ctx context.Context, userID string, req CreateBudgetRequest) (*Budget, error)
	
	// UpdateBudget modifies an existing budget's category allocations
	UpdateBudget(ctx context.Context, userID string, budgetID string, req UpdateBudgetRequest) (*Budget, error)
	
	// GetCategories retrieves all categories for a user's budget
	GetCategories(ctx context.Context, userID string) ([]BudgetCategory, error)
	
	// RecordSpending records a spending transaction and updates budget amounts
	RecordSpending(ctx context.Context, userID string, req SpendingRequest) (*SpendingTransaction, error)
	
	// GetSpendingSummary returns spending summary for a specific month
	GetSpendingSummary(ctx context.Context, userID string, month time.Time) (*SpendingSummary, error)
	
	// CheckBudgetAlerts generates alerts for categories exceeding thresholds
	CheckBudgetAlerts(ctx context.Context, userID string) ([]BudgetAlert, error)
}

// Repository defines the data access interface for budget operations
type Repository interface {
	// CreateBudget inserts a new budget record
	CreateBudget(ctx context.Context, budget *Budget) error
	
	// CreateBudgetCategory inserts a new budget category
	CreateBudgetCategory(ctx context.Context, category *BudgetCategory) error
	
	// GetBudgetByUserAndMonth retrieves a budget for a specific user and month
	GetBudgetByUserAndMonth(ctx context.Context, userID string, month time.Time) (*Budget, error)
	
	// GetBudgetByID retrieves a budget by its ID
	GetBudgetByID(ctx context.Context, budgetID string) (*Budget, error)
	
	// GetCategoriesByBudgetID retrieves all categories for a budget
	GetCategoriesByBudgetID(ctx context.Context, budgetID string) ([]BudgetCategory, error)
	
	// UpdateBudget updates a budget record
	UpdateBudget(ctx context.Context, budget *Budget) error
	
	// UpdateBudgetCategory updates a budget category
	UpdateBudgetCategory(ctx context.Context, category *BudgetCategory) error
	
	// CreateSpendingTransaction inserts a new spending transaction
	CreateSpendingTransaction(ctx context.Context, tx *SpendingTransaction) error
	
	// UpdateCategorySpentAmount updates the spent amount for a category
	UpdateCategorySpentAmount(ctx context.Context, categoryID string, amount float64) error
	
	// UpdateBudgetTotalSpent updates the total spent for a budget
	UpdateBudgetTotalSpent(ctx context.Context, budgetID string, amount float64) error
	
	// GetCategoryByID retrieves a category by its ID
	GetCategoryByID(ctx context.Context, categoryID string) (*BudgetCategory, error)
	
	// BeginTx starts a database transaction
	BeginTx(ctx context.Context) (Transaction, error)
}

// Transaction represents a database transaction
type Transaction interface {
	// CreateSpendingTransaction inserts a spending transaction within a transaction
	CreateSpendingTransaction(ctx context.Context, tx *SpendingTransaction) error
	
	// UpdateCategorySpentAmount updates category spent amount within a transaction
	UpdateCategorySpentAmount(ctx context.Context, categoryID string, amount float64) error
	
	// UpdateBudgetTotalSpent updates budget total spent within a transaction
	UpdateBudgetTotalSpent(ctx context.Context, budgetID string, amount float64) error
	
	// GetCategoryByID retrieves a category within a transaction
	GetCategoryByID(ctx context.Context, categoryID string) (*BudgetCategory, error)
	
	// Commit commits the transaction
	Commit() error
	
	// Rollback rolls back the transaction
	Rollback() error
}
