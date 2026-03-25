#!/bin/bash

# Production Readiness Verification Script
# InSavein Platform
# This script performs automated checks for production readiness

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

# Namespace
NAMESPACE="insavein"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}InSavein Production Readiness Verification${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to print section header
print_section() {
    echo ""
    echo -e "${BLUE}=== $1 ===${NC}"
    echo ""
}

# Function to print success
print_success() {
    echo -e "${GREEN}✓${NC} $1"
    ((PASSED++))
}

# Function to print failure
print_failure() {
    echo -e "${RED}✗${NC} $1"
    ((FAILED++))
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
    ((WARNINGS++))
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
print_section "Checking Prerequisites"

if command_exists kubectl; then
    print_success "kubectl is installed"
else
    print_failure "kubectl is not installed"
    exit 1
fi

if command_exists docker; then
    print_success "docker is installed"
else
    print_warning "docker is not installed (optional for local testing)"
fi

if command_exists psql; then
    print_success "psql is installed"
else
    print_warning "psql is not installed (optional for database checks)"
fi

# Check Kubernetes cluster connectivity
print_section "Kubernetes Cluster Connectivity"

if kubectl cluster-info >/dev/null 2>&1; then
    print_success "Connected to Kubernetes cluster"
    kubectl cluster-info | head -n 1
else
    print_failure "Cannot connect to Kubernetes cluster"
    exit 1
fi

# Check namespace
print_section "Namespace Verification"

if kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
    print_success "Namespace '$NAMESPACE' exists"
else
    print_failure "Namespace '$NAMESPACE' does not exist"
    exit 1
fi

# Check deployments
print_section "Service Deployment Status"

SERVICES=(
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

for service in "${SERVICES[@]}"; do
    if kubectl get deployment "$service" -n "$NAMESPACE" >/dev/null 2>&1; then
        READY=$(kubectl get deployment "$service" -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}')
        DESIRED=$(kubectl get deployment "$service" -n "$NAMESPACE" -o jsonpath='{.spec.replicas}')
        
        if [ "$READY" = "$DESIRED" ] && [ "$READY" != "" ]; then
            print_success "$service: $READY/$DESIRED replicas ready"
        else
            print_failure "$service: $READY/$DESIRED replicas ready (expected $DESIRED)"
        fi
    else
        print_failure "$service: Deployment not found"
    fi
done

# Check pods
print_section "Pod Health Status"

TOTAL_PODS=$(kubectl get pods -n "$NAMESPACE" --no-headers 2>/dev/null | wc -l)
RUNNING_PODS=$(kubectl get pods -n "$NAMESPACE" --field-selector=status.phase=Running --no-headers 2>/dev/null | wc -l)

if [ "$TOTAL_PODS" -eq "$RUNNING_PODS" ] && [ "$TOTAL_PODS" -gt 0 ]; then
    print_success "All pods running: $RUNNING_PODS/$TOTAL_PODS"
else
    print_failure "Not all pods running: $RUNNING_PODS/$TOTAL_PODS"
    echo ""
    echo "Pod status:"
    kubectl get pods -n "$NAMESPACE"
fi

# Check for pods with restarts
print_section "Pod Restart Check"

HIGH_RESTART_PODS=$(kubectl get pods -n "$NAMESPACE" --no-headers 2>/dev/null | awk '{if ($4 > 5) print $1, $4}')

if [ -z "$HIGH_RESTART_PODS" ]; then
    print_success "No pods with excessive restarts (>5)"
else
    print_warning "Pods with high restart count:"
    echo "$HIGH_RESTART_PODS"
fi

# Check database
print_section "Database Status"

if kubectl get statefulset postgres -n "$NAMESPACE" >/dev/null 2>&1; then
    READY=$(kubectl get statefulset postgres -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}')
    DESIRED=$(kubectl get statefulset postgres -n "$NAMESPACE" -o jsonpath='{.spec.replicas}')
    
    if [ "$READY" = "$DESIRED" ] && [ "$READY" != "" ]; then
        print_success "PostgreSQL StatefulSet: $READY/$DESIRED replicas ready"
    else
        print_failure "PostgreSQL StatefulSet: $READY/$DESIRED replicas ready"
    fi
else
    print_warning "PostgreSQL StatefulSet not found (may be external)"
fi

# Check services
print_section "Kubernetes Services"

for service in "${SERVICES[@]}"; do
    if kubectl get service "$service" -n "$NAMESPACE" >/dev/null 2>&1; then
        CLUSTER_IP=$(kubectl get service "$service" -n "$NAMESPACE" -o jsonpath='{.spec.clusterIP}')
        print_success "$service: Service exists (ClusterIP: $CLUSTER_IP)"
    else
        print_failure "$service: Service not found"
    fi
done

# Check Ingress
print_section "Ingress Configuration"

if kubectl get ingress -n "$NAMESPACE" >/dev/null 2>&1; then
    INGRESS_COUNT=$(kubectl get ingress -n "$NAMESPACE" --no-headers | wc -l)
    if [ "$INGRESS_COUNT" -gt 0 ]; then
        print_success "Ingress configured ($INGRESS_COUNT found)"
        kubectl get ingress -n "$NAMESPACE"
    else
        print_warning "No Ingress resources found"
    fi
else
    print_warning "Cannot check Ingress resources"
fi

# Check ConfigMaps and Secrets
print_section "Configuration Resources"

if kubectl get configmap -n "$NAMESPACE" >/dev/null 2>&1; then
    CM_COUNT=$(kubectl get configmap -n "$NAMESPACE" --no-headers | wc -l)
    print_success "ConfigMaps found: $CM_COUNT"
else
    print_warning "No ConfigMaps found"
fi

if kubectl get secret -n "$NAMESPACE" >/dev/null 2>&1; then
    SECRET_COUNT=$(kubectl get secret -n "$NAMESPACE" --no-headers | wc -l)
    print_success "Secrets found: $SECRET_COUNT"
else
    print_warning "No Secrets found"
fi

# Check HorizontalPodAutoscalers
print_section "Horizontal Pod Autoscalers"

HPA_COUNT=$(kubectl get hpa -n "$NAMESPACE" --no-headers 2>/dev/null | wc -l)

if [ "$HPA_COUNT" -gt 0 ]; then
    print_success "HPAs configured: $HPA_COUNT"
    kubectl get hpa -n "$NAMESPACE"
else
    print_warning "No HPAs found"
fi

# Check monitoring stack
print_section "Monitoring Stack"

if kubectl get deployment prometheus -n "$NAMESPACE" >/dev/null 2>&1; then
    READY=$(kubectl get deployment prometheus -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}')
    if [ "$READY" -gt 0 ]; then
        print_success "Prometheus is running"
    else
        print_failure "Prometheus is not ready"
    fi
else
    print_warning "Prometheus deployment not found"
fi

if kubectl get deployment grafana -n "$NAMESPACE" >/dev/null 2>&1; then
    READY=$(kubectl get deployment grafana -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}')
    if [ "$READY" -gt 0 ]; then
        print_success "Grafana is running"
    else
        print_failure "Grafana is not ready"
    fi
else
    print_warning "Grafana deployment not found"
fi

# Check for alert rules
if kubectl get prometheusrule -n "$NAMESPACE" >/dev/null 2>&1; then
    RULE_COUNT=$(kubectl get prometheusrule -n "$NAMESPACE" --no-headers 2>/dev/null | wc -l)
    if [ "$RULE_COUNT" -gt 0 ]; then
        print_success "Prometheus alert rules configured: $RULE_COUNT"
    else
        print_warning "No Prometheus alert rules found"
    fi
else
    print_warning "Cannot check Prometheus alert rules"
fi

# Check TLS certificates
print_section "TLS Certificates"

if command_exists kubectl && kubectl get certificate -n "$NAMESPACE" >/dev/null 2>&1; then
    CERT_COUNT=$(kubectl get certificate -n "$NAMESPACE" --no-headers 2>/dev/null | wc -l)
    if [ "$CERT_COUNT" -gt 0 ]; then
        print_success "TLS certificates found: $CERT_COUNT"
        
        # Check certificate status
        while IFS= read -r cert; do
            CERT_NAME=$(echo "$cert" | awk '{print $1}')
            READY=$(echo "$cert" | awk '{print $2}')
            
            if [ "$READY" = "True" ]; then
                print_success "Certificate '$CERT_NAME' is ready"
            else
                print_failure "Certificate '$CERT_NAME' is not ready"
            fi
        done < <(kubectl get certificate -n "$NAMESPACE" --no-headers 2>/dev/null)
    else
        print_warning "No TLS certificates found"
    fi
else
    print_warning "Cannot check TLS certificates (cert-manager may not be installed)"
fi

# Check resource quotas
print_section "Resource Management"

if kubectl get resourcequota -n "$NAMESPACE" >/dev/null 2>&1; then
    QUOTA_COUNT=$(kubectl get resourcequota -n "$NAMESPACE" --no-headers 2>/dev/null | wc -l)
    if [ "$QUOTA_COUNT" -gt 0 ]; then
        print_success "Resource quotas configured: $QUOTA_COUNT"
    else
        print_warning "No resource quotas found"
    fi
else
    print_warning "Cannot check resource quotas"
fi

# Check network policies
if kubectl get networkpolicy -n "$NAMESPACE" >/dev/null 2>&1; then
    NP_COUNT=$(kubectl get networkpolicy -n "$NAMESPACE" --no-headers 2>/dev/null | wc -l)
    if [ "$NP_COUNT" -gt 0 ]; then
        print_success "Network policies configured: $NP_COUNT"
    else
        print_warning "No network policies found"
    fi
else
    print_warning "Cannot check network policies"
fi

# Check for integration tests
print_section "Testing Infrastructure"

if [ -d "integration-tests" ]; then
    print_success "Integration tests directory exists"
    
    if [ -f "integration-tests/Makefile" ]; then
        print_success "Integration test Makefile found"
    else
        print_warning "Integration test Makefile not found"
    fi
else
    print_warning "Integration tests directory not found"
fi

# Check for performance tests
if [ -d "performance-tests" ]; then
    print_success "Performance tests directory exists"
    
    if [ -f "performance-tests/normal-load.js" ]; then
        print_success "Normal load test script found"
    else
        print_warning "Normal load test script not found"
    fi
    
    if [ -f "performance-tests/peak-load.js" ]; then
        print_success "Peak load test script found"
    else
        print_warning "Peak load test script not found"
    fi
else
    print_warning "Performance tests directory not found"
fi

# Check documentation
print_section "Documentation"

DOCS=(
    "README.md"
    "docs/API_DOCUMENTATION.md"
    "docs/DEPLOYMENT.md"
    "docs/DEVELOPER_GUIDE.md"
    "docs/OPERATIONS_RUNBOOK.md"
)

for doc in "${DOCS[@]}"; do
    if [ -f "$doc" ]; then
        print_success "$doc exists"
    else
        print_warning "$doc not found"
    fi
done

# Check CI/CD workflows
print_section "CI/CD Pipelines"

WORKFLOWS=(
    ".github/workflows/lint.yml"
    ".github/workflows/test.yml"
    ".github/workflows/security.yml"
    ".github/workflows/build-push.yml"
    ".github/workflows/deploy-staging.yml"
    ".github/workflows/deploy-production.yml"
)

for workflow in "${WORKFLOWS[@]}"; do
    if [ -f "$workflow" ]; then
        print_success "$(basename "$workflow") exists"
    else
        print_warning "$(basename "$workflow") not found"
    fi
done

# Summary
print_section "Verification Summary"

echo ""
echo -e "${GREEN}Passed:${NC}   $PASSED"
echo -e "${YELLOW}Warnings:${NC} $WARNINGS"
echo -e "${RED}Failed:${NC}   $FAILED"
echo ""

if [ "$FAILED" -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Production readiness verification completed successfully!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Review warnings and address if necessary"
    echo "2. Run integration tests: cd integration-tests && make test"
    echo "3. Run performance tests: cd performance-tests && make test-normal"
    echo "4. Run security scans: trivy image <image-name>"
    echo "5. Review PRODUCTION_READINESS_CHECKLIST.md"
    exit 0
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}Production readiness verification failed!${NC}"
    echo -e "${RED}========================================${NC}"
    echo ""
    echo "Please address the failed checks before proceeding to production."
    exit 1
fi
