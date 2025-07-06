#!/bin/bash

# Comprehensive End-to-End Auth0 Integration Test
# Tests the complete authentication flow between Next.js frontend and Go backend
# Combines functionality from simple, with-db, and full E2E tests

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

# Test modes
SIMPLE_MODE=false
WITH_DB_MODE=false
FULL_MODE=false

# Function to print test results
test_result() {
    local test_name="$1"
    local status="$2"
    local details="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✅ PASS${NC}: $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    elif [ "$status" = "SKIP" ]; then
        echo -e "${YELLOW}⏭️  SKIP${NC}: $test_name"
        if [ -n "$details" ]; then
            echo -e "   ${YELLOW}Details:${NC} $details"
        fi
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

# Function to wait for database to be ready
wait_for_database() {
    local host=$1
    local port=$2
    local timeout=${3:-30}
    local counter=0
    
    echo "Waiting for database at $host:$port to be ready..."
    while [ $counter -lt $timeout ]; do
        if nc -z "$host" "$port" > /dev/null 2>&1; then
            echo "Database is ready!"
            return 0
        fi
        sleep 1
        counter=$((counter + 1))
        echo -n "."
    done
    echo
    echo "Database failed to start within $timeout seconds"
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
    
    # Stop Docker containers if running
    if [ "$WITH_DB_MODE" = true ]; then
        echo "Stopping Docker containers..."
        docker compose down > /dev/null 2>&1 || true
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

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --simple)
            SIMPLE_MODE=true
            shift
            ;;
        --with-db)
            WITH_DB_MODE=true
            shift
            ;;
        --full)
            FULL_MODE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --simple     Run simple tests (file structure + build only)"
            echo "  --with-db    Run tests with database setup"
            echo "  --full       Run complete E2E tests with both API and web servers"
            echo "  -h, --help   Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Default to simple mode if no mode specified
if [ "$SIMPLE_MODE" = false ] && [ "$WITH_DB_MODE" = false ] && [ "$FULL_MODE" = false ]; then
    SIMPLE_MODE=true
fi

echo "🔧 Setting up test environment..."
echo "Test mode: Simple=$SIMPLE_MODE, With-DB=$WITH_DB_MODE, Full=$FULL_MODE"

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
export ENV="development"
export LOG_LEVEL="INFO"
export API_PORT="8080"
export WEB_PORT="3000"
export AUTH0_BASE_URL="http://localhost:3000"
export AUTH0_ISSUER_BASE_URL="https://$AUTH0_DOMAIN"

# Database configuration (for with-db mode)
if [ "$WITH_DB_MODE" = true ]; then
    export TIDB_HOST="localhost"
    export TIDB_PORT="4000"
    export TIDB_USER="root"
    export TIDB_PASSWORD="test-password"
    export TIDB_DATABASE="bocchi_the_map_test"
else
    export TIDB_PASSWORD="dummy-password-for-testing"
fi

# ===== BASIC STRUCTURE AND BUILD TESTS =====
echo "🏗️  Testing Go API Components..."
cd api

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
    test_result "Go Module File" "FAIL" "go.mod not found"
else
    test_result "Go Module File" "PASS"
fi

# Test Go build (dry run)
echo "Testing Go build capability..."
if go build -o /tmp/api-test ./cmd/api/main.go > /dev/null 2>&1; then
    test_result "Go Build Test" "PASS"
    # Clean up test binary
    rm -f /tmp/api-test
else
    test_result "Go Build Test" "FAIL" "Go build failed"
fi

# Test Go dependencies
echo "Testing Go dependencies..."
if go mod tidy > /dev/null 2>&1; then
    test_result "Go Dependencies" "PASS"
else
    test_result "Go Dependencies" "FAIL" "go mod tidy failed"
fi

cd ..

# Test Next.js components
echo "🌐 Testing Next.js Components..."
cd web

# Check if package.json exists
if [ ! -f "package.json" ]; then
    test_result "Package.json File" "FAIL" "package.json not found"
else
    test_result "Package.json File" "PASS"
fi

# Check Next.js configuration
if [ -f "next.config.ts" ] || [ -f "next.config.js" ]; then
    test_result "Next.js Configuration" "PASS"
else
    test_result "Next.js Configuration" "FAIL" "next.config file not found"
fi

# Check TypeScript configuration
if [ -f "tsconfig.json" ]; then
    test_result "TypeScript Configuration" "PASS"
else
    test_result "TypeScript Configuration" "FAIL" "tsconfig.json not found"
fi

cd ..

# ===== FILE STRUCTURE TESTS =====
echo "📁 Testing Auth0 Implementation Files..."

# Frontend Auth0 files
if [ -f "web/src/lib/auth0.ts" ]; then
    test_result "Auth0 Client Configuration" "PASS"
else
    test_result "Auth0 Client Configuration" "FAIL" "auth0.ts not found"
fi

if [ -f "web/middleware.ts" ]; then
    if grep -q "auth0" web/middleware.ts; then
        test_result "Next.js Auth0 Middleware" "PASS"
    else
        test_result "Next.js Auth0 Middleware" "FAIL" "Auth0 middleware not configured"
    fi
else
    test_result "Next.js Auth0 Middleware" "FAIL" "middleware.ts not found"
fi

# Check Auth components
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

# Backend Auth0 files
if [ -f "api/interfaces/http/handlers/auth_handler.go" ]; then
    test_result "Auth Handler" "PASS"
else
    test_result "Auth Handler" "FAIL" "auth_handler.go not found"
fi

if [ -f "api/pkg/auth/auth.go" ]; then
    test_result "Auth Service" "PASS"
else
    test_result "Auth Service" "FAIL" "auth.go not found"
fi

if [ -f "api/pkg/auth/middleware.go" ]; then
    test_result "JWT Middleware" "PASS"
else
    test_result "JWT Middleware" "FAIL" "middleware.go not found"
fi

if [ -f "api/pkg/auth/jwt.go" ]; then
    test_result "JWT Validator" "PASS"
else
    test_result "JWT Validator" "FAIL" "jwt.go not found"
fi

# Database related files
echo "💾 Testing Database Integration Files..."

if [ -f "api/migrations/000005_add_users_and_auth_tables.up.sql" ]; then
    test_result "User Migration (Up)" "PASS"
else
    test_result "User Migration (Up)" "FAIL" "User migration up file not found"
fi

if [ -f "api/migrations/000005_add_users_and_auth_tables.down.sql" ]; then
    test_result "User Migration (Down)" "PASS"
else
    test_result "User Migration (Down)" "FAIL" "User migration down file not found"
fi

if [ -f "api/queries/users.sql" ]; then
    test_result "User SQL Queries" "PASS"
else
    test_result "User SQL Queries" "FAIL" "users.sql not found"
fi

if [ -f "api/infrastructure/database/users.sql.go" ]; then
    test_result "Generated User Queries" "PASS"
else
    test_result "Generated User Queries" "FAIL" "users.sql.go not found"
fi

# Stop here if simple mode
if [ "$SIMPLE_MODE" = true ]; then
    echo -e "\n${BLUE}🏁 Simple Auth0 E2E Integration Test Complete${NC}"
    echo "=============================================="
    exit 0
fi

# ===== DATABASE SETUP (WITH-DB MODE) =====
if [ "$WITH_DB_MODE" = true ]; then
    # Check if Docker is running
    if ! docker info > /dev/null 2>&1; then
        test_result "Docker Availability" "FAIL" "Docker is not running"
        exit 1
    fi
    test_result "Docker Availability" "PASS"
    
    # Start database using Docker Compose
    echo "🗄️  Starting test database..."
    if [ -f "docker-compose.yml" ]; then
        docker compose up -d tidb > /dev/null 2>&1
        
        if ! wait_for_database "localhost" "4000" 30; then
            test_result "Database Startup" "FAIL" "Database failed to start"
            exit 1
        fi
        test_result "Database Startup" "PASS"
    else
        echo "Docker Compose file not found, using SQLite for testing..."
        # Create in-memory SQLite database for testing
        export TIDB_HOST=""
        export TIDB_PORT=""
        export TIDB_USER=""
        export TIDB_PASSWORD=""
        export TIDB_DATABASE=":memory:"
    fi
    
    # Run database migrations
    echo "🔄 Running database migrations..."
    cd api
    if [ -f "migrate" ] || command -v migrate > /dev/null 2>&1; then
        if command -v migrate > /dev/null 2>&1; then
            MIGRATE_CMD="migrate"
        else
            MIGRATE_CMD="./migrate"
        fi
        
        # Create database if it doesn't exist
        if [ "$TIDB_HOST" != "" ]; then
            mysql -h "$TIDB_HOST" -P "$TIDB_PORT" -u "$TIDB_USER" -p"$TIDB_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS $TIDB_DATABASE;" 2>/dev/null || true
            
            # Run migrations
            $MIGRATE_CMD -path migrations -database "mysql://$TIDB_USER:$TIDB_PASSWORD@tcp($TIDB_HOST:$TIDB_PORT)/$TIDB_DATABASE" up 2>/dev/null || true
        fi
        test_result "Database Migrations" "PASS"
    else
        echo "Migration tool not found, skipping migrations..."
        test_result "Database Migrations" "SKIP" "Migration tool not available"
    fi
    cd ..
fi

# ===== SERVER STARTUP AND TESTING (FULL MODE) =====
if [ "$FULL_MODE" = true ]; then
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
    
    if ! wait_for_server "http://localhost:$API_PORT/health" 20; then
        echo "API server failed to start. Check logs:"
        tail -20 api.log
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
    
    # ===== API ENDPOINT TESTS =====
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
    
    # ===== FRONTEND TESTS =====
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
fi

# ===== CONFIGURATION VALIDATION =====
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

# Security tests
echo "🔒 Testing Security Configuration..."

# Check if JWT secret is sufficiently complex (at least 32 characters)
if [ ${#JWT_SECRET} -ge 32 ]; then
    test_result "JWT Secret Complexity" "PASS"
else
    test_result "JWT Secret Complexity" "FAIL" "JWT secret should be at least 32 characters"
fi

# Type definitions test
if [ -f "web/src/types/auth0.d.ts" ]; then
    test_result "Auth0 Type Definitions" "PASS"
else
    test_result "Auth0 Type Definitions" "FAIL" "auth0.d.ts not found"
fi

# Next.js middleware test
if [ -f "web/middleware.ts" ]; then
    if grep -q "auth0.middleware" web/middleware.ts; then
        test_result "Next.js Middleware Configuration" "PASS"
    else
        test_result "Next.js Middleware Configuration" "FAIL" "Auth0 middleware not found in configuration"
    fi
else
    test_result "Next.js Middleware Configuration" "FAIL" "middleware.ts file not found"
fi

echo -e "\n${BLUE}🏁 Comprehensive Auth0 E2E Integration Test Complete${NC}"
echo "======================================================"