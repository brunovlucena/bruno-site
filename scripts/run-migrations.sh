#!/bin/bash

# Bruno Site Database Migration Script
# This script runs migrations against the local PostgreSQL database

set -e

# Database configuration
DB_HOST="127.0.0.1"
DB_PORT="5432"
DB_NAME="bruno_site"
DB_USER="postgres"
DB_PASSWORD="secure-password"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 Bruno Site Database Migration${NC}"
echo "=================================="

# Check if PostgreSQL is running
echo -e "${YELLOW}📋 Checking if PostgreSQL is running...${NC}"
if ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER > /dev/null 2>&1; then
    echo -e "${RED}❌ PostgreSQL is not running on $DB_HOST:$DB_PORT${NC}"
    echo -e "${YELLOW}💡 Make sure to start the services with: make start${NC}"
    exit 1
fi

echo -e "${GREEN}✅ PostgreSQL is running${NC}"

# Set PGPASSWORD environment variable
export PGPASSWORD=$DB_PASSWORD

# Run migrations in order
echo -e "${YELLOW}📦 Running migrations...${NC}"

# Migration 1: Initial schema
echo -e "${YELLOW}📋 Running initial schema migration...${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f api/migrations/001_initial_schema.sql
echo -e "${GREEN}✅ Initial schema migration completed${NC}"

# Migration 2: Populate data
echo -e "${YELLOW}📋 Running data population migration...${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f api/migrations/002_populate_data.sql
echo -e "${GREEN}✅ Data population migration completed${NC}"

# Migration 3: Add project active column
echo -e "${YELLOW}📋 Running project active migration...${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f api/migrations/003_add_project_active.sql
echo -e "${GREEN}✅ Project active migration completed${NC}"

# Verify the data
echo -e "${YELLOW}🔍 Verifying data...${NC}"
echo -e "${YELLOW}📊 Projects:${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT id, title, type, featured, \"order\" FROM projects ORDER BY \"order\";"

echo -e "${YELLOW}📊 Skills count:${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) as skills_count FROM skills;"

echo -e "${YELLOW}📊 Experience count:${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) as experience_count FROM experience;"

echo -e "${GREEN}✅ All migrations completed successfully!${NC}"
echo -e "${GREEN}🎉 Database is ready for Bruno Site${NC}"

# Unset PGPASSWORD for security
unset PGPASSWORD
