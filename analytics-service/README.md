# Analytics Service

The Analytics Service provides spending analysis, savings pattern detection, financial health scoring, and AI-assisted recommendations for the InSavein platform.

## Features

- **Spending Analysis**: Analyze spending patterns with category breakdown, top merchants, and daily averages
- **Savings Patterns**: Detect savings patterns (consistent, irregular, improving) with frequency analysis
- **Financial Health Score**: Calculate comprehensive financial health score (0-100) based on savings, budget adherence, and consistency
- **AI-Assisted Recommendations**: Generate actionable recommendations based on spending and savings behavior
- **Monthly Reports**: Comprehensive monthly financial reports

## Architecture

The service follows a clean architecture pattern:
- `cmd/server/main.go`: Application entry point
- `internal/analytics/`: Core business logic
  - `service.go`: Service interface
  - `analytics_service.go`: Service implementation
  - `postgres_repository.go`: Database access layer
  - `types.go`: Data structures
  - `memory_cache.go`: In-memory caching
- `internal/handlers/`: HTTP request handlers
- `internal/middleware/`: Authentication middleware
- `pkg/database/`: Database connection utilities

## API Endpoints

All endpoints require JWT authentication via `Authorization: Bearer <token>` header.

### GET /api/analytics/spending

Get spending analysis for a period.

**Query Parameters:**
- `period`: Time period (week, month, quarter, year) - default: 30 days

**Response:**
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
    }
  ],
  "top_merchants": [
    {
      "merchant_name": "Whole Foods",
      "amount": 200.00,
      "count": 5
    }
  ],
  "daily_average": 40.34,
  "comparison_to_previous": 15.5,
  "trends": []
}
```

### GET /api/analytics/patterns

Get savings patterns.

**Response:**
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

### GET /api/analytics/recommendations

Get AI-assisted recommendations.

**Response:**
```json
[
  {
    "id": "uuid",
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
  }
]
```

### GET /api/analytics/health

Get financial health score.

**Response:**
```json
{
  "overall_score": 75,
  "savings_score": 80,
  "budget_score": 70,
  "consistency_score": 75,
  "insights": [
    "Excellent savings habits!",
    "Your financial health is good, with room for improvement"
  ],
  "improvement_areas": [
    "Focus on staying within your budget limits"
  ]
}
```

## Configuration

Environment variables:

- `PORT`: Server port (default: 8008)
- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (default: insavein)
- `DB_SSLMODE`: SSL mode (default: disable)
- `DB_REPLICA_HOST`: Read replica host (optional, for read-heavy operations)
- `JWT_SECRET`: JWT signing secret (required)

## Running the Service

### Local Development

```bash
# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=insavein
export JWT_SECRET=your-secret-key
export PORT=8008

# Run the service
go run cmd/server/main.go
```

### Docker

```bash
docker build -t analytics-service .
docker run -p 8008:8008 \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=yourpassword \
  -e JWT_SECRET=your-secret-key \
  analytics-service
```

## Financial Health Score Calculation

The financial health score is calculated as a weighted average of three components:

1. **Savings Score (40% weight)**: Based on savings frequency and amount
   - Frequency: 0-50 points (30 transactions in 30 days = 50 points)
   - Amount: 0-50 points ($20+ average = 50 points)

2. **Budget Score (30% weight)**: Based on budget adherence
   - 0-80% used = 100 points
   - 80-100% used = 80-50 points
   - 100-120% used = 50-20 points
   - 120%+ used = 0-20 points

3. **Consistency Score (30% weight)**: Based on streak and regularity
   - Streak: 0-60 points (30+ day streak = 60 points)
   - Regularity: 0-40 points (current/longest streak ratio)

All scores are clamped to 0-100 range. Results are cached for 1 hour to reduce computation load.

## Requirements Validation

This service implements the following requirements:

- **Requirement 13.1**: Spending analysis with total, category breakdown, top merchants, daily average
- **Requirement 13.2**: Comparison to previous period with percentage change
- **Requirement 13.3**: Savings pattern detection (consistent, irregular, improving)
- **Requirement 13.4**: AI-assisted recommendations with priority levels
- **Requirement 13.5**: Read from database replicas for performance
- **Requirement 13.6**: Cache financial health scores for 1 hour
- **Requirement 14.1**: Financial health score calculation with component scores
- **Requirement 14.2**: Weighted average (savings 40%, budget 30%, consistency 30%)
- **Requirement 14.3**: All scores are integers 0-100
- **Requirement 15.1**: JWT authentication on all endpoints

## Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Dependencies

- `github.com/gorilla/mux`: HTTP router
- `github.com/lib/pq`: PostgreSQL driver
- `github.com/google/uuid`: UUID generation
- `github.com/golang-jwt/jwt/v5`: JWT token validation
