# Education Service

The Education Service is a microservice for delivering financial education content and tracking user progress in the InSavein platform.

## Features

- Retrieve financial education lessons with completion status
- Get detailed lesson content including text, videos, and resources
- Track lesson completion with timestamps
- Calculate education progress percentage
- Read from database replicas to reduce load on primary database

## API Endpoints

### Lessons

- `GET /api/education/lessons` - Get all lessons with completion status for authenticated user
- `GET /api/education/lessons/:id` - Get detailed lesson content
- `POST /api/education/lessons/:id/complete` - Mark a lesson as completed

### Progress

- `GET /api/education/progress` - Get user's education progress (total, completed, percentage)

## Requirements

- Go 1.21+
- PostgreSQL 14+ with read replicas
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
- `internal/education/` - Business logic and domain models
- `internal/handlers/` - HTTP request handlers
- `internal/middleware/` - HTTP middleware (authentication)
- `pkg/database/` - Database connection utilities

## Database Schema

The service uses the following tables:

- `lessons` - Financial education lesson content
- `education_progress` - User lesson completion tracking

## Read Replica Usage

The Education Service implements read/write splitting to optimize database performance:

- **Read Operations** (from replicas):
  - Get all lessons
  - Get lesson details
  - Get user completions
  - Count lessons
  - Count completed lessons

- **Write Operations** (to primary):
  - Mark lesson complete

This follows Requirement 11.6: "THE Education_Service SHALL read lesson content from database replicas to reduce load on primary database"

## Authentication

All API endpoints (except `/health`) require JWT authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Implementation Details

### Lesson Retrieval

Lessons are retrieved from read replicas with completion status merged from the user's progress records. This reduces load on the primary database while ensuring users see their current progress.

### Progress Calculation

Progress percentage is calculated as:
```
progress_percent = (completed_lessons / total_lessons) × 100
```

This follows Requirement 11.5.

### Completion Tracking

When a user marks a lesson as complete:
1. The lesson existence is verified
2. A record is inserted or updated in `education_progress` with `is_completed = true` and current timestamp
3. The operation uses `ON CONFLICT` to handle duplicate completions gracefully

## Requirements Implemented

- **11.1**: Return lessons with title, description, category, duration, difficulty, and completion status
- **11.2**: Return detailed lesson content including text, video URL, and resources
- **11.3**: Record lesson completion with timestamp
- **11.4**: Return total lessons, completed lessons, and progress percentage
- **11.5**: Calculate progress percentage as (completed_lessons / total_lessons) × 100
- **11.6**: Read lesson content from database replicas to reduce load on primary database
- **15.1**: JWT authentication on all API endpoints

## Port

The service runs on port **8085** by default.
