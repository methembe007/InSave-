-- Create Monthly Partitions Script
-- This script creates partitions for savings_transactions and spending_transactions tables
-- Run this script to create partitions for a specific date range

-- Function to create a partition for a given month
CREATE OR REPLACE FUNCTION create_partition_for_month(
    table_name TEXT,
    partition_date DATE
) RETURNS VOID AS $$
DECLARE
    partition_name TEXT;
    start_date DATE;
    end_date DATE;
BEGIN
    -- Calculate partition boundaries (first day of month to first day of next month)
    start_date := DATE_TRUNC('month', partition_date);
    end_date := start_date + INTERVAL '1 month';
    
    -- Generate partition name (e.g., savings_transactions_2024_01)
    partition_name := table_name || '_' || TO_CHAR(start_date, 'YYYY_MM');
    
    -- Check if partition already exists
    IF NOT EXISTS (
        SELECT 1 FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE c.relname = partition_name
        AND n.nspname = 'public'
    ) THEN
        -- Create the partition
        EXECUTE format(
            'CREATE TABLE %I PARTITION OF %I FOR VALUES FROM (%L) TO (%L)',
            partition_name,
            table_name,
            start_date,
            end_date
        );
        
        RAISE NOTICE 'Created partition: %', partition_name;
    ELSE
        RAISE NOTICE 'Partition already exists: %', partition_name;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Function to create partitions for a date range
CREATE OR REPLACE FUNCTION create_partitions_for_range(
    table_name TEXT,
    start_month DATE,
    end_month DATE
) RETURNS VOID AS $$
DECLARE
    current_month DATE;
BEGIN
    current_month := DATE_TRUNC('month', start_month);
    
    WHILE current_month <= DATE_TRUNC('month', end_month) LOOP
        PERFORM create_partition_for_month(table_name, current_month);
        current_month := current_month + INTERVAL '1 month';
    END LOOP;
    
    RAISE NOTICE 'Partition creation complete for % from % to %', 
        table_name, start_month, end_month;
END;
$$ LANGUAGE plpgsql;

-- Create partitions for the past 6 months, current month, and next 6 months
-- This ensures we have partitions ready for historical data and future inserts

DO $$
DECLARE
    start_date DATE;
    end_date DATE;
BEGIN
    -- Calculate date range (6 months ago to 6 months in the future)
    start_date := DATE_TRUNC('month', CURRENT_DATE - INTERVAL '6 months');
    end_date := DATE_TRUNC('month', CURRENT_DATE + INTERVAL '6 months');
    
    RAISE NOTICE 'Creating partitions for savings_transactions...';
    PERFORM create_partitions_for_range('savings_transactions', start_date, end_date);
    
    RAISE NOTICE 'Creating partitions for spending_transactions...';
    PERFORM create_partitions_for_range('spending_transactions', start_date, end_date);
    
    RAISE NOTICE 'All partitions created successfully!';
END $$;

-- Example: Create partitions for a specific year
-- SELECT create_partitions_for_range('savings_transactions', '2024-01-01', '2024-12-31');
-- SELECT create_partitions_for_range('spending_transactions', '2024-01-01', '2024-12-31');

-- Example: Create partition for a specific month
-- SELECT create_partition_for_month('savings_transactions', '2024-06-01');
-- SELECT create_partition_for_month('spending_transactions', '2024-06-01');
