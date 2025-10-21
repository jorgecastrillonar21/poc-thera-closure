#!/bin/bash

# TheraClosure Geolocation Service - Complete Workflow Test
# Tests the complete country -> state -> city workflow

set -e

BASE_URL="http://localhost:3003"
API_URL="$BASE_URL/api/v1"

echo "🌍 TheraClosure Geolocation Service - Complete Workflow Test"
echo "============================================================"

# Colors for output  
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Store created IDs for cleanup
COUNTRY_ID=""
STATE_ID=""
CITY_ID=""

cleanup() {
    echo -e "\n${YELLOW}🧹 Cleaning up test data...${NC}"
    
    if [ -n "$CITY_ID" ]; then
        curl -s -X DELETE "$API_URL/cities/$CITY_ID" > /dev/null
        echo "Deleted city: $CITY_ID"
    fi
    
    if [ -n "$STATE_ID" ]; then
        curl -s -X DELETE "$API_URL/states/$STATE_ID" > /dev/null  
        echo "Deleted state: $STATE_ID"
    fi
    
    if [ -n "$COUNTRY_ID" ]; then
        curl -s -X DELETE "$API_URL/countries/$COUNTRY_ID" > /dev/null
        echo "Deleted country: $COUNTRY_ID"
    fi
}

# Set trap for cleanup on script exit
trap cleanup EXIT

# Extract ID from JSON response
extract_id() {
    local response=$1
    local entity=$2
    echo "$response" | grep -o "\"id\":\"[^\"]*\"" | cut -d'"' -f4
}

# Test function with JSON response handling
test_with_response() {
    local method=$1
    local url=$2 
    local data=$3
    local expected_status=$4
    local description=$5

    echo -e "\n${BLUE}Testing: $description${NC}"
    
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
        echo -e "${GREEN}✓ SUCCESS${NC} (Status: $http_code)"
        echo "$body" | jq . 2>/dev/null || echo "$body"
        echo "$body"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC} (Expected: $expected_status, Got: $http_code)"
        echo "Response: $body"
        exit 1
    fi
}

echo -e "\n${YELLOW}Step 1: Create a Country${NC}"
country_response=$(test_with_response "POST" "$API_URL/countries" '{
  "name": "United States",
  "code": "USA", 
  "code2": "US",
  "region": "North America",
  "currency": "USD"
}' 201 "Creating United States")

COUNTRY_ID=$(extract_id "$country_response" "country")
echo -e "${GREEN}📍 Created Country ID: $COUNTRY_ID${NC}"

echo -e "\n${YELLOW}Step 2: Create States in the Country${NC}"
state_response=$(test_with_response "POST" "$API_URL/states" '{
  "country_id": "'"$COUNTRY_ID"'",
  "name": "California", 
  "code": "CA"
}' 201 "Creating California state")

STATE_ID=$(extract_id "$state_response" "state")
echo -e "${GREEN}📍 Created State ID: $STATE_ID${NC}"

# Create another state
test_with_response "POST" "$API_URL/states" '{
  "country_id": "'"$COUNTRY_ID"'",
  "name": "New York",
  "code": "NY"
}' 201 "Creating New York state"

echo -e "\n${YELLOW}Step 3: Create Cities in the State${NC}"
city_response=$(test_with_response "POST" "$API_URL/cities" '{
  "state_id": "'"$STATE_ID"'",
  "name": "Los Angeles",
  "zip_code": "90210",
  "latitude": 34.0522,
  "longitude": -118.2437
}' 201 "Creating Los Angeles city")

CITY_ID=$(extract_id "$city_response" "city")
echo -e "${GREEN}📍 Created City ID: $CITY_ID${NC}"

# Create another city
test_with_response "POST" "$API_URL/cities" '{
  "state_id": "'"$STATE_ID"'",
  "name": "San Francisco", 
  "zip_code": "94102",
  "latitude": 37.7749,
  "longitude": -122.4194
}' 201 "Creating San Francisco city"

echo -e "\n${YELLOW}Step 4: Test Hierarchical Queries${NC}"

test_with_response "GET" "$API_URL/countries/$COUNTRY_ID/states" "" 200 "List states in country"

test_with_response "GET" "$API_URL/states/$STATE_ID/cities" "" 200 "List cities in state"

test_with_response "GET" "$API_URL/countries/$COUNTRY_ID/hierarchy" "" 200 "Get country hierarchy"

test_with_response "GET" "$API_URL/states/$STATE_ID/hierarchy" "" 200 "Get state hierarchy"

echo -e "\n${YELLOW}Step 5: Test Search Functionality${NC}"

test_with_response "GET" "$API_URL/countries/search?q=united" "" 200 "Search for countries containing 'united'"

test_with_response "GET" "$API_URL/states/search?q=california&country_id=$COUNTRY_ID" "" 200 "Search for states containing 'california'"

test_with_response "GET" "$API_URL/cities/search?q=los&state_id=$STATE_ID" "" 200 "Search for cities containing 'los'"

echo -e "\n${YELLOW}Step 6: Test Update Operations${NC}"

test_with_response "PUT" "$API_URL/countries/$COUNTRY_ID" '{
  "name": "United States of America",
  "region": "North America - Updated",
  "active": true
}' 200 "Update country information"

test_with_response "PUT" "$API_URL/cities/$CITY_ID" '{
  "name": "Los Angeles (Updated)",
  "zip_code": "90210",
  "active": true
}' 200 "Update city information"

echo -e "\n${YELLOW}Step 7: Test Individual Entity Retrieval${NC}"

test_with_response "GET" "$API_URL/countries/$COUNTRY_ID" "" 200 "Get country by ID"

test_with_response "GET" "$API_URL/states/$STATE_ID" "" 200 "Get state by ID"  

test_with_response "GET" "$API_URL/cities/$CITY_ID" "" 200 "Get city by ID"

echo -e "\n${YELLOW}Step 8: Test Complete Hierarchy${NC}"

test_with_response "GET" "$API_URL/hierarchy?limit=1" "" 200 "Get complete hierarchy"

echo -e "\n${GREEN}🎉 Complete workflow test passed successfully!${NC}"
echo -e "${GREEN}✅ Country → State → City hierarchy created and tested${NC}"
echo -e "${GREEN}✅ Search functionality working${NC}"  
echo -e "${GREEN}✅ Update operations working${NC}"
echo -e "${GREEN}✅ Hierarchical queries working${NC}"

echo -e "\n${BLUE}📊 Test Summary:${NC}"
echo "Country ID: $COUNTRY_ID"
echo "State ID: $STATE_ID"  
echo "City ID: $CITY_ID"