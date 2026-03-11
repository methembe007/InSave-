package analytics

import "time"

// TimePeriod represents a time range for analysis
type TimePeriod struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// SpendingAnalysis represents spending analysis for a period
type SpendingAnalysis struct {
	Period               TimePeriod          `json:"period"`
	TotalSpending        float64             `json:"total_spending"`
	CategoryBreakdown    []CategorySpending  `json:"category_breakdown"`
	TopMerchants         []MerchantSpending  `json:"top_merchants"`
	DailyAverage         float64             `json:"daily_average"`
	ComparisonToPrevious float64             `json:"comparison_to_previous"`
	Trends               []SpendingTrend     `json:"trends"`
}

// CategorySpending represents spending in a category
type CategorySpending struct {
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
	Count        int     `json:"count"`
}

// MerchantSpending represents spending at a merchant
type MerchantSpending struct {
	MerchantName string  `json:"merchant_name"`
	Amount       float64 `json:"amount"`
	Count        int     `json:"count"`
}

// SpendingTrend represents a spending trend
type SpendingTrend struct {
	Category   string  `json:"category"`
	Trend      string  `json:"trend"` // "increasing" | "decreasing" | "stable"
	ChangeRate float64 `json:"change_rate"`
}

// SavingsPattern represents detected savings patterns
type SavingsPattern struct {
	PatternType   string   `json:"pattern_type"` // "consistent" | "irregular" | "improving"
	AverageAmount float64  `json:"average_amount"`
	Frequency     string   `json:"frequency"`
	BestDayOfWeek string   `json:"best_day_of_week"`
	Insights      []string `json:"insights"`
}

// Recommendation represents an AI-assisted recommendation
type Recommendation struct {
	ID               string   `json:"id"`
	Type             string   `json:"type"` // "savings" | "budget" | "spending"
	Priority         string   `json:"priority"` // "high" | "medium" | "low"
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	ActionItems      []string `json:"action_items"`
	PotentialSavings float64  `json:"potential_savings,omitempty"`
}

// FinancialHealthScore represents the user's financial health metrics
type FinancialHealthScore struct {
	OverallScore     int      `json:"overall_score"` // 0-100
	SavingsScore     int      `json:"savings_score"`
	BudgetScore      int      `json:"budget_score"`
	ConsistencyScore int      `json:"consistency_score"`
	Insights         []string `json:"insights"`
	ImprovementAreas []string `json:"improvement_areas"`
}

// MonthlyReport represents a comprehensive monthly financial report
type MonthlyReport struct {
	Month            time.Time            `json:"month"`
	TotalSaved       float64              `json:"total_saved"`
	TotalSpent       float64              `json:"total_spent"`
	NetSavings       float64              `json:"net_savings"`
	SpendingAnalysis SpendingAnalysis     `json:"spending_analysis"`
	SavingsPatterns  []SavingsPattern     `json:"savings_patterns"`
	FinancialHealth  FinancialHealthScore `json:"financial_health"`
	Recommendations  []Recommendation     `json:"recommendations"`
}
