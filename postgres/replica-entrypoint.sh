#!/bin/bash
set -e

# Environment variables
REPLICA_NAME=${REPLICA_NAME:-replica1}
PRIMARY_HOST=${PRIMARY_HOST:-postgres-primary}
PRIMARY_PORT=${PRIMARY_PORT:-5432}
REPLICATION_USER=${REPLICATION_USER:-replicator}
REPLICATION_PASSWORD=${REPLICATION_PASSWORD:-replicator_password}
PGDATA=${PGDATA:-/var/lib/postgresql/data}

echo "Starting replica initialization for ${REPLICA_NAME}..."

# Wait for primary to be ready
echo "Waiting for primary database at ${PRIMARY_HOST}:${PRIMARY_PORT}..."
until PGPASSWORD="${REPLICATION_PASSWORD}" pg_isready -h "${PRIMARY_HOST}" -p "${PRIMARY_PORT}" -U "${REPLICATION_USER}" 2>/dev/null; do
    echo "Primary is unavailable - sleeping"
    sleep 2
done
echo "Primary database is ready!"

# Check if data directory needs initialization
if [ ! -s "${PGDATA}/PG_VERSION" ]; then
    echo "Data directory is empty or uninitialized. Setting up replication..."
    
    # Clean the data directory
    rm -rf ${PGDATA}/*
    rm -rf ${PGDATA}/.[!.]*
    
    # Use pg_basebackup to clone from primary
    echo "Running pg_basebackup from ${PRIMARY_HOST}:${PRIMARY_PORT}..."
    PGPASSWORD="${REPLICATION_PASSWORD}" pg_basebackup \
        -h "${PRIMARY_HOST}" \
        -p "${PRIMARY_PORT}" \
        -U "${REPLICATION_USER}" \
        -D "${PGDATA}" \
        -Fp \
        -Xs \
        -P \
        -R \
        -S "${REPLICA_NAME}_slot"
    
    echo "Base backup complete!"
    
    # Ensure standby.signal exists
    touch ${PGDATA}/standby.signal
    
    # Configure replica connection in postgresql.auto.conf
    cat >> ${PGDATA}/postgresql.auto.conf <<EOF
primary_conninfo = 'host=${PRIMARY_HOST} port=${PRIMARY_PORT} user=${REPLICATION_USER} password=${REPLICATION_PASSWORD} application_name=${REPLICA_NAME}'
primary_slot_name = '${REPLICA_NAME}_slot'
EOF
    
    # Set proper permissions
    chmod 0700 ${PGDATA}
    chown -R postgres:postgres ${PGDATA}
    
    echo "Replica configuration complete!"
else
    echo "Data directory already initialized."
fi

# Start PostgreSQL
echo "Starting PostgreSQL in standby mode..."
exec postgres -c config_file=/etc/postgresql/postgresql.conf
