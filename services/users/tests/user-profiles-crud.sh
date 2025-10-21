#!/bin/bash

# User Profiles CRUD Test Script for TheraClosure Users Service
# Tests all user profile operations: Create, Read, Update, Delete

set -e

echo "👤 TheraClosure Users Service - User Profiles CRUD Test"
echo "======================================================"

SERVICE_URL="http://localhost:3002"
USER_ID=""

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

# Test 1: Create User Profile
echo ""
print_status "INFO" "Testing user profile creation..."

create_payload='{
  "first_name": "Dr. Jane",
  "last_name": "Smith",
  "email": "jane.smith@therapist.com",
  "phone": "+1-555-0123",
  "date_of_birth": "1985-06-15T00:00:00Z",
  "license_number": "LIC123456",
  "license_state": "CA",
  "license_expiration": "2025-12-31T00:00:00Z",
  "years_of_experience": 8,
  "practice_name": "Mindful Therapy Center",
  "practice_address": "123 Main St, San Francisco, CA 94102",
  "practice_phone": "+1-555-0456"
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  -H "Content-Type: application/json" \
  -d "$create_payload" \
  $SERVICE_URL/api/v1/users/profiles)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "201" ]; then
    print_status "PASS" "User profile created successfully"
    USER_ID=$(extract_json_field "$body" "user_id")
    echo "Created User ID: $USER_ID"
    echo "Response: $body"
else
    print_status "FAIL" "Failed to create user profile. Status: $http_code"
    echo "Response: $body"
fi

# Test 2: Get User Profile by ID
echo ""
print_status "INFO" "Testing get user profile by ID..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" $SERVICE_URL/api/v1/users/profiles/$USER_ID)
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Successfully retrieved user profile"
    firstName=$(extract_json_field "$body" "first_name")
    email=$(extract_json_field "$body" "email")
    echo "Retrieved: $firstName ($email)"
else
    print_status "FAIL" "Failed to get user profile. Status: $http_code"
    echo "Response: $body"
fi

# Test 3: Update User Profile
echo ""
print_status "INFO" "Testing user profile update..."

update_payload='{
  "first_name": "Dr. Jane",
  "last_name": "Smith-Johnson",
  "email": "jane.smith@therapist.com",
  "phone": "+1-555-0124",
  "years_of_experience": 9
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT \
  -H "Content-Type: application/json" \
  -d "$update_payload" \
  $SERVICE_URL/api/v1/users/profiles/$USER_ID)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "User profile updated successfully"
    lastName=$(extract_json_field "$body" "last_name")
    yearsExp=$(extract_json_number "$body" "years_of_experience")
    echo "Updated lastName: $lastName, Years Experience: $yearsExp"
else
    print_status "FAIL" "Failed to update user profile. Status: $http_code"
    echo "Response: $body"
fi

# Test 4: List User Profiles (Pagination)
echo ""
print_status "INFO" "Testing list user profiles with pagination..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" "$SERVICE_URL/api/v1/users/profiles?page=1&limit=10")
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Successfully retrieved user profiles list"
    total=$(echo $body | grep -o '"total":[0-9]*' | cut -d':' -f2)
    echo "Total profiles: $total"
else
    print_status "FAIL" "Failed to get user profiles list. Status: $http_code"
    echo "Response: $body"
fi

# Test 5: Search User Profiles
echo ""
print_status "INFO" "Testing search user profiles..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" "$SERVICE_URL/api/v1/users/profiles/search?q=anxiety&specialties=anxiety")
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Successfully searched user profiles"
    echo "Search results: $body"
else
    print_status "FAIL" "Failed to search user profiles. Status: $http_code"
    echo "Response: $body"
fi

# Test 6: Delete User Profile
echo ""
print_status "INFO" "Testing user profile deletion..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X DELETE $SERVICE_URL/api/v1/users/profiles/$USER_ID)
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')

if [ "$http_code" = "204" ]; then
    print_status "PASS" "User profile deleted successfully"
else
    print_status "FAIL" "Failed to delete user profile. Status: $http_code"
fi

# Test 7: Verify Deletion
echo ""
print_status "INFO" "Verifying user profile deletion..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" $SERVICE_URL/api/v1/users/profiles/$USER_ID)
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')

if [ "$http_code" = "404" ]; then
    print_status "PASS" "Confirmed user profile was deleted (404 Not Found)"
else
    print_status "FAIL" "User profile still exists after deletion. Status: $http_code"
fi

echo ""
print_status "PASS" "All user profile CRUD tests passed!"
echo ""