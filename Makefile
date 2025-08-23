# Bruno Site Makefile
# This Makefile manages the complete Bruno site system using Docker Compose

# Environment configuration
ENV ?= dev
DOCKER_COMPOSE_FILE = docker-compose.yml
REGISTRY ?= ghcr.io/brunovlucena/bruno-site

.PHONY: help start stop restart build build-push logs clean status api-logs frontend-logs db-logs psql redis-cli api-shell frontend-shell frontend-dev migrate test-api test-api-unit test-frontend-unit test-e2e test-load test-coverage update-deps format lint pf-api pf-redis pf-postgres tp-intercept tp-intercept-with-mounts tp-stop tp-connect tp-disconnect tp-status tp-list restart-fresh reconcile

# Default target
help:
	@echo "🚀 Bruno Site Management"
	@echo "========================"
	@echo ""
	@echo "Environment: $(ENV)"
	@echo ""
	@echo "Available commands:"
	@echo "  make start      - Start all services (dev environment - Docker Compose)"
	@echo "  make stop       - Stop all services"
	@echo "  make restart    - Restart all services"
	@echo "  make build      - Build all Docker images (dev environment)"
	@echo "  make build-push - Build and push Docker images to ghcr.io/brunovlucena/bruno-site"
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
	@echo "  make restart-fresh  - Restart with fresh database (clean + start)"
	@echo "  make pf-api        - Port forward API service for local testing (Kubernetes)"
	@echo "  make tp-intercept  - Intercept both API and frontend services (no volume mounts)"
	@echo "  make tp-intercept-with-mounts - Intercept with volume mounts (requires sshfs)"
	@echo "  make tp-stop       - Stop all active intercepts"
	@echo "  make tp-disconnect - Disconnect from Kubernetes cluster"
	@echo "  make tp-status     - Show Telepresence status"
	@echo "  make tp-list       - List active intercepts"
	@echo "  make reconcile     - Reconcile Flux HelmRelease for bruno-site"
	@echo "  make test-api      - Test API endpoints"
	@echo "  make test          - Run all tests (API, frontend, E2E)"
	@echo "  make format        - Format code"
	@echo "  make lint          - Lint code"
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

# Build and push images to registry (latest tag only)
build-push:
	@echo "🏗️ Building and pushing Docker images..."
	@echo "Environment: $(ENV)"
	@echo "Registry: $(REGISTRY)"
	@echo "Tag: latest"
	@echo "📦 Building API image..."
	@docker build -t $(REGISTRY)/api:latest ./api
	@echo "📦 Building Frontend image..."
	@cp frontend/Dockerfile frontend/Dockerfile.temp
	@docker build -t $(REGISTRY)/frontend:latest ./frontend
	@mv frontend/Dockerfile.temp frontend/Dockerfile
	@echo "🚀 Pushing images to registry..."
	@docker push $(REGISTRY)/api:latest
	@docker push $(REGISTRY)/frontend:latest
	@echo "✅ Images built and pushed successfully!"
	@echo "📋 Pushed images:"
	@echo "  API: $(REGISTRY)/api:latest"
	@echo "  Frontend: $(REGISTRY)/frontend:latest"

# Show logs from all services
logs:
	@echo "📋 Following logs from all services (Ctrl+C to stop):"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

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
	@echo "⏳ Starting Vite dev server..."
	@cd frontend && npm run dev

# Port forward API service for local testing
pf-api:
	@echo "🚪 Port forwarding API service for local testing..."
	@echo "💡 Access API at http://localhost:8080"
	@echo "💡 Health check: http://localhost:8080/health"
	@echo "💡 API endpoints: http://localhost:8080/api/v1/*"
	@kubectl port-forward --address 0.0.0.0 -n bruno svc/bruno-site-api 8080:8080

# Port forward Redis service for local testing
pf-redis:
	@echo "🔴 Port forwarding Redis service for local testing..."
	@echo "💡 Access Redis at localhost:6379"
	@kubectl port-forward --address 0.0.0.0 -n bruno svc/bruno-site-redis 6379:6379

# Port forward both database services for local testing
pf-postgres:
	@echo "🗄️ Port forwarding database services for local testing..."
	@echo "💡 PostgreSQL: localhost:5432 (bruno_site/postgres/secure-password)"
	@echo "💡 Redis: localhost:6379"
	@kubectl port-forward --address 0.0.0.0 -n bruno svc/bruno-site-postgres 5432:5432 & \
	kubectl port-forward --address 0.0.0.0 -n bruno svc/bruno-site-redis 6379:6379 & \
	wait

# Telepresence intercept both services for local development
tp-intercept:
	@echo "🔗 Setting up Telepresence intercept for both API and frontend development..."
	@echo "💡 This will route traffic from K8s to your local services"
	@echo "💡 Make sure both services are running locally first"
	@echo "💡 API: make run-api (in one terminal)"
	@echo "💡 Frontend: make frontend-dev (in another terminal)"
	@echo "💡 Press Ctrl+C to stop Telepresence intercept"
	@echo "🔗 Starting API intercept..."
	@telepresence intercept bruno-site-api --port 8080:8080 --mechanism tcp --mount false &
	@echo "🔗 Starting Frontend intercept..."
	@telepresence intercept bruno-site-frontend --port 3000:80 --mechanism tcp --mount false &
	@echo "✅ Both intercepts started. Press Ctrl+C to stop all intercepts."
	@wait

# Telepresence intercept with volume mounts (requires sshfs)
tp-intercept-with-mounts:
	@echo "🔗 Setting up Telepresence intercept with volume mounts..."
	@echo "💡 This requires sshfs to be installed"
	@echo "💡 Make sure both services are running locally first"
	@echo "💡 API: make run-api (in one terminal)"
	@echo "💡 Frontend: make frontend-dev (in another terminal)"
	@echo "💡 Press Ctrl+C to stop Telepresence intercept"
	@echo "🔗 Starting API intercept with mounts..."
	@telepresence intercept bruno-site-api --port 8080:8080 --mechanism tcp &
	@echo "🔗 Starting Frontend intercept with mounts..."
	@telepresence intercept bruno-site-frontend --port 3000:80 --mechanism tcp &
	@echo "✅ Both intercepts started with volume mounts. Press Ctrl+C to stop all intercepts."
	@wait

# Stop all Telepresence intercepts
tp-stop:
	@echo "🛑 Stopping all Telepresence intercepts..."
	@telepresence leave bruno-site-api || true
	@telepresence leave bruno-site-frontend || true
	@echo "✅ All intercepts stopped"

# Telepresence connect to cluster
tp-connect:
	@echo "🔗 Connecting to Kubernetes cluster with Telepresence..."
	@telepresence connect

# Telepresence disconnect from cluster
tp-disconnect:
	@echo "🔗 Disconnecting from Kubernetes cluster..."
	@telepresence quit

# Telepresence status
tp-status:
	@echo "📊 Telepresence status:"
	@telepresence status

# List active Telepresence intercepts
tp-list:
	@echo "📋 Active Telepresence intercepts:"
	@telepresence list

# Reconcile Flux HelmRelease for bruno-site
reconcile:
	@echo "🔄 Reconciling Flux Git source for bruno-site..."
	@flux reconcile source git bruno-site -n flux-system
	@echo "✅ Git source reconciliation completed"
	@echo "🔄 Reconciling Flux HelmRelease for bruno-site..."
	@flux reconcile helmrelease bruno-site -n bruno
	@echo "✅ HelmRelease reconciliation completed"

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
	@curl -s http://localhost:8080/api/about | jq . || curl -s http://localhost:8080/api/about
	@echo ""
	@echo "Contact:"
	@curl -s http://localhost:8080/api/contact | jq . || curl -s http://localhost:8080/api/contact
	@echo "Chat Health:"
	@curl -s http://localhost:8080/api/chat/health | jq . || curl -s http://localhost:8080/api/chat/health
	@echo ""
	@echo "Chat (POST test):"
	@curl -s -X POST http://localhost:8080/api/chat -H "Content-Type: application/json" -d '{"message": "Hello"}' | jq . || curl -s -X POST http://localhost:8080/api/chat -H "Content-Type: application/json" -d '{"message": "Hello"}'

# Run all tests
test: test-api-unit test-frontend-unit test-e2e

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



# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@cd api && go test -v -race -coverprofile=coverage.out ./...
	@cd api && go tool cover -func=coverage.out
	@cd frontend && npm run test:coverage
	@echo "✅ Coverage reports generated"

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