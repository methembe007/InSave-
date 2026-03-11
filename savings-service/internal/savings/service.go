package savings

import (
	"context"
	"time"
)

// Service defines the interface for savings operations
type Service interface {
	// GetSummary returns a summary of the user's savings
	GetSummary(ctx context.Context, userID string) (*SavingsSummary, error)
	
	// GetHistory returns the user's savings transaction history
	GetHistory(ctx context.Context, userID string, params HistoryParams) ([]SavingsTransaction, error)
	
	// CreateTransaction creates a new savings transaction
	CreateTransaction(ctx context.Context, userID string, req CreateTransactionRequest) (*SavingsTransaction, error)
	
	// GetStreak returns the user's current savings streak
	GetStreak(ctx context.Context, userID string) (*SavingsStreak, error)
	
	// UpdateStreak recalculates and updates the user's savings streak
	UpdateStreak(ctx context.Context, userID string) error
	
	// GetMonthlyStats returns statistics for a specific month
	GetMonthlyStats(ctx context.Context, userID string, month time.Time) (*MonthlyStats, error)
}
