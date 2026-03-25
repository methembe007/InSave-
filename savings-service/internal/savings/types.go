package savings

import "time"

// SavingsTransaction represents a single savings deposit
type SavingsTransaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

// SavingsSummary provides an overview of user's savings
type SavingsSummary struct {
	TotalSaved     float64   `json:"total_saved"`
	CurrentStreak  int       `json:"current_streak"`
	LongestStreak  int       `json:"longest_streak"`
	LastSavingDate time.Time `json:"last_saving_date"`
	MonthlyAverage float64   `json:"monthly_average"`
	ThisMonthSaved float64   `json:"this_month_saved"`
}

// SavingsStreak represents the user's savings streak information
type SavingsStreak struct {
	CurrentStreak int       `json:"current_streak"`
	LongestStreak int       `json:"longest_streak"`
	LastSaveDate  time.Time `json:"last_save_date"`
}

// MonthlyStats provides statistics for a specific month
type MonthlyStats struct {
	Month         time.Time `json:"month"`
	TotalSaved    float64   `json:"total_saved"`
	Count         int       `json:"count"`
	AverageAmount float64   `json:"average_amount"`
}

// CreateTransactionRequest represents the request to create a savings transaction
type CreateTransactionRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Currency    string  `json:"currency" validate:"required,len=3"`
	Description string  `json:"description" validate:"max=500"`
	Category    string  `json:"category" validate:"required,max=50"`
}

// HistoryParams represents parameters for fetching savings history
type HistoryParams struct {
	Limit  int
	Offset int
}
