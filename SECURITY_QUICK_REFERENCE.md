# Security Implementation - Quick Reference Card

Quick reference for developers implementing security features in InSavein services.

## Input Validation

### Add Validator to Handler

```go
import "github.com/go-playground/validator/v10"

type Handler struct {
    service  Service
    validate *validator.Validate
}

func NewHandler(service Service) *Handler {
    return &Handler{
        service:  service,
        validate: validator.New(),
    }
}
```

### Add Validation Tags to Structs

```go
type CreateRequest struct {
    Amount      float64 `json:"amount" validate:"required,gt=0"`
    Currency    string  `json:"currency" validate:"required,len=3"`
    Description string  `json:"description" validate:"max=500"`
    Email       string  `json:"email" validate:"required,email"`
}
```

### Common Validation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field must be present | `validate:"required"` |
| `email` | Valid email format | `validate:"email"` |
| `min=X` | Minimum length/value | `validate:"min=8"` |
| `max=X` | Maximum length/value | `validate:"max=100"` |
| `len=X` | Exact length | `validate:"len=3"` |
| `gt=X` | Greater than | `validate:"gt=0"` |
| `gte=X` | Greater than or equal | `validate:"gte=18"` |
| `lt=X` | Less than | `validate:"lt=100"` |
| `lte=X` | Less than or equal | `validate:"lte=1000"` |

### Validate in Handler

```go
func (h *Handler) CreateResource(w http.ResponseWriter, r *http.Request) {
    var req CreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    // Validate
    if err := h.validate.Struct(req); err != nil {
        respondValidationError(w, err)
        return
    }
    
    // Sanitize strings
    req.Description = sanitizeString(req.Description)
    
    // Process request...
}
```

### Sanitize Strings

```go
import "html"

func sanitizeString(s string) string {
    // HTML escape to prevent XSS
    s = html.EscapeString(s)
    
    // Remove SQL injection patterns
    s = strings.ReplaceAll(s, "--", "")
    s = strings.ReplaceAll(s, ";", "")
    
    return strings.TrimSpace(s)
}
```

## Authorization

### Check Resource Ownership

```go
func (h *Handler) GetResource(w http.ResponseWriter, r *http.Request) {
    // Get authenticated user ID
    userID, ok := r.Context().Value("user_id").(string)
    if !ok {
        respondError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
    
    // Get resource ID from URL
    resourceID := mux.Vars(r)["id"]
    
    // Verify ownership
    resource, err := h.service.GetResource(r.Context(), resourceID)
    if err != nil {
        respondError(w, http.StatusNotFound, "Resource not found")
        return
    }
    
    if resource.UserID != userID {
        respondError(w, http.StatusForbidden, "Forbidden: you do not have access to this resource")
        return
    }
    
    // Return resource...
}
```

### Check User Role

```go
func (h *Handler) AdminOnlyEndpoint(w http.ResponseWriter, r *http.Request) {
    // Get roles from context
    roles, ok := r.Context().Value("roles").([]string)
    if !ok {
        respondError(w, http.StatusForbidden, "Forbidden")
        return
    }
    
    // Check for admin role
    hasAdmin := false
    for _, role := range roles {
        if role == "admin" {
            hasAdmin = true
            break
        }
    }
    
    if !hasAdmin {
        respondError(w, http.StatusForbidden, "Forbidden: admin access required")
        return
    }
    
    // Process admin request...
}
```

## Error Responses

### Standard Error Response

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

func respondError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorResponse{
        Error:   http.StatusText(status),
        Message: message,
    })
}
```

### Validation Error Response

```go
type ValidationError struct {
    Field   string `json:"field"`
    Tag     string `json:"tag"`
    Value   string `json:"value,omitempty"`
    Message string `json:"message"`
}

type ValidationErrorResponse struct {
    Error   string            `json:"error"`
    Message string            `json:"message"`
    Errors  []ValidationError `json:"errors"`
}

func respondValidationError(w http.ResponseWriter, err error) {
    var errors []ValidationError
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            errors = append(errors, ValidationError{
                Field:   e.Field(),
                Tag:     e.Tag(),
                Value:   e.Param(),
                Message: getErrorMessage(e),
            })
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(ValidationErrorResponse{
        Error:   "Bad Request",
        Message: "Validation failed",
        Errors:  errors,
    })
}
```

## HTTP Status Codes

| Code | Status | When to Use |
|------|--------|-------------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful POST (resource created) |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Invalid input, validation failed |
| 401 | Unauthorized | Missing or invalid authentication |
| 403 | Forbidden | Authenticated but not authorized |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Duplicate resource, constraint violation |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Unexpected server error |

## Security Headers (Ingress)

Already configured in `k8s/ingress.yaml`:

```yaml
annotations:
  # Force HTTPS
  nginx.ingress.kubernetes.io/ssl-redirect: "true"
  
  # Rate limiting
  nginx.ingress.kubernetes.io/limit-rps: "100"
  nginx.ingress.kubernetes.io/limit-rpm: "100"
  
  # Security headers (in configuration-snippet)
  # - Strict-Transport-Security
  # - X-Frame-Options
  # - X-Content-Type-Options
  # - X-XSS-Protection
  # - Content-Security-Policy
  # - Referrer-Policy
```

## Rate Limiting

### Current Limits

- **Per User**: 100 requests/minute (burst: 20)
- **Per IP**: 1000 requests/minute (burst: 50)
- **Auth Endpoints**: 5 requests/minute (burst: 2)

### Rate Limit Response

```json
{
  "error": "Too Many Requests",
  "message": "Rate limit exceeded. Please try again later.",
  "retry_after": 60
}
```

## Testing

### Test Validation

```bash
# Test with invalid data
curl -X POST https://api.insavein.com/api/savings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount": -10, "currency": "USD"}'

# Expected: 400 Bad Request with validation errors
```

### Test Authorization

```bash
# Test without token
curl https://api.insavein.com/api/savings

# Expected: 401 Unauthorized

# Test with invalid token
curl -H "Authorization: Bearer invalid" https://api.insavein.com/api/savings

# Expected: 401 Unauthorized
```

### Test Rate Limiting

```bash
# Send 110 requests
for i in {1..110}; do
  curl -H "Authorization: Bearer $TOKEN" https://api.insavein.com/api/savings
done

# Expected: 429 after ~100 requests
```

## Common Patterns

### Handler with Full Security

```go
package handlers

import (
    "encoding/json"
    "html"
    "net/http"
    "strings"
    
    "github.com/go-playground/validator/v10"
    "github.com/gorilla/mux"
)

type Handler struct {
    service  Service
    validate *validator.Validate
}

func NewHandler(service Service) *Handler {
    return &Handler{
        service:  service,
        validate: validator.New(),
    }
}

func (h *Handler) CreateResource(w http.ResponseWriter, r *http.Request) {
    // 1. Get authenticated user
    userID, ok := r.Context().Value("user_id").(string)
    if !ok {
        respondError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
    
    // 2. Parse request
    var req CreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    // 3. Validate
    if err := h.validate.Struct(req); err != nil {
        respondValidationError(w, err)
        return
    }
    
    // 4. Sanitize
    req.Description = sanitizeString(req.Description)
    
    // 5. Process
    resource, err := h.service.CreateResource(r.Context(), userID, req)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    // 6. Respond
    respondJSON(w, http.StatusCreated, resource)
}

func (h *Handler) GetResource(w http.ResponseWriter, r *http.Request) {
    // 1. Get authenticated user
    userID, ok := r.Context().Value("user_id").(string)
    if !ok {
        respondError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
    
    // 2. Get resource ID
    resourceID := mux.Vars(r)["id"]
    
    // 3. Fetch resource
    resource, err := h.service.GetResource(r.Context(), resourceID)
    if err != nil {
        respondError(w, http.StatusNotFound, "Resource not found")
        return
    }
    
    // 4. Check ownership
    if resource.UserID != userID {
        respondError(w, http.StatusForbidden, "Forbidden")
        return
    }
    
    // 5. Respond
    respondJSON(w, http.StatusOK, resource)
}

func sanitizeString(s string) string {
    s = html.EscapeString(s)
    s = strings.ReplaceAll(s, "--", "")
    s = strings.ReplaceAll(s, ";", "")
    return strings.TrimSpace(s)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, map[string]string{
        "error":   http.StatusText(status),
        "message": message,
    })
}
```

## Checklist for New Endpoints

- [ ] Add validation tags to request struct
- [ ] Validate request in handler
- [ ] Sanitize string inputs
- [ ] Check user authentication (get user_id from context)
- [ ] Check authorization (verify ownership or role)
- [ ] Return appropriate HTTP status codes
- [ ] Return detailed error messages
- [ ] Use parameterized queries (prevent SQL injection)
- [ ] Test with invalid inputs
- [ ] Test without authentication
- [ ] Test with wrong user

## Resources

- **Deployment Guide**: `TASK_25_SECURITY_IMPLEMENTATION.md`
- **Middleware Docs**: `shared/middleware/README.md`
- **Completion Summary**: `TASK_25_COMPLETION_SUMMARY.md`
- **Test Script**: `test-security-implementation.sh`
- **Deployment Checklist**: `TASK_25_DEPLOYMENT_CHECKLIST.md`

## Support

For questions or issues:
- Review documentation in `shared/middleware/README.md`
- Check examples in `savings-service/internal/handlers/savings_handler_enhanced.go`
- Run tests with `test-security-implementation.sh`
- Contact security team: security@insavein.com
