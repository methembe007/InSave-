#!/bin/bash

# Integration test script for User Profile Service
# This script tests the service against a running database

set -e

echo "=== User Profile Service Integration Test ==="
echo ""

# Check if service is running
SERVICE_URL="${SERVICE_URL:-http://localhost:8081}"

echo "Testing health endpoints..."
curl -f "$SERVICE_URL/health" || { echo "Health check failed"; exit 1; }
echo "✓ Health check passed"

curl -f "$SERVICE_URL/health/live" || { echo "Liveness check failed"; exit 1; }
echo "✓ Liveness check passed"

curl -f "$SERVICE_URL/health/ready" || { echo "Readiness check failed"; exit 1; }
echo "✓ Readiness check passed"

echo ""
echo "=== All integration tests passed! ==="
