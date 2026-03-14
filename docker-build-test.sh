#!/bin/bash

# Docker Build and Test Script for InSavein Platform
# This script builds all Docker images and tests the deployment

set -e  # Exit on error

echo "=========================================="
echo "InSavein Platform - Docker Build & Test"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Check if Docker is running
echo "Checking Docker status..."
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker Desktop and try again."
    exit 1
fi
print_success "Docker is running"
echo ""

# Check if Docker Compose is available
echo "Checking Docker Compose..."
if ! docker-compose version > /dev/null 2>&1; then
    print_error "Docker Compose is not installed or not in PATH"
    exit 1
fi
print_success "Docker Compose is available"
echo ""

# Build all services
echo "=========================================="
echo "Building Docker Images"
echo "=========================================="
echo ""

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

for service in "${services[@]}"; do
    echo "Building $service..."
    if docker-compose build "$service"; then
        print_success "$service built successfully"
    else
        print_error "$service build failed"
        exit 1
    fi
    echo ""
done

print_success "All services built successfully!"
echo ""

# Start services
echo "=========================================="
echo "Starting Services"
echo "=========================================="
echo ""

print_info "Starting PostgreSQL databases..."
docker-compose up -d postgres-primary postgres-replica1 postgres-replica2

print_info "Waiting for databases to be healthy (30 seconds)..."
sleep 30

# Check database health
if docker-compose ps postgres-primary | grep -q "healthy"; then
    print_success "PostgreSQL primary is healthy"
else
    print_error "PostgreSQL primary is not healthy"
    docker-compose logs postgres-primary
    exit 1
fi

print_info "Starting microservices..."
docker-compose up -d auth-service user-service savings-service budget-service goal-service education-service notification-service analytics-service

print_info "Waiting for microservices to be healthy (20 seconds)..."
sleep 20

print_info "Starting frontend..."
docker-compose up -d frontend

print_info "Waiting for frontend to be healthy (15 seconds)..."
sleep 15

echo ""
echo "=========================================="
echo "Service Health Checks"
echo "=========================================="
echo ""

# Check health of all services
check_health() {
    local service=$1
    local port=$2
    
    if curl -f -s "http://localhost:$port/health" > /dev/null 2>&1; then
        print_success "$service (port $port) is healthy"
        return 0
    else
        print_error "$service (port $port) is not responding"
        return 1
    fi
}

# Check microservices
check_health "Auth Service" 8080
check_health "User Service" 8081
check_health "Savings Service" 8082
check_health "Budget Service" 8083
check_health "Goal Service" 8005
check_health "Education Service" 8085
check_health "Notification Service" 8086
check_health "Analytics Service" 8008

# Check frontend
if curl -f -s "http://localhost:3000" > /dev/null 2>&1; then
    print_success "Frontend (port 3000) is responding"
else
    print_error "Frontend (port 3000) is not responding"
fi

echo ""
echo "=========================================="
echo "Container Status"
echo "=========================================="
echo ""

docker-compose ps

echo ""
echo "=========================================="
echo "Basic Functionality Tests"
echo "=========================================="
echo ""

# Test 1: Register a user
print_info "Test 1: User Registration"
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test@example.com",
        "password": "testpassword123",
        "first_name": "Test",
        "last_name": "User",
        "date_of_birth": "1990-01-01"
    }')

if echo "$REGISTER_RESPONSE" | grep -q "access_token"; then
    print_success "User registration successful"
    ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
else
    print_error "User registration failed"
    echo "Response: $REGISTER_RESPONSE"
fi

# Test 2: Login
print_info "Test 2: User Login"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test@example.com",
        "password": "testpassword123"
    }')

if echo "$LOGIN_RESPONSE" | grep -q "access_token"; then
    print_success "User login successful"
    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
else
    print_error "User login failed"
    echo "Response: $LOGIN_RESPONSE"
fi

# Test 3: Get user profile
if [ -n "$ACCESS_TOKEN" ]; then
    print_info "Test 3: Get User Profile"
    PROFILE_RESPONSE=$(curl -s -X GET http://localhost:8081/api/user/profile \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if echo "$PROFILE_RESPONSE" | grep -q "email"; then
        print_success "User profile retrieval successful"
    else
        print_error "User profile retrieval failed"
        echo "Response: $PROFILE_RESPONSE"
    fi
fi

echo ""
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo ""

print_success "Docker build and deployment test completed!"
echo ""
print_info "Services are running at:"
echo "  - Frontend:              http://localhost:3000"
echo "  - Auth Service:          http://localhost:8080"
echo "  - User Service:          http://localhost:8081"
echo "  - Savings Service:       http://localhost:8082"
echo "  - Budget Service:        http://localhost:8083"
echo "  - Goal Service:          http://localhost:8005"
echo "  - Education Service:     http://localhost:8085"
echo "  - Notification Service:  http://localhost:8086"
echo "  - Analytics Service:     http://localhost:8008"
echo ""
print_info "To view logs: docker-compose logs -f [service-name]"
print_info "To stop all services: docker-compose down"
print_info "To stop and remove volumes: docker-compose down -v"
echo ""
