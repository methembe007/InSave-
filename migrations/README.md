# InSavein Platform Database Migrations

This directory contains PostgreSQL database migrations for the InSavein Platform using golang-migrate.

## Overview

The database schema includes:
- **Users**: User accounts with bcrypt password hashing (cost factor 12)
- **Savings Transactions**: Partitioned by month for performance
- **Budgets**: Monthly budget planning
- **Budget Categories**: Category allocations within budgets
- **Spending Transactions**: Partitioned by transaction date for performance
- **Goals**: Financial goal tracking
- **Goal Milestones**: Milestone tracking for goals
- **Notifications**: User notification system
- **Lessons**: Financial education content
- **Education Progress**: User lesson completion tracking

## Key Features

### Table Partitioning
Two tables are partitioned for optimal query performance:
- `savings_transactions`: Partitioned by `created_at` (monthly partitions)
- `spending_transactions`: Partitioned by `transaction_date` (monthly partitions)

Partitions are pre-created for 2024-2027. Additional partitions can be added as needed.

### Data Integrity
- All amounts stored as `DECIMAL(15,2)` with 2 decimal places precision
- CASCADE deletion for user data (Requirement 3.4)
- CHECK constraints for data validation
- UNIQUE constraints to prevent duplicates
- Foreign key relationships with appropriate ON DELETE actions

### Indexes
Comprehensive indexes for query optimization:
- User lookups by email
- Transaction queries by user_id and date
- Budget queries by user and month
- Goal queries by user and status
- Notification queries by user and read status

## Installation

### Prerequisites
- PostgreSQL 12 or higher
- golang-migrate CLI tool

### Install golang-migrate

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

**Windows:**
```bash
scoop install migrate
```

Or download from: https://github.com/golang-migrate/migrate/releases

## Usage

### Set Database URL
```bash
export DATABASE_URL="postgresql://username:password@localhost:5432/insavein?sslmode=disable"
```

### Run Migrations (Up)
Apply all pending migrations:
```bash
migrate -database ${DATABASE_URL} -path migrations up
```

Apply specific number of migrations:
```bash
migrate -database ${DATABASE_URL} -path migrations up 5
```

### Rollback Migrations (Down)
Rollback all migrations:
```bash
migrate -database ${DATABASE_URL} -path migrations down
```

Rollback specific number of migrations:
```bash
migrate -database ${DATABASE_URL} -path migrations down 1
```

### Check Migration Version
```bash
migrate -database ${DATABASE_URL} -path migrations version
```

### Force Migration Version (Use with caution)
If migrations are in a dirty state:
```bash
migrate -database ${DATABASE_URL} -path migrations force VERSION
```

## Migration Files

Migrations follow the naming convention: `{version}_{description}.{up|down}.sql`

| Version | Description | Tables Created |
|---------|-------------|----------------|
| 000001 | Create users table | users |
| 000002 | Create savings transactions table | savings_transactions (partitioned) |
| 000003 | Create budgets table | budgets |
| 000004 | Create budget categories table | budget_categories |
| 000005 | Create spending transactions table | spending_transactions (partitioned) |
| 000006 | Create goals table | goals |
| 000007 | Create goal milestones table | goal_milestones |
| 000008 | Create notifications table | notifications |
| 000009 | Create lessons table | lessons |
| 000010 | Create education progress table | education_progress |

## Adding New Partitions

As time progresses, you'll need to add new monthly partitions for the partitioned tables.

### Example: Adding 2028 Partitions

Create a new migration file:
```bash
migrate create -ext sql -dir migrations -seq add_2028_partitions
```

In the `.up.sql` file:
```sql
-- Savings transactions 2028 partitions
CREATE TABLE savings_transactions_2028_01 PARTITION OF savings_transactions
    FOR VALUES FROM ('2028-01-01') TO ('2028-02-01');
-- ... repeat for all months

-- Spending transactions 2028 partitions
CREATE TABLE spending_transactions_2028_01 PARTITION OF spending_transactions
    FOR VALUES FROM ('2028-01-01') TO ('2028-02-01');
-- ... repeat for all months
```

## Development Workflow

### Local Development
1. Start PostgreSQL locally
2. Create database: `createdb insavein`
3. Run migrations: `migrate -database ${DATABASE_URL} -path migrations up`

### Testing Migrations
1. Apply migrations: `migrate up`
2. Verify schema: `psql -d insavein -c "\dt"`
3. Test rollback: `migrate down 1`
4. Re-apply: `migrate up`

### CI/CD Integration
Add to your CI/CD pipeline:
```bash
# Run migrations in CI
migrate -database ${DATABASE_URL} -path migrations up

# Verify migrations are clean
migrate -database ${DATABASE_URL} -path migrations version
```

## Security Considerations

- **Password Hashing**: Use bcrypt with cost factor 12 (Requirement 1.2)
- **Connection Security**: Always use SSL/TLS in production
- **Least Privilege**: Grant only necessary permissions to application database user
- **Backup**: Regular automated backups (Requirement 25)
- **Encryption**: Enable database-level encryption at rest (Requirement 20)

## Performance Optimization

### Partitioning Strategy
- Monthly partitions reduce query scan time
- Partition pruning automatically excludes irrelevant partitions
- Indexes on each partition improve query performance

### Index Strategy
- Composite indexes for common query patterns
- Covering indexes where beneficial
- Regular ANALYZE to update statistics

### Monitoring
Monitor partition sizes and query performance:
```sql
-- Check partition sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'savings_transactions_%'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Check index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

## Troubleshooting

### Migration Failed
If a migration fails:
1. Check the error message
2. Fix the SQL in the migration file
3. Force the version: `migrate force VERSION`
4. Re-run: `migrate up`

### Dirty Database State
If migrations show "dirty" state:
1. Check what went wrong: `migrate version`
2. Manually fix the database if needed
3. Force to correct version: `migrate force VERSION`

### Connection Issues
- Verify DATABASE_URL is correct
- Check PostgreSQL is running
- Verify network connectivity
- Check firewall rules

## References

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Partitioning](https://www.postgresql.org/docs/current/ddl-partitioning.html)
- [PostgreSQL Indexes](https://www.postgresql.org/docs/current/indexes.html)
- InSavein Platform Requirements: `.kiro/specs/insavein-platform/requirements.md`
- InSavein Platform Design: `.kiro/specs/insavein-platform/design.md`
