#!/bin/bash

# Integration Test Runner Script for InSavein Platform
# This script sets up the test environment and runs integration tests

set -e

echo "========================================="
echo "InSavein Integration Test Runner"
echo "========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    print_error "docker-compose is not installed. Please install it and try again."
    exit 1
fi

# Navigate to integration-tests directory
cd "$(dirname "$0")"

print_status "Starting test environment..."
docker-compose -f docker-compose.test.yml up -d

print_status "Waiting for services to be healthy..."
sleep 10

# Wait for PostgreSQL to be ready
print_status "Checking PostgreSQL health..."
for i in {1..30}; do
    if docker-compose -f docker-compose.test.yml exec -T postgres-test pg_isready -U postgres > /dev/null 2>&1; then
        print_status "PostgreSQL is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "PostgreSQL failed to start"
        docker-compose -f docker-compose.test.yml logs postgres-test
        exit 1
    fi
    sleep 2
done

# Wait for services to be healthy
print_status "Waiting for microservices to be healthy..."
sleep 20

# Check service health
print_status "Checking service health..."
services=("auth-service-test:18080" "user-service-test:18081" "savings-service-test:18082" "budget-service-test:18083" "goal-service-test:18005")

for service in "${services[@]}"; do
    IFS=':' read -r name port <<< "$service"
    if curl -f -s "http://localhost:${port}/health" > /dev/null; then
        print_status "${name} is healthy"
    else
        print_warning "${name} health check failed (may not have /health endpoint)"
    fi
done

# Run tests
print_status "Running integration tests..."
echo ""

if [ "$1" == "coverage" ]; then
    print_status "Running tests with coverage..."
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    print_status "Coverage report generated: coverage.html"
elif [ -n "$1" ]; then
    print_status "Running specific test: $1"
    go test -v -run "$1" ./...
else
    go test -v ./...
fi

TEST_EXIT_CODE=$?

echo ""
if [ $TEST_EXIT_CODE -eq 0 ]; then
    print_status "All tests passed! ✓"
else
    print_error "Some tests failed! ✗"
fi

# Cleanup option
if [ "$2" == "cleanup" ] || [ "$1" == "cleanup" ]; then
    print_status "Cleaning up test environment..."
    docker-compose -f docker-compose.test.yml down -v
    print_status "Cleanup complete"
else
    print_warning "Test environment is still running. Use './run-tests.sh cleanup' to stop it."
    print_status "View logs: docker-compose -f docker-compose.test.yml logs [service-name]"
fi

exit $TEST_EXIT_CODE
