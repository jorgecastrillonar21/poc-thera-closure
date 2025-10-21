#!/bin/bash

# TheraClosure Geolocation Service - Smoke Tests
# Tests basic functionality of the geolocation service

set -e

BASE_URL="http://localhost:3003"
API_URL="$BASE_URL/api/v1"

echo "🧪 TheraClosure Geolocation Service - Smoke Tests"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper function to test HTTP endpoints
test_endpoint() {
    local method=$1
    local url=$2
    local data=$3
    local expected_status=$4
    local description=$5

    echo -n "Testing: $description... "
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$url")
    else
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" "$url")
    fi
    
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $response | sed -E 's/HTTPSTATUS\:[0-9]{3}$//')
    
    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ PASS${NC} ($http_code)"
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (expected $expected_status, got $http_code)"
        echo "Response: $body"
        return 1
    fi
}

# 1. Health Check
echo -e "\n${YELLOW}1. Health Check${NC}"
test_endpoint "GET" "$BASE_URL/health" "" 200 "Service health check"

# 2. Country Operations
echo -e "\n${YELLOW}2. Country Operations${NC}"

# Create a country
country_data='{
  "name": "Test Country",
  "code": "TST",
  "code2": "TC",
  "region": "Test Region",
  "currency": "TST"
}'

test_endpoint "POST" "$API_URL/countries" "$country_data" 201 "Create country"

# List countries
test_endpoint "GET" "$API_URL/countries?limit=10&offset=0" "" 200 "List countries"

# Search countries
test_endpoint "GET" "$API_URL/countries/search?q=test&limit=5" "" 200 "Search countries"

# 3. Bulk Operations Test
echo -e "\n${YELLOW}3. Bulk Operations${NC}"

bulk_countries='{
  "countries": [
    {
      "name": "Bulk Country 1",
      "code": "BC1",
      "code2": "B1",
      "region": "Test Region",
      "currency": "TST"
    },
    {
      "name": "Bulk Country 2", 
      "code": "BC2",
      "code2": "B2",
      "region": "Test Region",
      "currency": "TST"
    }
  ]
}'

test_endpoint "POST" "$API_URL/bulk/countries" "$bulk_countries" 201 "Bulk create countries"

# 4. Hierarchy Test
echo -e "\n${YELLOW}4. Hierarchy Operations${NC}"
test_endpoint "GET" "$API_URL/hierarchy?limit=5" "" 200 "Get complete hierarchy"

# 5. Invalid Data Tests
echo -e "\n${YELLOW}5. Validation Tests${NC}"

# Test missing required fields
invalid_country='{"name": ""}'
test_endpoint "POST" "$API_URL/countries" "$invalid_country" 400 "Create country with missing fields"

# Test invalid endpoints
test_endpoint "GET" "$API_URL/countries/invalid-id" "" 404 "Get non-existent country"

echo -e "\n${GREEN}🎉 Smoke tests completed!${NC}"
echo "If all tests passed, the Geolocation service is working correctly."