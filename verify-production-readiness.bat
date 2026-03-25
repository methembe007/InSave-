@echo off
REM Production Readiness Verification Script (Windows)
REM InSavein Platform

setlocal enabledelayedexpansion

set NAMESPACE=insavein
set PASSED=0
set FAILED=0
set WARNINGS=0

echo ========================================
echo InSavein Production Readiness Verification
echo ========================================
echo.

REM Check prerequisites
echo === Checking Prerequisites ===
echo.

where kubectl >nul 2>&1
if %errorlevel% equ 0 (
    echo [PASS] kubectl is installed
    set /a PASSED+=1
) else (
    echo [FAIL] kubectl is not installed
    set /a FAILED+=1
    goto :end
)

where docker >nul 2>&1
if %errorlevel% equ 0 (
    echo [PASS] docker is installed
    set /a PASSED+=1
) else (
    echo [WARN] docker is not installed ^(optional^)
    set /a WARNINGS+=1
)

where psql >nul 2>&1
if %errorlevel% equ 0 (
    echo [PASS] psql is installed
    set /a PASSED+=1
) else (
    echo [WARN] psql is not installed ^(optional^)
    set /a WARNINGS+=1
)

echo.
echo === Kubernetes Cluster Connectivity ===
echo.

kubectl cluster-info >nul 2>&1
if %errorlevel% equ 0 (
    echo [PASS] Connected to Kubernetes cluster
    set /a PASSED+=1
    kubectl cluster-info | findstr "Kubernetes"
) else (
    echo [FAIL] Cannot connect to Kubernetes cluster
    set /a FAILED+=1
    goto :end
)

echo.
echo === Namespace Verification ===
echo.

kubectl get namespace %NAMESPACE% >nul 2>&1
if %errorlevel% equ 0 (
    echo [PASS] Namespace '%NAMESPACE%' exists
    set /a PASSED+=1
) else (
    echo [FAIL] Namespace '%NAMESPACE%' does not exist
    set /a FAILED+=1
    goto :end
)

echo.
echo === Service Deployment Status ===
echo.

set SERVICES=auth-service user-service savings-service budget-service goal-service education-service notification-service analytics-service frontend

for %%s in (%SERVICES%) do (
    kubectl get deployment %%s -n %NAMESPACE% >nul 2>&1
    if !errorlevel! equ 0 (
        echo [PASS] %%s: Deployment exists
        set /a PASSED+=1
    ) else (
        echo [FAIL] %%s: Deployment not found
        set /a FAILED+=1
    )
)

echo.
echo === Pod Health Status ===
echo.

kubectl get pods -n %NAMESPACE% --field-selector=status.phase=Running >nul 2>&1
if %errorlevel% equ 0 (
    echo [PASS] Pods are running
    set /a PASSED+=1
    kubectl get pods -n %NAMESPACE%
) else (
    echo [FAIL] Not all pods are running
    set /a FAILED+=1
    kubectl get pods -n %NAMESPACE%
)

echo.
echo === Database Status ===
echo.

kubectl get statefulset postgres -n %NAMESPACE% >nul 2>&1
if %errorlevel% equ 0 (
    echo [PASS] PostgreSQL StatefulSet exists
    set /a PASSED+=1
) else (
    echo [WARN] PostgreSQL StatefulSet not found ^(may be external^)
    set /a WARNINGS+=1
)

echo.
echo === Monitoring Stack ===
echo.

kubectl get deployment prometheus -n %NAMESPACE% >nul 2>&1
if %errorlevel% equ 0 (
    echo [PASS] Prometheus is deployed
    set /a PASSED+=1
) else (
    echo [WARN] Prometheus deployment not found
    set /a WARNINGS+=1
)

kubectl get deployment grafana -n %NAMESPACE% >nul 2>&1
if %errorlevel% equ 0 (
    echo [PASS] Grafana is deployed
    set /a PASSED+=1
) else (
    echo [WARN] Grafana deployment not found
    set /a WARNINGS+=1
)

echo.
echo === Documentation ===
echo.

if exist "README.md" (
    echo [PASS] README.md exists
    set /a PASSED+=1
) else (
    echo [WARN] README.md not found
    set /a WARNINGS+=1
)

if exist "docs\API_DOCUMENTATION.md" (
    echo [PASS] API_DOCUMENTATION.md exists
    set /a PASSED+=1
) else (
    echo [WARN] API_DOCUMENTATION.md not found
    set /a WARNINGS+=1
)

if exist "docs\DEPLOYMENT.md" (
    echo [PASS] DEPLOYMENT.md exists
    set /a PASSED+=1
) else (
    echo [WARN] DEPLOYMENT.md not found
    set /a WARNINGS+=1
)

echo.
echo === CI/CD Pipelines ===
echo.

if exist ".github\workflows\lint.yml" (
    echo [PASS] lint.yml exists
    set /a PASSED+=1
) else (
    echo [WARN] lint.yml not found
    set /a WARNINGS+=1
)

if exist ".github\workflows\test.yml" (
    echo [PASS] test.yml exists
    set /a PASSED+=1
) else (
    echo [WARN] test.yml not found
    set /a WARNINGS+=1
)

if exist ".github\workflows\security.yml" (
    echo [PASS] security.yml exists
    set /a PASSED+=1
) else (
    echo [WARN] security.yml not found
    set /a WARNINGS+=1
)

:end
echo.
echo === Verification Summary ===
echo.
echo Passed:   %PASSED%
echo Warnings: %WARNINGS%
echo Failed:   %FAILED%
echo.

if %FAILED% equ 0 (
    echo ========================================
    echo Production readiness verification completed successfully!
    echo ========================================
    echo.
    echo Next steps:
    echo 1. Review warnings and address if necessary
    echo 2. Run integration tests: cd integration-tests ^&^& make test
    echo 3. Run performance tests: cd performance-tests ^&^& make test-normal
    echo 4. Run security scans: trivy image ^<image-name^>
    echo 5. Review PRODUCTION_READINESS_CHECKLIST.md
    exit /b 0
) else (
    echo ========================================
    echo Production readiness verification failed!
    echo ========================================
    echo.
    echo Please address the failed checks before proceeding to production.
    exit /b 1
)
