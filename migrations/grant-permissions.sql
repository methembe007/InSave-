-- Grant permissions to insavein_user for the insavein_db database
-- Run this as the postgres superuser

\c insavein_db

-- Grant all privileges on the database
GRANT ALL PRIVILEGES ON DATABASE insavein TO insavein_user;

-- Grant all privileges on the public schema
GRANT ALL PRIVILEGES ON SCHEMA public TO insavein_user;

-- Grant create privilege on the public schema
GRANT CREATE ON SCHEMA public TO insavein_user;

-- Grant usage on the public schema
GRANT USAGE ON SCHEMA public TO insavein_user;

-- Grant all privileges on all tables in public schema
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO insavein_user;

-- Grant all privileges on all sequences in public schema
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO insavein_user;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO insavein_user;

-- Set default privileges for future sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO insavein_user;

-- Make insavein_user the owner of the public schema (optional but recommended)
ALTER SCHEMA public OWNER TO insavein_user;

-- Verify permissions
\du insavein_user
\dn+ public
