-- Archive Old Partitions Script
-- This script detaches and optionally archives old partitions to manage database size
-- Partitions older than a specified retention period are detached from the parent table

-- Function to detach a partition (makes it a standalone table)
CREATE OR REPLACE FUNCTION detach_partition(
    table_name TEXT,
    partition_name TEXT
) RETURNS VOID AS $$
BEGIN
    -- Check if partition exists and is attached
    IF EXISTS (
        SELECT 1 FROM pg_inherits
        JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
        JOIN pg_class child ON pg_inherits.inhrelid = child.oid
        WHERE parent.relname = table_name
        AND child.relname = partition_name
    ) THEN
        -- Detach the partition
        EXECUTE format(
            'ALTER TABLE %I DETACH PARTITION %I',
            table_name,
            partition_name
        );
        
        RAISE NOTICE 'Detached partition: %', partition_name;
    ELSE
        RAISE NOTICE 'Partition not found or already detached: %', partition_name;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Function to archive partitions older than retention period
CREATE OR REPLACE FUNCTION archive_old_partitions(
    table_name TEXT,
    retention_months INTEGER DEFAULT 12
) RETURNS TABLE(partition_name TEXT, action TEXT) AS $$
DECLARE
    cutoff_date DATE;
    partition_record RECORD;
    partition_date DATE;
BEGIN
    -- Calculate cutoff date
    cutoff_date := DATE_TRUNC('month', CURRENT_DATE - (retention_months || ' months')::INTERVAL);
    
    RAISE NOTICE 'Archiving partitions older than % (retention: % months)', 
        cutoff_date, retention_months;
    
    -- Find all partitions for the table
    FOR partition_record IN
        SELECT 
            child.relname AS pname
        FROM pg_inherits
        JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
        JOIN pg_class child ON pg_inherits.inhrelid = child.oid
        WHERE parent.relname = table_name
        AND child.relname LIKE table_name || '_%'
        ORDER BY child.relname
    LOOP
        -- Extract date from partition name (format: table_YYYY_MM)
        BEGIN
            partition_date := TO_DATE(
                SUBSTRING(partition_record.pname FROM '(\d{4}_\d{2})$'),
                'YYYY_MM'
            );
            
            -- Check if partition is older than retention period
            IF partition_date < cutoff_date THEN
                -- Detach the partition
                PERFORM detach_partition(table_name, partition_record.pname);
                
                partition_name := partition_record.pname;
                action := 'DETACHED';
                RETURN NEXT;
            END IF;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Could not process partition: %', partition_record.pname;
        END;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Function to drop detached partitions (use with caution!)
CREATE OR REPLACE FUNCTION drop_detached_partition(
    partition_name TEXT
) RETURNS VOID AS $$
BEGIN
    -- Verify partition is detached (not inherited)
    IF NOT EXISTS (
        SELECT 1 FROM pg_inherits
        JOIN pg_class child ON pg_inherits.inhrelid = child.oid
        WHERE child.relname = partition_name
    ) THEN
        -- Drop the table
        EXECUTE format('DROP TABLE IF EXISTS %I', partition_name);
        RAISE NOTICE 'Dropped detached partition: %', partition_name;
    ELSE
        RAISE EXCEPTION 'Cannot drop partition % - it is still attached to a parent table', partition_name;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Function to export partition data to CSV before dropping
CREATE OR REPLACE FUNCTION export_partition_to_csv(
    partition_name TEXT,
    export_path TEXT
) RETURNS VOID AS $$
BEGIN
    -- Export partition data to CSV
    EXECUTE format(
        'COPY %I TO %L WITH (FORMAT CSV, HEADER true)',
        partition_name,
        export_path
    );
    
    RAISE NOTICE 'Exported partition % to %', partition_name, export_path;
END;
$$ LANGUAGE plpgsql;

-- Main archival workflow function
CREATE OR REPLACE FUNCTION archive_workflow(
    table_name TEXT,
    retention_months INTEGER DEFAULT 12,
    export_before_drop BOOLEAN DEFAULT true,
    export_base_path TEXT DEFAULT '/var/lib/postgresql/archives/'
) RETURNS VOID AS $$
DECLARE
    partition_record RECORD;
    export_file TEXT;
BEGIN
    RAISE NOTICE 'Starting archive workflow for % (retention: % months)', 
        table_name, retention_months;
    
    -- Archive old partitions (detach them)
    FOR partition_record IN
        SELECT * FROM archive_old_partitions(table_name, retention_months)
    LOOP
        RAISE NOTICE 'Processing: %', partition_record.partition_name;
        
        -- Export to CSV if requested
        IF export_before_drop THEN
            export_file := export_base_path || partition_record.partition_name || '.csv';
            
            BEGIN
                PERFORM export_partition_to_csv(partition_record.partition_name, export_file);
            EXCEPTION
                WHEN OTHERS THEN
                    RAISE WARNING 'Could not export partition %: %', 
                        partition_record.partition_name, SQLERRM;
            END;
        END IF;
    END LOOP;
    
    RAISE NOTICE 'Archive workflow completed for %', table_name;
END;
$$ LANGUAGE plpgsql;

-- Example usage:

-- 1. Archive partitions older than 12 months (default)
-- SELECT archive_workflow('savings_transactions');
-- SELECT archive_workflow('spending_transactions');

-- 2. Archive partitions older than 24 months
-- SELECT archive_workflow('savings_transactions', 24);
-- SELECT archive_workflow('spending_transactions', 24);

-- 3. Archive without exporting
-- SELECT archive_workflow('savings_transactions', 12, false);

-- 4. List all detached partitions
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
AND (tablename LIKE 'savings_transactions_%' OR tablename LIKE 'spending_transactions_%')
AND tablename NOT IN (
    SELECT child.relname
    FROM pg_inherits
    JOIN pg_class child ON pg_inherits.inhrelid = child.oid
)
ORDER BY tablename;

-- 5. Drop a specific detached partition (CAUTION: Data will be lost!)
-- SELECT drop_detached_partition('savings_transactions_2023_01');

-- 6. View current partition retention
SELECT 
    parent.relname AS table_name,
    child.relname AS partition_name,
    TO_DATE(SUBSTRING(child.relname FROM '(\d{4}_\d{2})$'), 'YYYY_MM') AS partition_month,
    pg_size_pretty(pg_total_relation_size(child.oid)) AS size,
    CASE 
        WHEN TO_DATE(SUBSTRING(child.relname FROM '(\d{4}_\d{2})$'), 'YYYY_MM') < 
             DATE_TRUNC('month', CURRENT_DATE - INTERVAL '12 months')
        THEN 'ELIGIBLE FOR ARCHIVE'
        ELSE 'ACTIVE'
    END AS status
FROM pg_inherits
JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
JOIN pg_class child ON pg_inherits.inhrelid = child.oid
WHERE parent.relname IN ('savings_transactions', 'spending_transactions')
ORDER BY parent.relname, child.relname;
