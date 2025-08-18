# Bruno Site Makefile
# This Makefile manages the complete Bruno site system using Docker Compose

.PHONY: help start stop restart build logs clean status api-logs frontend-logs db-logs psql redis-cli api-shell frontend-shell

# Default target
help:
	@echo "🚀 Bruno Site Management"
	@echo "========================"
	@echo ""
	@echo "Available commands:"
	@echo "  make start      - Start all services"
	@echo "  make stop       - Stop all services"
	@echo "  make restart    - Restart all services"
	@echo "  make build      - Build all Docker images"
	@echo "  make logs       - Show logs from all services"
	@echo "  make api-logs   - Show API logs only"
	@echo "  make frontend-logs - Show frontend logs only"
	@echo "  make db-logs    - Show database logs only"
	@echo "  make status     - Show service status"
	@echo "  make clean      - Stop and remove all containers/volumes"
	@echo "  make psql       - Connect to PostgreSQL database"
	@echo "  make migrate    - Run database migration"
	@echo "  make redis-cli  - Connect to Redis CLI"
	@echo "  make api-shell  - Open shell in API container"
	@echo "  make frontend-shell - Open shell in frontend container"
	@echo ""

# Start services
start:
	@echo "🚀 Starting Bruno Site..."
	@docker-compose up --build -d
	# @echo "⏳ Waiting for services to be ready..."
	# @timeout 60 bash -c 'until docker exec postgres pg_isready -U bruno_user -d bruno_site; do sleep 2; done' || true
	# @timeout 30 bash -c 'until docker exec redis redis-cli ping; do sleep 2; done' || true
	# @timeout 60 bash -c 'until curl -f http://localhost:8080/health; do sleep 3; done' || true
	# @timeout 60 bash -c 'until curl -f http://localhost:3000; do sleep 3; done' || true
	@echo "✅ Bruno site is running!"
	@echo ""
	@echo "📋 Access Information:"
	@echo "  Frontend: http://localhost:3000"
	@echo "  API Health: http://localhost:8080/health"
	@echo "  API Projects: http://localhost:8080/api/v1/projects"
	@echo "  Grafana: http://localhost:3002 (admin/admin)"
	@echo "  Prometheus: http://localhost:9090"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Redis: localhost:6379"

# Stop services
stop:
	@echo "🛑 Stopping Bruno Site..."
	@docker-compose down
	@echo "✅ Services stopped"

# Restart services
restart: stop start

# Build all images
build:
	@echo "🏗️ Building Docker images..."
	@docker-compose build
	@echo "✅ Images built successfully"

# Show logs from all services
logs:
	@echo "📋 Recent logs from all services:"
	@docker-compose logs --tail=50

# Show API logs
api-logs:
	@echo "📋 API logs:"
	@docker logs -f api --tail=50

# Show frontend logs
frontend-logs:
	@echo "📋 Frontend logs:"
	@docker logs -f frontend --tail=50

# Show database logs
db-logs:
	@echo "📋 Database logs:"
	@docker logs -f postgres --tail=50

# Show service status
status:
	@echo "📊 Service Status:"
	@docker-compose ps
	# @echo ""
	# @echo "🔍 Health Checks:"
	# @echo "  PostgreSQL: $$(docker exec postgres pg_isready -U bruno_user -d bruno_site > /dev/null 2>&1 && echo "✅ Healthy" || echo "❌ Unhealthy")"
	# @echo "  Redis: $$(docker exec redis redis-cli ping > /dev/null 2>&1 && echo "✅ Healthy" || echo "❌ Unhealthy")"
	# @echo "  API: $$(curl -f http://localhost:8080/health > /dev/null 2>&1 && echo "✅ Healthy" || echo "❌ Unhealthy")"
	# @echo "  Frontend: $$(curl -f http://localhost:3000 > /dev/null 2>&1 && echo "✅ Healthy" || echo "❌ Unhealthy")"

# Clean everything (stop and remove containers, volumes, networks)
clean:
	@echo "🧹 Cleaning up Bruno Site..."
	@docker-compose down -v --remove-orphans
	@docker system prune -f
	@echo "✅ Cleanup completed"

# Connect to PostgreSQL
psql:
	@echo "🗄️ Connecting to PostgreSQL..."
	@docker exec -it postgres psql -U postgres -d bruno_site

# Run database migration
migrate:
	@echo "🗄️ Running database migration..."
	@PGPASSWORD=secure-password psql -h 127.0.0.1 -p 5432 -U postgres -d bruno_site < scripts/update_projects.sql

# Connect to Redis CLI
redis-cli:
	@echo "⚡ Connecting to Redis CLI..."
	@docker exec -it redis redis-cli

# Open shell in API container
api-shell:
	@echo "🔧 Opening shell in API container..."
	@docker exec -it api /bin/sh

# Open shell in frontend container
frontend-shell:
	@echo "🔧 Opening shell in frontend container..."
	@docker exec -it frontend /bin/sh

# Test API endpoints
test-api:
	@echo "🧪 Testing API endpoints..."
	@echo "Health check:"
	@curl -s http://localhost:8080/health | jq . || curl -s http://localhost:8080/health
	@echo ""
	@echo "Projects:"
	@curl -s http://localhost:8080/api/v1/projects | jq . || curl -s http://localhost:8080/api/v1/projects
	@echo ""
	@echo "About:"
	@curl -s http://localhost:8080/api/v1/about | jq . || curl -s http://localhost:8080/api/v1/about
	@echo ""
	@echo "Contact:"
	@curl -s http://localhost:8080/api/v1/contact | jq . || curl -s http://localhost:8080/api/v1/contact

# Watch logs from all services
watch-logs:
	@echo "👀 Watching logs from all services (Ctrl+C to stop):"
	@docker-compose logs -f

# Update dependencies
update-deps:
	@echo "📦 Updating dependencies..."
	@cd api && go mod tidy
	@cd frontend && npm update
	@echo "✅ Dependencies updated"

# Format code
format:
	@echo "🎨 Formatting code..."
	@cd api && go fmt ./...
	@cd frontend && npm run format 2>/dev/null || echo "No format script found in frontend"
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@cd api && go vet ./...
	@cd frontend && npm run lint 2>/dev/null || echo "No lint script found in frontend"
	@echo "✅ Code linted"