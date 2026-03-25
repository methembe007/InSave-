#!/bin/bash

# Partition Maintenance Script
# This script provides convenient commands for managing database partitions

set -e

# Configuration
DB_NAME="${DB_NAME:-insavein}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to execute SQL
execute_sql() {
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$1"
}

# Function to execute SQL file
execute_sql_file() {
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$1"
}

# Command: setup
setup() {
    print_info "Setting up partition management..."
    
    print_info "Creating partition management functions..."
    execute_sql_file "create_monthly_partitions.sql"
    
    print_info "Setting up automatic partition creation..."
    execute_sql_file "auto_create_partitions.sql"
    
    print_info "Setting up archival functions..."
    execute_sql_file "archive_old_partitions.sql"
    
    print_success "Partition management setup complete!"
}

# Command: create
create() {
    local start_date="${1:-$(date -d '6 months ago' +%Y-%m-01)}"
    local end_date="${2:-$(date -d '6 months' +%Y-%m-01)}"
    
    print_info "Creating partitions from $start_date to $end_date..."
    
    execute_sql "SELECT create_partitions_for_range('savings_transactions', '$start_date', '$end_date');"
    execute_sql "SELECT create_partitions_for_range('spending_transactions', '$start_date', '$end_date');"
    
    print_success "Partitions created successfully!"
}

# Command: maintain
maintain() {
    print_info "Running partition maintenance..."
    
    execute_sql "SELECT maintain_partitions();"
    
    print_success "Partition maintenance complete!"
}

# Command: status
status() {
    print_info "Partition Status Report"
    echo ""
    
    print_info "Active Partitions:"
    execute_sql "
        SELECT 
            parent.relname AS table_name,
            child.relname AS partition_name,
            TO_DATE(SUBSTRING(child.relname FROM '(\d{4}_\d{2})$'), 'YYYY_MM') AS partition_month,
            pg_size_pretty(pg_total_relation_size(child.oid)) AS size
        FROM pg_inherits
        JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
        JOIN pg_class child ON pg_inherits.inhrelid = child.oid
        WHERE parent.relname IN ('savings_transactions', 'spending_transactions')
        ORDER BY parent.relname, partition_month DESC
        LIMIT 20;
    "
    
    echo ""
    print_info "Total Size by Table:"
    execute_sql "
        SELECT 
            tablename,
            pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size
        FROM pg_tables
        WHERE tablename IN ('savings_transactions', 'spending_transactions')
        ORDER BY tablename;
    "
}

# Command: archive
archive() {
    local retention_months="${1:-12}"
    
    print_warning "This will archive partitions older than $retention_months months"
    read -p "Continue? (y/N) " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Archiving old partitions..."
        
        execute_sql "SELECT archive_workflow('savings_transactions', $retention_months);"
        execute_sql "SELECT archive_workflow('spending_transactions', $retention_months);"
        
        print_success "Archival complete!"
    else
        print_info "Archival cancelled"
    fi
}

# Command: list-detached
list_detached() {
    print_info "Detached Partitions:"
    
    execute_sql "
        SELECT 
            tablename,
            pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
        FROM pg_tables
        WHERE schemaname = 'public'
        AND (tablename LIKE 'savings_transactions_%' OR tablename LIKE 'spending_transactions_%')
        AND tablename NOT IN (
            SELECT child.relname
            FROM pg_inherits
            JOIN pg_class child ON pg_inherits.inhrelid = child.oid
        )
        ORDER BY tablename;
    "
}

# Command: verify
verify() {
    print_info "Verifying partition health..."
    
    print_info "Checking for missing future partitions..."
    local next_month=$(date -d 'next month' +%Y-%m-01)
    local month_after=$(date -d '2 months' +%Y-%m-01)
    
    local missing=0
    
    for table in "savings_transactions" "spending_transactions"; do
        for month in "$next_month" "$month_after"; do
            local partition_name="${table}_$(date -d "$month" +%Y_%m)"
            local exists=$(execute_sql "SELECT COUNT(*) FROM pg_tables WHERE tablename = '$partition_name';" | grep -o '[0-9]*' | head -1)
            
            if [ "$exists" -eq 0 ]; then
                print_warning "Missing partition: $partition_name"
                missing=$((missing + 1))
            fi
        done
    done
    
    if [ $missing -eq 0 ]; then
        print_success "All required partitions exist!"
    else
        print_warning "Found $missing missing partitions. Run 'maintain' to create them."
    fi
}

# Command: help
show_help() {
    cat << EOF
Partition Maintenance Script

Usage: $0 <command> [options]

Commands:
    setup                   Initial setup of partition management functions
    create [start] [end]    Create partitions for date range (default: 6 months ago to 6 months ahead)
    maintain                Run daily maintenance (create future partitions)
    status                  Show partition status and sizes
    archive [months]        Archive partitions older than N months (default: 12)
    list-detached           List detached partitions
    verify                  Verify partition health
    help                    Show this help message

Environment Variables:
    DB_NAME                 Database name (default: insavein)
    DB_USER                 Database user (default: postgres)
    DB_HOST                 Database host (default: localhost)
    DB_PORT                 Database port (default: 5432)

Examples:
    # Initial setup
    $0 setup

    # Create partitions for 2024
    $0 create 2024-01-01 2024-12-31

    # Daily maintenance
    $0 maintain

    # Check status
    $0 status

    # Archive partitions older than 18 months
    $0 archive 18

    # Verify partition health
    $0 verify

EOF
}

# Main script
main() {
    local command="${1:-help}"
    
    case "$command" in
        setup)
            setup
            ;;
        create)
            create "$2" "$3"
            ;;
        maintain)
            maintain
            ;;
        status)
            status
            ;;
        archive)
            archive "$2"
            ;;
        list-detached)
            list_detached
            ;;
        verify)
            verify
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
