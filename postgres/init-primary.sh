#!/bin/bash
# Initialize PostgreSQL Primary Server with Replication User

set -e

echo "Initializing PostgreSQL Primary Server..."

# Wait for PostgreSQL to be ready
until pg_isready -U postgres; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# Create replication user
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create replication user if not exists
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'replicator') THEN
            CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD '${REPLICATION_PASSWORD:-replicator_password}';
        END IF;
    END
    \$\$;

    -- Grant necessary permissions
    GRANT CONNECT ON DATABASE ${POSTGRES_DB} TO replicator;
    
    -- Create replication slot for replica1
    SELECT pg_create_physical_replication_slot('replica1_slot');
    
    -- Create replication slot for replica2
    SELECT pg_create_physical_replication_slot('replica2_slot');
    
    -- Enable pg_stat_statements extension
    CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
    
    -- Create monitoring view for replication lag
    CREATE OR REPLACE VIEW replication_status AS
    SELECT 
        client_addr,
        state,
        sent_lsn,
        write_lsn,
        flush_lsn,
        replay_lsn,
        sync_state,
        EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp())) AS lag_seconds
    FROM pg_stat_replication;
    
    -- Grant access to monitoring view
    GRANT SELECT ON replication_status TO replicator;
EOSQL

echo "Primary server initialization complete!"
echo "Replication user 'replicator' created"
echo "Replication slots 'replica1_slot' and 'replica2_slot' created"
