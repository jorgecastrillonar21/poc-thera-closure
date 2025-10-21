#!/bin/bash

echo "=== Complete Auth Flow Test ==="

# Kill any existing processes
pkill -f "auth-service" 2>/dev/null

# Start the service
export DB_HOST=localhost
export DB_USER=theraclosure
export DB_PASSWORD=password123
export DB_NAME=theraclosure_auth
export JWT_SECRET=test-secret-key

echo "Starting auth service..."
./auth-service &
SERVICE_PID=$!

# Wait for service to start
sleep 3

echo "=== Step 1: Register New User ==="
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:3001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "completetest@example.com",
    "password": "securepass123",
    "firstName": "Complete",
    "lastName": "Test"
  }')

echo "$REGISTER_RESPONSE" | jq '.'

# Extract access token
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.accessToken')

echo -e "\n=== Step 2: Test Current User Endpoint ==="
curl -s -X GET http://localhost:3001/api/v1/auth/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'

echo -e "\n=== Step 3: Test Login with Same Credentials ==="
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "completetest@example.com",
    "password": "securepass123"
  }')

echo "$LOGIN_RESPONSE" | jq '.'

# Extract refresh token
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refreshToken')

echo -e "\n=== Step 4: Test Token Refresh ==="
curl -s -X POST http://localhost:3001/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refreshToken\": \"$REFRESH_TOKEN\"}" | jq '.'

# Clean up
echo -e "\n=== Cleaning up ==="
kill $SERVICE_PID 2>/dev/null
echo "Complete test finished!"