package budget

import "time"

// Budget represents a monthly budget with category allocations
type Budget struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	Month           time.Time        `json:"month"`
	TotalBudget     float64          `json:"total_budget"`
	Categories      []BudgetCategory `json:"categories"`
	TotalSpent      float64          `json:"total_spent"`
	RemainingBudget float64          `json:"remaining_budget"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// BudgetCategory represents a spending category within a budget
type BudgetCategory struct {
	ID              string  `json:"id"`
	BudgetID        string  `json:"budget_id"`
	Name            string  `json:"name"`
	AllocatedAmount float64 `json:"allocated_amount"`
	SpentAmount     float64 `json:"spent_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	Color           string  `json:"color"`
}

// SpendingTransaction represents a spending record
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

// BudgetAlert represents a budget threshold alert
type BudgetAlert struct {
	CategoryName   string  `json:"category_name"`
	PercentageUsed float64 `json:"percentage_used"`
	AlertType      string  `json:"alert_type"` // "warning" | "critical"
	Message        string  `json:"message"`
}

// CreateBudgetRequest represents the request to create a budget
type CreateBudgetRequest struct {
	Month      time.Time                  `json:"month" validate:"required"`
	Categories []CreateBudgetCategoryRequest `json:"categories" validate:"required,min=1"`
}

// CreateBudgetCategoryRequest represents a category in budget creation
type CreateBudgetCategoryRequest struct {
	Name            string  `json:"name" validate:"required"`
	AllocatedAmount float64 `json:"allocated_amount" validate:"required,gte=0"`
	Color           string  `json:"color"`
}

// UpdateBudgetRequest represents the request to update a budget
type UpdateBudgetRequest struct {
	Categories []UpdateBudgetCategoryRequest `json:"categories" validate:"required,min=1"`
}

// UpdateBudgetCategoryRequest represents a category update
type UpdateBudgetCategoryRequest struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	AllocatedAmount float64 `json:"allocated_amount" validate:"gte=0"`
	Color           string  `json:"color"`
}

// SpendingRequest represents the request to record spending
type SpendingRequest struct {
	CategoryID  string    `json:"category_id" validate:"required"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Description string    `json:"description"`
	Merchant    string    `json:"merchant"`
	Date        time.Time `json:"date" validate:"required"`
}

// SpendingSummary provides spending overview for a month
type SpendingSummary struct {
	Month        time.Time `json:"month"`
	TotalSpent   float64   `json:"total_spent"`
	TotalBudget  float64   `json:"total_budget"`
	Remaining    float64   `json:"remaining"`
	CategoryCount int      `json:"category_count"`
}
