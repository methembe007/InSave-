@echo off
REM Partition Maintenance Script for Windows
REM This script provides convenient commands for managing database partitions

setlocal enabledelayedexpansion

REM Configuration
if "%DB_NAME%"=="" set DB_NAME=insavein
if "%DB_USER%"=="" set DB_USER=postgres
if "%DB_HOST%"=="" set DB_HOST=localhost
if "%DB_PORT%"=="" set DB_PORT=5432

REM Get command
set COMMAND=%1

if "%COMMAND%"=="" (
    call :show_help
    exit /b 0
)

if /i "%COMMAND%"=="setup" (
    call :setup
) else if /i "%COMMAND%"=="create" (
    call :create %2 %3
) else if /i "%COMMAND%"=="maintain" (
    call :maintain
) else if /i "%COMMAND%"=="status" (
    call :status
) else if /i "%COMMAND%"=="archive" (
    call :archive %2
) else if /i "%COMMAND%"=="list-detached" (
    call :list_detached
) else if /i "%COMMAND%"=="verify" (
    call :verify
) else if /i "%COMMAND%"=="help" (
    call :show_help
) else (
    echo [ERROR] Unknown command: %COMMAND%
    echo.
    call :show_help
    exit /b 1
)

exit /b 0

:setup
echo [INFO] Setting up partition management...
echo.
echo [INFO] Creating partition management functions...
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -f create_monthly_partitions.sql
echo.
echo [INFO] Setting up automatic partition creation...
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -f auto_create_partitions.sql
echo.
echo [INFO] Setting up archival functions...
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -f archive_old_partitions.sql
echo.
echo [SUCCESS] Partition management setup complete!
exit /b 0

:create
set START_DATE=%~1
set END_DATE=%~2

if "%START_DATE%"=="" (
    REM Default: 6 months ago
    for /f "tokens=1-3 delims=/ " %%a in ('date /t') do (
        set START_DATE=%%c-%%a-01
    )
)

if "%END_DATE%"=="" (
    REM Default: 6 months ahead
    for /f "tokens=1-3 delims=/ " %%a in ('date /t') do (
        set END_DATE=%%c-%%a-01
    )
)

echo [INFO] Creating partitions from %START_DATE% to %END_DATE%...
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT create_partitions_for_range('savings_transactions', '%START_DATE%', '%END_DATE%');"
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT create_partitions_for_range('spending_transactions', '%START_DATE%', '%END_DATE%');"
echo [SUCCESS] Partitions created successfully!
exit /b 0

:maintain
echo [INFO] Running partition maintenance...
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT maintain_partitions();"
echo [SUCCESS] Partition maintenance complete!
exit /b 0

:status
echo [INFO] Partition Status Report
echo.
echo [INFO] Active Partitions:
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT parent.relname AS table_name, child.relname AS partition_name, TO_DATE(SUBSTRING(child.relname FROM '(\d{4}_\d{2})$'), 'YYYY_MM') AS partition_month, pg_size_pretty(pg_total_relation_size(child.oid)) AS size FROM pg_inherits JOIN pg_class parent ON pg_inherits.inhparent = parent.oid JOIN pg_class child ON pg_inherits.inhrelid = child.oid WHERE parent.relname IN ('savings_transactions', 'spending_transactions') ORDER BY parent.relname, partition_month DESC LIMIT 20;"
echo.
echo [INFO] Total Size by Table:
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size FROM pg_tables WHERE tablename IN ('savings_transactions', 'spending_transactions') ORDER BY tablename;"
exit /b 0

:archive
set RETENTION_MONTHS=%~1
if "%RETENTION_MONTHS%"=="" set RETENTION_MONTHS=12

echo [WARNING] This will archive partitions older than %RETENTION_MONTHS% months
set /p CONFIRM="Continue? (y/N): "

if /i "%CONFIRM%"=="y" (
    echo [INFO] Archiving old partitions...
    psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT archive_workflow('savings_transactions', %RETENTION_MONTHS%);"
    psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT archive_workflow('spending_transactions', %RETENTION_MONTHS%);"
    echo [SUCCESS] Archival complete!
) else (
    echo [INFO] Archival cancelled
)
exit /b 0

:list_detached
echo [INFO] Detached Partitions:
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size FROM pg_tables WHERE schemaname = 'public' AND (tablename LIKE 'savings_transactions_%%' OR tablename LIKE 'spending_transactions_%%') AND tablename NOT IN (SELECT child.relname FROM pg_inherits JOIN pg_class child ON pg_inherits.inhrelid = child.oid) ORDER BY tablename;"
exit /b 0

:verify
echo [INFO] Verifying partition health...
echo [INFO] Checking for missing future partitions...
REM Note: Date arithmetic in batch is complex, so we'll just run the SQL check
psql -h %DB_HOST% -p %DB_PORT% -U %DB_USER% -d %DB_NAME% -c "SELECT COUNT(*) as missing_partitions FROM generate_series(DATE_TRUNC('month', CURRENT_DATE), DATE_TRUNC('month', CURRENT_DATE + INTERVAL '2 months'), '1 month'::interval) AS month WHERE NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'savings_transactions_' || TO_CHAR(month, 'YYYY_MM'));"
echo [INFO] Run 'maintain' if any partitions are missing
exit /b 0

:show_help
echo Partition Maintenance Script
echo.
echo Usage: %~nx0 ^<command^> [options]
echo.
echo Commands:
echo     setup                   Initial setup of partition management functions
echo     create [start] [end]    Create partitions for date range
echo     maintain                Run daily maintenance (create future partitions)
echo     status                  Show partition status and sizes
echo     archive [months]        Archive partitions older than N months (default: 12)
echo     list-detached           List detached partitions
echo     verify                  Verify partition health
echo     help                    Show this help message
echo.
echo Environment Variables:
echo     DB_NAME                 Database name (default: insavein)
echo     DB_USER                 Database user (default: postgres)
echo     DB_HOST                 Database host (default: localhost)
echo     DB_PORT                 Database port (default: 5432)
echo.
echo Examples:
echo     # Initial setup
echo     %~nx0 setup
echo.
echo     # Create partitions for 2024
echo     %~nx0 create 2024-01-01 2024-12-31
echo.
echo     # Daily maintenance
echo     %~nx0 maintain
echo.
echo     # Check status
echo     %~nx0 status
echo.
echo     # Archive partitions older than 18 months
echo     %~nx0 archive 18
echo.
echo     # Verify partition health
echo     %~nx0 verify
exit /b 0
