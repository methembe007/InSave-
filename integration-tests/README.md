# InSavein Platform - Integration Tests

This directory contains comprehensive integration tests for the InSavein platform that validate cross-service workflows and end-to-end functionality.

## Overview

The integration tests verify the following workflows:
- **User Registration Flow**: Registration → Profile → Preferences
- **Savings Flow**: Transaction → Streak → Notifications
- **Budget Alert Flow**: Spending → Threshold → Alert → Notification
- **Goal Progress Flow**: Create Goal → Contributions → Milestones

## Requirements Validated

### User Registration and Profile Management
- **Requirement 1.1**: User Registration and Authentication
- **Requirement 3.2**: User Profile Management
- **Requirement 3.3**: User Preferences

### Savings Transactions and Streaks
- **Requirement 4.1**: Savings Transaction Recording
- **Requirement 4.3**: Transaction Validation
- **Requirement 5.1**: Savings Streak Calculation

### Budget Alerts and Notifications
- **Requirement 7.1**: Spending Transaction Recording
- **Requirement 8.1**: Budget Alert Generation
- **Requirement 12.1**: Notification Delivery

### Goal Management and Milestones
- **Requirement 9.1**: Financial Goal Management
- **Requirement 10.1**: Goal Progress Tracking
- **Requirement 10.4**: Milestone Completion

## Test Environment

The integration tests use a dedicated test environment with:
- Isolated PostgreSQL test database
- All 8 microservices running in test mode
- Test data seeding for consistent test execution
- Docker Compose orchestration

## Prerequisites

- Docker and Docker Compose installed
- Go 1.21 or later
- All service Docker images built

## Setup

1. **Build service images** (if not already built):
   ```bash
   docker-compose build
   ```

2. **Install Go dependencies**:
   ```bash
   cd integration-tests
   go mod download
   ```

## Running Tests

### Start Test Environment

```bash
# From the integration-tests directory
docker-compose -f docker-compose.test.yml up -d

# Wait for services to be healthy (about 30 seconds)
docker-compose -f docker-compose.test.yml ps
```

### Run All Integration Tests

```bash
# From the integration-tests directory
go test -v ./...
```

### Run Specific Test Suite

```bash
# User registration flow
go test -v -run TestUserRegistrationFlow

# Savings flow
go test -v -run TestSavingsFlow

# Budget alert flow
go test -v -run TestBudgetAlertFlow

# Goal progress flow
go test -v -run TestGoalProgressFlow
```

### Run with Coverage

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Test Structure

```
integration-tests/
├── docker-compose.test.yml       # Test environment configuration
├── test-data/
│   └── init-test-db.sql         # Database schema and seed data
├── helpers/
│   ├── client.go                # HTTP client utilities
│   └── types.go                 # Shared type definitions
├── user_registration_test.go    # User registration flow tests
├── savings_flow_test.go         # Savings transaction tests
├── budget_alert_flow_test.go    # Budget alert tests
└── goal_progress_flow_test.go   # Goal progress tests
```

## Test Data

The test database is seeded with:
- 2 test users (test1@example.com, test2@example.com)
- Password: `TestPassword123!` (bcrypt hashed)
- Sample savings transactions
- Sample budget with categories
- Sample goal with milestones

## Service Endpoints (Test Environment)

- Auth Service: http://localhost:18080
- User Service: http://localhost:18081
- Savings Service: http://localhost:18082
- Budget Service: http://localhost:18083
- Goal Service: http://localhost:18005
- Notification Service: http://localhost:18086
- PostgreSQL: localhost:5435

## Cleanup

```bash
# Stop and remove test containers
docker-compose -f docker-compose.test.yml down

# Remove test volumes (clean slate)
docker-compose -f docker-compose.test.yml down -v
```

## Troubleshooting

### Services Not Starting

Check service logs:
```bash
docker-compose -f docker-compose.test.yml logs [service-name]
```

### Database Connection Issues

Verify PostgreSQL is healthy:
```bash
docker-compose -f docker-compose.test.yml ps postgres-test
```

### Test Failures

1. Ensure all services are healthy before running tests
2. Check service logs for errors
3. Verify test database is properly seeded
4. Run tests individually to isolate issues

## CI/CD Integration

These tests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run Integration Tests
  run: |
    docker-compose -f integration-tests/docker-compose.test.yml up -d
    sleep 30  # Wait for services
    cd integration-tests && go test -v ./...
    docker-compose -f integration-tests/docker-compose.test.yml down
```

## Best Practices

1. **Isolation**: Each test should be independent and not rely on other tests
2. **Cleanup**: Tests should clean up after themselves
3. **Idempotency**: Tests should produce the same results when run multiple times
4. **Assertions**: Use descriptive assertion messages
5. **Timeouts**: Set appropriate timeouts for async operations

## Contributing

When adding new integration tests:
1. Follow the existing test structure
2. Use the helper utilities in `helpers/` package
3. Add appropriate requirement validation comments
4. Update this README with new test coverage
5. Ensure tests pass in isolation and as a suite
