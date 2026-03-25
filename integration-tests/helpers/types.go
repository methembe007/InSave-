package helpers

import "time"

// Auth types
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	User         UserSummary `json:"user"`
}

type UserSummary struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// User types
type UserProfile struct {
	ID              string                 `json:"id"`
	Email           string                 `json:"email"`
	FirstName       string                 `json:"first_name"`
	LastName        string                 `json:"last_name"`
	DateOfBirth     string                 `json:"date_of_birth"`
	ProfileImageURL string                 `json:"profile_image_url"`
	Preferences     map[string]interface{} `json:"preferences"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	DateOfBirth     string `json:"date_of_birth,omitempty"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
}

type UserPreferences struct {
	Currency             string `json:"currency"`
	NotificationsEnabled bool   `json:"notifications_enabled"`
	EmailNotifications   bool   `json:"email_notifications"`
	PushNotifications    bool   `json:"push_notifications"`
	SavingsReminders     bool   `json:"savings_reminders"`
	ReminderTime         string `json:"reminder_time"`
	Theme                string `json:"theme"`
}

// Savings types
type CreateSavingsRequest struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

type SavingsTransaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

type SavingsSummary struct {
	TotalSaved     float64   `json:"total_saved"`
	CurrentStreak  int       `json:"current_streak"`
	LongestStreak  int       `json:"longest_streak"`
	LastSavingDate time.Time `json:"last_saving_date"`
	MonthlyAverage float64   `json:"monthly_average"`
	ThisMonthSaved float64   `json:"this_month_saved"`
}

type SavingsStreak struct {
	CurrentStreak int       `json:"current_streak"`
	LongestStreak int       `json:"longest_streak"`
	LastSaveDate  time.Time `json:"last_save_date"`
}

// Budget types
type CreateBudgetRequest struct {
	Month       string           `json:"month"`
	TotalBudget float64          `json:"total_budget"`
	Categories  []BudgetCategory `json:"categories"`
}

type Budget struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	Month           string           `json:"month"`
	TotalBudget     float64          `json:"total_budget"`
	TotalSpent      float64          `json:"total_spent"`
	RemainingBudget float64          `json:"remaining_budget"`
	Categories      []BudgetCategory `json:"categories"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type BudgetCategory struct {
	ID              string  `json:"id"`
	BudgetID        string  `json:"budget_id,omitempty"`
	Name            string  `json:"name"`
	AllocatedAmount float64 `json:"allocated_amount"`
	SpentAmount     float64 `json:"spent_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	Color           string  `json:"color"`
}

type SpendingRequest struct {
	CategoryID  string  `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Merchant    string  `json:"merchant"`
	Date        string  `json:"date"`
}

type SpendingTransaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	BudgetID        string    `json:"budget_id"`
	CategoryID      string    `json:"category_id"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	Merchant        string    `json:"merchant"`
	TransactionDate string    `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
}

type BudgetAlert struct {
	CategoryName   string  `json:"category_name"`
	PercentageUsed float64 `json:"percentage_used"`
	AlertType      string  `json:"alert_type"` // "warning" | "critical"
	Message        string  `json:"message"`
}

// Goal types
type CreateGoalRequest struct {
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	TargetAmount float64     `json:"target_amount"`
	Currency     string      `json:"currency"`
	TargetDate   string      `json:"target_date"`
	Milestones   []Milestone `json:"milestones,omitempty"`
}

type Goal struct {
	ID              string      `json:"id"`
	UserID          string      `json:"user_id"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	TargetAmount    float64     `json:"target_amount"`
	CurrentAmount   float64     `json:"current_amount"`
	Currency        string      `json:"currency"`
	TargetDate      string      `json:"target_date"`
	Status          string      `json:"status"`
	ProgressPercent float64     `json:"progress_percent"`
	Milestones      []Milestone `json:"milestones,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type Milestone struct {
	ID          string     `json:"id"`
	GoalID      string     `json:"goal_id,omitempty"`
	Title       string     `json:"title"`
	Amount      float64    `json:"amount"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Order       int        `json:"order"`
}

type ContributionRequest struct {
	Amount float64 `json:"amount"`
}

// Notification types
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
