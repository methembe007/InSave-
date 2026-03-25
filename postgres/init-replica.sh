#!/bin/bash
# Initialize PostgreSQL Replica Server

set -e

REPLICA_NAME=${REPLICA_NAME:-replica1}
PRIMARY_HOST=${PRIMARY_HOST:-postgres-primary}
PRIMARY_PORT=${PRIMARY_PORT:-5432}
REPLICATION_USER=${REPLICATION_USER:-replicator}
REPLICATION_PASSWORD=${REPLICATION_PASSWORD:-replicator_password}

echo "Initializing PostgreSQL Replica: ${REPLICA_NAME}..."

# Wait for primary to be ready
echo "Waiting for primary database to be ready..."
until PGPASSWORD="${REPLICATION_PASSWORD}" pg_isready -h "${PRIMARY_HOST}" -p "${PRIMARY_PORT}" -U "${REPLICATION_USER}"; do
    echo "Primary is unavailable - sleeping"
    sleep 2
done

echo "Primary database is ready!"

# Check if this is first initialization (no PG_VERSION file)
if [ ! -f "${PGDATA}/PG_VERSION" ]; then
    echo "First initialization detected. Setting up replication from primary..."
    
    # Remove any existing data
    rm -rf ${PGDATA}/*
    
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
    
    # Create standby.signal file (PostgreSQL 12+)
    touch ${PGDATA}/standby.signal
    
    # Configure replica connection
    cat >> ${PGDATA}/postgresql.auto.conf <<EOF
primary_conninfo = 'host=${PRIMARY_HOST} port=${PRIMARY_PORT} user=${REPLICATION_USER} password=${REPLICATION_PASSWORD} application_name=${REPLICA_NAME}'
primary_slot_name = '${REPLICA_NAME}_slot'
EOF
    
    # Set proper permissions
    chmod 0700 ${PGDATA}
    
    echo "Replica configuration complete!"
else
    echo "Data directory already initialized. Skipping setup."
fi
