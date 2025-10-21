#!/bin/bash

# TheraClosure Users Service - Complete Test Suite Runner
# Runs all test scripts in the correct order with proper setup and cleanup

set -e

echo "🚀 TheraClosure Users Service - Complete Test Suite"
echo "=================================================="
echo "Running comprehensive tests for all service functionality"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Test configuration
SERVICE_URL="http://localhost:3002"
TESTS_DIR="$(dirname "$0")"
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to print colored output
print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_status() {
    if [ "$1" = "PASS" ]; then
        echo -e "${GREEN}✅ $2${NC}"
    elif [ "$1" = "FAIL" ]; then
        echo -e "${RED}❌ $2${NC}"
    elif [ "$1" = "INFO" ]; then
        echo -e "${YELLOW}ℹ️  $2${NC}"
    elif [ "$1" = "RUNNING" ]; then
        echo -e "${PURPLE}🔄 $2${NC}"
    fi
}

# Function to run a test script
run_test() {
    local test_name="$1"
    local test_script="$2"
    local test_description="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    print_header "$test_name"
    echo -e "${CYAN}Description: $test_description${NC}"
    echo ""
    
    print_status "RUNNING" "Starting $test_name..."
    echo ""
    
    if [ -f "$TESTS_DIR/$test_script" ] && [ -x "$TESTS_DIR/$test_script" ]; then
        if "$TESTS_DIR/$test_script"; then
            PASSED_TESTS=$((PASSED_TESTS + 1))
            print_status "PASS" "$test_name completed successfully!"
            echo ""
            echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
            echo ""
        else
            FAILED_TESTS=$((FAILED_TESTS + 1))
            print_status "FAIL" "$test_name failed!"
            echo ""
            echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
            echo ""
            return 1
        fi
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
        print_status "FAIL" "Test script $test_script not found or not executable!"
        echo ""
        return 1
    fi
}

# Pre-flight checks
print_header "PRE-FLIGHT CHECKS"

# Check if service is running
print_status "INFO" "Checking if Users Service is running..."
if curl -s -f $SERVICE_URL/health > /dev/null 2>&1; then
    service_response=$(curl -s $SERVICE_URL/health)
    print_status "PASS" "Users Service is running and healthy"
    echo "Service Response: $service_response"
else
    print_status "FAIL" "Users Service is not running at $SERVICE_URL"
    echo ""
    print_status "INFO" "Please start the Users Service first:"
    echo "cd /path/to/users/service && go run cmd/main.go"
    echo ""
    exit 1
fi

# Check if test scripts exist
print_status "INFO" "Checking test scripts..."
test_scripts=("health-check.sh" "user-profiles-crud.sh" "enrollment-workflow.sh")
for script in "${test_scripts[@]}"; do
    if [ -f "$TESTS_DIR/$script" ] && [ -x "$TESTS_DIR/$script" ]; then
        print_status "PASS" "$script found and executable"
    else
        print_status "FAIL" "$script not found or not executable"
        exit 1
    fi
done

echo ""
print_status "PASS" "All pre-flight checks passed!"
echo ""

# Run the test suite
print_header "STARTING TEST SUITE EXECUTION"
echo ""

# Test 1: Health Check
run_test "Health Check Test" "health-check.sh" "Verifies service health endpoint and basic connectivity"

# Brief pause between tests
sleep 2

# Test 2: User Profiles CRUD
run_test "User Profiles CRUD Test" "user-profiles-crud.sh" "Tests complete user profile lifecycle: Create, Read, Update, Delete, Search"

# Brief pause between tests
sleep 2

# Test 3: Enrollment Workflow
run_test "Enrollment Workflow Test" "enrollment-workflow.sh" "Tests complete 5-step enrollment process from start to completion"

# Final results
print_header "TEST SUITE RESULTS"

echo -e "${CYAN}📊 Test Execution Summary${NC}"
echo -e "${BLUE}=========================${NC}"
echo ""
echo -e "Total Tests Run:    ${YELLOW}$TOTAL_TESTS${NC}"
echo -e "Tests Passed:       ${GREEN}$PASSED_TESTS${NC}"
echo -e "Tests Failed:       ${RED}$FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! 🎉${NC}"
    echo -e "${GREEN}The TheraClosure Users Service is fully functional and ready for integration.${NC}"
    echo ""
    
    # Additional service info
    print_header "SERVICE INFORMATION"
    echo -e "${CYAN}Service URL:${NC} $SERVICE_URL"
    echo -e "${CYAN}Health Check:${NC} $SERVICE_URL/health"
    echo -e "${CYAN}API Endpoints:${NC}"
    echo "  • GET    /api/users/profiles           - List user profiles"
    echo "  • POST   /api/users/profiles           - Create user profile"
    echo "  • GET    /api/users/profiles/:id       - Get user profile by ID"
    echo "  • PUT    /api/users/profiles/:id       - Update user profile"
    echo "  • DELETE /api/users/profiles/:id       - Delete user profile"
    echo "  • GET    /api/users/profiles/search    - Search user profiles"
    echo "  • POST   /api/users/enrollments        - Start enrollment"
    echo "  • GET    /api/users/enrollments/:id    - Get enrollment status"
    echo "  • PUT    /api/users/enrollments/:id/steps/:step - Complete enrollment step"
    echo "  • GET    /api/users/enrollments/user/:userId    - List user enrollments"
    echo ""
    
    exit 0
else
    echo -e "${RED}❌ $FAILED_TESTS TEST(S) FAILED ❌${NC}"
    echo -e "${RED}Please review the failed tests and fix any issues before proceeding.${NC}"
    echo ""
    
    # Debugging information
    print_header "DEBUGGING INFORMATION"
    echo -e "${CYAN}Service Status:${NC}"
    curl -s $SERVICE_URL/health 2>/dev/null || echo "Service not responding"
    echo ""
    echo -e "${CYAN}Service Logs:${NC}"
    echo "Check the service console output for detailed error messages."
    echo ""
    echo -e "${CYAN}Common Issues:${NC}"
    echo "• Database connection problems"
    echo "• Missing environment variables"
    echo "• Port conflicts"
    echo "• Database schema not migrated"
    echo ""
    
    exit 1
fi