package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// CacheConfig holds caching configuration
type CacheConfig struct {
	// Session caching
	SessionTTL time.Duration // 15 minutes

	// Financial health score caching
	FinancialHealthTTL time.Duration // 1 hour

	// Education content caching
	EducationContentTTL time.Duration // 24 hours

	// User profile caching
	UserProfileTTL time.Duration // 30 minutes

	// Budget caching
	BudgetTTL time.Duration // 15 minutes

	// Savings summary caching
	SavingsSummaryTTL time.Duration // 5 minutes

	// Goals caching
	GoalsTTL time.Duration // 10 minutes
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		SessionTTL:          15 * time.Minute,
		FinancialHealthTTL:  1 * time.Hour,
		EducationContentTTL: 24 * time.Hour,
		UserProfileTTL:      30 * time.Minute,
		BudgetTTL:           15 * time.Minute,
		SavingsSummaryTTL:   5 * time.Minute,
		GoalsTTL:            10 * time.Minute,
	}
}

// CacheManager provides high-level caching operations
type CacheManager struct {
	cache  Cache
	config *CacheConfig
}

// NewCacheManager creates a new cache manager
func NewCacheManager(cache Cache, config *CacheConfig) *CacheManager {
	if config == nil {
		config = DefaultCacheConfig()
	}
	return &CacheManager{
		cache:  cache,
		config: config,
	}
}

// Session Cache Operations

// GetSession retrieves a session from cache
func (m *CacheManager) GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	key := CacheKey("session", sessionID)
	var session map[string]interface{}
	if err := m.cache.Get(ctx, key, &session); err != nil {
		return nil, err
	}
	return session, nil
}

// SetSession stores a session in cache
func (m *CacheManager) SetSession(ctx context.Context, sessionID string, session map[string]interface{}) error {
	key := CacheKey("session", sessionID)
	return m.cache.Set(ctx, key, session, m.config.SessionTTL)
}

// DeleteSession removes a session from cache
func (m *CacheManager) DeleteSession(ctx context.Context, sessionID string) error {
	key := CacheKey("session", sessionID)
	return m.cache.Delete(ctx, key)
}

// Financial Health Score Cache Operations

// GetFinancialHealth retrieves financial health score from cache
func (m *CacheManager) GetFinancialHealth(ctx context.Context, userID string) (interface{}, error) {
	key := UserCacheKey("financial_health", userID, "score")
	var health interface{}
	if err := m.cache.Get(ctx, key, &health); err != nil {
		return nil, err
	}
	return health, nil
}

// SetFinancialHealth stores financial health score in cache
func (m *CacheManager) SetFinancialHealth(ctx context.Context, userID string, health interface{}) error {
	key := UserCacheKey("financial_health", userID, "score")
	return m.cache.Set(ctx, key, health, m.config.FinancialHealthTTL)
}

// InvalidateFinancialHealth removes financial health score from cache
func (m *CacheManager) InvalidateFinancialHealth(ctx context.Context, userID string) error {
	key := UserCacheKey("financial_health", userID, "score")
	return m.cache.Delete(ctx, key)
}

// Education Content Cache Operations

// GetLesson retrieves a lesson from cache
func (m *CacheManager) GetLesson(ctx context.Context, lessonID string) (interface{}, error) {
	key := CacheKey("lesson", lessonID)
	var lesson interface{}
	if err := m.cache.Get(ctx, key, &lesson); err != nil {
		return nil, err
	}
	return lesson, nil
}

// SetLesson stores a lesson in cache
func (m *CacheManager) SetLesson(ctx context.Context, lessonID string, lesson interface{}) error {
	key := CacheKey("lesson", lessonID)
	return m.cache.Set(ctx, key, lesson, m.config.EducationContentTTL)
}

// GetLessonList retrieves lesson list from cache
func (m *CacheManager) GetLessonList(ctx context.Context, category string) (interface{}, error) {
	key := CacheKey("lessons", category)
	var lessons interface{}
	if err := m.cache.Get(ctx, key, &lessons); err != nil {
		return nil, err
	}
	return lessons, nil
}

// SetLessonList stores lesson list in cache
func (m *CacheManager) SetLessonList(ctx context.Context, category string, lessons interface{}) error {
	key := CacheKey("lessons", category)
	return m.cache.Set(ctx, key, lessons, m.config.EducationContentTTL)
}

// User Profile Cache Operations

// GetUserProfile retrieves user profile from cache
func (m *CacheManager) GetUserProfile(ctx context.Context, userID string) (interface{}, error) {
	key := UserCacheKey("profile", userID, "data")
	var profile interface{}
	if err := m.cache.Get(ctx, key, &profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// SetUserProfile stores user profile in cache
func (m *CacheManager) SetUserProfile(ctx context.Context, userID string, profile interface{}) error {
	key := UserCacheKey("profile", userID, "data")
	return m.cache.Set(ctx, key, profile, m.config.UserProfileTTL)
}

// InvalidateUserProfile removes user profile from cache
func (m *CacheManager) InvalidateUserProfile(ctx context.Context, userID string) error {
	key := UserCacheKey("profile", userID, "data")
	return m.cache.Delete(ctx, key)
}

// Budget Cache Operations

// GetBudget retrieves budget from cache
func (m *CacheManager) GetBudget(ctx context.Context, userID, month string) (interface{}, error) {
	key := UserCacheKey("budget", userID, month)
	var budget interface{}
	if err := m.cache.Get(ctx, key, &budget); err != nil {
		return nil, err
	}
	return budget, nil
}

// SetBudget stores budget in cache
func (m *CacheManager) SetBudget(ctx context.Context, userID, month string, budget interface{}) error {
	key := UserCacheKey("budget", userID, month)
	return m.cache.Set(ctx, key, budget, m.config.BudgetTTL)
}

// InvalidateBudget removes budget from cache
func (m *CacheManager) InvalidateBudget(ctx context.Context, userID, month string) error {
	key := UserCacheKey("budget", userID, month)
	return m.cache.Delete(ctx, key)
}

// InvalidateUserBudgets removes all budgets for a user
func (m *CacheManager) InvalidateUserBudgets(ctx context.Context, userID string) error {
	pattern := UserCacheKey("budget", userID, "*")
	return m.cache.DeletePattern(ctx, pattern)
}

// Savings Summary Cache Operations

// GetSavingsSummary retrieves savings summary from cache
func (m *CacheManager) GetSavingsSummary(ctx context.Context, userID string) (interface{}, error) {
	key := UserCacheKey("savings", userID, "summary")
	var summary interface{}
	if err := m.cache.Get(ctx, key, &summary); err != nil {
		return nil, err
	}
	return summary, nil
}

// SetSavingsSummary stores savings summary in cache
func (m *CacheManager) SetSavingsSummary(ctx context.Context, userID string, summary interface{}) error {
	key := UserCacheKey("savings", userID, "summary")
	return m.cache.Set(ctx, key, summary, m.config.SavingsSummaryTTL)
}

// InvalidateSavingsSummary removes savings summary from cache
func (m *CacheManager) InvalidateSavingsSummary(ctx context.Context, userID string) error {
	key := UserCacheKey("savings", userID, "summary")
	return m.cache.Delete(ctx, key)
}

// Goals Cache Operations

// GetGoals retrieves goals from cache
func (m *CacheManager) GetGoals(ctx context.Context, userID string) (interface{}, error) {
	key := UserCacheKey("goals", userID, "list")
	var goals interface{}
	if err := m.cache.Get(ctx, key, &goals); err != nil {
		return nil, err
	}
	return goals, nil
}

// SetGoals stores goals in cache
func (m *CacheManager) SetGoals(ctx context.Context, userID string, goals interface{}) error {
	key := UserCacheKey("goals", userID, "list")
	return m.cache.Set(ctx, key, goals, m.config.GoalsTTL)
}

// InvalidateGoals removes goals from cache
func (m *CacheManager) InvalidateGoals(ctx context.Context, userID string) error {
	key := UserCacheKey("goals", userID, "list")
	return m.cache.Delete(ctx, key)
}

// Utility Functions

// GenerateCacheKey generates a cache key from multiple parts
func GenerateCacheKey(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte(part))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// InvalidateUserCache removes all cached data for a user
func (m *CacheManager) InvalidateUserCache(ctx context.Context, userID string) error {
	patterns := []string{
		UserCacheKey("*", userID, "*"),
	}

	for _, pattern := range patterns {
		if err := m.cache.DeletePattern(ctx, pattern); err != nil {
			return fmt.Errorf("failed to invalidate user cache: %w", err)
		}
	}

	return nil
}
