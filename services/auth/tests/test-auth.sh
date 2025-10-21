#!/bin/bash

echo "=== Starting Auth Service Test ==="

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

echo "=== Testing Health Endpoint ==="
curl -s http://localhost:3001/api/v1/health | jq '.'

echo -e "\n=== Testing User Registration ==="
curl -s -X POST http://localhost:3001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "securepass123",
    "firstName": "Test",
    "lastName": "User"
  }' | jq '.'

echo -e "\n=== Testing User Login ==="
curl -s -X POST http://localhost:3001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "securepass123"
  }' | jq '.'

# Clean up
echo -e "\n=== Cleaning up ==="
kill $SERVICE_PID 2>/dev/null
echo "Test completed!"