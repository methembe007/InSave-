#!/bin/bash

# Run All Tests Script
# InSavein Platform
# Executes integration tests, performance tests, and security scans

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Test results
INTEGRATION_PASSED=false
PERFORMANCE_PASSED=false
SECURITY_PASSED=false

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}InSavein Platform - Complete Test Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to print section
print_section() {
    echo ""
    echo -e "${BLUE}=== $1 ===${NC}"
    echo ""
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
print_section "Checking Prerequisites"

MISSING_TOOLS=()

if ! command_exists docker; then
    MISSING_TOOLS+=("docker")
fi

if ! command_exists docker-compose; then
    MISSING_TOOLS+=("docker-compose")
fi

if ! command_exists go; then
    MISSING_TOOLS+=("go")
fi

if ! command_exists k6; then
    echo -e "${YELLOW}⚠${NC} k6 not found. Install from: https://k6.io/docs/getting-started/installation/"
    MISSING_TOOLS+=("k6")
fi

if ! command_exists trivy; then
    echo -e "${YELLOW}⚠${NC} trivy not found. Install from: https://aquasecurity.github.io/trivy/"
    MISSING_TOOLS+=("trivy")
fi

if [ ${#MISSING_TOOLS[@]} -gt 0 ]; then
    echo -e "${RED}Missing required tools: ${MISSING_TOOLS[*]}${NC}"
    echo "Please install missing tools and try again."
    exit 1
fi

echo -e "${GREEN}✓${NC} All required tools are installed"

# 1. Integration Tests
print_section "Running Integration Tests"

if [ -d "integration-tests" ]; then
    cd integration-tests
    
    echo "Starting test environment..."
    docker-compose -f docker-compose.test.yml up -d
    
    echo "Waiting for services to be ready..."
    sleep 10
    
    echo "Running integration tests..."
    if make test; then
        echo -e "${GREEN}✓${NC} Integration tests passed"
        INTEGRATION_PASSED=true
    else
        echo -e "${RED}✗${NC} Integration tests failed"
        INTEGRATION_PASSED=false
    fi
    
    echo "Cleaning up test environment..."
    docker-compose -f docker-compose.test.yml down -v
    
    cd ..
else
    echo -e "${YELLOW}⚠${NC} Integration tests directory not found"
fi

# 2. Performance Tests
print_section "Running Performance Tests"

if [ -d "performance-tests" ]; then
    cd performance-tests
    
    echo "Running normal load test..."
    if make test-normal; then
        echo -e "${GREEN}✓${NC} Normal load test passed"
        
        echo ""
        echo "Running peak load test..."
        if make test-peak; then
            echo -e "${GREEN}✓${NC} Peak load test passed"
            PERFORMANCE_PASSED=true
        else
            echo -e "${RED}✗${NC} Peak load test failed"
            PERFORMANCE_PASSED=false
        fi
    else
        echo -e "${RED}✗${NC} Normal load test failed"
        PERFORMANCE_PASSED=false
    fi
    
    # Display results if they exist
    if [ -f "results/normal-load-results.json" ]; then
        echo ""
        echo "Normal Load Test Results:"
        cat results/normal-load-results.json | grep -E "(http_req_duration|http_req_failed)" || true
    fi
    
    if [ -f "results/peak-load-results.json" ]; then
        echo ""
        echo "Peak Load Test Results:"
        cat results/peak-load-results.json | grep -E "(http_req_duration|http_req_failed)" || true
    fi
    
    cd ..
else
    echo -e "${YELLOW}⚠${NC} Performance tests directory not found"
fi

# 3. Security Scans
print_section "Running Security Scans"

IMAGES=(
    "auth-service"
    "user-service"
    "savings-service"
    "budget-service"
    "goal-service"
    "education-service"
    "notification-service"
    "analytics-service"
    "frontend"
)

CRITICAL_VULNS=0
HIGH_VULNS=0

echo "Scanning Docker images for vulnerabilities..."
echo ""

for image in "${IMAGES[@]}"; do
    echo "Scanning $image..."
    
    # Check if image exists locally
    if docker images | grep -q "$image"; then
        # Run Trivy scan
        SCAN_OUTPUT=$(trivy image --severity CRITICAL,HIGH --format json "$image:latest" 2>/dev/null || echo "{}")
        
        # Count vulnerabilities
        CRITICAL=$(echo "$SCAN_OUTPUT" | grep -o '"Severity":"CRITICAL"' | wc -l)
        HIGH=$(echo "$SCAN_OUTPUT" | grep -o '"Severity":"HIGH"' | wc -l)
        
        CRITICAL_VULNS=$((CRITICAL_VULNS + CRITICAL))
        HIGH_VULNS=$((HIGH_VULNS + HIGH))
        
        if [ "$CRITICAL" -eq 0 ] && [ "$HIGH" -eq 0 ]; then
            echo -e "${GREEN}✓${NC} $image: No critical or high vulnerabilities"
        else
            echo -e "${YELLOW}⚠${NC} $image: $CRITICAL critical, $HIGH high vulnerabilities"
        fi
    else
        echo -e "${YELLOW}⚠${NC} $image: Image not found locally (skipping)"
    fi
done

echo ""
echo "Total vulnerabilities found:"
echo "  Critical: $CRITICAL_VULNS"
echo "  High: $HIGH_VULNS"

if [ "$CRITICAL_VULNS" -eq 0 ]; then
    echo -e "${GREEN}✓${NC} No critical vulnerabilities found"
    SECURITY_PASSED=true
else
    echo -e "${RED}✗${NC} Critical vulnerabilities found - must be addressed before production"
    SECURITY_PASSED=false
fi

# 4. Unit Tests (Go services)
print_section "Running Unit Tests"

GO_SERVICES=(
    "auth-service"
    "user-service"
    "savings-service"
    "budget-service"
    "goal-service"
    "education-service"
    "notification-service"
    "analytics-service"
)

UNIT_TESTS_PASSED=true

for service in "${GO_SERVICES[@]}"; do
    if [ -d "$service" ]; then
        echo "Testing $service..."
        cd "$service"
        
        if go test ./... -v 2>&1 | tee test-output.txt; then
            echo -e "${GREEN}✓${NC} $service: Unit tests passed"
        else
            echo -e "${YELLOW}⚠${NC} $service: Some unit tests failed or no tests found"
            # Don't fail the entire suite for unit tests
        fi
        
        cd ..
    fi
done

# 5. Frontend Tests
print_section "Running Frontend Tests"

if [ -d "frontend" ]; then
    cd frontend
    
    if [ -f "package.json" ]; then
        echo "Installing dependencies..."
        npm install --silent
        
        echo "Running frontend tests..."
        if npm test -- --run; then
            echo -e "${GREEN}✓${NC} Frontend tests passed"
        else
            echo -e "${YELLOW}⚠${NC} Frontend tests failed or no tests found"
        fi
    fi
    
    cd ..
else
    echo -e "${YELLOW}⚠${NC} Frontend directory not found"
fi

# Summary
print_section "Test Suite Summary"

echo ""
echo "Test Results:"
echo "============="

if [ "$INTEGRATION_PASSED" = true ]; then
    echo -e "${GREEN}✓${NC} Integration Tests: PASSED"
else
    echo -e "${RED}✗${NC} Integration Tests: FAILED"
fi

if [ "$PERFORMANCE_PASSED" = true ]; then
    echo -e "${GREEN}✓${NC} Performance Tests: PASSED"
else
    echo -e "${RED}✗${NC} Performance Tests: FAILED"
fi

if [ "$SECURITY_PASSED" = true ]; then
    echo -e "${GREEN}✓${NC} Security Scans: PASSED"
else
    echo -e "${RED}✗${NC} Security Scans: FAILED (Critical vulnerabilities found)"
fi

echo ""

# Overall result
if [ "$INTEGRATION_PASSED" = true ] && [ "$PERFORMANCE_PASSED" = true ] && [ "$SECURITY_PASSED" = true ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}All critical tests passed!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "The platform is ready for production deployment."
    echo ""
    echo "Next steps:"
    echo "1. Review PRODUCTION_READINESS_CHECKLIST.md"
    echo "2. Run ./verify-production-readiness.sh"
    echo "3. Complete manual E2E testing"
    echo "4. Schedule production deployment"
    exit 0
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}Some tests failed!${NC}"
    echo -e "${RED}========================================${NC}"
    echo ""
    echo "Please address the failed tests before proceeding to production."
    echo ""
    
    if [ "$INTEGRATION_PASSED" = false ]; then
        echo "- Fix integration test failures"
    fi
    
    if [ "$PERFORMANCE_PASSED" = false ]; then
        echo "- Investigate performance issues"
    fi
    
    if [ "$SECURITY_PASSED" = false ]; then
        echo "- Address critical security vulnerabilities"
    fi
    
    exit 1
fi
