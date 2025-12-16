# Testing Guide

## Quick Start

Run all tests:

```bash
go test ./...
```

Run tests with detailed output:

```bash
go test ./... -v
```

Run specific layer tests:

```bash
# Handler tests only
go test ./internal/handler -v

# Service tests only
go test ./internal/service -v

# Repository integration tests
go test ./internal/repository -v

# Response utility tests
go test ./internal/response -v
```

## Test Coverage Breakdown

### Handler Tests (10 tests)

Located in `internal/handler/*_test.go`

These test HTTP endpoints and request/response handling:

- User registration and login
- Room CRUD operations
- Booking creation
- Payment processing
- Room maintenance requests

Run with:

```bash
go test ./internal/handler -v
```

### Service Tests (9 tests)

Located in `internal/service/*_test.go`

These test business logic and validation:

- User service: account creation, login logic
- Room service: room management, availability
- Booking service: booking validation
- Payment service: payment processing
- Room maintenance service: maintenance request handling

Run with:

```bash
go test ./internal/service -v
```

### Repository Integration Tests (5 tests)

Located in `internal/repository/*_integration_test.go`

These test database operations with a real PostgreSQL database.

**Prerequisites**: Set up PostgreSQL and configure connection string

```bash
# Option 1: Using DB_URL environment variable
export DB_URL="postgres://user:password@localhost:5432/industry_api_test"
go test ./internal/repository -v

# Option 2: Using TEST_DB_URL environment variable
export TEST_DB_URL="postgres://user:password@localhost:5432/industry_api_test"
go test ./internal/repository -v
```

If database is not configured, tests will automatically skip with a message like:

```
--- SKIP: TestUserRepository_CreateAndGet (0.00s)
    user_repo_integration_test.go:23: DB_URL or TEST_DB_URL not set; skipping integration tests
```

### Infrastructure Tests (2 tests)

Located in `internal/response/response_test.go` and `internal/service/login_test.go`

These test response formatting and JWT token generation.

Run with:

```bash
go test ./internal/response -v
go test ./internal/service -v -run TestGenerateJWT
```

## Running Tests with Different Options

### Run tests in parallel:

```bash
go test ./... -parallel 4
```

### Run only tests matching a pattern:

```bash
# Run only user-related tests
go test ./... -run User -v

# Run only validation tests
go test ./... -run Validation -v

# Run only success path tests
go test ./... -run Success -v
```

### Generate coverage report:

```bash
# Text report to console
go test ./... -cover

# HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
# Then open coverage.html in browser
```

### Run tests with timeout:

```bash
go test ./... -timeout 30s
```

### Run tests without cache:

```bash
go test ./... -count=1
```

### Verbose output with race condition detection:

```bash
go test ./... -v -race
```

## Test Structure

Each test follows this pattern:

1. **Setup**: Create mock dependencies
2. **Execute**: Call the function being tested
3. **Assert**: Verify the results

Example from handler test:

```go
func TestAddRoomHandler_Success(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    mr := &mockRoomSvcRepo{
        add: func(ctx context.Context, room *models.Room) error {
            room.ID = 1
            return nil
        },
        list: func(ctx context.Context) ([]*models.Room, error) {
            return nil, errors.New("not-impl")
        },
    }
    h := NewRoomHandler(service.NewRoomService(mr))

    // Execute
    reqBody := models.RoomRequest{...}
    b, _ := json.Marshal(reqBody)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request, _ = http.NewRequest("POST", "/api/v1/rooms/add", bytes.NewBuffer(b))
    h.AddRoom(c)

    // Assert
    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d", w.Code)
    }
}
```

## Mocking Strategy

Tests use lightweight mock repositories that implement service interfaces:

```go
type mockUserRepo struct {
    GetUserByEmailFn func(ctx context.Context, email string) (*models.User, error)
    CreateUserFn     func(ctx context.Context, user *models.User) error
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    return m.GetUserByEmailFn(ctx, email)
}
```

This approach:

- Requires no external mock libraries
- Keeps tests simple and readable
- Makes dependencies explicit
- Allows easy test-specific behavior

## Common Test Issues & Solutions

### Issue: Integration tests timeout or hang

**Solution**: Check database connection

```bash
# Test PostgreSQL connection
psql -h localhost -U postgres -d postgres -c "SELECT 1"
```

### Issue: Tests fail due to JWT errors

**Solution**: Ensure JWT_SECRET is set

```bash
export JWT_SECRET="your-secret-key"
go test ./... -v
```

### Issue: Port already in use (handler tests)

**Solution**: Uses random ports via httptest, shouldn't be an issue. If it occurs:

```bash
# Restart the test
go clean -testcache
go test ./... -v
```

### Issue: Data not cleaning up between tests

**Solution**: Integration tests include cleanup in defer blocks

```go
defer func() {
    // Cleanup database state
}()
```

## Continuous Integration (CI) Setup

To run tests in CI/CD pipeline:

```yaml
# Example GitHub Actions workflow
- name: Run tests
  run: |
    export JWT_SECRET="test-secret"
    go test ./... -v -race -timeout 10s -coverprofile=coverage.out
```

For integration tests:

```yaml
# Add PostgreSQL service
services:
  postgres:
    image: postgres:latest
    env:
      POSTGRES_PASSWORD: postgres
    options: >-
      --health-cmd pg_isready
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5

- name: Run integration tests
  env:
    TEST_DB_URL: "postgres://postgres:postgres@localhost:5432/test_db"
  run: go test ./internal/repository -v
```

## Test Best Practices Used

1. **Isolated Tests**: Each test is independent and doesn't rely on other tests
2. **Clear Names**: Test names describe what is being tested (TestFeature_Scenario)
3. **Mock Dependencies**: Services are tested with mock repositories
4. **Error Paths**: Tests cover both success and error cases
5. **Cleanup**: Resources are cleaned up after each test
6. **No Global State**: Tests don't depend on global variables

## Adding New Tests

When adding a new handler, service, or repository:

1. **Create test file**: `your_feature_test.go` in same directory
2. **Import testing package**: `import "testing"`
3. **Follow naming**: `TestYourFeature_Scenario`
4. **Mock dependencies**: Create simple mock structs
5. **Run tests**: `go test ./...`
6. **Check coverage**: `go test ./... -cover`

Example template:

```go
package handler

import (
    "testing"
    // other imports
)

func TestYourHandler_Scenario(t *testing.T) {
    // Setup
    // Execute
    // Assert
}
```

## Performance Testing

For performance-critical code:

```bash
# Run benchmarks
go test ./... -bench=. -benchmem

# Run specific benchmark
go test ./internal/service -bench=BenchmarkAddUser -benchmem
```

## Debugging Failed Tests

To debug a failing test:

```bash
# Run with print statements visible
go test -v ./path/to/test -run TestName

# With race detector
go test -race ./path/to/test -run TestName

# With CPU profiling
go test ./path/to/test -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## Resources

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Gin Testing](https://github.com/gin-gonic/gin#testing)
- [Go Test Flags](https://golang.org/cmd/go/#hdr-Testing_flags)
- [httptest Package](https://golang.org/pkg/net/http/httptest/)
