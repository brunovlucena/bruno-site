# Bruno Site Makefile
# This Makefile manages the complete Bruno site system using Docker Compose

# Environment configuration
ENV ?= dev
DOCKER_COMPOSE_FILE = docker-compose.yml

.PHONY: help start stop restart build build-dev build-prd logs clean status api-logs frontend-logs db-logs psql redis-cli api-shell frontend-shell frontend-dev frontend-dev-stop pf-api

# Default target
help:
	@echo "🚀 Bruno Site Management"
	@echo "========================"
	@echo ""
	@echo "Environment: $(ENV)"
	@echo ""
	@echo "Available commands:"
	@echo "  make start      - Start all services (dev environment - Docker Compose)"
	@echo "  make start-prd  - Show production deployment info (Kubernetes)"
	@echo "  make stop       - Stop all services"
	@echo "  make restart    - Restart all services"
	@echo "  make build      - Build all Docker images (dev environment)"
	@echo "  make build-prd  - Show production build info (Kubernetes)"
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
	@echo "  make frontend-dev   - Run frontend in development mode (hot reload)"
	@echo "  make frontend-dev-stop - Stop frontend dev gracefully (prevents browser flicker)"
	@echo "  make restart-fresh  - Restart with fresh database (clean + start)"
	@echo "  make pf-api        - Port forward API service for local testing (Kubernetes)"
	@echo ""

# Start services (development)
start:
	@echo "🚀 Starting Bruno Site (Development)..."
	@echo "Environment: $(ENV)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d
	@echo "⏳ Waiting for database to be ready..."
	@timeout 60 bash -c 'until docker exec postgres pg_isready -U postgres -d bruno_site; do sleep 2; done' || true
	@echo "🗄️ Running database migrations..."
	@make migrate || echo "⚠️ Migration failed, but continuing..."
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
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "✅ Services stopped"

# Show production deployment info (Kubernetes)
start-prd:
	@echo "🚀 Production Deployment Information (Kubernetes)"
	@echo "================================================="
	@echo "📋 Production deployment uses Kubernetes with Helm charts"
	@echo "📁 Chart location: ./chart/"
	@echo "🔧 Deployment process:"
	@echo "  1. Ensure Kubernetes cluster is running"
	@echo "  2. Install/update Helm chart"
	@echo "  3. Monitor deployment status"
	@echo ""
	@echo "🚀 To deploy to production:"
	@echo "  helm upgrade --install bruno-site ./chart --namespace bruno"
	@echo ""
	@echo "📊 To check deployment status:"
	@echo "  kubectl get pods -n bruno"
	@echo "  kubectl get services -n bruno"
	@echo ""
	@echo "🔍 To view logs:"
	@echo "  kubectl logs -f deployment/bruno-api -n bruno"
	@echo "  kubectl logs -f deployment/bruno-frontend -n bruno"

# Restart services
restart: stop start

# Restart with fresh database (clean and start)
restart-fresh: clean start

# Build all images (development)
build:
	@echo "🏗️ Building Docker images (Development)..."
	@echo "Environment: $(ENV)"
	@cp frontend/Dockerfile.dev frontend/Dockerfile
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build
	@echo "✅ Images built successfully"

# Show production build info (Kubernetes)
build-prd:
	@echo "🏗️ Production Build Information (Kubernetes)"
	@echo "============================================="
	@echo "📋 Production deployment uses Kubernetes with Helm charts"
	@echo "📁 Chart location: ./chart/"
	@echo "🔧 Build process:"
	@echo "  1. Docker images are built with production Dockerfiles"
	@echo "  2. Images are pushed to container registry"
	@echo "  3. Kubernetes manifests are applied via Helm"
	@echo ""
	@echo "📦 To build for production:"
	@echo "  docker build -f api/Dockerfile -t your-registry/bruno-api:prd ./api"
	@echo "  docker build -f frontend/Dockerfile -t your-registry/bruno-frontend:prd ./frontend"
	@echo ""
	@echo "🚀 To deploy to production:"
	@echo "  helm upgrade --install bruno-site ./chart --namespace bruno"

# Show logs from all services
logs:
	@echo "📋 Recent logs from all services:"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs --tail=50

# Show API logs
api-logs:
	@echo "📋 API logs:"
	@docker logs -f bruno-api --tail=50

# Show frontend logs
frontend-logs:
	@echo "📋 Frontend logs:"
	@docker logs -f bruno-frontend --tail=50

# Show database logs
db-logs:
	@echo "📋 Database logs:"
	@docker logs -f bruno-postgres --tail=50

# Show service status
status:
	@echo "📊 Service Status:"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) ps
	# @echo ""
	# @echo "🔍 Health Checks:"
	# @echo "  PostgreSQL: $$(docker exec postgres pg_isready -U bruno_user -d bruno_site > /dev/null 2>&1 && echo "✅ Healthy" || echo "❌ Unhealthy")"
	# @echo "  Redis: $$(docker exec redis redis-cli ping > /dev/null 2>&1 && echo "✅ Healthy" || echo "❌ Unhealthy")"
	# @echo "  API: $$(curl -f http://localhost:8080/health > /dev/null 2>&1 && echo "✅ Healthy" || echo "❌ Unhealthy")"
	# @echo "  Frontend: $$(curl -f http://localhost:3000 > /dev/null 2>&1 && echo "✅ Healthy" || echo "❌ Unhealthy")"

# Clean everything (stop and remove containers, volumes, networks)
clean:
	@echo "🧹 Cleaning up Bruno Site..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans
	@docker system prune -f
	@echo "✅ Cleanup completed"

# Connect to PostgreSQL
psql:
	@echo "🗄️ Connecting to PostgreSQL..."
	@docker exec -it bruno-postgres psql -U postgres -d bruno_site

# Run database migration
migrate:
	@echo "🗄️ Running database migration..."
	@chmod +x scripts/run-migrations.sh
	@./scripts/run-migrations.sh

# Connect to Redis CLI
redis-cli:
	@echo "⚡ Connecting to Redis CLI..."
	@docker exec -it bruno-redis redis-cli

# Open shell in API container
api-shell:
	@echo "🔧 Opening shell in API container..."
	@docker exec -it bruno-api /bin/sh

# Open shell in frontend container
frontend-shell:
	@echo "🔧 Opening shell in frontend container..."
	@docker exec -it bruno-frontend /bin/sh

# Run frontend in development mode (hot reload)
frontend-dev:
	@echo "🚀 Starting frontend in development mode..."
	@echo "📋 This will run the frontend with hot reload on http://localhost:5173"
	@echo "🔗 It will connect to the API running in Docker on http://localhost:8080"
	@echo "💡 To stop gracefully, use: make frontend-dev-stop"
	@echo "⏳ Starting Vite dev server..."
	@cd frontend && npm run dev

# Stop frontend dev gracefully (prevents browser flicker)
frontend-dev-stop:
	@echo "🛑 Stopping frontend dev gracefully..."
	@echo "💡 This prevents the browser from flickering when stopping the dev server"
	@echo "📋 Close your browser tab first, then press Ctrl+C in the terminal"
	@echo "⏳ Or use Ctrl+C directly in the terminal where frontend-dev is running"
	@echo ""
	@echo "🔍 If you need to force kill the process:"
	@echo "   pkill -f 'vite'"
	@echo "   or"
	@echo "   lsof -ti:5173 | xargs kill -9"

# Test API endpoints
test-api:
	@echo "🧪 Testing API endpoints..."
	@echo "Health check:"
	@curl -s http://localhost:8080/health | jq . || curl -s http://localhost:8080/health
	@echo ""
	@echo "Projects:"
	@curl -s http://localhost:8080/api/projects | jq . || curl -s http://localhost:8080/api/projects
	@echo ""
	@echo "About:"
	@curl -s http://localhost:8080/api/about | jq . || curl -s http://localhost:8080/api/about
	@echo ""
	@echo "Contact:"
	@curl -s http://localhost:8080/api/contact | jq . || curl -s http://localhost:8080/api/contact

# Setup domain and SSL
setup-domain:
	@echo "🌐 Setting up domain and SSL certificates..."
	@./scripts/setup-domain.sh

# Check SSL certificate status
check-ssl:
	@echo "🔒 Checking SSL certificate status..."
	@kubectl get certificate -n bruno
	@echo ""
	@echo "📋 Certificate details:"
	@kubectl describe certificate -n bruno bruno-site-tls || echo "Certificate not found yet"

# Port forward nginx-ingress for local testing
port-forward-nginx:
	@echo "🚪 Port forwarding nginx-ingress for local testing..."
	@echo "💡 Access your site at http://localhost (port 80) or https://localhost (port 443)"
	@echo "💡 Make sure to add 'localhost lucena.cloud' to your /etc/hosts file"
	@kubectl port-forward --address 0.0.0.0 -n nginx-ingress svc/nginx-ingress-ingress-nginx-controller 80:80 443:443

# Port forward API service for local testing
pf-api:
	@echo "🚪 Port forwarding API service for local testing..."
	@echo "💡 Access API at http://localhost:8080"
	@echo "💡 Health check: http://localhost:8080/health"
	@echo "💡 API endpoints: http://localhost:8080/api/v1/*"
	@kubectl port-forward --address 0.0.0.0 -n bruno svc/bruno-site-api 8080:8080

# Run all tests
test: test-api-unit test-frontend-unit test-e2e test-load

# Run API unit tests
test-api-unit:
	@echo "🧪 Running API unit tests..."
	@cd api && go test -v -race -coverprofile=coverage.out ./...
	@cd api && go tool cover -html=coverage.out -o coverage.html
	@echo "✅ API unit tests completed"

# Run frontend unit tests
test-frontend-unit:
	@echo "🧪 Running frontend unit tests..."
	@cd frontend && npm install --legacy-peer-deps && npm run test -- --run --coverage
	@echo "✅ Frontend unit tests completed"

# Run E2E tests
test-e2e:
	@echo "🧪 Running E2E tests..."
	@cd frontend && npm install --legacy-peer-deps && npm run test:e2e
	@echo "✅ E2E tests completed"

# Run load tests
test-load:
	@echo "🧪 Running load tests..."
	@k6 run tests/k6/load-test.js
	@echo "✅ Load tests completed"

# Run tests in watch mode
test-watch:
	@echo "🧪 Running tests in watch mode..."
	@cd frontend && npm run test:watch

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@cd api && go test -v -race -coverprofile=coverage.out ./...
	@cd api && go tool cover -func=coverage.out
	@cd frontend && npm run test:coverage
	@echo "✅ Coverage reports generated"

# Watch logs from all services
watch-logs:
	@echo "👀 Watching logs from all services (Ctrl+C to stop):"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

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