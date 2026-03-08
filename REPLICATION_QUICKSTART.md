# PostgreSQL Replication Quick Start Guide

## Overview

This guide will help you quickly set up and verify the PostgreSQL replication cluster for the InSavein Platform.

## Prerequisites

- Docker and Docker Compose installed
- Make utility installed
- Ports 5432, 5433, 5434, and 6432 available

## Quick Setup (5 minutes)

### Step 1: Start the Replication Cluster

```bash
make replication-up
```

This command will:
1. Start the PostgreSQL primary server (port 5432)
2. Start 2 read replicas (ports 5433, 5434)
3. Configure streaming replication
4. Set up PgBouncer connection pooler (port 6432)

**Expected output:**
```
✓ PostgreSQL cluster is running!
  Primary:   localhost:5432
  Replica 1: localhost:5433
  Replica 2: localhost:5434
  PgBouncer: localhost:6432
```

### Step 2: Run Database Migrations

```bash
make migrate-up
```

This applies all database schema migrations to the primary database. The changes will automatically replicate to the read replicas.

### Step 3: Verify Replication

```bash
make replication-status
```

**Expected output:**
- 2 active replication connections (replica1 and replica2)
- Lag should be < 1 second
- Both replicas should show "streaming" state

### Step 4: Test Replication

```bash
make replication-test
```

This will:
1. Insert test data on the primary
2. Verify data appears on both replicas
3. Confirm replicas are read-only

**Expected output:**
```
✓ Replica 1 is read-only
✓ Replica 2 is read-only
✓ Replication test complete!
```

## Connection Strings

### For Application Services

**Write Operations** (Auth, User, Savings, Budget, Goal Services):
```
postgresql://postgres:postgres@localhost:5432/insavein
```

**Read Operations - Education Service** (via Replica 1):
```
postgresql://postgres:postgres@localhost:5433/insavein
```

**Read Operations - Analytics Service** (via Replica 2):
```
postgresql://postgres:postgres@localhost:5434/insavein
```

**Via PgBouncer** (recommended for production):
```
# Write operations
postgresql://postgres:postgres@localhost:6432/insavein_primary

# Read operations (load balanced)
postgresql://postgres:postgres@localhost:6432/insavein_read
```

## Monitoring

### Check Replication Lag

```bash
make replication-status
```

### Start Continuous Monitoring

```bash
make monitor-start
make monitor-logs
```

The monitor checks replication lag every 10 seconds and alerts if:
- Lag exceeds 1 second (WARNING)
- Lag exceeds 5 seconds (CRITICAL)

### View PgBouncer Statistics

```bash
make pgbouncer-stats
```

## Common Commands

```bash
# Start cluster
make replication-up

# Stop cluster
make replication-down

# Check status
make replication-status

# Test replication
make replication-test

# View logs
make replication-logs

# Monitor lag
make monitor-start
make monitor-logs

# PgBouncer stats
make pgbouncer-stats
```

## Troubleshooting

### Replicas Not Connecting

1. Check primary is running:
   ```bash
   docker ps | grep postgres-primary
   ```

2. Check replica logs:
   ```bash
   docker logs postgres-replica1
   docker logs postgres-replica2
   ```

3. Restart replicas:
   ```bash
   docker-compose restart postgres-replica1 postgres-replica2
   ```

### High Replication Lag

1. Check primary load:
   ```bash
   docker stats postgres-primary
   ```

2. Check network connectivity:
   ```bash
   docker exec postgres-replica1 ping postgres-primary
   ```

3. View detailed replication status:
   ```bash
   make replication-status
   ```

### PgBouncer Connection Issues

1. Regenerate userlist:
   ```bash
   make pgbouncer-setup
   ```

2. Check PgBouncer logs:
   ```bash
   docker logs pgbouncer
   ```

3. Test direct connection:
   ```bash
   docker exec pgbouncer psql -h postgres-primary -U postgres -d insavein
   ```

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                   Application Services                       │
│  Auth | User | Savings | Budget | Goal | Education | Analytics│
└────────────┬────────────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────────┐
│                  PgBouncer (Port 6432)                       │
│         Connection Pooling & Load Balancing                  │
└─────┬──────────────────────────────────┬────────────────────┘
      │                                   │
      │ Writes                            │ Reads
      ▼                                   ▼
┌──────────────────┐            ┌─────────────────────────────┐
│ PostgreSQL       │            │   Read Replicas             │
│ Primary          │──Streaming─┤   - Replica 1 (Port 5433)  │
│ (Port 5432)      │ Replication│   - Replica 2 (Port 5434)  │
└──────────────────┘            └─────────────────────────────┘
```

## Requirements Validation

✅ **Requirement 11.6**: Education Service reads from database replicas
- Education Service connects to Replica 1 (port 5433)

✅ **Requirement 13.5**: Analytics Service reads from database replicas
- Analytics Service connects to Replica 2 (port 5434)

✅ **Requirement 19.1**: System responds to health checks within 1 second
- All services have health checks configured
- PgBouncer provides fast connection pooling

✅ **Requirement 22.1**: At least 2 read replicas
- Configured with 2 read replicas

✅ **Requirement 22.5**: Monitor replication lag, alert if exceeds 1 second
- Automated monitoring with configurable thresholds

## Next Steps

1. **Configure Application Services**: Update service connection strings to use appropriate endpoints
2. **Set Up Monitoring**: Integrate replication lag monitoring with your alerting system
3. **Test Failover**: Practice promoting a replica to primary
4. **Configure Backups**: Set up automated backup procedures
5. **Security Hardening**: Change default passwords and configure TLS

## Additional Resources

- [Full Replication Setup Documentation](./REPLICATION_SETUP.md)
- [PostgreSQL Replication Documentation](https://www.postgresql.org/docs/15/runtime-config-replication.html)
- [PgBouncer Documentation](https://www.pgbouncer.org/)

## Support

For issues or questions:
1. Check the [Troubleshooting](#troubleshooting) section
2. Review logs: `make replication-logs`
3. Check replication status: `make replication-status`
4. Refer to the full documentation: [REPLICATION_SETUP.md](./REPLICATION_SETUP.md)
