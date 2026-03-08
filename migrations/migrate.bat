@echo off
REM InSavein Platform Database Migration Script for Windows
REM This script provides convenient commands for managing database migrations

setlocal enabledelayedexpansion

REM Default values
set MIGRATIONS_DIR=migrations
if "%DATABASE_URL%"=="" set DATABASE_URL=postgresql://postgres:postgres@localhost:5432/insavein?sslmode=disable

REM Check if migrate is installed
where migrate >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: golang-migrate is not installed
    echo Install it with: scoop install migrate
    echo Or download from: https://github.com/golang-migrate/migrate/releases
    exit /b 1
)

REM Parse command
set COMMAND=%1
if "%COMMAND%"=="" set COMMAND=help

if "%COMMAND%"=="up" goto migrate_up
if "%COMMAND%"=="down" goto migrate_down
if "%COMMAND%"=="create" goto create_migration
if "%COMMAND%"=="version" goto show_version
if "%COMMAND%"=="force" goto force_version
if "%COMMAND%"=="status" goto show_status
if "%COMMAND%"=="validate" goto validate_migrations
if "%COMMAND%"=="help" goto usage
goto unknown_command

:migrate_up
echo Applying migrations...
if "%2"=="" (
    migrate -database "%DATABASE_URL%" -path %MIGRATIONS_DIR% up
) else (
    migrate -database "%DATABASE_URL%" -path %MIGRATIONS_DIR% up %2
)
if %ERRORLEVEL% equ 0 (
    echo [32m✓ Migrations applied successfully[0m
) else (
    echo [31m✗ Migration failed[0m
    exit /b 1
)
goto end

:migrate_down
echo Rolling back migrations...
if "%2"=="" (
    set /p confirm="Are you sure you want to rollback ALL migrations? (yes/no): "
    if not "!confirm!"=="yes" (
        echo Rollback cancelled
        goto end
    )
    migrate -database "%DATABASE_URL%" -path %MIGRATIONS_DIR% down
) else (
    migrate -database "%DATABASE_URL%" -path %MIGRATIONS_DIR% down %2
)
if %ERRORLEVEL% equ 0 (
    echo [32m✓ Migrations rolled back successfully[0m
) else (
    echo [31m✗ Rollback failed[0m
    exit /b 1
)
goto end

:create_migration
if "%2"=="" (
    echo Error: Migration name is required
    echo Usage: migrate.bat create migration_name
    exit /b 1
)
echo Creating migration: %2
migrate create -ext sql -dir %MIGRATIONS_DIR% -seq %2
if %ERRORLEVEL% equ 0 (
    echo [32m✓ Migration files created[0m
) else (
    echo [31m✗ Failed to create migration[0m
    exit /b 1
)
goto end

:show_version
echo Current migration version:
migrate -database "%DATABASE_URL%" -path %MIGRATIONS_DIR% version
goto end

:force_version
if "%2"=="" (
    echo Error: Version number is required
    echo Usage: migrate.bat force version_number
    exit /b 1
)
echo Forcing migration version to: %2
set /p confirm="Are you sure? This should only be used to fix dirty state (yes/no): "
if not "%confirm%"=="yes" (
    echo Force cancelled
    goto end
)
migrate -database "%DATABASE_URL%" -path %MIGRATIONS_DIR% force %2
if %ERRORLEVEL% equ 0 (
    echo [32m✓ Version forced to %2[0m
) else (
    echo [31m✗ Force failed[0m
    exit /b 1
)
goto end

:show_status
echo Migration Status:
echo Database URL: %DATABASE_URL%
echo Migrations Directory: %MIGRATIONS_DIR%
echo.
echo Migration files:
dir /b %MIGRATIONS_DIR%\*.up.sql 2>nul | find /c ".up.sql"
dir /b %MIGRATIONS_DIR%\*.down.sql 2>nul | find /c ".down.sql"
echo.
echo Current version:
migrate -database "%DATABASE_URL%" -path %MIGRATIONS_DIR% version
goto end

:validate_migrations
echo Validating migration files...
set VALID=1
for %%f in (%MIGRATIONS_DIR%\*.up.sql) do (
    set UP_FILE=%%f
    set BASE_NAME=%%~nf
    set DOWN_FILE=%MIGRATIONS_DIR%\!BASE_NAME:.up=.down!.sql
    if not exist "!DOWN_FILE!" (
        echo [31m✗ Missing down migration for: %%f[0m
        set VALID=0
    )
)
if %VALID% equ 1 (
    echo [32m✓ All migration files are valid[0m
) else (
    echo [31m✗ Validation failed[0m
    exit /b 1
)
goto end

:usage
echo Usage: migrate.bat [command] [options]
echo.
echo Commands:
echo   up [N]              Apply all or N pending migrations
echo   down [N]            Rollback all or N migrations
echo   create ^<name^>       Create a new migration file
echo   version             Show current migration version
echo   force ^<version^>     Force set migration version (use with caution)
echo   status              Show migration status
echo   validate            Validate migration files
echo   help                Show this help message
echo.
echo Environment Variables:
echo   DATABASE_URL        PostgreSQL connection string
echo                       Default: postgresql://postgres:postgres@localhost:5432/insavein?sslmode=disable
echo.
echo Examples:
echo   migrate.bat up                    # Apply all pending migrations
echo   migrate.bat up 1                  # Apply next migration
echo   migrate.bat down 1                # Rollback last migration
echo   migrate.bat create add_indexes    # Create new migration files
echo   migrate.bat version               # Show current version
echo   migrate.bat status                # Show migration status
goto end

:unknown_command
echo Error: Unknown command '%COMMAND%'
echo.
goto usage

:end
endlocal
