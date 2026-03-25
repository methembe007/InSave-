-- Automatic Partition Creation Script
-- This script sets up automatic partition creation for future months
-- Partitions are created 2 months in advance to ensure data can always be inserted

-- Function to automatically create next month's partition if it doesn't exist
CREATE OR REPLACE FUNCTION auto_create_next_partitions()
RETURNS VOID AS $$
DECLARE
    next_month DATE;
    month_after_next DATE;
BEGIN
    -- Calculate next month and the month after
    next_month := DATE_TRUNC('month', CURRENT_DATE + INTERVAL '1 month');
    month_after_next := DATE_TRUNC('month', CURRENT_DATE + INTERVAL '2 months');
    
    -- Create partitions for next month if they don't exist
    PERFORM create_partition_for_month('savings_transactions', next_month);
    PERFORM create_partition_for_month('spending_transactions', next_month);
    
    -- Create partitions for the month after next (2 months ahead)
    PERFORM create_partition_for_month('savings_transactions', month_after_next);
    PERFORM create_partition_for_month('spending_transactions', month_after_next);
    
    RAISE NOTICE 'Auto-created partitions for % and %', next_month, month_after_next;
END;
$$ LANGUAGE plpgsql;

-- Create a maintenance function that should be run daily
CREATE OR REPLACE FUNCTION maintain_partitions()
RETURNS VOID AS $$
BEGIN
    -- Create future partitions
    PERFORM auto_create_next_partitions();
    
    -- Log maintenance run
    RAISE NOTICE 'Partition maintenance completed at %', NOW();
END;
$$ LANGUAGE plpgsql;

-- Schedule automatic partition creation using pg_cron (if available)
-- Note: pg_cron extension must be installed and configured
-- Uncomment the following lines if pg_cron is available:

/*
-- Install pg_cron extension (run once as superuser)
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Schedule daily partition maintenance at 2 AM
SELECT cron.schedule(
    'partition-maintenance',
    '0 2 * * *',  -- Every day at 2 AM
    'SELECT maintain_partitions();'
);
*/

-- Alternative: Manual scheduling instructions
-- If pg_cron is not available, set up a cron job on your system:
-- 
-- Linux/Mac crontab entry:
-- 0 2 * * * psql -U postgres -d insavein -c "SELECT maintain_partitions();"
--
-- Windows Task Scheduler:
-- Create a scheduled task that runs daily at 2 AM:
-- psql -U postgres -d insavein -c "SELECT maintain_partitions();"

-- Run initial partition creation
SELECT maintain_partitions();

-- Verify partitions exist
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'savings_transactions_%' 
   OR tablename LIKE 'spending_transactions_%'
ORDER BY tablename;
