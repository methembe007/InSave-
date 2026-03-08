# InSavein Platform - Database Setup Guide

This guide walks you through setting up the PostgreSQL database for the InSavein Platform.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Quick Start with Docker](#quick-start-with-docker)
- [Manual PostgreSQL Setup](#manual-postgresql-setup)
- [Running Migrations](#running-migrations)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software
- **PostgreSQL 12+** (or Docker for containerized setup)
- **golang-migrate** CLI tool
- **Make** (optional, for convenience commands)

### Install golang-migrate

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
chmod +x /usr/local/bin/migrate
```

**Windows:**
```powershell
scoop install migrate
```

Or download from: https://github.com/golang-migrate/migrate/releases

## Quick Start with Docker

The fastest way to get started is using Docker Compose:

### 1. Start PostgreSQL
```bash
docker-compose up -d postgres
```

This will:
- Start PostgreSQL 15 on port 5432
- Create database `insavein`
- Set up user `postgres` with password `postgres`

### 2. Wait for PostgreSQL to be ready
```bash
docker-compose ps
```

Wait until the postgres service shows as "healthy".

### 3. Run migrations
```bash
# Using Make (recommended)
make migrate-up

# Or using the migration script
./migrations/migrate.sh up

# Or using migrate directly
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/insavein?sslmode=disable"
migrate -database $DATABASE_URL -path migrations up
```

### 4. Verify setup
```bash
make migrate-status
```

### Optional: Start pgAdmin
If you want a GUI for database management:
```bash
docker-compose --profile tools up -d pgadmin
```

Access pgAdmin at: http://localhost:5050
- Email: admin@insavein.com
- Password: admin

## Manual PostgreSQL Setup

If you prefer to install PostgreSQL directly:

### 1. Install PostgreSQL

**macOS:**
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install postgresql-15 postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

**Windows:**
Download and install from: https://www.postgresql.org/download/windows/

### 2. Create Database
```bash
# Switch to postgres user (Linux/macOS)
sudo -u postgres psql

# Or connect directly (if configured)
psql -U postgres
```

In the PostgreSQL shell:
```sql
CREATE DATABASE insavein;
CREATE USER insavein_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE insavein TO insavein_user;
\q
```

### 3. Set Database URL
```bash
# Linux/macOS
export DATABASE_URL="postgresql://insavein_user:your_secure_password@localhost:5432/insavein?sslmode=disable"

# Windows PowerShell
$env:DATABASE_URL="postgresql://insavein_user:your_secure_password@localhost:5432/insavein?sslmode=disable"
```

### 4. Run Migrations
```bash
make migrate-up
```

## Running Migrations

### Using Make (Recommended)

```bash
# Apply all migrations
make migrate-up

# Apply next migration only
make migrate-up-1

# Rollback last migration
make migrate-down

# Show current version
make migrate-version

# Show migration status
make migrate-status

# Create new migration
make migrate-create name=add_new_feature

# Validate migration files
make validate-migrations
```

### Using Migration Script

**Linux/macOS:**
```bash
./migrations/migrate.sh up
./migrations/migrate.sh down 1
./migrations/migrate.sh version
./migrations/migrate.sh status
```

**Windows:**
```cmd
migrations\migrate.bat up
migrations\migrate.bat down 1
migrations\migrate.bat version
migrations\migrate.bat status
```

### Using migrate CLI Directly

```bash
# Set database URL
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/insavein?sslmode=disable"

# Apply all migrations
migrate -database $DATABASE_URL -path migrations up

# Rollback last migration
migrate -database $DATABASE_URL -path migrations down 1

# Show version
migrate -database $DATABASE_URL -path migrations version
```

## Verification

### 1. Check Migration Version
```bash
make migrate-version
```

Expected output: `10` (current latest version)

### 2. Verify Tables
```bash
psql $DATABASE_URL -c "\dt"
```

Expected tables:
- users
- savings_transactions (parent table)
- savings_transactions_YYYY_MM (partitions)
- budgets
- budget_categories
- spending_transactions (parent table)
- spending_transactions_YYYY_MM (partitions)
- goals
- goal_milestones
- notifications
- lessons
- education_progress

### 3. Check Partitions
```sql
-- Connect to database
psql $DATABASE_URL

-- List all partitions
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'savings_transactions_%' 
   OR tablename LIKE 'spending_transactions_%'
ORDER BY tablename;
```

### 4. Verify Indexes
```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;
```

### 5. Test Insert
```sql
-- Insert test user
INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth)
VALUES ('test@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIiIkYiIkY', 'Test', 'User', '1990-01-01')
RETURNING id;

-- Insert test savings transaction (use the returned user ID)
INSERT INTO savings_transactions (user_id, amount, currency, description, category)
VALUES ('USER_ID_HERE', 10.50, 'USD', 'Test savings', 'general');

-- Verify
SELECT * FROM savings_transactions;

-- Clean up
DELETE FROM users WHERE email = 'test@example.com';
```

## Database Schema Overview

### Core Tables

**users**
- Stores user accounts with bcrypt password hashing (cost factor 12)
- JSONB preferences field for flexible user settings
- CASCADE deletion for all related data

**savings_transactions** (Partitioned)
- Partitioned by `created_at` (monthly)
- Stores all savings deposits
- DECIMAL(15,2) for precise amount storage

**budgets**
- Monthly budget plans
- Unique constraint on (user_id, month)
- Tracks total budget and total spent

**budget_categories**
- Category allocations within budgets
- Tracks allocated and spent amounts
- Color field for UI visualization

**spending_transactions** (Partitioned)
- Partitioned by `transaction_date` (monthly)
- Links to budgets and categories
- Merchant tracking for analytics

**goals**
- Long-term financial goals
- Status tracking (active, completed, paused)
- Progress calculation support

**goal_milestones**
- Intermediate checkpoints for goals
- Ordered milestones with completion tracking

**notifications**
- User notification system
- Read/unread status tracking
- Type-based categorization

**lessons**
- Financial education content
- Difficulty levels and categories
- JSONB resources field

**education_progress**
- Tracks user lesson completion
- Unique constraint on (user_id, lesson_id)

### Partitioning Strategy

Both transaction tables use monthly partitioning:
- **savings_transactions**: Partitioned by `created_at`
- **spending_transactions**: Partitioned by `transaction_date`

Partitions are pre-created for 2024-2027. To add new partitions:

```sql
-- Example: Add January 2028 partition for savings
CREATE TABLE savings_transactions_2028_01 PARTITION OF savings_transactions
    FOR VALUES FROM ('2028-01-01') TO ('2028-02-01');

-- Example: Add January 2028 partition for spending
CREATE TABLE spending_transactions_2028_01 PARTITION OF spending_transactions
    FOR VALUES FROM ('2028-01-01') TO ('2028-02-01');
```

## Troubleshooting

### Migration Failed

**Problem:** Migration fails with error

**Solution:**
1. Check the error message carefully
2. Verify database connection: `psql $DATABASE_URL`
3. Check migration version: `make migrate-version`
4. If in dirty state, force version: `make migrate-force version=X`
5. Re-run migration: `make migrate-up`

### Connection Refused

**Problem:** Cannot connect to PostgreSQL

**Solution:**
1. Check if PostgreSQL is running:
   ```bash
   # Docker
   docker-compose ps
   
   # System service
   sudo systemctl status postgresql
   ```
2. Verify port 5432 is not in use by another process
3. Check DATABASE_URL is correct
4. Verify firewall settings

### Permission Denied

**Problem:** Permission denied when creating tables

**Solution:**
1. Ensure user has proper privileges:
   ```sql
   GRANT ALL PRIVILEGES ON DATABASE insavein TO your_user;
   GRANT ALL PRIVILEGES ON SCHEMA public TO your_user;
   ```
2. Or use superuser for initial setup

### Dirty Database State

**Problem:** Migration shows "dirty" state

**Solution:**
1. Check what went wrong: `make migrate-version`
2. Manually inspect database state
3. Fix any partial changes
4. Force to correct version: `make migrate-force version=X`
5. Re-run migrations

### Partition Not Found

**Problem:** Insert fails with "no partition of relation found"

**Solution:**
1. Check the date of the transaction
2. Verify partition exists for that month:
   ```sql
   SELECT tablename FROM pg_tables 
   WHERE tablename LIKE 'savings_transactions_%';
   ```
3. Create missing partition if needed

## Production Considerations

### Security
- Use strong passwords (not default `postgres`)
- Enable SSL/TLS connections (`sslmode=require`)
- Restrict network access with firewall rules
- Use separate users with minimal privileges
- Enable database-level encryption at rest

### Performance
- Monitor partition sizes regularly
- Add indexes for new query patterns
- Run ANALYZE after bulk inserts
- Configure connection pooling (max 20 per service)
- Use read replicas for analytics queries

### Backup
- Set up automated daily backups
- Test backup restoration monthly
- Store backups in separate geographic location
- Retain backups for 30 days minimum

### Monitoring
- Monitor replication lag (should be < 1 second)
- Track query performance
- Alert on failed migrations
- Monitor disk space usage

## Next Steps

After setting up the database:

1. **Configure Backend Services**: Update microservice configurations with DATABASE_URL
2. **Set up Connection Pooling**: Configure PgBouncer or similar
3. **Enable Replication**: Set up read replicas for production
4. **Configure Backups**: Implement automated backup strategy
5. **Set up Monitoring**: Configure Prometheus metrics and Grafana dashboards

## Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Partitioning Guide](https://www.postgresql.org/docs/current/ddl-partitioning.html)
- InSavein Requirements: `.kiro/specs/insavein-platform/requirements.md`
- InSavein Design: `.kiro/specs/insavein-platform/design.md`
