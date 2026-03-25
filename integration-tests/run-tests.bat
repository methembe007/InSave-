@echo off
REM Integration Test Runner Script for InSavein Platform (Windows)
REM This script sets up the test environment and runs integration tests

echo =========================================
echo InSavein Integration Test Runner
echo =========================================
echo.

REM Check if Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is not running. Please start Docker and try again.
    exit /b 1
)

REM Check if Docker Compose is available
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] docker-compose is not installed. Please install it and try again.
    exit /b 1
)

REM Navigate to integration-tests directory
cd /d "%~dp0"

echo [INFO] Starting test environment...
docker-compose -f docker-compose.test.yml up -d

echo [INFO] Waiting for services to be healthy...
timeout /t 10 /nobreak >nul

REM Wait for PostgreSQL to be ready
echo [INFO] Checking PostgreSQL health...
set /a counter=0
:wait_postgres
docker-compose -f docker-compose.test.yml exec -T postgres-test pg_isready -U postgres >nul 2>&1
if errorlevel 1 (
    set /a counter+=1
    if %counter% geq 30 (
        echo [ERROR] PostgreSQL failed to start
        docker-compose -f docker-compose.test.yml logs postgres-test
        exit /b 1
    )
    timeout /t 2 /nobreak >nul
    goto wait_postgres
)
echo [INFO] PostgreSQL is ready

REM Wait for services to be healthy
echo [INFO] Waiting for microservices to be healthy...
timeout /t 20 /nobreak >nul

REM Check service health
echo [INFO] Checking service health...
curl -f -s "http://localhost:18080/health" >nul 2>&1 && echo [INFO] auth-service is healthy || echo [WARN] auth-service health check failed
curl -f -s "http://localhost:18081/health" >nul 2>&1 && echo [INFO] user-service is healthy || echo [WARN] user-service health check failed
curl -f -s "http://localhost:18082/health" >nul 2>&1 && echo [INFO] savings-service is healthy || echo [WARN] savings-service health check failed
curl -f -s "http://localhost:18083/health" >nul 2>&1 && echo [INFO] budget-service is healthy || echo [WARN] budget-service health check failed
curl -f -s "http://localhost:18005/health" >nul 2>&1 && echo [INFO] goal-service is healthy || echo [WARN] goal-service health check failed

REM Run tests
echo [INFO] Running integration tests...
echo.

if "%1"=="coverage" (
    echo [INFO] Running tests with coverage...
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    echo [INFO] Coverage report generated: coverage.html
) else if not "%1"=="" (
    echo [INFO] Running specific test: %1
    go test -v -run "%1" ./...
) else (
    go test -v ./...
)

set TEST_EXIT_CODE=%errorlevel%

echo.
if %TEST_EXIT_CODE%==0 (
    echo [INFO] All tests passed!
) else (
    echo [ERROR] Some tests failed!
)

REM Cleanup option
if "%2"=="cleanup" (
    echo [INFO] Cleaning up test environment...
    docker-compose -f docker-compose.test.yml down -v
    echo [INFO] Cleanup complete
) else if "%1"=="cleanup" (
    echo [INFO] Cleaning up test environment...
    docker-compose -f docker-compose.test.yml down -v
    echo [INFO] Cleanup complete
) else (
    echo [WARN] Test environment is still running. Use 'run-tests.bat cleanup' to stop it.
    echo [INFO] View logs: docker-compose -f docker-compose.test.yml logs [service-name]
)

exit /b %TEST_EXIT_CODE%
