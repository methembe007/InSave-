-- Database Optimization Verification Script
-- Run this after applying performance indexes migration

\echo '=== Database Optimization Verification ==='
\echo ''

-- ============================================================================
-- 1. Verify Indexes Exist
-- ============================================================================
\echo '1. Checking Performance Indexes...'
\echo ''

SELECT 
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND indexname LIKE 'idx_%'
ORDER BY tablename, indexname;

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- 2. Check Index Usage Statistics
-- ============================================================================
\echo '2. Index Usage Statistics (Top 20)...'
\echo ''

SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC
LIMIT 20;

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- 3. Find Unused Indexes
-- ============================================================================
\echo '3. Unused Indexes (candidates for removal)...'
\echo ''

SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND idx_scan = 0
  AND indexname NOT LIKE '%_pkey'
  AND indexname NOT LIKE '%_unique';

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- 4. Table Statistics
-- ============================================================================
\echo '4. Table Statistics...'
\echo ''

SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) as table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename)) as indexes_size,
    n_live_tup as row_count,
    n_dead_tup as dead_rows,
    last_vacuum,
    last_autovacuum,
    last_analyze,
    last_autoanalyze
FROM pg_stat_user_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- 5. Sequential Scans (should be minimal after optimization)
-- ============================================================================
\echo '5. Tables with High Sequential Scans...'
\echo ''

SELECT 
    schemaname,
    tablename,
    seq_scan,
    seq_tup_read,
    idx_scan,
    CASE 
        WHEN seq_scan = 0 THEN 0
        ELSE ROUND((seq_tup_read::numeric / seq_scan), 2)
    END as avg_seq_tup_read,
    CASE 
        WHEN (seq_scan + idx_scan) = 0 THEN 0
        ELSE ROUND((seq_scan::numeric / (seq_scan + idx_scan) * 100), 2)
    END as seq_scan_pct
FROM pg_stat_user_tables
WHERE schemaname = 'public'
  AND seq_scan > 0
ORDER BY seq_scan DESC
LIMIT 20;

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- 6. Verify Materialized View
-- ============================================================================
\echo '6. Materialized View Status...'
\echo ''

SELECT 
    schemaname,
    matviewname,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||matviewname)) as size,
    ispopulated
FROM pg_matviews
WHERE schemaname = 'public';

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- 7. Check Partition Setup
-- ============================================================================
\echo '7. Partition Information...'
\echo ''

SELECT 
    parent.relname AS parent_table,
    child.relname AS partition_name,
    pg_get_expr(child.relpartbound, child.oid) AS partition_bounds,
    pg_size_pretty(pg_total_relation_size(child.oid)) as size
FROM pg_inherits
JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
JOIN pg_class child ON pg_inherits.inhrelid = child.oid
WHERE parent.relname IN ('savings_transactions', 'spending_transactions')
ORDER BY parent.relname, child.relname;

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- 8. Connection Statistics
-- ============================================================================
\echo '8. Connection Statistics...'
\echo ''

SELECT 
    datname,
    numbackends as active_connections,
    xact_commit as commits,
    xact_rollback as rollbacks,
    ROUND((xact_rollback::numeric / NULLIF(xact_commit + xact_rollback, 0) * 100), 2) as rollback_pct,
    blks_read as disk_reads,
    blks_hit as cache_hits,
    ROUND((blks_hit::numeric / NULLIF(blks_hit + blks_read, 0) * 100), 2) as cache_hit_ratio
FROM pg_stat_database
WHERE datname = 'insavein';

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- 9. Slow Queries (requires pg_stat_statements)
-- ============================================================================
\echo '9. Slowest Queries (Top 10)...'
\echo ''

SELECT 
    LEFT(query, 80) as query_preview,
    calls,
    ROUND(total_exec_time::numeric / 1000, 2) as total_time_sec,
    ROUND(mean_exec_time::numeric, 2) as mean_time_ms,
    ROUND(max_exec_time::numeric, 2) as max_time_ms,
    ROUND(stddev_exec_time::numeric, 2) as stddev_time_ms
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
  AND query NOT LIKE '%pg_catalog%'
ORDER BY mean_exec_time DESC
LIMIT 10;

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- 10. Query Performance Test
-- ============================================================================
\echo '10. Testing Query Performance...'
\echo ''

-- Test 1: User savings summary (should use idx_savings_user_created_desc)
\echo 'Test 1: User Savings Summary Query'
EXPLAIN (ANALYZE, BUFFERS)
SELECT 
    user_id,
    COUNT(*) as transaction_count,
    SUM(amount) as total_amount,
    MAX(created_at) as last_transaction
FROM savings_transactions
WHERE user_id = (SELECT id FROM users LIMIT 1)
  AND created_at >= NOW() - INTERVAL '30 days'
GROUP BY user_id;

\echo ''

-- Test 2: Budget with categories (should use JOIN efficiently)
\echo 'Test 2: Budget with Categories Query'
EXPLAIN (ANALYZE, BUFFERS)
SELECT 
    b.id, b.user_id, b.month, b.total_budget, b.total_spent,
    bc.id as category_id, bc.name, bc.allocated_amount, bc.spent_amount
FROM budgets b
LEFT JOIN budget_categories bc ON b.id = bc.budget_id
WHERE b.user_id = (SELECT id FROM users LIMIT 1)
  AND b.month = DATE_TRUNC('month', NOW())::date;

\echo ''

-- Test 3: Unread notifications (should use partial index)
\echo 'Test 3: Unread Notifications Query'
EXPLAIN (ANALYZE, BUFFERS)
SELECT id, type, title, message, created_at
FROM notifications
WHERE user_id = (SELECT id FROM users LIMIT 1)
  AND is_read = false
ORDER BY created_at DESC
LIMIT 20;

\echo ''
\echo '---'
\echo ''

-- ============================================================================
-- Summary
-- ============================================================================
\echo '=== Optimization Verification Complete ==='
\echo ''
\echo 'Review the results above and check:'
\echo '  1. All performance indexes exist'
\echo '  2. Indexes are being used (idx_scan > 0)'
\echo '  3. Sequential scans are minimal'
\echo '  4. Cache hit ratio > 95%'
\echo '  5. Query plans use indexes (Index Scan, not Seq Scan)'
\echo ''
