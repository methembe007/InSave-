#!/bin/bash

# InSavein Platform Database Migration Script
# This script provides convenient commands for managing database migrations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
MIGRATIONS_DIR="migrations"
DATABASE_URL="${DATABASE_URL:-postgresql://postgres:postgres@localhost:5432/insavein?sslmode=disable}"

# Print colored message
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Check if migrate is installed
check_migrate() {
    if ! command -v migrate &> /dev/null; then
        print_message "$RED" "Error: golang-migrate is not installed"
        echo "Install it with:"
        echo "  macOS: brew install golang-migrate"
        echo "  Linux: curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && sudo mv migrate /usr/local/bin/"
        echo "  Windows: scoop install migrate"
        exit 1
    fi
}

# Show usage
usage() {
    cat << EOF
Usage: $0 [command] [options]

Commands:
    up [N]              Apply all or N pending migrations
    down [N]            Rollback all or N migrations
    create <name>       Create a new migration file
    version             Show current migration version
    force <version>     Force set migration version (use with caution)
    status              Show migration status
    validate            Validate migration files
    help                Show this help message

Environment Variables:
    DATABASE_URL        PostgreSQL connection string
                        Default: postgresql://postgres:postgres@localhost:5432/insavein?sslmode=disable

Examples:
    $0 up                    # Apply all pending migrations
    $0 up 1                  # Apply next migration
    $0 down 1                # Rollback last migration
    $0 create add_indexes    # Create new migration files
    $0 version               # Show current version
    $0 status                # Show migration status

EOF
}

# Apply migrations up
migrate_up() {
    local steps=${1:-}
    print_message "$GREEN" "Applying migrations..."
    if [ -z "$steps" ]; then
        migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" up
    else
        migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" up "$steps"
    fi
    print_message "$GREEN" "✓ Migrations applied successfully"
}

# Rollback migrations
migrate_down() {
    local steps=${1:-}
    print_message "$YELLOW" "Rolling back migrations..."
    if [ -z "$steps" ]; then
        read -p "Are you sure you want to rollback ALL migrations? (yes/no): " confirm
        if [ "$confirm" != "yes" ]; then
            print_message "$YELLOW" "Rollback cancelled"
            exit 0
        fi
        migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" down
    else
        migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" down "$steps"
    fi
    print_message "$GREEN" "✓ Migrations rolled back successfully"
}

# Create new migration
create_migration() {
    local name=$1
    if [ -z "$name" ]; then
        print_message "$RED" "Error: Migration name is required"
        echo "Usage: $0 create <migration_name>"
        exit 1
    fi
    print_message "$GREEN" "Creating migration: $name"
    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$name"
    print_message "$GREEN" "✓ Migration files created"
}

# Show current version
show_version() {
    print_message "$GREEN" "Current migration version:"
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" version
}

# Force version
force_version() {
    local version=$1
    if [ -z "$version" ]; then
        print_message "$RED" "Error: Version number is required"
        echo "Usage: $0 force <version>"
        exit 1
    fi
    print_message "$YELLOW" "Forcing migration version to: $version"
    read -p "Are you sure? This should only be used to fix dirty state (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        print_message "$YELLOW" "Force cancelled"
        exit 0
    fi
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" force "$version"
    print_message "$GREEN" "✓ Version forced to $version"
}

# Show migration status
show_status() {
    print_message "$GREEN" "Migration Status:"
    echo "Database URL: $DATABASE_URL"
    echo "Migrations Directory: $MIGRATIONS_DIR"
    echo ""
    
    # Count migration files
    local up_files=$(ls -1 "$MIGRATIONS_DIR"/*.up.sql 2>/dev/null | wc -l)
    local down_files=$(ls -1 "$MIGRATIONS_DIR"/*.down.sql 2>/dev/null | wc -l)
    echo "Migration files found: $up_files up, $down_files down"
    
    # Show current version
    echo ""
    echo "Current version:"
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" version || true
}

# Validate migrations
validate_migrations() {
    print_message "$GREEN" "Validating migration files..."
    
    # Check for matching up/down files
    local valid=true
    for up_file in "$MIGRATIONS_DIR"/*.up.sql; do
        local base=$(basename "$up_file" .up.sql)
        local down_file="$MIGRATIONS_DIR/${base}.down.sql"
        if [ ! -f "$down_file" ]; then
            print_message "$RED" "✗ Missing down migration for: $up_file"
            valid=false
        fi
    done
    
    if [ "$valid" = true ]; then
        print_message "$GREEN" "✓ All migration files are valid"
    else
        print_message "$RED" "✗ Validation failed"
        exit 1
    fi
}

# Main script
main() {
    check_migrate
    
    local command=${1:-help}
    shift || true
    
    case "$command" in
        up)
            migrate_up "$@"
            ;;
        down)
            migrate_down "$@"
            ;;
        create)
            create_migration "$@"
            ;;
        version)
            show_version
            ;;
        force)
            force_version "$@"
            ;;
        status)
            show_status
            ;;
        validate)
            validate_migrations
            ;;
        help|--help|-h)
            usage
            ;;
        *)
            print_message "$RED" "Error: Unknown command '$command'"
            echo ""
            usage
            exit 1
            ;;
    esac
}

main "$@"
