#!/bin/bash

echo "=== Refresh Token Debug Test ==="

# Kill any existing processes
pkill -f "auth-service" 2>/dev/null

# Start the service
export DB_HOST=localhost
export DB_USER=theraclosure
export DB_PASSWORD=password123
export DB_NAME=theraclosure_auth
export JWT_SECRET=test-secret-key
export REDIS_HOST=localhost
export REDIS_PORT=6379

echo "Starting auth service..."
./auth-service &
SERVICE_PID=$!

# Wait for service to start
sleep 3

# Generate unique email
UNIQUE_EMAIL="refresh-test-$(date +%s)@example.com"

echo "=== Step 1: Register User with unique email: $UNIQUE_EMAIL ==="
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:3001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$UNIQUE_EMAIL\",
    \"password\": \"securepass123\",
    \"firstName\": \"Refresh\",
    \"lastName\": \"Test\"
  }")

echo "$REGISTER_RESPONSE" | jq '.'

# Extract tokens
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.accessToken')
REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.refreshToken')

echo -e "\n=== Debug: Extracted Tokens ==="
echo "Access Token: $ACCESS_TOKEN"
echo "Refresh Token: $REFRESH_TOKEN"

echo -e "\n=== Step 2: Test Token Refresh ==="
REFRESH_RESPONSE=$(curl -s -X POST http://localhost:3001/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refreshToken\": \"$REFRESH_TOKEN\"}")

echo "$REFRESH_RESPONSE" | jq '.'

# Clean up
echo -e "\n=== Cleaning up ==="
kill $SERVICE_PID 2>/dev/null
echo "Debug test finished!"