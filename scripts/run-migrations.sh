#!/bin/bash

# Bruno Site Database Migration Script
# This script runs database migrations on the PostgreSQL pod

set -e  # Exit on any error

NAMESPACE="bruno"
DATABASE="bruno-site"
USER="postgres"

echo "🗄️ Starting database migration process..."

# Get the postgres pod name
echo "🔍 Finding PostgreSQL pod..."
POSTGRES_POD=$(kubectl get pods -n $NAMESPACE -l app=postgres -o jsonpath='{.items[0].metadata.name}')
echo "✅ Found PostgreSQL pod: $POSTGRES_POD"

# Wait for postgres to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod/$POSTGRES_POD -n $NAMESPACE --timeout=300s
echo "✅ PostgreSQL is ready!"

# Test database connection
echo "🔍 Testing database connection..."
kubectl exec -n $NAMESPACE $POSTGRES_POD -- psql -U $USER -d $DATABASE -c "SELECT version();"
echo "✅ Database connection successful!"

# Copy migration files to postgres pod
echo "📋 Copying migration files to PostgreSQL pod..."
kubectl cp ./api/migrations $POSTGRES_POD:/tmp/migrations -n $NAMESPACE
echo "✅ Migration files copied successfully!"

# Run migrations in order
echo "🗄️ Running database migrations..."

echo "📋 Running initial schema migration..."
kubectl exec -n $NAMESPACE $POSTGRES_POD -- psql -U $USER -d $DATABASE -f /tmp/migrations/001_initial_schema.sql

echo "📋 Running data population migration..."
kubectl exec -n $NAMESPACE $POSTGRES_POD -- psql -U $USER -d $DATABASE -f /tmp/migrations/002_populate_data.sql

echo "📋 Running project active migration..."
kubectl exec -n $NAMESPACE $POSTGRES_POD -- psql -U $USER -d $DATABASE -f /tmp/migrations/003_add_project_active.sql

echo "✅ Database migrations completed!"

# Verify migration results
echo "🔍 Verifying migration results..."
kubectl exec -n $NAMESPACE $POSTGRES_POD -- psql -U $USER -d $DATABASE -c "\dt"

echo "📊 Checking data counts..."
kubectl exec -n $NAMESPACE $POSTGRES_POD -- psql -U $USER -d $DATABASE -c "SELECT COUNT(*) as projects_count FROM projects;"
kubectl exec -n $NAMESPACE $POSTGRES_POD -- psql -U $USER -d $DATABASE -c "SELECT COUNT(*) as skills_count FROM skills;"
kubectl exec -n $NAMESPACE $POSTGRES_POD -- psql -U $USER -d $DATABASE -c "SELECT COUNT(*) as experience_count FROM experience;"

echo "🎉 Database migration completed successfully!"
echo "🔄 You may need to restart the API deployment to pick up the database changes:"
echo "   kubectl rollout restart deployment/api -n bruno"
