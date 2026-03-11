# Goal Service

The Goal Service is a microservice for managing financial goals and milestones in the InSavein platform.

## Features

- Create and manage financial goals with target amounts and dates
- Track goal progress with contributions
- Milestone tracking with automatic completion
- Concurrency-safe progress updates with row-level locking
- Automatic goal completion when target is reached

## API Endpoints

### Goals

- `POST /api/goals` - Create a new goal
- `GET /api/goals` - Get all active goals for the authenticated user
- `GET /api/goals/:id` - Get a specific goal with milestones
- `PUT /api/goals/:id` - Update a goal
- `DELETE /api/goals/:id` - Delete a goal (cascades to milestones)

### Progress

- `POST /api/goals/:id/progress` - Add a contribution to a goal

### Milestones

- `GET /api/goals/:id/milestones` - Get all milestones for a goal

## Requirements

- Go 1.21+
- PostgreSQL 14+
- Environment variables configured (see .env.example)

## Setup

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

## Development

### Build
```bash
make build
```

### Run Tests
```bash
make test
```

### Clean Build Artifacts
```bash
make clean
```

## Docker

### Build Docker Image
```bash
make docker-build
```

### Run Docker Container
```bash
make docker-run
```

## Architecture

The service follows a clean architecture pattern:

- `cmd/server/` - Application entry point
- `internal/goal/` - Business logic and domain models
- `internal/handlers/` - HTTP request handlers
- `internal/middleware/` - HTTP middleware (authentication)
- `pkg/database/` - Database connection utilities

## Database Schema

The service uses the following tables:

- `goals` - Financial goal records
- `goal_milestones` - Milestone checkpoints for goals

## Authentication

All API endpoints (except `/health`) require JWT authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Implementation Details

### Goal Progress Updates

Goal progress updates use database transactions with row-level locking (`FOR UPDATE`) to prevent race conditions when multiple contributions are made concurrently.

### Milestone Completion

When progress is updated, the service automatically:
1. Checks all uncompleted milestones
2. Marks milestones as completed if the current amount reaches their threshold
3. Processes milestones in ascending order by amount
4. Stops at the first unreached milestone

### Goal Completion

Goals are automatically marked as "completed" when the current amount reaches or exceeds the target amount.

## Requirements Implemented

- **9.1**: Goal creation with title, description, target amount, target date, and currency
- **9.2**: Initialize current_amount to 0 and status to "active"
- **9.3**: Retrieve active goals filtered by status
- **9.4**: Update goal fields
- **9.5**: Delete goal with cascade to milestones
- **9.6**: Calculate progress percentage as (current_amount / target_amount) × 100
- **10.1**: Add contributions to increase current_amount
- **10.2**: Change status to "completed" when target reached
- **10.3**: Use row-level locking for concurrency control
- **10.4**: Update milestones when reached
- **10.5**: Set completed_at timestamp on milestone completion
- **10.6**: Process milestones in ascending order, stop at first unreached
- **15.1**: JWT authentication on all API endpoints
- **15.4**: Authorization check (users can only access their own goals)
- **16.4**: Database transactions for atomic updates
