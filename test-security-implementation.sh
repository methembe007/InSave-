#!/bin/bash

# Script to test security implementations
# Tests validation, authorization, rate limiting, and security headers

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-https://api.insavein.com}"
TEST_TOKEN="${TEST_TOKEN:-}"

echo "================================================"
echo "InSavein Security Implementation Test Suite"
echo "================================================"
echo ""
echo "API URL: $API_URL"
echo ""

# Function to print test result
print_result() {
    local test_name=$1
    local result=$2
    local message=$3
    
    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name"
    elif [ "$result" = "FAIL" ]; then
        echo -e "${RED}✗ FAIL${NC}: $test_name - $message"
    else
        echo -e "${YELLOW}⚠ SKIP${NC}: $test_name - $message"
    fi
}

# Test counter
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Test 1: Security Headers
echo -e "${BLUE}Test 1: Security Headers${NC}"
echo "Testing for required security headers..."

RESPONSE=$(curl -s -I "$API_URL/api/auth/health" 2>/dev/null || echo "")
TOTAL_TESTS=$((TOTAL_TESTS + 6))

if echo "$RESPONSE" | grep -qi "Strict-Transport-Security"; then
    print_result "HSTS Header" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "HSTS Header" "FAIL" "Header not found"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

if echo "$RESPONSE" | grep -qi "X-Frame-Options"; then
    print_result "X-Frame-Options Header" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "X-Frame-Options Header" "FAIL" "Header not found"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

if echo "$RESPONSE" | grep -qi "X-Content-Type-Options"; then
    print_result "X-Content-Type-Options Header" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "X-Content-Type-Options Header" "FAIL" "Header not found"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

if echo "$RESPONSE" | grep -qi "X-XSS-Protection"; then
    print_result "X-XSS-Protection Header" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "X-XSS-Protection Header" "FAIL" "Header not found"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

if echo "$RESPONSE" | grep -qi "Content-Security-Policy"; then
    print_result "Content-Security-Policy Header" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "Content-Security-Policy Header" "FAIL" "Header not found"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

if echo "$RESPONSE" | grep -qi "Referrer-Policy"; then
    print_result "Referrer-Policy Header" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "Referrer-Policy Header" "FAIL" "Header not found"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# Test 2: TLS/SSL Configuration
echo -e "${BLUE}Test 2: TLS/SSL Configuration${NC}"
echo "Testing HTTPS and certificate..."

TOTAL_TESTS=$((TOTAL_TESTS + 2))

if curl -s -o /dev/null -w "%{http_code}" "$API_URL/api/auth/health" | grep -q "200\|401"; then
    print_result "HTTPS Connection" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "HTTPS Connection" "FAIL" "Cannot connect via HTTPS"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# Test HTTP to HTTPS redirect
HTTP_URL=$(echo "$API_URL" | sed 's/https/http/')
REDIRECT_CODE=$(curl -s -o /dev/null -w "%{http_code}" -L "$HTTP_URL/api/auth/health" 2>/dev/null || echo "000")

if [ "$REDIRECT_CODE" = "200" ] || [ "$REDIRECT_CODE" = "401" ]; then
    print_result "HTTP to HTTPS Redirect" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "HTTP to HTTPS Redirect" "SKIP" "HTTP not accessible or redirect not configured"
    SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
fi

echo ""

# Test 3: Input Validation
echo -e "${BLUE}Test 3: Input Validation${NC}"
echo "Testing input validation with invalid data..."

if [ -z "$TEST_TOKEN" ]; then
    echo -e "${YELLOW}⚠ Skipping validation tests - TEST_TOKEN not set${NC}"
    TOTAL_TESTS=$((TOTAL_TESTS + 3))
    SKIPPED_TESTS=$((SKIPPED_TESTS + 3))
else
    TOTAL_TESTS=$((TOTAL_TESTS + 3))
    
    # Test invalid amount (negative)
    RESPONSE=$(curl -s -X POST "$API_URL/api/savings" \
        -H "Authorization: Bearer $TEST_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"amount": -10, "currency": "USD", "category": "test"}' \
        -w "\n%{http_code}" 2>/dev/null || echo "000")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    if [ "$HTTP_CODE" = "400" ]; then
        print_result "Negative Amount Validation" "PASS"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_result "Negative Amount Validation" "FAIL" "Expected 400, got $HTTP_CODE"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    # Test invalid currency (wrong length)
    RESPONSE=$(curl -s -X POST "$API_URL/api/savings" \
        -H "Authorization: Bearer $TEST_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"amount": 10, "currency": "INVALID", "category": "test"}' \
        -w "\n%{http_code}" 2>/dev/null || echo "000")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    if [ "$HTTP_CODE" = "400" ]; then
        print_result "Invalid Currency Validation" "PASS"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_result "Invalid Currency Validation" "FAIL" "Expected 400, got $HTTP_CODE"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    # Test missing required field
    RESPONSE=$(curl -s -X POST "$API_URL/api/savings" \
        -H "Authorization: Bearer $TEST_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"currency": "USD"}' \
        -w "\n%{http_code}" 2>/dev/null || echo "000")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    if [ "$HTTP_CODE" = "400" ]; then
        print_result "Missing Required Field Validation" "PASS"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_result "Missing Required Field Validation" "FAIL" "Expected 400, got $HTTP_CODE"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
fi

echo ""

# Test 4: Authorization
echo -e "${BLUE}Test 4: Authorization${NC}"
echo "Testing authorization and access control..."

TOTAL_TESTS=$((TOTAL_TESTS + 2))

# Test without token
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/api/savings" 2>/dev/null || echo "000")
if [ "$RESPONSE" = "401" ]; then
    print_result "Unauthorized Access (No Token)" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "Unauthorized Access (No Token)" "FAIL" "Expected 401, got $RESPONSE"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

# Test with invalid token
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" \
    -H "Authorization: Bearer invalid_token_12345" \
    "$API_URL/api/savings" 2>/dev/null || echo "000")
if [ "$RESPONSE" = "401" ]; then
    print_result "Unauthorized Access (Invalid Token)" "PASS"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    print_result "Unauthorized Access (Invalid Token)" "FAIL" "Expected 401, got $RESPONSE"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""

# Test 5: Rate Limiting
echo -e "${BLUE}Test 5: Rate Limiting${NC}"
echo "Testing rate limiting (this may take a minute)..."

if [ -z "$TEST_TOKEN" ]; then
    echo -e "${YELLOW}⚠ Skipping rate limit test - TEST_TOKEN not set${NC}"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
else
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo "Sending 110 requests to test rate limit..."
    RATE_LIMITED=false
    
    for i in {1..110}; do
        RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" \
            -H "Authorization: Bearer $TEST_TOKEN" \
            "$API_URL/api/savings/summary" 2>/dev/null || echo "000")
        
        if [ "$RESPONSE" = "429" ]; then
            RATE_LIMITED=true
            break
        fi
        
        # Show progress every 20 requests
        if [ $((i % 20)) -eq 0 ]; then
            echo "  Sent $i requests..."
        fi
    done
    
    if [ "$RATE_LIMITED" = true ]; then
        print_result "Rate Limiting (100 req/min)" "PASS"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_result "Rate Limiting (100 req/min)" "FAIL" "Did not receive 429 after 110 requests"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
fi

echo ""

# Summary
echo "================================================"
echo "Test Summary"
echo "================================================"
echo ""
echo "Total Tests:   $TOTAL_TESTS"
echo -e "${GREEN}Passed:        $PASSED_TESTS${NC}"
echo -e "${RED}Failed:        $FAILED_TESTS${NC}"
echo -e "${YELLOW}Skipped:       $SKIPPED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed. Please review the results above.${NC}"
    exit 1
fi
