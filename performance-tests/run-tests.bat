@echo off
REM Performance Test Runner Script for InSavein Platform (Windows)

setlocal enabledelayedexpansion

set "BASE_URL=%BASE_URL%"
if "%BASE_URL%"=="" set "BASE_URL=http://localhost:8080"

set "TEST_TYPE=%1"
if "%TEST_TYPE%"=="" set "TEST_TYPE=normal"

set "OUTPUT_DIR=.\results"

echo === InSavein Performance Testing ===
echo.

REM Check if k6 is installed
where k6 >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: k6 is not installed
    echo Install k6 from: https://k6.io/docs/getting-started/installation/
    exit /b 1
)

REM Create output directory
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

REM Get timestamp
for /f "tokens=2 delims==" %%I in ('wmic os get localdatetime /value') do set datetime=%%I
set "TIMESTAMP=%datetime:~0,8%_%datetime:~8,6%"

echo Configuration:
echo   Base URL: %BASE_URL%
echo   Test Type: %TEST_TYPE%
echo   Output Directory: %OUTPUT_DIR%
echo.

if "%TEST_TYPE%"=="normal" (
    call :run_test "normal-load.js" "normal-load"
) else if "%TEST_TYPE%"=="peak" (
    call :run_test "peak-load.js" "peak-load"
) else if "%TEST_TYPE%"=="stress" (
    call :run_test "stress-test.js" "stress-test"
) else if "%TEST_TYPE%"=="all" (
    echo Running all tests sequentially...
    echo.
    call :run_test "normal-load.js" "normal-load"
    echo.
    timeout /t 60 /nobreak
    call :run_test "peak-load.js" "peak-load"
    echo.
    timeout /t 60 /nobreak
    call :run_test "stress-test.js" "stress-test"
) else (
    echo Error: Unknown test type '%TEST_TYPE%'
    echo Usage: %0 [normal^|peak^|stress^|all]
    exit /b 1
)

echo.
echo === Performance Testing Complete ===
echo.
echo Results are available in: %OUTPUT_DIR%
echo.
echo To analyze results:
echo   - Review the JSON output files
echo   - Check p95/p99 latencies against targets (p95^<500ms, p99^<1000ms)
echo   - Verify error rate is below 0.1%%
echo   - Monitor system resources during tests

exit /b 0

:run_test
set "test_file=%~1"
set "test_name=%~2"
set "output_file=%OUTPUT_DIR%\%test_name%_%TIMESTAMP%.json"

echo Running %test_name%...

k6 run -e BASE_URL=%BASE_URL% --out json=%output_file% %test_file%

if %ERRORLEVEL% EQU 0 (
    echo [OK] %test_name% completed successfully
    echo   Results saved to: %output_file%
) else (
    echo [FAIL] %test_name% failed
)

exit /b 0
