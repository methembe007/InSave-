#!/bin/bash
# Create application user for InSavein services

set -e

echo "Creating application user..."

# Wait for PostgreSQL to be ready
until pg_isready -U postgres; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# Create application user
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create application user if not exists
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'insavein_user') THEN
            CREATE ROLE insavein_user WITH LOGIN PASSWORD 'insavein_password';
        END IF;
    END
    \$\$;

    -- Grant necessary permissions
    GRANT ALL PRIVILEGES ON DATABASE ${POSTGRES_DB} TO insavein_user;
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO insavein_user;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO insavein_user;
    
    -- Grant default privileges for future tables
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO insavein_user;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO insavein_user;
EOSQL

echo "Application user 'insavein_user' created successfully!"
