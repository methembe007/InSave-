# InSavein Platform API Documentation

**Version**: 1.0.0  
**Base URL**: `https://api.insavein.com`  
**Authentication**: JWT Bearer Token

## Table of Contents

1. [Authentication](#authentication)
2. [Auth Service API](#auth-service-api)
3. [User Service API](#user-service-api)
4. [Savings Service API](#savings-service-api)
5. [Budget Service API](#budget-service-api)
6. [Goal Service API](#goal-service-api)
7. [Education Service API](#education-service-api)
8. [Notification Service API](#notification-service-api)
9. [Analytics Service API](#analytics-service-api)
10. [Error Codes](#error-codes)
11. [Rate Limiting](#rate-limiting)

---

## Authentication

All API endpoints (except registration and login) require authentication using JWT Bearer tokens.

### Headers

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

### Token Expiration

- **Access Token**: 15 minutes
- **Refresh Token**: 7 days

### Refreshing Tokens

When the access token expires, use the refresh token to obtain new tokens:

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## Auth Service API

**Base URL**: `http://localhost:8080` (development)

### POST /api/auth/register

Register a new user account.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "Jane",
  "last_name": "Doe",
  "date_of_birth": "1998-05-15"
}
```

**Validation Rules**:
- `email`: Valid email format, unique
- `password`: Minimum 8 characters
- `first_name`, `last_name`: Required, non-empty
- `date_of_birth`: Valid date in YYYY-MM-DD format

**Success Response** (201 Created):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "first_name": "Jane",
    "last_name": "Doe"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
- `409 Conflict`: Email already registered
- `500 Internal Server Error`: Server error

### POST /api/auth/login

Authenticate user and receive JWT tokens.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Success Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "first_name": "Jane",
    "last_name": "Doe"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Missing required fields
- `401 Unauthorized`: Invalid credentials
- `429 Too Many Requests`: Rate limit exceeded (5 attempts in 15 minutes)

### POST /api/auth/refresh

Refresh access token using refresh token.

**Request Body**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid or expired refresh token

### POST /api/auth/logout

Invalidate refresh token and log out user.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "message": "Logged out successfully"
}
```

---

## User Service API

**Base URL**: `http://localhost:8081` (development)

### GET /api/users/profile

Get current user's profile.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "first_name": "Jane",
  "last_name": "Doe",
  "date_of_birth": "1998-05-15",
  "profile_image_url": "https://cdn.insavein.com/profiles/user123.jpg",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-15T10:30:00Z"
}
```

### PUT /api/users/profile

Update user profile.

**Headers**: Requires authentication

**Request Body**:
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "profile_image_url": "https://cdn.insavein.com/profiles/newimage.jpg"
}
```

**Success Response** (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "first_name": "Jane",
  "last_name": "Smith",
  "date_of_birth": "1998-05-15",
  "profile_image_url": "https://cdn.insavein.com/profiles/newimage.jpg",
  "updated_at": "2026-01-15T11:00:00Z"
}
```

### GET /api/users/preferences

Get user preferences.

**Headers**: Requires authentication

**Success Response** (200 OK):
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

### PUT /api/users/preferences

Update user preferences.

**Headers**: Requires authentication

**Request Body**:
```json
{
  "currency": "USD",
  "notifications_enabled": true,
  "savings_reminders": true,
  "reminder_time": "09:00",
  "theme": "dark"
}
```

**Success Response** (200 OK):
```json
{
  "message": "Preferences updated successfully"
}
```

---

## Savings Service API

**Base URL**: `http://localhost:8082` (development)

### GET /api/savings/summary

Get savings summary for current user.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "total_saved": 1250.50,
  "current_streak": 15,
  "longest_streak": 30,
  "last_saving_date": "2026-01-15T00:00:00Z",
  "monthly_average": 125.00,
  "this_month_saved": 250.00
}
```

### GET /api/savings/history

Get savings transaction history.

**Headers**: Requires authentication

**Query Parameters**:
- `limit` (optional): Number of transactions (default: 50, max: 100)
- `offset` (optional): Pagination offset (default: 0)
- `start_date` (optional): Filter by start date (YYYY-MM-DD)
- `end_date` (optional): Filter by end date (YYYY-MM-DD)

**Example Request**:
```http
GET /api/savings/history?limit=10&offset=0&start_date=2026-01-01
```

**Success Response** (200 OK):
```json
{
  "transactions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "amount": 25.00,
      "currency": "USD",
      "description": "Daily savings",
      "category": "emergency",
      "created_at": "2026-01-15T10:00:00Z"
    }
  ],
  "total": 150,
  "limit": 10,
  "offset": 0
}
```

### POST /api/savings/transactions

Create a new savings transaction.

**Headers**: Requires authentication

**Request Body**:
```json
{
  "amount": 25.00,
  "currency": "USD",
  "description": "Daily savings",
  "category": "emergency"
}
```

**Validation Rules**:
- `amount`: Must be positive (> 0)
- `currency`: Valid 3-letter ISO code
- `category`: One of: emergency, goal, general

**Success Response** (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 25.00,
  "currency": "USD",
  "description": "Daily savings",
  "category": "emergency",
  "created_at": "2026-01-15T10:00:00Z"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid amount or missing required fields
- `401 Unauthorized`: Invalid or missing token

### GET /api/savings/streak

Get current savings streak.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "current_streak": 15,
  "longest_streak": 30,
  "last_save_date": "2026-01-15T00:00:00Z"
}
```

---

## Budget Service API

**Base URL**: `http://localhost:8083` (development)

### GET /api/budgets/current

Get current month's budget.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "month": "2026-01-01T00:00:00Z",
  "total_budget": 500.00,
  "total_spent": 320.50,
  "remaining_budget": 179.50,
  "categories": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "budget_id": "550e8400-e29b-41d4-a716-446655440002",
      "name": "Food",
      "allocated_amount": 200.00,
      "spent_amount": 150.00,
      "remaining_amount": 50.00,
      "color": "#FF6B6B"
    }
  ],
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-15T10:00:00Z"
}
```

### POST /api/budgets

Create a new budget.

**Headers**: Requires authentication

**Request Body**:
```json
{
  "month": "2026-02",
  "total_budget": 500.00,
  "categories": [
    {
      "name": "Food",
      "allocated_amount": 200.00,
      "color": "#FF6B6B"
    },
    {
      "name": "Transport",
      "allocated_amount": 100.00,
      "color": "#4ECDC4"
    }
  ]
}
```

**Validation Rules**:
- `total_budget`: Must be non-negative
- `categories`: At least one category required
- `allocated_amount`: Must be non-negative
- Sum of category allocations should not exceed total_budget

**Success Response** (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440004",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "month": "2026-02-01T00:00:00Z",
  "total_budget": 500.00,
  "total_spent": 0.00,
  "remaining_budget": 500.00,
  "categories": [...]
}
```

**Error Responses**:
- `400 Bad Request`: Invalid data or validation error
- `409 Conflict`: Budget already exists for this month

### POST /api/budgets/spending

Record a spending transaction.

**Headers**: Requires authentication

**Request Body**:
```json
{
  "category_id": "550e8400-e29b-41d4-a716-446655440003",
  "amount": 25.50,
  "description": "Grocery shopping",
  "merchant": "SuperMart",
  "date": "2026-01-15"
}
```

**Validation Rules**:
- `amount`: Must be positive
- `date`: Cannot be in the future
- `category_id`: Must exist and belong to user's budget

**Success Response** (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440005",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "budget_id": "550e8400-e29b-41d4-a716-446655440002",
  "category_id": "550e8400-e29b-41d4-a716-446655440003",
  "amount": 25.50,
  "description": "Grocery shopping",
  "merchant": "SuperMart",
  "date": "2026-01-15T00:00:00Z",
  "created_at": "2026-01-15T10:30:00Z"
}
```

### GET /api/budgets/alerts

Get budget alerts for current month.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "alerts": [
    {
      "category_name": "Food",
      "percentage_used": 95.5,
      "alert_type": "critical",
      "message": "You've exceeded your Food budget by 5.5%"
    },
    {
      "category_name": "Transport",
      "percentage_used": 85.0,
      "alert_type": "warning",
      "message": "You've used 85.0% of your Transport budget"
    }
  ]
}
```

---

## Goal Service API

**Base URL**: `http://localhost:8005` (development)

### GET /api/goals

Get all active goals.

**Headers**: Requires authentication

**Query Parameters**:
- `status` (optional): Filter by status (active, completed, paused)

**Success Response** (200 OK):
```json
{
  "goals": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440006",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Emergency Fund",
      "description": "Build 3 months of expenses",
      "target_amount": 1500.00,
      "current_amount": 750.00,
      "currency": "USD",
      "target_date": "2026-12-31T00:00:00Z",
      "status": "active",
      "progress_percent": 50.0,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-15T10:00:00Z"
    }
  ]
}
```

### POST /api/goals

Create a new financial goal.

**Headers**: Requires authentication

**Request Body**:
```json
{
  "title": "Emergency Fund",
  "description": "Build 3 months of expenses",
  "target_amount": 1500.00,
  "currency": "USD",
  "target_date": "2026-12-31"
}
```

**Validation Rules**:
- `target_amount`: Must be positive
- `target_date`: Must be in the future

**Success Response** (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440006",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Emergency Fund",
  "description": "Build 3 months of expenses",
  "target_amount": 1500.00,
  "current_amount": 0.00,
  "currency": "USD",
  "target_date": "2026-12-31T00:00:00Z",
  "status": "active",
  "progress_percent": 0.0,
  "created_at": "2026-01-15T10:00:00Z"
}
```

### POST /api/goals/{goal_id}/contribute

Add contribution to a goal.

**Headers**: Requires authentication

**Request Body**:
```json
{
  "amount": 50.00
}
```

**Success Response** (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440006",
  "current_amount": 800.00,
  "progress_percent": 53.33,
  "status": "active"
}
```

### GET /api/goals/{goal_id}/milestones

Get milestones for a goal.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "milestones": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440007",
      "goal_id": "550e8400-e29b-41d4-a716-446655440006",
      "title": "First $500",
      "amount": 500.00,
      "is_completed": true,
      "completed_at": "2026-01-10T00:00:00Z",
      "order": 1
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440008",
      "goal_id": "550e8400-e29b-41d4-a716-446655440006",
      "title": "Halfway There",
      "amount": 750.00,
      "is_completed": true,
      "completed_at": "2026-01-15T00:00:00Z",
      "order": 2
    }
  ]
}
```

---

## Education Service API

**Base URL**: `http://localhost:8085` (development)

### GET /api/education/lessons

Get all available lessons.

**Headers**: Requires authentication

**Query Parameters**:
- `category` (optional): Filter by category
- `difficulty` (optional): Filter by difficulty (beginner, intermediate, advanced)

**Success Response** (200 OK):
```json
{
  "lessons": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440009",
      "title": "Introduction to Budgeting",
      "description": "Learn the basics of creating and maintaining a budget",
      "category": "budgeting",
      "duration_minutes": 15,
      "difficulty": "beginner",
      "tags": ["budgeting", "basics", "finance"],
      "is_completed": false,
      "order": 1
    }
  ]
}
```

### GET /api/education/lessons/{lesson_id}

Get detailed lesson content.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440009",
  "title": "Introduction to Budgeting",
  "description": "Learn the basics of creating and maintaining a budget",
  "category": "budgeting",
  "duration_minutes": 15,
  "difficulty": "beginner",
  "content": "# Introduction to Budgeting\n\nBudgeting is...",
  "video_url": "https://cdn.insavein.com/videos/lesson1.mp4",
  "resources": [
    {
      "title": "Budget Template",
      "url": "https://cdn.insavein.com/resources/budget-template.pdf",
      "type": "pdf"
    }
  ],
  "quiz": [
    {
      "question": "What is the 50/30/20 rule?",
      "options": ["A", "B", "C", "D"],
      "correct_answer": 0
    }
  ]
}
```

### POST /api/education/lessons/{lesson_id}/complete

Mark a lesson as completed.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "message": "Lesson marked as complete",
  "progress": {
    "total_lessons": 50,
    "completed_lessons": 15,
    "progress_percent": 30.0
  }
}
```

---

## Notification Service API

**Base URL**: `http://localhost:8086` (development)

### GET /api/notifications

Get user notifications.

**Headers**: Requires authentication

**Query Parameters**:
- `limit` (optional): Number of notifications (default: 20, max: 100)
- `unread_only` (optional): Filter unread notifications (true/false)

**Success Response** (200 OK):
```json
{
  "notifications": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "type": "budget_alert",
      "title": "Budget Alert",
      "message": "You've used 85% of your Food budget",
      "is_read": false,
      "created_at": "2026-01-15T10:00:00Z"
    }
  ]
}
```

### PUT /api/notifications/{notification_id}/read

Mark notification as read.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "message": "Notification marked as read"
}
```

---

## Analytics Service API

**Base URL**: `http://localhost:8008` (development)

### GET /api/analytics/spending

Get spending analysis.

**Headers**: Requires authentication

**Query Parameters**:
- `period`: Time period (week, month, quarter, year)
- `start_date` (optional): Custom start date
- `end_date` (optional): Custom end date

**Success Response** (200 OK):
```json
{
  "period": {
    "type": "month",
    "start": "2026-01-01T00:00:00Z",
    "end": "2026-01-31T23:59:59Z"
  },
  "total_spending": 450.00,
  "category_breakdown": [
    {
      "category": "Food",
      "amount": 200.00,
      "percentage": 44.4
    }
  ],
  "top_merchants": [
    {
      "merchant": "SuperMart",
      "amount": 150.00,
      "transaction_count": 8
    }
  ],
  "daily_average": 15.00,
  "comparison_to_previous": -10.5,
  "trends": [
    {
      "category": "Food",
      "trend": "decreasing",
      "change_percent": -5.0
    }
  ]
}
```

### GET /api/analytics/financial-health

Get financial health score.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "overall_score": 75,
  "savings_score": 80,
  "budget_score": 70,
  "consistency_score": 75,
  "insights": [
    "Your savings rate is above average",
    "You're staying within budget most months"
  ],
  "improvement_areas": [
    "Try to increase your emergency fund",
    "Consider reducing discretionary spending"
  ]
}
```

### GET /api/analytics/recommendations

Get personalized recommendations.

**Headers**: Requires authentication

**Success Response** (200 OK):
```json
{
  "recommendations": [
    {
      "id": "rec_001",
      "type": "savings",
      "priority": "high",
      "title": "Increase Emergency Fund",
      "description": "Your emergency fund covers only 1 month of expenses",
      "action_items": [
        "Save an additional $50 per week",
        "Redirect 10% of income to emergency fund"
      ],
      "potential_savings": 200.00
    }
  ]
}
```

---

## Error Codes

All error responses follow this format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "Additional context"
    }
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_REQUEST` | 400 | Request validation failed |
| `UNAUTHORIZED` | 401 | Authentication required or failed |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource already exists |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

### Validation Errors

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Validation failed",
    "details": {
      "amount": "must be greater than 0",
      "email": "invalid email format"
    }
  }
}
```

---

## Rate Limiting

**Limits**:
- 100 requests per minute per user
- 1000 requests per minute per IP address
- Burst allowance: 20 requests

**Rate Limit Headers**:
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640000000
```

**Rate Limit Exceeded Response** (429):
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests",
    "retry_after": 60
  }
}
```

---

## Pagination

List endpoints support pagination:

**Query Parameters**:
- `limit`: Number of items (default: 50, max: 100)
- `offset`: Number of items to skip (default: 0)

**Response Format**:
```json
{
  "data": [...],
  "pagination": {
    "total": 150,
    "limit": 50,
    "offset": 0,
    "has_more": true
  }
}
```

---

## Webhooks (Future Feature)

Coming soon: Webhook support for real-time event notifications.

---

**Last Updated**: 2026-01-15  
**API Version**: 1.0.0
