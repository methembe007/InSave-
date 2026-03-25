#!/bin/bash

# Validation script to check integration test setup

echo "========================================="
echo "Integration Test Setup Validation"
echo "========================================="
echo ""

# Check required files
echo "Checking required files..."
files=(
    "docker-compose.test.yml"
    "go.mod"
    "test-data/init-test-db.sql"
    "helpers/client.go"
    "helpers/types.go"
    "user_registration_test.go"
    "savings_flow_test.go"
    "budget_alert_flow_test.go"
    "goal_progress_flow_test.go"
)

all_files_exist=true
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "✓ $file"
    else
        echo "✗ $file (MISSING)"
        all_files_exist=false
    fi
done

echo ""

# Check Docker
echo "Checking Docker..."
if docker info > /dev/null 2>&1; then
    echo "✓ Docker is running"
else
    echo "✗ Docker is not running"
    exit 1
fi

# Check Docker Compose
echo "Checking Docker Compose..."
if command -v docker-compose &> /dev/null; then
    echo "✓ docker-compose is installed"
else
    echo "✗ docker-compose is not installed"
    exit 1
fi

# Check Go
echo "Checking Go..."
if command -v go &> /dev/null; then
    go_version=$(go version)
    echo "✓ Go is installed: $go_version"
else
    echo "✗ Go is not installed"
    exit 1
fi

echo ""

if [ "$all_files_exist" = true ]; then
    echo "========================================="
    echo "✓ All checks passed!"
    echo "========================================="
    echo ""
    echo "You can now run the integration tests:"
    echo "  ./run-tests.sh"
    echo "  or"
    echo "  make full-test"
    exit 0
else
    echo "========================================="
    echo "✗ Some files are missing"
    echo "========================================="
    exit 1
fi
