package savings

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

type savingsService struct {
	repo Repository
}

// NewService creates a new savings service
func NewService(repo Repository) Service {
	return &savingsService{
		repo: repo,
	}
}

// CreateTransaction creates a new savings transaction
func (s *savingsService) CreateTransaction(ctx context.Context, userID string, req CreateTransactionRequest) (*SavingsTransaction, error) {
	// Validate amount is positive
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	
	// Round amount to 2 decimal places
	roundedAmount := math.Round(req.Amount*100) / 100
	
	// Create transaction
	tx := &SavingsTransaction{
		ID:          uuid.New().String(),
		UserID:      userID,
		Amount:      roundedAmount,
		Currency:    req.Currency,
		Description: req.Description,
		Category:    req.Category,
		CreatedAt:   time.Now().UTC(),
	}
	
	// Insert into database
	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	
	// Trigger asynchronous streak update
	go func() {
		// Use a background context to avoid cancellation
		bgCtx := context.Background()
		if err := s.UpdateStreak(bgCtx, userID); err != nil {
			// Log error but don't fail the transaction
			// In production, this should use proper logging
			fmt.Printf("failed to update streak for user %s: %v\n", userID, err)
		}
	}()
	
	return tx, nil
}

// GetHistory returns the user's savings transaction history
func (s *savingsService) GetHistory(ctx context.Context, userID string, params HistoryParams) ([]SavingsTransaction, error) {
	// Set default limit if not provided
	if params.Limit <= 0 {
		params.Limit = 50
	}
	
	transactions, err := s.repo.GetTransactionsByUser(ctx, userID, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}
	
	return transactions, nil
}

// GetSummary returns a summary of the user's savings
func (s *savingsService) GetSummary(ctx context.Context, userID string) (*SavingsSummary, error) {
	// Get total saved
	totalSaved, err := s.repo.GetTotalSaved(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total saved: %w", err)
	}
	
	// Get current month total
	now := time.Now().UTC()
	thisMonthSaved, err := s.repo.GetMonthlyTotal(ctx, userID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get this month total: %w", err)
	}
	
	// Get monthly average
	monthlyAverage, err := s.repo.GetMonthlyAverage(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly average: %w", err)
	}
	
	// Get last saving date
	lastSavingDate, err := s.repo.GetLastSavingDate(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last saving date: %w", err)
	}
	
	// Get streak information
	currentStreak, longestStreak, err := s.repo.GetUserStreak(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get streak: %w", err)
	}
	
	summary := &SavingsSummary{
		TotalSaved:     totalSaved,
		CurrentStreak:  currentStreak,
		LongestStreak:  longestStreak,
		LastSavingDate: lastSavingDate,
		MonthlyAverage: monthlyAverage,
		ThisMonthSaved: thisMonthSaved,
	}
	
	return summary, nil
}

// GetStreak returns the user's current savings streak
func (s *savingsService) GetStreak(ctx context.Context, userID string) (*SavingsStreak, error) {
	currentStreak, longestStreak, err := s.repo.GetUserStreak(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get streak: %w", err)
	}
	
	lastSaveDate, err := s.repo.GetLastSavingDate(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last save date: %w", err)
	}
	
	streak := &SavingsStreak{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastSaveDate:  lastSaveDate,
	}
	
	return streak, nil
}

// GetMonthlyStats returns statistics for a specific month
func (s *savingsService) GetMonthlyStats(ctx context.Context, userID string, month time.Time) (*MonthlyStats, error) {
	totalSaved, err := s.repo.GetMonthlyTotal(ctx, userID, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly total: %w", err)
	}
	
	// Get count of transactions in the month
	firstDay := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Second)
	
	transactions, err := s.repo.GetTransactionsByUser(ctx, userID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	
	count := 0
	for _, tx := range transactions {
		if tx.CreatedAt.After(firstDay) && tx.CreatedAt.Before(lastDay) {
			count++
		}
	}
	
	averageAmount := 0.0
	if count > 0 {
		averageAmount = totalSaved / float64(count)
	}
	
	stats := &MonthlyStats{
		Month:         month,
		TotalSaved:    totalSaved,
		Count:         count,
		AverageAmount: averageAmount,
	}
	
	return stats, nil
}

// UpdateStreak recalculates and updates the user's savings streak
func (s *savingsService) UpdateStreak(ctx context.Context, userID string) error {
	// Get all saving dates
	dates, err := s.repo.GetAllSavingDates(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get saving dates: %w", err)
	}
	
	// If no transactions, set streaks to 0
	if len(dates) == 0 {
		return s.repo.UpdateUserStreak(ctx, userID, 0, 0)
	}
	
	// Calculate current streak
	currentStreak := 0
	longestStreak := 0
	tempStreak := 1
	
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastSaveDate := dates[0].Truncate(24 * time.Hour)
	
	// Check if last save was today or yesterday
	daysSinceLastSave := int(today.Sub(lastSaveDate).Hours() / 24)
	
	if daysSinceLastSave > 1 {
		// Streak is broken
		currentStreak = 0
	} else {
		// Start counting streak
		currentStreak = 1
		
		// Count consecutive days backwards
		for i := 1; i < len(dates); i++ {
			prevDate := dates[i-1].Truncate(24 * time.Hour)
			currDate := dates[i].Truncate(24 * time.Hour)
			dayDiff := int(prevDate.Sub(currDate).Hours() / 24)
			
			if dayDiff == 1 {
				// Consecutive day
				tempStreak++
				currentStreak = tempStreak
			} else if dayDiff > 1 {
				// Streak broken
				if tempStreak > longestStreak {
					longestStreak = tempStreak
				}
				tempStreak = 1
			}
			// dayDiff == 0 means multiple saves same day, continue streak
		}
	}
	
	// Final longest streak calculation
	if tempStreak > longestStreak {
		longestStreak = tempStreak
	}
	if currentStreak > longestStreak {
		longestStreak = currentStreak
	}
	
	// Get existing longest streak to ensure we never decrease it
	_, existingLongest, err := s.repo.GetUserStreak(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing streak: %w", err)
	}
	
	if existingLongest > longestStreak {
		longestStreak = existingLongest
	}
	
	// Update database
	return s.repo.UpdateUserStreak(ctx, userID, currentStreak, longestStreak)
}
