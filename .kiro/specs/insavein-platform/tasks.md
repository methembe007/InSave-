# Implementation Plan: InSavein Platform

## Overview

This implementation plan breaks down the InSavein financial discipline platform into discrete, actionable coding tasks. The platform consists of a TanStack Start frontend (TypeScript/React), 8 Golang microservices, PostgreSQL database with replication, Kubernetes orchestration, and comprehensive observability stack. Tasks are organized to enable incremental development with early validation through testing.

## Implementation Strategy

- Build infrastructure foundation first (database, Kubernetes configs)
- Implement backend microservices in dependency order (Auth → User → domain services)
- Develop frontend components alongside backend APIs
- Integrate observability and monitoring throughout
- Deploy CI/CD pipeline for automated testing and deployment

## Tasks

- [ ] 1. Infrastructure and Database Setup
  - [x] 1.1 Set up PostgreSQL database schema and migrations
    - Create all database tables (users, savings_transactions, budgets, budget_categories, spending_transactions, goals, goal_milestones, notifications, lessons, education_progress)
    - Implement table partitioning for savings_transactions and spending_transactions by month
    - Create all indexes for query optimization
    - Write migration scripts using golang-migrate or similar tool
    - _Requirements: 1.1, 4.1, 6.1, 7.1, 9.1, 11.1, 12.1_
  
  - [x] 1.2 Configure PostgreSQL replication setup
    - Set up primary-replica replication configuration
    - Configure 2 read replicas for read-heavy operations
    - Implement connection pooling with PgBouncer
    - Test replication lag monitoring
    - _Requirements: 11.6, 13.5, 19.1_
  
  - [x] 1.3 Create Kubernetes namespace and base configurations
    - Create insavein namespace
    - Set up ConfigMaps for environment variables
    - Create Secrets for database credentials, JWT secret, API keys
    - Configure resource quotas and limits
    - _Requirements: 18.1, 20.1_


- [x] 2. Auth Service Implementation
  - [x] 2.1 Create Auth Service project structure and core interfaces
    - Initialize Go module for auth-service
    - Define Service interface with Register, Login, RefreshToken, ValidateToken, Logout methods
    - Create request/response structs (RegisterRequest, LoginRequest, AuthResponse, TokenClaims)
    - Set up project structure (cmd, internal, pkg directories)
    - _Requirements: 1.1, 1.5, 2.1_
  
  - [x] 2.2 Implement user registration with password hashing
    - Implement Register method with bcrypt password hashing (cost factor 12)
    - Add email uniqueness validation
    - Add password length validation (minimum 8 characters)
    - Implement database insertion for new users
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 20.2_
  
  - [ ]* 2.3 Write property test for password security
    - **Property 1: Password Security**
    - **Validates: Requirements 1.2, 20.1, 20.2**
    - Test that passwords are always hashed with bcrypt cost ≥ 12
    - Test that plaintext passwords never appear in responses
  
  - [x] 2.4 Implement JWT token generation and validation
    - Implement Login method with credential verification
    - Generate JWT access tokens (15 min expiry) and refresh tokens (7 days expiry)
    - Use HMAC-SHA256 for token signing
    - Include user_id, email, roles in token payload
    - Implement ValidateToken method with signature and expiration checks
    - _Requirements: 1.5, 1.6, 2.1, 2.2, 2.5, 2.6, 15.5_
  
  - [ ]* 2.5 Write property tests for JWT token operations
    - **Property 2: JWT Token Validity**
    - **Validates: Requirements 2.1, 2.2, 2.5, 2.6, 15.5**
    - **Property 3: Token Refresh Round Trip**
    - **Validates: Requirements 1.5, 2.3**
  
  - [x] 2.6 Implement token refresh and logout functionality
    - Implement RefreshToken method to issue new tokens
    - Implement Logout method with token invalidation
    - Add token revocation list (in-memory or Redis)
    - _Requirements: 1.5, 2.3, 2.4_
  
  - [x] 2.7 Add rate limiting for login attempts
    - Implement rate limiter (5 attempts per 15 minutes per email)
    - Add temporary blocking on exceeded attempts
    - Return appropriate error messages
    - _Requirements: 1.7, 1.8_
  
  - [x] 2.8 Create HTTP handlers and routes for Auth Service
    - Implement POST /api/auth/register handler
    - Implement POST /api/auth/login handler
    - Implement POST /api/auth/refresh handler
    - Implement POST /api/auth/logout handler
    - Add input validation middleware
    - _Requirements: 1.1, 1.5, 17.1, 17.2_
  
  - [ ]* 2.9 Write unit tests for Auth Service
    - Test registration with valid/invalid data
    - Test login with correct/incorrect credentials
    - Test token validation edge cases
    - Test rate limiting behavior


- [x] 3. User Profile Service Implementation
  - [x] 3.1 Create User Service project structure and interfaces
    - Initialize Go module for user-service
    - Define Service interface with GetProfile, UpdateProfile, GetPreferences, UpdatePreferences, DeleteAccount methods
    - Create UserProfile and UserPreferences structs
    - _Requirements: 3.1, 3.2, 3.3_
  
  - [x] 3.2 Implement profile retrieval and update operations
    - Implement GetProfile method to fetch user data
    - Implement UpdateProfile method with field validation
    - Add authorization check (users can only access own profile)
    - _Requirements: 3.1, 3.2, 3.5, 15.4_
  
  - [x] 3.3 Implement user preferences management
    - Implement GetPreferences method
    - Implement UpdatePreferences method for currency, notifications, theme settings
    - Store preferences as JSONB in database
    - _Requirements: 3.3, 12.6_
  
  - [ ]* 3.4 Write property test for profile update round trip
    - **Property 27: Profile Update Round Trip**
    - **Validates: Requirements 3.2, 3.3**
    - Test that reading profile after update returns updated values
  
  - [x] 3.5 Implement account deletion with cascade
    - Implement DeleteAccount method
    - Ensure cascade deletion of all user data (transactions, budgets, goals)
    - Use database transaction for atomicity
    - _Requirements: 3.4, 16.1, 16.2_
  
  - [ ]* 3.6 Write property test for cascade deletion
    - **Property 21: Cascade Deletion Completeness**
    - **Validates: Requirements 3.4, 9.5**
  
  - [x] 3.7 Create HTTP handlers and routes for User Service
    - Implement GET /api/user/profile handler
    - Implement PUT /api/user/profile handler
    - Implement GET /api/user/preferences handler
    - Implement PUT /api/user/preferences handler
    - Implement DELETE /api/user/account handler
    - Add authentication middleware
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 15.1, 15.2_
  
  - [ ]* 3.8 Write unit tests for User Service
    - Test profile CRUD operations
    - Test authorization enforcement
    - Test preferences update


- [x] 4. Savings Service Implementation
  - [x] 4.1 Create Savings Service project structure and interfaces
    - Initialize Go module for savings-service
    - Define Service interface with GetSummary, GetHistory, CreateTransaction, GetStreak, UpdateStreak methods
    - Create SavingsTransaction, SavingsSummary, SavingsStreak structs
    - _Requirements: 4.1, 4.4, 4.5, 5.1_
  
  - [x] 4.2 Implement savings transaction creation
    - Implement CreateTransaction method with amount validation (must be > 0)
    - Store amount as decimal with 2 decimal places precision
    - Insert transaction into partitioned savings_transactions table
    - Trigger asynchronous streak update
    - _Requirements: 4.1, 4.2, 4.3, 4.6_
  
  - [ ]* 4.3 Write property tests for savings transaction integrity
    - **Property 4: Savings Transaction Integrity**
    - **Validates: Requirements 4.1, 4.2, 4.6**
    - Test amount is always positive
    - Test user_id references existing user
    - Test created_at is not in future
    - Test decimal precision is exactly 2 places
  
  - [x] 4.4 Implement savings streak calculation algorithm
    - Implement UpdateStreak method with consecutive day counting logic
    - Set current streak to 0 if last save > 1 day ago
    - Calculate current streak by counting backward from last save date
    - Handle multiple transactions on same day (count as one day)
    - Update longest streak if current exceeds it
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6_
  
  - [ ]* 4.5 Write property tests for streak calculation
    - **Property 5: Savings Streak Invariant**
    - **Validates: Requirements 5.2, 5.5, 5.6**
    - **Property 6: Streak Calculation Determinism**
    - **Validates: Requirements 5.1, 5.3, 5.4**
    - **Property 30: Streak Update Trigger**
    - **Validates: Requirements 4.3**
  
  - [x] 4.6 Implement savings history and summary retrieval
    - Implement GetHistory method with date descending order
    - Implement GetSummary method calculating total saved, streaks, monthly stats
    - Query partitioned tables efficiently
    - _Requirements: 4.4, 4.5_
  
  - [x] 4.7 Create HTTP handlers and routes for Savings Service
    - Implement POST /api/savings/transactions handler
    - Implement GET /api/savings/history handler with pagination
    - Implement GET /api/savings/summary handler
    - Implement GET /api/savings/streak handler
    - Add authentication and authorization middleware
    - _Requirements: 4.1, 4.4, 4.5, 15.1, 15.4_
  
  - [ ]* 4.8 Write unit tests for Savings Service
    - Test transaction creation with valid/invalid amounts
    - Test streak calculation with various date sequences
    - Test summary aggregation accuracy


- [x] 5. Budget Service Implementation
  - [x] 5.1 Create Budget Service project structure and interfaces
    - Initialize Go module for budget-service
    - Define Service interface with GetCurrentBudget, CreateBudget, UpdateBudget, GetCategories, RecordSpending, CheckBudgetAlerts methods
    - Create Budget, BudgetCategory, SpendingTransaction, BudgetAlert structs
    - _Requirements: 6.1, 6.2, 7.1, 8.1_
  
  - [x] 5.2 Implement budget creation and management
    - Implement CreateBudget method with category allocations
    - Enforce unique constraint on (user_id, month)
    - Validate all amounts are non-negative
    - Implement UpdateBudget method for modifying allocations
    - _Requirements: 6.1, 6.2, 6.4, 6.5, 6.6_
  
  - [x] 5.3 Implement spending transaction recording with atomic updates
    - Implement RecordSpending method with amount validation (must be > 0)
    - Reject future-dated transactions
    - Use database transaction to atomically update spending_transactions, budget_categories.spent_amount, and budgets.total_spent
    - Implement rollback on any failure
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 16.1, 16.2, 16.3_
  
  - [ ]* 5.4 Write property tests for budget consistency
    - **Property 7: Budget Consistency Invariant**
    - **Validates: Requirements 6.5, 6.6, 7.2, 7.3, 16.6**
    - **Property 8: Spending Transaction Atomicity**
    - **Validates: Requirements 7.6, 16.2, 16.3, 16.5**
    - **Property 26: Transaction Rollback on Failure**
    - **Validates: Requirements 16.1, 16.2**
  
  - [x] 5.5 Implement budget alert detection algorithm
    - Implement CheckBudgetAlerts method
    - Generate warning alerts for categories at 80-99% spent
    - Generate critical alerts for categories at 100%+ spent
    - Skip categories with zero allocated amount
    - Sort alerts: critical before warning, then by percentage descending
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6_
  
  - [ ]* 5.6 Write property tests for budget alerts
    - **Property 9: Budget Alert Thresholds**
    - **Validates: Requirements 8.1, 8.2, 8.3, 8.6**
    - **Property 10: Alert Sorting Order**
    - **Validates: Requirements 8.4, 8.5**
  
  - [x] 5.7 Create HTTP handlers and routes for Budget Service
    - Implement POST /api/budget handler
    - Implement GET /api/budget/current handler
    - Implement PUT /api/budget/:id handler
    - Implement POST /api/budget/spending handler
    - Implement GET /api/budget/alerts handler
    - Add authentication and authorization middleware
    - _Requirements: 6.1, 6.3, 6.4, 7.1, 8.1, 15.1, 15.4_
  
  - [ ]* 5.8 Write unit tests for Budget Service
    - Test budget creation with valid/invalid data
    - Test spending recording and atomic updates
    - Test alert generation at various thresholds


- [x] 6. Goal Service Implementation
  - [x] 6.1 Create Goal Service project structure and interfaces
    - Initialize Go module for goal-service
    - Define Service interface with GetActiveGoals, GetGoal, CreateGoal, UpdateGoal, DeleteGoal, GetMilestones, UpdateProgress methods
    - Create Goal, GoalDetail, Milestone structs
    - _Requirements: 9.1, 9.2, 9.3, 10.1_
  
  - [x] 6.2 Implement goal CRUD operations
    - Implement CreateGoal method with validation
    - Initialize current_amount to 0 and status to "active"
    - Implement GetActiveGoals to filter by status
    - Implement UpdateGoal and DeleteGoal with cascade for milestones
    - Calculate progress_percent as (current_amount / target_amount) × 100
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5, 9.6_
  
  - [x] 6.3 Implement goal progress update with concurrency control
    - Implement UpdateProgress method with database row-level locking (FOR UPDATE)
    - Increase current_amount by contribution amount
    - Change status to "completed" when current_amount >= target_amount
    - Use database transaction for atomicity
    - _Requirements: 10.1, 10.2, 10.3, 16.4_
  
  - [ ]* 6.4 Write property tests for goal progress
    - **Property 11: Goal Progress Calculation**
    - **Validates: Requirements 9.6, 10.2**
    - **Property 12: Goal Contribution Monotonicity**
    - **Validates: Requirements 10.1**
    - **Property 14: Goal Update Concurrency Safety**
    - **Validates: Requirements 10.3, 16.4**
  
  - [x] 6.5 Implement milestone tracking and completion
    - Implement GetMilestones method
    - In UpdateProgress, check and mark milestones as completed
    - Process milestones in ascending order by amount
    - Set completed_at timestamp when milestone reached
    - Stop at first unreached milestone
    - _Requirements: 10.4, 10.5, 10.6_
  
  - [ ]* 6.6 Write property test for milestone completion order
    - **Property 13: Milestone Completion Order**
    - **Validates: Requirements 10.4, 10.5, 10.6**
  
  - [x] 6.7 Create HTTP handlers and routes for Goal Service
    - Implement POST /api/goals handler
    - Implement GET /api/goals handler (active goals)
    - Implement GET /api/goals/:id handler
    - Implement PUT /api/goals/:id handler
    - Implement DELETE /api/goals/:id handler
    - Implement POST /api/goals/:id/progress handler
    - Implement GET /api/goals/:id/milestones handler
    - Add authentication and authorization middleware
    - _Requirements: 9.1, 9.3, 9.4, 9.5, 10.1, 15.1, 15.4_
  
  - [ ]* 6.8 Write unit tests for Goal Service
    - Test goal creation and CRUD operations
    - Test progress updates with concurrent requests
    - Test milestone completion logic


- [x] 7. Education Service Implementation
  - [x] 7.1 Create Education Service project structure and interfaces
    - Initialize Go module for education-service
    - Define Service interface with GetLessons, GetLesson, MarkLessonComplete, GetUserProgress methods
    - Create Lesson, LessonDetail, EducationProgress structs
    - _Requirements: 11.1, 11.2, 11.3, 11.4_
  
  - [x] 7.2 Implement lesson retrieval from read replicas
    - Implement GetLessons method querying database replicas
    - Return lessons with completion status for authenticated user
    - Implement GetLesson method for detailed content
    - _Requirements: 11.1, 11.2, 11.6_
  
  - [x] 7.3 Implement lesson completion tracking
    - Implement MarkLessonComplete method
    - Insert or update education_progress record with completion timestamp
    - _Requirements: 11.3_
  
  - [x] 7.4 Implement education progress calculation
    - Implement GetUserProgress method
    - Calculate progress_percent as (completed_lessons / total_lessons) × 100
    - Return total, completed, and percentage
    - _Requirements: 11.4, 11.5_
  
  - [ ]* 7.5 Write property test for education progress
    - **Property 17: Education Progress Calculation**
    - **Validates: Requirements 11.4, 11.5**
  
  - [x] 7.6 Create HTTP handlers and routes for Education Service
    - Implement GET /api/education/lessons handler
    - Implement GET /api/education/lessons/:id handler
    - Implement POST /api/education/lessons/:id/complete handler
    - Implement GET /api/education/progress handler
    - Add authentication middleware
    - _Requirements: 11.1, 11.2, 11.3, 11.4, 15.1_
  
  - [ ]* 7.7 Write unit tests for Education Service
    - Test lesson retrieval
    - Test completion tracking
    - Test progress calculation


- [x] 8. Notification Service Implementation
  - [x] 8.1 Create Notification Service project structure and interfaces
    - Initialize Go module for notification-service
    - Define Service interface with SendEmail, SendPushNotification, ScheduleReminder, GetUserNotifications, MarkAsRead methods
    - Create EmailRequest, PushNotificationRequest, ReminderRequest, Notification structs
    - _Requirements: 12.1, 12.2, 12.3, 12.4_
  
  - [x] 8.2 Implement email notification delivery
    - Implement SendEmail method integrating with SendGrid or AWS SES
    - Support template-based emails
    - Handle email delivery errors gracefully
    - _Requirements: 12.1_
  
  - [x] 8.3 Implement push notification delivery
    - Implement SendPushNotification method integrating with Firebase Cloud Messaging
    - Support both mobile and web push notifications
    - _Requirements: 12.2_
  
  - [x] 8.4 Implement notification preference enforcement
    - Check user preferences before sending notifications
    - Skip sending if notifications disabled for user
    - _Requirements: 12.6_
  
  - [ ]* 8.5 Write property test for notification preferences
    - **Property 22: Notification Preference Respect**
    - **Validates: Requirements 12.6**
  
  - [x] 8.6 Implement notification history and read status
    - Implement GetUserNotifications method with date descending order
    - Implement MarkAsRead method to update is_read flag
    - _Requirements: 12.4, 12.5_
  
  - [x] 8.7 Create HTTP handlers and routes for Notification Service
    - Implement GET /api/notifications handler
    - Implement PUT /api/notifications/:id/read handler
    - Add authentication middleware
    - _Requirements: 12.4, 12.5, 15.1_
  
  - [ ]* 8.8 Write unit tests for Notification Service
    - Test email sending with mocked provider
    - Test push notification delivery
    - Test preference enforcement


- [x] 9. Analytics Service Implementation
  - [x] 9.1 Create Analytics Service project structure and interfaces
    - Initialize Go module for analytics-service
    - Define Service interface with GetSpendingAnalysis, GetSavingsPatterns, GetRecommendations, GetFinancialHealth methods
    - Create SpendingAnalysis, SavingsPattern, Recommendation, FinancialHealthScore structs
    - _Requirements: 13.1, 13.2, 13.3, 13.4, 14.1_
  
  - [x] 9.2 Implement spending analysis from read replicas
    - Implement GetSpendingAnalysis method querying database replicas
    - Calculate total spending, category breakdown, top merchants, daily average
    - Compare to previous period and calculate percentage change
    - _Requirements: 13.1, 13.2, 13.5_
  
  - [x] 9.3 Implement savings pattern detection
    - Implement GetSavingsPatterns method
    - Determine pattern type (consistent, irregular, improving)
    - Calculate average amount, frequency, best day of week
    - Generate insights based on patterns
    - _Requirements: 13.3_
  
  - [x] 9.4 Implement financial health score calculation
    - Implement GetFinancialHealth method
    - Calculate savings score (40% weight) based on frequency and amount
    - Calculate budget score (30% weight) based on adherence
    - Calculate consistency score (30% weight) based on streak and regularity
    - Compute overall score as weighted average
    - Ensure all scores are integers 0-100
    - Cache results for 1 hour
    - _Requirements: 14.1, 14.2, 14.3, 13.6_
  
  - [ ]* 9.5 Write property tests for financial health score
    - **Property 15: Financial Health Score Bounds**
    - **Validates: Requirements 14.3**
    - **Property 16: Financial Health Score Weighted Average**
    - **Validates: Requirements 14.1, 14.2**
  
  - [x] 9.6 Implement AI-assisted recommendations
    - Implement GetRecommendations method
    - Generate actionable recommendations based on spending and savings patterns
    - Assign priority levels (high, medium, low)
    - Calculate potential savings for each recommendation
    - _Requirements: 13.4_
  
  - [x] 9.7 Create HTTP handlers and routes for Analytics Service
    - Implement GET /api/analytics/spending handler with period parameter
    - Implement GET /api/analytics/patterns handler
    - Implement GET /api/analytics/recommendations handler
    - Implement GET /api/analytics/health handler
    - Add authentication middleware
    - _Requirements: 13.1, 13.3, 13.4, 14.1, 15.1_
  
  - [ ]* 9.8 Write unit tests for Analytics Service
    - Test spending analysis calculations
    - Test pattern detection logic
    - Test financial health score computation


- [ ] 10. Checkpoint - Backend Services Complete
  - Ensure all backend microservices build successfully
  - Verify all database migrations run without errors
  - Run all unit tests and property tests
  - Test inter-service communication if applicable
  - Ask the user if questions arise

- [x] 11. Frontend Application Setup
  - [x] 11.1 Initialize TanStack Start project
    - Create new TanStack Start project with TypeScript
    - Configure Tailwind CSS for styling
    - Set up project structure (routes, components, lib, hooks)
    - Configure environment variables
    - _Requirements: 18.1_
  
  - [x] 11.2 Create API client library
    - Implement API client with typed interfaces for all services
    - Create auth, user, savings, budget, goals, education, analytics, notifications modules
    - Add request/response type definitions
    - Implement error handling and retry logic
    - _Requirements: 15.1, 15.2, 17.1_
  
  - [x] 11.3 Implement authentication context and token management
    - Create AuthContext with login, logout, register, refreshToken methods
    - Implement token storage (httpOnly cookies or secure storage)
    - Add automatic token refresh on expiry
    - Create useAuth hook for components
    - _Requirements: 1.5, 2.3, 15.1, 15.2_
  
  - [x] 11.4 Create protected route wrapper
    - Implement ProtectedRoute component
    - Check authentication status before rendering
    - Redirect to login if unauthenticated
    - _Requirements: 15.1, 15.2, 15.4_
  
  - [ ]* 11.5 Write property test for authentication requirement
    - **Property 20: Authentication Requirement**
    - **Validates: Requirements 15.1, 15.2**


- [x] 12. Frontend Authentication Pages
  - [x] 12.1 Create registration page
    - Implement registration form with email, password, first name, last name, date of birth fields
    - Add client-side validation (email format, password length ≥ 8)
    - Integrate with auth.register API
    - Handle success (redirect to dashboard) and error states
    - _Requirements: 1.1, 1.3, 1.4, 17.1, 17.2_
  
  - [x] 12.2 Create login page
    - Implement login form with email and password fields
    - Integrate with auth.login API
    - Store tokens securely
    - Handle success (redirect to dashboard) and error states
    - _Requirements: 1.6, 1.7_
  
  - [x] 12.3 Implement logout functionality
    - Create logout button/action
    - Call auth.logout API
    - Clear stored tokens
    - Redirect to login page
    - _Requirements: 2.4_
  
  - [ ]* 12.4 Write unit tests for authentication pages
    - Test form validation
    - Test successful registration/login flows
    - Test error handling


- [x] 13. Frontend Dashboard Implementation
  - [x] 13.1 Create dashboard layout and navigation
    - Implement main dashboard layout with sidebar navigation
    - Add navigation links to Savings, Budget, Goals, Education sections
    - Create responsive design for mobile and desktop
    - _Requirements: 18.1_
  
  - [x] 13.2 Implement dashboard summary component
    - Fetch data from savings.getSummary, budget.getCurrentBudget, goals.getActiveGoals, analytics.getFinancialHealth
    - Display total saved, current streak, budget status, active goals count
    - Show financial health score with visual indicator
    - Use TanStack Query for data fetching and caching
    - _Requirements: 4.5, 6.3, 9.3, 14.1_
  
  - [x] 13.3 Create quick stats cards
    - Display savings this month
    - Display budget remaining
    - Display goals progress
    - Display current streak with visual indicator
    - _Requirements: 4.5, 5.2, 6.3, 9.3_
  
  - [x] 13.4 Implement recent activity feed
    - Fetch recent savings and spending transactions
    - Display in chronological order
    - Show transaction type, amount, description, date
    - _Requirements: 4.4, 7.1_
  
  - [ ]* 13.5 Write unit tests for dashboard components
    - Test data fetching and display
    - Test loading and error states
    - Test responsive layout


- [x] 14. Frontend Savings Tracker Implementation
  - [x] 14.1 Create savings tracker page
    - Implement page layout with summary cards and transaction form
    - Display total saved, current streak, longest streak
    - Show monthly savings chart using recharts
    - _Requirements: 4.5, 5.2, 5.5_
  
  - [x] 14.2 Implement savings transaction form
    - Create form with amount, description, category fields
    - Add client-side validation (amount > 0)
    - Integrate with savings.createTransaction API
    - Implement optimistic UI updates
    - Show success message and update streak display
    - _Requirements: 4.1, 4.2, 17.1, 17.2_
  
  - [x] 14.3 Create savings history list
    - Fetch and display savings transactions with pagination
    - Show amount, description, category, date for each transaction
    - Order by date descending
    - _Requirements: 4.4_
  
  - [x] 14.4 Implement streak visualization
    - Create visual streak counter with flame/fire icon
    - Show current streak prominently
    - Display longest streak as achievement
    - Add motivational messages based on streak length
    - _Requirements: 5.2, 5.5_
  
  - [ ]* 14.5 Write unit tests for savings tracker
    - Test form submission and validation
    - Test optimistic updates
    - Test history display


- [x] 15. Frontend Budget Planner Implementation
  - [x] 15.1 Create budget planner page
    - Implement page layout with budget overview and category management
    - Display total budget, total spent, remaining budget
    - Show budget alerts prominently
    - _Requirements: 6.3, 8.1_
  
  - [x] 15.2 Implement budget creation form
    - Create form for setting monthly budget and category allocations
    - Add fields for category name, allocated amount, color
    - Support adding/removing categories dynamically
    - Integrate with budget.createBudget API
    - _Requirements: 6.1, 6.2, 17.1, 17.2_
  
  - [x] 15.3 Create budget category cards
    - Display each category with allocated, spent, remaining amounts
    - Show progress bar with color coding (green < 80%, yellow 80-99%, red ≥ 100%)
    - Display percentage used
    - _Requirements: 6.2, 8.1, 8.2_
  
  - [x] 15.4 Implement spending transaction form
    - Create form with amount, category, description, merchant, date fields
    - Add validation (amount > 0, date not in future)
    - Integrate with budget.recordSpending API
    - Update category spent amounts optimistically
    - _Requirements: 7.1, 7.4, 7.5, 17.1, 17.2_
  
  - [x] 15.5 Create budget alerts display
    - Fetch and display budget alerts
    - Show critical alerts in red, warning alerts in yellow
    - Sort by severity and percentage
    - _Requirements: 8.1, 8.2, 8.4, 8.5_
  
  - [x] 15.6 Implement spending history with category filter
    - Display spending transactions with filtering by category
    - Show amount, description, merchant, date
    - Order by date descending
    - _Requirements: 7.1_
  
  - [ ]* 15.7 Write unit tests for budget planner
    - Test budget creation and updates
    - Test spending recording
    - Test alert display and sorting


- [x] 16. Frontend Goal Manager Implementation
  - [x] 16.1 Create goals page
    - Implement page layout with active goals list and creation form
    - Display goal cards with progress visualization
    - _Requirements: 9.3_
  
  - [x] 16.2 Implement goal creation form
    - Create form with title, description, target amount, target date fields
    - Add validation (target amount > 0, target date in future)
    - Integrate with goals.createGoal API
    - _Requirements: 9.1, 17.1, 17.2_
  
  - [x] 16.3 Create goal card component
    - Display goal title, description, target amount, target date
    - Show current amount and progress percentage
    - Display progress bar with percentage
    - Show status (active, completed, paused)
    - _Requirements: 9.1, 9.6_
  
  - [x] 16.4 Implement goal contribution form
    - Create form to add contribution amount to a goal
    - Integrate with goals.updateProgress API
    - Update goal display optimistically
    - Show celebration animation when goal completed
    - _Requirements: 10.1, 10.2_
  
  - [x] 16.5 Create milestone display
    - Fetch and display milestones for each goal
    - Show milestone title, amount, completion status
    - Mark completed milestones with checkmark
    - Display completion date for completed milestones
    - _Requirements: 10.4, 10.5_
  
  - [x] 16.6 Implement goal edit and delete functionality
    - Add edit button to modify goal details
    - Add delete button with confirmation dialog
    - Integrate with goals.updateGoal and goals.deleteGoal APIs
    - _Requirements: 9.4, 9.5_
  
  - [ ]* 16.7 Write unit tests for goal manager
    - Test goal creation and CRUD operations
    - Test contribution functionality
    - Test milestone display


- [x] 17. Frontend Education Section Implementation
  - [x] 17.1 Create education page
    - Implement page layout with lesson list and progress tracker
    - Display education progress percentage
    - _Requirements: 11.1, 11.4_
  
  - [x] 17.2 Create lesson list component
    - Fetch and display lessons with title, description, duration, difficulty
    - Show completion status (checkmark for completed)
    - Filter by category
    - Order by lesson order
    - _Requirements: 11.1_
  
  - [x] 17.3 Implement lesson detail page
    - Display lesson content (text, video, resources)
    - Show quiz questions if available
    - Add "Mark as Complete" button
    - Integrate with education.markLessonComplete API
    - _Requirements: 11.2, 11.3_
  
  - [x] 17.4 Create progress tracker component
    - Display total lessons, completed lessons, progress percentage
    - Show progress bar
    - Display current streak if applicable
    - _Requirements: 11.4, 11.5_
  
  - [ ]* 17.5 Write unit tests for education section
    - Test lesson list display
    - Test lesson completion
    - Test progress calculation


- [x] 18. Frontend Analytics Dashboard Implementation
  - [x] 18.1 Create analytics page
    - Implement page layout with spending analysis, patterns, and recommendations
    - Display financial health score prominently
    - _Requirements: 13.1, 14.1_
  
  - [x] 18.2 Implement spending analysis visualization
    - Fetch spending analysis data for selected period
    - Create pie chart for category breakdown using recharts
    - Display total spending, daily average, comparison to previous period
    - Show top merchants list
    - _Requirements: 13.1, 13.2_
  
  - [x] 18.3 Create financial health score display
    - Display overall score with visual gauge/meter
    - Show component scores (savings, budget, consistency)
    - Display insights and improvement areas
    - _Requirements: 14.1, 14.2, 14.3_
  
  - [x] 18.4 Implement savings patterns display
    - Fetch and display savings patterns
    - Show pattern type, average amount, frequency
    - Display insights as cards or list
    - _Requirements: 13.3_
  
  - [x] 18.5 Create recommendations list
    - Fetch and display AI-assisted recommendations
    - Show priority level with color coding
    - Display potential savings for each recommendation
    - Show action items as checklist
    - _Requirements: 13.4_
  
  - [ ]* 18.6 Write unit tests for analytics dashboard
    - Test data visualization rendering
    - Test financial health score display
    - Test recommendations display


- [x] 19. Frontend User Profile and Settings
  - [x] 19.1 Create profile page
    - Display user profile information (email, name, date of birth, profile image)
    - Add edit profile form
    - Integrate with user.updateProfile API
    - _Requirements: 3.1, 3.2_
  
  - [x] 19.2 Create settings page
    - Display user preferences (currency, notifications, theme)
    - Add toggle switches for notification settings
    - Add time picker for reminder time
    - Integrate with user.updatePreferences API
    - _Requirements: 3.3_
  
  - [x] 19.3 Implement account deletion
    - Add "Delete Account" button with confirmation dialog
    - Show warning about data loss
    - Integrate with user.deleteAccount API
    - Logout and redirect after successful deletion
    - _Requirements: 3.4_
  
  - [ ]* 19.4 Write unit tests for profile and settings
    - Test profile update
    - Test preferences update
    - Test account deletion flow


- [x] 20. Frontend Notifications Implementation
  - [x] 20.1 Create notifications dropdown/panel
    - Fetch user notifications
    - Display notification list with title, message, timestamp
    - Show unread count badge
    - Order by date descending
    - _Requirements: 12.4_
  
  - [x] 20.2 Implement mark as read functionality
    - Add click handler to mark notification as read
    - Integrate with notifications.markAsRead API
    - Update unread count
    - _Requirements: 12.5_
  
  - [x] 20.3 Add notification bell icon to header
    - Display bell icon with unread count badge
    - Show dropdown on click
    - _Requirements: 12.4_
  
  - [ ]* 20.4 Write unit tests for notifications
    - Test notification display
    - Test mark as read functionality
    - Test unread count updates


- [x] 21. Checkpoint - Frontend Application Complete
  - Ensure all frontend pages render correctly
  - Verify all API integrations work
  - Test authentication flow end-to-end
  - Test responsive design on mobile and desktop
  - Run all frontend unit tests
  - Ask the user if questions arise


- [x] 22. Docker Containerization
  - [x] 22.1 Create Dockerfiles for all Go microservices
    - Write multi-stage Dockerfile for each service (auth, user, savings, budget, goal, education, notification, analytics)
    - Use golang:1.21-alpine as builder, alpine:3.19 as runtime
    - Create non-root user for security
    - Add health check endpoints
    - _Requirements: 18.1, 20.3_
  
  - [x] 22.2 Create Dockerfile for TanStack Start frontend
    - Write multi-stage Dockerfile with Node.js 20
    - Build production bundle
    - Create non-root user
    - Expose port 3000
    - _Requirements: 18.1_
  
  - [x] 22.3 Create docker-compose.yml for local development
    - Define services for all microservices, frontend, PostgreSQL
    - Configure networking and volumes
    - Set environment variables
    - Add health checks
    - _Requirements: 18.1_
  
  - [x] 22.4 Test Docker builds and local deployment
    - Build all Docker images
    - Run docker-compose up
    - Verify all services start and communicate
    - Test basic functionality


- [x] 23. Kubernetes Deployment Configuration
  - [x] 23.1 Create Kubernetes Deployment manifests for all services
    - Write Deployment YAML for each microservice (auth, user, savings, budget, goal, education, notification, analytics)
    - Configure replicas (3 for critical services, 2 for others)
    - Set resource requests and limits (CPU, memory)
    - Add liveness and readiness probes
    - Configure security context (non-root, read-only filesystem)
    - _Requirements: 18.1, 18.2, 19.1, 19.2, 20.3_
  
  - [x] 23.2 Create Kubernetes Service manifests
    - Write Service YAML for each microservice
    - Configure ClusterIP type for internal services
    - Expose HTTP and metrics ports
    - _Requirements: 18.1_
  
  - [x] 23.3 Create HorizontalPodAutoscaler manifests
    - Write HPA YAML for each service
    - Configure CPU and memory-based scaling
    - Set min/max replicas (e.g., 3-20 for savings-service)
    - _Requirements: 18.2_
  
  - [x] 23.4 Create Ingress manifest
    - Write Ingress YAML for API routing
    - Configure path-based routing to services
    - Add TLS configuration with cert-manager annotations
    - Configure rate limiting
    - _Requirements: 18.1, 20.4_
  
  - [x] 23.5 Create PostgreSQL StatefulSet manifest
    - Write StatefulSet YAML for PostgreSQL
    - Configure persistent volume claims
    - Set resource limits
    - Add health checks
    - _Requirements: 18.1, 19.1_
  
  - [x] 23.6 Create ConfigMaps and Secrets
    - Create ConfigMap for non-sensitive configuration
    - Create Secret for database credentials, JWT secret, API keys
    - _Requirements: 20.1, 20.2_


- [x] 24. Observability Stack Setup
  - [x] 24.1 Implement Prometheus metrics in all Go services
    - Add prometheus client library to each service
    - Expose /metrics endpoint on port 9090
    - Instrument HTTP request metrics (count, duration, status)
    - Add custom business metrics (transactions created, active users, etc.)
    - _Requirements: 19.3, 19.4_
  
  - [x] 24.2 Create Prometheus deployment and configuration
    - Write Prometheus Deployment and Service manifests
    - Configure scrape configs for all services
    - Set retention period and storage
    - _Requirements: 19.3_
  
  - [x] 24.3 Create Grafana deployment and dashboards
    - Write Grafana Deployment and Service manifests
    - Create dashboard for service health (uptime, request rate, error rate, latency)
    - Create dashboard for business metrics (registrations, transactions, goals)
    - Create dashboard for infrastructure (pod status, resource usage)
    - _Requirements: 19.3_
  
  - [x] 24.4 Implement structured logging in all services
    - Add structured logging library (e.g., zap or logrus) to each service
    - Log in JSON format with timestamp, level, service, trace_id, message
    - Log at appropriate levels (DEBUG, INFO, WARN, ERROR, FATAL)
    - _Requirements: 19.5_
  
  - [x] 24.5 Implement OpenTelemetry distributed tracing
    - Add OpenTelemetry SDK to each Go service
    - Instrument HTTP handlers with tracing
    - Propagate trace context across services
    - Configure Jaeger exporter
    - _Requirements: 19.5_
  
  - [x] 24.6 Create alerting rules
    - Define Prometheus alerting rules for critical conditions
    - Configure alerts for service down, high error rate, high latency, resource exhaustion
    - Set up alert routing (PagerDuty for critical, Slack for warnings)
    - _Requirements: 19.6_


- [ ] 25. Security Implementation
  - [ ] 25.1 Implement input validation middleware
    - Add validation middleware to all HTTP handlers
    - Validate request body against schemas using go-playground/validator
    - Sanitize inputs to prevent SQL injection and XSS
    - Return 400 Bad Request with detailed errors for invalid input
    - _Requirements: 17.1, 17.2, 17.3, 17.5, 17.6_
  
  - [ ]* 25.2 Write property test for input validation
    - **Property 18: Input Validation Rejection**
    - **Validates: Requirements 4.2, 6.6, 7.4, 7.5, 17.1, 17.2, 17.3, 17.5**
    - **Property 25: Input Sanitization**
    - **Validates: Requirements 17.6**
  
  - [ ] 25.3 Implement authorization middleware
    - Add authorization middleware to check user ownership of resources
    - Verify JWT token user_id matches resource owner
    - Return 403 Forbidden for unauthorized access
    - _Requirements: 3.5, 15.4_
  
  - [ ]* 25.4 Write property test for authorization enforcement
    - **Property 19: Authorization Enforcement**
    - **Validates: Requirements 3.5, 15.4**
  
  - [ ] 25.5 Configure TLS/SSL for all services
    - Set up cert-manager in Kubernetes
    - Configure Let's Encrypt issuer
    - Add TLS certificates to Ingress
    - Enable HSTS headers
    - _Requirements: 20.4, 20.5_
  
  - [ ] 25.6 Add security headers to NGINX/Ingress
    - Configure X-Frame-Options, X-Content-Type-Options, X-XSS-Protection
    - Set Content-Security-Policy
    - Enable Strict-Transport-Security
    - _Requirements: 20.5_
  
  - [ ] 25.7 Implement rate limiting at API gateway
    - Configure rate limiting in NGINX Ingress
    - Set per-user limit (100 req/min) and per-IP limit (1000 req/min)
    - Return 429 Too Many Requests on limit exceeded
    - _Requirements: 20.6_


- [ ] 26. CI/CD Pipeline Implementation
  - [ ] 26.1 Create GitHub Actions workflow for linting
    - Write workflow to run golangci-lint on Go code
    - Run ESLint on TypeScript/React code
    - Trigger on push and pull requests
    - _Requirements: 18.3_
  
  - [ ] 26.2 Create GitHub Actions workflow for testing
    - Write workflow to run Go unit tests with coverage
    - Run frontend tests with Vitest
    - Upload coverage reports to Codecov
    - Set up PostgreSQL service for integration tests
    - _Requirements: 18.3_
  
  - [ ] 26.3 Create GitHub Actions workflow for security scanning
    - Add Trivy vulnerability scanner for filesystem and containers
    - Add Snyk security scan for dependencies
    - Upload results to GitHub Security
    - _Requirements: 20.7_
  
  - [ ] 26.4 Create GitHub Actions workflow for building and pushing Docker images
    - Write workflow to build Docker images for all services
    - Tag images with branch name and commit SHA
    - Push to GitHub Container Registry
    - Use Docker Buildx for multi-platform builds
    - _Requirements: 18.3_
  
  - [ ] 26.5 Create GitHub Actions workflow for deployment to staging
    - Write workflow to deploy to staging environment on develop branch
    - Update Kubernetes deployments with new image tags
    - Wait for rollout completion
    - Run smoke tests
    - _Requirements: 18.3, 18.4_
  
  - [ ] 26.6 Create GitHub Actions workflow for deployment to production
    - Write workflow to deploy to production on main branch
    - Require manual approval
    - Use rolling update strategy
    - Run smoke tests
    - Send Slack notification on success/failure
    - _Requirements: 18.3, 18.4_


- [ ] 27. Integration Testing
  - [ ] 27.1 Set up integration test environment
    - Create Docker Compose setup for integration tests
    - Include all services and PostgreSQL test database
    - Configure test data seeding
    - _Requirements: 18.3_
  
  - [ ] 27.2 Write integration tests for user registration flow
    - Test complete flow: register → create profile → set preferences
    - Verify database state after each step
    - Test Auth Service + User Service integration
    - _Requirements: 1.1, 3.2, 3.3_
  
  - [ ] 27.3 Write integration tests for savings flow
    - Test: create transaction → verify streak update → check notifications
    - Test Savings Service + Notification Service integration
    - _Requirements: 4.1, 4.3, 5.1_
  
  - [ ] 27.4 Write integration tests for budget alert flow
    - Test: record spending → exceed threshold → verify alert → check notification
    - Test Budget Service + Analytics Service + Notification Service integration
    - _Requirements: 7.1, 8.1, 12.1_
  
  - [ ] 27.5 Write integration tests for goal progress flow
    - Test: create goal with milestones → add contributions → verify milestone completion
    - Test Goal Service with database transactions
    - _Requirements: 9.1, 10.1, 10.4_
  
  - [ ]* 27.6 Write integration tests for concurrent goal updates
    - Test concurrent contributions to same goal
    - Verify row-level locking prevents race conditions
    - _Requirements: 10.3, 16.4_


- [ ] 28. Performance Testing and Optimization
  - [ ] 28.1 Create k6 load test scripts
    - Write k6 script for normal load (1,000 users, 10 req/sec)
    - Write k6 script for peak load (10,000 users, 50 req/sec)
    - Write k6 script for stress test (gradually increase to 50,000 users)
    - _Requirements: 18.5_
  
  - [ ] 28.2 Run load tests and analyze results
    - Execute load tests against staging environment
    - Measure p95 and p99 response times
    - Measure error rate and throughput
    - Identify bottlenecks
    - _Requirements: 18.5_
  
  - [ ] 28.3 Optimize database queries
    - Review slow query logs
    - Add missing indexes
    - Optimize N+1 queries
    - Verify partition pruning works correctly
    - _Requirements: 18.5_
  
  - [ ] 28.4 Implement caching strategy
    - Add Redis for session caching (15 min TTL)
    - Cache financial health scores (1 hour TTL)
    - Cache education content (24 hours TTL)
    - Implement cache invalidation on mutations
    - _Requirements: 18.5_
  
  - [ ] 28.5 Configure connection pooling
    - Set max connections per service to 20
    - Configure idle timeout and connection lifetime
    - Set up PgBouncer for connection pooling
    - _Requirements: 18.5_


- [ ] 29. Database Seeding and Sample Data
  - [ ] 29.1 Create database seed scripts
    - Write seed script for sample users
    - Create sample savings transactions with various dates for streak testing
    - Create sample budgets with categories
    - Create sample spending transactions
    - Create sample goals with milestones
    - Create sample education lessons
    - _Requirements: 18.1_
  
  - [ ] 29.2 Create data migration scripts
    - Write scripts for creating monthly partitions
    - Create script for partition maintenance (auto-create future partitions)
    - Write script for archiving old partitions
    - _Requirements: 18.1_


- [ ] 30. Documentation
  - [ ] 30.1 Write API documentation
    - Document all API endpoints with request/response examples
    - Use OpenAPI/Swagger specification
    - Include authentication requirements
    - Document error codes and messages
    - _Requirements: 18.1_
  
  - [ ] 30.2 Write deployment documentation
    - Document Kubernetes cluster setup
    - Document environment variables and secrets
    - Document deployment process
    - Document rollback procedures
    - _Requirements: 18.1, 18.4_
  
  - [ ] 30.3 Write developer setup guide
    - Document local development setup with Docker Compose
    - Document how to run tests
    - Document code structure and conventions
    - Document how to add new services
    - _Requirements: 18.1_
  
  - [ ] 30.4 Create runbook for operations
    - Document monitoring and alerting
    - Document common issues and resolutions
    - Document backup and restore procedures
    - Document scaling procedures
    - _Requirements: 19.6_


- [ ] 31. Final Checkpoint and Production Readiness
  - Verify all services are deployed and healthy
  - Run complete end-to-end test suite
  - Verify all monitoring dashboards are working
  - Verify all alerts are configured
  - Run security scan and address any critical issues
  - Verify backup and restore procedures
  - Verify TLS certificates are valid
  - Run load tests and verify performance targets met
  - Review deployment checklist
  - Ask the user if questions arise

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP delivery
- Each task references specific requirements for traceability
- Property-based tests validate universal correctness properties from the design
- Unit tests validate specific examples and edge cases
- Integration tests validate cross-service workflows
- Checkpoints ensure incremental validation and provide opportunities for user feedback
- Implementation follows dependency order: infrastructure → backend → frontend → deployment
- All code examples use TypeScript for frontend and Go for backend as specified in the design
- Security, observability, and testing are integrated throughout the implementation

