#!/bin/bash

# 🔐 Bruno Site - Secure Password Generator
# This script generates secure passwords for the Bruno Site application

set -e

echo "🔐 Bruno Site - Secure Password Generator"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to generate secure password
generate_password() {
    local length=${1:-32}
    local password=$(openssl rand -base64 $length | tr -d "=+/" | cut -c1-$length)
    echo "$password"
}

# Function to check if openssl is available
check_openssl() {
    if ! command -v openssl &> /dev/null; then
        echo -e "${RED}❌ Error: openssl is not installed${NC}"
        echo "Please install openssl to generate secure passwords:"
        echo "  Ubuntu/Debian: sudo apt-get install openssl"
        echo "  macOS: brew install openssl"
        echo "  CentOS/RHEL: sudo yum install openssl"
        exit 1
    fi
}

# Function to display password requirements
show_requirements() {
    echo -e "${BLUE}📋 Password Requirements:${NC}"
    echo "  • Minimum 16 characters"
    echo "  • Mix of uppercase, lowercase, numbers, and special characters"
    echo "  • Unique for each service"
    echo "  • No common words or patterns"
    echo ""
}

# Function to generate and display passwords
generate_all_passwords() {
    echo -e "${GREEN}🔑 Generating secure passwords...${NC}"
    echo ""
    
    # Generate passwords
    POSTGRES_PASSWORD=$(generate_password 32)
    METRICS_PASSWORD=$(generate_password 32)
    
    echo -e "${YELLOW}📝 Copy these passwords to your .env file:${NC}"
    echo ""
    echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD"
    echo "METRICS_PASSWORD=$METRICS_PASSWORD"
    echo ""
    
    # Save to file
    cat > .env.generated << EOF
# 🔐 Bruno Site - Generated Passwords
# Generated on: $(date)
# WARNING: Keep this file secure and delete after use!

# Database Configuration
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
DATABASE_URL=postgresql://postgres:\${POSTGRES_PASSWORD}@localhost:5432/bruno_site?sslmode=disable

# Security Configuration
METRICS_USERNAME=metrics_admin
METRICS_PASSWORD=$METRICS_PASSWORD
METRICS_ENABLED=true

# Other Configuration
PORT=8080
GEMMA_MODEL=gemma3n:e4b
OLLAMA_URL=http://192.168.0.3:11434
EOF
    
    echo -e "${GREEN}✅ Passwords saved to .env.generated${NC}"
    echo -e "${YELLOW}⚠️  IMPORTANT:${NC}"
    echo "  1. Copy the passwords to your .env file"
    echo "  2. Delete .env.generated after use"
    echo "  3. Never commit passwords to version control"
    echo ""
}

# Function to validate existing .env file
validate_env_file() {
    if [ -f ".env" ]; then
        echo -e "${BLUE}🔍 Checking existing .env file...${NC}"
        
        # Check for default passwords
        if grep -q "secure-password\|secure_password_change_me\|your_secure" .env; then
            echo -e "${RED}❌ Found default passwords in .env file!${NC}"
            echo "Please update with secure passwords before running the application."
            echo ""
            return 1
        fi
        
        # Check for empty passwords
        if grep -q "POSTGRES_PASSWORD=$" .env || grep -q "METRICS_PASSWORD=$" .env; then
            echo -e "${RED}❌ Found empty passwords in .env file!${NC}"
            echo "Please set secure passwords before running the application."
            echo ""
            return 1
        fi
        
        echo -e "${GREEN}✅ .env file looks secure${NC}"
        echo ""
        return 0
    else
        echo -e "${YELLOW}⚠️  No .env file found${NC}"
        echo "Run this script to generate secure passwords."
        echo ""
        return 1
    fi
}

# Main script logic
main() {
    check_openssl
    show_requirements
    
    # Check if user wants to validate existing .env
    if [ "$1" = "--validate" ]; then
        validate_env_file
        exit $?
    fi
    
    # Check if .env already exists and is secure
    if validate_env_file; then
        echo -e "${YELLOW}Do you want to generate new passwords anyway? (y/N)${NC}"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            echo "Exiting..."
            exit 0
        fi
    fi
    
    generate_all_passwords
    
    echo -e "${GREEN}🎉 Password generation complete!${NC}"
    echo ""
    echo -e "${BLUE}Next steps:${NC}"
    echo "  1. Copy passwords from .env.generated to .env"
    echo "  2. Delete .env.generated"
    echo "  3. Start the application: docker-compose up -d"
    echo ""
    echo -e "${YELLOW}For more information, see:${NC}"
    echo "  📖 Security Setup Guide: ./SETUP_SECURITY.md"
    echo "  📖 Security Documentation: ./SECURITY.md"
}

# Run main function with all arguments
main "$@"
