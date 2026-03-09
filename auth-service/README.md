# Auth Service

Authentication and authorization microservice for the InSavein platform.

## Features

- User registration with bcrypt password hashing (cost factor 12)
- User login with JWT token generation
- Token refresh mechanism
- Token validation
- Logout with token revocation
- Rate limiting (5 attempts per 15 minutes per email)
- In-memory token revocation list
- Health check endpoints for Kubernetes

## API Endpoints

### Public Endpoints

#### POST /api/auth/register
Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe",
  "date_of_birth": "1995-01-15"
}
```

**Response (201 Created):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

#### POST /api/auth/login
Authenticate user and receive tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

#### POST /api/auth/refresh
Refresh access token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900
}
```

#### GET /api/auth/validate
Validate a JWT token (for internal service use).

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response (200 OK):**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "roles": ["user"],
  "exp": 1705318200
}
```

### Protected Endpoints

#### POST /api/auth/logout
Logout and revoke refresh token.

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200 OK):**
```json
{
  "message": "logged out successfully"
}
```

### Health Check Endpoints

#### GET /health
Overall health check including database connectivity.

**Response (200 OK):**
```json
{
  "status": "healthy",
  "service": "auth-service"
}
```

#### GET /health/live
Liveness probe for Kubernetes.

**Response (200 OK):**
```json
{
  "status": "alive"
}
```

#### GET /health/ready
Readiness probe for Kubernetes.

**Response (200 OK):**
```json
{
  "status": "ready"
}
```

## Configuration

The service is configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | HTTP server port | 8080 |
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | Database user | insavein_user |
| DB_PASSWORD | Database password | insavein_password |
| DB_NAME | Database name | insavein_db |
| DB_SSLMODE | SSL mode for database connection | disable |
| JWT_SECRET | Secret key for JWT signing | (must be set in production) |

## Security Features

### Password Hashing
- Uses bcrypt with cost factor 12
- Passwords are never stored in plaintext
- Password hashes are never returned in API responses

### JWT Tokens
- Access tokens expire in 15 minutes
- Refresh tokens expire in 7 days
- Tokens are signed with HMAC-SHA256
- Token payload includes user_id, email, and roles

### Rate Limiting
- Maximum 5 login attempts per 15 minutes per email
- Temporary blocking for 15 minutes after exceeding limit
- Automatic cleanup of expired rate limit entries

### Token Revocation
- In-memory token revocation list
- Tokens are revoked on logout
- Expired tokens are automatically cleaned up

## Development

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 14 or higher

### Running Locally

1. Install dependencies:
```bash
go mod download
```

2. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=insavein_user
export DB_PASSWORD=insavein_password
export DB_NAME=insavein_db
export JWT_SECRET=your-secret-key
```

3. Run the service:
```bash
go run cmd/server/main.go
```

### Building Docker Image

```bash
docker build -t insavein/auth-service:latest .
```

### Running with Docker

```bash
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=insavein_user \
  -e DB_PASSWORD=insavein_password \
  -e DB_NAME=insavein_db \
  -e JWT_SECRET=your-secret-key \
  insavein/auth-service:latest
```

## Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Project Structure

```
auth-service/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── auth/
│   │   ├── service.go        # Service interface
│   │   ├── auth_service.go   # Service implementation
│   │   ├── types.go          # Request/response types
│   │   ├── repository.go     # Repository interface
│   │   ├── postgres_repository.go  # PostgreSQL implementation
│   │   ├── rate_limiter.go   # Rate limiting
│   │   └── token_store.go    # Token revocation store
│   ├── handlers/
│   │   └── auth_handler.go   # HTTP handlers
│   └── middleware/
│       └── auth_middleware.go # Authentication middleware
├── pkg/
│   └── database/
│       └── postgres.go       # Database connection
├── Dockerfile
├── go.mod
└── README.md
```

## Requirements Mapping

This implementation satisfies the following requirements from the spec:

- **Requirement 1.1**: User registration with email and password
- **Requirement 1.2**: Password hashing with bcrypt cost factor 12
- **Requirement 1.3**: Email uniqueness validation
- **Requirement 1.4**: Password length validation (minimum 8 characters)
- **Requirement 1.5**: JWT token generation with 15-minute access token and 7-day refresh token
- **Requirement 1.6**: JWT tokens signed with HMAC-SHA256
- **Requirement 1.7**: Rate limiting (5 attempts per 15 minutes)
- **Requirement 1.8**: Temporary blocking on exceeded attempts
- **Requirement 2.1**: JWT token signature validation
- **Requirement 2.2**: Token expiration validation
- **Requirement 2.3**: Token refresh functionality
- **Requirement 2.4**: Token revocation on logout
- **Requirement 2.5**: Token payload includes user_id, email, and roles
- **Requirement 2.6**: Invalid signature rejection
- **Requirement 15.5**: Token validation for authenticated requests
- **Requirement 17.1**: Input validation with error messages
- **Requirement 17.2**: Required field validation
- **Requirement 20.2**: Password hashes never returned in responses

## License

Copyright © 2024 InSavein Platform
