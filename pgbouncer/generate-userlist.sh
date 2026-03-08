#!/bin/bash
# Generate PgBouncer userlist with proper password hashes

set -e

POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}
REPLICATION_PASSWORD=${REPLICATION_PASSWORD:-replicator_password}

echo "Generating PgBouncer userlist..."

# Function to generate SCRAM-SHA-256 hash
generate_scram_hash() {
    local username=$1
    local password=$2
    
    # Use PostgreSQL to generate the hash
    docker exec postgres-primary psql -U postgres -t -c \
        "SELECT concat('\"', rolname, '\" \"', rolpassword, '\"') FROM pg_authid WHERE rolname='${username}';" | tr -d ' '
}

# Wait for primary to be ready
echo "Waiting for primary database..."
until docker exec postgres-primary pg_isready -U postgres > /dev/null 2>&1; do
    sleep 2
done

# Generate userlist
{
    generate_scram_hash "postgres" "${POSTGRES_PASSWORD}"
    generate_scram_hash "replicator" "${REPLICATION_PASSWORD}"
} > pgbouncer/userlist.txt

echo "Userlist generated successfully!"
cat pgbouncer/userlist.txt
