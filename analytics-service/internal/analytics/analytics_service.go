package analytics

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

type analyticsService struct {
	repo  Repository
	cache Cache
}

// Cache defines the caching interface for analytics
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
}

// NewService creates a new analytics service
func NewService(repo Repository, cache Cache) Service {
	return &analyticsService{
		repo:  repo,
		cache: cache,
	}
}

// GetSpendingAnalysis analyzes spending patterns for a given period
func (s *analyticsService) GetSpendingAnalysis(ctx context.Context, userID string, period TimePeriod) (*SpendingAnalysis, error) {
	// Get spending transactions for the period
	transactions, err := s.repo.GetSpendingTransactions(ctx, userID, period.Start, period.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending transactions: %w", err)
	}
	
	// Calculate total spending
	totalSpending := 0.0
	categoryMap := make(map[string]*CategorySpending)
	merchantMap := make(map[string]*MerchantSpending)
	
	for _, tx := range transactions {
		totalSpending += tx.Amount
		
		// Aggregate by category
		if cat, exists := categoryMap[tx.CategoryID]; exists {
			cat.Amount += tx.Amount
			cat.Count++
		} else {
			categoryMap[tx.CategoryID] = &CategorySpending{
				CategoryName: tx.CategoryID, // Will be replaced with actual name
				Amount:       tx.Amount,
				Count:        1,
			}
		}
		
		// Aggregate by merchant
		if tx.Merchant != "" {
			if merch, exists := merchantMap[tx.Merchant]; exists {
				merch.Amount += tx.Amount
				merch.Count++
			} else {
				merchantMap[tx.Merchant] = &MerchantSpending{
					MerchantName: tx.Merchant,
					Amount:       tx.Amount,
					Count:        1,
				}
			}
		}
	}
	
	// Convert maps to slices and calculate percentages
	categoryBreakdown := make([]CategorySpending, 0, len(categoryMap))
	for _, cat := range categoryMap {
		cat.Percentage = (cat.Amount / totalSpending) * 100
		categoryBreakdown = append(categoryBreakdown, *cat)
	}
	
	// Sort categories by amount descending
	sort.Slice(categoryBreakdown, func(i, j int) bool {
		return categoryBreakdown[i].Amount > categoryBreakdown[j].Amount
	})
	
	// Get top 5 merchants
	topMerchants := make([]MerchantSpending, 0, len(merchantMap))
	for _, merch := range merchantMap {
		topMerchants = append(topMerchants, *merch)
	}
	sort.Slice(topMerchants, func(i, j int) bool {
		return topMerchants[i].Amount > topMerchants[j].Amount
	})
	if len(topMerchants) > 5 {
		topMerchants = topMerchants[:5]
	}
	
	// Calculate daily average
	days := period.End.Sub(period.Start).Hours() / 24
	if days == 0 {
		days = 1
	}
	dailyAverage := totalSpending / days
	
	// Calculate comparison to previous period
	previousPeriod := TimePeriod{
		Start: period.Start.AddDate(0, 0, -int(days)),
		End:   period.Start,
	}
	previousTransactions, err := s.repo.GetSpendingTransactions(ctx, userID, previousPeriod.Start, previousPeriod.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous period transactions: %w", err)
	}
	
	previousTotal := 0.0
	for _, tx := range previousTransactions {
		previousTotal += tx.Amount
	}
	
	comparisonToPrevious := 0.0
	if previousTotal > 0 {
		comparisonToPrevious = ((totalSpending - previousTotal) / previousTotal) * 100
	}
	
	analysis := &SpendingAnalysis{
		Period:               period,
		TotalSpending:        math.Round(totalSpending*100) / 100,
		CategoryBreakdown:    categoryBreakdown,
		TopMerchants:         topMerchants,
		DailyAverage:         math.Round(dailyAverage*100) / 100,
		ComparisonToPrevious: math.Round(comparisonToPrevious*100) / 100,
		Trends:               []SpendingTrend{}, // TODO: Implement trend detection
	}
	
	return analysis, nil
}

// GetSavingsPatterns detects and returns savings patterns
func (s *analyticsService) GetSavingsPatterns(ctx context.Context, userID string) ([]SavingsPattern, error) {
	// Get last 90 days of savings transactions
	end := time.Now().UTC()
	start := end.AddDate(0, 0, -90)
	
	transactions, err := s.repo.GetSavingsTransactions(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get savings transactions: %w", err)
	}
	
	if len(transactions) == 0 {
		return []SavingsPattern{}, nil
	}
	
	// Calculate average amount
	totalAmount := 0.0
	dayOfWeekCounts := make(map[time.Weekday]int)
	dayOfWeekAmounts := make(map[time.Weekday]float64)
	
	for _, tx := range transactions {
		totalAmount += tx.Amount
		weekday := tx.CreatedAt.Weekday()
		dayOfWeekCounts[weekday]++
		dayOfWeekAmounts[weekday] += tx.Amount
	}
	
	averageAmount := totalAmount / float64(len(transactions))
	
	// Determine best day of week
	bestDay := time.Sunday
	maxCount := 0
	for day, count := range dayOfWeekCounts {
		if count > maxCount {
			maxCount = count
			bestDay = day
		}
	}
	
	// Determine frequency
	frequency := "irregular"
	avgDaysPerTransaction := 90.0 / float64(len(transactions))
	if avgDaysPerTransaction <= 1.5 {
		frequency = "daily"
	} else if avgDaysPerTransaction <= 7 {
		frequency = "weekly"
	} else if avgDaysPerTransaction <= 14 {
		frequency = "bi-weekly"
	} else if avgDaysPerTransaction <= 31 {
		frequency = "monthly"
	}
	
	// Determine pattern type
	patternType := "irregular"
	if len(transactions) >= 30 {
		// Check consistency by comparing first half to second half
		midpoint := len(transactions) / 2
		firstHalfAvg := 0.0
		secondHalfAvg := 0.0
		
		for i := 0; i < midpoint; i++ {
			firstHalfAvg += transactions[i].Amount
		}
		firstHalfAvg /= float64(midpoint)
		
		for i := midpoint; i < len(transactions); i++ {
			secondHalfAvg += transactions[i].Amount
		}
		secondHalfAvg /= float64(len(transactions) - midpoint)
		
		// If amounts are similar, it's consistent
		if math.Abs(firstHalfAvg-secondHalfAvg)/firstHalfAvg < 0.2 {
			patternType = "consistent"
		} else if secondHalfAvg > firstHalfAvg {
			patternType = "improving"
		}
	}
	
	// Generate insights
	insights := []string{}
	if patternType == "consistent" {
		insights = append(insights, fmt.Sprintf("You're saving consistently with an average of $%.2f per transaction", averageAmount))
	} else if patternType == "improving" {
		insights = append(insights, "Your savings amounts are increasing over time - great progress!")
	} else {
		insights = append(insights, "Your savings pattern is irregular. Try setting a regular savings schedule.")
	}
	
	if frequency == "daily" {
		insights = append(insights, "You're saving almost every day - excellent habit!")
	}
	
	insights = append(insights, fmt.Sprintf("You save most frequently on %ss", bestDay.String()))
	
	pattern := SavingsPattern{
		PatternType:   patternType,
		AverageAmount: math.Round(averageAmount*100) / 100,
		Frequency:     frequency,
		BestDayOfWeek: bestDay.String(),
		Insights:      insights,
	}
	
	return []SavingsPattern{pattern}, nil
}

// GetRecommendations generates AI-assisted recommendations
func (s *analyticsService) GetRecommendations(ctx context.Context, userID string) ([]Recommendation, error) {
	recommendations := []Recommendation{}
	
	// Get spending analysis for last 30 days
	end := time.Now().UTC()
	start := end.AddDate(0, 0, -30)
	period := TimePeriod{Start: start, End: end}
	
	spendingAnalysis, err := s.GetSpendingAnalysis(ctx, userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending analysis: %w", err)
	}
	
	// Get savings patterns
	savingsPatterns, err := s.GetSavingsPatterns(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get savings patterns: %w", err)
	}
	
	// Recommendation 1: High spending category
	if len(spendingAnalysis.CategoryBreakdown) > 0 {
		topCategory := spendingAnalysis.CategoryBreakdown[0]
		if topCategory.Percentage > 40 {
			recommendations = append(recommendations, Recommendation{
				ID:          uuid.New().String(),
				Type:        "spending",
				Priority:    "high",
				Title:       fmt.Sprintf("Reduce %s spending", topCategory.CategoryName),
				Description: fmt.Sprintf("You're spending %.1f%% of your budget on %s. Consider reducing this category.", topCategory.Percentage, topCategory.CategoryName),
				ActionItems: []string{
					fmt.Sprintf("Set a lower budget for %s", topCategory.CategoryName),
					"Track your spending more carefully in this category",
					"Look for cheaper alternatives",
				},
				PotentialSavings: topCategory.Amount * 0.2, // 20% reduction potential
			})
		}
	}
	
	// Recommendation 2: Savings pattern improvement
	if len(savingsPatterns) > 0 {
		pattern := savingsPatterns[0]
		if pattern.PatternType == "irregular" {
			recommendations = append(recommendations, Recommendation{
				ID:          uuid.New().String(),
				Type:        "savings",
				Priority:    "medium",
				Title:       "Establish a regular savings schedule",
				Description: "Your savings pattern is irregular. Setting up automatic savings can help build consistency.",
				ActionItems: []string{
					"Set up automatic transfers on payday",
					"Start with a small, manageable amount",
					"Gradually increase your savings rate",
				},
				PotentialSavings: pattern.AverageAmount * 4, // Potential monthly increase
			})
		}
	}
	
	// Recommendation 3: Spending increase alert
	if spendingAnalysis.ComparisonToPrevious > 20 {
		recommendations = append(recommendations, Recommendation{
			ID:          uuid.New().String(),
			Type:        "budget",
			Priority:    "high",
			Title:       "Spending increased significantly",
			Description: fmt.Sprintf("Your spending increased by %.1f%% compared to the previous period.", spendingAnalysis.ComparisonToPrevious),
			ActionItems: []string{
				"Review your recent transactions",
				"Identify unnecessary expenses",
				"Set stricter budget limits",
			},
			PotentialSavings: spendingAnalysis.TotalSpending * 0.15,
		})
	}
	
	// Sort by priority (high > medium > low)
	sort.Slice(recommendations, func(i, j int) bool {
		priorityOrder := map[string]int{"high": 3, "medium": 2, "low": 1}
		return priorityOrder[recommendations[i].Priority] > priorityOrder[recommendations[j].Priority]
	})
	
	return recommendations, nil
}

// GenerateMonthlyReport creates a comprehensive monthly financial report
func (s *analyticsService) GenerateMonthlyReport(ctx context.Context, userID string, month time.Time) (*MonthlyReport, error) {
	// Normalize to first day of month
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	
	// Get period for the month
	start := month
	end := month.AddDate(0, 1, 0).Add(-time.Second)
	period := TimePeriod{Start: start, End: end}
	
	// Get spending analysis
	spendingAnalysis, err := s.GetSpendingAnalysis(ctx, userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending analysis: %w", err)
	}
	
	// Get savings for the month
	savingsTransactions, err := s.repo.GetSavingsTransactions(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get savings transactions: %w", err)
	}
	
	totalSaved := 0.0
	for _, tx := range savingsTransactions {
		totalSaved += tx.Amount
	}
	
	// Get savings patterns
	savingsPatterns, err := s.GetSavingsPatterns(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get savings patterns: %w", err)
	}
	
	// Get financial health
	financialHealth, err := s.GetFinancialHealth(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get financial health: %w", err)
	}
	
	// Get recommendations
	recommendations, err := s.GetRecommendations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}
	
	report := &MonthlyReport{
		Month:            month,
		TotalSaved:       math.Round(totalSaved*100) / 100,
		TotalSpent:       spendingAnalysis.TotalSpending,
		NetSavings:       math.Round((totalSaved-spendingAnalysis.TotalSpending)*100) / 100,
		SpendingAnalysis: *spendingAnalysis,
		SavingsPatterns:  savingsPatterns,
		FinancialHealth:  *financialHealth,
		Recommendations:  recommendations,
	}
	
	return report, nil
}

// GetFinancialHealth calculates the user's financial health score
func (s *analyticsService) GetFinancialHealth(ctx context.Context, userID string) (*FinancialHealthScore, error) {
	// Check cache first (1 hour TTL)
	cacheKey := fmt.Sprintf("financial_health:%s", userID)
	if cached, found := s.cache.Get(cacheKey); found {
		if score, ok := cached.(*FinancialHealthScore); ok {
			return score, nil
		}
	}
	
	// Get data for last 30 days
	end := time.Now().UTC()
	start := end.AddDate(0, 0, -30)
	
	// Check if user has at least 30 days of history
	savingsTransactions, err := s.repo.GetSavingsTransactions(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get savings transactions: %w", err)
	}
	
	// Get user's first transaction to check history length
	allTimeSavings, err := s.repo.GetSavingsTransactions(ctx, userID, time.Time{}, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get all savings transactions: %w", err)
	}
	
	if len(allTimeSavings) > 0 {
		firstTransaction := allTimeSavings[len(allTimeSavings)-1]
		daysSinceFirst := int(end.Sub(firstTransaction.CreatedAt).Hours() / 24)
		if daysSinceFirst < 30 {
			return nil, fmt.Errorf("insufficient data: user has only %d days of history, need at least 30 days", daysSinceFirst)
		}
	} else {
		return nil, fmt.Errorf("insufficient data: no transaction history found")
	}
	
	// Calculate Savings Score (40% weight)
	savingsScore := s.calculateSavingsScore(ctx, userID, savingsTransactions)
	
	// Calculate Budget Score (30% weight)
	budgetScore, err := s.calculateBudgetScore(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate budget score: %w", err)
	}
	
	// Calculate Consistency Score (30% weight)
	consistencyScore, err := s.calculateConsistencyScore(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate consistency score: %w", err)
	}
	
	// Calculate overall score as weighted average
	overallScore := int(math.Round(
		float64(savingsScore)*0.4 +
			float64(budgetScore)*0.3 +
			float64(consistencyScore)*0.3,
	))
	
	// Ensure scores are within bounds
	overallScore = clamp(overallScore, 0, 100)
	savingsScore = clamp(savingsScore, 0, 100)
	budgetScore = clamp(budgetScore, 0, 100)
	consistencyScore = clamp(consistencyScore, 0, 100)
	
	// Generate insights
	insights := []string{}
	improvementAreas := []string{}
	
	if savingsScore >= 80 {
		insights = append(insights, "Excellent savings habits!")
	} else if savingsScore < 50 {
		improvementAreas = append(improvementAreas, "Increase your savings frequency and amounts")
	}
	
	if budgetScore >= 80 {
		insights = append(insights, "Great budget adherence!")
	} else if budgetScore < 50 {
		improvementAreas = append(improvementAreas, "Focus on staying within your budget limits")
	}
	
	if consistencyScore >= 80 {
		insights = append(insights, "You're maintaining excellent consistency!")
	} else if consistencyScore < 50 {
		improvementAreas = append(improvementAreas, "Build a more consistent savings routine")
	}
	
	if overallScore >= 80 {
		insights = append(insights, "Your financial health is excellent!")
	} else if overallScore >= 60 {
		insights = append(insights, "Your financial health is good, with room for improvement")
	} else {
		insights = append(insights, "Focus on building better financial habits")
	}
	
	score := &FinancialHealthScore{
		OverallScore:     overallScore,
		SavingsScore:     savingsScore,
		BudgetScore:      budgetScore,
		ConsistencyScore: consistencyScore,
		Insights:         insights,
		ImprovementAreas: improvementAreas,
	}
	
	// Cache the result for 1 hour
	s.cache.Set(cacheKey, score, time.Hour)
	
	return score, nil
}

// calculateSavingsScore calculates the savings component score (0-100)
func (s *analyticsService) calculateSavingsScore(ctx context.Context, userID string, transactions []SavingsTransaction) int {
	if len(transactions) == 0 {
		return 0
	}
	
	// Calculate frequency score (0-50 points)
	// 30 transactions in 30 days = 50 points
	frequencyScore := float64(len(transactions)) / 30.0 * 50.0
	if frequencyScore > 50 {
		frequencyScore = 50
	}
	
	// Calculate amount score (0-50 points)
	totalAmount := 0.0
	for _, tx := range transactions {
		totalAmount += tx.Amount
	}
	averageAmount := totalAmount / float64(len(transactions))
	
	// $10 average = 25 points, $20+ average = 50 points
	amountScore := (averageAmount / 20.0) * 50.0
	if amountScore > 50 {
		amountScore = 50
	}
	
	return int(math.Round(frequencyScore + amountScore))
}

// calculateBudgetScore calculates the budget adherence score (0-100)
func (s *analyticsService) calculateBudgetScore(ctx context.Context, userID string, start, end time.Time) (int, error) {
	// Get current month's budget
	month := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	budget, err := s.repo.GetBudgetForMonth(ctx, userID, month)
	if err != nil {
		return 50, nil // Default score if no budget
	}
	
	if budget == nil || budget.TotalBudget == 0 {
		return 50, nil // Default score if no budget set
	}
	
	// Calculate percentage of budget used
	percentageUsed := (budget.TotalSpent / budget.TotalBudget) * 100
	
	// Score based on budget adherence
	// 0-80% used = 100 points
	// 80-100% used = 80-50 points
	// 100-120% used = 50-20 points
	// 120%+ used = 0-20 points
	
	var score float64
	if percentageUsed <= 80 {
		score = 100
	} else if percentageUsed <= 100 {
		score = 100 - ((percentageUsed - 80) / 20 * 20)
	} else if percentageUsed <= 120 {
		score = 50 - ((percentageUsed - 100) / 20 * 30)
	} else {
		score = 20 - math.Min(20, (percentageUsed-120)/10)
	}
	
	return int(math.Round(score)), nil
}

// calculateConsistencyScore calculates the consistency score (0-100)
func (s *analyticsService) calculateConsistencyScore(ctx context.Context, userID string) (int, error) {
	// Get streak information
	currentStreak, longestStreak, err := s.repo.GetUserStreak(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get streak: %w", err)
	}
	
	// Streak score (0-60 points)
	// 30+ day streak = 60 points
	streakScore := float64(currentStreak) / 30.0 * 60.0
	if streakScore > 60 {
		streakScore = 60
	}
	
	// Regularity score (0-40 points)
	// Based on ratio of current to longest streak
	regularityScore := 0.0
	if longestStreak > 0 {
		regularityScore = (float64(currentStreak) / float64(longestStreak)) * 40.0
	}
	
	return int(math.Round(streakScore + regularityScore)), nil
}

// clamp ensures a value is within the specified range
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
