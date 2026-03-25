# Troubleshooting Auth Check Failures

## Error: "Auth check failed: {}"

This error appears in the browser console when the frontend cannot connect to the backend services.

### Root Cause
The frontend is trying to verify authentication by calling the user service API, but the service is not running or not accessible.

### Quick Fix

**Step 1: Start the backend services**

Run the quick start script:
```bash
QUICK_START.bat
```

Or manually start the auth and user services:
```bash
# Terminal 1 - Auth Service
cd auth-service
go run cmd/server/main.go

# Terminal 2 - User Service  
cd user-service
go run cmd/server/main.go
```

**Step 2: Verify services are running**

Test auth service:
```bash
curl http://localhost:8080/health
```

Test user service:
```bash
curl http://localhost:8081/health
```

Expected response:
```json
{"status":"healthy","service":"auth-service"}
```

**Step 3: Refresh your browser**

Once the services are running, refresh your frontend application.

### Service Port Configuration

Make sure these services are running on the correct ports:

| Service | Port | Status |
|---------|------|--------|
| Auth | 8080 | Required for login/register |
| User | 8081 | Required for profile/auth check |
| Savings | 8082 | Optional |
| Budget | 8083 | Optional |
| Goal | 8005 | Optional |
| Education | 8085 | Optional |
| Notification | 8086 | Optional |
| Analytics | 8008 | Optional |

### What Changed

I've improved the error handling to:
1. Show better error messages in the console
2. Only clear tokens on actual auth failures (not network errors)
3. Provide helpful messages when services are unavailable

### Still Having Issues?

1. Check if PostgreSQL is running
2. Verify database migrations are complete
3. Check service logs for errors
4. Ensure no port conflicts (use `netstat -ano | findstr "8080"`)

### Prevention

To avoid this error in the future:
1. Always start backend services before using the frontend
2. Use the QUICK_START.bat script for easy startup
3. Check service health endpoints before testing
