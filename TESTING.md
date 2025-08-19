# Bruno Site Testing Strategy

This document outlines the comprehensive testing strategy for the Bruno Site project, including unit tests, integration tests, E2E tests, and load tests.

## 🧪 Test Types

### 1. API Tests (Go)
- **Location**: `api/main_test.go`
- **Framework**: Go's built-in testing + testify
- **Coverage**: Unit tests for all endpoints, database operations, and error handling
- **Run**: `make test-api-unit`

**Features:**
- HTTP endpoint testing with httptest
- Database integration tests
- Mock handlers for isolated testing
- Benchmark tests for performance
- Race condition detection

### 2. Frontend Tests (React/TypeScript)
- **Location**: `frontend/src/test/`
- **Framework**: Vitest + React Testing Library
- **Coverage**: Component testing, API mocking, user interactions
- **Run**: `make test-frontend-unit`

**Features:**
- Component rendering tests
- User interaction testing
- API service mocking with MSW
- Accessibility testing
- Responsive design testing

### 3. E2E Tests (Playwright)
- **Location**: `tests/e2e/`
- **Framework**: Playwright
- **Coverage**: Full user journey testing across browsers
- **Run**: `make test-e2e`

**Features:**
- Cross-browser testing (Chrome, Firefox, Safari)
- Mobile device testing
- Visual regression testing
- Performance testing
- Accessibility validation

### 4. Load Tests (k6)
- **Location**: `tests/k6/load-test.js`
- **Framework**: k6
- **Coverage**: Performance and stress testing
- **Run**: `make test-load`

**Features:**
- Multi-stage load testing
- Performance thresholds
- Error rate monitoring
- Response time analysis

## 🚀 Running Tests

### Quick Start
```bash
# Run all tests
make test

# Run specific test types
make test-api-unit      # API unit tests
make test-frontend-unit # Frontend unit tests
make test-e2e          # E2E tests
make test-load         # Load tests

# Run tests with coverage
make test-coverage

# Run tests in watch mode
make test-watch
```

### Individual Commands
```bash
# API Tests
cd api
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Frontend Tests
cd frontend
npm run test
npm run test:coverage
npm run test:ui

# E2E Tests
cd frontend
npm run test:e2e
npm run test:e2e:ui

# Load Tests
k6 run tests/k6/load-test.js
```

## 📊 Test Coverage

### API Coverage
- **Endpoints**: 100% coverage for all HTTP endpoints
- **Database**: Integration tests for all database operations
- **Error Handling**: Comprehensive error scenario testing
- **Performance**: Benchmark tests for critical paths

### Frontend Coverage
- **Components**: Unit tests for all React components
- **User Interactions**: Click, type, and navigation testing
- **API Integration**: Mocked API calls and error handling
- **Accessibility**: ARIA compliance and keyboard navigation

### E2E Coverage
- **User Journeys**: Complete user workflows
- **Cross-browser**: Chrome, Firefox, Safari testing
- **Mobile**: Responsive design validation
- **Performance**: Page load and interaction timing

### Load Testing Coverage
- **Concurrent Users**: Up to 20 concurrent users
- **Response Times**: 95% under 500ms threshold
- **Error Rates**: Less than 10% error rate
- **Endpoints**: All critical API endpoints

## 🔧 Test Configuration

### API Test Setup
```go
// Test database connection
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("postgres", "postgres://postgres:secure-password@localhost:5432/bruno_site_test?sslmode=disable")
    if err != nil {
        t.Skipf("Skipping test - database not available: %v", err)
    }
    return db
}
```

### Frontend Test Setup
```typescript
// Vitest configuration
export default defineConfig({
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
    }
  }
})
```

### E2E Test Setup
```typescript
// Playwright configuration
export default defineConfig({
  testDir: './tests/e2e',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'firefox', use: { ...devices['Desktop Firefox'] } },
    { name: 'webkit', use: { ...devices['Desktop Safari'] } },
  ]
})
```

## 🏗️ CI/CD Integration

### GitHub Actions Workflow
The test suite is integrated into the CI/CD pipeline with the following jobs:

1. **test-api**: Go unit and integration tests
2. **test-frontend**: React unit tests and linting
3. **test-e2e**: Playwright E2E tests
4. **test-load**: k6 load testing
5. **security-scan**: Trivy vulnerability scanning

### Test Requirements
- All tests must pass before deployment
- Coverage reports are generated and uploaded
- Test artifacts are preserved for analysis
- Security scans are mandatory

## 📈 Performance Benchmarks

### API Performance
- Health check: < 200ms
- Projects endpoint: < 300ms
- Analytics tracking: < 400ms
- Metrics endpoint: < 500ms

### Frontend Performance
- Initial page load: < 2s
- Component rendering: < 100ms
- API calls: < 500ms
- Bundle size: < 500KB

### Load Test Thresholds
- 95% of requests under 500ms
- Error rate below 10%
- Support for 20 concurrent users
- Graceful degradation under load

## 🐛 Debugging Tests

### API Test Debugging
```bash
# Run specific test with verbose output
go test -v -run TestHealthEndpoint

# Run with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Test Debugging
```bash
# Run tests in watch mode
npm run test:watch

# Run specific test file
npm run test -- App.test.tsx

# Debug with UI
npm run test:ui
```

### E2E Test Debugging
```bash
# Run with UI mode
npm run test:e2e:ui

# Run specific test
npx playwright test home.spec.ts

# Debug mode
npx playwright test --debug
```

## 📝 Best Practices

### Writing Tests
1. **Arrange-Act-Assert**: Structure tests clearly
2. **Descriptive Names**: Use clear test names
3. **Isolation**: Each test should be independent
4. **Mocking**: Mock external dependencies
5. **Coverage**: Aim for high test coverage

### Test Data
1. **Fixtures**: Use consistent test data
2. **Factories**: Create test data factories
3. **Cleanup**: Clean up after tests
4. **Randomization**: Use random data when appropriate

### Performance
1. **Fast Tests**: Keep tests fast
2. **Parallel Execution**: Run tests in parallel
3. **Resource Management**: Clean up resources
4. **Caching**: Cache test dependencies

## 🔍 Monitoring and Reporting

### Coverage Reports
- HTML coverage reports generated
- Uploaded to Codecov for tracking
- Coverage thresholds enforced
- Trend analysis available

### Test Results
- Test results stored as artifacts
- Failure screenshots captured
- Performance metrics tracked
- Historical data maintained

### Alerts
- Test failure notifications
- Coverage drop alerts
- Performance regression alerts
- Security vulnerability alerts

## 🚨 Troubleshooting

### Common Issues
1. **Database Connection**: Ensure test database is running
2. **Port Conflicts**: Check for port conflicts
3. **Dependencies**: Install all test dependencies
4. **Environment**: Set correct environment variables

### Debug Commands
```bash
# Check test database
docker exec postgres pg_isready -U postgres -d bruno_site_test

# Check Redis
docker exec redis redis-cli ping

# Check API health
curl http://localhost:8080/health

# Check frontend
curl http://localhost:3000
```

This comprehensive testing strategy ensures the Bruno Site is robust, reliable, and performant across all environments.
