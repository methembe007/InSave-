#!/bin/bash
# Verification script for PostgreSQL Replication Setup
# Task 1.2: Configure PostgreSQL replication setup

set -e

echo "=========================================="
echo "PostgreSQL Replication Setup Verification"
echo "Task 1.2 - InSavein Platform"
echo "=========================================="
echo ""

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
PASSED=0
FAILED=0

# Test function
test_check() {
    local test_name=$1
    local test_command=$2
    
    echo -n "Testing: $test_name... "
    
    if eval "$test_command" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ PASS${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC}"
        ((FAILED++))
        return 1
    fi
}

echo "=== File Structure Verification ==="
echo ""

# Check configuration files
test_check "Primary PostgreSQL config exists" "[ -f postgres/primary/postgresql.conf ]"
test_check "Primary pg_hba.conf exists" "[ -f postgres/primary/pg_hba.conf ]"
test_check "Primary init script exists" "[ -f postgres/init-primary.sh ]"
test_check "Replica PostgreSQL config exists" "[ -f postgres/replica/postgresql.conf ]"
test_check "Replica init script exists" "[ -f postgres/init-replica.sh ]"
test_check "PgBouncer config exists" "[ -f pgbouncer/pgbouncer.ini ]"
test_check "PgBouncer userlist exists" "[ -f pgbouncer/userlist.txt ]"
test_check "PgBouncer setup script exists" "[ -f pgbouncer/generate-userlist.sh ]"
test_check "Monitoring script exists" "[ -f monitoring/check-replication-lag.sh ]"
test_check "Docker Compose file exists" "[ -f docker-compose.yml ]"

echo ""
echo "=== Configuration Content Verification ==="
echo ""

# Check primary configuration
test_check "Primary config has wal_level=replica" "grep -q 'wal_level = replica' postgres/primary/postgresql.conf"
test_check "Primary config has max_wal_senders" "grep -q 'max_wal_senders' postgres/primary/postgresql.conf"
test_check "Primary config has replication slots" "grep -q 'max_replication_slots' postgres/primary/postgresql.conf"
test_check "Primary pg_hba allows replication" "grep -q 'replication.*replicator' postgres/primary/pg_hba.conf"

# Check replica configuration
test_check "Replica config has hot_standby=on" "grep -q 'hot_standby = on' postgres/replica/postgresql.conf"
test_check "Replica config has hot_standby_feedback" "grep -q 'hot_standby_feedback' postgres/replica/postgresql.conf"

# Check PgBouncer configuration
test_check "PgBouncer has primary database" "grep -q 'insavein_primary' pgbouncer/pgbouncer.ini"
test_check "PgBouncer has replica1 database" "grep -q 'insavein_replica1' pgbouncer/pgbouncer.ini"
test_check "PgBouncer has replica2 database" "grep -q 'insavein_replica2' pgbouncer/pgbouncer.ini"
test_check "PgBouncer has read load balancing" "grep -q 'insavein_read' pgbouncer/pgbouncer.ini"
test_check "PgBouncer pool mode is transaction" "grep -q 'pool_mode = transaction' pgbouncer/pgbouncer.ini"

echo ""
echo "=== Docker Compose Verification ==="
echo ""

# Check Docker Compose services
test_check "Docker Compose has postgres-primary" "grep -q 'postgres-primary:' docker-compose.yml"
test_check "Docker Compose has postgres-replica1" "grep -q 'postgres-replica1:' docker-compose.yml"
test_check "Docker Compose has postgres-replica2" "grep -q 'postgres-replica2:' docker-compose.yml"
test_check "Docker Compose has pgbouncer" "grep -q 'pgbouncer:' docker-compose.yml"
test_check "Docker Compose has replication-monitor" "grep -q 'replication-monitor:' docker-compose.yml"

# Check volumes
test_check "Docker Compose has primary volume" "grep -q 'postgres_primary_data:' docker-compose.yml"
test_check "Docker Compose has replica1 volume" "grep -q 'postgres_replica1_data:' docker-compose.yml"
test_check "Docker Compose has replica2 volume" "grep -q 'postgres_replica2_data:' docker-compose.yml"

# Check ports
test_check "Primary exposed on port 5432" "grep -q '5432:5432' docker-compose.yml"
test_check "Replica1 exposed on port 5433" "grep -q '5433:5432' docker-compose.yml"
test_check "Replica2 exposed on port 5434" "grep -q '5434:5432' docker-compose.yml"
test_check "PgBouncer exposed on port 6432" "grep -q '6432:5432' docker-compose.yml"

echo ""
echo "=== Documentation Verification ==="
echo ""

test_check "Replication setup documentation exists" "[ -f REPLICATION_SETUP.md ]"
test_check "Quick start guide exists" "[ -f REPLICATION_QUICKSTART.md ]"
test_check "Task completion summary exists" "[ -f TASK_1.2_COMPLETION_SUMMARY.md ]"

# Check documentation content
test_check "Setup doc mentions requirement 11.6" "grep -q '11.6' REPLICATION_SETUP.md"
test_check "Setup doc mentions requirement 13.5" "grep -q '13.5' REPLICATION_SETUP.md"
test_check "Setup doc mentions requirement 19.1" "grep -q '19.1' REPLICATION_SETUP.md"
test_check "Setup doc has connection strings" "grep -q 'Connection Strings' REPLICATION_SETUP.md"
test_check "Setup doc has troubleshooting" "grep -q 'Troubleshooting' REPLICATION_SETUP.md"

echo ""
echo "=== Makefile Commands Verification ==="
echo ""

test_check "Makefile has replication-up command" "grep -q 'replication-up:' Makefile"
test_check "Makefile has replication-down command" "grep -q 'replication-down:' Makefile"
test_check "Makefile has replication-status command" "grep -q 'replication-status:' Makefile"
test_check "Makefile has replication-test command" "grep -q 'replication-test:' Makefile"
test_check "Makefile has pgbouncer-setup command" "grep -q 'pgbouncer-setup:' Makefile"
test_check "Makefile has pgbouncer-stats command" "grep -q 'pgbouncer-stats:' Makefile"
test_check "Makefile has monitor-start command" "grep -q 'monitor-start:' Makefile"
test_check "Makefile has monitor-logs command" "grep -q 'monitor-logs:' Makefile"

echo ""
echo "=== Script Validation ==="
echo ""

# Check init scripts have required commands
test_check "Primary init creates replication user" "grep -q 'CREATE ROLE replicator' postgres/init-primary.sh"
test_check "Primary init creates replication slots" "grep -q 'pg_create_physical_replication_slot' postgres/init-primary.sh"
test_check "Replica init uses pg_basebackup" "grep -q 'pg_basebackup' postgres/init-replica.sh"
test_check "Replica init creates standby.signal" "grep -q 'standby.signal' postgres/init-replica.sh"

# Check monitoring script
test_check "Monitor checks replication status" "grep -q 'pg_stat_replication' monitoring/check-replication-lag.sh"
test_check "Monitor checks lag threshold" "grep -q 'LAG_THRESHOLD' monitoring/check-replication-lag.sh"
test_check "Monitor checks replication slots" "grep -q 'pg_replication_slots' monitoring/check-replication-lag.sh"

echo ""
echo "=========================================="
echo "Verification Summary"
echo "=========================================="
echo ""
echo -e "Tests Passed: ${GREEN}${PASSED}${NC}"
echo -e "Tests Failed: ${RED}${FAILED}${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All verification checks passed!${NC}"
    echo ""
    echo "Task 1.2 implementation is complete and verified."
    echo ""
    echo "Next steps:"
    echo "1. Start the cluster: make replication-up"
    echo "2. Run migrations: make migrate-up"
    echo "3. Verify replication: make replication-status"
    echo "4. Test replication: make replication-test"
    echo ""
    exit 0
else
    echo -e "${RED}✗ Some verification checks failed!${NC}"
    echo ""
    echo "Please review the failed checks above and fix any issues."
    echo ""
    exit 1
fi
