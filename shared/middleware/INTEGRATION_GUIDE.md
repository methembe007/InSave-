# Shared Middleware Integration Guide

## Problem Solved

The `shared/middleware` package is now a proper Go module with all dependencies installed. The validator import error has been resolved.

## What Was Done

1. Created `go.mod` file in `shared/middleware/`
2. Installed `github.com/go-playground/validator/v10` dependency
3. Verified the package compiles successfully

## How to Use in Your Services

Since the shared middleware is now a local module, you have two options for using it in your services:

### Option 1: Copy Files Directly (Recommended for Now)

Copy the middleware files directly into each service:

```bash
# For each service (auth-service, user-service, savings-service, etc.)
mkdir -p auth-service/internal/middleware
cp shared/middleware/validation.go auth-service/internal/middleware/
cp shared/middleware/authorization.go auth-service/internal/middleware/
```

Then update the package declaration in the copied files:

```go
// Change from:
package middleware

// To:
package middleware
```

And update imports in your handlers:

```go
import (
    "github.com/insavein/auth-service/internal/middleware"
)
```

### Option 2: Use Go Workspace (Advanced)

Create a Go workspace to use the shared module:

```bash
# In project root
go work init
go work use ./shared/middleware
go work use ./auth-service
go work use ./user-service
# ... add all services
```

Then in each service's `go.mod`, add:

```go
require github.com/insavein/shared/middleware v0.0.0

replace github.com/insavein/shared/middleware => ../shared/middleware
```

## Quick Integration Steps

### Step 1: Add Validator to Service

```bash
cd auth-service
go get github.com/go-playground/validator/v10
```

### Step 2: Copy Middleware Files

```bash
# Windows
copy ..\shared\middleware\validation.go internal\middleware\
copy ..\shared\middleware\authorization.go internal\middleware\

# Linux/Mac
cp ../shared/middleware/validation.go internal/middleware/
cp ../shared/middleware/authorization.go internal/middleware/
```

### Step 3: Update Handler

```go
package handlers

import (
    "encoding/json"
    "net/http"
    
    "github.com/insavein/auth-service/internal/middleware"
    "github.com/insavein/auth-service/internal/auth"
)

type AuthHandler struct {
    service   auth.Service
    validator *middleware.ValidationMiddleware
}

func NewAuthHandler(service auth.Service) *AuthHandler {
    return &AuthHandler{
        service:   service,
        validator: middleware.NewValidationMiddleware(),
    }
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req auth.RegisterRequest
    
    // Decode request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    // Validate
    if err := h.validator.Validate(&req); err != nil {
        respondValidationError(w, err)
        return
    }
    
    // Process...
}
```

### Step 4: Add Validation Tags

```go
type RegisterRequest struct {
    Email       string `json:"email" validate:"required,email"`
    Password    string `json:"password" validate:"required,min=8"`
    FirstName   string `json:"first_name" validate:"required"`
    LastName    string `json:"last_name" validate:"required"`
    DateOfBirth string `json:"date_of_birth" validate:"required"`
}
```

## Verification

Test that the middleware works:

```bash
cd auth-service
go build ./...
```

If successful, you should see no errors.

## Troubleshooting

### Import Error Still Occurs

If you still see import errors after copying files:

1. Make sure validator is installed in the service:
   ```bash
   cd auth-service
   go get github.com/go-playground/validator/v10
   go mod tidy
   ```

2. Verify the import path matches your service:
   ```go
   import "github.com/insavein/auth-service/internal/middleware"
   ```

3. Rebuild:
   ```bash
   go build ./...
   ```

### Module Not Found

If you see "module not found" errors:

1. Check your `go.mod` file exists in the service directory
2. Run `go mod tidy` in the service directory
3. Verify the module path in `go.mod` matches your imports

## Next Steps

1. Copy middleware files to each service
2. Add validator dependency to each service
3. Update handlers to use validation
4. Add validation tags to request structs
5. Test with invalid inputs

See `shared/middleware/README.md` for detailed usage examples.
