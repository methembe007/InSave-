package budget

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

type budgetService struct {
	repo Repository
}

// NewService creates a new budget service
func NewService(repo Repository) Service {
	return &budgetService{
		repo: repo,
	}
}

// GetCurrentBudget retrieves the budget for the current month
func (s *budgetService) GetCurrentBudget(ctx context.Context, userID string) (*Budget, error) {
	now := time.Now().UTC()
	// Get first day of current month
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	
	budget, err := s.repo.GetBudgetByUserAndMonth(ctx, userID, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get current budget: %w", err)
	}
	
	// Get categories for the budget
	if budget != nil {
		categories, err := s.repo.GetCategoriesByBudgetID(ctx, budget.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get categories: %w", err)
		}
		
		// Calculate remaining amounts for each category
		for i := range categories {
			categories[i].RemainingAmount = categories[i].AllocatedAmount - categories[i].SpentAmount
		}
		
		budget.Categories = categories
		budget.RemainingBudget = budget.TotalBudget - budget.TotalSpent
	}
	
	return budget, nil
}

// CreateBudget creates a new monthly budget with category allocations
func (s *budgetService) CreateBudget(ctx context.Context, userID string, req CreateBudgetRequest) (*Budget, error) {
	// Validate all amounts are non-negative
	totalBudget := 0.0
	for _, cat := range req.Categories {
		if cat.AllocatedAmount < 0 {
			return nil, fmt.Errorf("allocated amount must be non-negative")
		}
		totalBudget += cat.AllocatedAmount
	}
	
	// Normalize month to first day of month
	month := time.Date(req.Month.Year(), req.Month.Month(), 1, 0, 0, 0, 0, time.UTC)
	
	// Check if budget already exists for this user and month (unique constraint)
	existing, err := s.repo.GetBudgetByUserAndMonth(ctx, userID, month)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing budget: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("budget already exists for this month")
	}
	
	// Create budget
	budget := &Budget{
		ID:          uuid.New().String(),
		UserID:      userID,
		Month:       month,
		TotalBudget: math.Round(totalBudget*100) / 100,
		TotalSpent:  0,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	
	if err := s.repo.CreateBudget(ctx, budget); err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}
	
	// Create categories
	categories := make([]BudgetCategory, 0, len(req.Categories))
	for _, catReq := range req.Categories {
		category := BudgetCategory{
			ID:              uuid.New().String(),
			BudgetID:        budget.ID,
			Name:            catReq.Name,
			AllocatedAmount: math.Round(catReq.AllocatedAmount*100) / 100,
			SpentAmount:     0,
			RemainingAmount: math.Round(catReq.AllocatedAmount*100) / 100,
			Color:           catReq.Color,
		}
		
		if category.Color == "" {
			category.Color = "#000000"
		}
		
		if err := s.repo.CreateBudgetCategory(ctx, &category); err != nil {
			return nil, fmt.Errorf("failed to create category: %w", err)
		}
		
		categories = append(categories, category)
	}
	
	budget.Categories = categories
	budget.RemainingBudget = budget.TotalBudget
	
	return budget, nil
}

// UpdateBudget modifies an existing budget's category allocations
func (s *budgetService) UpdateBudget(ctx context.Context, userID string, budgetID string, req UpdateBudgetRequest) (*Budget, error) {
	// Get existing budget
	budget, err := s.repo.GetBudgetByID(ctx, budgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}
	
	if budget == nil {
		return nil, fmt.Errorf("budget not found")
	}
	
	// Verify ownership
	if budget.UserID != userID {
		return nil, fmt.Errorf("unauthorized: budget does not belong to user")
	}
	
	// Validate all amounts are non-negative and update categories
	totalBudget := 0.0
	for _, catReq := range req.Categories {
		if catReq.AllocatedAmount < 0 {
			return nil, fmt.Errorf("allocated amount must be non-negative")
		}
		
		// Get existing category to preserve spent amount
		category, err := s.repo.GetCategoryByID(ctx, catReq.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get category: %w", err)
		}
		
		if category == nil || category.BudgetID != budgetID {
			return nil, fmt.Errorf("category not found or does not belong to budget")
		}
		
		// Update category
		category.Name = catReq.Name
		category.AllocatedAmount = math.Round(catReq.AllocatedAmount*100) / 100
		category.RemainingAmount = category.AllocatedAmount - category.SpentAmount
		if catReq.Color != "" {
			category.Color = catReq.Color
		}
		
		if err := s.repo.UpdateBudgetCategory(ctx, category); err != nil {
			return nil, fmt.Errorf("failed to update category: %w", err)
		}
		
		totalBudget += category.AllocatedAmount
	}
	
	// Update budget totals
	budget.TotalBudget = math.Round(totalBudget*100) / 100
	budget.RemainingBudget = budget.TotalBudget - budget.TotalSpent
	budget.UpdatedAt = time.Now().UTC()
	
	if err := s.repo.UpdateBudget(ctx, budget); err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}
	
	// Get updated categories
	categories, err := s.repo.GetCategoriesByBudgetID(ctx, budgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	
	budget.Categories = categories
	
	return budget, nil
}

// GetCategories retrieves all categories for a user's current budget
func (s *budgetService) GetCategories(ctx context.Context, userID string) ([]BudgetCategory, error) {
	budget, err := s.GetCurrentBudget(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current budget: %w", err)
	}
	
	if budget == nil {
		return []BudgetCategory{}, nil
	}
	
	return budget.Categories, nil
}

// RecordSpending records a spending transaction and updates budget amounts atomically
func (s *budgetService) RecordSpending(ctx context.Context, userID string, req SpendingRequest) (*SpendingTransaction, error) {
	// Validate amount is positive
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	
	// Reject future-dated transactions
	now := time.Now().UTC()
	if req.Date.After(now) {
		return nil, fmt.Errorf("cannot record spending for future dates")
	}
	
	// Round amount to 2 decimal places
	roundedAmount := math.Round(req.Amount*100) / 100
	
	// Begin database transaction for atomicity
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	// Ensure rollback on error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	
	// Get category to verify it exists and get budget ID
	category, err := tx.GetCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}
	
	// Create spending transaction
	spending := &SpendingTransaction{
		ID:          uuid.New().String(),
		UserID:      userID,
		BudgetID:    category.BudgetID,
		CategoryID:  req.CategoryID,
		Amount:      roundedAmount,
		Description: req.Description,
		Merchant:    req.Merchant,
		Date:        req.Date,
		CreatedAt:   time.Now().UTC(),
	}
	
	// Insert spending transaction
	if err = tx.CreateSpendingTransaction(ctx, spending); err != nil {
		return nil, fmt.Errorf("failed to create spending transaction: %w", err)
	}
	
	// Update category spent amount
	if err = tx.UpdateCategorySpentAmount(ctx, req.CategoryID, roundedAmount); err != nil {
		return nil, fmt.Errorf("failed to update category spent amount: %w", err)
	}
	
	// Update budget total spent
	if err = tx.UpdateBudgetTotalSpent(ctx, category.BudgetID, roundedAmount); err != nil {
		return nil, fmt.Errorf("failed to update budget total spent: %w", err)
	}
	
	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return spending, nil
}

// GetSpendingSummary returns spending summary for a specific month
func (s *budgetService) GetSpendingSummary(ctx context.Context, userID string, month time.Time) (*SpendingSummary, error) {
	// Normalize to first day of month
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	
	budget, err := s.repo.GetBudgetByUserAndMonth(ctx, userID, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}
	
	if budget == nil {
		return &SpendingSummary{
			Month:         month,
			TotalSpent:    0,
			TotalBudget:   0,
			Remaining:     0,
			CategoryCount: 0,
		}, nil
	}
	
	categories, err := s.repo.GetCategoriesByBudgetID(ctx, budget.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	
	summary := &SpendingSummary{
		Month:         month,
		TotalSpent:    budget.TotalSpent,
		TotalBudget:   budget.TotalBudget,
		Remaining:     budget.TotalBudget - budget.TotalSpent,
		CategoryCount: len(categories),
	}
	
	return summary, nil
}

// CheckBudgetAlerts generates alerts for categories exceeding thresholds
func (s *budgetService) CheckBudgetAlerts(ctx context.Context, userID string) ([]BudgetAlert, error) {
	budget, err := s.GetCurrentBudget(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current budget: %w", err)
	}
	
	if budget == nil {
		return []BudgetAlert{}, nil
	}
	
	alerts := []BudgetAlert{}
	
	// Check each category against thresholds
	for _, category := range budget.Categories {
		// Skip categories with zero allocated amount
		if category.AllocatedAmount == 0 {
			continue
		}
		
		percentageUsed := (category.SpentAmount / category.AllocatedAmount) * 100
		
		// Critical alert: 100% or more spent
		if percentageUsed >= 100 {
			alerts = append(alerts, BudgetAlert{
				CategoryName:   category.Name,
				PercentageUsed: percentageUsed,
				AlertType:      "critical",
				Message:        fmt.Sprintf("You've exceeded your %s budget by %.1f%%", category.Name, percentageUsed-100),
			})
		} else if percentageUsed >= 80 {
			// Warning alert: 80-99% spent
			alerts = append(alerts, BudgetAlert{
				CategoryName:   category.Name,
				PercentageUsed: percentageUsed,
				AlertType:      "warning",
				Message:        fmt.Sprintf("You've used %.1f%% of your %s budget", percentageUsed, category.Name),
			})
		}
	}
	
	// Sort alerts: critical before warning, then by percentage descending
	sort.Slice(alerts, func(i, j int) bool {
		if alerts[i].AlertType == "critical" && alerts[j].AlertType == "warning" {
			return true
		}
		if alerts[i].AlertType == "warning" && alerts[j].AlertType == "critical" {
			return false
		}
		return alerts[i].PercentageUsed > alerts[j].PercentageUsed
	})
	
	return alerts, nil
}
