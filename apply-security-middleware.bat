@echo off
REM Script to apply security middleware to all InSavein services
REM This script adds go-playground/validator dependency to all services

setlocal enabledelayedexpansion

echo ================================================
echo InSavein Security Middleware Application Script
echo ================================================
echo.

REM List of services
set SERVICES=auth-service user-service savings-service budget-service goal-service education-service notification-service analytics-service

echo Step 1: Adding validator dependency to all services...
echo.

set BUILD_ERRORS=0

for %%s in (%SERVICES%) do (
    if exist "%%s" (
        echo Processing %%s...
        
        cd %%s
        
        if exist "go.mod" (
            echo   - Adding go-playground/validator/v10...
            go get github.com/go-playground/validator/v10
            
            echo   - Running go mod tidy...
            go mod tidy
            
            echo   [32m✓ %%s updated successfully[0m
        ) else (
            echo   [31m✗ go.mod not found in %%s[0m
        )
        
        cd ..
        echo.
    ) else (
        echo   [31m✗ Directory %%s not found[0m
        echo.
    )
)

echo ================================================
echo Step 2: Verifying installations...
echo ================================================
echo.

for %%s in (%SERVICES%) do (
    if exist "%%s\go.mod" (
        findstr /C:"github.com/go-playground/validator/v10" "%%s\go.mod" >nul
        if !errorlevel! equ 0 (
            echo [32m✓ %%s: validator installed[0m
        ) else (
            echo [31m✗ %%s: validator NOT installed[0m
        )
    )
)

echo.
echo ================================================
echo Step 3: Building services to verify...
echo ================================================
echo.

for %%s in (%SERVICES%) do (
    if exist "%%s" (
        echo Building %%s...
        
        cd %%s
        
        go build -o nul .\cmd\server 2>nul
        if !errorlevel! equ 0 (
            echo   [32m✓ %%s builds successfully[0m
        ) else (
            echo   [31m✗ %%s build failed[0m
            set /a BUILD_ERRORS+=1
        )
        
        cd ..
        echo.
    )
)

echo ================================================
echo Summary
echo ================================================
echo.

if !BUILD_ERRORS! equ 0 (
    echo [32m✓ All services updated and building successfully![0m
    echo.
    echo Next steps:
    echo 1. Update handler files to use validator
    echo 2. Add validation tags to request structs
    echo 3. Rebuild Docker images
    echo 4. Deploy to Kubernetes
    echo.
    echo See TASK_25_SECURITY_IMPLEMENTATION.md for detailed instructions.
) else (
    echo [31m✗ !BUILD_ERRORS! service(s) failed to build[0m
    echo.
    echo Please check the errors above and fix them before proceeding.
)

echo.
echo ================================================
echo Done!
echo ================================================

endlocal
pause
