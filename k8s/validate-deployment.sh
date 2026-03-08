#!/bin/bash

# InSavein Kubernetes Deployment Validator
# This script validates that all base configurations are correctly deployed

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
PASSED=0
FAILED=0
WARNINGS=0

echo "=================================================="
echo "InSavein Kubernetes Deployment Validator"
echo "=================================================="
echo ""

# Function to check if a resource exists
check_resource() {
    local resource_type=$1
    local resource_name=$2
    local namespace=$3
    local description=$4
    
    echo -n "Checking $description... "
    
    if [ -z "$namespace" ]; then
        if kubectl get "$resource_type" "$resource_name" &> /dev/null; then
            echo -e "${GREEN}✓ PASS${NC}"
            ((PASSED++))
            return 0
        else
            echo -e "${RED}✗ FAIL${NC}"
            echo "  Resource not found: $resource_type/$resource_name"
            ((FAILED++))
            return 1
        fi
    else
        if kubectl get "$resource_type" "$resource_name" -n "$namespace" &> /dev/null; then
            echo -e "${GREEN}✓ PASS${NC}"
            ((PASSED++))
            return 0
        else
            echo -e "${RED}✗ FAIL${NC}"
            echo "  Resource not found: $resource_type/$resource_name in namespace $namespace"
            ((FAILED++))
            return 1
        fi
    fi
}

# Function to check resource count
check_resource_count() {
    local resource_type=$1
    local namespace=$2
    local expected_count=$3
    local description=$4
    
    echo -n "Checking $description... "
    
    local actual_count
    if [ -z "$namespace" ]; then
        actual_count=$(kubectl get "$resource_type" --no-headers 2>/dev/null | wc -l)
    else
        actual_count=$(kubectl get "$resource_type" -n "$namespace" --no-headers 2>/dev/null | wc -l)
    fi
    
    if [ "$actual_count" -ge "$expected_count" ]; then
        echo -e "${GREEN}✓ PASS${NC} (found $actual_count)"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC}"
        echo "  Expected at least $expected_count, found $actual_count"
        ((FAILED++))
        return 1
    fi
}

# Function to check for placeholder values in secrets
check_secrets_placeholders() {
    echo -n "Checking for placeholder values in secrets... "
    
    if kubectl get secret insavein-secrets -n insavein -o yaml 2>/dev/null | grep -q "CHANGE_ME"; then
        echo -e "${YELLOW}⚠ WARNING${NC}"
        echo "  Secrets contain placeholder values (CHANGE_ME)"
        echo "  Run 'make generate-secrets' to generate secure values"
        ((WARNINGS++))
        return 1
    else
        echo -e "${GREEN}✓ PASS${NC}"
        ((PASSED++))
        return 0
    fi
}

# Function to check ConfigMap values
check_configmap_values() {
    local key=$1
    local expected_pattern=$2
    local description=$3
    
    echo -n "Checking $description... "
    
    local value
    value=$(kubectl get configmap insavein-config -n insavein -o jsonpath="{.data.$key}" 2>/dev/null)
    
    if [ -z "$value" ]; then
        echo -e "${RED}✗ FAIL${NC}"
        echo "  ConfigMap key '$key' not found"
        ((FAILED++))
        return 1
    fi
    
    if [[ "$value" =~ $expected_pattern ]]; then
        echo -e "${GREEN}✓ PASS${NC} ($value)"
        ((PASSED++))
        return 0
    else
        echo -e "${YELLOW}⚠ WARNING${NC}"
        echo "  Value '$value' doesn't match expected pattern '$expected_pattern'"
        ((WARNINGS++))
        return 1
    fi
}

echo -e "${BLUE}=== Checking Cluster Access ===${NC}"
echo ""

if ! kubectl cluster-info &> /dev/null; then
    echo -e "${RED}✗ FAIL: Cannot connect to Kubernetes cluster${NC}"
    echo "Please check your kubectl configuration"
    exit 1
fi

echo -e "${GREEN}✓ Cluster access verified${NC}"
echo ""

echo -e "${BLUE}=== Checking Namespace ===${NC}"
echo ""

check_resource "namespace" "insavein" "" "namespace 'insavein'"
echo ""

echo -e "${BLUE}=== Checking Priority Classes ===${NC}"
echo ""

check_resource "priorityclass" "insavein-critical" "" "priority class 'insavein-critical'"
check_resource "priorityclass" "insavein-high" "" "priority class 'insavein-high'"
check_resource "priorityclass" "insavein-medium" "" "priority class 'insavein-medium'"
check_resource "priorityclass" "insavein-low" "" "priority class 'insavein-low'"
echo ""

echo -e "${BLUE}=== Checking ConfigMap ===${NC}"
echo ""

check_resource "configmap" "insavein-config" "insavein" "ConfigMap 'insavein-config'"

# Check specific ConfigMap values
check_configmap_values "DB_HOST" "postgres" "database host configuration"
check_configmap_values "DB_PORT" "5432" "database port configuration"
check_configmap_values "RATE_LIMIT_PER_USER" "100" "rate limit per user (Requirement 18.1)"
check_configmap_values "RATE_LIMIT_PER_IP" "1000" "rate limit per IP (Requirement 18.1)"
check_configmap_values "JWT_ACCESS_TOKEN_EXPIRY" "15m" "JWT access token expiry"
check_configmap_values "JWT_REFRESH_TOKEN_EXPIRY" "168h" "JWT refresh token expiry"
echo ""

echo -e "${BLUE}=== Checking Secrets ===${NC}"
echo ""

check_resource "secret" "insavein-secrets" "insavein" "Secret 'insavein-secrets'"
check_resource "secret" "postgres-credentials" "insavein" "Secret 'postgres-credentials'"
check_secrets_placeholders
echo ""

echo -e "${BLUE}=== Checking Resource Quotas ===${NC}"
echo ""

check_resource "resourcequota" "insavein-resource-quota" "insavein" "ResourceQuota 'insavein-resource-quota'"

# Check quota details
echo -n "Checking resource quota limits... "
QUOTA_CPU=$(kubectl get resourcequota insavein-resource-quota -n insavein -o jsonpath='{.spec.hard.requests\.cpu}' 2>/dev/null)
if [ "$QUOTA_CPU" = "50" ]; then
    echo -e "${GREEN}✓ PASS${NC} (CPU: $QUOTA_CPU)"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠ WARNING${NC}"
    echo "  Expected CPU quota: 50, found: $QUOTA_CPU"
    ((WARNINGS++))
fi
echo ""

echo -e "${BLUE}=== Checking Limit Ranges ===${NC}"
echo ""

check_resource "limitrange" "insavein-limit-range" "insavein" "LimitRange 'insavein-limit-range'"
check_resource "limitrange" "insavein-database-limits" "insavein" "LimitRange 'insavein-database-limits'"
echo ""

echo -e "${BLUE}=== Checking Network Policies ===${NC}"
echo ""

check_resource_count "networkpolicy" "insavein" 7 "network policies"

check_resource "networkpolicy" "insavein-default-deny" "insavein" "default deny policy"
check_resource "networkpolicy" "allow-frontend-to-services" "insavein" "frontend to services policy"
check_resource "networkpolicy" "allow-services-to-database" "insavein" "services to database policy"
check_resource "networkpolicy" "allow-ingress-to-frontend" "insavein" "ingress to frontend policy"
check_resource "networkpolicy" "allow-prometheus-scraping" "insavein" "Prometheus scraping policy"
check_resource "networkpolicy" "allow-dns-egress" "insavein" "DNS egress policy"
check_resource "networkpolicy" "allow-external-apis" "insavein" "external APIs policy"
echo ""

echo -e "${BLUE}=== Resource Usage Summary ===${NC}"
echo ""

# Show resource quota usage
echo "Resource Quota Usage:"
kubectl describe resourcequota insavein-resource-quota -n insavein 2>/dev/null | grep -A 10 "Used"
echo ""

echo "=================================================="
echo -e "${BLUE}Validation Summary${NC}"
echo "=================================================="
echo ""
echo -e "${GREEN}Passed:${NC}   $PASSED"
echo -e "${RED}Failed:${NC}   $FAILED"
echo -e "${YELLOW}Warnings:${NC} $WARNINGS"
echo ""

if [ $FAILED -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✓ All checks passed! Deployment is valid.${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Deploy PostgreSQL database (primary + replicas)"
    echo "2. Deploy backend microservices"
    echo "3. Deploy frontend application"
    echo "4. Set up monitoring and observability"
    exit 0
elif [ $FAILED -eq 0 ]; then
    echo -e "${YELLOW}⚠ Deployment is valid but has warnings.${NC}"
    echo "Please review the warnings above."
    exit 0
else
    echo -e "${RED}✗ Deployment validation failed!${NC}"
    echo "Please fix the errors above and try again."
    exit 1
fi
