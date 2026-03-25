# Partition Management Scripts

This directory contains SQL scripts for managing table partitions in the InSavein database. The `savings_transactions` and `spending_transactions` tables are partitioned by month to maintain query performance as data grows.

## Overview

The partition management system provides:
- **Automatic partition creation** for future months
- **Partition maintenance** to ensure partitions exist before data arrives
- **Archival system** to detach and optionally export old partitions
- **Monitoring queries** to track partition health and size

## Scripts

### 1. create_monthly_partitions.sql

Creates partitions for a specified date range.

**When to use:**
- Initial setup after database migration
- Creating partitions for historical data import
- Manually creating partitions for specific months

**Usage:**
```bash
psql -U postgres -d insavein -f create_monthly_partitions.sql
```

**Functions provided:**
- `create_partition_for_month(table_name, partition_date)` - Create single partition
- `create_partitions_for_range(table_name, start_month, end_month)` - Create multiple partitions

**Examples:**
```sql
-- Create partitions for entire year 2024
SELECT create_partitions_for_range('savings_transactions', '2024-01-01', '2024-12-31');
SELECT create_partitions_for_range('spending_transactions', '2024-01-01', '2024-12-31');

-- Create partition for specific month
SELECT create_partition_for_month('savings_transactions', '2024-06-01');
```

### 2. auto_create_partitions.sql

Sets up automatic partition creation for future months.

**When to use:**
- After initial database setup
- To ensure partitions are always ready for new data

**Usage:**
```bash
psql -U postgres -d insavein -f auto_create_partitions.sql
```

**Functions provided:**
- `auto_create_next_partitions()` - Create partitions for next 2 months
- `maintain_partitions()` - Main maintenance function

**Scheduling options:**

**Option 1: pg_cron (recommended)**
```sql
-- Install extension (run once as superuser)
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Schedule daily at 2 AM
SELECT cron.schedule(
    'partition-maintenance',
    '0 2 * * *',
    'SELECT maintain_partitions();'
);
```

**Option 2: System cron (Linux/Mac)**
```bash
# Add to crontab
0 2 * * * psql -U postgres -d insavein -c "SELECT maintain_partitions();"
```

**Option 3: Windows Task Scheduler**
- Create scheduled task running daily at 2 AM
- Action: `psql -U postgres -d insavein -c "SELECT maintain_partitions();"`

### 3. archive_old_partitions.sql

Manages archival of old partitions to control database size.

**When to use:**
- Regularly (monthly or quarterly) to manage database growth
- When database size becomes a concern
- Before major upgrades or migrations

**Usage:**
```bash
psql -U postgres -d insavein -f archive_old_partitions.sql
```

**Functions provided:**
- `detach_partition(table_name, partition_name)` - Detach single partition
- `archive_old_partitions(table_name, retention_months)` - Detach old partitions
- `drop_detached_partition(partition_name)` - Drop detached partition (CAUTION!)
- `export_partition_to_csv(partition_name, export_path)` - Export before dropping
- `archive_workflow(table_name, retention_months, export_before_drop, export_base_path)` - Complete workflow

**Examples:**
```sql
-- Archive partitions older than 12 months (default)
SELECT archive_workflow('savings_transactions');
SELECT archive_workflow('spending_transactions');

-- Archive partitions older than 24 months
SELECT archive_workflow('savings_transactions', 24);

-- Archive without exporting (not recommended)
SELECT archive_workflow('savings_transactions', 12, false);

-- List detached partitions
SELECT tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
AND (tablename LIKE 'savings_transactions_%' OR tablename LIKE 'spending_transactions_%')
AND tablename NOT IN (
    SELECT child.relname FROM pg_inherits
    JOIN pg_class child ON pg_inherits.inhrelid = child.oid
);

-- Drop specific detached partition (CAUTION!)
SELECT drop_detached_partition('savings_transactions_2023_01');
```

## Partition Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│                    Partition Lifecycle                       │
└─────────────────────────────────────────────────────────────┘

1. CREATION (2 months before needed)
   ↓
   auto_create_next_partitions() runs daily
   ↓
2. ACTIVE (current month ± 12 months)
   ↓
   Data is inserted and queried normally
   ↓
3. ARCHIVAL (older than retention period)
   ↓
   archive_workflow() detaches partition
   ↓
4. DETACHED (standalone table)
   ↓
   Optional: export_partition_to_csv()
   ↓
5. DROPPED (optional, permanent deletion)
   ↓
   drop_detached_partition()
```

## Monitoring Queries

### View all partitions with status
```sql
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
ORDER BY parent.relname, partition_month DESC;
```

### Check partition sizes
```sql
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_total_relation_size(schemaname||'.'||tablename) AS bytes
FROM pg_tables
WHERE tablename LIKE 'savings_transactions_%' 
   OR tablename LIKE 'spending_transactions_%'
ORDER BY bytes DESC;
```

### Verify future partitions exist
```sql
SELECT 
    table_name,
    partition_name,
    partition_month
FROM (
    SELECT 
        parent.relname AS table_name,
        child.relname AS partition_name,
        TO_DATE(SUBSTRING(child.relname FROM '(\d{4}_\d{2})$'), 'YYYY_MM') AS partition_month
    FROM pg_inherits
    JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
    JOIN pg_class child ON pg_inherits.inhrelid = child.oid
    WHERE parent.relname IN ('savings_transactions', 'spending_transactions')
) partitions
WHERE partition_month >= DATE_TRUNC('month', CURRENT_DATE)
ORDER BY table_name, partition_month;
```

### Count rows per partition
```sql
DO $$
DECLARE
    partition_record RECORD;
    row_count BIGINT;
BEGIN
    FOR partition_record IN
        SELECT child.relname AS partition_name
        FROM pg_inherits
        JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
        JOIN pg_class child ON pg_inherits.inhrelid = child.oid
        WHERE parent.relname IN ('savings_transactions', 'spending_transactions')
        ORDER BY child.relname
    LOOP
        EXECUTE format('SELECT COUNT(*) FROM %I', partition_record.partition_name) INTO row_count;
        RAISE NOTICE 'Partition %: % rows', partition_record.partition_name, row_count;
    END LOOP;
END $$;
```

## Best Practices

### Retention Policy
- **Recommended**: 12-24 months of active data
- **Minimum**: 6 months for analytics and reporting
- **Maximum**: Depends on storage capacity and query performance

### Maintenance Schedule
- **Daily**: Run `maintain_partitions()` to create future partitions
- **Monthly**: Review partition sizes and growth trends
- **Quarterly**: Run archival workflow for old partitions
- **Annually**: Review retention policy and adjust if needed

### Before Archiving
1. **Verify backups** are current and tested
2. **Export partitions** to CSV or external storage
3. **Test restore process** from exports
4. **Document** which partitions were archived and when

### Performance Considerations
- Partitions improve query performance when queries filter by date
- Too many partitions (>100) can slow down query planning
- Archive old partitions to keep partition count manageable
- Monitor partition sizes - aim for 10-50 GB per partition

### Disaster Recovery
- Detached partitions can be reattached if needed:
  ```sql
  ALTER TABLE savings_transactions 
  ATTACH PARTITION savings_transactions_2023_01 
  FOR VALUES FROM ('2023-01-01') TO ('2023-02-01');
  ```
- Keep CSV exports for at least 90 days after archival
- Store exports in geographically separate location

## Troubleshooting

### Partition creation fails
```sql
-- Check if partition already exists
SELECT tablename FROM pg_tables 
WHERE tablename = 'savings_transactions_2024_06';

-- Check parent table exists
SELECT tablename FROM pg_tables 
WHERE tablename = 'savings_transactions';
```

### Cannot insert data
```sql
-- Verify partition exists for the date
SELECT * FROM pg_tables 
WHERE tablename LIKE 'savings_transactions_' || TO_CHAR(CURRENT_DATE, 'YYYY_MM');

-- Create missing partition
SELECT create_partition_for_month('savings_transactions', CURRENT_DATE);
```

### Partition detachment fails
```sql
-- Check if partition has dependencies
SELECT 
    conname AS constraint_name,
    contype AS constraint_type
FROM pg_constraint
WHERE conrelid = 'savings_transactions_2023_01'::regclass;

-- Drop constraints if safe, then retry detachment
```

### Export fails with permission error
```bash
# Ensure PostgreSQL has write permissions to export directory
sudo chown postgres:postgres /var/lib/postgresql/archives/
sudo chmod 755 /var/lib/postgresql/archives/
```

## Initial Setup Checklist

- [ ] Run `create_monthly_partitions.sql` to create initial partitions
- [ ] Run `auto_create_partitions.sql` to set up automatic creation
- [ ] Schedule `maintain_partitions()` to run daily
- [ ] Set up monitoring alerts for partition health
- [ ] Document retention policy
- [ ] Test archival workflow on non-production environment
- [ ] Set up backup for detached partitions
- [ ] Schedule quarterly archival reviews
