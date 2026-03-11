# Goal Service API Examples

This document provides example requests and responses for the Goal Service API.

## Authentication

All endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### 1. Create Goal

**Request:**
```http
POST /api/goals
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "Emergency Fund",
  "description": "Build a 6-month emergency fund",
  "target_amount": 10000.00,
  "currency": "USD",
  "target_date": "2025-12-31T00:00:00Z",
  "milestones": [
    {
      "title": "First $2,500",
      "amount": 2500.00,
      "order": 1
    },
    {
      "title": "Halfway There",
      "amount": 5000.00,
      "order": 2
    },
    {
      "title": "Almost Done",
      "amount": 7500.00,
      "order": 3
    }
  ]
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Emergency Fund",
  "description": "Build a 6-month emergency fund",
  "target_amount": 10000.00,
  "current_amount": 0.00,
  "currency": "USD",
  "target_date": "2025-12-31T00:00:00Z",
  "status": "active",
  "progress_percent": 0.00,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

---

### 2. Get Active Goals

**Request:**
```http
GET /api/goals
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Emergency Fund",
    "description": "Build a 6-month emergency fund",
    "target_amount": 10000.00,
    "current_amount": 3500.00,
    "currency": "USD",
    "target_date": "2025-12-31T00:00:00Z",
    "status": "active",
    "progress_percent": 35.00,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-02-20T14:15:00Z"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Vacation Fund",
    "description": "Save for summer vacation",
    "target_amount": 3000.00,
    "current_amount": 1200.00,
    "currency": "USD",
    "target_date": "2024-06-01T00:00:00Z",
    "status": "active",
    "progress_percent": 40.00,
    "created_at": "2024-01-20T09:00:00Z",
    "updated_at": "2024-02-18T16:30:00Z"
  }
]
```

---

### 3. Get Specific Goal

**Request:**
```http
GET /api/goals/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Emergency Fund",
  "description": "Build a 6-month emergency fund",
  "target_amount": 10000.00,
  "current_amount": 3500.00,
  "currency": "USD",
  "target_date": "2025-12-31T00:00:00Z",
  "status": "active",
  "progress_percent": 35.00,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-02-20T14:15:00Z",
  "milestones": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440010",
      "goal_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "First $2,500",
      "amount": 2500.00,
      "is_completed": true,
      "completed_at": "2024-02-10T12:00:00Z",
      "order": 1
    },
    {
      "id": "880e8400-e29b-41d4-a716-446655440011",
      "goal_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Halfway There",
      "amount": 5000.00,
      "is_completed": false,
      "completed_at": null,
      "order": 2
    },
    {
      "id": "990e8400-e29b-41d4-a716-446655440012",
      "goal_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Almost Done",
      "amount": 7500.00,
      "is_completed": false,
      "completed_at": null,
      "order": 3
    }
  ]
}
```

---

### 4. Update Goal

**Request:**
```http
PUT /api/goals/550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "Emergency Fund - Updated",
  "target_amount": 12000.00,
  "target_date": "2026-01-31T00:00:00Z"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Emergency Fund - Updated",
  "description": "Build a 6-month emergency fund",
  "target_amount": 12000.00,
  "current_amount": 3500.00,
  "currency": "USD",
  "target_date": "2026-01-31T00:00:00Z",
  "status": "active",
  "progress_percent": 29.17,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-02-21T10:00:00Z"
}
```

---

### 5. Update Progress (Add Contribution)

**Request:**
```http
POST /api/goals/550e8400-e29b-41d4-a716-446655440000/progress
Content-Type: application/json
Authorization: Bearer <token>

{
  "amount": 500.00
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Emergency Fund",
  "description": "Build a 6-month emergency fund",
  "target_amount": 10000.00,
  "current_amount": 4000.00,
  "currency": "USD",
  "target_date": "2025-12-31T00:00:00Z",
  "status": "active",
  "progress_percent": 40.00,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-02-21T11:30:00Z"
}
```

**Note:** When this contribution brings the current_amount to 4000, no new milestones are completed (next milestone is at 5000).

---

### 6. Update Progress - Goal Completion

**Request:**
```http
POST /api/goals/550e8400-e29b-41d4-a716-446655440000/progress
Content-Type: application/json
Authorization: Bearer <token>

{
  "amount": 6000.00
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Emergency Fund",
  "description": "Build a 6-month emergency fund",
  "target_amount": 10000.00,
  "current_amount": 10000.00,
  "currency": "USD",
  "target_date": "2025-12-31T00:00:00Z",
  "status": "completed",
  "progress_percent": 100.00,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-03-01T15:45:00Z"
}
```

**Note:** Status automatically changed to "completed" when current_amount >= target_amount. All remaining milestones are also marked as completed.

---

### 7. Get Milestones

**Request:**
```http
GET /api/goals/550e8400-e29b-41d4-a716-446655440000/milestones
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440010",
    "goal_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "First $2,500",
    "amount": 2500.00,
    "is_completed": true,
    "completed_at": "2024-02-10T12:00:00Z",
    "order": 1
  },
  {
    "id": "880e8400-e29b-41d4-a716-446655440011",
    "goal_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Halfway There",
    "amount": 5000.00,
    "is_completed": true,
    "completed_at": "2024-03-01T15:45:00Z",
    "order": 2
  },
  {
    "id": "990e8400-e29b-41d4-a716-446655440012",
    "goal_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Almost Done",
    "amount": 7500.00,
    "is_completed": true,
    "completed_at": "2024-03-01T15:45:00Z",
    "order": 3
  }
]
```

---

### 8. Delete Goal

**Request:**
```http
DELETE /api/goals/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Goal deleted successfully"
}
```

**Note:** All associated milestones are automatically deleted due to database cascade.

---

## Error Responses

### 400 Bad Request - Validation Error

```json
{
  "error": "Key: 'CreateGoalRequest.TargetAmount' Error:Field validation for 'TargetAmount' failed on the 'gt' tag"
}
```

### 401 Unauthorized - Missing/Invalid Token

```json
{
  "error": "Invalid or expired token"
}
```

### 403 Forbidden - Not Authorized

```json
{
  "error": "goal does not belong to user"
}
```

### 404 Not Found

```json
{
  "error": "goal not found"
}
```

### 500 Internal Server Error

```json
{
  "error": "failed to create goal: database connection error"
}
```

---

## Testing with cURL

### Create a Goal
```bash
curl -X POST http://localhost:8005/api/goals \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Emergency Fund",
    "description": "Build a 6-month emergency fund",
    "target_amount": 10000.00,
    "currency": "USD",
    "target_date": "2025-12-31T00:00:00Z"
  }'
```

### Get Active Goals
```bash
curl -X GET http://localhost:8005/api/goals \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update Progress
```bash
curl -X POST http://localhost:8005/api/goals/GOAL_ID/progress \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 500.00
  }'
```

---

## Notes

1. **Currency**: Defaults to "USD" if not provided
2. **Progress Percentage**: Automatically calculated as (current_amount / target_amount) × 100
3. **Goal Completion**: Status automatically changes to "completed" when current_amount >= target_amount
4. **Milestone Completion**: Milestones are automatically marked as completed when the goal's current_amount reaches their threshold
5. **Concurrency**: Multiple simultaneous progress updates are handled safely with database row-level locking
6. **Authorization**: Users can only access and modify their own goals
