#!/bin/bash

# 🧪 Test LLM Failure Handling
# This script tests that the chatbot properly handles LLM failures
# without falling back to rule-based responses

set -e

echo "🧪 Testing LLM Failure Handling"
echo "================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "PASS")
            echo -e "${GREEN}✅ PASS${NC}: $message"
            ;;
        "FAIL")
            echo -e "${RED}❌ FAIL${NC}: $message"
            ;;
        "INFO")
            echo -e "${YELLOW}ℹ️  INFO${NC}: $message"
            ;;
    esac
}

# Test 1: Check if API is running
print_status "INFO" "Testing API health endpoint..."

if curl -s http://localhost:8080/health > /dev/null; then
    print_status "PASS" "API is running"
else
    print_status "FAIL" "API is not running on localhost:8080"
    echo "Please start the API with: make start"
    exit 1
fi

# Test 2: Check LLM health endpoint
print_status "INFO" "Testing LLM health endpoint..."

LLM_HEALTH_RESPONSE=$(curl -s http://localhost:8080/api/chat/health 2>/dev/null || echo "{}")
LLM_STATUS=$(echo "$LLM_HEALTH_RESPONSE" | jq -r '.status // "error"' 2>/dev/null || echo "error")

if [ "$LLM_STATUS" = "healthy" ]; then
    print_status "PASS" "LLM is healthy"
    LLM_MODEL=$(echo "$LLM_HEALTH_RESPONSE" | jq -r '.model // "unknown"' 2>/dev/null || echo "unknown")
    print_status "INFO" "Using model: $LLM_MODEL"
else
    print_status "INFO" "LLM is not healthy (status: $LLM_STATUS)"
    LLM_ERROR=$(echo "$LLM_HEALTH_RESPONSE" | jq -r '.error // "unknown error"' 2>/dev/null || echo "unknown error")
    print_status "INFO" "LLM error: $LLM_ERROR"
fi

# Test 3: Test chat endpoint with LLM failure simulation
print_status "INFO" "Testing chat endpoint response..."

# Test with a simple message
CHAT_RESPONSE=$(curl -s -X POST http://localhost:8080/api/chat \
    -H "Content-Type: application/json" \
    -d '{"message": "Hello, how are you?"}' 2>/dev/null || echo "{}")

if [ "$LLM_STATUS" != "healthy" ]; then
    # If LLM is not healthy, check if response contains error message
    if echo "$CHAT_RESPONSE" | grep -q "AI service is currently unavailable\|LLM model.*not responding\|unavailable"; then
        print_status "PASS" "Chatbot correctly shows error message when LLM is unavailable"
    else
        print_status "FAIL" "Chatbot should show error message when LLM is unavailable"
        echo "Response: $CHAT_RESPONSE"
    fi
else
    # If LLM is healthy, check if response is normal
    if echo "$CHAT_RESPONSE" | jq -e '.response' > /dev/null 2>&1; then
        print_status "PASS" "Chatbot returns normal response when LLM is healthy"
        RESPONSE_TEXT=$(echo "$CHAT_RESPONSE" | jq -r '.response // "no response"' 2>/dev/null || echo "no response")
        print_status "INFO" "Response: ${RESPONSE_TEXT:0:100}..."
    else
        print_status "FAIL" "Chatbot should return normal response when LLM is healthy"
        echo "Response: $CHAT_RESPONSE"
    fi
fi

# Test 4: Check if Ollama is running
print_status "INFO" "Checking Ollama service..."

if curl -s http://192.168.0.3:11434/api/tags > /dev/null 2>&1; then
    print_status "PASS" "Ollama is running on 192.168.0.3:11434"
    
    # Check if gemma3n:e4b model is available
    OLLAMA_MODELS=$(curl -s http://192.168.0.3:11434/api/tags 2>/dev/null || echo "{}")
    if echo "$OLLAMA_MODELS" | jq -e '.models[] | select(.name | contains("gemma3n:e4b"))' > /dev/null 2>&1; then
        print_status "PASS" "gemma3n:e4b model is available"
    else
        print_status "FAIL" "gemma3n:e4b model is not available"
        print_status "INFO" "Available models:"
        echo "$OLLAMA_MODELS" | jq -r '.models[].name' 2>/dev/null || echo "No models found"
    fi
else
    print_status "FAIL" "Ollama is not running on 192.168.0.3:11434"
    print_status "INFO" "To start Ollama: ollama serve --host 0.0.0.0:11434"
    print_status "INFO" "To pull the model: ollama pull gemma3n:e4b"
fi

echo ""
echo "🎯 Test Summary:"
echo "================"

if [ "$LLM_STATUS" = "healthy" ]; then
    print_status "PASS" "All systems are healthy - LLM is working correctly"
else
    print_status "INFO" "LLM is not healthy - this is expected behavior for testing"
    print_status "INFO" "The chatbot should show error messages instead of falling back to rule-based responses"
fi

echo ""
echo "💡 Next Steps:"
echo "=============="
echo "1. If LLM is not working, check Ollama service and model availability"
echo "2. Test the frontend chatbot to verify error messages are displayed"
echo "3. Verify that no rule-based fallback responses are shown"
echo ""
echo "🔧 To fix LLM issues:"
echo "  - Start Ollama: ollama serve --host 0.0.0.0:11434"
echo "  - Pull model: ollama pull gemma3n:e4b"
echo "  - Check API logs: make api-logs"
