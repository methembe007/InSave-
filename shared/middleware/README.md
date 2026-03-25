# Shared Middleware Package

This package provides reusable middleware for input validation, authorization, and security across all InSavein microservices.

## Features

### 1. Input Validation Middleware (`validation.go`)

Provides comprehensive input validation and sanitization to prevent SQL injection and XSS attacks.

**Features:**
- Struct validation using `go-playground/validator`
- Automatic HTML escaping to prevent XSS
- SQL injection pattern detection and removal
- Detailed validation error messages
- HTTP 400 Bad Request responses with field-level errors

**Usage Example:**

```go
package handlers

import (
    "github.com/insavein/shared/middleware"
    "github.com/go-playground/validator/v10"
)

type CreateSavingsRequest struct {
    Amount      float64 `json:"amount" validate:"required,gt=0"`
    Currency    string  `json:"currency" validate:"required,len=3"`
    Description string  `json:"description" validate:"max=500"`
    Category    string  `json:"category" validate:"required"`
}

type SavingsHandler struct {
    service   SavingsService
    validator *middleware.ValidationMiddleware
}

func NewSavingsHandler(service SavingsService) *SavingsHandler {
    return &SavingsHandler{
        service:   service,
        validator: middleware.NewValidationMiddleware(),
    }
}

func (h *SavingsHandler) CreateSavings(w http.ResponseWriter, r *http.Request) {
    var req CreateSavingsRequest
    
    // Decode and validate
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Validate
    if err := h.validator.Validate(&req); err != nil {
        // Validation errors are automatically formatted
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Process request...
}
```

**Validation Tags:**

- `required` - Field must be present
- `email` - Must be valid email format
- `min=X` - Minimum length/value
- `max=X` - Maximum length/value
- `gt=X` - Greater than X
- `gte=X` - Greater than or equal to X
- `lt=X` - Less than X
- `lte=X` - Less than or equal to X
- `len=X` - Exact length

### 2. Authorization Middleware (`authorization.go`)

Ensures users can only access resources they own and enforces role-based access control.

**Features:**
- Resource ownership verification
- Self-or-admin access patterns
- Role-based access control
- HTTP 403 Forbidden responses for unauthorized access

**Usage Example:**

```go
package handlers

import (
    "context"
    "github.com/insavein/shared/middleware"
    "github.com/gorilla/mux"
)

type GoalHandler struct {
    service       GoalService
    authz         *middleware.AuthorizationMiddleware
}

func NewGoalHandler(service GoalService) *GoalHandler {
    return &GoalHandler{
        service: service,
        authz:   middleware.NewAuthorizationMiddleware(),
    }
}

// Setup routes with authorization
func (h *GoalHandler) RegisterRoutes(r *mux.Router) {
    // Ownership checker function
    checkGoalOwnership := func(ctx context.Context, userID string, goalID string) (bool, error) {
        goal, err := h.service.GetGoal(ctx, goalID)
        if err != nil {
            return false, err
        }
        return goal.UserID == userID, nil
    }
    
    // Protected route - requires goal ownership
    r.HandleFunc("/api/goals/{id}", 
        h.authz.RequireOwnership(checkGoalOwnership, "id")(h.UpdateGoal),
    ).Methods("PUT")
    
    // Protected route - requires self or admin
    r.HandleFunc("/api/users/{id}", 
        h.authz.RequireSelfOrAdmin("id")(h.GetUserProfile),
    ).Methods("GET")
    
    // Protected route - requires admin role
    r.HandleFunc("/api/admin/users", 
        h.authz.RequireRole("admin")(h.ListAllUsers),
    ).Methods("GET")
}
```

**Authorization Methods:**

1. **RequireOwnership** - Verifies user owns the resource
   - Takes a checker function and resource ID parameter name
   - Returns 403 if user doesn't own the resource

2. **RequireSelfOrAdmin** - Allows access to own data or admin access
   - Takes target user ID parameter name
   - Returns 403 if not self and not admin

3. **RequireRole** - Requires specific role
   - Takes required role name
   - Returns 403 if user doesn't have the role

## Integration with Services

### Step 1: Add to go.mod

```bash
# In each service directory
go get github.com/go-playground/validator/v10
```

### Step 2: Import in Handlers

```go
import (
    "github.com/insavein/shared/middleware"
)
```

### Step 3: Initialize Middleware

```go
validator := middleware.NewValidationMiddleware()
authz := middleware.NewAuthorizationMiddleware()
```

### Step 4: Apply to Routes

```go
// With gorilla/mux
r.HandleFunc("/api/resource", 
    authMiddleware.Authenticate(
        validator.ValidateRequest(&RequestType{})(
            authz.RequireOwnership(checker, "id")(handler),
        ),
    ),
).Methods("POST")
```

## Security Best Practices

### Input Validation
- ✅ Always validate all user inputs
- ✅ Use strict validation rules (min, max, format)
- ✅ Sanitize strings to prevent XSS
- ✅ Return detailed validation errors (400 Bad Request)
- ❌ Never trust client-side validation alone

### Authorization
- ✅ Always check resource ownership
- ✅ Verify JWT token user_id matches resource owner
- ✅ Use role-based access for admin functions
- ✅ Return 403 Forbidden (not 404) for unauthorized access
- ❌ Never expose whether a resource exists to unauthorized users

### SQL Injection Prevention
- ✅ Use parameterized queries (prepared statements)
- ✅ Sanitize inputs (remove SQL keywords)
- ✅ Use ORM/query builders when possible
- ✅ Validate input types and formats
- ❌ Never concatenate user input into SQL queries

### XSS Prevention
- ✅ HTML escape all user-generated content
- ✅ Use Content-Security-Policy headers
- ✅ Validate and sanitize all inputs
- ✅ Encode output based on context (HTML, JS, URL)
- ❌ Never render unsanitized user input

## Testing

### Unit Tests

```go
func TestValidationMiddleware(t *testing.T) {
    validator := middleware.NewValidationMiddleware()
    
    type TestRequest struct {
        Email string `validate:"required,email"`
        Age   int    `validate:"required,gte=18"`
    }
    
    // Valid request
    req := TestRequest{Email: "test@example.com", Age: 25}
    err := validator.Validate(&req)
    assert.NoError(t, err)
    
    // Invalid email
    req = TestRequest{Email: "invalid", Age: 25}
    err = validator.Validate(&req)
    assert.Error(t, err)
    
    // Age too low
    req = TestRequest{Email: "test@example.com", Age: 15}
    err = validator.Validate(&req)
    assert.Error(t, err)
}
```

## Requirements Mapping

This middleware implementation satisfies the following requirements:

- **Requirement 17.1**: Return HTTP 400 Bad Request with detailed validation errors
- **Requirement 17.2**: Include field name in error messages
- **Requirement 17.3**: Specify valid range in error messages
- **Requirement 17.5**: Validate all inputs against defined schemas
- **Requirement 17.6**: Sanitize inputs to prevent SQL injection and XSS
- **Requirement 3.5**: Only allow users to access their own data
- **Requirement 15.4**: Return HTTP 403 Forbidden for unauthorized access
- **Requirement 20.4**: Encrypt data in transit (TLS)
- **Requirement 20.5**: Security headers (HSTS, CSP, etc.)
- **Requirement 20.6**: Rate limiting (100 req/min per user, 1000 req/min per IP)
