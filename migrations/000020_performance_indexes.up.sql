-- Performance Optimization Indexes
-- Task 28.3: Add missing indexes for query optimization

-- ============================================================================
-- SAVINGS TRANSACTIONS INDEXES
-- ============================================================================

-- Index for user savings summary queries (most common query)
CREATE INDEX IF NOT EXISTS idx_savings_user_created_desc 
ON savings_transactions(user_id, created_at DESC);

-- Index for category-based queries
CREATE INDEX IF NOT EXISTS idx_savings_user_category 
ON savings_transactions(user_id, category);

-- Index for date range queries (monthly stats)
CREATE INDEX IF NOT EXISTS idx_savings_created_date 
ON savings_transactions(DATE(created_at));

-- Composite index for streak calculations
CREATE INDEX IF NOT EXISTS idx_savings_user_date_only 
ON savings_transactions(user_id, DATE(created_at));

-- ============================================================================
-- SPENDING TRANSACTIONS INDEXES
-- ============================================================================

-- Index for budget spending queries
CREATE INDEX IF NOT EXISTS idx_spending_budget_category 
ON spending_transactions(budget_id, category_id);

-- Index for user spending history
CREATE INDEX IF NOT EXISTS idx_spending_user_date_desc 
ON spending_transactions(user_id, transaction_date DESC);

-- Index for merchant analysis
CREATE INDEX IF NOT EXISTS idx_spending_user_merchant 
ON spending_transactions(user_id, merchant);

-- Index for date range analytics
CREATE INDEX IF NOT EXISTS idx_spending_date_range 
ON spending_transactions(transaction_date);

-- ============================================================================
-- BUDGET INDEXES
-- ============================================================================

-- Index for current budget lookup (most common)
CREATE INDEX IF NOT EXISTS idx_budgets_user_month_desc 
ON budgets(user_id, month DESC);

-- Index for budget categories lookup
CREATE INDEX IF NOT EXISTS idx_budget_categories_budget_spent 
ON budget_categories(budget_id, spent_amount);

-- ============================================================================
-- GOALS INDEXES
-- ============================================================================

-- Index for active goals query
CREATE INDEX IF NOT EXISTS idx_goals_user_status_date 
ON goals(user_id, status, target_date);

-- Index for goal progress tracking
CREATE INDEX IF NOT EXISTS idx_goals_status_current_target 
ON goals(status, current_amount, target_amount);

-- Index for milestone queries
CREATE INDEX IF NOT EXISTS idx_milestones_goal_order 
ON goal_milestones(goal_id, "order", is_completed);

-- Index for completed milestones
CREATE INDEX IF NOT EXISTS idx_milestones_goal_completed 
ON goal_milestones(goal_id, is_completed, amount);

-- ============================================================================
-- NOTIFICATIONS INDEXES
-- ============================================================================

-- Index for unread notifications (most common query)
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread 
ON notifications(user_id, is_read, created_at DESC) 
WHERE is_read = false;

-- Index for notification history
CREATE INDEX IF NOT EXISTS idx_notifications_user_created 
ON notifications(user_id, created_at DESC);

-- Index for notification type filtering
CREATE INDEX IF NOT EXISTS idx_notifications_user_type 
ON notifications(user_id, type, created_at DESC);

-- ============================================================================
-- EDUCATION PROGRESS INDEXES
-- ============================================================================

-- Index for user progress queries
CREATE INDEX IF NOT EXISTS idx_education_progress_user_completed 
ON education_progress(user_id, is_completed);

-- Index for lesson completion lookup
CREATE INDEX IF NOT EXISTS idx_education_progress_user_lesson 
ON education_progress(user_id, lesson_id);

-- Composite index for progress calculation
CREATE INDEX IF NOT EXISTS idx_education_progress_user_stats 
ON education_progress(user_id, is_completed, completed_at);

-- ============================================================================
-- USERS INDEXES (Additional)
-- ============================================================================

-- Index for user preferences queries
CREATE INDEX IF NOT EXISTS idx_users_preferences 
ON users USING GIN (preferences);

-- ============================================================================
-- ANALYTICS OPTIMIZATION
-- ============================================================================

-- Materialized view for financial health score calculation
-- This pre-computes expensive aggregations
CREATE MATERIALIZED VIEW IF NOT EXISTS mv_user_financial_stats AS
SELECT 
    u.id as user_id,
    -- Savings stats
    COALESCE(SUM(st.amount), 0) as total_saved,
    COUNT(DISTINCT DATE(st.created_at)) as savings_days,
    COALESCE(AVG(st.amount), 0) as avg_savings_amount,
    -- Budget stats
    COALESCE(SUM(b.total_budget), 0) as total_budget,
    COALESCE(SUM(b.total_spent), 0) as total_spent,
    -- Spending stats
    COUNT(DISTINCT sp.id) as spending_transactions,
    COALESCE(AVG(sp.amount), 0) as avg_spending_amount,
    -- Goals stats
    COUNT(DISTINCT g.id) as total_goals,
    COUNT(DISTINCT CASE WHEN g.status = 'completed' THEN g.id END) as completed_goals,
    -- Last activity
    GREATEST(
        MAX(st.created_at),
        MAX(sp.created_at),
        MAX(g.updated_at)
    ) as last_activity
FROM users u
LEFT JOIN savings_transactions st ON u.id = st.user_id 
    AND st.created_at >= NOW() - INTERVAL '90 days'
LEFT JOIN budgets b ON u.id = b.user_id 
    AND b.month >= DATE_TRUNC('month', NOW() - INTERVAL '6 months')
LEFT JOIN spending_transactions sp ON u.id = sp.user_id 
    AND sp.transaction_date >= NOW() - INTERVAL '90 days'
LEFT JOIN goals g ON u.id = g.user_id
GROUP BY u.id;

-- Index on materialized view
CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_financial_stats_user 
ON mv_user_financial_stats(user_id);

-- ============================================================================
-- PARTITION MAINTENANCE
-- ============================================================================

-- Function to create future partitions automatically
CREATE OR REPLACE FUNCTION create_monthly_partitions()
RETURNS void AS $$
DECLARE
    start_date date;
    end_date date;
    partition_name text;
BEGIN
    -- Create partitions for next 3 months
    FOR i IN 0..2 LOOP
        start_date := DATE_TRUNC('month', NOW() + (i || ' months')::interval);
        end_date := start_date + INTERVAL '1 month';
        
        -- Savings transactions partition
        partition_name := 'savings_transactions_' || TO_CHAR(start_date, 'YYYY_MM');
        IF NOT EXISTS (
            SELECT 1 FROM pg_tables 
            WHERE tablename = partition_name
        ) THEN
            EXECUTE format(
                'CREATE TABLE IF NOT EXISTS %I PARTITION OF savings_transactions 
                FOR VALUES FROM (%L) TO (%L)',
                partition_name, start_date, end_date
            );
            RAISE NOTICE 'Created partition: %', partition_name;
        END IF;
        
        -- Spending transactions partition
        partition_name := 'spending_transactions_' || TO_CHAR(start_date, 'YYYY_MM');
        IF NOT EXISTS (
            SELECT 1 FROM pg_tables 
            WHERE tablename = partition_name
        ) THEN
            EXECUTE format(
                'CREATE TABLE IF NOT EXISTS %I PARTITION OF spending_transactions 
                FOR VALUES FROM (%L) TO (%L)',
                partition_name, start_date, end_date
            );
            RAISE NOTICE 'Created partition: %', partition_name;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- QUERY OPTIMIZATION VIEWS
-- ============================================================================

-- View for recent user activity (commonly queried)
CREATE OR REPLACE VIEW v_recent_user_activity AS
SELECT 
    user_id,
    'savings' as activity_type,
    amount,
    description,
    created_at as activity_date
FROM savings_transactions
WHERE created_at >= NOW() - INTERVAL '30 days'
UNION ALL
SELECT 
    user_id,
    'spending' as activity_type,
    amount,
    description,
    created_at as activity_date
FROM spending_transactions
WHERE transaction_date >= NOW() - INTERVAL '30 days'
ORDER BY activity_date DESC;

-- ============================================================================
-- STATISTICS UPDATE
-- ============================================================================

-- Analyze all tables to update query planner statistics
ANALYZE users;
ANALYZE savings_transactions;
ANALYZE spending_transactions;
ANALYZE budgets;
ANALYZE budget_categories;
ANALYZE goals;
ANALYZE goal_milestones;
ANALYZE notifications;
ANALYZE lessons;
ANALYZE education_progress;

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON INDEX idx_savings_user_created_desc IS 'Optimizes user savings summary and history queries';
COMMENT ON INDEX idx_spending_budget_category IS 'Optimizes budget spending queries and alerts';
COMMENT ON INDEX idx_notifications_user_unread IS 'Partial index for unread notifications - most common query';
COMMENT ON MATERIALIZED VIEW mv_user_financial_stats IS 'Pre-computed financial statistics for analytics service';
COMMENT ON FUNCTION create_monthly_partitions IS 'Automatically creates future monthly partitions';
