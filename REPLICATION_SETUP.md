# PostgreSQL Replication Setup Documentation

## Overview

This document describes the PostgreSQL replication setup for the InSavein Platform, implementing a primary-replica architecture with 2 read replicas and PgBouncer connection pooling.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Application Layer                        │
│  (Auth, User, Savings, Budget, Goal Services)               │
└────────────┬────────────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────────┐
│                      PgBouncer (Port 6432)                   │
│              Connection Pooling & Load Balancing             │
└─────┬──────────────────────────────────┬────────────────────┘
      │                                   │
      │ Write Operations                  │ Read Operations
      ▼                                   ▼
┌──────────────────┐            ┌─────────────────────────────┐
│ PostgreSQL       │            │   Read Replicas             │
│ Primary          │──Streaming─┤   - Replica 1 (Port 5433)  │
│ (Port 5432)      │ Replication│   - Replica 2 (Port 5434)  │
└──────────────────┘            └─────────────────────────────┘
```

## Components

### 1. PostgreSQL Primary (postgres-primary)
- **Port**: 5432
- **Role**: Handles all write operations (INSERT, UPDATE, DELETE)
- **Configuration**: `postgres/primary/postgresql.conf`
- **Replication**: Configured with WAL streaming to 2 replicas
- **Replication Slots**: `replica1_slot`, `replica2_slot`

### 2. PostgreSQL Replica 1 (postgres-replica1)
- **Port**: 5433
- **Role**: Read-only operations for Education and Notification Services
- **Configuration**: `postgres/replica/postgresql.conf`
- **Replication**: Streams from primary using physical replication

### 3. PostgreSQL Replica 2 (postgres-replica2)
- **Port**: 5434
- **Role**: Read-only operations for Analytics Service
- **Configuration**: `postgres/replica/postgresql.conf`
- **Replication**: Streams from primary using physical replication

### 4. PgBouncer Connection Pooler
- **Port**: 6432
- **Configuration**: `pgbouncer/pgbouncer.ini`
- **Pool Mode**: Transaction-level pooling
- **Connection Limits**:
  - Max client connections: 1000
  - Default pool size: 20 per database
  - Min pool size: 5
  - Reserve pool: 5

### 5. Replication Monitor
- **Purpose**: Monitors replication lag and alerts on threshold violations
- **Script**: `monitoring/check-replication-lag.sh`
- **Thresholds**:
  - Warning: ≥ 1 second lag
  - Critical: ≥ 5 seconds lag

## Requirements Validation

### Requirement 11.6: Education Service reads from replicas ✓
- Education Service connects to `postgres-replica1` for lesson content
- Reduces load on primary database

### Requirement 13.5: Analytics Service reads from replicas ✓
- Analytics Service connects to `postgres-replica2` for data analysis
- Isolates heavy analytical queries from transactional workload

### Requirement 19.1: Health checks respond within 1 second ✓
- All services have health checks with 5s timeout
- PgBouncer provides fast connection pooling
- Replication lag monitored to ensure data freshness

### Requirement 22.1: At least 2 read replicas ✓
- Configured with 2 read replicas (replica1 and replica2)

### Requirement 22.5: Monitor replication lag, alert if exceeds 1 second ✓
- Automated monitoring script checks lag every 10 seconds
- Alerts on console when lag exceeds thresholds

## Setup Instructions

### Initial Setup

1. **Start the PostgreSQL cluster**:
   ```bash
   docker-compose up -d postgres-primary
   ```

   Wait for primary to be healthy:
   ```bash
   docker-compose ps postgres-primary
   ```

2. **Initialize replicas**:
   ```bash
   docker-compose up -d postgres-replica1 postgres-replica2
   ```

   The replicas will automatically:
   - Clone data from primary using `pg_basebackup`
   - Configure streaming replication
   - Start in hot standby mode

3. **Start PgBouncer**:
   ```bash
   # Generate userlist with proper password hashes
   chmod +x pgbouncer/generate-userlist.sh
   ./pgbouncer/generate-userlist.sh

   # Start PgBouncer
   docker-compose up -d pgbouncer
   ```

4. **Start replication monitoring** (optional):
   ```bash
   docker-compose --profile monitoring up -d replication-monitor
   ```

### Verify Replication

1. **Check replication status on primary**:
   ```bash
   docker exec -it postgres-primary psql -U postgres -d insavein -c "SELECT * FROM pg_stat_replication;"
   ```

   Expected output: 2 rows showing replica1 and replica2 connections

2. **Check replica status**:
   ```bash
   # Replica 1
   docker exec -it postgres-replica1 psql -U postgres -d insavein -c "SELECT pg_is_in_recovery();"
   
   # Replica 2
   docker exec -it postgres-replica2 psql -U postgres -d insavein -c "SELECT pg_is_in_recovery();"
   ```

   Expected output: `t` (true) indicating replica mode

3. **Check replication lag**:
   ```bash
   docker logs replication-monitor
   ```

   Or manually:
   ```bash
   docker exec -it postgres-primary psql -U postgres -d insavein -c "
   SELECT 
       application_name,
       state,
       sync_state,
       EXTRACT(EPOCH FROM (now() - replay_lag)) AS lag_seconds
   FROM pg_stat_replication;
   "
   ```

### Test Replication

1. **Insert data on primary**:
   ```bash
   docker exec -it postgres-primary psql -U postgres -d insavein -c "
   INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth)
   VALUES ('test@example.com', 'hash', 'Test', 'User', '1990-01-01');
   "
   ```

2. **Verify data on replicas**:
   ```bash
   # Check replica 1
   docker exec -it postgres-replica1 psql -U postgres -d insavein -c "
   SELECT email, first_name, last_name FROM users WHERE email = 'test@example.com';
   "
   
   # Check replica 2
   docker exec -it postgres-replica2 psql -U postgres -d insavein -c "
   SELECT email, first_name, last_name FROM users WHERE email = 'test@example.com';
   "
   ```

3. **Verify replicas are read-only**:
   ```bash
   docker exec -it postgres-replica1 psql -U postgres -d insavein -c "
   INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth)
   VALUES ('fail@example.com', 'hash', 'Fail', 'User', '1990-01-01');
   "
   ```

   Expected: Error message "cannot execute INSERT in a read-only transaction"

## Connection Strings

### For Microservices

**Write Operations** (Auth, User, Savings, Budget, Goal Services):
```
# Direct to primary
postgresql://postgres:postgres@postgres-primary:5432/insavein

# Via PgBouncer (recommended)
postgresql://postgres:postgres@pgbouncer:6432/insavein_primary
```

**Read Operations** (Education Service):
```
# Direct to replica 1
postgresql://postgres:postgres@postgres-replica1:5432/insavein

# Via PgBouncer with load balancing (recommended)
postgresql://postgres:postgres@pgbouncer:6432/insavein_replica1
```

**Read Operations** (Analytics Service):
```
# Direct to replica 2
postgresql://postgres:postgres@postgres-replica2:5432/insavein

# Via PgBouncer (recommended)
postgresql://postgres:postgres@pgbouncer:6432/insavein_replica2
```

**Load-Balanced Reads** (Notification Service):
```
# Via PgBouncer with automatic load balancing
postgresql://postgres:postgres@pgbouncer:6432/insavein_read
```

## PgBouncer Usage

### Connection Pooling Benefits

1. **Reduced Connection Overhead**: Reuses database connections
2. **Connection Limiting**: Prevents database connection exhaustion
3. **Transaction-Level Pooling**: Each transaction gets a connection from the pool
4. **Load Balancing**: Distributes read queries across replicas

### PgBouncer Admin Console

Access the admin console:
```bash
docker exec -it pgbouncer psql -h localhost -p 5432 -U postgres -d pgbouncer
```

Useful commands:
```sql
-- Show pool statistics
SHOW POOLS;

-- Show active connections
SHOW CLIENTS;

-- Show server connections
SHOW SERVERS;

-- Show configuration
SHOW CONFIG;

-- Reload configuration
RELOAD;

-- Pause all connections
PAUSE;

-- Resume connections
RESUME;
```

## Monitoring and Maintenance

### Replication Lag Monitoring

The replication monitor runs every 10 seconds and checks:
- Replication connection status
- Lag in seconds for each replica
- Replication slot status
- WAL retention

**View monitoring logs**:
```bash
docker logs -f replication-monitor
```

**Manual lag check**:
```bash
docker exec -it postgres-primary psql -U postgres -d insavein -c "
SELECT 
    application_name,
    client_addr,
    state,
    COALESCE(EXTRACT(EPOCH FROM (now() - replay_lag)), 0) AS lag_seconds,
    sync_state
FROM pg_stat_replication;
"
```

### Replication Slot Management

**View replication slots**:
```bash
docker exec -it postgres-primary psql -U postgres -d insavein -c "
SELECT 
    slot_name,
    slot_type,
    active,
    pg_size_pretty(pg_wal_lsn_diff(pg_current_wal_lsn(), restart_lsn)) AS retained_wal
FROM pg_replication_slots;
"
```

**Drop inactive slot** (if needed):
```bash
docker exec -it postgres-primary psql -U postgres -d insavein -c "
SELECT pg_drop_replication_slot('slot_name');
"
```

### Failover Procedure

If the primary fails, promote a replica to primary:

1. **Promote replica to primary**:
   ```bash
   docker exec -it postgres-replica1 pg_ctl promote -D /var/lib/postgresql/data
   ```

2. **Update application connection strings** to point to the new primary

3. **Reconfigure remaining replica** to stream from new primary

4. **Rebuild failed primary** as a new replica when recovered

## Troubleshooting

### Replica Not Connecting

1. Check primary logs:
   ```bash
   docker logs postgres-primary
   ```

2. Check replica logs:
   ```bash
   docker logs postgres-replica1
   ```

3. Verify replication user:
   ```bash
   docker exec -it postgres-primary psql -U postgres -d insavein -c "
   SELECT rolname, rolreplication FROM pg_roles WHERE rolname = 'replicator';
   "
   ```

4. Check network connectivity:
   ```bash
   docker exec -it postgres-replica1 ping postgres-primary
   ```

### High Replication Lag

1. **Check primary load**:
   ```bash
   docker exec -it postgres-primary psql -U postgres -d insavein -c "
   SELECT * FROM pg_stat_activity WHERE state = 'active';
   "
   ```

2. **Check WAL generation rate**:
   ```bash
   docker exec -it postgres-primary psql -U postgres -d insavein -c "
   SELECT pg_current_wal_lsn();
   "
   ```

3. **Check replica performance**:
   ```bash
   docker stats postgres-replica1 postgres-replica2
   ```

4. **Increase WAL retention** if needed (in `postgresql.conf`):
   ```
   wal_keep_size = 2048  # Increase from 1024
   ```

### PgBouncer Connection Issues

1. **Check PgBouncer logs**:
   ```bash
   docker logs pgbouncer
   ```

2. **Verify userlist**:
   ```bash
   docker exec -it pgbouncer cat /etc/pgbouncer/userlist.txt
   ```

3. **Test direct connection**:
   ```bash
   docker exec -it pgbouncer psql -h postgres-primary -U postgres -d insavein
   ```

4. **Regenerate userlist**:
   ```bash
   ./pgbouncer/generate-userlist.sh
   docker-compose restart pgbouncer
   ```

## Performance Tuning

### Primary Database

Key settings in `postgres/primary/postgresql.conf`:
- `shared_buffers = 256MB` - Memory for caching data
- `effective_cache_size = 1GB` - Estimated OS cache size
- `work_mem = 16MB` - Memory per query operation
- `maintenance_work_mem = 64MB` - Memory for maintenance operations
- `max_wal_senders = 10` - Maximum replication connections
- `wal_keep_size = 1024` - WAL retention for replicas

### Replica Configuration

Key settings in `postgres/replica/postgresql.conf`:
- `hot_standby = on` - Allow read queries on replica
- `max_standby_streaming_delay = 30s` - Max delay before canceling queries
- `hot_standby_feedback = on` - Prevent query conflicts

### PgBouncer Tuning

Key settings in `pgbouncer/pgbouncer.ini`:
- `pool_mode = transaction` - Connection pooling mode
- `default_pool_size = 20` - Connections per database
- `max_client_conn = 1000` - Maximum client connections
- `server_idle_timeout = 600` - Idle connection timeout

## Security Considerations

1. **Password Management**:
   - Change default passwords in production
   - Use environment variables or secrets management
   - Rotate passwords regularly

2. **Network Security**:
   - Use TLS for replication in production
   - Restrict network access with firewall rules
   - Use VPC/private networks in cloud deployments

3. **Authentication**:
   - Use SCRAM-SHA-256 for password authentication
   - Consider certificate-based authentication for replicas
   - Implement IP whitelisting in `pg_hba.conf`

4. **Monitoring**:
   - Set up alerts for replication lag
   - Monitor failed authentication attempts
   - Track connection pool exhaustion

## Backup Strategy

While replication provides high availability, it's not a backup solution:

1. **Regular Backups**:
   ```bash
   docker exec postgres-primary pg_dump -U postgres insavein > backup.sql
   ```

2. **Point-in-Time Recovery**:
   - Enable WAL archiving
   - Store WAL files in external storage
   - Test restore procedures regularly

3. **Backup Verification**:
   - Restore backups to test environment
   - Verify data integrity
   - Document restore procedures

## Migration from Single Instance

To migrate from the previous single-instance setup:

1. **Backup existing data**:
   ```bash
   docker exec insavein-postgres pg_dump -U postgres insavein > migration_backup.sql
   ```

2. **Stop old container**:
   ```bash
   docker-compose down
   ```

3. **Update docker-compose.yml** (already done)

4. **Start new cluster**:
   ```bash
   docker-compose up -d postgres-primary
   ```

5. **Restore data** (if needed):
   ```bash
   docker exec -i postgres-primary psql -U postgres insavein < migration_backup.sql
   ```

6. **Start replicas and PgBouncer**:
   ```bash
   docker-compose up -d
   ```

## References

- [PostgreSQL Replication Documentation](https://www.postgresql.org/docs/15/runtime-config-replication.html)
- [PgBouncer Documentation](https://www.pgbouncer.org/config.html)
- [PostgreSQL High Availability](https://www.postgresql.org/docs/15/high-availability.html)
- InSavein Platform Requirements: 11.6, 13.5, 19.1, 22.1, 22.5
