#!/bin/bash

# Stripe Integration Test Script
# This script runs comprehensive Stripe integration tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Stripe Integration Test Suite ===${NC}"
echo

# Check if Stripe test key is configured
if [ -z "${STRIPE_TEST_SECRET_KEY}" ]; then
    echo -e "${YELLOW}Warning: STRIPE_TEST_SECRET_KEY environment variable not set${NC}"
    echo -e "${YELLOW}Using default test key for basic validation tests${NC}"
    echo
fi

# Build the service first
echo -e "${GREEN}Building payments service...${NC}"
if ! go build ./cmd/main.go; then
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi
echo -e "${GREEN}Build successful${NC}"
echo

# Run unit tests first
echo -e "${GREEN}Running unit tests...${NC}"
go test -short ./internal/adapters/stripe/ -v
echo

# Run integration tests if Stripe key is available
if [ -n "${STRIPE_TEST_SECRET_KEY}" ]; then
    echo -e "${GREEN}Running Stripe integration tests...${NC}"
    go test ./internal/adapters/stripe/ -v -run "TestStripeClient"
    
    echo
    echo -e "${GREEN}Running performance benchmarks...${NC}"
    go test ./internal/adapters/stripe/ -bench=. -benchmem -run="^$"
else
    echo -e "${YELLOW}Skipping integration tests (no Stripe test key configured)${NC}"
    echo -e "${YELLOW}To run integration tests, set STRIPE_TEST_SECRET_KEY environment variable${NC}"
fi

echo
echo -e "${GREEN}=== Test Results Summary ===${NC}"

# Run test coverage
echo -e "${GREEN}Generating test coverage report...${NC}"
go test ./internal/adapters/stripe/ -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

echo -e "${GREEN}Coverage report generated: coverage.html${NC}"
echo

# Clean up
rm -f coverage.out main

echo -e "${GREEN}All tests completed successfully!${NC}"
echo
echo -e "${GREEN}Next steps:${NC}"
echo "1. Review test coverage report"
echo "2. Configure STRIPE_TEST_SECRET_KEY for full integration testing"
echo "3. Test webhook endpoints with Stripe CLI"
echo "4. Verify payment flows in test environment"