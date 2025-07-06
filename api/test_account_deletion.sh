#!/bin/bash

# Simple test script for account deletion functionality
# This script tests the DELETE /api/v1/users/me endpoint

set -e

API_BASE_URL="http://localhost:8080"
AUTH_TOKEN=""

echo "Testing account deletion functionality..."

# Function to check if the API is running
check_api() {
    if ! curl -s "$API_BASE_URL/health" > /dev/null 2>&1; then
        echo "API is not running at $API_BASE_URL"
        echo "Please start the API server first with: go run cmd/api/main.go"
        exit 1
    fi
    echo "✓ API is running"
}

# Function to test delete endpoint
test_delete_endpoint() {
    echo "Testing DELETE /api/v1/users/me endpoint..."
    
    # Test without authentication (should return 401)
    echo "  - Testing without authentication..."
    response=$(curl -s -w "%{http_code}" -o /dev/null -X DELETE "$API_BASE_URL/api/v1/users/me")
    if [ "$response" = "401" ]; then
        echo "  ✓ Returns 401 without authentication"
    else
        echo "  ✗ Expected 401, got $response"
        exit 1
    fi
    
    # Test with invalid token (should return 401)
    echo "  - Testing with invalid token..."
    response=$(curl -s -w "%{http_code}" -o /dev/null \
        -X DELETE "$API_BASE_URL/api/v1/users/me" \
        -H "Authorization: Bearer invalid_token")
    if [ "$response" = "401" ]; then
        echo "  ✓ Returns 401 with invalid token"
    else
        echo "  ✗ Expected 401, got $response"
        exit 1
    fi
    
    echo "  ✓ Authentication tests passed"
}

# Function to test endpoint availability
test_endpoint_availability() {
    echo "Testing endpoint availability..."
    
    # Check if the endpoint exists (should not return 404)
    response=$(curl -s -w "%{http_code}" -o /dev/null \
        -X DELETE "$API_BASE_URL/api/v1/users/me" \
        -H "Authorization: Bearer dummy")
    
    if [ "$response" = "404" ]; then
        echo "  ✗ Endpoint not found (404)"
        exit 1
    elif [ "$response" = "401" ]; then
        echo "  ✓ Endpoint exists (returns 401 for invalid auth)"
    else
        echo "  ? Endpoint exists (returns $response)"
    fi
}

# Main test execution
main() {
    echo "=== Account Deletion API Test ==="
    check_api
    test_endpoint_availability
    test_delete_endpoint
    echo "=== All tests passed! ==="
}

main "$@"