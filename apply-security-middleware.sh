#!/bin/bash

# Script to apply security middleware to all InSavein services
# This script adds go-playground/validator dependency to all services

set -e

echo "================================================"
echo "InSavein Security Middleware Application Script"
echo "================================================"
echo ""

# List of services
SERVICES=(
    "auth-service"
    "user-service"
    "savings-service"
    "budget-service"
    "goal-service"
    "education-service"
    "notification-service"
    "analytics-service"
)

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Step 1: Adding validator dependency to all services..."
echo ""

for service in "${SERVICES[@]}"; do
    if [ -d "$service" ]; then
        echo -e "${YELLOW}Processing $service...${NC}"
        
        cd "$service"
        
        # Check if go.mod exists
        if [ -f "go.mod" ]; then
            # Add validator dependency
            echo "  - Adding go-playground/validator/v10..."
            go get github.com/go-playground/validator/v10
            
            # Tidy dependencies
            echo "  - Running go mod tidy..."
            go mod tidy
            
            echo -e "${GREEN}  ✓ $service updated successfully${NC}"
        else
            echo -e "${RED}  ✗ go.mod not found in $service${NC}"
        fi
        
        cd ..
        echo ""
    else
        echo -e "${RED}  ✗ Directory $service not found${NC}"
        echo ""
    fi
done

echo "================================================"
echo "Step 2: Verifying installations..."
echo "================================================"
echo ""

for service in "${SERVICES[@]}"; do
    if [ -d "$service" ] && [ -f "$service/go.mod" ]; then
        cd "$service"
        
        if grep -q "github.com/go-playground/validator/v10" go.mod; then
            echo -e "${GREEN}✓ $service: validator installed${NC}"
        else
            echo -e "${RED}✗ $service: validator NOT installed${NC}"
        fi
        
        cd ..
    fi
done

echo ""
echo "================================================"
echo "Step 3: Building services to verify..."
echo "================================================"
echo ""

BUILD_ERRORS=0

for service in "${SERVICES[@]}"; do
    if [ -d "$service" ]; then
        echo -e "${YELLOW}Building $service...${NC}"
        
        cd "$service"
        
        if go build -o /dev/null ./cmd/server 2>&1; then
            echo -e "${GREEN}  ✓ $service builds successfully${NC}"
        else
            echo -e "${RED}  ✗ $service build failed${NC}"
            BUILD_ERRORS=$((BUILD_ERRORS + 1))
        fi
        
        cd ..
        echo ""
    fi
done

echo "================================================"
echo "Summary"
echo "================================================"
echo ""

if [ $BUILD_ERRORS -eq 0 ]; then
    echo -e "${GREEN}✓ All services updated and building successfully!${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Update handler files to use validator"
    echo "2. Add validation tags to request structs"
    echo "3. Rebuild Docker images"
    echo "4. Deploy to Kubernetes"
    echo ""
    echo "See TASK_25_SECURITY_IMPLEMENTATION.md for detailed instructions."
else
    echo -e "${RED}✗ $BUILD_ERRORS service(s) failed to build${NC}"
    echo ""
    echo "Please check the errors above and fix them before proceeding."
fi

echo ""
echo "================================================"
echo "Done!"
echo "================================================"
