#!/bin/bash

# Comprehensive End-to-End Auth0 Integration Test
# Tests the complete authentication flow between Next.js frontend and Go backend

set -e  # Exit on any error

echo "🧪 Starting Comprehensive Auth0 E2E Integration Test"
echo "=================================================="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to print test results
test_result() {
    local test_name="$1"
    local status="$2"
    local details="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✅ PASS${NC}: $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAIL${NC}: $test_name"
        if [ -n "$details" ]; then
            echo -e "   ${YELLOW}Details:${NC} $details"
        fi
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo
}

# Function to check if port is in use
is_port_in_use() {
    local port=$1
    lsof -i :$port > /dev/null 2>&1
}

# Function to wait for server to be ready
wait_for_server() {
    local url=$1
    local timeout=${2:-30}
    local counter=0
    
    echo "Waiting for server at $url to be ready..."
    while [ $counter -lt $timeout ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            echo "Server is ready!"
            return 0
        fi
        sleep 1
        counter=$((counter + 1))
        echo -n "."
    done
    echo
    echo "Server failed to start within $timeout seconds"
    return 1
}

# Cleanup function
cleanup() {
    echo -e "\n${BLUE}🧹 Cleaning up...${NC}"
    
    # Kill background processes
    if [ -n "$API_PID" ]; then
        echo "Stopping API server (PID: $API_PID)..."
        kill $API_PID 2>/dev/null || true
        wait $API_PID 2>/dev/null || true
    fi
    
    if [ -n "$WEB_PID" ]; then
        echo "Stopping Web server (PID: $WEB_PID)..."
        kill $WEB_PID 2>/dev/null || true
        wait $WEB_PID 2>/dev/null || true
    fi
    
    # Print final results
    echo -e "\n${BLUE}📊 Test Results Summary${NC}"
    echo "======================="
    echo -e "Total Tests: $TOTAL_TESTS"
    echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "\n${GREEN}🎉 All tests passed!${NC}"
        exit 0
    else
        echo -e "\n${RED}💥 Some tests failed!${NC}"
        exit 1
    fi
}

# Setup cleanup on script exit
trap cleanup EXIT

echo "🔧 Setting up test environment..."

# Check if required directories exist
if [ ! -d "api" ]; then
    test_result "Directory Structure Check" "FAIL" "api directory not found"
    exit 1
fi

if [ ! -d "web" ]; then
    test_result "Directory Structure Check" "FAIL" "web directory not found"
    exit 1
fi

test_result "Directory Structure Check" "PASS"

# Set environment variables for testing
echo "Setting up environment variables..."
export JWT_SECRET="test-jwt-secret-1234567890abcdefghijklmnopqrstuvwxyz-with-special-chars!@#$%"
export AUTH0_DOMAIN="bocchi-the-map-dev.us.auth0.com"
export AUTH0_AUDIENCE="bocchi-the-map-api"
export AUTH0_CLIENT_ID="test-client-id"
export AUTH0_CLIENT_SECRET="test-client-secret"
export AUTH0_SECRET="test-auth0-secret-1234567890abcdefghijklmnopqrstuvwxyz"
export AUTH0_SCOPE="openid profile email"
export TIDB_PASSWORD="dummy-password-for-testing"
export ENV="development"
export LOG_LEVEL="INFO"
export API_PORT="8080"
export WEB_PORT="3000"
export AUTH0_BASE_URL="http://localhost:3000"
export AUTH0_ISSUER_BASE_URL="https://$AUTH0_DOMAIN"

# Check if ports are available
echo "Checking port availability..."
if is_port_in_use $API_PORT; then
    test_result "Port Availability Check (API)" "FAIL" "Port $API_PORT is already in use"
    exit 1
fi

if is_port_in_use $WEB_PORT; then
    test_result "Port Availability Check (Web)" "FAIL" "Port $WEB_PORT is already in use"
    exit 1
fi

test_result "Port Availability Check" "PASS"

# Build and start API server
echo "🏗️  Building Go API server..."
cd api
if ! go build -o ../bin/api ./cmd/api/main.go; then
    test_result "API Build" "FAIL" "Go build failed"
    exit 1
fi
test_result "API Build" "PASS"

echo "🚀 Starting API server on port $API_PORT..."
../bin/api > ../api.log 2>&1 &
API_PID=$!
cd ..

if ! wait_for_server "http://localhost:$API_PORT/health" 15; then
    test_result "API Server Startup" "FAIL" "API server failed to start"
    exit 1
fi
test_result "API Server Startup" "PASS"

# Install and start Next.js server
echo "📦 Installing Next.js dependencies..."
cd web
if ! npm install > ../web-install.log 2>&1; then
    test_result "Web Dependencies Install" "FAIL" "npm install failed"
    exit 1
fi
test_result "Web Dependencies Install" "PASS"

echo "🚀 Starting Next.js development server on port $WEB_PORT..."
npm run dev > ../web.log 2>&1 &
WEB_PID=$!
cd ..

if ! wait_for_server "http://localhost:$WEB_PORT" 30; then
    test_result "Web Server Startup" "FAIL" "Next.js server failed to start"
    exit 1
fi
test_result "Web Server Startup" "PASS"

# Test API endpoints
echo "🔍 Testing API Authentication Endpoints..."

# Test health endpoint
echo "Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s -w "%{http_code}" http://localhost:$API_PORT/health)
HEALTH_CODE=${HEALTH_RESPONSE: -3}
if [ "$HEALTH_CODE" = "200" ]; then
    test_result "API Health Endpoint" "PASS"
else
    test_result "API Health Endpoint" "FAIL" "Expected 200, got $HEALTH_CODE"
fi

# Test auth status endpoint (should work without authentication)
echo "Testing auth status endpoint without authentication..."
AUTH_STATUS_RESPONSE=$(curl -s -w "%{http_code}" http://localhost:$API_PORT/api/v1/auth/status)
AUTH_STATUS_CODE=${AUTH_STATUS_RESPONSE: -3}
if [ "$AUTH_STATUS_CODE" = "200" ]; then
    # Check if response indicates not authenticated
    if echo "$AUTH_STATUS_RESPONSE" | grep -q '"authenticated":false'; then
        test_result "Auth Status (Unauthenticated)" "PASS"
    else
        test_result "Auth Status (Unauthenticated)" "FAIL" "Should indicate not authenticated"
    fi
else
    test_result "Auth Status (Unauthenticated)" "FAIL" "Expected 200, got $AUTH_STATUS_CODE"
fi

# Test token validation with invalid token
echo "Testing token validation with invalid token..."
VALIDATE_RESPONSE=$(curl -s -w "%{http_code}" \
  -H "Content-Type: application/json" \
  -d '{"token": "invalid.jwt.token"}' \
  http://localhost:$API_PORT/api/v1/auth/validate)
VALIDATE_CODE=${VALIDATE_RESPONSE: -3}
if [ "$VALIDATE_CODE" = "200" ]; then
    # Check if response indicates invalid token
    if echo "$VALIDATE_RESPONSE" | grep -q '"valid":false'; then
        test_result "Token Validation (Invalid Token)" "PASS"
    else
        test_result "Token Validation (Invalid Token)" "FAIL" "Should indicate invalid token"
    fi
else
    test_result "Token Validation (Invalid Token)" "FAIL" "Expected 200, got $VALIDATE_CODE"
fi

# Test protected endpoint without auth (should fail)
echo "Testing protected endpoint without authentication..."
PROTECTED_RESPONSE=$(curl -s -w "%{http_code}" http://localhost:$API_PORT/api/v1/spots)
PROTECTED_CODE=${PROTECTED_RESPONSE: -3}
if [ "$PROTECTED_CODE" = "401" ]; then
    test_result "Protected Endpoint (No Auth)" "PASS"
else
    test_result "Protected Endpoint (No Auth)" "FAIL" "Expected 401, got $PROTECTED_CODE"
fi

# Test Next.js frontend endpoints
echo "🌐 Testing Next.js Frontend Integration..."

# Test main page
echo "Testing main page accessibility..."
MAIN_PAGE_RESPONSE=$(curl -s -w "%{http_code}" http://localhost:$WEB_PORT/)
MAIN_PAGE_CODE=${MAIN_PAGE_RESPONSE: -3}
if [ "$MAIN_PAGE_CODE" = "200" ]; then
    test_result "Main Page Accessibility" "PASS"
else
    test_result "Main Page Accessibility" "FAIL" "Expected 200, got $MAIN_PAGE_CODE"
fi

# Test Auth0 route configuration
echo "Testing Auth0 route configuration..."
AUTH_ROUTE_RESPONSE=$(curl -s -w "%{http_code}" http://localhost:$WEB_PORT/api/auth/login)
AUTH_ROUTE_CODE=${AUTH_ROUTE_RESPONSE: -3}
# Auth0 login should redirect (302) or return a page (200)
if [ "$AUTH_ROUTE_CODE" = "200" ] || [ "$AUTH_ROUTE_CODE" = "302" ] || [ "$AUTH_ROUTE_CODE" = "307" ]; then
    test_result "Auth0 Login Route" "PASS"
else
    test_result "Auth0 Login Route" "FAIL" "Expected 200/302/307, got $AUTH_ROUTE_CODE"
fi

# Test logout route
echo "Testing Auth0 logout route..."
LOGOUT_ROUTE_RESPONSE=$(curl -s -w "%{http_code}" http://localhost:$WEB_PORT/api/auth/logout)
LOGOUT_ROUTE_CODE=${LOGOUT_ROUTE_RESPONSE: -3}
if [ "$LOGOUT_ROUTE_CODE" = "200" ] || [ "$LOGOUT_ROUTE_CODE" = "302" ] || [ "$LOGOUT_ROUTE_CODE" = "307" ]; then
    test_result "Auth0 Logout Route" "PASS"
else
    test_result "Auth0 Logout Route" "FAIL" "Expected 200/302/307, got $LOGOUT_ROUTE_CODE"
fi

# Configuration validation tests
echo "🔧 Testing Configuration Validation..."

# Check if environment variables are properly set
ENV_VARS_OK=true
MISSING_VARS=""

check_env_var() {
    local var_name="$1"
    local var_value="${!var_name}"
    if [ -z "$var_value" ]; then
        ENV_VARS_OK=false
        MISSING_VARS="$MISSING_VARS $var_name"
    fi
}

check_env_var "AUTH0_DOMAIN"
check_env_var "AUTH0_CLIENT_ID"
check_env_var "AUTH0_CLIENT_SECRET"
check_env_var "AUTH0_SECRET"
check_env_var "JWT_SECRET"

if [ "$ENV_VARS_OK" = true ]; then
    test_result "Environment Variables" "PASS"
else
    test_result "Environment Variables" "FAIL" "Missing variables:$MISSING_VARS"
fi

# Test middleware configuration
echo "Testing Next.js middleware configuration..."
if [ -f "web/middleware.ts" ]; then
    if grep -q "auth0.middleware" web/middleware.ts; then
        test_result "Middleware Configuration" "PASS"
    else
        test_result "Middleware Configuration" "FAIL" "Auth0 middleware not found in configuration"
    fi
else
    test_result "Middleware Configuration" "FAIL" "middleware.ts file not found"
fi

# Component integration tests
echo "🧩 Testing Component Integration..."

# Check if Auth components exist
if [ -f "web/src/components/auth/auth-button.tsx" ]; then
    test_result "Auth Button Component" "PASS"
else
    test_result "Auth Button Component" "FAIL" "auth-button.tsx not found"
fi

if [ -f "web/src/components/auth/user-profile.tsx" ]; then
    test_result "User Profile Component" "PASS"
else
    test_result "User Profile Component" "FAIL" "user-profile.tsx not found"
fi

if [ -f "web/src/components/auth/auth-guard.tsx" ]; then
    test_result "Auth Guard Component" "PASS"
else
    test_result "Auth Guard Component" "FAIL" "auth-guard.tsx not found"
fi

# Backend auth handler tests
echo "🔐 Testing Backend Auth Handler..."

if [ -f "api/interfaces/http/handlers/auth_handler.go" ]; then
    test_result "Auth Handler Exists" "PASS"
else
    test_result "Auth Handler Exists" "FAIL" "auth_handler.go not found"
fi

# JWT middleware tests
if [ -f "api/pkg/auth/middleware.go" ]; then
    test_result "JWT Middleware Exists" "PASS"
else
    test_result "JWT Middleware Exists" "FAIL" "middleware.go not found"
fi

# Database integration tests
echo "💾 Testing Database Integration..."

# Check if user-related migrations exist
if [ -f "api/migrations/000005_add_users_and_auth_tables.up.sql" ]; then
    test_result "User Migration Files" "PASS"
else
    test_result "User Migration Files" "FAIL" "User migration not found"
fi

# Check if user queries exist
if [ -f "api/queries/users.sql" ]; then
    test_result "User SQL Queries" "PASS"
else
    test_result "User SQL Queries" "FAIL" "users.sql not found"
fi

# Security tests
echo "🔒 Testing Security Configuration..."

# Check if JWT secret is sufficiently complex (at least 32 characters)
if [ ${#JWT_SECRET} -ge 32 ]; then
    test_result "JWT Secret Complexity" "PASS"
else
    test_result "JWT Secret Complexity" "FAIL" "JWT secret should be at least 32 characters"
fi

# Integration flow simulation
echo "🔄 Testing Integration Flow Simulation..."

# Simulate API call from frontend to backend
echo "Simulating frontend-to-backend API call..."
API_CLIENT_RESPONSE=$(curl -s -w "%{http_code}" \
  -H "Content-Type: application/json" \
  -H "X-Requested-With: XMLHttpRequest" \
  http://localhost:$API_PORT/api/v1/auth/status)
API_CLIENT_CODE=${API_CLIENT_RESPONSE: -3}

if [ "$API_CLIENT_CODE" = "200" ]; then
    test_result "Frontend-to-Backend API Call" "PASS"
else
    test_result "Frontend-to-Backend API Call" "FAIL" "Expected 200, got $API_CLIENT_CODE"
fi

# Type definitions test
echo "📝 Testing TypeScript Definitions..."

if [ -f "web/src/types/auth0.d.ts" ]; then
    test_result "Auth0 Type Definitions" "PASS"
else
    test_result "Auth0 Type Definitions" "FAIL" "auth0.d.ts not found"
fi

echo -e "\n${BLUE}🏁 Comprehensive Auth0 E2E Integration Test Complete${NC}"
echo "======================================================"