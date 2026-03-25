#!/bin/bash

# Performance Test Runner Script for InSavein Platform
# This script runs k6 load tests and generates reports

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
BASE_URL="${BASE_URL:-http://localhost:8080}"
TEST_TYPE="${1:-normal}"
OUTPUT_DIR="./results"

echo -e "${GREEN}=== InSavein Performance Testing ===${NC}\n"

# Check if k6 is installed
if ! command -v k6 &> /dev/null; then
    echo -e "${RED}Error: k6 is not installed${NC}"
    echo "Install k6 from: https://k6.io/docs/getting-started/installation/"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Get timestamp for results
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

echo "Configuration:"
echo "  Base URL: $BASE_URL"
echo "  Test Type: $TEST_TYPE"
echo "  Output Directory: $OUTPUT_DIR"
echo ""

# Function to run a test
run_test() {
    local test_file=$1
    local test_name=$2
    local output_file="${OUTPUT_DIR}/${test_name}_${TIMESTAMP}.json"
    
    echo -e "${YELLOW}Running $test_name...${NC}"
    
    if k6 run \
        -e BASE_URL="$BASE_URL" \
        --out json="$output_file" \
        "$test_file"; then
        echo -e "${GREEN}✓ $test_name completed successfully${NC}"
        echo "  Results saved to: $output_file"
        return 0
    else
        echo -e "${RED}✗ $test_name failed${NC}"
        return 1
    fi
}

# Run the specified test
case "$TEST_TYPE" in
    normal)
        run_test "normal-load.js" "normal-load"
        ;;
    peak)
        run_test "peak-load.js" "peak-load"
        ;;
    stress)
        run_test "stress-test.js" "stress-test"
        ;;
    all)
        echo -e "${YELLOW}Running all tests sequentially...${NC}\n"
        run_test "normal-load.js" "normal-load"
        echo ""
        sleep 60  # Wait 1 minute between tests
        run_test "peak-load.js" "peak-load"
        echo ""
        sleep 60
        run_test "stress-test.js" "stress-test"
        ;;
    *)
        echo -e "${RED}Error: Unknown test type '$TEST_TYPE'${NC}"
        echo "Usage: $0 [normal|peak|stress|all]"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}=== Performance Testing Complete ===${NC}"
echo ""
echo "Results are available in: $OUTPUT_DIR"
echo ""
echo "To analyze results:"
echo "  - Review the JSON output files"
echo "  - Check p95/p99 latencies against targets (p95<500ms, p99<1000ms)"
echo "  - Verify error rate is below 0.1%"
echo "  - Monitor system resources during tests"
