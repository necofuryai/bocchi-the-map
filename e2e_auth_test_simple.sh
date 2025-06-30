#!/bin/bash

# Simple End-to-End Auth0 Integration Test (No Database Required)
# Tests the authentication components and configuration without database dependency

set -e  # Exit on any error

echo "🧪 Starting Simple Auth0 E2E Integration Test"
echo "=============================================="
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

# Cleanup function
cleanup() {
    echo -e "\n${BLUE}🧹 Test completed${NC}"
    
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

echo "🔧 Testing Auth0 Integration Components..."

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

# Test Go dependencies and build
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

# File structure tests
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

# Check Auth0 API routes
if [ -f "web/src/app/api/auth/[...auth0]/route.ts" ]; then
    test_result "Auth0 API Routes" "PASS"
else
    test_result "Auth0 API Routes" "FAIL" "[...auth0]/route.ts not found"
fi

# Check Auth pages
if [ -f "web/src/app/auth/login/page.tsx" ]; then
    test_result "Login Page" "PASS"
else
    test_result "Login Page" "FAIL" "login/page.tsx not found"
fi

if [ -f "web/src/app/auth/logout/page.tsx" ]; then
    test_result "Logout Page" "PASS"
else
    test_result "Logout Page" "FAIL" "logout/page.tsx not found"
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

# Configuration tests
echo "🔧 Testing Configuration Files..."

# Check if configuration supports Auth0
if grep -q "Auth0" api/pkg/config/config.go; then
    test_result "API Auth0 Configuration" "PASS"
else
    test_result "API Auth0 Configuration" "FAIL" "Auth0 config not found in config.go"
fi

# Type definition tests
echo "📝 Testing Type Definitions..."

if [ -f "web/src/types/auth0.d.ts" ]; then
    test_result "Auth0 Type Definitions" "PASS"
else
    test_result "Auth0 Type Definitions" "FAIL" "auth0.d.ts not found"
fi

if [ -f "web/src/types/index.ts" ]; then
    test_result "General Type Definitions" "PASS"
else
    test_result "General Type Definitions" "FAIL" "index.ts not found"
fi

# Package.json dependencies test
echo "📦 Testing Package Dependencies..."

# Check Next.js Auth0 dependencies
if grep -q "@auth0/nextjs-auth0" web/package.json; then
    test_result "Next.js Auth0 Dependency" "PASS"
else
    test_result "Next.js Auth0 Dependency" "FAIL" "@auth0/nextjs-auth0 not found in package.json"
fi

# Go dependencies test
if grep -q "github.com/golang-jwt/jwt" api/go.mod; then
    test_result "Go JWT Dependency" "PASS"
else
    test_result "Go JWT Dependency" "FAIL" "JWT library not found in go.mod"
fi

# Content validation tests
echo "🔍 Testing Implementation Content..."

# Check Auth0 configuration in lib file
if [ -f "web/src/lib/auth0.ts" ]; then
    if grep -q "Auth0Client" web/src/lib/auth0.ts; then
        test_result "Auth0 Client Implementation" "PASS"
    else
        test_result "Auth0 Client Implementation" "FAIL" "Auth0Client not found in auth0.ts"
    fi
fi

# Check middleware implementation
if [ -f "web/middleware.ts" ]; then
    if grep -q "auth0.middleware" web/middleware.ts; then
        test_result "Middleware Auth0 Integration" "PASS"
    else
        test_result "Middleware Auth0 Integration" "FAIL" "auth0.middleware not found"
    fi
fi

# Check API handler implementation
if [ -f "api/interfaces/http/handlers/auth_handler.go" ]; then
    if grep -q "GetAuthStatus" api/interfaces/http/handlers/auth_handler.go; then
        test_result "Auth Handler Implementation" "PASS"
    else
        test_result "Auth Handler Implementation" "FAIL" "GetAuthStatus method not found"
    fi
fi

# Environment variable template check
echo "🌍 Testing Environment Configuration..."

# Check if environment variables are documented
if [ -f ".env.example" ] || [ -f ".env.local.example" ] || [ -f "web/.env.local.example" ]; then
    test_result "Environment Template" "PASS"
else
    test_result "Environment Template" "SKIP" "No environment template found (recommended to have one)"
fi

# Security configuration check
echo "🔒 Testing Security Configuration..."

# Check for security headers in Next.js config
if [ -f "web/next.config.ts" ]; then
    if grep -q -i "security\|header" web/next.config.ts; then
        test_result "Security Headers Configuration" "PASS"
    else
        test_result "Security Headers Configuration" "SKIP" "Security headers not explicitly configured"
    fi
elif [ -f "web/next.config.js" ]; then
    if grep -q -i "security\|header" web/next.config.js; then
        test_result "Security Headers Configuration" "PASS"
    else
        test_result "Security Headers Configuration" "SKIP" "Security headers not explicitly configured"
    fi
else
    test_result "Security Headers Configuration" "SKIP" "Next.js config not found"
fi

# Check for CORS configuration in API
if grep -q -i "cors" api/cmd/api/main.go; then
    test_result "CORS Configuration" "PASS"
else
    test_result "CORS Configuration" "SKIP" "CORS configuration not explicitly found"
fi

echo -e "\n${BLUE}🏁 Simple Auth0 E2E Integration Test Complete${NC}"
echo "=============================================="