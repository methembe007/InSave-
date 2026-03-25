-- Rollback performance optimization indexes

-- Drop materialized view
DROP MATERIALIZED VIEW IF EXISTS mv_user_financial_stats;

-- Drop views
DROP VIEW IF EXISTS v_recent_user_activity;

-- Drop function
DROP FUNCTION IF EXISTS create_monthly_partitions();

-- Drop savings transaction indexes
DROP INDEX IF EXISTS idx_savings_user_created_desc;
DROP INDEX IF EXISTS idx_savings_user_category;
DROP INDEX IF EXISTS idx_savings_created_date;
DROP INDEX IF EXISTS idx_savings_user_date_only;

-- Drop spending transaction indexes
DROP INDEX IF EXISTS idx_spending_budget_category;
DROP INDEX IF EXISTS idx_spending_user_date_desc;
DROP INDEX IF EXISTS idx_spending_user_merchant;
DROP INDEX IF EXISTS idx_spending_date_range;

-- Drop budget indexes
DROP INDEX IF EXISTS idx_budgets_user_month_desc;
DROP INDEX IF EXISTS idx_budget_categories_budget_spent;

-- Drop goal indexes
DROP INDEX IF EXISTS idx_goals_user_status_date;
DROP INDEX IF EXISTS idx_goals_status_current_target;
DROP INDEX IF EXISTS idx_milestones_goal_order;
DROP INDEX IF EXISTS idx_milestones_goal_completed;

-- Drop notification indexes
DROP INDEX IF EXISTS idx_notifications_user_unread;
DROP INDEX IF EXISTS idx_notifications_user_created;
DROP INDEX IF EXISTS idx_notifications_user_type;

-- Drop education progress indexes
DROP INDEX IF EXISTS idx_education_progress_user_completed;
DROP INDEX IF EXISTS idx_education_progress_user_lesson;
DROP INDEX IF EXISTS idx_education_progress_user_stats;

-- Drop user indexes
DROP INDEX IF EXISTS idx_users_preferences;
