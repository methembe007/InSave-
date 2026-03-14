@echo off
echo ========================================
echo InSavein Platform - Quick Start
echo ========================================
echo.

echo Step 1: Starting PostgreSQL Database...
docker-compose up -d postgres
timeout /t 5 /nobreak >nul

echo.
echo Step 2: Running Database Migrations...
cd migrations
call migrate.bat up
cd ..

echo.
echo Step 3: Starting Auth Service (Port 8081)...
start "Auth Service" cmd /k "cd auth-service && go run cmd/server/main.go"

echo.
echo ========================================
echo Services Starting...
echo ========================================
echo.
echo Auth Service: http://localhost:8081
echo Frontend: http://localhost:3000
echo.
echo Check auth service health:
echo   curl http://localhost:8081/health
echo.
echo Press any key to start additional services...
pause >nul

echo.
echo Starting User Service (Port 8082)...
start "User Service" cmd /k "cd user-service && go run cmd/server/main.go"

echo.
echo Starting Savings Service (Port 8083)...
start "Savings Service" cmd /k "cd savings-service && go run cmd/server/main.go"

echo.
echo Starting Budget Service (Port 8084)...
start "Budget Service" cmd /k "cd budget-service && go run cmd/server/main.go"

echo.
echo Starting Goal Service (Port 8085)...
start "Goal Service" cmd /k "cd goal-service && go run cmd/server/main.go"

echo.
echo Starting Education Service (Port 8086)...
start "Education Service" cmd /k "cd education-service && go run cmd/server/main.go"

echo.
echo Starting Notification Service (Port 8087)...
start "Notification Service" cmd /k "cd notification-service && go run cmd/server/main.go"

echo.
echo Starting Analytics Service (Port 8088)...
start "Analytics Service" cmd /k "cd analytics-service && go run cmd/server/main.go"

echo.
echo ========================================
echo All Services Started!
echo ========================================
echo.
echo You can now use the InSavein application at:
echo   http://localhost:3000
echo.
echo To stop all services, close the terminal windows.
echo.
pause
