# Budget Service

The Budget Service is a microservice for the InSavein platform that handles monthly budget management, spending tracking, and budget alerts.

## Features

- **Budget Management**: Create and update monthly budgets with category allocations
- **Spending Tracking**: Record spending transactions with atomic updates
- **Budget Alerts**: Automatic detection of budget threshold breaches (80% warning, 100% critical)
- **Category Management**: Organize spending into customizable categories
- **Spending Summary**: View spending summaries by month

## API Endpoints

### Budget Operations

- `POST /api/budget` - Create a new monthly budget
- `GET /api/budget/current` - Get the current month's budget
- `PUT /api/budget/:id` - Update an existing budget
- `GET /api/budget/categories` - Get all budget categories
- `GET /api/budget/summary?month=YYYY-MM` - Get spending summary for a month

### Spending Operations

- `POST /api/budget/spending` - Record a spending transaction
- `GET /api/budget/alerts` - Get budget alerts for current month

### Health Checks

- `GET /health` - Overall health check
- `GET /health/live` - Liveness probe
- `GET /health/ready` - Readiness probe

## Requirements

- Go 1.25.4 or higher
- PostgreSQL 14 or higher
- JWT secret for authentication

## Environment Variables

See `.env.example` for required environment variables:

- `PORT` - Server port (default: 8083)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSLMODE` - SSL mode for database connection
- `JWT_SECRET` - Secret key for JWT validation

## Getting Started

1. Install dependencies:
   ```bash
   make deps
   ```

2. Copy environment file:
   ```bash
   cp .env.example .env
   ```

3. Update `.env` with your configuration

4. Run the service:
   ```bash
   make run
   ```

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Run with Coverage

```bash
make test-coverage
```

### Format Code

```bash
make fmt
```

## API Examples

### Create Budget

```bash
curl -X POST http://localhost:8083/api/budget \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "month": "2024-01-01T00:00:00Z",
    "categories": [
      {
        "name": "Groceries",
        "allocated_amount": 500.00,
        "color": "#4CAF50"
      },
      {
        "name": "Transportation",
        "allocated_amount": 200.00,
        "color": "#2196F3"
      }
    ]
  }'
```

### Record Spending

```bash
curl -X POST http://localhost:8083/api/budget/spending \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": "CATEGORY_UUID",
    "amount": 45.50,
    "description": "Weekly groceries",
    "merchant": "Local Supermarket",
    "date": "2024-01-15T10:30:00Z"
  }'
```

### Get Budget Alerts

```bash
curl -X GET http://localhost:8083/api/budget/alerts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Architecture

The service follows a clean architecture pattern:

- `cmd/server/` - Application entry point
- `internal/budget/` - Core business logic
  - `service.go` - Service interface definitions
  - `budget_service.go` - Service implementation
  - `postgres_repository.go` - Database operations
  - `types.go` - Data models
- `internal/handlers/` - HTTP request handlers
- `internal/middleware/` - Authentication middleware
- `pkg/database/` - Database connection utilities

## Key Features

### Atomic Spending Updates

When recording spending, the service uses database transactions to ensure:
1. Spending transaction is created
2. Category spent amount is updated
3. Budget total spent is updated

All three operations succeed or fail together, maintaining data consistency.

### Budget Alert Detection

The service automatically detects:
- **Warning alerts**: Categories at 80-99% of allocated amount
- **Critical alerts**: Categories at 100%+ of allocated amount

Alerts are sorted by severity (critical first) and percentage (highest first).

### Unique Budget Constraint

Each user can have only one budget per month, enforced at the database level with a unique constraint on (user_id, month).

## License

Copyright © 2024 InSavein Platform
