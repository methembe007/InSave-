package savings

import (
	"context"
	"time"
)

// Repository defines the interface for savings data persistence
type Repository interface {
	// CreateTransaction inserts a new savings transaction
	CreateTransaction(ctx context.Context, tx *SavingsTransaction) error
	
	// GetTransactionsByUser retrieves all transactions for a user
	GetTransactionsByUser(ctx context.Context, userID string, limit, offset int) ([]SavingsTransaction, error)
	
	// GetTotalSaved calculates the total amount saved by a user
	GetTotalSaved(ctx context.Context, userID string) (float64, error)
	
	// GetMonthlyTotal calculates the total saved in a specific month
	GetMonthlyTotal(ctx context.Context, userID string, month time.Time) (float64, error)
	
	// GetMonthlyAverage calculates the average monthly savings
	GetMonthlyAverage(ctx context.Context, userID string) (float64, error)
	
	// GetLastSavingDate returns the date of the last savings transaction
	GetLastSavingDate(ctx context.Context, userID string) (time.Time, error)
	
	// GetAllSavingDates returns all unique dates when user saved money
	GetAllSavingDates(ctx context.Context, userID string) ([]time.Time, error)
	
	// UpdateUserStreak updates the streak information in user preferences
	UpdateUserStreak(ctx context.Context, userID string, currentStreak, longestStreak int) error
	
	// GetUserStreak retrieves the current streak information from user preferences
	GetUserStreak(ctx context.Context, userID string) (currentStreak, longestStreak int, err error)
}
