#!/bin/bash

# Dockerfile Verification Script
# Validates that all required Dockerfiles exist and have correct structure

set -e

echo "=========================================="
echo "Dockerfile Verification"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Services that should have Dockerfiles
services=(
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

all_valid=true

# Check each service
for service in "${services[@]}"; do
    dockerfile="$service/Dockerfile"
    
    if [ ! -f "$dockerfile" ]; then
        print_error "$service: Dockerfile not found"
        all_valid=false
        continue
    fi
    
    # Check for required Dockerfile elements
    has_from=$(grep -c "^FROM" "$dockerfile" || true)
    has_workdir=$(grep -c "WORKDIR" "$dockerfile" || true)
    has_copy=$(grep -c "COPY" "$dockerfile" || true)
    has_expose=$(grep -c "EXPOSE" "$dockerfile" || true)
    has_cmd=$(grep -c "CMD" "$dockerfile" || true)
    
    if [ "$has_from" -lt 1 ]; then
        print_error "$service: Missing FROM instruction"
        all_valid=false
        continue
    fi
    
    if [ "$has_workdir" -lt 1 ]; then
        print_error "$service: Missing WORKDIR instruction"
        all_valid=false
        continue
    fi
    
    if [ "$has_copy" -lt 1 ]; then
        print_error "$service: Missing COPY instruction"
        all_valid=false
        continue
    fi
    
    if [ "$has_expose" -lt 1 ]; then
        print_error "$service: Missing EXPOSE instruction"
        all_valid=false
        continue
    fi
    
    if [ "$has_cmd" -lt 1 ]; then
        print_error "$service: Missing CMD instruction"
        all_valid=false
        continue
    fi
    
    # Check for multi-stage build (Go services)
    if [[ "$service" != "frontend" ]]; then
        has_builder=$(grep -c "AS builder" "$dockerfile" || true)
        if [ "$has_builder" -lt 1 ]; then
            print_error "$service: Missing multi-stage build (AS builder)"
            all_valid=false
            continue
        fi
        
        # Check for non-root user
        has_user=$(grep -c "USER" "$dockerfile" || true)
        if [ "$has_user" -lt 1 ]; then
            print_error "$service: Missing USER instruction (should run as non-root)"
            all_valid=false
            continue
        fi
        
        # Check for health check
        has_healthcheck=$(grep -c "HEALTHCHECK" "$dockerfile" || true)
        if [ "$has_healthcheck" -lt 1 ]; then
            print_error "$service: Missing HEALTHCHECK instruction"
            all_valid=false
            continue
        fi
    fi
    
    print_success "$service: Dockerfile is valid"
done

echo ""
echo "=========================================="
echo "Docker Compose Verification"
echo "=========================================="
echo ""

if [ ! -f "docker-compose.yml" ]; then
    print_error "docker-compose.yml not found"
    all_valid=false
else
    # Check if all services are defined in docker-compose.yml
    for service in "${services[@]}"; do
        if grep -q "^  $service:" docker-compose.yml; then
            print_success "$service: Defined in docker-compose.yml"
        else
            print_error "$service: Not defined in docker-compose.yml"
            all_valid=false
        fi
    done
    
    # Check for database services
    if grep -q "postgres-primary:" docker-compose.yml; then
        print_success "PostgreSQL primary: Defined in docker-compose.yml"
    else
        print_error "PostgreSQL primary: Not defined in docker-compose.yml"
        all_valid=false
    fi
    
    if grep -q "postgres-replica1:" docker-compose.yml; then
        print_success "PostgreSQL replica1: Defined in docker-compose.yml"
    else
        print_error "PostgreSQL replica1: Not defined in docker-compose.yml"
        all_valid=false
    fi
    
    if grep -q "postgres-replica2:" docker-compose.yml; then
        print_success "PostgreSQL replica2: Defined in docker-compose.yml"
    else
        print_error "PostgreSQL replica2: Not defined in docker-compose.yml"
        all_valid=false
    fi
fi

echo ""
echo "=========================================="
echo "Summary"
echo "=========================================="
echo ""

if [ "$all_valid" = true ]; then
    print_success "All Dockerfiles and docker-compose.yml are valid!"
    echo ""
    print_info "Next steps:"
    echo "  1. Start Docker Desktop"
    echo "  2. Run: ./docker-build-test.sh (Linux/Mac) or docker-build-test.bat (Windows)"
    echo "  3. Or run: docker-compose up -d"
    exit 0
else
    print_error "Some validation checks failed. Please fix the issues above."
    exit 1
fi
