# Design Document: InSavein Platform

## Overview

InSavein is a financial discipline application designed to help young people facing high unemployment in 2026 build savings habits, develop financial discipline, and work toward long-term financial independence. The platform enables users to track micro-savings, categorize spending, plan budgets, maintain financial discipline streaks, and receive AI-assisted savings strategies. Built as a scalable distributed system using TanStack Start (frontend), Golang microservices (backend), PostgreSQL (database), and Kubernetes orchestration, the platform is designed to serve millions of users with high availability and fault tolerance.

## Architecture

### System Architecture Overview

```mermaid
graph TB
    subgraph "Client Layer"
        WEB[Web Browser]
        MOBILE[Mobile Browser]
    end
    
    subgraph "Edge Layer"
        CF[Cloudflare<br/>CDN + WAF + DDoS]
    end
    
    subgraph "Load Balancing Layer"
        NGINX[NGINX Load Balancer<br/>TLS Termination]
    end
    
    subgraph "Kubernetes Cluster"
        subgraph "Frontend Services"
            TSS[TanStack Start Server<br/>SSR + API Routes]
        end
        
        subgraph "Backend Microservices"
            AUTH[Auth Service<br/>Go]
            USER[User Profile Service<br/>Go]
            SAVINGS[Savings Service<br/>Go]
            BUDGET[Budget Service<br/>Go]
            GOAL[Goal Service<br/>Go]
            EDU[Education Service<br/>Go]
            NOTIF[Notification Service<br/>Go]
            ANALYTICS[Analytics Service<br/>Go]
        end
        
        subgraph "Service Mesh"
            INGRESS[Ingress Controller]
            HPA[Horizontal Pod Autoscaler]
        end
    end
    
    subgraph "Data Layer"
        PGPRIMARY[(PostgreSQL Primary)]
        PGREPLICA1[(PostgreSQL Replica 1)]
        PGREPLICA2[(PostgreSQL Replica 2)]
    end
    
    subgraph "Observability"
        PROM[Prometheus]
        GRAF[Grafana]
        OTEL[OpenTelemetry]
    end
    
    WEB --> CF
    MOBILE --> CF
    CF --> NGINX
    NGINX --> INGRESS
    INGRESS --> TSS
    TSS --> AUTH
    TSS --> USER
    TSS --> SAVINGS
    TSS --> BUDGET
    TSS --> GOAL
    TSS --> EDU
    TSS --> NOTIF
    TSS --> ANALYTICS
    
    AUTH --> PGPRIMARY
    USER --> PGPRIMARY
    SAVINGS --> PGPRIMARY
    BUDGET --> PGPRIMARY
    GOAL --> PGPRIMARY
    EDU --> PGREPLICA1
    NOTIF --> PGREPLICA1
    ANALYTICS --> PGREPLICA2
    
    PGPRIMARY -.->|Replication| PGREPLICA1
    PGPRIMARY -.->|Replication| PGREPLICA2
    
    AUTH -.-> PROM
    USER -.-> PROM
    SAVINGS -.-> PROM
    BUDGET -.-> PROM
    GOAL -.-> PROM
    EDU -.-> PROM
    NOTIF -.-> PROM
    ANALYTICS -.-> PROM
    PROM --> GRAF
    
    HPA -.->|Scale| AUTH
    HPA -.->|Scale| USER
    HPA -.->|Scale| SAVINGS
    HPA -.->|Scale| BUDGET
```


### Main User Workflow

```mermaid
sequenceDiagram
    participant U as User
    participant CF as Cloudflare
    participant TSS as TanStack Start
    participant AUTH as Auth Service
    participant SAVINGS as Savings Service
    participant BUDGET as Budget Service
    participant GOAL as Goal Service
    participant DB as PostgreSQL
    
    U->>CF: Access InSavein App
    CF->>TSS: Forward Request (CDN Cache Check)
    TSS->>AUTH: Authenticate User
    AUTH->>DB: Verify JWT Token
    DB-->>AUTH: User Session Valid
    AUTH-->>TSS: Auth Success + User Context
    
    TSS->>SAVINGS: Get Savings Summary
    TSS->>BUDGET: Get Current Budget Status
    TSS->>GOAL: Get Active Goals
    
    par Parallel Data Fetch
        SAVINGS->>DB: Query Savings Data
        BUDGET->>DB: Query Budget Data
        GOAL->>DB: Query Goals Data
    end
    
    DB-->>SAVINGS: Savings Records
    DB-->>BUDGET: Budget Records
    DB-->>GOAL: Goal Records
    
    SAVINGS-->>TSS: Savings Summary
    BUDGET-->>TSS: Budget Status
    GOAL-->>TSS: Goals Progress
    
    TSS->>TSS: Server-Side Render Dashboard
    TSS-->>CF: HTML + Hydration Data
    CF-->>U: Rendered Dashboard
    
    U->>TSS: Record New Savings Entry
    TSS->>SAVINGS: Create Savings Transaction
    SAVINGS->>DB: Insert Transaction
    DB-->>SAVINGS: Transaction Created
    SAVINGS->>SAVINGS: Update Streak Counter
    SAVINGS-->>TSS: Success + Updated Streak
    TSS-->>U: Optimistic UI Update
```

## Components and Interfaces

### Frontend: TanStack Start Application

**Purpose**: Server-first React application providing user interface, server-side rendering, and API route handling

**Technology Stack**:
- TanStack Start (SSR framework)
- React 18+ (UI library)
- TypeScript (type safety)
- TanStack Router (routing)
- TanStack Query (data fetching)
- Tailwind CSS (styling)

**Key Routes**:

```typescript
// app/routes/index.tsx
export default function HomePage() {
  return <LandingPage />
}

// app/routes/dashboard.tsx
export default function DashboardPage() {
  const { data: user } = useUser()
  const { data: savings } = useSavings()
  const { data: budget } = useBudget()
  const { data: goals } = useGoals()
  
  return <Dashboard user={user} savings={savings} budget={budget} goals={goals} />
}

// app/routes/savings.tsx
export default function SavingsPage() {
  return <SavingsTracker />
}

// app/routes/budget.tsx
export default function BudgetPage() {
  return <BudgetPlanner />
}

// app/routes/goals.tsx
export default function GoalsPage() {
  return <GoalManager />
}
```


**API Client Interface**:

```typescript
// lib/api/client.ts
interface ApiClient {
  auth: {
    register(data: RegisterRequest): Promise<AuthResponse>
    login(data: LoginRequest): Promise<AuthResponse>
    logout(): Promise<void>
    refreshToken(): Promise<TokenResponse>
  }
  
  user: {
    getProfile(): Promise<UserProfile>
    updateProfile(data: UpdateProfileRequest): Promise<UserProfile>
    getPreferences(): Promise<UserPreferences>
    updatePreferences(data: UserPreferences): Promise<UserPreferences>
  }
  
  savings: {
    getSummary(): Promise<SavingsSummary>
    getHistory(params: HistoryParams): Promise<SavingsTransaction[]>
    createTransaction(data: CreateSavingsRequest): Promise<SavingsTransaction>
    getStreak(): Promise<SavingsStreak>
  }
  
  budget: {
    getCurrentBudget(): Promise<Budget>
    createBudget(data: CreateBudgetRequest): Promise<Budget>
    updateBudget(id: string, data: UpdateBudgetRequest): Promise<Budget>
    getCategories(): Promise<BudgetCategory[]>
    recordSpending(data: SpendingRequest): Promise<SpendingTransaction>
  }
  
  goals: {
    getActiveGoals(): Promise<Goal[]>
    createGoal(data: CreateGoalRequest): Promise<Goal>
    updateGoal(id: string, data: UpdateGoalRequest): Promise<Goal>
    deleteGoal(id: string): Promise<void>
    getMilestones(goalId: string): Promise<Milestone[]>
  }
  
  education: {
    getLessons(): Promise<Lesson[]>
    getLesson(id: string): Promise<LessonDetail>
    markLessonComplete(id: string): Promise<void>
  }
  
  analytics: {
    getSpendingAnalysis(period: TimePeriod): Promise<SpendingAnalysis>
    getSavingsPatterns(): Promise<SavingsPattern[]>
    getRecommendations(): Promise<Recommendation[]>
  }
}
```

### Backend Microservices

#### 1. Auth Service

**Purpose**: Handle user authentication, registration, JWT token management, and password operations

**Interface**:

```go
// internal/auth/service.go
package auth

type Service interface {
    Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
    Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
    ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
    Logout(ctx context.Context, userID string) error
    ResetPassword(ctx context.Context, req ResetPasswordRequest) error
}

type RegisterRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Password        string `json:"password" validate:"required,min=8"`
    FirstName       string `json:"first_name" validate:"required"`
    LastName        string `json:"last_name" validate:"required"`
    DateOfBirth     string `json:"date_of_birth" validate:"required"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token"`
    ExpiresIn    int64        `json:"expires_in"`
    User         UserSummary  `json:"user"`
}

type TokenClaims struct {
    UserID    string   `json:"user_id"`
    Email     string   `json:"email"`
    Roles     []string `json:"roles"`
    ExpiresAt int64    `json:"exp"`
}
```


#### 2. User Profile Service

**Purpose**: Manage user profile data, settings, and preferences

**Interface**:

```go
// internal/user/service.go
package user

type Service interface {
    GetProfile(ctx context.Context, userID string) (*UserProfile, error)
    UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*UserProfile, error)
    GetPreferences(ctx context.Context, userID string) (*UserPreferences, error)
    UpdatePreferences(ctx context.Context, userID string, prefs UserPreferences) error
    DeleteAccount(ctx context.Context, userID string) error
}

type UserProfile struct {
    ID              string    `json:"id"`
    Email           string    `json:"email"`
    FirstName       string    `json:"first_name"`
    LastName        string    `json:"last_name"`
    DateOfBirth     string    `json:"date_of_birth"`
    ProfileImageURL string    `json:"profile_image_url"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type UserPreferences struct {
    Currency              string   `json:"currency"`
    NotificationsEnabled  bool     `json:"notifications_enabled"`
    EmailNotifications    bool     `json:"email_notifications"`
    PushNotifications     bool     `json:"push_notifications"`
    SavingsReminders      bool     `json:"savings_reminders"`
    ReminderTime          string   `json:"reminder_time"`
    Theme                 string   `json:"theme"`
}
```

#### 3. Savings Service

**Purpose**: Track savings deposits, maintain savings history, calculate streaks, and manage savings goals

**Interface**:

```go
// internal/savings/service.go
package savings

type Service interface {
    GetSummary(ctx context.Context, userID string) (*SavingsSummary, error)
    GetHistory(ctx context.Context, userID string, params HistoryParams) ([]SavingsTransaction, error)
    CreateTransaction(ctx context.Context, userID string, req CreateTransactionRequest) (*SavingsTransaction, error)
    GetStreak(ctx context.Context, userID string) (*SavingsStreak, error)
    UpdateStreak(ctx context.Context, userID string) error
    GetMonthlyStats(ctx context.Context, userID string, month time.Time) (*MonthlyStats, error)
}

type SavingsSummary struct {
    TotalSaved       float64   `json:"total_saved"`
    CurrentStreak    int       `json:"current_streak"`
    LongestStreak    int       `json:"longest_streak"`
    LastSavingDate   time.Time `json:"last_saving_date"`
    MonthlyAverage   float64   `json:"monthly_average"`
    ThisMonthSaved   float64   `json:"this_month_saved"`
}

type SavingsTransaction struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    Amount      float64   `json:"amount"`
    Currency    string    `json:"currency"`
    Description string    `json:"description"`
    Category    string    `json:"category"`
    CreatedAt   time.Time `json:"created_at"`
}

type SavingsStreak struct {
    CurrentStreak int       `json:"current_streak"`
    LongestStreak int       `json:"longest_streak"`
    LastSaveDate  time.Time `json:"last_save_date"`
}
```


#### 4. Budget Service

**Purpose**: Manage monthly budgets, spending categories, track spending, and generate alerts

**Interface**:

```go
// internal/budget/service.go
package budget

type Service interface {
    GetCurrentBudget(ctx context.Context, userID string) (*Budget, error)
    CreateBudget(ctx context.Context, userID string, req CreateBudgetRequest) (*Budget, error)
    UpdateBudget(ctx context.Context, userID string, budgetID string, req UpdateBudgetRequest) (*Budget, error)
    GetCategories(ctx context.Context, userID string) ([]BudgetCategory, error)
    RecordSpending(ctx context.Context, userID string, req SpendingRequest) (*SpendingTransaction, error)
    GetSpendingSummary(ctx context.Context, userID string, month time.Time) (*SpendingSummary, error)
    CheckBudgetAlerts(ctx context.Context, userID string) ([]BudgetAlert, error)
}

type Budget struct {
    ID              string           `json:"id"`
    UserID          string           `json:"user_id"`
    Month           time.Time        `json:"month"`
    TotalBudget     float64          `json:"total_budget"`
    Categories      []BudgetCategory `json:"categories"`
    TotalSpent      float64          `json:"total_spent"`
    RemainingBudget float64          `json:"remaining_budget"`
    CreatedAt       time.Time        `json:"created_at"`
    UpdatedAt       time.Time        `json:"updated_at"`
}

type BudgetCategory struct {
    ID            string    `json:"id"`
    BudgetID      string    `json:"budget_id"`
    Name          string    `json:"name"`
    AllocatedAmount float64 `json:"allocated_amount"`
    SpentAmount   float64   `json:"spent_amount"`
    RemainingAmount float64 `json:"remaining_amount"`
    Color         string    `json:"color"`
}

type SpendingTransaction struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    BudgetID    string    `json:"budget_id"`
    CategoryID  string    `json:"category_id"`
    Amount      float64   `json:"amount"`
    Description string    `json:"description"`
    Merchant    string    `json:"merchant"`
    Date        time.Time `json:"date"`
    CreatedAt   time.Time `json:"created_at"`
}

type BudgetAlert struct {
    CategoryName    string  `json:"category_name"`
    PercentageUsed  float64 `json:"percentage_used"`
    AlertType       string  `json:"alert_type"` // "warning" | "critical"
    Message         string  `json:"message"`
}
```

#### 5. Goal Service

**Purpose**: Manage long-term financial goals, track milestones, and calculate progress

**Interface**:

```go
// internal/goal/service.go
package goal

type Service interface {
    GetActiveGoals(ctx context.Context, userID string) ([]Goal, error)
    GetGoal(ctx context.Context, userID string, goalID string) (*GoalDetail, error)
    CreateGoal(ctx context.Context, userID string, req CreateGoalRequest) (*Goal, error)
    UpdateGoal(ctx context.Context, userID string, goalID string, req UpdateGoalRequest) (*Goal, error)
    DeleteGoal(ctx context.Context, userID string, goalID string) error
    GetMilestones(ctx context.Context, goalID string) ([]Milestone, error)
    UpdateProgress(ctx context.Context, goalID string, amount float64) error
}

type Goal struct {
    ID              string    `json:"id"`
    UserID          string    `json:"user_id"`
    Title           string    `json:"title"`
    Description     string    `json:"description"`
    TargetAmount    float64   `json:"target_amount"`
    CurrentAmount   float64   `json:"current_amount"`
    Currency        string    `json:"currency"`
    TargetDate      time.Time `json:"target_date"`
    Status          string    `json:"status"` // "active" | "completed" | "paused"
    ProgressPercent float64   `json:"progress_percent"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type Milestone struct {
    ID          string    `json:"id"`
    GoalID      string    `json:"goal_id"`
    Title       string    `json:"title"`
    Amount      float64   `json:"amount"`
    IsCompleted bool      `json:"is_completed"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
    Order       int       `json:"order"`
}
```


#### 6. Education Service

**Purpose**: Deliver financial education content, track lesson completion, and provide motivational content

**Interface**:

```go
// internal/education/service.go
package education

type Service interface {
    GetLessons(ctx context.Context, userID string) ([]Lesson, error)
    GetLesson(ctx context.Context, userID string, lessonID string) (*LessonDetail, error)
    MarkLessonComplete(ctx context.Context, userID string, lessonID string) error
    GetUserProgress(ctx context.Context, userID string) (*EducationProgress, error)
    GetRecommendedLessons(ctx context.Context, userID string) ([]Lesson, error)
}

type Lesson struct {
    ID              string   `json:"id"`
    Title           string   `json:"title"`
    Description     string   `json:"description"`
    Category        string   `json:"category"`
    DurationMinutes int      `json:"duration_minutes"`
    Difficulty      string   `json:"difficulty"` // "beginner" | "intermediate" | "advanced"
    Tags            []string `json:"tags"`
    IsCompleted     bool     `json:"is_completed"`
    Order           int      `json:"order"`
}

type LessonDetail struct {
    Lesson
    Content         string    `json:"content"`
    VideoURL        string    `json:"video_url,omitempty"`
    Resources       []Resource `json:"resources"`
    Quiz            []QuizQuestion `json:"quiz,omitempty"`
}

type EducationProgress struct {
    TotalLessons      int     `json:"total_lessons"`
    CompletedLessons  int     `json:"completed_lessons"`
    ProgressPercent   float64 `json:"progress_percent"`
    CurrentStreak     int     `json:"current_streak"`
}
```

#### 7. Notification Service

**Purpose**: Send email notifications, push notifications, and reminders to users

**Interface**:

```go
// internal/notification/service.go
package notification

type Service interface {
    SendEmail(ctx context.Context, req EmailRequest) error
    SendPushNotification(ctx context.Context, req PushNotificationRequest) error
    ScheduleReminder(ctx context.Context, req ReminderRequest) error
    GetUserNotifications(ctx context.Context, userID string) ([]Notification, error)
    MarkAsRead(ctx context.Context, userID string, notificationID string) error
}

type EmailRequest struct {
    To          string            `json:"to"`
    Subject     string            `json:"subject"`
    TemplateID  string            `json:"template_id"`
    TemplateData map[string]interface{} `json:"template_data"`
}

type PushNotificationRequest struct {
    UserID  string `json:"user_id"`
    Title   string `json:"title"`
    Body    string `json:"body"`
    Data    map[string]string `json:"data"`
}

type ReminderRequest struct {
    UserID      string    `json:"user_id"`
    Type        string    `json:"type"` // "savings" | "budget" | "goal"
    ScheduledAt time.Time `json:"scheduled_at"`
    Message     string    `json:"message"`
}

type Notification struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Type      string    `json:"type"`
    Title     string    `json:"title"`
    Message   string    `json:"message"`
    IsRead    bool      `json:"is_read"`
    CreatedAt time.Time `json:"created_at"`
}
```


#### 8. Analytics Service

**Purpose**: Analyze spending patterns, generate savings insights, and provide AI-assisted recommendations

**Interface**:

```go
// internal/analytics/service.go
package analytics

type Service interface {
    GetSpendingAnalysis(ctx context.Context, userID string, period TimePeriod) (*SpendingAnalysis, error)
    GetSavingsPatterns(ctx context.Context, userID string) ([]SavingsPattern, error)
    GetRecommendations(ctx context.Context, userID string) ([]Recommendation, error)
    GenerateMonthlyReport(ctx context.Context, userID string, month time.Time) (*MonthlyReport, error)
    GetFinancialHealth(ctx context.Context, userID string) (*FinancialHealthScore, error)
}

type SpendingAnalysis struct {
    Period              TimePeriod           `json:"period"`
    TotalSpending       float64              `json:"total_spending"`
    CategoryBreakdown   []CategorySpending   `json:"category_breakdown"`
    TopMerchants        []MerchantSpending   `json:"top_merchants"`
    DailyAverage        float64              `json:"daily_average"`
    ComparisonToPrevious float64             `json:"comparison_to_previous"`
    Trends              []SpendingTrend      `json:"trends"`
}

type SavingsPattern struct {
    PatternType     string    `json:"pattern_type"` // "consistent" | "irregular" | "improving"
    AverageAmount   float64   `json:"average_amount"`
    Frequency       string    `json:"frequency"`
    BestDayOfWeek   string    `json:"best_day_of_week"`
    Insights        []string  `json:"insights"`
}

type Recommendation struct {
    ID          string   `json:"id"`
    Type        string   `json:"type"` // "savings" | "budget" | "spending"
    Priority    string   `json:"priority"` // "high" | "medium" | "low"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    ActionItems []string `json:"action_items"`
    PotentialSavings float64 `json:"potential_savings,omitempty"`
}

type FinancialHealthScore struct {
    OverallScore      int     `json:"overall_score"` // 0-100
    SavingsScore      int     `json:"savings_score"`
    BudgetScore       int     `json:"budget_score"`
    ConsistencyScore  int     `json:"consistency_score"`
    Insights          []string `json:"insights"`
    ImprovementAreas  []string `json:"improvement_areas"`
}
```

## Data Models

### Database Schema

```mermaid
erDiagram
    USERS ||--o{ SAVINGS_TRANSACTIONS : has
    USERS ||--o{ BUDGETS : creates
    USERS ||--o{ GOALS : sets
    USERS ||--o{ NOTIFICATIONS : receives
    USERS ||--o{ EDUCATION_PROGRESS : tracks
    
    BUDGETS ||--o{ BUDGET_CATEGORIES : contains
    BUDGET_CATEGORIES ||--o{ SPENDING_TRANSACTIONS : tracks
    
    GOALS ||--o{ GOAL_MILESTONES : has
    
    LESSONS ||--o{ EDUCATION_PROGRESS : tracked_by
    
    USERS {
        uuid id PK
        string email UK
        string password_hash
        string first_name
        string last_name
        date date_of_birth
        string profile_image_url
        jsonb preferences
        timestamp created_at
        timestamp updated_at
    }
    
    SAVINGS_TRANSACTIONS {
        uuid id PK
        uuid user_id FK
        decimal amount
        string currency
        string description
        string category
        timestamp created_at
    }
    
    BUDGETS {
        uuid id PK
        uuid user_id FK
        date month
        decimal total_budget
        decimal total_spent
        timestamp created_at
        timestamp updated_at
    }
    
    BUDGET_CATEGORIES {
        uuid id PK
        uuid budget_id FK
        string name
        decimal allocated_amount
        decimal spent_amount
        string color
    }
    
    SPENDING_TRANSACTIONS {
        uuid id PK
        uuid user_id FK
        uuid budget_id FK
        uuid category_id FK
        decimal amount
        string description
        string merchant
        date transaction_date
        timestamp created_at
    }
    
    GOALS {
        uuid id PK
        uuid user_id FK
        string title
        text description
        decimal target_amount
        decimal current_amount
        string currency
        date target_date
        string status
        timestamp created_at
        timestamp updated_at
    }
    
    GOAL_MILESTONES {
        uuid id PK
        uuid goal_id FK
        string title
        decimal amount
        boolean is_completed
        timestamp completed_at
        int order
    }
    
    NOTIFICATIONS {
        uuid id PK
        uuid user_id FK
        string type
        string title
        text message
        boolean is_read
        timestamp created_at
    }
    
    LESSONS {
        uuid id PK
        string title
        text description
        string category
        int duration_minutes
        string difficulty
        text content
        string video_url
        jsonb resources
        int order
    }
    
    EDUCATION_PROGRESS {
        uuid id PK
        uuid user_id FK
        uuid lesson_id FK
        boolean is_completed
        timestamp completed_at
    }
```


### PostgreSQL Schema Definition

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    profile_image_url TEXT,
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Savings transactions table (partitioned by month)
CREATE TABLE savings_transactions (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    description TEXT,
    category VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_savings_user_date ON savings_transactions(user_id, created_at DESC);
CREATE INDEX idx_savings_category ON savings_transactions(category);

-- Budgets table
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    month DATE NOT NULL,
    total_budget DECIMAL(15, 2) NOT NULL CHECK (total_budget >= 0),
    total_spent DECIMAL(15, 2) DEFAULT 0 CHECK (total_spent >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, month)
);

CREATE INDEX idx_budgets_user_month ON budgets(user_id, month DESC);

-- Budget categories table
CREATE TABLE budget_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    allocated_amount DECIMAL(15, 2) NOT NULL CHECK (allocated_amount >= 0),
    spent_amount DECIMAL(15, 2) DEFAULT 0 CHECK (spent_amount >= 0),
    color VARCHAR(7) DEFAULT '#000000'
);

CREATE INDEX idx_budget_categories_budget ON budget_categories(budget_id);

-- Spending transactions table (partitioned by date)
CREATE TABLE spending_transactions (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    category_id UUID REFERENCES budget_categories(id) ON DELETE SET NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    description TEXT,
    merchant VARCHAR(255),
    transaction_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (id, transaction_date)
) PARTITION BY RANGE (transaction_date);

CREATE INDEX idx_spending_user_date ON spending_transactions(user_id, transaction_date DESC);
CREATE INDEX idx_spending_category ON spending_transactions(category_id);

-- Goals table
CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    target_amount DECIMAL(15, 2) NOT NULL CHECK (target_amount > 0),
    current_amount DECIMAL(15, 2) DEFAULT 0 CHECK (current_amount >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    target_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'paused')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_goals_user_status ON goals(user_id, status);
CREATE INDEX idx_goals_target_date ON goals(target_date);

-- Goal milestones table
CREATE TABLE goal_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    goal_id UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    "order" INT NOT NULL
);

CREATE INDEX idx_milestones_goal ON goal_milestones(goal_id, "order");

-- Notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_read ON notifications(user_id, is_read, created_at DESC);

-- Lessons table
CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    duration_minutes INT NOT NULL,
    difficulty VARCHAR(20) CHECK (difficulty IN ('beginner', 'intermediate', 'advanced')),
    content TEXT NOT NULL,
    video_url TEXT,
    resources JSONB DEFAULT '[]',
    "order" INT NOT NULL
);

CREATE INDEX idx_lessons_category ON lessons(category, "order");

-- Education progress table
CREATE TABLE education_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, lesson_id)
);

CREATE INDEX idx_education_progress_user ON education_progress(user_id, is_completed);
```

## Algorithmic Pseudocode

### Main Processing Algorithm: User Dashboard Data Aggregation

```typescript
async function aggregateDashboardData(userId: string): Promise<DashboardData> {
  // INPUT: userId of type string (UUID)
  // OUTPUT: DashboardData containing all dashboard information
  // PRECONDITIONS:
  //   - userId is a valid UUID
  //   - User exists in the database
  //   - User is authenticated
  // POSTCONDITIONS:
  //   - Returns complete dashboard data
  //   - All data is current and consistent
  //   - No side effects on database state
  
  // Step 1: Validate user authentication
  const tokenClaims = await validateToken(context.token)
  if (tokenClaims.userId !== userId) {
    throw new UnauthorizedError("Token user mismatch")
  }
  
  // Step 2: Fetch data in parallel from all services
  const [
    savingsSummary,
    currentBudget,
    activeGoals,
    recentNotifications,
    financialHealth
  ] = await Promise.all([
    savingsService.getSummary(userId),
    budgetService.getCurrentBudget(userId),
    goalService.getActiveGoals(userId),
    notificationService.getRecent(userId, 10),
    analyticsService.getFinancialHealth(userId)
  ])
  
  // Step 3: Calculate derived metrics
  const dashboardData: DashboardData = {
    user: tokenClaims,
    savings: savingsSummary,
    budget: currentBudget,
    goals: activeGoals,
    notifications: recentNotifications,
    financialHealth: financialHealth,
    quickStats: calculateQuickStats(savingsSummary, currentBudget, activeGoals)
  }
  
  // POSTCONDITION CHECK: Ensure all required fields are present
  if (!dashboardData.savings || !dashboardData.budget) {
    throw new Error("Incomplete dashboard data")
  }
  
  return dashboardData
}
```

**Preconditions:**
- `userId` is a valid UUID format
- User exists in the database
- User authentication token is valid
- All microservices are available and responsive

**Postconditions:**
- Returns `DashboardData` object with all required fields populated
- Data is consistent across all services
- No mutations to database state
- Response time < 500ms for optimal UX

**Loop Invariants:** N/A (parallel execution, no loops)


### Savings Streak Calculation Algorithm

```typescript
async function updateSavingsStreak(userId: string): Promise<SavingsStreak> {
  // INPUT: userId of type string (UUID)
  // OUTPUT: SavingsStreak with current and longest streak counts
  // PRECONDITIONS:
  //   - userId is valid and user exists
  //   - Savings transactions table is accessible
  // POSTCONDITIONS:
  //   - Returns accurate streak counts
  //   - Current streak reflects consecutive days with savings
  //   - Longest streak is historical maximum
  //   - Database updated with new streak values
  
  // Step 1: Get all savings transactions ordered by date
  const transactions = await db.query(`
    SELECT DATE(created_at) as save_date
    FROM savings_transactions
    WHERE user_id = $1
    ORDER BY created_at DESC
  `, [userId])
  
  if (transactions.length === 0) {
    return { currentStreak: 0, longestStreak: 0, lastSaveDate: null }
  }
  
  // Step 2: Calculate current streak
  let currentStreak = 0
  let longestStreak = 0
  let tempStreak = 1
  const today = new Date()
  const lastSaveDate = transactions[0].save_date
  
  // Check if last save was today or yesterday
  const daysSinceLastSave = daysBetween(lastSaveDate, today)
  if (daysSinceLastSave > 1) {
    currentStreak = 0
  } else {
    currentStreak = 1
    
    // Step 3: Count consecutive days backwards
    // LOOP INVARIANT: tempStreak represents consecutive days up to current index
    for (let i = 1; i < transactions.length; i++) {
      const prevDate = transactions[i - 1].save_date
      const currDate = transactions[i].save_date
      const dayDiff = daysBetween(currDate, prevDate)
      
      if (dayDiff === 1) {
        tempStreak++
        currentStreak = tempStreak
      } else if (dayDiff > 1) {
        // Streak broken
        longestStreak = Math.max(longestStreak, tempStreak)
        tempStreak = 1
      }
      // dayDiff === 0 means multiple saves same day, continue streak
    }
  }
  
  // Step 4: Final longest streak calculation
  longestStreak = Math.max(longestStreak, tempStreak, currentStreak)
  
  // Step 5: Update database
  await db.query(`
    UPDATE users
    SET preferences = jsonb_set(
      jsonb_set(preferences, '{current_streak}', $1::text::jsonb),
      '{longest_streak}', $2::text::jsonb
    )
    WHERE id = $3
  `, [currentStreak, longestStreak, userId])
  
  return {
    currentStreak,
    longestStreak,
    lastSaveDate
  }
}
```

**Preconditions:**
- `userId` is a valid UUID
- User exists in database
- `savings_transactions` table is accessible and indexed
- Transactions are ordered by `created_at` DESC

**Postconditions:**
- Returns `SavingsStreak` with accurate counts
- `currentStreak` = 0 if no saves in last 2 days
- `currentStreak` ≥ 1 if saved today or yesterday
- `longestStreak` ≥ `currentStreak`
- Database `users.preferences` updated with new streak values

**Loop Invariants:**
- `tempStreak` represents the count of consecutive days up to index `i`
- `longestStreak` holds the maximum streak found so far
- All dates are processed in descending chronological order
- `dayDiff` is always calculated between consecutive transaction dates


### Budget Alert Detection Algorithm

```typescript
async function checkBudgetAlerts(userId: string): Promise<BudgetAlert[]> {
  // INPUT: userId of type string (UUID)
  // OUTPUT: Array of BudgetAlert objects for categories exceeding thresholds
  // PRECONDITIONS:
  //   - userId is valid
  //   - Current month budget exists
  //   - Budget categories are defined
  // POSTCONDITIONS:
  //   - Returns alerts for categories at 80% (warning) or 100% (critical)
  //   - Alerts are sorted by severity (critical first)
  //   - No duplicate alerts
  
  const alerts: BudgetAlert[] = []
  
  // Step 1: Get current month's budget
  const budget = await budgetService.getCurrentBudget(userId)
  if (!budget) {
    return alerts
  }
  
  // Step 2: Check each category against thresholds
  // LOOP INVARIANT: All processed categories have been checked for alerts
  for (const category of budget.categories) {
    const percentageUsed = (category.spentAmount / category.allocatedAmount) * 100
    
    // Critical alert: 100% or more spent
    if (percentageUsed >= 100) {
      alerts.push({
        categoryName: category.name,
        percentageUsed: percentageUsed,
        alertType: 'critical',
        message: `You've exceeded your ${category.name} budget by ${(percentageUsed - 100).toFixed(1)}%`
      })
    }
    // Warning alert: 80-99% spent
    else if (percentageUsed >= 80) {
      alerts.push({
        categoryName: category.name,
        percentageUsed: percentageUsed,
        alertType: 'warning',
        message: `You've used ${percentageUsed.toFixed(1)}% of your ${category.name} budget`
      })
    }
  }
  
  // Step 3: Sort alerts by severity (critical first, then by percentage)
  alerts.sort((a, b) => {
    if (a.alertType === 'critical' && b.alertType === 'warning') return -1
    if (a.alertType === 'warning' && b.alertType === 'critical') return 1
    return b.percentageUsed - a.percentageUsed
  })
  
  return alerts
}
```

**Preconditions:**
- `userId` is a valid UUID
- User has a budget for the current month
- Budget categories have `allocatedAmount > 0`
- `spentAmount` is accurately tracked

**Postconditions:**
- Returns array of `BudgetAlert` objects
- Alerts only for categories ≥ 80% spent
- Critical alerts (≥100%) appear before warning alerts (80-99%)
- Within same severity, sorted by percentage descending
- No alerts for categories < 80% spent

**Loop Invariants:**
- All categories up to index `i` have been evaluated for alerts
- `alerts` array contains only valid alert objects
- No duplicate alerts for the same category
- `percentageUsed` is correctly calculated for each category


### Goal Progress Calculation Algorithm

```typescript
async function updateGoalProgress(goalId: string, contributionAmount: float64): Promise<Goal> {
  // INPUT: goalId (UUID), contributionAmount (positive decimal)
  // OUTPUT: Updated Goal object with new progress
  // PRECONDITIONS:
  //   - goalId is valid and goal exists
  //   - contributionAmount > 0
  //   - Goal status is 'active'
  // POSTCONDITIONS:
  //   - Goal currentAmount increased by contributionAmount
  //   - progressPercent accurately reflects new amount
  //   - If target reached, status changed to 'completed'
  //   - Milestones updated if thresholds crossed
  //   - Database transaction committed atomically
  
  // Step 1: Start database transaction
  const tx = await db.beginTransaction()
  
  try {
    // Step 2: Get goal with row lock
    const goal = await tx.query(`
      SELECT * FROM goals WHERE id = $1 AND status = 'active' FOR UPDATE
    `, [goalId])
    
    if (!goal) {
      throw new Error("Goal not found or not active")
    }
    
    // Step 3: Calculate new amounts
    const newCurrentAmount = goal.currentAmount + contributionAmount
    const newProgressPercent = (newCurrentAmount / goal.targetAmount) * 100
    const isCompleted = newCurrentAmount >= goal.targetAmount
    
    // Step 4: Update goal
    await tx.query(`
      UPDATE goals
      SET current_amount = $1,
          status = $2,
          updated_at = NOW()
      WHERE id = $3
    `, [
      newCurrentAmount,
      isCompleted ? 'completed' : 'active',
      goalId
    ])
    
    // Step 5: Update milestones
    const milestones = await tx.query(`
      SELECT * FROM goal_milestones
      WHERE goal_id = $1 AND is_completed = false
      ORDER BY amount ASC
    `, [goalId])
    
    // LOOP INVARIANT: All milestones up to newCurrentAmount are marked completed
    for (const milestone of milestones) {
      if (newCurrentAmount >= milestone.amount) {
        await tx.query(`
          UPDATE goal_milestones
          SET is_completed = true, completed_at = NOW()
          WHERE id = $1
        `, [milestone.id])
      } else {
        // Milestones are ordered, so we can break early
        break
      }
    }
    
    // Step 6: Commit transaction
    await tx.commit()
    
    return {
      ...goal,
      currentAmount: newCurrentAmount,
      progressPercent: newProgressPercent,
      status: isCompleted ? 'completed' : 'active'
    }
    
  } catch (error) {
    await tx.rollback()
    throw error
  }
}
```

**Preconditions:**
- `goalId` is a valid UUID
- Goal exists and has `status = 'active'`
- `contributionAmount > 0`
- `goal.targetAmount > 0`
- Database supports transactions with row-level locking

**Postconditions:**
- `goal.currentAmount` increased by exactly `contributionAmount`
- `goal.progressPercent = (currentAmount / targetAmount) * 100`
- If `currentAmount >= targetAmount`, then `status = 'completed'`
- All milestones with `amount <= currentAmount` are marked `is_completed = true`
- Database changes are atomic (all or nothing)
- `goal.updated_at` reflects current timestamp

**Loop Invariants:**
- Milestones are processed in ascending order by amount
- All milestones with `amount <= newCurrentAmount` are marked completed
- Once a milestone with `amount > newCurrentAmount` is found, loop terminates
- No milestone is processed twice


### JWT Token Validation Algorithm

```go
func ValidateToken(tokenString string, secretKey []byte) (*TokenClaims, error) {
    // INPUT: tokenString (JWT string), secretKey (signing key)
    // OUTPUT: TokenClaims if valid, error otherwise
    // PRECONDITIONS:
    //   - tokenString is non-empty
    //   - secretKey is the correct signing key
    // POSTCONDITIONS:
    //   - Returns TokenClaims if signature valid and not expired
    //   - Returns error if invalid signature, expired, or malformed
    //   - No side effects
    
    // Step 1: Parse token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })
    
    if err != nil {
        return nil, fmt.Errorf("token parse error: %w", err)
    }
    
    // Step 2: Extract claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token claims")
    }
    
    // Step 3: Validate expiration
    exp, ok := claims["exp"].(float64)
    if !ok {
        return nil, errors.New("missing expiration claim")
    }
    
    if time.Now().Unix() > int64(exp) {
        return nil, errors.New("token expired")
    }
    
    // Step 4: Extract required fields
    userID, ok := claims["user_id"].(string)
    if !ok || userID == "" {
        return nil, errors.New("missing or invalid user_id")
    }
    
    email, ok := claims["email"].(string)
    if !ok || email == "" {
        return nil, errors.New("missing or invalid email")
    }
    
    // Step 5: Extract roles (optional field)
    var roles []string
    if rolesInterface, ok := claims["roles"].([]interface{}); ok {
        for _, role := range rolesInterface {
            if roleStr, ok := role.(string); ok {
                roles = append(roles, roleStr)
            }
        }
    }
    
    return &TokenClaims{
        UserID:    userID,
        Email:     email,
        Roles:     roles,
        ExpiresAt: int64(exp),
    }, nil
}
```

**Preconditions:**
- `tokenString` is a non-empty string
- `secretKey` is the correct HMAC signing key used to create the token
- Token was created using HS256, HS384, or HS512 algorithm

**Postconditions:**
- Returns `*TokenClaims` with valid user data if token is valid
- Returns `error` if:
  - Token signature is invalid
  - Token is expired (`exp` < current time)
  - Token is malformed or missing required claims
  - Signing method is not HMAC
- No database queries or side effects
- Function is idempotent (same input always produces same output)

**Loop Invariants:**
- When extracting roles: all processed interface{} elements are converted to strings
- Roles array contains only valid string values


## Key Functions with Formal Specifications

### Function: CreateSavingsTransaction

```go
func (s *SavingsService) CreateTransaction(ctx context.Context, userID string, req CreateTransactionRequest) (*SavingsTransaction, error)
```

**Preconditions:**
- `userID` is a valid UUID and user exists in database
- `req.Amount > 0`
- `req.Currency` is a valid 3-letter ISO currency code
- Database connection is available

**Postconditions:**
- New savings transaction record created in database
- Transaction ID is a valid UUID
- `created_at` timestamp is set to current time
- User's savings streak is updated asynchronously
- Returns `*SavingsTransaction` with all fields populated
- If error occurs, no database changes are committed

**Side Effects:**
- Inserts row into `savings_transactions` table
- Triggers async streak calculation
- May trigger notification if milestone reached

### Function: RecordSpending

```go
func (s *BudgetService) RecordSpending(ctx context.Context, userID string, req SpendingRequest) (*SpendingTransaction, error)
```

**Preconditions:**
- `userID` is valid and user exists
- `req.Amount > 0`
- `req.CategoryID` references existing budget category
- Budget for current month exists
- `req.TransactionDate` is not in the future

**Postconditions:**
- New spending transaction created
- Budget category `spent_amount` increased by `req.Amount`
- Budget `total_spent` increased by `req.Amount`
- If category exceeds 80% of allocated amount, alert is generated
- Returns `*SpendingTransaction` with complete data
- All updates are atomic (transaction-based)

**Side Effects:**
- Inserts row into `spending_transactions` table
- Updates `budget_categories.spent_amount`
- Updates `budgets.total_spent`
- May create notification if budget threshold exceeded

### Function: AuthenticateUser

```go
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
```

**Preconditions:**
- `req.Email` is non-empty and valid email format
- `req.Password` is non-empty
- Database connection is available

**Postconditions:**
- If credentials valid:
  - Returns `*AuthResponse` with access token and refresh token
  - Access token expires in 15 minutes
  - Refresh token expires in 7 days
  - Tokens are signed with HMAC-SHA256
- If credentials invalid:
  - Returns error "Invalid credentials"
  - No tokens generated
- Password is never logged or exposed
- Failed login attempts are rate-limited

**Side Effects:**
- Queries `users` table for email lookup
- Compares bcrypt password hash
- No database writes on successful login
- May increment failed login counter (rate limiting)

### Function: CalculateFinancialHealth

```go
func (s *AnalyticsService) GetFinancialHealth(ctx context.Context, userID string) (*FinancialHealthScore, error)
```

**Preconditions:**
- `userID` is valid and user exists
- User has at least 30 days of transaction history
- Database replicas are available for read queries

**Postconditions:**
- Returns `*FinancialHealthScore` with scores 0-100
- `OverallScore` is weighted average of component scores:
  - `SavingsScore` (40% weight): Based on savings frequency and amount
  - `BudgetScore` (30% weight): Based on budget adherence
  - `ConsistencyScore` (30% weight): Based on streak and regularity
- Scores are deterministic for same input data
- Includes actionable insights and improvement areas
- Read-only operation, no database mutations

**Side Effects:**
- Queries multiple tables (savings, budgets, spending) from read replicas
- No database writes
- Results may be cached for 1 hour

## Example Usage

### Frontend: Creating a Savings Entry

```typescript
// components/SavingsForm.tsx
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/lib/api/client'

export function SavingsForm() {
  const queryClient = useQueryClient()
  
  const savingsMutation = useMutation({
    mutationFn: (data: CreateSavingsRequest) => 
      apiClient.savings.createTransaction(data),
    onSuccess: (newTransaction) => {
      // Optimistic update
      queryClient.setQueryData(['savings', 'summary'], (old: SavingsSummary) => ({
        ...old,
        totalSaved: old.totalSaved + newTransaction.amount,
        thisMonthSaved: old.thisMonthSaved + newTransaction.amount
      }))
      
      // Refetch to get updated streak
      queryClient.invalidateQueries(['savings', 'streak'])
    }
  })
  
  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    const formData = new FormData(e.target as HTMLFormElement)
    
    savingsMutation.mutate({
      amount: parseFloat(formData.get('amount') as string),
      currency: 'USD',
      description: formData.get('description') as string,
      category: formData.get('category') as string
    })
  }
  
  return (
    <form onSubmit={handleSubmit}>
      <input name="amount" type="number" step="0.01" required />
      <input name="description" type="text" />
      <select name="category">
        <option value="emergency">Emergency Fund</option>
        <option value="goal">Goal Savings</option>
        <option value="general">General Savings</option>
      </select>
      <button type="submit" disabled={savingsMutation.isPending}>
        {savingsMutation.isPending ? 'Saving...' : 'Save Money'}
      </button>
    </form>
  )
}
```

### Backend: Savings Service Implementation

```go
// internal/savings/service.go
package savings

import (
    "context"
    "database/sql"
    "time"
    
    "github.com/google/uuid"
)

type service struct {
    db     *sql.DB
    logger Logger
}

func (s *service) CreateTransaction(
    ctx context.Context,
    userID string,
    req CreateTransactionRequest,
) (*SavingsTransaction, error) {
    // Validate input
    if req.Amount <= 0 {
        return nil, ErrInvalidAmount
    }
    
    // Create transaction
    tx := &SavingsTransaction{
        ID:          uuid.New().String(),
        UserID:      userID,
        Amount:      req.Amount,
        Currency:    req.Currency,
        Description: req.Description,
        Category:    req.Category,
        CreatedAt:   time.Now(),
    }
    
    // Insert into database
    query := `
        INSERT INTO savings_transactions 
        (id, user_id, amount, currency, description, category, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    
    _, err := s.db.ExecContext(ctx, query,
        tx.ID, tx.UserID, tx.Amount, tx.Currency,
        tx.Description, tx.Category, tx.CreatedAt,
    )
    
    if err != nil {
        s.logger.Error("failed to insert savings transaction", "error", err)
        return nil, err
    }
    
    // Trigger async streak update
    go s.updateStreakAsync(userID)
    
    s.logger.Info("savings transaction created",
        "user_id", userID,
        "amount", req.Amount,
        "transaction_id", tx.ID,
    )
    
    return tx, nil
}

func (s *service) updateStreakAsync(userID string) {
    ctx := context.Background()
    if err := s.UpdateStreak(ctx, userID); err != nil {
        s.logger.Error("failed to update streak", "user_id", userID, "error", err)
    }
}
```


### Complete Workflow: User Registration to First Savings

```typescript
// Example: Complete user journey from registration to first savings entry

// Step 1: User Registration
const registerResponse = await apiClient.auth.register({
  email: "user@example.com",
  password: "SecurePass123!",
  firstName: "Jane",
  lastName: "Doe",
  dateOfBirth: "1998-05-15"
})

// Response includes access token and user data
const { accessToken, user } = registerResponse

// Step 2: Set user preferences
await apiClient.user.updatePreferences({
  currency: "USD",
  notificationsEnabled: true,
  savingsReminders: true,
  reminderTime: "09:00"
})

// Step 3: Create first budget
const budget = await apiClient.budget.createBudget({
  month: "2026-01",
  totalBudget: 500,
  categories: [
    { name: "Food", allocatedAmount: 200, color: "#FF6B6B" },
    { name: "Transport", allocatedAmount: 100, color: "#4ECDC4" },
    { name: "Entertainment", allocatedAmount: 50, color: "#45B7D1" },
    { name: "Other", allocatedAmount: 150, color: "#96CEB4" }
  ]
})

// Step 4: Set a financial goal
const goal = await apiClient.goals.createGoal({
  title: "Emergency Fund",
  description: "Build 3 months of expenses",
  targetAmount: 1500,
  currency: "USD",
  targetDate: "2026-12-31"
})

// Step 5: Make first savings entry
const savingsTransaction = await apiClient.savings.createTransaction({
  amount: 10,
  currency: "USD",
  description: "First savings - starting my journey!",
  category: "emergency"
})

// Step 6: Check updated streak
const streak = await apiClient.savings.getStreak()
// { currentStreak: 1, longestStreak: 1, lastSaveDate: "2026-01-15" }

// Step 7: View dashboard summary
const dashboard = await apiClient.dashboard.getSummary()
// Returns complete dashboard with savings, budget, goals, and health score
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Password Security

*For any* user registration or password update, the password must be hashed using bcrypt with cost factor at least 12, and the plaintext password must never be stored or returned in API responses.

**Validates: Requirements 1.2, 20.1, 20.2**

### Property 2: JWT Token Validity

*For any* valid JWT token, the signature must verify correctly with the secret key, the expiration time must be in the future, and the token must contain user_id, email, and roles in the payload.

**Validates: Requirements 2.1, 2.2, 2.5, 2.6, 15.5**

### Property 3: Token Refresh Round Trip

*For any* valid refresh token, using it to obtain new tokens must succeed and produce valid access and refresh tokens with correct expiration times (15 minutes for access, 7 days for refresh).

**Validates: Requirements 1.5, 2.3**

### Property 4: Savings Transaction Integrity

*For any* savings transaction, the amount must be positive, the user_id must reference an existing user, the created_at timestamp must not be in the future, and all amounts must be stored with exactly 2 decimal places precision.

**Validates: Requirements 4.1, 4.2, 4.6**

### Property 5: Savings Streak Invariant

*For any* user, the current streak must be non-negative, must never exceed the longest streak, and must be zero if the last savings date is more than 1 day ago.

**Validates: Requirements 5.2, 5.5, 5.6**

### Property 6: Streak Calculation Determinism

*For any* sequence of savings transaction dates, calculating the streak multiple times must always produce the same current streak and longest streak values.

**Validates: Requirements 5.1, 5.3, 5.4**

### Property 7: Budget Consistency Invariant

*For any* budget, the total_spent must equal the sum of all category spent_amounts, all amounts must be non-negative, and the budget must be unique per user per month.

**Validates: Requirements 6.5, 6.6, 7.2, 7.3, 16.6**

### Property 8: Spending Transaction Atomicity

*For any* spending transaction, creating it must atomically update both the spending_transactions table and the budget_categories spent_amount, or roll back all changes if any step fails.

**Validates: Requirements 7.6, 16.2, 16.3, 16.5**

### Property 9: Budget Alert Thresholds

*For any* budget category, a warning alert must be generated when spent reaches 80-99% of allocated, a critical alert when spent reaches 100% or more, and no alerts for categories with zero allocated amount.

**Validates: Requirements 8.1, 8.2, 8.3, 8.6**

### Property 10: Alert Sorting Order

*For any* collection of budget alerts, they must be sorted with critical alerts before warning alerts, and within the same type, sorted by percentage descending.

**Validates: Requirements 8.4, 8.5**

### Property 11: Goal Progress Calculation

*For any* goal, the progress percentage must equal (current_amount / target_amount) × 100, must be between 0 and 100, and the status must be "completed" if and only if current_amount >= target_amount.

**Validates: Requirements 9.6, 10.2**

### Property 12: Goal Contribution Monotonicity

*For any* goal and positive contribution amount, adding the contribution must increase current_amount by exactly that amount and never decrease it.

**Validates: Requirements 10.1**

### Property 13: Milestone Completion Order

*For any* goal with milestones, if a milestone at amount M2 is completed, then all milestones at amounts M1 < M2 must also be completed.

**Validates: Requirements 10.4, 10.5, 10.6**

### Property 14: Goal Update Concurrency Safety

*For any* two concurrent goal update operations on the same goal, the database must use row-level locking to ensure only one succeeds at a time, preventing race conditions.

**Validates: Requirements 10.3, 16.4**

### Property 15: Financial Health Score Bounds

*For any* financial health score calculation, all component scores (savings, budget, consistency) and the overall score must be integers between 0 and 100 inclusive.

**Validates: Requirements 14.3**

### Property 16: Financial Health Score Weighted Average

*For any* financial health score, the overall score must equal the weighted average of savings score (40%), budget score (30%), and consistency score (30%).

**Validates: Requirements 14.1, 14.2**

### Property 17: Education Progress Calculation

*For any* user's education progress, the progress percentage must equal (completed_lessons / total_lessons) × 100 and must be between 0 and 100.

**Validates: Requirements 11.4, 11.5**

### Property 18: Input Validation Rejection

*For any* API request with invalid input (negative amounts, missing required fields, out-of-range values, future dates for past transactions), the system must return HTTP 400 Bad Request with detailed validation errors and make no database changes.

**Validates: Requirements 4.2, 6.6, 7.4, 7.5, 17.1, 17.2, 17.3, 17.5**

### Property 19: Authorization Enforcement

*For any* authenticated request, users must only be able to access and modify their own data, and attempts to access another user's data must return HTTP 403 Forbidden.

**Validates: Requirements 3.5, 15.4**

### Property 20: Authentication Requirement

*For any* API request without a valid token or with an expired token, the system must return HTTP 401 Unauthorized and not process the request.

**Validates: Requirements 15.1, 15.2**

### Property 21: Cascade Deletion Completeness

*For any* user account deletion or goal deletion, all associated data (transactions, budgets, goals, milestones) must be removed atomically.

**Validates: Requirements 3.4, 9.5**

### Property 22: Notification Preference Respect

*For any* user with notifications disabled in preferences, the system must not send notifications to that user regardless of triggering events.

**Validates: Requirements 12.6**

### Property 23: Data Ordering Consistency

*For any* request for historical data (savings history, notifications, spending transactions), the results must be ordered by creation date descending.

**Validates: Requirements 4.4, 12.4**

### Property 24: Health Check Response Correctness

*For any* health check request, the system must return HTTP 200 OK if healthy (database connected), HTTP 503 Service Unavailable if unhealthy (database disconnected), and respond within 1 second.

**Validates: Requirements 19.1, 19.2, 19.6**

### Property 25: Input Sanitization

*For any* user input containing potentially malicious content (SQL injection attempts, XSS payloads), the system must sanitize the input and prevent the attack from executing.

**Validates: Requirements 17.6**

### Property 26: Transaction Rollback on Failure

*For any* multi-table database operation that fails partway through, all changes made within that transaction must be rolled back, leaving the database in its original state.

**Validates: Requirements 16.1, 16.2**

### Property 27: Profile Update Round Trip

*For any* user profile or preferences update, reading the profile immediately after the update must return the updated values.

**Validates: Requirements 3.2, 3.3**

### Property 28: Duplicate Email Rejection

*For any* registration attempt with an email that already exists in the system, the registration must fail with an error indicating the email is already in use.

**Validates: Requirements 1.3**

### Property 29: Password Complexity Enforcement

*For any* registration or password change with a password shorter than 8 characters, the system must reject it with a validation error.

**Validates: Requirements 1.4**

### Property 30: Streak Update Trigger

*For any* savings transaction creation, the system must trigger an asynchronous streak update for that user.

**Validates: Requirements 4.3**

## Error Handling

### Error Scenario 1: Invalid Authentication Token

**Condition**: User provides expired or invalid JWT token
**Response**: 
- HTTP 401 Unauthorized
- Error message: "Invalid or expired token"
- Client redirected to login page
**Recovery**: 
- User must re-authenticate
- Refresh token can be used if still valid
- Session state cleared from client

### Error Scenario 2: Budget Category Overspending

**Condition**: User attempts to record spending that exceeds category budget
**Response**:
- Transaction is still recorded (not blocked)
- HTTP 200 OK with warning flag
- Budget alert generated
- Notification sent to user
**Recovery**:
- User can adjust budget allocation
- User can review spending and make corrections
- Analytics service provides recommendations

### Error Scenario 3: Database Connection Failure

**Condition**: PostgreSQL primary database becomes unavailable
**Response**:
- Read operations automatically routed to replicas
- Write operations return HTTP 503 Service Unavailable
- Circuit breaker activated after 3 consecutive failures
- Error logged to monitoring system
**Recovery**:
- Kubernetes health checks detect failure
- Automatic failover to replica promoted to primary
- Circuit breaker reopens after 30 seconds
- Retry logic with exponential backoff

### Error Scenario 4: Microservice Timeout

**Condition**: Microservice doesn't respond within 5 seconds
**Response**:
- Request times out
- HTTP 504 Gateway Timeout returned
- Partial data returned if available (graceful degradation)
- Error logged with trace ID
**Recovery**:
- Client retries with exponential backoff
- Load balancer routes subsequent requests to healthy instances
- HPA scales up service if consistent timeouts detected

### Error Scenario 5: Invalid Input Data

**Condition**: User submits form with invalid data (negative amount, invalid date, etc.)
**Response**:
- HTTP 400 Bad Request
- Detailed validation errors returned
- No database changes made
**Recovery**:
- Client displays validation errors to user
- User corrects input and resubmits
- Optimistic UI updates are rolled back

### Error Scenario 6: Concurrent Goal Update Conflict

**Condition**: Two simultaneous requests attempt to update same goal
**Response**:
- Database row-level locking prevents race condition
- First transaction commits successfully
- Second transaction detects conflict and retries
**Recovery**:
- Automatic retry with fresh data
- If retry fails, return HTTP 409 Conflict
- Client refetches current state and allows user to retry


## Testing Strategy

### Unit Testing Approach

**Frontend (TanStack Start + React)**
- Test framework: Vitest + React Testing Library
- Coverage target: 80% for business logic, 60% for UI components
- Key test areas:
  - API client functions (mocked responses)
  - Form validation logic
  - State management hooks
  - Utility functions (date calculations, currency formatting)
  - Component rendering and user interactions

**Backend (Golang Microservices)**
- Test framework: Go standard testing package + testify
- Coverage target: 85% for service layer, 70% for handlers
- Key test areas:
  - Service method logic with mocked database
  - Input validation
  - Error handling paths
  - JWT token generation and validation
  - Business logic (streak calculation, budget alerts, goal progress)

**Example Unit Test (Go)**:
```go
func TestCreateSavingsTransaction(t *testing.T) {
    // Setup
    mockDB := &MockDatabase{}
    service := NewSavingsService(mockDB, logger)
    
    // Test case: Valid transaction
    t.Run("valid transaction creates record", func(t *testing.T) {
        req := CreateTransactionRequest{
            Amount:      50.00,
            Currency:    "USD",
            Description: "Weekly savings",
            Category:    "emergency",
        }
        
        result, err := service.CreateTransaction(context.Background(), "user-123", req)
        
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Equal(t, 50.00, result.Amount)
        assert.Equal(t, "user-123", result.UserID)
    })
    
    // Test case: Invalid amount
    t.Run("negative amount returns error", func(t *testing.T) {
        req := CreateTransactionRequest{
            Amount:   -10.00,
            Currency: "USD",
        }
        
        result, err := service.CreateTransaction(context.Background(), "user-123", req)
        
        assert.Error(t, err)
        assert.Nil(t, result)
        assert.Equal(t, ErrInvalidAmount, err)
    })
}
```

### Property-Based Testing Approach

**Property Test Library**: 
- Frontend: fast-check (TypeScript)
- Backend: gopter (Go)

**Key Properties to Test**:

1. **Savings Streak Calculation**
   - Property: For any sequence of savings dates, currentStreak ≤ longestStreak
   - Property: If last save date is > 1 day ago, currentStreak = 0
   - Property: Streak calculation is deterministic (same input → same output)

2. **Budget Calculations**
   - Property: Sum of category spent amounts always equals budget total spent
   - Property: Budget remaining = total budget - total spent
   - Property: Category percentage = (spent / allocated) × 100

3. **Goal Progress**
   - Property: Progress percentage is always between 0 and 100
   - Property: Adding positive contribution always increases current amount
   - Property: Goal completion status correctly reflects target achievement

4. **JWT Token Operations**
   - Property: Token generation followed by validation always succeeds (if not expired)
   - Property: Expired tokens always fail validation
   - Property: Tampered tokens always fail validation

**Example Property Test (TypeScript with fast-check)**:
```typescript
import fc from 'fast-check'

describe('Budget calculations', () => {
  it('category totals always equal budget total spent', () => {
    fc.assert(
      fc.property(
        fc.array(fc.record({
          name: fc.string(),
          allocated: fc.float({ min: 0, max: 1000 }),
          spent: fc.float({ min: 0, max: 1000 })
        })),
        (categories) => {
          const budget = {
            categories,
            totalSpent: categories.reduce((sum, cat) => sum + cat.spent, 0)
          }
          
          const calculatedTotal = calculateBudgetTotal(budget)
          
          // Property: calculated total matches sum of categories
          expect(calculatedTotal).toBeCloseTo(budget.totalSpent, 2)
        }
      )
    )
  })
})
```

**Example Property Test (Go with gopter)**:
```go
func TestStreakCalculationProperties(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("current streak never exceeds longest streak", 
        prop.ForAll(
            func(dates []time.Time) bool {
                streak := calculateStreak(dates)
                return streak.CurrentStreak <= streak.LongestStreak
            },
            gen.SliceOf(gen.TimeRange(
                time.Now().AddDate(0, -6, 0),
                time.Now(),
            )),
        ),
    )
    
    properties.Property("no saves in 2+ days means zero current streak",
        prop.ForAll(
            func(lastSaveDate time.Time) bool {
                if time.Since(lastSaveDate) > 48*time.Hour {
                    streak := calculateStreak([]time.Time{lastSaveDate})
                    return streak.CurrentStreak == 0
                }
                return true
            },
            gen.TimeRange(
                time.Now().AddDate(0, 0, -30),
                time.Now().AddDate(0, 0, -3),
            ),
        ),
    )
    
    properties.TestingRun(t)
}
```

### Integration Testing Approach

**Test Environment**: 
- Docker Compose with all services
- Test PostgreSQL database
- Mock external services (email, push notifications)

**Key Integration Tests**:

1. **End-to-End User Registration Flow**
   - Register user → Create profile → Set preferences → Verify database state
   - Tests: Auth Service + User Service + Database

2. **Savings Entry with Streak Update**
   - Create savings transaction → Verify streak calculation → Check notifications
   - Tests: Savings Service + Notification Service + Database

3. **Budget Alert Generation**
   - Record spending → Exceed threshold → Verify alert created → Check notification sent
   - Tests: Budget Service + Analytics Service + Notification Service

4. **Goal Progress with Milestones**
   - Create goal with milestones → Add contributions → Verify milestone completion
   - Tests: Goal Service + Database transactions

**Example Integration Test**:
```go
func TestSavingsFlowIntegration(t *testing.T) {
    // Setup test environment
    db := setupTestDatabase(t)
    defer db.Close()
    
    authService := auth.NewService(db)
    savingsService := savings.NewService(db)
    notificationService := notification.NewService(db)
    
    // Step 1: Create user
    user, err := authService.Register(context.Background(), auth.RegisterRequest{
        Email:     "test@example.com",
        Password:  "password123",
        FirstName: "Test",
        LastName:  "User",
    })
    require.NoError(t, err)
    
    // Step 2: Create savings transaction
    tx, err := savingsService.CreateTransaction(context.Background(), user.ID, savings.CreateTransactionRequest{
        Amount:   25.00,
        Currency: "USD",
        Category: "emergency",
    })
    require.NoError(t, err)
    assert.Equal(t, 25.00, tx.Amount)
    
    // Step 3: Verify streak updated
    streak, err := savingsService.GetStreak(context.Background(), user.ID)
    require.NoError(t, err)
    assert.Equal(t, 1, streak.CurrentStreak)
    
    // Step 4: Verify notification created (if first save)
    notifications, err := notificationService.GetUserNotifications(context.Background(), user.ID)
    require.NoError(t, err)
    assert.GreaterOrEqual(t, len(notifications), 1)
}
```

### Load Testing Strategy

**Tool**: k6 (Grafana k6)

**Test Scenarios**:
1. **Normal Load**: 1,000 concurrent users, 10 req/sec per user
2. **Peak Load**: 10,000 concurrent users, 50 req/sec per user
3. **Stress Test**: Gradually increase to 50,000 users until system degrades

**Performance Targets**:
- API response time p95 < 500ms
- API response time p99 < 1000ms
- Error rate < 0.1%
- Database query time p95 < 100ms
- Throughput: 100,000 requests/minute

## Performance Considerations

### Caching Strategy

1. **Frontend Caching**
   - TanStack Query cache: 5 minutes for dashboard data
   - Browser cache: Static assets (1 year), API responses (no-cache)
   - Service Worker: Offline support for read-only views

2. **CDN Caching (Cloudflare)**
   - Static assets: Cache for 1 year with versioned URLs
   - API responses: No caching (dynamic data)
   - HTML pages: Cache for 5 minutes with stale-while-revalidate

3. **Application-Level Caching**
   - Redis cache for frequently accessed data:
     - User sessions: 15 minutes TTL
     - Financial health scores: 1 hour TTL
     - Education content: 24 hours TTL
   - Cache invalidation on data mutations

### Database Optimization

1. **Indexing Strategy**
   - Primary keys: UUID with B-tree index
   - Foreign keys: All indexed
   - Query-specific indexes:
     - `(user_id, created_at DESC)` for transaction history
     - `(user_id, status)` for active goals
     - `(user_id, month)` for budget lookups

2. **Partitioning**
   - `savings_transactions`: Monthly partitions
   - `spending_transactions`: Monthly partitions
   - Automatic partition creation via cron job
   - Old partitions archived after 2 years

3. **Connection Pooling**
   - Max connections per service: 20
   - Idle timeout: 5 minutes
   - Connection lifetime: 30 minutes
   - PgBouncer for connection pooling

4. **Read Replicas**
   - 2 read replicas for read-heavy operations
   - Analytics queries routed to replicas
   - Education content served from replicas
   - Replication lag monitoring (target < 1 second)

### Horizontal Scaling

1. **Kubernetes HPA Configuration**
   ```yaml
   apiVersion: autoscaling/v2
   kind: HorizontalPodAutoscaler
   metadata:
     name: savings-service-hpa
   spec:
     scaleTargetRef:
       apiVersion: apps/v1
       kind: Deployment
       name: savings-service
     minReplicas: 3
     maxReplicas: 20
     metrics:
     - type: Resource
       resource:
         name: cpu
         target:
           type: Utilization
           averageUtilization: 70
     - type: Resource
       resource:
         name: memory
         target:
           type: Utilization
           averageUtilization: 80
   ```

2. **Service Replica Counts**
   - Auth Service: 3-10 replicas
   - Savings Service: 3-20 replicas
   - Budget Service: 3-15 replicas
   - Analytics Service: 2-10 replicas
   - Other services: 2-5 replicas

### API Rate Limiting

- Per user: 100 requests/minute
- Per IP: 1000 requests/minute
- Burst allowance: 20 requests
- Rate limit headers included in responses
- 429 Too Many Requests on limit exceeded


## Security Considerations

### Authentication & Authorization

1. **JWT Token Security**
   - Algorithm: HMAC-SHA256 (HS256)
   - Access token lifetime: 15 minutes
   - Refresh token lifetime: 7 days
   - Token rotation on refresh
   - Secure token storage: httpOnly cookies or secure storage
   - Token revocation list for logout

2. **Password Security**
   - Hashing: bcrypt with cost factor 12
   - Minimum length: 8 characters
   - Complexity requirements: uppercase, lowercase, number
   - Password reset via time-limited tokens (1 hour expiry)
   - Rate limiting on login attempts: 5 attempts per 15 minutes

3. **Role-Based Access Control (RBAC)**
   - Roles: `user`, `admin`, `support`
   - Permissions checked at API gateway and service level
   - User can only access their own data
   - Admin can access all data (audit logged)

### API Security

1. **Input Validation**
   - All inputs validated against schemas
   - SQL injection prevention via parameterized queries
   - XSS prevention via input sanitization
   - CSRF protection via tokens
   - Request size limits: 1MB for JSON, 10MB for file uploads

2. **TLS/SSL**
   - TLS 1.3 required for all connections
   - Certificate management via Let's Encrypt
   - HSTS headers enabled (max-age: 31536000)
   - Certificate pinning for mobile apps

3. **API Rate Limiting**
   - Implemented at NGINX and application level
   - Per-user and per-IP limits
   - Exponential backoff for repeated violations
   - DDoS protection via Cloudflare

### Data Security

1. **Encryption at Rest**
   - PostgreSQL: Transparent Data Encryption (TDE)
   - Sensitive fields: Additional application-level encryption
   - Encryption keys stored in Kubernetes Secrets
   - Key rotation every 90 days

2. **Encryption in Transit**
   - All service-to-service communication over TLS
   - Database connections encrypted
   - Internal Kubernetes network policies

3. **Data Privacy**
   - PII (Personally Identifiable Information) identified and protected
   - GDPR compliance: Right to access, right to deletion
   - Data retention policy: 7 years for financial data
   - Audit logs for all data access

### Security Headers

```nginx
# NGINX security headers
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';" always;
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

### Vulnerability Management

1. **Dependency Scanning**
   - Automated scanning via Snyk or Dependabot
   - Weekly vulnerability reports
   - Critical vulnerabilities patched within 24 hours

2. **Container Security**
   - Base images: Minimal Alpine or Distroless
   - Image scanning via Trivy
   - No root user in containers
   - Read-only root filesystem where possible

3. **Security Audits**
   - Quarterly penetration testing
   - Annual third-party security audit
   - Bug bounty program for responsible disclosure

## Infrastructure

### Kubernetes Architecture

```yaml
# Example Deployment: Savings Service
apiVersion: apps/v1
kind: Deployment
metadata:
  name: savings-service
  namespace: insavein
spec:
  replicas: 3
  selector:
    matchLabels:
      app: savings-service
  template:
    metadata:
      labels:
        app: savings-service
        version: v1
    spec:
      containers:
      - name: savings-service
        image: insavein/savings-service:v1.2.3
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: postgres-credentials
              key: connection-string
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-secret
              key: secret
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          readOnlyRootFilesystem: true
          allowPrivilegeEscalation: false
---
apiVersion: v1
kind: Service
metadata:
  name: savings-service
  namespace: insavein
spec:
  selector:
    app: savings-service
  ports:
  - name: http
    port: 80
    targetPort: 8080
  - name: metrics
    port: 9090
    targetPort: 9090
  type: ClusterIP
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: savings-service-hpa
  namespace: insavein
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: savings-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Ingress Configuration

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: insavein-ingress
  namespace: insavein
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - api.insavein.com
    secretName: insavein-tls
  rules:
  - host: api.insavein.com
    http:
      paths:
      - path: /api/auth
        pathType: Prefix
        backend:
          service:
            name: auth-service
            port:
              number: 80
      - path: /api/savings
        pathType: Prefix
        backend:
          service:
            name: savings-service
            port:
              number: 80
      - path: /api/budget
        pathType: Prefix
        backend:
          service:
            name: budget-service
            port:
              number: 80
      - path: /api/goals
        pathType: Prefix
        backend:
          service:
            name: goal-service
            port:
              number: 80
```

### Docker Configuration

**Multi-Stage Dockerfile for Go Services**:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run binary
CMD ["./main"]
```

**Dockerfile for TanStack Start Frontend**:

```dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./
RUN npm ci

# Copy source code
COPY . .

# Build application
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy built application
COPY --from=builder /app/.output ./

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 3000

ENV NODE_ENV=production

CMD ["node", "server/index.mjs"]
```

### PostgreSQL High Availability

**Primary-Replica Configuration**:

```yaml
# PostgreSQL StatefulSet
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: insavein
spec:
  serviceName: postgres
  replicas: 3
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:16-alpine
        ports:
        - containerPort: 5432
          name: postgres
        env:
        - name: POSTGRES_DB
          value: insavein
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: postgres-credentials
              key: username
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-credentials
              key: password
        - name: PGDATA
          value: /var/lib/postgresql/data/pgdata
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "2Gi"
            cpu: "1000m"
          limits:
            memory: "4Gi"
            cpu: "2000m"
  volumeClaimTemplates:
  - metadata:
      name: postgres-storage
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 100Gi
```

**Replication Configuration** (postgresql.conf):
```conf
# Replication settings
wal_level = replica
max_wal_senders = 10
max_replication_slots = 10
hot_standby = on
hot_standby_feedback = on

# Performance tuning
shared_buffers = 2GB
effective_cache_size = 6GB
maintenance_work_mem = 512MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 10MB
min_wal_size = 1GB
max_wal_size = 4GB
max_connections = 200
```


## CI/CD Pipeline

### GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml
name: Build and Deploy

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  lint:
    name: Lint Code
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          args: --timeout=5m
      
      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
      
      - name: Install frontend dependencies
        working-directory: ./frontend
        run: npm ci
      
      - name: Run ESLint
        working-directory: ./frontend
        run: npm run lint

  test:
    name: Run Tests
    runs-on: ubuntu-latest
    needs: lint
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: insavein_test
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Run Go tests
        env:
          DATABASE_URL: postgres://test:test@localhost:5432/insavein_test?sslmode=disable
        run: |
          go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
      
      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
      
      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
      
      - name: Install frontend dependencies
        working-directory: ./frontend
        run: npm ci
      
      - name: Run frontend tests
        working-directory: ./frontend
        run: npm run test:ci

  security-scan:
    name: Security Scanning
    runs-on: ubuntu-latest
    needs: lint
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'
      
      - name: Upload Trivy results to GitHub Security
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: 'trivy-results.sarif'
      
      - name: Run Snyk security scan
        uses: snyk/actions/golang@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}

  build-and-push:
    name: Build and Push Docker Images
    runs-on: ubuntu-latest
    needs: [test, security-scan]
    if: github.event_name == 'push'
    strategy:
      matrix:
        service:
          - auth-service
          - user-service
          - savings-service
          - budget-service
          - goal-service
          - education-service
          - notification-service
          - analytics-service
          - frontend
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      
      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/${{ matrix.service }}
          tags: |
            type=ref,event=branch
            type=sha,prefix={{branch}}-
            type=semver,pattern={{version}}
      
      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: ./${{ matrix.service }}
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          platforms: linux/amd64,linux/arm64

  deploy-staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    needs: build-and-push
    if: github.ref == 'refs/heads/develop'
    environment:
      name: staging
      url: https://staging.insavein.com
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up kubectl
        uses: azure/setup-kubectl@v3
      
      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBE_CONFIG_STAGING }}" | base64 -d > kubeconfig
          export KUBECONFIG=./kubeconfig
      
      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/auth-service \
            auth-service=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/auth-service:develop-${{ github.sha }} \
            -n insavein-staging
          kubectl set image deployment/savings-service \
            savings-service=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/savings-service:develop-${{ github.sha }} \
            -n insavein-staging
          # ... repeat for all services
      
      - name: Wait for rollout
        run: |
          kubectl rollout status deployment/auth-service -n insavein-staging
          kubectl rollout status deployment/savings-service -n insavein-staging
          # ... repeat for all services
      
      - name: Run smoke tests
        run: |
          npm ci
          npm run test:smoke -- --env=staging

  deploy-production:
    name: Deploy to Production
    runs-on: ubuntu-latest
    needs: build-and-push
    if: github.ref == 'refs/heads/main'
    environment:
      name: production
      url: https://insavein.com
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up kubectl
        uses: azure/setup-kubectl@v3
      
      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBE_CONFIG_PRODUCTION }}" | base64 -d > kubeconfig
          export KUBECONFIG=./kubeconfig
      
      - name: Deploy to Kubernetes (Rolling Update)
        run: |
          kubectl set image deployment/auth-service \
            auth-service=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/auth-service:main-${{ github.sha }} \
            -n insavein-production
          kubectl set image deployment/savings-service \
            savings-service=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/savings-service:main-${{ github.sha }} \
            -n insavein-production
          # ... repeat for all services
      
      - name: Wait for rollout
        run: |
          kubectl rollout status deployment/auth-service -n insavein-production --timeout=10m
          kubectl rollout status deployment/savings-service -n insavein-production --timeout=10m
          # ... repeat for all services
      
      - name: Run smoke tests
        run: |
          npm ci
          npm run test:smoke -- --env=production
      
      - name: Notify deployment
        if: success()
        uses: slackapi/slack-github-action@v1
        with:
          webhook-url: ${{ secrets.SLACK_WEBHOOK }}
          payload: |
            {
              "text": "✅ InSavein deployed to production successfully!",
              "blocks": [
                {
                  "type": "section",
                  "text": {
                    "type": "mrkdwn",
                    "text": "*Deployment Successful*\n\nCommit: ${{ github.sha }}\nBranch: ${{ github.ref }}\nAuthor: ${{ github.actor }}"
                  }
                }
              ]
            }
```

### Deployment Strategy

**Rolling Update Configuration**:
- Max surge: 25% (1 extra pod per 4 pods)
- Max unavailable: 0 (zero downtime)
- Update strategy: Rolling update with health checks
- Rollback: Automatic on failed health checks

**Blue-Green Deployment** (for major releases):
1. Deploy new version to "green" environment
2. Run comprehensive tests on green
3. Switch traffic from blue to green via Ingress
4. Monitor for issues
5. Keep blue environment for quick rollback

**Canary Deployment** (for risky changes):
1. Deploy new version to 10% of pods
2. Monitor metrics for 30 minutes
3. If stable, increase to 50%
4. Monitor for 30 minutes
5. If stable, complete rollout to 100%

## Observability

### Prometheus Metrics

**Service-Level Metrics**:
```go
// Example metrics in Go service
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    savingsTransactionsTotal = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "savings_transactions_total",
            Help: "Total number of savings transactions created",
        },
    )
    
    activeUsersGauge = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users",
            Help: "Number of currently active users",
        },
    )
)
```

**Key Metrics to Monitor**:
- Request rate (requests/second)
- Error rate (errors/second, percentage)
- Request duration (p50, p95, p99)
- Database query duration
- Active connections
- Memory usage
- CPU usage
- Disk I/O
- Network I/O

### Grafana Dashboards

**Dashboard 1: Service Health**
- Service uptime
- Request rate per service
- Error rate per service
- Response time percentiles
- Active connections

**Dashboard 2: Business Metrics**
- New user registrations (per hour/day)
- Savings transactions created
- Total amount saved (aggregated)
- Active goals count
- Budget alerts generated
- Financial health score distribution

**Dashboard 3: Infrastructure**
- Kubernetes pod status
- Node resource usage
- Database connections
- Replication lag
- Disk usage
- Network throughput

### Logging Strategy

**Structured Logging Format** (JSON):
```json
{
  "timestamp": "2026-01-15T10:30:45Z",
  "level": "info",
  "service": "savings-service",
  "trace_id": "abc123def456",
  "user_id": "user-789",
  "message": "Savings transaction created",
  "amount": 25.00,
  "currency": "USD",
  "transaction_id": "tx-xyz789"
}
```

**Log Levels**:
- DEBUG: Detailed diagnostic information
- INFO: General informational messages
- WARN: Warning messages (non-critical issues)
- ERROR: Error messages (failures that don't crash service)
- FATAL: Critical errors (service crashes)

**Log Aggregation**:
- Centralized logging via Loki or ELK stack
- Log retention: 30 days for INFO, 90 days for ERROR
- Log search and filtering capabilities
- Alert on ERROR/FATAL logs

### Distributed Tracing

**OpenTelemetry Integration**:
```go
// Example tracing in Go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func (s *SavingsService) CreateTransaction(ctx context.Context, userID string, req CreateTransactionRequest) (*SavingsTransaction, error) {
    tracer := otel.Tracer("savings-service")
    ctx, span := tracer.Start(ctx, "CreateTransaction")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("user_id", userID),
        attribute.Float64("amount", req.Amount),
    )
    
    // Business logic...
    
    return tx, nil
}
```

**Trace Visualization**:
- Jaeger UI for trace exploration
- Trace retention: 7 days
- Sampling rate: 10% in production, 100% in staging
- Trace critical paths: User registration, savings creation, goal updates

## Dependencies

### External Services

1. **Email Service**: SendGrid or AWS SES
   - Transactional emails (registration, password reset)
   - Notification emails (budget alerts, goal milestones)
   - Rate limit: 100 emails/second

2. **Push Notification Service**: Firebase Cloud Messaging (FCM)
   - Mobile push notifications
   - Web push notifications
   - Rate limit: 1000 notifications/second

3. **Object Storage**: AWS S3 or MinIO
   - User profile images
   - Education content (videos, PDFs)
   - Backup storage

4. **CDN**: Cloudflare
   - Static asset delivery
   - DDoS protection
   - WAF (Web Application Firewall)
   - Bot protection

### Third-Party Libraries

**Frontend (TypeScript/React)**:
- `@tanstack/react-router` - Routing
- `@tanstack/react-query` - Data fetching and caching
- `react` - UI library
- `tailwindcss` - Styling
- `zod` - Schema validation
- `date-fns` - Date manipulation
- `recharts` - Data visualization
- `react-hook-form` - Form handling

**Backend (Go)**:
- `github.com/gin-gonic/gin` - HTTP framework
- `github.com/golang-jwt/jwt` - JWT handling
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/prometheus/client_golang` - Metrics
- `go.opentelemetry.io/otel` - Distributed tracing
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/go-playground/validator` - Input validation

**Infrastructure**:
- Kubernetes 1.28+
- PostgreSQL 16
- NGINX Ingress Controller
- Prometheus + Grafana
- cert-manager (TLS certificates)

### Development Tools

- Go 1.21+
- Node.js 20+
- Docker 24+
- kubectl
- Helm 3
- k9s (Kubernetes CLI)
- Postman/Insomnia (API testing)

## Deployment Checklist

### Pre-Deployment

- [ ] All tests passing (unit, integration, e2e)
- [ ] Security scan completed with no critical vulnerabilities
- [ ] Database migrations tested and ready
- [ ] Environment variables configured
- [ ] Secrets created in Kubernetes
- [ ] TLS certificates valid
- [ ] Monitoring dashboards configured
- [ ] Alerting rules configured
- [ ] Backup strategy verified
- [ ] Rollback plan documented

### Post-Deployment

- [ ] Health checks passing for all services
- [ ] Smoke tests passing
- [ ] Metrics flowing to Prometheus
- [ ] Logs appearing in aggregation system
- [ ] Database connections stable
- [ ] API response times within SLA
- [ ] Error rate < 0.1%
- [ ] User registration flow tested
- [ ] Savings creation flow tested
- [ ] Notification delivery verified

### Monitoring Alerts

**Critical Alerts** (PagerDuty):
- Service down (all replicas unhealthy)
- Database connection failure
- Error rate > 5%
- Response time p99 > 5 seconds
- Disk usage > 90%

**Warning Alerts** (Slack):
- Error rate > 1%
- Response time p95 > 1 second
- Memory usage > 80%
- CPU usage > 80%
- Replication lag > 5 seconds

## Conclusion

The InSavein platform is designed as a highly scalable, fault-tolerant distributed system capable of serving millions of users. The architecture leverages modern technologies (TanStack Start, Golang microservices, PostgreSQL, Kubernetes) to provide a robust foundation for financial discipline and savings tracking. With comprehensive security measures, observability, and automated CI/CD pipelines, the platform is production-ready and maintainable for long-term growth.
