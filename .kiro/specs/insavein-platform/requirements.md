# Requirements Document: InSavein Platform

## Introduction

InSavein is a financial discipline application designed to help young people facing high unemployment build savings habits, develop financial discipline, and work toward long-term financial independence. The platform enables users to track micro-savings, categorize spending, plan budgets, maintain financial discipline streaks, and receive AI-assisted savings strategies. This document specifies the functional and non-functional requirements derived from the technical design.

## Glossary

- **System**: The complete InSavein platform including frontend, backend microservices, and database
- **Auth_Service**: Authentication and authorization microservice
- **User_Service**: User profile management microservice
- **Savings_Service**: Savings transaction tracking microservice
- **Budget_Service**: Budget planning and spending tracking microservice
- **Goal_Service**: Financial goal management microservice
- **Education_Service**: Financial education content delivery microservice
- **Notification_Service**: Email and push notification delivery microservice
- **Analytics_Service**: Financial analysis and recommendation microservice
- **User**: A registered person using the InSavein platform
- **Transaction**: A savings or spending record with amount, date, and category
- **Streak**: Consecutive days with at least one savings transaction
- **Budget**: Monthly spending plan with category allocations
- **Goal**: Long-term financial target with milestones
- **JWT_Token**: JSON Web Token used for authentication
- **Milestone**: Intermediate checkpoint toward a financial goal
- **Alert**: Notification triggered by budget threshold breach
- **Financial_Health_Score**: Calculated metric (0-100) representing user's financial discipline

## Requirements

### Requirement 1: User Registration and Authentication

**User Story:** As a new user, I want to register for an account with my email and password, so that I can access the InSavein platform securely.

#### Acceptance Criteria

1. WHEN a user submits valid registration data, THE Auth_Service SHALL create a new user account with hashed password
2. WHEN a user registers, THE System SHALL hash the password using bcrypt with cost factor 12
3. WHEN a user provides an email already in use, THE Auth_Service SHALL return an error indicating duplicate email
4. WHEN a user submits a password shorter than 8 characters, THE Auth_Service SHALL reject the registration
5. WHEN registration succeeds, THE Auth_Service SHALL return an access token valid for 15 minutes and a refresh token valid for 7 days
6. WHEN a registered user submits valid credentials, THE Auth_Service SHALL return JWT tokens signed with HMAC-SHA256
7. WHEN a user submits invalid credentials, THE Auth_Service SHALL return an error without revealing whether email or password was incorrect
8. WHEN a user attempts login 5 times unsuccessfully within 15 minutes, THE Auth_Service SHALL temporarily block further attempts

### Requirement 2: Token Management and Session Security

**User Story:** As a user, I want my session to remain secure with automatic token refresh, so that I don't have to log in repeatedly while maintaining security.

#### Acceptance Criteria

1. WHEN the System validates a JWT token, THE Auth_Service SHALL verify the signature using the secret key
2. WHEN a token's expiration time is in the past, THE Auth_Service SHALL reject the token as expired
3. WHEN a user provides a valid refresh token, THE Auth_Service SHALL issue new access and refresh tokens
4. WHEN a user logs out, THE Auth_Service SHALL invalidate the refresh token
5. THE Auth_Service SHALL include user_id, email, and roles in every JWT token payload
6. WHEN a token signature is invalid, THE Auth_Service SHALL reject the token and return an authentication error

### Requirement 3: User Profile Management

**User Story:** As a user, I want to manage my profile information and preferences, so that I can personalize my experience on the platform.

#### Acceptance Criteria

1. WHEN a user requests their profile, THE User_Service SHALL return profile data including email, name, date of birth, and profile image URL
2. WHEN a user updates their profile, THE User_Service SHALL validate all fields and persist changes to the database
3. WHEN a user updates preferences, THE User_Service SHALL store currency, notification settings, reminder time, and theme preferences
4. WHEN a user requests account deletion, THE User_Service SHALL remove all user data including transactions, budgets, and goals
5. THE User_Service SHALL only allow users to access and modify their own profile data

### Requirement 4: Savings Transaction Recording

**User Story:** As a user, I want to record savings transactions with amount and description, so that I can track my savings progress over time.

#### Acceptance Criteria

1. WHEN a user creates a savings transaction with valid data, THE Savings_Service SHALL insert a record with amount, currency, description, category, and timestamp
2. WHEN a user attempts to create a transaction with amount ≤ 0, THE Savings_Service SHALL reject the transaction
3. WHEN a savings transaction is created, THE Savings_Service SHALL trigger an asynchronous streak update
4. WHEN a user requests savings history, THE Savings_Service SHALL return transactions ordered by creation date descending
5. WHEN a user requests savings summary, THE Savings_Service SHALL return total saved, current streak, longest streak, last saving date, monthly average, and current month total
6. THE Savings_Service SHALL store all amounts as decimal values with 2 decimal places precision

### Requirement 5: Savings Streak Calculation

**User Story:** As a user, I want to see my savings streak count, so that I can stay motivated to save consistently.

#### Acceptance Criteria

1. WHEN calculating streak, THE Savings_Service SHALL count consecutive days with at least one savings transaction
2. WHEN the last savings date is more than 1 day ago, THE Savings_Service SHALL set current streak to 0
3. WHEN the last savings date is today or yesterday, THE Savings_Service SHALL calculate current streak by counting consecutive days backward
4. WHEN multiple transactions occur on the same day, THE Savings_Service SHALL count it as one day in the streak
5. THE Savings_Service SHALL maintain longest streak as the maximum streak ever achieved by the user
6. THE Savings_Service SHALL ensure current streak never exceeds longest streak

### Requirement 6: Budget Creation and Management

**User Story:** As a user, I want to create monthly budgets with spending categories, so that I can plan and control my expenses.

#### Acceptance Criteria

1. WHEN a user creates a budget, THE Budget_Service SHALL store total budget amount, month, and category allocations
2. WHEN a user creates budget categories, THE Budget_Service SHALL store name, allocated amount, and color for each category
3. WHEN a user requests current budget, THE Budget_Service SHALL return the budget for the current month
4. WHEN a user updates a budget, THE Budget_Service SHALL modify category allocations and recalculate totals
5. THE Budget_Service SHALL enforce unique constraint on user_id and month combination
6. THE Budget_Service SHALL ensure all budget amounts are non-negative

### Requirement 7: Spending Transaction Recording

**User Story:** As a user, I want to record spending transactions against my budget categories, so that I can track where my money goes.

#### Acceptance Criteria

1. WHEN a user records spending, THE Budget_Service SHALL create a transaction with amount, category, description, merchant, and date
2. WHEN spending is recorded, THE Budget_Service SHALL increment the category's spent amount by the transaction amount
3. WHEN spending is recorded, THE Budget_Service SHALL increment the budget's total spent by the transaction amount
4. WHEN a user attempts to record spending with amount ≤ 0, THE Budget_Service SHALL reject the transaction
5. WHEN a user records spending for a future date, THE Budget_Service SHALL reject the transaction
6. THE Budget_Service SHALL ensure all spending updates are atomic using database transactions

### Requirement 8: Budget Alert Generation

**User Story:** As a user, I want to receive alerts when I approach or exceed my budget limits, so that I can adjust my spending behavior.

#### Acceptance Criteria

1. WHEN a category's spent amount reaches 80% of allocated amount, THE Budget_Service SHALL generate a warning alert
2. WHEN a category's spent amount reaches 100% of allocated amount, THE Budget_Service SHALL generate a critical alert
3. WHEN generating alerts, THE Budget_Service SHALL calculate percentage as (spent_amount / allocated_amount) × 100
4. WHEN multiple alerts exist, THE Budget_Service SHALL sort them with critical alerts before warning alerts
5. WHEN alerts of the same type exist, THE Budget_Service SHALL sort them by percentage descending
6. THE Budget_Service SHALL only generate alerts for categories with allocated_amount > 0

### Requirement 9: Financial Goal Management

**User Story:** As a user, I want to set financial goals with target amounts and dates, so that I can work toward long-term objectives.

#### Acceptance Criteria

1. WHEN a user creates a goal, THE Goal_Service SHALL store title, description, target amount, target date, and currency
2. WHEN a goal is created, THE Goal_Service SHALL initialize current amount to 0 and status to "active"
3. WHEN a user requests active goals, THE Goal_Service SHALL return all goals with status "active"
4. WHEN a user updates a goal, THE Goal_Service SHALL modify the specified fields and update the timestamp
5. WHEN a user deletes a goal, THE Goal_Service SHALL remove the goal and all associated milestones
6. THE Goal_Service SHALL calculate progress percentage as (current_amount / target_amount) × 100

### Requirement 10: Goal Progress Tracking

**User Story:** As a user, I want to contribute toward my goals and see progress updates, so that I can track my journey to achieving them.

#### Acceptance Criteria

1. WHEN a user adds a contribution to a goal, THE Goal_Service SHALL increase current_amount by the contribution amount
2. WHEN current_amount reaches or exceeds target_amount, THE Goal_Service SHALL change status to "completed"
3. WHEN updating goal progress, THE Goal_Service SHALL use database row-level locking to prevent race conditions
4. WHEN a contribution is added, THE Goal_Service SHALL update all milestones that have been reached
5. WHEN a milestone amount is reached, THE Goal_Service SHALL mark it as completed with completion timestamp
6. THE Goal_Service SHALL process milestones in ascending order by amount and stop at the first unreached milestone

### Requirement 11: Financial Education Content Delivery

**User Story:** As a user, I want to access financial education lessons, so that I can improve my financial literacy.

#### Acceptance Criteria

1. WHEN a user requests lessons, THE Education_Service SHALL return a list with title, description, category, duration, difficulty, and completion status
2. WHEN a user requests a specific lesson, THE Education_Service SHALL return detailed content including text, video URL, resources, and quiz questions
3. WHEN a user marks a lesson as complete, THE Education_Service SHALL record completion with timestamp
4. WHEN a user requests progress, THE Education_Service SHALL return total lessons, completed lessons, and progress percentage
5. THE Education_Service SHALL calculate progress percentage as (completed_lessons / total_lessons) × 100
6. THE Education_Service SHALL read lesson content from database replicas to reduce load on primary database

### Requirement 12: Notification Delivery

**User Story:** As a user, I want to receive notifications about important events, so that I stay informed about my financial activity.

#### Acceptance Criteria

1. WHEN the System needs to send an email, THE Notification_Service SHALL use the configured email provider with template and recipient data
2. WHEN the System needs to send a push notification, THE Notification_Service SHALL deliver it to the user's registered devices
3. WHEN a reminder is scheduled, THE Notification_Service SHALL store it with user_id, type, scheduled time, and message
4. WHEN a user requests their notifications, THE Notification_Service SHALL return all notifications ordered by creation date descending
5. WHEN a user marks a notification as read, THE Notification_Service SHALL update the is_read flag to true
6. WHERE notification preferences are disabled, THE Notification_Service SHALL not send notifications to that user

### Requirement 13: Spending Analysis and Insights

**User Story:** As a user, I want to see analysis of my spending patterns, so that I can understand my financial behavior and make improvements.

#### Acceptance Criteria

1. WHEN a user requests spending analysis, THE Analytics_Service SHALL calculate total spending, category breakdown, top merchants, and daily average for the specified period
2. WHEN analyzing spending, THE Analytics_Service SHALL compare current period to previous period and calculate percentage change
3. WHEN identifying patterns, THE Analytics_Service SHALL determine pattern type as "consistent", "irregular", or "improving"
4. WHEN generating recommendations, THE Analytics_Service SHALL provide actionable items with priority levels
5. THE Analytics_Service SHALL read data from database replicas to avoid impacting write performance
6. THE Analytics_Service SHALL cache financial health scores for 1 hour to reduce computation load

### Requirement 14: Financial Health Score Calculation

**User Story:** As a user, I want to see my financial health score, so that I can gauge my overall financial discipline.

#### Acceptance Criteria

1. WHEN calculating financial health, THE Analytics_Service SHALL compute savings score, budget score, and consistency score
2. WHEN computing overall score, THE Analytics_Service SHALL use weighted average with savings (40%), budget (30%), and consistency (30%)
3. THE Analytics_Service SHALL ensure all scores are integers between 0 and 100 inclusive
4. WHEN a user has less than 30 days of history, THE Analytics_Service SHALL return an error indicating insufficient data
5. WHEN calculating scores, THE Analytics_Service SHALL include insights and improvement areas in the response
6. THE Analytics_Service SHALL base savings score on frequency and amount of savings transactions

### Requirement 15: API Request Authentication

**User Story:** As a system administrator, I want all API requests to be authenticated, so that user data remains secure.

#### Acceptance Criteria

1. WHEN a request is received without a valid token, THE System SHALL return HTTP 401 Unauthorized
2. WHEN a request is received with an expired token, THE System SHALL return HTTP 401 Unauthorized with message "Invalid or expired token"
3. WHEN a request is received with a valid token, THE System SHALL extract user_id from token claims and pass it to the service
4. WHEN a user attempts to access another user's data, THE System SHALL return HTTP 403 Forbidden
5. THE System SHALL validate token signature on every authenticated request
6. THE System SHALL include token expiration time in error responses for expired tokens

### Requirement 16: Database Transaction Integrity

**User Story:** As a system administrator, I want all multi-step operations to be atomic, so that data remains consistent even during failures.

#### Acceptance Criteria

1. WHEN an operation involves multiple database tables, THE System SHALL use database transactions to ensure atomicity
2. WHEN a transaction fails, THE System SHALL roll back all changes made within that transaction
3. WHEN a transaction succeeds, THE System SHALL commit all changes atomically
4. WHEN updating goal progress, THE System SHALL lock the goal row to prevent concurrent modification
5. WHEN recording spending, THE System SHALL update both spending_transactions and budget_categories tables atomically
6. THE System SHALL ensure budget total_spent always equals the sum of category spent_amounts

### Requirement 17: Input Validation and Error Handling

**User Story:** As a user, I want clear error messages when I submit invalid data, so that I can correct my input and try again.

#### Acceptance Criteria

1. WHEN a user submits invalid data, THE System SHALL return HTTP 400 Bad Request with detailed validation errors
2. WHEN a required field is missing, THE System SHALL include the field name in the error message
3. WHEN a field value is out of range, THE System SHALL specify the valid range in the error message
4. WHEN a database constraint is violated, THE System SHALL return a user-friendly error message without exposing internal details
5. THE System SHALL validate all inputs against defined schemas before processing
6. THE System SHALL sanitize all user inputs to prevent SQL injection and XSS attacks

### Requirement 18: API Rate Limiting

**User Story:** As a system administrator, I want to rate limit API requests, so that the system remains available and prevents abuse.

#### Acceptance Criteria

1. WHEN a user exceeds 100 requests per minute, THE System SHALL return HTTP 429 Too Many Requests
2. WHEN an IP address exceeds 1000 requests per minute, THE System SHALL return HTTP 429 Too Many Requests
3. WHEN rate limits are enforced, THE System SHALL include rate limit headers in responses
4. THE System SHALL allow burst requests up to 20 above the rate limit
5. WHEN a user is rate limited, THE System SHALL include retry-after header indicating when to retry
6. THE System SHALL reset rate limit counters every minute

### Requirement 19: Service Health Monitoring

**User Story:** As a system administrator, I want health check endpoints on all services, so that I can monitor system availability.

#### Acceptance Criteria

1. WHEN a health check request is received, THE System SHALL return HTTP 200 OK if the service is healthy
2. WHEN a service cannot connect to the database, THE System SHALL return HTTP 503 Service Unavailable
3. THE System SHALL provide separate liveness and readiness probes for Kubernetes
4. WHEN a liveness probe fails, THE System SHALL indicate the service needs to be restarted
5. WHEN a readiness probe fails, THE System SHALL indicate the service should not receive traffic
6. THE System SHALL respond to health checks within 1 second

### Requirement 20: Data Encryption and Security

**User Story:** As a user, I want my sensitive data to be encrypted, so that my financial information remains private.

#### Acceptance Criteria

1. THE System SHALL never store passwords in plaintext
2. THE System SHALL never return password hashes in API responses
3. WHEN storing data at rest, THE System SHALL use database-level encryption
4. WHEN transmitting data, THE System SHALL require TLS 1.3 for all connections
5. THE System SHALL encrypt database connections between services and PostgreSQL
6. THE System SHALL rotate encryption keys every 90 days

### Requirement 21: Horizontal Scaling and Load Balancing

**User Story:** As a system administrator, I want services to scale automatically based on load, so that the system handles traffic spikes gracefully.

#### Acceptance Criteria

1. WHEN CPU utilization exceeds 70%, THE System SHALL scale up the number of service replicas
2. WHEN memory utilization exceeds 80%, THE System SHALL scale up the number of service replicas
3. WHEN load decreases, THE System SHALL scale down to minimum replica count
4. THE System SHALL maintain at least 3 replicas for Auth_Service, Savings_Service, and Budget_Service
5. THE System SHALL distribute requests across healthy replicas using round-robin load balancing
6. WHEN a replica fails health checks, THE System SHALL stop routing traffic to that replica

### Requirement 22: Database Replication and Failover

**User Story:** As a system administrator, I want database replication with automatic failover, so that the system remains available during database failures.

#### Acceptance Criteria

1. THE System SHALL maintain at least 2 read replicas for the PostgreSQL database
2. WHEN the primary database fails, THE System SHALL promote a replica to primary within 30 seconds
3. THE System SHALL route read-heavy operations to replicas to reduce primary database load
4. THE System SHALL route all write operations to the primary database
5. THE System SHALL monitor replication lag and alert when it exceeds 1 second
6. THE System SHALL ensure replication lag remains below 5 seconds under normal load

### Requirement 23: Observability and Metrics

**User Story:** As a system administrator, I want comprehensive metrics and logs, so that I can troubleshoot issues and monitor performance.

#### Acceptance Criteria

1. WHEN a request is processed, THE System SHALL record request count, duration, and status code as Prometheus metrics
2. WHEN an error occurs, THE System SHALL log it with timestamp, service name, trace ID, and error details in JSON format
3. THE System SHALL expose metrics endpoint on port 9090 for Prometheus scraping
4. THE System SHALL include trace IDs in all logs to enable distributed tracing
5. THE System SHALL retain INFO logs for 30 days and ERROR logs for 90 days
6. THE System SHALL sample 10% of requests for distributed tracing in production

### Requirement 24: Deployment and Rollback

**User Story:** As a system administrator, I want zero-downtime deployments with quick rollback capability, so that updates don't disrupt users.

#### Acceptance Criteria

1. WHEN deploying a new version, THE System SHALL use rolling updates with max surge 25% and max unavailable 0
2. WHEN a new replica fails health checks, THE System SHALL automatically roll back the deployment
3. THE System SHALL wait for new replicas to pass readiness checks before terminating old replicas
4. WHEN a deployment is rolled back, THE System SHALL restore the previous version within 2 minutes
5. THE System SHALL run smoke tests after deployment and alert on failures
6. THE System SHALL maintain the previous version's container images for quick rollback

### Requirement 25: Data Backup and Recovery

**User Story:** As a system administrator, I want automated database backups, so that data can be recovered in case of catastrophic failure.

#### Acceptance Criteria

1. THE System SHALL create full database backups daily at 2 AM UTC
2. THE System SHALL retain daily backups for 30 days
3. THE System SHALL create incremental backups every 6 hours
4. THE System SHALL store backups in geographically separate locations
5. THE System SHALL test backup restoration monthly to verify integrity
6. WHEN a backup fails, THE System SHALL alert administrators immediately

## Non-Functional Requirements

### Performance Requirements

1. THE System SHALL respond to API requests with p95 latency less than 500ms
2. THE System SHALL respond to API requests with p99 latency less than 1000ms
3. THE System SHALL support at least 100,000 requests per minute
4. THE System SHALL maintain error rate below 0.1% under normal load
5. THE System SHALL execute database queries with p95 latency less than 100ms

### Scalability Requirements

1. THE System SHALL support at least 1 million registered users
2. THE System SHALL handle at least 10,000 concurrent users
3. THE System SHALL scale to 50,000 concurrent users during peak load
4. THE System SHALL partition transaction tables by month to maintain query performance
5. THE System SHALL use connection pooling with maximum 20 connections per service

### Availability Requirements

1. THE System SHALL maintain 99.9% uptime (less than 43 minutes downtime per month)
2. THE System SHALL recover from service failures within 30 seconds using Kubernetes health checks
3. THE System SHALL provide graceful degradation when non-critical services are unavailable
4. THE System SHALL cache frequently accessed data to maintain availability during database issues

### Security Requirements

1. THE System SHALL enforce HTTPS for all client connections
2. THE System SHALL implement CORS policies to prevent unauthorized cross-origin requests
3. THE System SHALL include security headers (X-Frame-Options, X-Content-Type-Options, CSP, HSTS)
4. THE System SHALL log all authentication failures for security auditing
5. THE System SHALL comply with GDPR requirements for data privacy and user rights

### Maintainability Requirements

1. THE System SHALL achieve at least 80% code coverage for business logic
2. THE System SHALL use structured logging in JSON format for all services
3. THE System SHALL document all API endpoints with OpenAPI specifications
4. THE System SHALL use semantic versioning for all service releases
5. THE System SHALL maintain separate staging and production environments

### Compatibility Requirements

1. THE System SHALL support modern web browsers (Chrome, Firefox, Safari, Edge) released within the last 2 years
2. THE System SHALL support mobile browsers on iOS 14+ and Android 10+
3. THE System SHALL provide responsive UI that works on screen sizes from 320px to 4K
4. THE System SHALL use UTF-8 encoding for all text data to support international characters

