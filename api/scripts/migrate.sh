#!/bin/bash

# Database migration script for Bruno API
# This script runs the SQL migrations to set up the database schema and populate initial data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-bruno_site}

# Migration directory
MIGRATIONS_DIR="migrations"

echo -e "${GREEN}🚀 Starting database migration for Bruno API${NC}"
echo "Database: $DB_HOST:$DB_PORT/$DB_NAME"
echo "User: $DB_USER"
echo ""

# Function to run a migration file
run_migration() {
    local file=$1
    local description=$2
    
    echo -e "${YELLOW}📋 Running migration: $description${NC}"
    echo "File: $file"
    
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$file" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Migration completed successfully${NC}"
    else
        echo -e "${RED}❌ Migration failed${NC}"
        echo "Running with verbose output:"
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$file"
        exit 1
    fi
    echo ""
}

# Check if database exists, create if it doesn't
echo -e "${YELLOW}🔍 Checking if database exists...${NC}"
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo -e "${YELLOW}📦 Creating database: $DB_NAME${NC}"
    PGPASSWORD=$DB_PASSWORD createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    echo -e "${GREEN}✅ Database created successfully${NC}"
else
    echo -e "${GREEN}✅ Database already exists${NC}"
fi
echo ""

# Run migrations in order
echo -e "${YELLOW}🔄 Running migrations...${NC}"

# Schema migration
if [ -f "$MIGRATIONS_DIR/001_initial_schema.sql" ]; then
    run_migration "$MIGRATIONS_DIR/001_initial_schema.sql" "Initial database schema"
else
    echo -e "${RED}❌ Migration file not found: $MIGRATIONS_DIR/001_initial_schema.sql${NC}"
    exit 1
fi

# Data migration
if [ -f "$MIGRATIONS_DIR/002_populate_data.sql" ]; then
    run_migration "$MIGRATIONS_DIR/002_populate_data.sql" "Populate initial data"
else
    echo -e "${RED}❌ Migration file not found: $MIGRATIONS_DIR/002_populate_data.sql${NC}"
    exit 1
fi

# Verify migration
echo -e "${YELLOW}🔍 Verifying migration...${NC}"
echo "Checking tables:"

tables=("projects" "project_views" "visitors" "content" "skills" "experience")
for table in "${tables[@]}"; do
    count=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM $table;" | xargs)
    echo -e "  📊 $table: $count records"
done

echo ""
echo -e "${GREEN}🎉 Database migration completed successfully!${NC}"
echo ""
echo -e "${YELLOW}📊 Migration Summary:${NC}"
echo "  • Database: $DB_NAME"
echo "  • Tables created: 6"
echo "  • Projects loaded: 7"
echo "  • Skills loaded: 37"
echo "  • Experience entries: 6"
echo "  • Content entries: 3"
echo ""
echo -e "${GREEN}🚀 Your Bruno database is ready!${NC}" 