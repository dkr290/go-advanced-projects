#!/bin/bash

set -e

BASE_URL="${1:-http://localhost:8080}"
API_URL="${BASE_URL}/api/v1/simpleapis"

echo "Testing CRD API Deployer at: ${BASE_URL}"
echo "=========================================="

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

test_endpoint() {
    local method=$1
    local url=$2
    local data=$3
    local expected_status=$4
    local test_name=$5
    
    echo -e "\n${BLUE}Testing: ${test_name}${NC}"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "${method}" "${url}" \
            -H "Content-Type: application/json" \
            -d "${data}")
    else
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "${method}" "${url}")
    fi
    
    http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d: -f2)
    response_body=$(echo "$response" | sed '/HTTP_CODE:/d')
    
    echo "Response Code: ${http_code}"
    echo "Response Body: ${response_body}"
    
    if [ "$http_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ Test passed${NC}"
    else
        echo -e "\033[0;31m✗ Test failed (expected ${expected_status}, got ${http_code})\033[0m"
    fi
}

echo -e "${YELLOW}Checking if API is ready...${NC}"
for i in {1..30}; do
    if curl -s "${BASE_URL}/docs" > /dev/null 2>&1; then
        echo -e "${GREEN}API is ready!${NC}"
        break
    fi
    echo "Waiting for API... (${i}/30)"
    sleep 2
done

echo -e "\n${YELLOW}=== Test 1: Create minimal SimpleAPI ===${NC}"
test_endpoint "POST" "${API_URL}" '{
  "kind": "Simpleapi",
  "name": "test-minimal",
  "namespace": "default", 
  "image": "nginx",
  "version": "latest"
}' "201" "Create minimal SimpleAPI"

echo -e "\n${YELLOW}=== Test 2: Get API documentation ===${NC}"
test_endpoint "GET" "${BASE_URL}/docs" "" "200" "Get API documentation"

echo -e "\n${GREEN}=========================================="
echo -e "Basic tests completed!${NC}"
echo -e "\n${BLUE}You can view the API documentation at: ${BASE_URL}/docs${NC}"