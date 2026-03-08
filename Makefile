# InSavein Platform Makefile
# Provides convenient commands for database migrations and development tasks

.PHONY: help migrate-up migrate-down migrate-create migrate-version migrate-status db-setup db-reset

# Default database URL (override with environment variable)
DATABASE_URL ?= postgresql://postgres:postgres@localhost:5432/insavein?sslmode=disable
MIGRATIONS_DIR = migrations

help: ## Show this help message
	@echo "InSavein Platform - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Environment Variables:"
	@echo "  DATABASE_URL    PostgreSQL connection string"
	@echo "                  Default: $(DATABASE_URL)"

migrate-up: ## Apply all pending migrations
	@echo "Applying migrations..."
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) up
	@echo "✓ Migrations applied successfully"

migrate-up-1: ## Apply next migration
	@echo "Applying next migration..."
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) up 1
	@echo "✓ Migration applied successfully"

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) down 1
	@echo "✓ Migration rolled back successfully"

migrate-down-all: ## Rollback all migrations (DESTRUCTIVE)
	@echo "WARNING: This will rollback ALL migrations!"
	@read -p "Are you sure? (yes/no): " confirm && [ "$$confirm" = "yes" ]
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) down
	@echo "✓ All migrations rolled back"

migrate-create: ## Create new migration (usage: make migrate-create name=add_indexes)
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required"; \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
	@echo "✓ Migration files created"

migrate-version: ## Show current migration version
	@echo "Current migration version:"
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) version

migrate-status: ## Show migration status
	@echo "Migration Status:"
	@echo "Database URL: $(DATABASE_URL)"
	@echo "Migrations Directory: $(MIGRATIONS_DIR)"
	@echo ""
	@echo "Migration files:"
	@ls -1 $(MIGRATIONS_DIR)/*.up.sql 2>/dev/null | wc -l | xargs echo "  Up migrations:"
	@ls -1 $(MIGRATIONS_DIR)/*.down.sql 2>/dev/null | wc -l | xargs echo "  Down migrations:"
	@echo ""
	@echo "Current version:"
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) version || true

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "Error: version parameter is required"; \
		echo "Usage: make migrate-force version=VERSION"; \
		exit 1; \
	fi
	@echo "WARNING: Forcing migration version to $(version)"
	@read -p "Are you sure? This should only be used to fix dirty state (yes/no): " confirm && [ "$$confirm" = "yes" ]
	migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) force $(version)
	@echo "✓ Version forced to $(version)"

db-setup: ## Create database and run migrations
	@echo "Setting up database..."
	createdb insavein || echo "Database may already exist"
	@$(MAKE) migrate-up
	@echo "✓ Database setup complete"

db-reset: ## Drop and recreate database (DESTRUCTIVE)
	@echo "WARNING: This will destroy all data in the database!"
	@read -p "Are you sure? (yes/no): " confirm && [ "$$confirm" = "yes" ]
	dropdb insavein || true
	createdb insavein
	@$(MAKE) migrate-up
	@echo "✓ Database reset complete"

db-shell: ## Open PostgreSQL shell
	psql "$(DATABASE_URL)"

db-dump: ## Dump database schema
	@echo "Dumping database schema..."
	pg_dump --schema-only "$(DATABASE_URL)" > schema_dump.sql
	@echo "✓ Schema dumped to schema_dump.sql"

db-backup: ## Create database backup
	@echo "Creating database backup..."
	@mkdir -p backups
	pg_dump "$(DATABASE_URL)" > backups/backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "✓ Backup created in backups/"

check-migrate: ## Check if golang-migrate is installed
	@which migrate > /dev/null || (echo "Error: golang-migrate is not installed" && \
		echo "Install it with:" && \
		echo "  macOS: brew install golang-migrate" && \
		echo "  Linux: curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && sudo mv migrate /usr/local/bin/" && \
		exit 1)
	@echo "✓ golang-migrate is installed"

validate-migrations: ## Validate migration files
	@echo "Validating migration files..."
	@for up_file in $(MIGRATIONS_DIR)/*.up.sql; do \
		base=$$(basename $$up_file .up.sql); \
		down_file="$(MIGRATIONS_DIR)/$${base}.down.sql"; \
		if [ ! -f "$$down_file" ]; then \
			echo "✗ Missing down migration for: $$up_file"; \
			exit 1; \
		fi; \
	done
	@echo "✓ All migration files are valid"

.DEFAULT_GOAL := help

# PostgreSQL Replication Management Commands

.PHONY: replication-up replication-down replication-status replication-test pgbouncer-setup pgbouncer-stats monitor-start monitor-logs

replication-up: ## Start PostgreSQL cluster with replication (primary + 2 replicas + PgBouncer)
	@echo "Starting PostgreSQL cluster with replication..."
	@docker-compose up -d postgres-primary
	@echo "Waiting for primary to be ready..."
	@sleep 5
	@docker-compose up -d postgres-replica1 postgres-replica2
	@echo "Waiting for replicas to be ready..."
	@sleep 10
	@echo "Setting up PgBouncer..."
	@chmod +x pgbouncer/generate-userlist.sh
	@./pgbouncer/generate-userlist.sh || echo "Note: Run 'make pgbouncer-setup' after primary is fully initialized"
	@docker-compose up -d pgbouncer
	@echo ""
	@echo "✓ PostgreSQL cluster is running!"
	@echo "  Primary:   localhost:5432"
	@echo "  Replica 1: localhost:5433"
	@echo "  Replica 2: localhost:5434"
	@echo "  PgBouncer: localhost:6432"
	@echo ""
	@echo "Run 'make replication-status' to check cluster health"

replication-down: ## Stop PostgreSQL replication cluster
	@echo "Stopping PostgreSQL cluster..."
	@docker-compose down
	@echo "✓ Cluster stopped"

replication-status: ## Check replication status and lag
	@echo "=== Replication Status ==="
	@docker exec postgres-primary psql -U postgres -d insavein -c "SELECT application_name, client_addr, state, sync_state, COALESCE(EXTRACT(EPOCH FROM (now() - replay_lag)), 0) AS lag_seconds FROM pg_stat_replication;" || echo "Primary not ready"
	@echo ""
	@echo "=== Replication Slots ==="
	@docker exec postgres-primary psql -U postgres -d insavein -c "SELECT slot_name, slot_type, active, pg_size_pretty(pg_wal_lsn_diff(pg_current_wal_lsn(), restart_lsn)) AS retained_wal FROM pg_replication_slots;" || echo "Primary not ready"
	@echo ""
	@echo "=== Container Health ==="
	@docker-compose ps postgres-primary postgres-replica1 postgres-replica2 pgbouncer

replication-test: ## Test replication by inserting and verifying data
	@echo "Testing replication..."
	@echo "1. Inserting test data on primary..."
	@docker exec postgres-primary psql -U postgres -d insavein -c "INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth) VALUES ('replication-test@example.com', 'test_hash', 'Replication', 'Test', '1990-01-01') ON CONFLICT (email) DO NOTHING;"
	@echo ""
	@echo "2. Waiting for replication (2 seconds)..."
	@sleep 2
	@echo ""
	@echo "3. Checking data on Replica 1..."
	@docker exec postgres-replica1 psql -U postgres -d insavein -c "SELECT email, first_name, last_name FROM users WHERE email = 'replication-test@example.com';"
	@echo ""
	@echo "4. Checking data on Replica 2..."
	@docker exec postgres-replica2 psql -U postgres -d insavein -c "SELECT email, first_name, last_name FROM users WHERE email = 'replication-test@example.com';"
	@echo ""
	@echo "5. Verifying replicas are read-only..."
	@docker exec postgres-replica1 psql -U postgres -d insavein -c "INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth) VALUES ('should-fail@example.com', 'hash', 'Should', 'Fail', '1990-01-01');" 2>&1 | grep -q "read-only" && echo "✓ Replica 1 is read-only" || echo "✗ Replica 1 is NOT read-only"
	@docker exec postgres-replica2 psql -U postgres -d insavein -c "INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth) VALUES ('should-fail@example.com', 'hash', 'Should', 'Fail', '1990-01-01');" 2>&1 | grep -q "read-only" && echo "✓ Replica 2 is read-only" || echo "✗ Replica 2 is NOT read-only"
	@echo ""
	@echo "✓ Replication test complete!"

pgbouncer-setup: ## Generate PgBouncer userlist with password hashes
	@echo "Generating PgBouncer userlist..."
	@chmod +x pgbouncer/generate-userlist.sh
	@./pgbouncer/generate-userlist.sh
	@docker-compose restart pgbouncer
	@echo "✓ PgBouncer userlist updated and service restarted"

pgbouncer-stats: ## Show PgBouncer connection pool statistics
	@echo "=== PgBouncer Pool Statistics ==="
	@docker exec pgbouncer psql -h localhost -p 5432 -U postgres -d pgbouncer -c "SHOW POOLS;"
	@echo ""
	@echo "=== PgBouncer Client Connections ==="
	@docker exec pgbouncer psql -h localhost -p 5432 -U postgres -d pgbouncer -c "SHOW CLIENTS;"
	@echo ""
	@echo "=== PgBouncer Server Connections ==="
	@docker exec pgbouncer psql -h localhost -p 5432 -U postgres -d pgbouncer -c "SHOW SERVERS;"

monitor-start: ## Start replication lag monitoring
	@echo "Starting replication lag monitor..."
	@chmod +x monitoring/check-replication-lag.sh
	@docker-compose --profile monitoring up -d replication-monitor
	@echo "✓ Monitor started. View logs with 'make monitor-logs'"

monitor-logs: ## Show replication monitor logs
	@docker logs -f replication-monitor

replication-logs: ## Show logs from all replication components
	@docker-compose logs -f postgres-primary postgres-replica1 postgres-replica2 pgbouncer

replication-clean: ## Remove all replication volumes (WARNING: deletes data)
	@echo "WARNING: This will delete all database data from primary and replicas!"
	@read -p "Are you sure? (yes/no): " confirm && [ "$confirm" = "yes" ]
	@docker-compose down -v
	@echo "✓ All replication volumes removed"
