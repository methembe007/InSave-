package analytics

import (
	"context"
	"time"
)

// Service defines the analytics service interface
type Service interface {
	// GetSpendingAnalysis analyzes spending patterns for a given period
	GetSpendingAnalysis(ctx context.Context, userID string, period TimePeriod) (*SpendingAnalysis, error)
	
	// GetSavingsPatterns detects and returns savings patterns
	GetSavingsPatterns(ctx context.Context, userID string) ([]SavingsPattern, error)
	
	// GetRecommendations generates AI-assisted recommendations
	GetRecommendations(ctx context.Context, userID string) ([]Recommendation, error)
	
	// GenerateMonthlyReport creates a comprehensive monthly financial report
	GenerateMonthlyReport(ctx context.Context, userID string, month time.Time) (*MonthlyReport, error)
	
	// GetFinancialHealth calculates the user's financial health score
	GetFinancialHealth(ctx context.Context, userID string) (*FinancialHealthScore, error)
}

// Repository defines the data access interface for analytics operations
type Repository interface {
	// GetSpendingTransactions retrieves spending transactions for a period
	GetSpendingTransactions(ctx context.Context, userID string, start, end time.Time) ([]SpendingTransaction, error)
	
	// GetSavingsTransactions retrieves savings transactions for a period
	GetSavingsTransactions(ctx context.Context, userID string, start, end time.Time) ([]SavingsTransaction, error)
	
	// GetBudgetForMonth retrieves the budget for a specific month
	GetBudgetForMonth(ctx context.Context, userID string, month time.Time) (*Budget, error)
	
	// GetUserStreak retrieves the user's current and longest streak
	GetUserStreak(ctx context.Context, userID string) (currentStreak, longestStreak int, err error)
	
	// GetTotalSaved retrieves the total amount saved by the user
	GetTotalSaved(ctx context.Context, userID string) (float64, error)
	
	// GetMonthlyAverage retrieves the average monthly savings
	GetMonthlyAverage(ctx context.Context, userID string) (float64, error)
}

// SpendingTransaction represents a spending transaction from the database
type SpendingTransaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	BudgetID    string    `json:"budget_id"`
	CategoryID  string    `json:"category_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Merchant    string    `json:"merchant"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

// SavingsTransaction represents a savings transaction from the database
type SavingsTransaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

// Budget represents a budget from the database
type Budget struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	Month           time.Time        `json:"month"`
	TotalBudget     float64          `json:"total_budget"`
	TotalSpent      float64          `json:"total_spent"`
	RemainingBudget float64          `json:"remaining_budget"`
	Categories      []BudgetCategory `json:"categories"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// BudgetCategory represents a budget category
type BudgetCategory struct {
	ID              string  `json:"id"`
	BudgetID        string  `json:"budget_id"`
	Name            string  `json:"name"`
	AllocatedAmount float64 `json:"allocated_amount"`
	SpentAmount     float64 `json:"spent_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	Color           string  `json:"color"`
}
