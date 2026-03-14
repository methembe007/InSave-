@echo off
REM Docker Build and Test Script for InSavein Platform (Windows)
REM This script builds all Docker images and tests the deployment

setlocal enabledelayedexpansion

echo ==========================================
echo InSavein Platform - Docker Build ^& Test
echo ==========================================
echo.

REM Check if Docker is running
echo Checking Docker status...
docker info >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is not running. Please start Docker Desktop and try again.
    exit /b 1
)
echo [OK] Docker is running
echo.

REM Check if Docker Compose is available
echo Checking Docker Compose...
docker-compose version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker Compose is not installed or not in PATH
    exit /b 1
)
echo [OK] Docker Compose is available
echo.

REM Build all services
echo ==========================================
echo Building Docker Images
echo ==========================================
echo.

set services=auth-service user-service savings-service budget-service goal-service education-service notification-service analytics-service frontend

for %%s in (%services%) do (
    echo Building %%s...
    docker-compose build %%s
    if errorlevel 1 (
        echo [ERROR] %%s build failed
        exit /b 1
    )
    echo [OK] %%s built successfully
    echo.
)

echo [OK] All services built successfully!
echo.

REM Start services
echo ==========================================
echo Starting Services
echo ==========================================
echo.

echo [INFO] Starting PostgreSQL databases...
docker-compose up -d postgres-primary postgres-replica1 postgres-replica2

echo [INFO] Waiting for databases to be healthy (30 seconds)...
timeout /t 30 /nobreak >nul

echo [INFO] Starting microservices...
docker-compose up -d auth-service user-service savings-service budget-service goal-service education-service notification-service analytics-service

echo [INFO] Waiting for microservices to be healthy (20 seconds)...
timeout /t 20 /nobreak >nul

echo [INFO] Starting frontend...
docker-compose up -d frontend

echo [INFO] Waiting for frontend to be healthy (15 seconds)...
timeout /t 15 /nobreak >nul

echo.
echo ==========================================
echo Service Health Checks
echo ==========================================
echo.

REM Check health of all services
curl -f -s http://localhost:8080/health >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Auth Service (port 8080) is not responding
) else (
    echo [OK] Auth Service (port 8080) is healthy
)

curl -f -s http://localhost:8081/health >nul 2>&1
if errorlevel 1 (
    echo [ERROR] User Service (port 8081) is not responding
) else (
    echo [OK] User Service (port 8081) is healthy
)

curl -f -s http://localhost:8082/health >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Savings Service (port 8082) is not responding
) else (
    echo [OK] Savings Service (port 8082) is healthy
)

curl -f -s http://localhost:8083/health >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Budget Service (port 8083) is not responding
) else (
    echo [OK] Budget Service (port 8083) is healthy
)

curl -f -s http://localhost:8005/health >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Goal Service (port 8005) is not responding
) else (
    echo [OK] Goal Service (port 8005) is healthy
)

curl -f -s http://localhost:8085/health >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Education Service (port 8085) is not responding
) else (
    echo [OK] Education Service (port 8085) is healthy
)

curl -f -s http://localhost:8086/health >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Notification Service (port 8086) is not responding
) else (
    echo [OK] Notification Service (port 8086) is healthy
)

curl -f -s http://localhost:8008/health >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Analytics Service (port 8008) is not responding
) else (
    echo [OK] Analytics Service (port 8008) is healthy
)

curl -f -s http://localhost:3000 >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Frontend (port 3000) is not responding
) else (
    echo [OK] Frontend (port 3000) is responding
)

echo.
echo ==========================================
echo Container Status
echo ==========================================
echo.

docker-compose ps

echo.
echo ==========================================
echo Basic Functionality Tests
echo ==========================================
echo.

echo [INFO] Test 1: User Registration
curl -s -X POST http://localhost:8080/api/auth/register ^
    -H "Content-Type: application/json" ^
    -d "{\"email\":\"test@example.com\",\"password\":\"testpassword123\",\"first_name\":\"Test\",\"last_name\":\"User\",\"date_of_birth\":\"1990-01-01\"}" > register_response.json

findstr /C:"access_token" register_response.json >nul
if errorlevel 1 (
    echo [ERROR] User registration failed
    type register_response.json
) else (
    echo [OK] User registration successful
)

echo [INFO] Test 2: User Login
curl -s -X POST http://localhost:8080/api/auth/login ^
    -H "Content-Type: application/json" ^
    -d "{\"email\":\"test@example.com\",\"password\":\"testpassword123\"}" > login_response.json

findstr /C:"access_token" login_response.json >nul
if errorlevel 1 (
    echo [ERROR] User login failed
    type login_response.json
) else (
    echo [OK] User login successful
)

REM Clean up temp files
del register_response.json login_response.json 2>nul

echo.
echo ==========================================
echo Test Summary
echo ==========================================
echo.

echo [OK] Docker build and deployment test completed!
echo.
echo [INFO] Services are running at:
echo   - Frontend:              http://localhost:3000
echo   - Auth Service:          http://localhost:8080
echo   - User Service:          http://localhost:8081
echo   - Savings Service:       http://localhost:8082
echo   - Budget Service:        http://localhost:8083
echo   - Goal Service:          http://localhost:8005
echo   - Education Service:     http://localhost:8085
echo   - Notification Service:  http://localhost:8086
echo   - Analytics Service:     http://localhost:8008
echo.
echo [INFO] To view logs: docker-compose logs -f [service-name]
echo [INFO] To stop all services: docker-compose down
echo [INFO] To stop and remove volumes: docker-compose down -v
echo.

endlocal
