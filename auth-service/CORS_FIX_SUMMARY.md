# CORS Fix Summary

## Problem
Registration and login were failing with "failed to fetch" error because the auth service was missing CORS (Cross-Origin Resource Sharing) configuration.

## Root Cause
The auth service was using standard `http.ServeMux` without CORS middleware, while the frontend runs on a different port (5173) than the backend (8080). Browsers block cross-origin requests without proper CORS headers.

## Inconsistencies Found

### 1. Missing CORS Configuration
- **Other services** (notification, goal, education) had CORS configured with `github.com/rs/cors`
- **Auth service** had NO CORS configuration
- This caused browser to block all requests from frontend to auth service

### 2. Port Configuration Mismatch
- `frontend/.env`: Auth service on port 8080 ✓
- `frontend/.env.example`: Auth service on port 8081 ✗ (FIXED)
- `auth-service/.env`: Service runs on port 8080 ✓

### 3. Missing Dependency
- `auth-service/go.mod` was missing `github.com/rs/cors` package

## Changes Made

### 1. Added CORS Dependency
```go
// auth-service/go.mod
require (
    github.com/rs/cors v1.10.1  // Added
)
```

### 2. Updated Main Server File
```go
// auth-service/cmd/server/main.go
import (
    "github.com/rs/cors"  // Added import
)

// Setup CORS (added before server creation)
c := cors.New(cors.Options{
    AllowedOrigins:   []string{"*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: true,
    MaxAge:           300,
})

// Wrap handler with CORS
server := &http.Server{
    Handler: c.Handler(mux),  // Changed from: Handler: mux
}
```

### 3. Fixed Port Configuration
Updated `frontend/.env.example` to match actual port configuration.

## Verification
CORS is now working correctly:
```
Status: 204
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: POST
Access-Control-Allow-Origin: *
Access-Control-Max-Age: 300
```

## Next Steps
1. Test user registration from the frontend
2. Test user login from the frontend
3. Verify token refresh works correctly

## Security Note
Currently using `AllowedOrigins: []string{"*"}` for development. In production, this should be restricted to specific domains:
```go
AllowedOrigins: []string{"https://yourdomain.com"},
```
