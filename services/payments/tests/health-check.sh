#!/bin/bash

# Health check for payments service
echo "Testing payments service health..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Set service URL
SERVICE_URL="http://localhost:3004"

# Function to make HTTP requests and check responses
test_endpoint() {
    local method=$1
    local endpoint=$2
    local expected_status=$3
    local data=$4
    local description=$5
    
    echo -e "${YELLOW}Testing: ${description}${NC}"
    echo "  ${method} ${SERVICE_URL}${endpoint}"
    
    if [[ -n "$data" ]]; then
        response=$(curl -s -w "\n%{http_code}" -X ${method} \
            -H "Content-Type: application/json" \
            -d "${data}" \
            "${SERVICE_URL}${endpoint}")
    else
        response=$(curl -s -w "\n%{http_code}" -X ${method} \
            "${SERVICE_URL}${endpoint}")
    fi
    
    # Split response body and status code
    body=$(echo "$response" | head -n -1)
    status_code=$(echo "$response" | tail -n 1)
    
    if [[ "$status_code" == "$expected_status" ]]; then
        echo -e "  ${GREEN}✓ Status: ${status_code} (Expected: ${expected_status})${NC}"
        if [[ -n "$body" ]]; then
            echo "  Response: $(echo "$body" | jq -c . 2>/dev/null || echo "$body")"
        fi
    else
        echo -e "  ${RED}✗ Status: ${status_code} (Expected: ${expected_status})${NC}"
        if [[ -n "$body" ]]; then
            echo "  Response: $(echo "$body" | jq -c . 2>/dev/null || echo "$body")"
        fi
    fi
    
    echo
}

# Check if service is running
echo -e "${YELLOW}Checking if payments service is running...${NC}"
if ! curl -s --connect-timeout 5 "${SERVICE_URL}/health" > /dev/null; then
    echo -e "${RED}✗ Payments service is not running at ${SERVICE_URL}${NC}"
    echo "Please start the service first:"
    echo "  cd services/payments && go run cmd/main.go"
    exit 1
fi

echo -e "${GREEN}✓ Payments service is running${NC}"
echo

# Test health endpoint
test_endpoint "GET" "/health" "200" "" "Health check"
test_endpoint "GET" "/api/v1/health" "200" "" "API health check"

# Test customer endpoints
echo -e "${YELLOW}Testing Customer Endpoints...${NC}"

# Create customer (will fail without proper validation)
customer_data='{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "test@example.com",
  "name": "Test Customer"
}'
test_endpoint "POST" "/api/v1/customers" "500" "$customer_data" "Create customer (expected to fail without Stripe key)"

# List customers
test_endpoint "GET" "/api/v1/customers" "200" "" "List customers"

# Test with invalid customer ID
test_endpoint "GET" "/api/v1/customers/invalid-id" "404" "" "Get invalid customer"

# Test subscription endpoints
echo -e "${YELLOW}Testing Subscription Endpoints...${NC}"
test_endpoint "GET" "/api/v1/subscriptions" "200" "" "List subscriptions"

# Test payment endpoints  
echo -e "${YELLOW}Testing Payment Endpoints...${NC}"
test_endpoint "GET" "/api/v1/payments" "200" "" "List payments"

# Test payment intent endpoints
echo -e "${YELLOW}Testing Payment Intent Endpoints...${NC}"

payment_intent_data='{
  "customer_id": "123e4567-e89b-12d3-a456-426614174000",
  "amount": 2000,
  "currency": "usd",
  "description": "Test payment"
}'
test_endpoint "POST" "/api/v1/payment-intents" "500" "$payment_intent_data" "Create payment intent (expected to fail without Stripe key)"

echo -e "${GREEN}Health check completed!${NC}"
echo -e "${YELLOW}Note: Some endpoints may fail without proper Stripe configuration${NC}"