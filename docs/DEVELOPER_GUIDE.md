# InSavein Platform Developer Guide

Complete guide for setting up local development environment and contributing to the InSavein platform.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Prerequisites](#prerequisites)
3. [Local Development Setup](#local-development-setup)
4. [Project Structure](#project-structure)
5. [Running Tests](#running-tests)
6. [Code Conventions](#code-conventions)
7. [Adding New Services](#adding-new-services)
8. [Database Management](#database-management)
9. [Debugging](#debugging)
10. [Contributing](#contributing)

---

## Quick Start

Get up and running in 5 minutes:

```bash
# 1. Clone repository
git clone https://github.com/insavein/platform.git
cd platform

# 2. Start all services with Docker Compose
docker-compose up -d

# 3. Run database migrations
cd migrations && ./migrate.sh && cd ..

# 4. Access the application
# Frontend: http://localhost:3000
# API: http://localhost:8080-8086, 8008
```

---

## Prerequisites

### Required Software

**Core Tools**:
- **Git** 2.30+
- **Docker** 20.10+
- **Docker Compose** 2.0+

**Backend Development** (optional, if not using Docker):
- **Go** 1.21+
- **PostgreSQL** 15+

**Frontend Development** (optional, if not using Docker):
- **Node.js** 20+
- **npm** 10+

### Recommended Tools

- **VS Code** with extensions:
  - Go (golang.go)
  - ESLint
  - Prettier
  - Docker
  - Kubernetes
- **Postman** or **Insomnia** for API testing
- **pgAdmin** or **DBeaver** for database management
- **k9s** for Kubernetes management (if deploying to k8s)

---

## Local Development Setup

### Option 1: Docker Compose (Recommended)

**Advantages**:
- Consistent environment across team
- No need to install Go, Node.js, PostgreSQL locally
- Easy to start/stop all services
- Matches production environment

**Setup**:

```bash
# 1. Clone repository
git clone https://github.com/insavein/platform.git
cd platform

# 2. Copy environment files
cp auth-service/.env.example auth-service/.env
cp user-service/.env.example user-service/.env
cp savings-service/.env.example savings-service/.env
cp budget-service/.env.example budget-service/.env
cp goal-service/.env.example goal-service/.env
cp education-service/.env.example education-service/.env
cp notification-service/.env.example notification-service/.env
cp analytics-service/.env.example analytics-service/.env
cp frontend/.env.example frontend/.env

# 3. Start all services
docker-compose up -d

# 4. Check service status
docker-compose ps

# 5. View logs
docker-compose logs -f

# 6. Run migrations
docker-compose exec postgres-primary psql -U insavein_user -d insavein
# Then run: \i /docker-entrypoint-initdb.d/init.sql
# Or use migrate tool:
cd migrations && ./migrate.sh
```

**Access Services**:
- Frontend: http://localhost:3000
- Auth Service: http://localhost:8080
- User Service: http://localhost:8081
- Savings Service: http://localhost:8082
- Budget Service: http://localhost:8083
- Goal Service: http://localhost:8005
- Education Service: http://localhost:8085
- Notification Service: http://localhost:8086
- Analytics Service: http://localhost:8008
- PostgreSQL Primary: localhost:5432
- PostgreSQL Replica 1: localhost:5433
- PostgreSQL Replica 2: localhost:5434
- PgBouncer: localhost:6432

**Development Workflow**:

```bash
# Make code changes in your editor

# Rebuild specific service
docker-compose build auth-service
docker-compose up -d auth-service

# View logs
docker-compose logs -f auth-service

# Restart service
docker-compose restart auth-service

# Stop all services
docker-compose down

# Stop and remove volumes (⚠️ deletes data)
docker-compose down -v
```

### Option 2: Manual Setup

**Advantages**:
- Faster iteration (no Docker rebuild)
- Direct debugging
- More control over environment

**Setup**:

#### 1. Install Dependencies

**Go**:
```bash
# Download from https://golang.org/dl/
# Or use package manager:
# macOS: brew install go
# Ubuntu: sudo apt install golang-go
# Windows: choco install golang

# Verify
go version  # Should be 1.21+
```

**Node.js**:
```bash
# Download from https://nodejs.org/
# Or use nvm:
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 20
nvm use 20

# Verify
node --version  # Should be 20+
npm --version   # Should be 10+
```

**PostgreSQL**:
```bash
# macOS: brew install postgresql@15
# Ubuntu: sudo apt install postgresql-15
# Windows: Download from https://www.postgresql.org/download/

# Start PostgreSQL
# macOS: brew services start postgresql@15
# Ubuntu: sudo systemctl start postgresql
# Windows: Start from Services

# Create database
psql -U postgres
CREATE DATABASE insavein;
CREATE USER insavein_user WITH PASSWORD 'insavein_password';
GRANT ALL PRIVILEGES ON DATABASE insavein TO insavein_user;
\q
```

#### 2. Setup Backend Services

```bash
# Clone repository
git clone https://github.com/insavein/platform.git
cd platform

# Setup each service
cd auth-service
cp .env.example .env
# Edit .env with your database credentials
go mod download
go run cmd/server/main.go

# Repeat for other services in separate terminals
```

#### 3. Setup Frontend

```bash
cd frontend
cp .env.example .env
# Edit .env with API URLs
npm install
npm run dev
```

#### 4. Run Migrations

```bash
cd migrations
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=insavein_user
export DB_PASSWORD=insavein_password
export DB_NAME=insavein
./migrate.sh
```

---

## Project Structure

```
insavein-platform/
├── .github/                    # GitHub Actions workflows
│   └── workflows/
│       ├── test.yml           # Run tests on PR
│       ├── lint.yml           # Code linting
│       ├── build-push.yml     # Build and push Docker images
│       └── deploy-*.yml       # Deployment workflows
│
├── auth-service/              # Authentication microservice
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # Service entry point
│   ├── internal/
│   │   ├── auth/              # Business logic
│   │   │   ├── service.go     # Service interface
│   │   │   ├── auth_service.go # Implementation
│   │   │   ├── repository.go  # Data access interface
│   │   │   ├── postgres_repository.go # PostgreSQL implementation
│   │   │   └── types.go       # Domain types
│   │   ├── handlers/          # HTTP handlers
│   │   │   └── auth_handler.go
│   │   └── middleware/        # HTTP middleware
│   │       ├── auth_middleware.go
│   │       └── metrics.go
│   ├── pkg/                   # Shared packages
│   │   └── database/
│   │       └── postgres.go
│   ├── .env.example           # Environment template
│   ├── Dockerfile             # Container image
│   ├── go.mod                 # Go dependencies
│   ├── go.sum                 # Dependency checksums
│   ├── Makefile               # Build commands
│   └── README.md              # Service documentation
│
├── user-service/              # User profile microservice
├── savings-service/           # Savings tracking microservice
├── budget-service/            # Budget planning microservice
├── goal-service/              # Goal management microservice
├── education-service/         # Education content microservice
├── notification-service/      # Notification delivery microservice
├── analytics-service/         # Analytics microservice
│   # (Same structure as auth-service)
│
├── frontend/                  # TanStack Start frontend
│   ├── src/
│   │   ├── routes/            # Page routes
│   │   │   ├── index.tsx      # Landing page
│   │   │   ├── login.tsx      # Login page
│   │   │   ├── register.tsx   # Registration page
│   │   │   ├── dashboard.tsx  # Dashboard
│   │   │   ├── savings.tsx    # Savings tracker
│   │   │   ├── budget.tsx     # Budget planner
│   │   │   └── goals.tsx      # Goal manager
│   │   ├── components/        # React components
│   │   │   ├── DashboardLayout.tsx
│   │   │   ├── SavingsTracker.tsx
│   │   │   ├── BudgetPlanner.tsx
│   │   │   └── GoalManager.tsx
│   │   ├── lib/               # Utilities
│   │   │   ├── api/           # API client
│   │   │   │   ├── client.ts  # HTTP client
│   │   │   │   ├── auth.ts    # Auth API
│   │   │   │   ├── savings.ts # Savings API
│   │   │   │   └── ...
│   │   │   ├── auth/          # Auth utilities
│   │   │   │   ├── context.tsx # Auth context
│   │   │   │   └── storage.ts  # Token storage
│   │   │   └── hooks/         # Custom hooks
│   │   └── test/              # Test utilities
│   ├── public/                # Static assets
│   ├── .env.example           # Environment template
│   ├── Dockerfile             # Container image
│   ├── package.json           # Dependencies
│   ├── tsconfig.json          # TypeScript config
│   └── README.md              # Frontend documentation
│
├── shared/                    # Shared code across services
│   ├── middleware/            # Shared middleware
│   │   ├── authorization.go   # RBAC middleware
│   │   ├── validation.go      # Input validation
│   │   └── README.md
│   └── cache/                 # Caching utilities
│       ├── redis_cache.go
│       └── cache_middleware.go
│
├── migrations/                # Database migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_savings_transactions_table.up.sql
│   ├── ...
│   ├── seed/                  # Seed data
│   │   ├── 001_sample_users.sql
│   │   └── ...
│   ├── migrate.sh             # Migration script
│   └── README.md
│
├── k8s/                       # Kubernetes manifests
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secrets.yaml
│   ├── postgres-statefulset.yaml
│   ├── *-deployment.yaml      # Service deployments
│   ├── ingress.yaml
│   ├── prometheus-deployment.yaml
│   ├── grafana-deployment.yaml
│   └── README.md
│
├── integration-tests/         # Integration tests
│   ├── user_registration_test.go
│   ├── goal_progress_flow_test.go
│   └── README.md
│
├── performance-tests/         # Load tests (k6)
│   ├── normal-load.js
│   ├── peak-load.js
│   └── README.md
│
├── docs/                      # Documentation
│   ├── API_DOCUMENTATION.md   # API reference
│   ├── DEPLOYMENT.md          # Deployment guide
│   ├── DEVELOPER_GUIDE.md     # This file
│   └── OPERATIONS_RUNBOOK.md  # Operations guide
│
├── docker-compose.yml         # Local development setup
├── Makefile                   # Project-wide commands
├── README.md                  # Project overview
└── .gitignore                 # Git ignore rules
```

### Service Structure Pattern

All Go microservices follow this structure:

```
<service-name>/
├── cmd/server/main.go         # Entry point
├── internal/
│   ├── <domain>/              # Business logic
│   │   ├── service.go         # Interface
│   │   ├── <domain>_service.go # Implementation
│   │   ├── repository.go      # Data interface
│   │   ├── postgres_repository.go # Data implementation
│   │   └── types.go           # Domain types
│   ├── handlers/              # HTTP handlers
│   └── middleware/            # Middleware
├── pkg/                       # Shared packages
├── .env.example
├── Dockerfile
├── go.mod
├── Makefile
└── README.md
```

---

## Running Tests

### Backend Tests (Go)

**Unit Tests**:
```bash
# Run tests for specific service
cd auth-service
go test ./...

# Run with coverage
go test ./... -cover

# Run with verbose output
go test ./... -v

# Run specific test
go test ./internal/auth -run TestAuthService_Login

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Integration Tests**:
```bash
cd integration-tests

# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run tests
go test -v ./...

# Clean up
docker-compose -f docker-compose.test.yml down -v
```

### Frontend Tests (TypeScript/React)

```bash
cd frontend

# Run all tests
npm test

# Run with coverage
npm test -- --coverage

# Run in watch mode
npm test -- --watch

# Run specific test file
npm test -- src/lib/auth/__tests__/storage.test.ts

# Run E2E tests (if configured)
npm run test:e2e
```

### Performance Tests (k6)

```bash
cd performance-tests

# Install k6
# macOS: brew install k6
# Ubuntu: sudo apt install k6
# Windows: choco install k6

# Run normal load test
k6 run normal-load.js

# Run peak load test
k6 run peak-load.js

# Run with custom parameters
k6 run --vus 100 --duration 5m normal-load.js
```

### Test Coverage Goals

- **Unit Tests**: 80%+ coverage for business logic
- **Integration Tests**: Cover critical user flows
- **E2E Tests**: Cover main user journeys
- **Performance Tests**: Validate SLAs (p95 < 500ms)

---

## Code Conventions

### Go Code Style

**Follow**:
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

**Key Conventions**:

```go
// Package names: lowercase, single word
package auth

// Interface names: noun or noun phrase
type Service interface {
    Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

// Struct names: PascalCase
type AuthService struct {
    repo Repository
    logger Logger
}

// Function names: PascalCase (exported), camelCase (unexported)
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
    // Implementation
}

func (s *AuthService) validateCredentials(email, password string) error {
    // Implementation
}

// Constants: PascalCase or SCREAMING_SNAKE_CASE
const (
    DefaultTimeout = 30 * time.Second
    MAX_RETRIES    = 3
)

// Error handling: always check errors
result, err := someFunction()
if err != nil {
    return nil, fmt.Errorf("failed to do something: %w", err)
}

// Context: always first parameter
func DoSomething(ctx context.Context, param string) error {
    // Implementation
}
```

**Linting**:
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

### TypeScript/React Code Style

**Follow**:
- [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)

**Key Conventions**:

```typescript
// Component names: PascalCase
export function SavingsTracker() {
  // Implementation
}

// Hook names: camelCase, start with "use"
export function useSavings() {
  // Implementation
}

// Type names: PascalCase
export type SavingsSummary = {
  totalSaved: number
  currentStreak: number
}

// Interface names: PascalCase, prefix with "I" optional
export interface ApiClient {
  auth: AuthApi
  savings: SavingsApi
}

// Constants: SCREAMING_SNAKE_CASE
export const API_BASE_URL = 'http://localhost:8080'
export const MAX_RETRIES = 3

// Async/await: prefer over promises
async function fetchSavings(): Promise<SavingsSummary> {
  const response = await apiClient.savings.getSummary()
  return response.data
}

// Error handling: try/catch
try {
  const data = await fetchSavings()
} catch (error) {
  console.error('Failed to fetch savings:', error)
}
```

**Linting**:
```bash
cd frontend

# Run ESLint
npm run lint

# Auto-fix issues
npm run lint -- --fix

# Run Prettier
npm run format
```

### Git Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding/updating tests
- `chore`: Maintenance tasks

**Examples**:
```
feat(auth): add password reset functionality

Implement password reset flow with email verification.
Includes rate limiting to prevent abuse.

Closes #123

fix(savings): correct streak calculation for timezone edge cases

The streak calculation was failing when users saved money
around midnight in different timezones.

docs: update API documentation with new endpoints

chore(deps): upgrade Go to 1.21.5
```

---

## Adding New Services

### Step 1: Create Service Structure

```bash
# Create service directory
mkdir new-service
cd new-service

# Initialize Go module
go mod init github.com/insavein/platform/new-service

# Create directory structure
mkdir -p cmd/server
mkdir -p internal/newdomain
mkdir -p internal/handlers
mkdir -p internal/middleware
mkdir -p pkg/database

# Create main.go
cat > cmd/server/main.go << 'EOF'
package main

import (
    "log"
    "net/http"
)

func main() {
    log.Println("Starting new-service on :8090")
    http.ListenAndServe(":8090", nil)
}
EOF
```

### Step 2: Implement Service Interface

```go
// internal/newdomain/service.go
package newdomain

import "context"

type Service interface {
    DoSomething(ctx context.Context, req Request) (*Response, error)
}

type Request struct {
    // Fields
}

type Response struct {
    // Fields
}
```

### Step 3: Implement Repository

```go
// internal/newdomain/repository.go
package newdomain

import "context"

type Repository interface {
    Create(ctx context.Context, entity Entity) error
    GetByID(ctx context.Context, id string) (*Entity, error)
}

// internal/newdomain/postgres_repository.go
package newdomain

import (
    "context"
    "database/sql"
)

type postgresRepository struct {
    db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
    return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, entity Entity) error {
    // Implementation
    return nil
}
```

### Step 4: Add HTTP Handlers

```go
// internal/handlers/handler.go
package handlers

import (
    "encoding/json"
    "net/http"
    
    "github.com/insavein/platform/new-service/internal/newdomain"
)

type Handler struct {
    service newdomain.Service
}

func NewHandler(service newdomain.Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Implementation
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

### Step 5: Add to Docker Compose

```yaml
# docker-compose.yml
services:
  new-service:
    build:
      context: ./new-service
      dockerfile: Dockerfile
    ports:
      - "8090:8090"
    environment:
      - PORT=8090
      - DB_HOST=postgres-primary
      - DB_PORT=5432
      - DB_USER=insavein_user
      - DB_PASSWORD=insavein_password
      - DB_NAME=insavein
    depends_on:
      - postgres-primary
    networks:
      - insavein-network
```

### Step 6: Create Kubernetes Deployment

```yaml
# k8s/new-service-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: new-service
  namespace: insavein
spec:
  replicas: 2
  selector:
    matchLabels:
      app: new-service
  template:
    metadata:
      labels:
        app: new-service
    spec:
      containers:
      - name: new-service
        image: insavein/new-service:latest
        ports:
        - containerPort: 8090
        env:
        - name: PORT
          value: "8090"
        # Add other env vars
---
apiVersion: v1
kind: Service
metadata:
  name: new-service
  namespace: insavein
spec:
  selector:
    app: new-service
  ports:
  - port: 8090
    targetPort: 8090
```

### Step 7: Add Tests

```go
// internal/newdomain/service_test.go
package newdomain_test

import (
    "context"
    "testing"
    
    "github.com/insavein/platform/new-service/internal/newdomain"
)

func TestService_DoSomething(t *testing.T) {
    // Setup
    service := newdomain.NewService(mockRepo)
    
    // Execute
    result, err := service.DoSomething(context.Background(), req)
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result == nil {
        t.Fatal("expected result, got nil")
    }
}
```

---

## Database Management

### Migrations

**Create Migration**:
```bash
cd migrations

# Create new migration
migrate create -ext sql -dir . -seq create_new_table

# This creates:
# 000XXX_create_new_table.up.sql
# 000XXX_create_new_table.down.sql
```

**Run Migrations**:
```bash
# Up (apply migrations)
./migrate.sh

# Down (rollback last migration)
./migrate.sh down 1

# Goto specific version
./migrate.sh goto 5

# Force version (if stuck)
./migrate.sh force 5
```

### Seed Data

```bash
cd migrations/seed

# Run all seeds
psql -U insavein_user -d insavein -f run_seeds.sql

# Run specific seed
psql -U insavein_user -d insavein -f 001_sample_users.sql
```

### Database Access

```bash
# Connect to database
psql -h localhost -p 5432 -U insavein_user -d insavein

# Or with Docker
docker-compose exec postgres-primary psql -U insavein_user -d insavein

# Common queries
\dt                    # List tables
\d users               # Describe table
SELECT * FROM users;   # Query data
```

---

## Debugging

### Backend Debugging (Go)

**VS Code Launch Configuration** (`.vscode/launch.json`):
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Auth Service",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/auth-service/cmd/server",
      "env": {
        "PORT": "8080",
        "DB_HOST": "localhost",
        "DB_PORT": "5432"
      }
    }
  ]
}
```

**Delve (CLI Debugger)**:
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug service
cd auth-service
dlv debug cmd/server/main.go

# Set breakpoint
(dlv) break main.main
(dlv) continue
```

### Frontend Debugging (React)

**Browser DevTools**:
- Chrome DevTools: F12
- React DevTools: Install extension
- Network tab: Monitor API calls
- Console: View logs and errors

**VS Code Debugging**:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Frontend",
      "type": "chrome",
      "request": "launch",
      "url": "http://localhost:3000",
      "webRoot": "${workspaceFolder}/frontend/src"
    }
  ]
}
```

### Docker Debugging

```bash
# View logs
docker-compose logs -f auth-service

# Execute command in container
docker-compose exec auth-service sh

# Inspect container
docker inspect <container-id>

# Check resource usage
docker stats
```

---

## Contributing

### Workflow

1. **Fork repository**
2. **Create feature branch**: `git checkout -b feature/amazing-feature`
3. **Make changes**
4. **Run tests**: `go test ./...` and `npm test`
5. **Run linters**: `golangci-lint run` and `npm run lint`
6. **Commit changes**: Follow commit message conventions
7. **Push to branch**: `git push origin feature/amazing-feature`
8. **Open Pull Request**

### Pull Request Checklist

- [ ] Tests added/updated
- [ ] Tests passing
- [ ] Linters passing
- [ ] Documentation updated
- [ ] Commit messages follow conventions
- [ ] No merge conflicts
- [ ] Reviewed by at least one team member

### Code Review Guidelines

**As Author**:
- Keep PRs small and focused
- Provide context in PR description
- Respond to feedback promptly
- Update PR based on feedback

**As Reviewer**:
- Review within 24 hours
- Be constructive and respectful
- Focus on code quality, not style
- Approve when satisfied

---

## Additional Resources

- [README.md](../README.md) - Project overview
- [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) - API reference
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Deployment guide
- [OPERATIONS_RUNBOOK.md](./OPERATIONS_RUNBOOK.md) - Operations guide
- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://react.dev/)
- [TanStack Start Documentation](https://tanstack.com/start)

---

**Last Updated**: 2026-01-15  
**Version**: 1.0.0
