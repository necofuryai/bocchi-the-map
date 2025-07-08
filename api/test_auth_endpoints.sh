#!/bin/bash

# Test Auth Endpoints Script
# This script tests the Go API authentication endpoints

echo "=== Testing Go API Authentication Endpoints ==="
echo

# Set required environment variables
export JWT_SECRET="test-jwt-secret-1234567890abcdefghijklmnopqrstuvwxyz-with-special-chars!@#$%"
export AUTH0_DOMAIN="test-domain.auth0.com"
export AUTH0_AUDIENCE="bocchi-the-map-api"
export AUTH0_CLIENT_ID="test-client-id-12345"
export AUTH0_CLIENT_SECRET="test-client-secret-67890"
export TIDB_PASSWORD="dummy-password-for-testing"
export ENV="development"
export LOG_LEVEL="INFO"
export PORT="8080"

# Build the API server
echo "Building API server..."
go build -o bin/api ./cmd/api/main.go
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"
echo

# Start the API server in background
echo "Starting API server on port 8080..."
./bin/api &
API_PID=$!

# Wait for server to start
sleep 3

echo "API server started with PID: $API_PID"
echo

# Function to cleanup on exit
cleanup() {
    echo "Stopping API server..."
    kill $API_PID 2>/dev/null
    wait $API_PID 2>/dev/null
    echo "✅ API server stopped"
}

# Setup cleanup on script exit
trap cleanup EXIT

# Test health endpoint first
echo "=== Testing Health Endpoint ==="
curl -s -w "Status: %{http_code}\n" http://localhost:8080/health
echo

# Test auth status endpoint (should work without authentication)
echo "=== Testing Auth Status Endpoint (without auth) ==="
curl -s -w "Status: %{http_code}\n" http://localhost:8080/api/v1/auth/status
echo

# Test validate token endpoint with invalid token
echo "=== Testing Token Validation Endpoint (invalid token) ==="
curl -s -w "Status: %{http_code}\n" \
  -H "Content-Type: application/json" \
  -d '{"token": "invalid.jwt.token"}' \
  http://localhost:8080/api/v1/auth/validate
echo

# Test protected endpoint without auth (should fail)
echo "=== Testing Protected Endpoint (without auth) ==="
curl -s -w "Status: %{http_code}\n" http://localhost:8080/api/v1/spots
echo

echo "=== Authentication endpoint tests completed ==="