# Notification Service

The Notification Service handles email notifications, push notifications, and reminders for the InSavein platform.

## Features

- **Email Notifications**: Send template-based emails via SendGrid or AWS SES
- **Push Notifications**: Deliver push notifications to mobile and web via Firebase Cloud Messaging
- **Notification Preferences**: Respect user notification preferences before sending
- **Notification History**: Track and retrieve user notification history
- **Read Status**: Mark notifications as read

## Requirements Implemented

- **12.1**: Email notification delivery with template support
- **12.2**: Push notification delivery via Firebase Cloud Messaging
- **12.3**: Reminder scheduling
- **12.4**: Notification history retrieval (date descending order)
- **12.5**: Mark notifications as read
- **12.6**: Notification preference enforcement
- **15.1**: API request authentication

## API Endpoints

### GET /api/notifications
Retrieve all notifications for the authenticated user.

**Authentication**: Required (JWT token)

**Response**:
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "type": "push",
    "title": "Savings Reminder",
    "message": "Don't forget to save today!",
    "is_read": false,
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

### PUT /api/notifications/:id/read
Mark a notification as read.

**Authentication**: Required (JWT token)

**Response**:
```json
{
  "message": "Notification marked as read"
}
```

## Setup

1. Copy `.env.example` to `.env` and configure:
   ```bash
   cp .env.example .env
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the service:
   ```bash
   go run cmd/server/main.go
   ```

## Configuration

### Email Providers

The service supports multiple email providers:

- **mock**: Simulates email sending (for development)
- **sendgrid**: SendGrid API integration (requires API key)
- **aws_ses**: AWS Simple Email Service (requires AWS credentials)

Set `EMAIL_PROVIDER` in `.env` to choose the provider.

### Push Notification Providers

The service supports Firebase Cloud Messaging:

- **mock**: Simulates push notifications (for development)
- **fcm**: Firebase Cloud Messaging (requires server key and project ID)

Set `PUSH_PROVIDER` in `.env` to choose the provider.

## Architecture

```
notification-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/
│   │   └── notification_handler.go  # HTTP request handlers
│   ├── middleware/
│   │   └── auth_middleware.go   # JWT authentication
│   └── notification/
│       ├── service.go           # Service interface
│       ├── notification_service.go  # Service implementation
│       ├── postgres_repository.go   # Database operations
│       ├── email_provider.go    # Email delivery
│       ├── push_provider.go     # Push notification delivery
│       └── types.go             # Data structures
└── pkg/
    └── database/
        └── postgres.go          # Database connection
```

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o notification-service cmd/server/main.go
```

## Integration Notes

### Email Templates

When using SendGrid or AWS SES, create email templates in your provider's dashboard and reference them by `template_id` in the `EmailRequest`.

### Firebase Cloud Messaging

To enable FCM:
1. Create a Firebase project
2. Generate a server key from Firebase Console
3. Configure `FCM_SERVER_KEY` and `FCM_PROJECT_ID` in `.env`
4. Store user device tokens in the database

### User Preferences

The service checks user preferences before sending notifications. Preferences are stored in the `users.preferences` JSONB field:

```json
{
  "notifications_enabled": true,
  "email_notifications": true,
  "push_notifications": true
}
```

## Dependencies

- `github.com/gorilla/mux`: HTTP router
- `github.com/lib/pq`: PostgreSQL driver
- `github.com/golang-jwt/jwt/v5`: JWT token handling
- `github.com/rs/cors`: CORS middleware
- `github.com/joho/godotenv`: Environment variable loading
