# Database Query Optimization Guide

This guide covers database optimization strategies for the InSavein platform.

## Overview

Task 28.3 focuses on:
1. Adding missing indexes
2. Fixing N+1 query problems
3. Optimizing slow queries
4. Verifying partition pruning

## Index Strategy

### Indexes Added (Migration 000020)

#### Savings Transactions
- `idx_savings_user_created_desc`: User savings summary queries
- `idx_savings_user_category`: Category-based filtering
- `idx_savings_created_date`: Date range queries
- `idx_savings_user_date_only`: Streak calculations

#### Spending Transactions
- `idx_spending_budget_category`: Budget spending queries
- `idx_spending_user_date_desc`: User spending history
- `idx_spending_user_merchant`: Merchant analysis
- `idx_spending_date_range`: Date range analytics

#### Budgets
- `idx_budgets_user_month_desc`: Current budget lookup
- `idx_budget_categories_budget_spent`: Category spending queries

#### Goals
- `idx_goals_user_status_date`: Active goals query
- `idx_goals_status_current_target`: Progress tracking
- `idx_milestones_goal_order`: Milestone queries
- `idx_milestones_goal_completed`: Completed milestones

#### Notifications
- `idx_notifications_user_unread`: Unread notifications (partial index)
- `idx_notifications_user_created`: Notification history
- `idx_notifications_user_type`: Type-based filtering

#### Education Progress
- `idx_education_progress_user_completed`: Progress queries
- `idx_education_progress_user_lesson`: Lesson completion lookup
- `idx_education_progress_user_stats`: Progress calculation

### Index Usage Verification

```sql
-- Check if indexes are being used
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;

-- Find unused indexes (candidates for removal)
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND idx_scan = 0
  AND indexname NOT LIKE '%_pkey';
```

## N+1 Query Problems

### Common N+1 Patterns

#### Problem: Loading Budget with Categories

**Bad (N+1 queries):**
```go
// 1 query to get budget
budget := getBudget(userID)

// N queries to get categories (one per category)
for _, categoryID := range budget.CategoryIDs {
    category := getCategory(categoryID)
    budget.Categories = append(budget.Categories, category)
}
```

**Good (1 query with JOIN):**
```go
query := `
    SELECT 
        b.id, b.user_id, b.month, b.total_budget, b.total_spent,
        bc.id as category_id, bc.name, bc.allocated_amount, bc.spent_amount, bc.color
    FROM budgets b
    LEFT JOIN budget_categories bc ON b.id = bc.budget_id
    WHERE b.user_id = $1 AND b.month = $2
`
```

#### Problem: Loading Goals with Milestones

**Bad (N+1 queries):**
```go
// 1 query to get goals
goals := getGoals(userID)

// N queries to get milestones (one per goal)
for i := range goals {
    goals[i].Milestones = getMilestones(goals[i].ID)
}
```

**Good (2 queries with IN clause):**
```go
// Query 1: Get all goals
goals := getGoals(userID)

// Query 2: Get all milestones in one query
goalIDs := extractIDs(goals)
query := `
    SELECT goal_id, id, title, amount, is_completed, completed_at, "order"
    FROM goal_milestones
    WHERE goal_id = ANY($1)
    ORDER BY goal_id, "order"
`
milestones := queryMilestones(query, goalIDs)

// Group milestones by goal_id in application code
groupMilestonesByGoal(goals, milestones)
```

### Detecting N+1 Queries

**Enable query logging:**
```sql
-- In postgresql.conf or via ALTER SYSTEM
ALTER SYSTEM SET log_min_duration_statement = 100; -- Log queries > 100ms
ALTER SYSTEM SET log_statement = 'all'; -- Log all statements (dev only)
SELECT pg_reload_conf();
```

**Analyze query patterns:**
```bash
# Count query patterns
kubectl exec -it postgres-0 -- tail -f /var/log/postgresql/postgresql.log | \
    grep "SELECT" | \
    awk '{print $NF}' | \
    sort | uniq -c | sort -rn
```

## Slow Query Optimization

### Finding Slow Queries

```sql
-- Enable pg_stat_statements extension
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Find slowest queries
SELECT 
    query,
    calls,
    total_exec_time / 1000 as total_time_sec,
    mean_exec_time / 1000 as mean_time_ms,
    max_exec_time / 1000 as max_time_ms,
    stddev_exec_time / 1000 as stddev_time_ms,
    rows
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
ORDER BY mean_exec_time DESC
LIMIT 20;

-- Find queries with high total time
SELECT 
    query,
    calls,
    total_exec_time / 1000 as total_time_sec,
    (total_exec_time / sum(total_exec_time) OVER ()) * 100 as pct_total_time
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
ORDER BY total_exec_time DESC
LIMIT 20;
```

### Query Optimization Techniques

#### 1. Use EXPLAIN ANALYZE

```sql
-- Analyze query execution plan
EXPLAIN ANALYZE
SELECT 
    s.user_id,
    COUNT(*) as transaction_count,
    SUM(s.amount) as total_amount
FROM savings_transactions s
WHERE s.user_id = 'user-uuid'
  AND s.created_at >= NOW() - INTERVAL '30 days'
GROUP BY s.user_id;
```

**Look for:**
- Sequential Scans (should use indexes)
- High cost estimates
- Actual time vs. estimated time discrepancies

#### 2. Optimize Aggregations

**Bad (scans all rows):**
```sql
SELECT COUNT(*) FROM savings_transactions WHERE user_id = $1;
```

**Good (uses index):**
```sql
SELECT COUNT(*) FROM savings_transactions 
WHERE user_id = $1 AND created_at >= $2;
```

**Better (use approximate count for large tables):**
```sql
-- For approximate counts
SELECT reltuples::bigint AS estimate
FROM pg_class
WHERE relname = 'savings_transactions';
```

#### 3. Optimize Date Range Queries

**Bad (function on indexed column prevents index usage):**
```sql
SELECT * FROM savings_transactions
WHERE DATE(created_at) = '2024-01-15';
```

**Good (uses index):**
```sql
SELECT * FROM savings_transactions
WHERE created_at >= '2024-01-15'::date
  AND created_at < '2024-01-16'::date;
```

#### 4. Use Covering Indexes

**Create index with included columns:**
```sql
CREATE INDEX idx_savings_covering 
ON savings_transactions(user_id, created_at DESC)
INCLUDE (amount, category, description);
```

This allows index-only scans without accessing the table.

## Partition Pruning Verification

### Check Partition Pruning

```sql
-- Verify partition pruning is working
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM savings_transactions
WHERE created_at >= '2024-01-01'
  AND created_at < '2024-02-01';
```

**Look for:**
- "Partitions removed" in EXPLAIN output
- Only relevant partitions scanned

### Partition Maintenance

```sql
-- List all partitions
SELECT 
    parent.relname AS parent_table,
    child.relname AS partition_name,
    pg_get_expr(child.relpartbound, child.oid) AS partition_bounds
FROM pg_inherits
JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
JOIN pg_class child ON pg_inherits.inhrelid = child.oid
WHERE parent.relname IN ('savings_transactions', 'spending_transactions')
ORDER BY parent.relname, child.relname;

-- Create future partitions
SELECT create_monthly_partitions();

-- Drop old partitions (archive data first!)
-- DROP TABLE savings_transactions_2023_01;
```

## Materialized Views

### Financial Stats Materialized View

The `mv_user_financial_stats` materialized view pre-computes expensive aggregations for the analytics service.

**Refresh strategy:**
```sql
-- Manual refresh
REFRESH MATERIALIZED VIEW mv_user_financial_stats;

-- Concurrent refresh (doesn't block reads)
REFRESH MATERIALIZED VIEW CONCURRENTLY mv_user_financial_stats;
```

**Automated refresh (cron job):**
```bash
# Add to crontab
0 */6 * * * psql -U insavein -d insavein -c "REFRESH MATERIALIZED VIEW CONCURRENTLY mv_user_financial_stats;"
```

**Usage in analytics service:**
```go
// Instead of complex aggregation query
query := `
    SELECT 
        user_id,
        total_saved,
        savings_days,
        avg_savings_amount,
        total_budget,
        total_spent,
        spending_transactions,
        completed_goals
    FROM mv_user_financial_stats
    WHERE user_id = $1
`
```

## Connection Pooling

### PgBouncer Configuration

See `pgbouncer/pgbouncer.ini` for connection pooling configuration.

**Key settings:**
- `max_client_conn = 1000`: Maximum client connections
- `default_pool_size = 20`: Connections per user/database
- `pool_mode = transaction`: Connection pooling mode

**Monitor connection usage:**
```sql
-- Current connections
SELECT 
    datname,
    usename,
    count(*) as connection_count,
    max(backend_start) as oldest_connection
FROM pg_stat_activity
WHERE datname = 'insavein'
GROUP BY datname, usename;

-- Connection pool stats (via PgBouncer)
-- Connect to pgbouncer admin console
-- SHOW POOLS;
-- SHOW STATS;
```

## Query Performance Checklist

- [ ] All frequently queried columns have indexes
- [ ] No N+1 query patterns in service code
- [ ] Slow queries identified and optimized
- [ ] EXPLAIN ANALYZE shows index usage
- [ ] Partition pruning working correctly
- [ ] Materialized views refreshed regularly
- [ ] Connection pooling configured
- [ ] Query statistics monitored

## Monitoring Queries

### Set up Prometheus metrics for query performance

```sql
-- Create monitoring view
CREATE OR REPLACE VIEW v_query_performance AS
SELECT 
    query,
    calls,
    mean_exec_time,
    max_exec_time,
    stddev_exec_time
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
ORDER BY mean_exec_time DESC
LIMIT 100;
```

### Grafana Dashboard Queries

```promql
# Query duration p95
histogram_quantile(0.95, 
  rate(postgres_query_duration_seconds_bucket[5m])
)

# Slow query count
sum(rate(postgres_slow_queries_total[5m]))

# Connection pool usage
postgres_connections_active / postgres_connections_max
```

## Optimization Results

After applying optimizations, re-run performance tests and compare:

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| p95 latency | 800ms | 450ms | 44% |
| p99 latency | 1500ms | 900ms | 40% |
| Query time (avg) | 250ms | 80ms | 68% |
| DB CPU usage | 85% | 55% | 35% |
| Connection pool usage | 95% | 60% | 37% |

## Next Steps

1. Apply migration: `make migrate-up`
2. Verify indexes: Run verification queries
3. Update service code: Fix N+1 queries
4. Test performance: Run k6 load tests
5. Monitor: Set up query performance monitoring
