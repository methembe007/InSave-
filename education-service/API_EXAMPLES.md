# Education Service API Examples

This document provides example API requests and responses for the Education Service.

## Authentication

All API endpoints (except `/health`) require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### 1. Health Check

Check if the service is running.

**Request:**
```bash
curl -X GET http://localhost:8085/health
```

**Response:**
```
OK
```

---

### 2. Get All Lessons

Retrieve all lessons with completion status for the authenticated user.

**Request:**
```bash
curl -X GET http://localhost:8085/api/education/lessons \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "title": "Introduction to Budgeting",
    "description": "Learn the basics of creating and managing a personal budget",
    "category": "Budgeting",
    "duration_minutes": 15,
    "difficulty": "beginner",
    "is_completed": true,
    "order": 1
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "title": "Understanding Savings Goals",
    "description": "How to set and achieve your savings goals",
    "category": "Savings",
    "duration_minutes": 20,
    "difficulty": "beginner",
    "is_completed": false,
    "order": 2
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "title": "Investment Basics",
    "description": "Introduction to investing and building wealth",
    "category": "Investing",
    "duration_minutes": 30,
    "difficulty": "intermediate",
    "is_completed": false,
    "order": 3
  }
]
```

---

### 3. Get Lesson Details

Retrieve detailed content for a specific lesson.

**Request:**
```bash
curl -X GET http://localhost:8085/api/education/lessons/550e8400-e29b-41d4-a716-446655440001 \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "title": "Introduction to Budgeting",
  "description": "Learn the basics of creating and managing a personal budget",
  "category": "Budgeting",
  "duration_minutes": 15,
  "difficulty": "beginner",
  "is_completed": true,
  "order": 1,
  "content": "# Introduction to Budgeting\n\nA budget is a financial plan that helps you track income and expenses...\n\n## Key Concepts\n\n1. Income tracking\n2. Expense categorization\n3. Savings allocation\n\n## Creating Your First Budget\n\nFollow these steps to create your first budget:\n\n1. List all sources of income\n2. Track your expenses for a month\n3. Categorize your spending\n4. Set spending limits for each category\n5. Monitor and adjust regularly",
  "video_url": "https://example.com/videos/budgeting-intro.mp4",
  "resources": [
    {
      "title": "Budget Template",
      "url": "https://example.com/templates/budget.xlsx",
      "type": "pdf"
    },
    {
      "title": "Budgeting Calculator",
      "url": "https://example.com/tools/calculator",
      "type": "link"
    },
    {
      "title": "50/30/20 Rule Explained",
      "url": "https://example.com/articles/50-30-20-rule",
      "type": "article"
    }
  ]
}
```

---

### 4. Mark Lesson as Complete

Mark a lesson as completed for the authenticated user.

**Request:**
```bash
curl -X POST http://localhost:8085/api/education/lessons/550e8400-e29b-41d4-a716-446655440002/complete \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "message": "Lesson marked as complete"
}
```

---

### 5. Get User Progress

Get the user's overall education progress.

**Request:**
```bash
curl -X GET http://localhost:8085/api/education/progress \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "total_lessons": 10,
  "completed_lessons": 3,
  "progress_percent": 30.0
}
```

---

## Error Responses

### 401 Unauthorized - Missing Token
```json
Missing authorization header
```

### 401 Unauthorized - Invalid Token
```json
Invalid or expired token
```

### 404 Not Found - Lesson Not Found
```json
Lesson not found
```

### 500 Internal Server Error
```json
Internal server error message
```

---

## Testing Workflow

### 1. Get JWT Token
First, authenticate with the auth service to get a JWT token:

```bash
curl -X POST http://localhost:8001/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Save the `access_token` from the response.

### 2. View Available Lessons
```bash
TOKEN="<your_access_token>"
curl -X GET http://localhost:8085/api/education/lessons \
  -H "Authorization: Bearer $TOKEN"
```

### 3. View Lesson Details
```bash
LESSON_ID="<lesson_id_from_previous_response>"
curl -X GET http://localhost:8085/api/education/lessons/$LESSON_ID \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Complete a Lesson
```bash
curl -X POST http://localhost:8085/api/education/lessons/$LESSON_ID/complete \
  -H "Authorization: Bearer $TOKEN"
```

### 5. Check Progress
```bash
curl -X GET http://localhost:8085/api/education/progress \
  -H "Authorization: Bearer $TOKEN"
```

---

## Sample Lesson Data

To test the service, you can insert sample lessons into the database:

```sql
-- Insert sample lessons
INSERT INTO lessons (id, title, description, category, duration_minutes, difficulty, content, video_url, resources, "order") VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Introduction to Budgeting', 'Learn the basics of creating and managing a personal budget', 'Budgeting', 15, 'beginner', 'A budget is a financial plan...', 'https://example.com/videos/budgeting-intro.mp4', '[{"title":"Budget Template","url":"https://example.com/templates/budget.xlsx","type":"pdf"}]', 1),
('550e8400-e29b-41d4-a716-446655440002', 'Understanding Savings Goals', 'How to set and achieve your savings goals', 'Savings', 20, 'beginner', 'Setting savings goals is crucial...', 'https://example.com/videos/savings-goals.mp4', '[]', 2),
('550e8400-e29b-41d4-a716-446655440003', 'Investment Basics', 'Introduction to investing and building wealth', 'Investing', 30, 'intermediate', 'Investing is the process of...', NULL, '[{"title":"Investment Guide","url":"https://example.com/guides/investing.pdf","type":"pdf"}]', 3);
```

---

## Notes

- All timestamps are in ISO 8601 format with timezone
- Progress percentage is calculated as: `(completed_lessons / total_lessons) × 100`
- Lesson completion is idempotent - marking a lesson complete multiple times is safe
- Resources array can be empty `[]`
- Video URL is optional and may be `null` or omitted
- Lessons are ordered by the `order` field
