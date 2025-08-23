#!/bin/bash

# 🤖 LLM Chatbot Test Script
# This script tests the new LLM-powered chatbot functionality

set -e

API_BASE_URL="http://localhost:8080"
HEALTH_ENDPOINT="$API_BASE_URL/api/chat/health"
CHAT_ENDPOINT="$API_BASE_URL/api/chat"

echo "🚀 Testing LLM Chatbot Integration"
echo "=================================="

# Test 1: Health Check
echo "📊 1. Testing LLM Health Check..."
if curl -s -f "$HEALTH_ENDPOINT" > /dev/null; then
    echo "✅ LLM Health Check: PASSED"
    curl -s "$HEALTH_ENDPOINT" | jq '.'
else
    echo "❌ LLM Health Check: FAILED"
    echo "💡 Make sure Ollama/LMStudio is running and the API server is started"
    exit 1
fi

echo ""

# Test 2: Simple Chat Request
echo "💬 2. Testing Simple Chat Request..."
RESPONSE=$(curl -s -X POST "$CHAT_ENDPOINT" \
    -H "Content-Type: application/json" \
    -d '{"message": "Hello, tell me about Bruno"}')

if echo "$RESPONSE" | jq -e '.response' > /dev/null 2>&1; then
    echo "✅ Simple Chat: PASSED"
    echo "Response: $(echo "$RESPONSE" | jq -r '.response' | head -c 100)..."
else
    echo "❌ Simple Chat: FAILED"
    echo "Response: $RESPONSE"
fi

echo ""

# Test 3: Skills Query
echo "🛠️ 3. Testing Skills Query..."
RESPONSE=$(curl -s -X POST "$CHAT_ENDPOINT" \
    -H "Content-Type: application/json" \
    -d '{"message": "What are Bruno'\''s key skills in cloud technologies?"}')

if echo "$RESPONSE" | jq -e '.response' > /dev/null 2>&1; then
    echo "✅ Skills Query: PASSED"
    echo "Response: $(echo "$RESPONSE" | jq -r '.response' | head -c 100)..."
else
    echo "❌ Skills Query: FAILED"
    echo "Response: $RESPONSE"
fi

echo ""

# Test 4: Experience Query
echo "💼 4. Testing Experience Query..."
RESPONSE=$(curl -s -X POST "$CHAT_ENDPOINT" \
    -H "Content-Type: application/json" \
    -d '{"message": "Tell me about Bruno'\''s experience with Kubernetes"}')

if echo "$RESPONSE" | jq -e '.response' > /dev/null 2>&1; then
    echo "✅ Experience Query: PASSED"
    echo "Response: $(echo "$RESPONSE" | jq -r '.response' | head -c 100)..."
else
    echo "❌ Experience Query: FAILED"
    echo "Response: $RESPONSE"
fi

echo ""

# Test 5: Contact Query
echo "📞 5. Testing Contact Query..."
RESPONSE=$(curl -s -X POST "$CHAT_ENDPOINT" \
    -H "Content-Type: application/json" \
    -d '{"message": "How can I contact Bruno?"}')

if echo "$RESPONSE" | jq -e '.response' > /dev/null 2>&1; then
    echo "✅ Contact Query: PASSED"
    echo "Response: $(echo "$RESPONSE" | jq -r '.response' | head -c 100)..."
else
    echo "❌ Contact Query: FAILED"
    echo "Response: $RESPONSE"
fi

echo ""
echo "🎉 LLM Chatbot Test Complete!"
echo ""
echo "📋 Summary:"
echo "- Health Check: API connectivity"
echo "- Simple Chat: Basic LLM functionality"  
echo "- Skills Query: PostgreSQL skills data integration"
echo "- Experience Query: PostgreSQL experience data integration"
echo "- Contact Query: PostgreSQL contact data integration"
echo ""
echo "💡 If any tests failed, check:"
echo "   1. API server is running (port 8080)"
echo "   2. PostgreSQL is running with data"
echo "   3. Ollama/LMStudio is running with Gemma3 model"
echo "   4. Environment variables are configured correctly"
