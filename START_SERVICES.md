# Quick Start Guide - InSavein Services

## Problem: "Failed to fetch" when registering

This error occurs because the backend services aren't running. The frontend needs the auth service (and other services) to be running to handle API requests.

## Solution: Start the Backend Services

### Prerequisites
1. ✅ PostgreSQL database running (port 5432)
2. ✅ Database migrations completed
3. ✅ Go installed (for running services)

### Start Auth Service (Required for Registration/Login)

**Option 1: Using Go directly**
```bash
cd auth-service
go run cmd/server/main.go
```

**Option 2: Using Make**
```bash
cd auth-service
make run
```

The auth service will start on **port 8081** and you should see:
```
Connected to database successfully
Auth service starting on port 8081
```

### Verify Auth Service is Running

Open a new terminal and test:
```bash
curl http://localhost:8081/health
```

Expected response:
```json
{"status":"healthy","service":"auth-service"}
```

### Start Other Services (Optional, for full functionality)

Once auth is working, you can start the other services:

**User Service (Port 8082)**
```bash
cd user-service
go run cmd/server/main.go
```

**Savings Service (Port 8083)**
```bash
cd savings-service
go run cmd/server/main.go
```

**Budget Service (Port 8084)**
```bash
cd budget-service
go run cmd/server/main.go
```

**Goal Service (Port 8085)**
```bash
cd goal-service
go run cmd/server/main.go
```

**Education Service (Port 8086)**
```bash
cd education-service
go run cmd/server/main.go
```

**Notification Service (Port 8087)**
```bash
cd notification-service
go run cmd/server/main.go
```

**Analytics Service (Port 8088)**
```bash
cd analytics-service
go run cmd/server/main.go
```

## Quick Test After Starting Auth Service

1. Start the auth service (see above)
2. Go back to your frontend (http://localhost:3000)
3. Try registering a new user
4. You should now see proper validation or success messages instead of "Failed to fetch"

## Troubleshooting

### "Failed to connect to database"
- Make sure PostgreSQL is running
- Check database credentials in `auth-service/.env`
- Verify database exists: `psql -U insavein_user -d insavein_db -c "\dt"`

### "Port already in use"
- Check if another process is using port 8081
- Windows: `netstat -ano | findstr "8081"`
- Kill the process or change the PORT in `.env`

### "Cannot find module"
- Run `go mod download` in the service directory
- Run `go mod tidy` to clean up dependencies

## Environment Configuration

I've created `auth-service/.env` with the correct port (8081) to match your frontend configuration.

**Auth Service:** Port 8081 ✅
**Frontend expects:** Port 8081 ✅

## Next Steps

1. Start the auth service
2. Test registration/login
3. Start other services as needed
4. Enjoy building with InSavein! 🎉
