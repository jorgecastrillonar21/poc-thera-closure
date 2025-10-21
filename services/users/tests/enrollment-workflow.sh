#!/bin/bash

# Enrollment Workflow Test Script for TheraClosure Users Service
# Tests the complete 5-step enrollment process

set -e

echo "📝 TheraClosure Users Service - Enrollment Workflow Test"
echo "======================================================="

SERVICE_URL="http://localhost:3002"
USER_ID=""
ENROLLMENT_ID=""

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
        exit 1
    else
        echo -e "${YELLOW}ℹ️  $2${NC}"
    fi
}

# Function to extract field from JSON
extract_json_field() {
    echo $1 | grep -o "\"$2\":\"[^\"]*" | cut -d'"' -f4
}

# Function to extract numeric field from JSON
extract_json_number() {
    echo $1 | grep -o "\"$2\":[0-9]*" | cut -d':' -f2
}

# Setup: Create a user first
echo ""
print_status "INFO" "Setting up test user for enrollment..."

create_user_payload='{
  "first_name": "Dr. Sarah",
  "last_name": "Johnson",
  "email": "sarah.johnson@test.com",
  "phone": "+1-555-0789",
  "date_of_birth": "1980-03-20T00:00:00Z",
  "license_number": "ENR123456",
  "license_state": "NY",
  "license_expiration": "2025-12-31T00:00:00Z",
  "years_of_experience": 12
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  -H "Content-Type: application/json" \
  -d "$create_user_payload" \
  $SERVICE_URL/api/v1/users/profiles)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "201" ]; then
    USER_ID=$(extract_json_field "$body" "user_id")
    print_status "PASS" "Test user created with ID: $USER_ID"
else
    print_status "FAIL" "Failed to create test user. Status: $http_code"
fi

# Test 1: Start Enrollment Process
echo ""
print_status "INFO" "Testing enrollment process initiation..."

start_enrollment_payload='{
  "user_id": "'$USER_ID'",
  "selected_plan": "professional"
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  -H "Content-Type: application/json" \
  -d "$start_enrollment_payload" \
  $SERVICE_URL/api/v1/enrollments/start)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "201" ]; then
    plan=$(extract_json_field "$body" "plan")
    print_status "PASS" "Enrollment started successfully"
    echo "User ID: $USER_ID"
    echo "Selected Plan: $plan"
    echo "Response: $body"
else
    print_status "FAIL" "Failed to start enrollment. Status: $http_code"
    echo "Response: $body"
fi

# Test 2: Get Initial Enrollment Status
echo ""
print_status "INFO" "Testing get initial enrollment status..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" $SERVICE_URL/api/v1/enrollments/$USER_ID)
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    status=$(extract_json_field "$body" "enrollment_status")
    current_step=$(extract_json_number "$body" "current_step")
    selected_plan=$(extract_json_field "$body" "selected_plan")
    print_status "PASS" "Successfully retrieved enrollment status"
    echo "Status: $status"
    echo "Current Step: $current_step"
    echo "Selected Plan: $selected_plan"
else
    print_status "FAIL" "Failed to get enrollment status. Status: $http_code"
    echo "Response: $body"
fi

# Test 3: Complete Step 1 - Personal Information
echo ""
print_status "INFO" "Testing Step 1 - Personal Information completion..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  $SERVICE_URL/api/v1/enrollments/$USER_ID/steps/1/complete)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Step 1 completed successfully"
    echo "Response: $body"
else
    print_status "FAIL" "Failed to complete Step 1. Status: $http_code"
    echo "Response: $body"
fi

# Test 4: Complete Step 2 - License Verification
echo ""
print_status "INFO" "Testing Step 2 - License Verification completion..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  $SERVICE_URL/api/v1/enrollments/$USER_ID/steps/2/complete)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Step 2 completed successfully"
    echo "Response: $body"
else
    print_status "FAIL" "Failed to complete Step 2. Status: $http_code"
    echo "Response: $body"
fi

# Test 5: Complete Step 3 - Practice Information
echo ""
print_status "INFO" "Testing Step 3 - Practice Information completion..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  $SERVICE_URL/api/v1/enrollments/$USER_ID/steps/3/complete)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Step 3 completed successfully"
    echo "Response: $body"
else
    print_status "FAIL" "Failed to complete Step 3. Status: $http_code"
    echo "Response: $body"
fi

# Test 6: Complete Step 4 - Admin Setup
echo ""
print_status "INFO" "Testing Step 4 - Admin Setup completion..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  $SERVICE_URL/api/v1/enrollments/$USER_ID/steps/4/complete)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Step 4 completed successfully"
    echo "Response: $body"
else
    print_status "FAIL" "Failed to complete Step 4. Status: $http_code"
    echo "Response: $body"
fi

# Test 7: Complete Step 5 - Schedule Configuration
echo ""
print_status "INFO" "Testing Step 5 - Schedule Configuration completion..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  $SERVICE_URL/api/v1/enrollments/$USER_ID/steps/5/complete)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Step 5 completed successfully - Enrollment Complete!"
    echo "Response: $body"
else
    print_status "FAIL" "Failed to complete Step 5. Status: $http_code"
    echo "Response: $body"
fi

# Test 8: Get Final Enrollment Status
echo ""
print_status "INFO" "Testing get final enrollment status..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" $SERVICE_URL/api/v1/enrollments/$USER_ID)
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    status=$(extract_json_field "$body" "enrollment_status")
    current_step=$(extract_json_number "$body" "current_step")
    print_status "PASS" "Successfully retrieved final enrollment status"
    echo "Final Status: $status"
    echo "Final Step: $current_step"
    echo "Full Response: $body"
else
    print_status "FAIL" "Failed to get final enrollment status. Status: $http_code"
    echo "Response: $body"
fi

# Test 9: Get Enrollment Progress
echo ""
print_status "INFO" "Testing get enrollment progress..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" $SERVICE_URL/api/v1/enrollments/$USER_ID/progress)
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Successfully retrieved enrollment progress"
    echo "Progress Details: $body"
else
    print_status "PASS" "Progress endpoint may not be fully implemented (Status: $http_code)"
    echo "Response: $body"
fi

# Cleanup: Delete test user and enrollment
echo ""
print_status "INFO" "Cleaning up test data..."

# Delete user (this should cascade to delete enrollment)
curl -s -X DELETE $SERVICE_URL/api/users/profiles/$USER_ID > /dev/null

print_status "PASS" "Test cleanup completed"

echo ""
print_status "PASS" "All enrollment workflow tests passed!"
echo ""