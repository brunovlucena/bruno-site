#!/bin/bash

# 🤖 Bruno Site - Gemma3n:e4b Setup Script
# This script helps set up the new Gemma3n:e4b model for the chatbot

set -e

echo "🚀 Setting up Gemma3n:e4b for Bruno Site Chatbot"
echo "================================================"

# Check if Ollama is installed
if ! command -v ollama &> /dev/null; then
    echo "❌ Ollama is not installed. Please install it first:"
    echo "   curl -fsSL https://ollama.ai/install.sh | sh"
    exit 1
fi

# Check if Ollama is running
if ! curl -s http://192.168.0.3:11434/api/tags &> /dev/null; then
    echo "❌ Ollama is not running on 192.168.0.3:11434. Please start it first:"
    echo "   ollama serve --host 0.0.0.0:11434"
    exit 1
fi

echo "✅ Ollama is running on 192.168.0.3:11434"

# Pull the Gemma3n:e4b model
echo "📥 Pulling Gemma3n:e4b model..."
ollama pull gemma3n:e4b

echo "✅ Model downloaded successfully!"

# Test the model
echo "🧪 Testing the model..."
echo "Testing with a simple prompt..."

# Create a temporary test file
cat > /tmp/test_prompt.txt << EOF
You are Bruno's AI assistant. Answer this question briefly: What is your name and what model are you using?
EOF

# Test the model
response=$(ollama run gemma3n:e4b < /tmp/test_prompt.txt)

echo "🤖 Model response:"
echo "$response"

# Clean up
rm -f /tmp/test_prompt.txt

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Start the Bruno site: make start"
echo "2. Test the chatbot at: http://localhost:3000"
echo "3. Check API health: http://localhost:8080/api/chat/health"
echo ""
echo "💡 The chatbot is now configured to use gemma3n:e4b"
