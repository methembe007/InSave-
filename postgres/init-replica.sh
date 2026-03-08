#!/bin/bash
# Initialize PostgreSQL Replica Server

set -e

REPLICA_NAME=${REPLICA_NAME:-replica1}
PRIMARY_HOST=${PRIMARY_HOST:-postgres-primary}
PRIMARY_PORT=${PRIMARY_PORT:-5432}
REPLICATION_USER=${REPLICATION_USER:-replicator}
REPLICATION_PASSWORD=${REPLICATION_PASSWORD:-replicator_password}

echo "Initializing PostgreSQL Replica: ${REPLICA_NAME}..."

# Check if data directory is empty
if [ -z "$(ls -A /var/lib/postgresql/data)" ]; then
    echo "Data directory is empty. Setting up replication from primary..."
    
    # Remove any existing data
    rm -rf /var/lib/postgresql/data/*
    
    # Use pg_basebackup to clone from primary
    echo "Running pg_basebackup from ${PRIMARY_HOST}:${PRIMARY_PORT}..."
    PGPASSWORD="${REPLICATION_PASSWORD}" pg_basebackup \
        -h "${PRIMARY_HOST}" \
        -p "${PRIMARY_PORT}" \
        -U "${REPLICATION_USER}" \
        -D /var/lib/postgresql/data \
        -Fp \
        -Xs \
        -P \
        -R \
        -S "${REPLICA_NAME}_slot"
    
    echo "Base backup complete!"
    
    # Create standby.signal file (PostgreSQL 12+)
    touch /var/lib/postgresql/data/standby.signal
    
    # Configure replica connection
    cat >> /var/lib/postgresql/data/postgresql.auto.conf <<EOF
primary_conninfo = 'host=${PRIMARY_HOST} port=${PRIMARY_PORT} user=${REPLICATION_USER} password=${REPLICATION_PASSWORD} application_name=${REPLICA_NAME}'
primary_slot_name = '${REPLICA_NAME}_slot'
EOF
    
    echo "Replica configuration complete!"
else
    echo "Data directory already initialized. Skipping setup."
fi

echo "Starting replica server..."
