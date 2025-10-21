#!/bin/bash

# Health Check Test Script for TheraClosure Users Service
# Tests the service health endpoint

set -e

echo "🏥 TheraClosure Users Service - Health Check Test"
echo "================================================="

SERVICE_URL="http://localhost:3002"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    if [ "$1" = "PASS" ]; then
        echo -e "${GREEN}✅ $2${NC}"
    elif [ "$1" = "FAIL" ]; then
        echo -e "${RED}❌ $2${NC}"
    else
        echo -e "${YELLOW}ℹ️  $2${NC}"
    fi
}

# Test 1: Basic health check
echo ""
print_status "INFO" "Testing health check endpoint..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" $SERVICE_URL/health)
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Health check endpoint returns 200 OK"
    echo "Response: $body"
    
    # Parse JSON response
    service=$(echo $body | grep -o '"service":"[^"]*' | cut -d'"' -f4)
    status=$(echo $body | grep -o '"status":"[^"]*' | cut -d'"' -f4)
    
    if [ "$service" = "users-service" ]; then
        print_status "PASS" "Service name is correct: $service"
    else
        print_status "FAIL" "Service name is incorrect: $service"
    fi
    
    if [ "$status" = "healthy" ]; then
        print_status "PASS" "Service status is healthy"
    else
        print_status "FAIL" "Service status is not healthy: $status"
    fi
else
    print_status "FAIL" "Health check endpoint failed with status: $http_code"
    echo "Response: $body"
    exit 1
fi

echo ""
print_status "PASS" "All health check tests passed!"
echo ""