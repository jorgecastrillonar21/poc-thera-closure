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
  "firstName": "Dr. Sarah",
  "lastName": "Johnson",
  "email": "sarah.johnson@test.com",
  "phoneNumber": "+1-555-0789",
  "dateOfBirth": "1980-03-20",
  "gender": "female",
  "licenseNumber": "ENR123456",
  "licenseState": "NY",
  "licenseExpiryDate": "2025-12-31",
  "specialties": ["family-therapy", "couples-counseling"],
  "yearsOfExperience": 12,
  "education": "M.A. Marriage and Family Therapy, NYU"
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  -H "Content-Type: application/json" \
  -d "$create_user_payload" \
  $SERVICE_URL/api/users/profiles)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "201" ]; then
    USER_ID=$(extract_json_field "$body" "id")
    print_status "PASS" "Test user created with ID: $USER_ID"
else
    print_status "FAIL" "Failed to create test user. Status: $http_code"
fi

# Test 1: Start Enrollment Process
echo ""
print_status "INFO" "Testing enrollment process initiation..."

start_enrollment_payload='{
  "userID": "'$USER_ID'",
  "enrollmentType": "individual-practice"
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  -H "Content-Type: application/json" \
  -d "$start_enrollment_payload" \
  $SERVICE_URL/api/users/enrollments)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "201" ]; then
    ENROLLMENT_ID=$(extract_json_field "$body" "id")
    current_step=$(extract_json_number "$body" "currentStep")
    status=$(extract_json_field "$body" "status")
    print_status "PASS" "Enrollment started successfully"
    echo "Enrollment ID: $ENROLLMENT_ID"
    echo "Current Step: $current_step"
    echo "Status: $status"
else
    print_status "FAIL" "Failed to start enrollment. Status: $http_code"
    echo "Response: $body"
fi

# Test 2: Complete Step 1 - Personal Information
echo ""
print_status "INFO" "Testing Step 1 - Personal Information completion..."

step1_payload='{
  "stepData": {
    "personalInfo": {
      "address": "456 Therapy Lane, New York, NY 10001",
      "emergencyContact": "John Johnson",
      "emergencyPhone": "+1-555-0987"
    }
  }
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT \
  -H "Content-Type: application/json" \
  -d "$step1_payload" \
  $SERVICE_URL/api/users/enrollments/$ENROLLMENT_ID/steps/1)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    current_step=$(extract_json_number "$body" "currentStep")
    print_status "PASS" "Step 1 completed successfully"
    echo "Advanced to Step: $current_step"
else
    print_status "FAIL" "Failed to complete Step 1. Status: $http_code"
    echo "Response: $body"
fi

# Test 3: Complete Step 2 - License Verification
echo ""
print_status "INFO" "Testing Step 2 - License Verification completion..."

step2_payload='{
  "stepData": {
    "licenseVerification": {
      "verificationMethod": "manual-review",
      "documentsUploaded": ["license-copy.pdf", "certification.pdf"]
    }
  }
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT \
  -H "Content-Type: application/json" \
  -d "$step2_payload" \
  $SERVICE_URL/api/users/enrollments/$ENROLLMENT_ID/steps/2)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    current_step=$(extract_json_number "$body" "currentStep")
    print_status "PASS" "Step 2 completed successfully"
    echo "Advanced to Step: $current_step"
else
    print_status "FAIL" "Failed to complete Step 2. Status: $http_code"
    echo "Response: $body"
fi

# Test 4: Complete Step 3 - Background Check
echo ""
print_status "INFO" "Testing Step 3 - Background Check completion..."

step3_payload='{
  "stepData": {
    "backgroundCheck": {
      "consentProvided": true,
      "checkType": "standard",
      "providerUsed": "trusted-background-services"
    }
  }
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT \
  -H "Content-Type: application/json" \
  -d "$step3_payload" \
  $SERVICE_URL/api/users/enrollments/$ENROLLMENT_ID/steps/3)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    current_step=$(extract_json_number "$body" "currentStep")
    print_status "PASS" "Step 3 completed successfully"
    echo "Advanced to Step: $current_step"
else
    print_status "FAIL" "Failed to complete Step 3. Status: $http_code"
    echo "Response: $body"
fi

# Test 5: Complete Step 4 - Insurance Setup
echo ""
print_status "INFO" "Testing Step 4 - Insurance Setup completion..."

step4_payload='{
  "stepData": {
    "insuranceSetup": {
      "acceptsInsurance": true,
      "networks": ["Anthem", "Blue Cross Blue Shield"],
      "malpracticeInsurance": {
        "provider": "Professional Risk Management",
        "policyNumber": "PRM-789456123",
        "expiryDate": "2025-06-30"
      }
    }
  }
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT \
  -H "Content-Type: application/json" \
  -d "$step4_payload" \
  $SERVICE_URL/api/users/enrollments/$ENROLLMENT_ID/steps/4)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    current_step=$(extract_json_number "$body" "currentStep")
    print_status "PASS" "Step 4 completed successfully"
    echo "Advanced to Step: $current_step"
else
    print_status "FAIL" "Failed to complete Step 4. Status: $http_code"
    echo "Response: $body"
fi

# Test 6: Complete Step 5 - Final Review
echo ""
print_status "INFO" "Testing Step 5 - Final Review completion..."

step5_payload='{
  "stepData": {
    "finalReview": {
      "agreementSigned": true,
      "termsAccepted": true,
      "profileComplete": true
    }
  }
}'

response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT \
  -H "Content-Type: application/json" \
  -d "$step5_payload" \
  $SERVICE_URL/api/users/enrollments/$ENROLLMENT_ID/steps/5)

http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    status=$(extract_json_field "$body" "status")
    completion_date=$(extract_json_field "$body" "completedAt")
    print_status "PASS" "Step 5 completed successfully"
    echo "Final Status: $status"
    echo "Completion Date: $completion_date"
else
    print_status "FAIL" "Failed to complete Step 5. Status: $http_code"
    echo "Response: $body"
fi

# Test 7: Get Enrollment Status
echo ""
print_status "INFO" "Testing get enrollment status..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" $SERVICE_URL/api/users/enrollments/$ENROLLMENT_ID)
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    status=$(extract_json_field "$body" "status")
    current_step=$(extract_json_number "$body" "currentStep")
    print_status "PASS" "Successfully retrieved enrollment status"
    echo "Status: $status"
    echo "Current Step: $current_step"
else
    print_status "FAIL" "Failed to get enrollment status. Status: $http_code"
    echo "Response: $body"
fi

# Test 8: List User Enrollments
echo ""
print_status "INFO" "Testing list user enrollments..."

response=$(curl -s -w "HTTPSTATUS:%{http_code}" $SERVICE_URL/api/users/enrollments/user/$USER_ID)
http_code=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')

if [ "$http_code" = "200" ]; then
    print_status "PASS" "Successfully retrieved user enrollments"
    echo "Enrollments: $body"
else
    print_status "FAIL" "Failed to get user enrollments. Status: $http_code"
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