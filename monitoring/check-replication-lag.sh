#!/bin/sh
# Monitor PostgreSQL Replication Lag
# Alerts if lag exceeds 1 second (per requirement 22.5)

set -e

PRIMARY_HOST=${PRIMARY_HOST:-postgres-primary}
POSTGRES_USER=${POSTGRES_USER:-postgres}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}
LAG_THRESHOLD_SECONDS=${LAG_THRESHOLD_SECONDS:-1}
LAG_CRITICAL_SECONDS=${LAG_CRITICAL_SECONDS:-5}

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Replication Lag Monitor - $(date)"
echo "=========================================="

# Check replication status on primary
REPLICATION_STATUS=$(PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${PRIMARY_HOST}" -U "${POSTGRES_USER}" -d insavein -t -A -c "
SELECT 
    application_name,
    client_addr,
    state,
    sent_lsn,
    write_lsn,
    flush_lsn,
    replay_lsn,
    COALESCE(EXTRACT(EPOCH FROM (now() - replay_lag)), 0) AS lag_seconds,
    sync_state
FROM pg_stat_replication
ORDER BY application_name;
" 2>/dev/null)

if [ -z "$REPLICATION_STATUS" ]; then
    echo "${RED}ERROR: No replication connections found!${NC}"
    echo "Primary server may not be configured for replication or replicas are not connected."
    exit 1
fi

# Parse and display replication status
echo "$REPLICATION_STATUS" | while IFS='|' read -r app_name client_addr state sent_lsn write_lsn flush_lsn replay_lsn lag_seconds sync_state; do
    if [ -n "$app_name" ]; then
        echo ""
        echo "Replica: $app_name"
        echo "  Address: $client_addr"
        echo "  State: $state"
        echo "  Sync State: $sync_state"
        
        # Convert lag to float comparison
        lag_int=$(echo "$lag_seconds" | cut -d'.' -f1)
        
        if [ -z "$lag_int" ] || [ "$lag_int" = "" ]; then
            lag_int=0
        fi
        
        # Check lag thresholds
        if [ "$lag_int" -ge "$LAG_CRITICAL_SECONDS" ]; then
            echo "  ${RED}Lag: ${lag_seconds}s (CRITICAL - exceeds ${LAG_CRITICAL_SECONDS}s threshold!)${NC}"
            echo "  ${RED}ACTION REQUIRED: Replication lag is critically high!${NC}"
        elif [ "$lag_int" -ge "$LAG_THRESHOLD_SECONDS" ]; then
            echo "  ${YELLOW}Lag: ${lag_seconds}s (WARNING - exceeds ${LAG_THRESHOLD_SECONDS}s threshold)${NC}"
        else
            echo "  ${GREEN}Lag: ${lag_seconds}s (OK)${NC}"
        fi
        
        # Display LSN information
        echo "  LSN Info:"
        echo "    Sent: $sent_lsn"
        echo "    Write: $write_lsn"
        echo "    Flush: $flush_lsn"
        echo "    Replay: $replay_lsn"
    fi
done

# Check replication slots
echo ""
echo "=========================================="
echo "Replication Slots Status"
echo "=========================================="

SLOT_STATUS=$(PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${PRIMARY_HOST}" -U "${POSTGRES_USER}" -d insavein -t -A -c "
SELECT 
    slot_name,
    slot_type,
    active,
    COALESCE(pg_size_pretty(pg_wal_lsn_diff(pg_current_wal_lsn(), restart_lsn)), '0 bytes') AS retained_wal
FROM pg_replication_slots
ORDER BY slot_name;
" 2>/dev/null)

if [ -n "$SLOT_STATUS" ]; then
    echo "$SLOT_STATUS" | while IFS='|' read -r slot_name slot_type active retained_wal; do
        if [ -n "$slot_name" ]; then
            echo ""
            echo "Slot: $slot_name"
            echo "  Type: $slot_type"
            if [ "$active" = "t" ]; then
                echo "  Status: ${GREEN}Active${NC}"
            else
                echo "  Status: ${RED}Inactive${NC}"
            fi
            echo "  Retained WAL: $retained_wal"
        fi
    done
else
    echo "${YELLOW}No replication slots found${NC}"
fi

echo ""
echo "=========================================="
echo "Health Check Complete"
echo "=========================================="
