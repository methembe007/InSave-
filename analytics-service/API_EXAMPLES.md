# Analytics Service API Examples

This document provides example API requests and responses for testing the Analytics Service.

## Authentication

All API endpoints require a JWT token in the Authorization header:

```bash
Authorization: Bearer <your-jwt-token>
```

## 1. Get Spending Analysis

### Request - Last 30 Days (Default)

```bash
curl -X GET http://localhost:8008/api/analytics/spending \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Request - Last Week

```bash
curl -X GET "http://localhost:8008/api/analytics/spending?period=week" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Request - Last Month

```bash
curl -X GET "http://localhost:8008/api/analytics/spending?period=month" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Request - Last Quarter

```bash
curl -X GET "http://localhost:8008/api/analytics/spending?period=quarter" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Request - Last Year

```bash
curl -X GET "http://localhost:8008/api/analytics/spending?period=year" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Response Example

```json
{
  "period": {
    "start": "2024-01-01T00:00:00Z",
    "end": "2024-01-31T23:59:59Z"
  },
  "total_spending": 1250.50,
  "category_breakdown": [
    {
      "category_name": "groceries",
      "amount": 450.00,
      "percentage": 36.0,
      "count": 12
    },
    {
      "category_name": "transportation",
      "amount": 300.00,
      "percentage": 24.0,
      "count": 8
    },
    {
      "category_name": "entertainment",
      "amount": 250.50,
      "percentage": 20.0,
      "count": 15
    },
    {
      "category_name": "utilities",
      "amount": 150.00,
      "percentage": 12.0,
      "count": 3
    },
    {
      "category_name": "other",
      "amount": 100.00,
      "percentage": 8.0,
      "count": 5
    }
  ],
  "top_merchants": [
    {
      "merchant_name": "Whole Foods",
      "amount": 200.00,
      "count": 5
    },
    {
      "merchant_name": "Shell Gas Station",
      "amount": 150.00,
      "count": 4
    },
    {
      "merchant_name": "Netflix",
      "amount": 120.00,
      "count": 3
    },
    {
      "merchant_name": "Amazon",
      "amount": 100.50,
      "count": 8
    },
    {
      "merchant_name": "Starbucks",
      "amount": 80.00,
      "count": 10
    }
  ],
  "daily_average": 40.34,
  "comparison_to_previous": 15.5,
  "trends": []
}
```

## 2. Get Savings Patterns

### Request

```bash
curl -X GET http://localhost:8008/api/analytics/patterns \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Response Example - Consistent Pattern

```json
[
  {
    "pattern_type": "consistent",
    "average_amount": 25.50,
    "frequency": "weekly",
    "best_day_of_week": "Friday",
    "insights": [
      "You're saving consistently with an average of $25.50 per transaction",
      "You save most frequently on Fridays"
    ]
  }
]
```

### Response Example - Improving Pattern

```json
[
  {
    "pattern_type": "improving",
    "average_amount": 30.75,
    "frequency": "bi-weekly",
    "best_day_of_week": "Monday",
    "insights": [
      "Your savings amounts are increasing over time - great progress!",
      "You save most frequently on Mondays"
    ]
  }
]
```

### Response Example - Irregular Pattern

```json
[
  {
    "pattern_type": "irregular",
    "average_amount": 15.25,
    "frequency": "irregular",
    "best_day_of_week": "Wednesday",
    "insights": [
      "Your savings pattern is irregular. Try setting a regular savings schedule.",
      "You save most frequently on Wednesdays"
    ]
  }
]
```

## 3. Get Recommendations

### Request

```bash
curl -X GET http://localhost:8008/api/analytics/recommendations \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Response Example

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "spending",
    "priority": "high",
    "title": "Reduce groceries spending",
    "description": "You're spending 45.0% of your budget on groceries. Consider reducing this category.",
    "action_items": [
      "Set a lower budget for groceries",
      "Track your spending more carefully in this category",
      "Look for cheaper alternatives"
    ],
    "potential_savings": 90.00
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "type": "budget",
    "priority": "high",
    "title": "Spending increased significantly",
    "description": "Your spending increased by 22.5% compared to the previous period.",
    "action_items": [
      "Review your recent transactions",
      "Identify unnecessary expenses",
      "Set stricter budget limits"
    ],
    "potential_savings": 187.58
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "type": "savings",
    "priority": "medium",
    "title": "Establish a regular savings schedule",
    "description": "Your savings pattern is irregular. Setting up automatic savings can help build consistency.",
    "action_items": [
      "Set up automatic transfers on payday",
      "Start with a small, manageable amount",
      "Gradually increase your savings rate"
    ],
    "potential_savings": 61.00
  }
]
```

## 4. Get Financial Health Score

### Request

```bash
curl -X GET http://localhost:8008/api/analytics/health \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Response Example - Excellent Health

```json
{
  "overall_score": 85,
  "savings_score": 90,
  "budget_score": 85,
  "consistency_score": 80,
  "insights": [
    "Excellent savings habits!",
    "Great budget adherence!",
    "You're maintaining excellent consistency!",
    "Your financial health is excellent!"
  ],
  "improvement_areas": []
}
```

### Response Example - Good Health

```json
{
  "overall_score": 68,
  "savings_score": 75,
  "budget_score": 65,
  "consistency_score": 60,
  "insights": [
    "Your financial health is good, with room for improvement"
  ],
  "improvement_areas": [
    "Focus on staying within your budget limits",
    "Build a more consistent savings routine"
  ]
}
```

### Response Example - Needs Improvement

```json
{
  "overall_score": 42,
  "savings_score": 45,
  "budget_score": 40,
  "consistency_score": 40,
  "insights": [
    "Focus on building better financial habits"
  ],
  "improvement_areas": [
    "Increase your savings frequency and amounts",
    "Focus on staying within your budget limits",
    "Build a more consistent savings routine"
  ]
}
```

### Error Response - Insufficient Data

```json
{
  "error": "Insufficient data: need at least 30 days of transaction history"
}
```

## 5. Health Check

### Request

```bash
curl -X GET http://localhost:8008/health
```

### Response

```json
{
  "status": "healthy"
}
```

## Error Responses

### 401 Unauthorized - Missing Token

```json
{
  "error": "Missing authorization header"
}
```

### 401 Unauthorized - Invalid Token

```json
{
  "error": "Invalid or expired token"
}
```

### 400 Bad Request - Insufficient Data

```json
{
  "error": "Insufficient data: need at least 30 days of transaction history"
}
```

### 500 Internal Server Error

```json
{
  "error": "Failed to get spending analysis"
}
```

## Testing with curl

### Complete Example with JWT Token

```bash
# Set your JWT token
export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDU2NzgtMTIzNC0xMjM0LTEyMzQtMTIzNDU2Nzg5MDEyIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwiZXhwIjoxNzM1MDAwMDAwfQ.signature"

# Test health check
curl -X GET http://localhost:8008/health

# Test spending analysis
curl -X GET http://localhost:8008/api/analytics/spending \
  -H "Authorization: Bearer $JWT_TOKEN"

# Test savings patterns
curl -X GET http://localhost:8008/api/analytics/patterns \
  -H "Authorization: Bearer $JWT_TOKEN"

# Test recommendations
curl -X GET http://localhost:8008/api/analytics/recommendations \
  -H "Authorization: Bearer $JWT_TOKEN"

# Test financial health
curl -X GET http://localhost:8008/api/analytics/health \
  -H "Authorization: Bearer $JWT_TOKEN"
```

## Testing with Postman

1. **Import Collection**: Create a new collection named "Analytics Service"

2. **Set Environment Variables**:
   - `base_url`: `http://localhost:8008`
   - `jwt_token`: Your JWT token

3. **Add Requests**:
   - GET `{{base_url}}/health`
   - GET `{{base_url}}/api/analytics/spending?period=month`
   - GET `{{base_url}}/api/analytics/patterns`
   - GET `{{base_url}}/api/analytics/recommendations`
   - GET `{{base_url}}/api/analytics/health`

4. **Set Authorization**: For all `/api/analytics/*` requests:
   - Type: Bearer Token
   - Token: `{{jwt_token}}`

## Score Interpretation

### Overall Score Ranges
- **80-100**: Excellent financial health
- **60-79**: Good financial health with room for improvement
- **40-59**: Fair financial health, needs attention
- **0-39**: Poor financial health, requires immediate action

### Component Scores

**Savings Score (40% weight)**:
- Based on frequency and amount of savings
- 30 transactions in 30 days with $20 average = 100 points

**Budget Score (30% weight)**:
- Based on budget adherence
- 0-80% budget used = 100 points
- 100% budget used = 50 points
- 120%+ budget used = 0-20 points

**Consistency Score (30% weight)**:
- Based on savings streak and regularity
- 30+ day streak = 60 points
- Current streak = longest streak = 40 points
