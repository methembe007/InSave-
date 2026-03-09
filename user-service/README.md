# User Profile Service

The User Profile Service is a microservice responsible for managing user profile information, preferences, and account operations for the InSavein platform.

## Features

- **Profile Management**: Get and update user profile information (name, date of birth, profile image)
- **Preferences Management**: Manage user preferences (currency, notifications, theme)
- **Account Deletion**: Delete user account with cascade deletion of all associated data
- **JWT Authentication**: Secure endpoints with JWT token validation
- **Authorization**: Users can only access and modify their own profile data

## API Endpoints

### Profile Operations

#### Get Profile
```
GET /api/user/profile
Authorization: Bearer <token>
```

Response:
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "date_of_birth": "1990-01-01",
  "profile_image_url": "https://example.com/image.jpg",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Update Profile
```
PUT /api/user/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "date_of_birth": "1990-01-01",
  "profile_image_url": "https://example.com/image.jpg"
}
```

### Preferences Operations

#### Get Preferences
```
GET /api/user/preferences
Authorization: Bearer <token>
```

Response:
```json
{
  "currency": "USD",
  "notifications_enabled": true,
  "email_notifications": true,
  "push_notifications": true,
  "savings_reminders": true,
  "reminder_time": "09:00",
  "theme": "light"
}
```

#### Update Preferences
```
PUT /api/user/preferences
Authorization: Bearer <token>
Content-Type: application/json

{
  "currency": "EUR",
  "notifications_enabled": true,
  "email_notifications": false,
  "push_notifications": true,
  "savings_reminders": true,
  "reminder_time": "10:00",
  "theme": "dark"
}
```

### Account Operations

#### Delete Account
```
DELETE /api/user/account
Authorization: Bearer <token>
```

Response:
```json
{
  "message": "account deleted successfully"
}
```

### Health Check

```
GET /health
GET /health/live
GET /health/ready
```

## Configuration

The service is configured using environment variables. Copy `.env.example` to `.env` and update the values:

```bash
# Server Configuration
PORT=8081

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=insavein_user
DB_PASSWORD=insavein_password
DB_NAME=insavein
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
```

## Development

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Make (optional)

### Setup

1. Install dependencies:
```bash
make deps
# or
go mod download
```

2. Run database migrations (from project root):
```bash
cd ../migrations
./migrate.sh up
```

3. Start the service:
```bash
make run
# or
go run cmd/server/main.go
```

### Testing

Run tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

### Building

Build the binary:
```bash
make build
```

Build Docker image:
```bash
make docker-build
```

## Architecture

The service follows a clean architecture pattern:

```
user-service/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── user/            # Core business logic
│   │   ├── service.go           # Service interface
│   │   ├── user_service.go      # Service implementation
│   │   ├── repository.go        # Repository interface
│   │   ├── postgres_repository.go # PostgreSQL implementation
│   │   └── types.go             # Domain types
│   ├── handlers/        # HTTP handlers
│   └── middleware/      # HTTP middleware (auth)
└── pkg/
    └── database/        # Database connection utilities
```

## Requirements Mapping

This service implements the following requirements from the design document:

- **Requirement 3.1**: Profile retrieval with email, name, date of birth, profile image
- **Requirement 3.2**: Profile update with field validation
- **Requirement 3.3**: Preferences management (currency, notifications, theme)
- **Requirement 3.4**: Account deletion with cascade
- **Requirement 3.5**: Authorization check (users can only access own profile)
- **Requirement 15.1**: JWT authentication for all endpoints
- **Requirement 15.2**: Token validation and user context
- **Requirement 15.4**: Authorization enforcement
- **Requirement 16.1**: Database transactions for atomicity
- **Requirement 16.2**: Rollback on failure

## Security

- All endpoints require JWT authentication
- Users can only access and modify their own data
- Passwords are never returned in API responses
- Database connections use connection pooling
- Graceful shutdown ensures no data loss

## License

Copyright © 2024 InSavein Platform
