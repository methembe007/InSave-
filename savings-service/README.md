# Savings Service

The Savings Service is a microservice for the InSavein platform that handles savings transaction tracking, streak calculation, and savings history management.

## Features

- **Savings Transaction Management**: Create and track savings deposits
- **Streak Calculation**: Automatically calculate and maintain savings streaks
- **Savings History**: Retrieve transaction history with pagination
- **Savings Summary**: Get comprehensive savings statistics
- **JWT Authentication**: Secure endpoints with JWT token validation

## API Endpoints

### POST /api/savings/transactions
Create a new savings transaction.

**Request Body:**
```json
{
  "amount": 50.00,
  "currency": "USD",
  "description": "Weekly savings",
  "category": "general"
}
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "amount": 50.00,
  "currency": "USD",
  "description": "Weekly savings",
  "category": "general",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### GET /api/savings/history
Get savings transaction history with pagination.

**Query Parameters:**
- `limit` (optional): Number of transactions to return (default: 50)
- `offset` (optional): Number of transactions to skip (default: 0)

**Response:**
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "amount": 50.00,
    "currency": "USD",
    "description": "Weekly savings",
    "category": "general",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### GET /api/savings/summary
Get comprehensive savings summary.

**Response:**
```json
{
  "total_saved": 1000.00,
  "current_streak": 7,
  "longest_streak": 14,
  "last_saving_date": "2024-01-01T00:00:00Z",
  "monthly_average": 200.00,
  "this_month_saved": 250.00
}
```

### GET /api/savings/streak
Get current savings streak information.

**Response:**
```json
{
  "current_streak": 7,
  "longest_streak": 14,
  "last_save_date": "2024-01-01T00:00:00Z"
}
```

## Health Check Endpoints

- `GET /health` - Overall health status
- `GET /health/live` - Liveness probe
- `GET /health/ready` - Readiness probe

## Configuration

The service is configured using environment variables. See `.env.example` for available options.

### Environment Variables

- `PORT`: Server port (default: 8082)
- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSLMODE`: SSL mode (default: disable)
- `JWT_SECRET`: Secret key for JWT validation

## Running the Service

### Local Development

1. Copy `.env.example` to `.env` and configure:
```bash
cp .env.example .env
```

2. Install dependencies:
```bash
make deps
```

3. Run the service:
```bash
make run
```

### Building

```bash
make build
```

### Testing

```bash
make test
```

## Architecture

The service follows a clean architecture pattern:

- `cmd/server`: Application entry point
- `internal/savings`: Core business logic
  - `service.go`: Service interface
  - `savings_service.go`: Service implementation
  - `repository.go`: Repository interface
  - `postgres_repository.go`: PostgreSQL implementation
  - `types.go`: Data structures
- `internal/handlers`: HTTP handlers
- `internal/middleware`: Authentication middleware
- `pkg/database`: Database connection utilities

## Streak Calculation Algorithm

The service implements an intelligent streak calculation algorithm:

1. Retrieves all unique saving dates for the user
2. Checks if the last save was today or yesterday
3. If more than 1 day has passed, sets current streak to 0
4. Otherwise, counts consecutive days backward from the last save
5. Handles multiple transactions on the same day (counts as one day)
6. Updates longest streak if current streak exceeds it
7. Ensures longest streak never decreases

## Dependencies

- Go 1.23+
- PostgreSQL 14+
- github.com/lib/pq: PostgreSQL driver
- github.com/google/uuid: UUID generation
- github.com/golang-jwt/jwt/v5: JWT token handling
- github.com/gorilla/mux: HTTP routing

## License

Copyright © 2024 InSavein Platform
