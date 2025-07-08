#!/bin/bash

# Simple Auth Endpoints Test Script
# This script tests the authentication endpoints with a lightweight server

echo "=== Testing Authentication Endpoints (Simple Server) ==="
echo

# Set required environment variables
export JWT_SECRET="test-jwt-secret-1234567890abcdefghijklmnopqrstuvwxyz-with-special-chars!@#$%"
export AUTH0_DOMAIN="test-domain.auth0.com"
export AUTH0_AUDIENCE="bocchi-the-map-api"
export AUTH0_CLIENT_ID="test-client-id-12345"
export ENV="development"
export LOG_LEVEL="INFO" 
export PORT="8080"

echo "Building test server..."
go build -o bin/test-server ./cmd/test-server/main.go
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"
echo

# Start the test server in background
echo "Starting test server on port 8080..."
./bin/test-server &
SERVER_PID=$!

# Wait for server to start
sleep 5

echo "Test server started with PID: $SERVER_PID"
echo

# Function to cleanup on exit
cleanup() {
    echo "Stopping test server..."
    kill $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
    echo "✅ Test server stopped"
}

# Setup cleanup on script exit
trap cleanup EXIT

# Test health endpoint
echo "=== 1. Testing Health Endpoint ==="
echo "Request: GET /health"
response=$(curl -s http://localhost:8080/health)
status_code=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
echo "Status: $status_code"
echo "Response: $response"
echo

# Test auth status endpoint (should work without authentication)
echo "=== 2. Testing Auth Status Endpoint (no auth) ==="
echo "Request: GET /api/v1/auth/status"
response=$(curl -s http://localhost:8080/api/v1/auth/status)
status_code=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/auth/status)
echo "Status: $status_code"
echo "Response: $response"
echo

# Test validate token endpoint with empty token
echo "=== 3. Testing Token Validation (empty token) ==="
echo "Request: POST /api/v1/auth/validate with empty token"
response=$(curl -s -H "Content-Type: application/json" -d '{"token": ""}' http://localhost:8080/api/v1/auth/validate)
status_code=$(curl -s -o /dev/null -w "%{http_code}" -H "Content-Type: application/json" -d '{"token": ""}' http://localhost:8080/api/v1/auth/validate)
echo "Status: $status_code"
echo "Response: $response"
echo

# Test validate token endpoint with invalid token
echo "=== 4. Testing Token Validation (invalid token) ==="
echo "Request: POST /api/v1/auth/validate with invalid token"
response=$(curl -s -H "Content-Type: application/json" -d '{"token": "invalid.jwt.token"}' http://localhost:8080/api/v1/auth/validate)
status_code=$(curl -s -o /dev/null -w "%{http_code}" -H "Content-Type: application/json" -d '{"token": "invalid.jwt.token"}' http://localhost:8080/api/v1/auth/validate)
echo "Status: $status_code"
echo "Response: $response"
echo

# Test validate token endpoint with longer token (should pass basic validation)
echo "=== 5. Testing Token Validation (valid format token) ==="
echo "Request: POST /api/v1/auth/validate with valid format token"
response=$(curl -s -H "Content-Type: application/json" -d '{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.test"}' http://localhost:8080/api/v1/auth/validate)
status_code=$(curl -s -o /dev/null -w "%{http_code}" -H "Content-Type: application/json" -d '{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.test"}' http://localhost:8080/api/v1/auth/validate)
echo "Status: $status_code"
echo "Response: $response"
echo

# Test protected endpoint without authentication
echo "=== 6. Testing Protected Endpoint (no auth) ==="
echo "Request: GET /api/v1/protected without Authorization header"
response=$(curl -s http://localhost:8080/api/v1/protected)
status_code=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/protected)
echo "Status: $status_code"
echo "Response: $response"
echo

# Test protected endpoint with Authorization header
echo "=== 7. Testing Protected Endpoint (with auth header) ==="
echo "Request: GET /api/v1/protected with Authorization header"
response=$(curl -s -H "Authorization: Bearer test-token-123" http://localhost:8080/api/v1/protected)
status_code=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer test-token-123" http://localhost:8080/api/v1/protected)
echo "Status: $status_code"
echo "Response: $response"
echo

echo "=== Authentication endpoint tests completed ==="
echo
echo "Expected Results:"
echo "- Health endpoint: 200 status with OK response"
echo "- Auth status: 200 status with authenticated=false"
echo "- Token validation (empty): 200 status with valid=false, error message"
echo "- Token validation (invalid): 200 status with valid=false, error message"
echo "- Token validation (valid format): 200 status with valid=true"
echo "- Protected endpoint (no auth): 200 status (should be 401 with proper auth middleware)"
echo "- Protected endpoint (with auth): 200 status"