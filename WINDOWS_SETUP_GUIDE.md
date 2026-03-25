# InSavein Platform - Windows Installation & Setup Guide

Complete step-by-step guide for setting up the InSavein Platform on Windows using Command Prompt.

## Table of Contents
- [Prerequisites Installation](#prerequisites-installation)
- [Quick Start (Docker - Recommended)](#quick-start-docker---recommended)
- [Manual Setup (Development)](#manual-setup-development)
- [Verification & Testing](#verification--testing)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites Installation

### 1. Install Docker Desktop (Recommended Method)

Docker Desktop includes everything you need to run the platform in containers.

**Download & Install:**
1. Visit: https://www.docker.com/products/docker-desktop/
2. Download Docker Desktop for Windows
3. Run the installer (requires admin privileges)
4. Restart your computer when prompted
5. Launch Docker Desktop from Start Menu
6. Wait for Docker to start (whale icon in system tray)

**Verify Installation:**
```cmd
docker --version
docker-compose --version
```

Expected output:
```
Docker version 24.0.0 or higher
Docker Compose version 2.0.0 or higher
```

**System Requirements:**
- Windows 10 64-bit: Pro, Enterprise, or Education (Build 19041 or higher)
- OR Windows 11 64-bit
- WSL 2 feature enabled (Docker Desktop will enable this)
- 8GB RAM minimum (16GB recommended)
- 20GB free disk space

---

### 2. Install Git for Windows

**Download & Install:**
1. Visit: https://git-scm.com/download/win
2. Download the 64-bit installer
3. Run installer with default settings
4. Select "Git Bash" as default terminal

**Verify Installation:**
```cmd
git --version
```

---

### 3. Install Go (For Manual Development Setup)

**Download & Install:**
1. Visit: https://go.dev/dl/
2. Download Windows installer (go1.21.x.windows-amd64.msi or higher)
3. Run installer with default settings
4. Restart Command Prompt after installation

**Verify Installation:**
```cmd
go version
```

Expected: `go version go1.21.0 or higher`

**Configure Go Environment (if needed):**
```cmd
REM Add to PATH permanently via System Properties > Environment Variables
REM Or set for current session:
set PATH=%PATH%;C:\Go\bin
set GOPATH=%USERPROFILE%\go
set PATH=%PATH%;%GOPATH%\bin
```

---

### 4. Install Node.js (For Frontend Development)

**Download & Install:**
1. Visit: https://nodejs.org/
2. Download LTS version (20.x or higher)
3. Run installer with default settings
4. Restart terminal after installation

**Verify Installation:**
```cmd
node --version
npm --version
```

Expected:
```
v20.0.0 or higher
10.0.0 or higher
```

---

### 5. Install PostgreSQL (For Manual Setup Only)

**Skip this if using Docker (recommended)**

**Download & Install:**
1. Visit: https://www.postgresql.org/download/windows/
2. Download PostgreSQL 15 installer
3. Run installer:
   - Set password for postgres user (remember this!)
   - Port: 5432 (default)
   - Locale: Default
4. Install Stack Builder components (optional)

**Verify Installation:**
```cmd
psql --version
```

**Add to PATH (if needed):**
```cmd
REM Add via System Properties > Environment Variables
REM Or set for current session:
set PATH=%PATH%;C:\Program Files\PostgreSQL\15\bin
```

---

### 6. Install golang-migrate (For Database Migrations)

**Using Scoop (Recommended):**

First, install Scoop if you don't have it:
```cmd
REM Open Command Prompt (not as Administrator)
REM Run this command to install Scoop:
powershell -Command "Set-ExecutionPolicy RemoteSigned -Scope CurrentUser; irm get.scoop.sh | iex"
```

Then install migrate:
```cmd
scoop install migrate
```

**Manual Installation:**
1. Visit: https://github.com/golang-migrate/migrate/releases
2. Download `migrate.windows-amd64.zip`
3. Extract to `C:\Program Files\migrate\`
4. Add to PATH via System Properties > Environment Variables
   - Add `C:\Program Files\migrate` to PATH

**Verify Installation:**
```cmd
migrate -version
```

---

### 7. Install Make (Optional, for convenience)

**Using Chocolatey:**

First, install Chocolatey (run Command Prompt as Administrator):
```cmd
@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "[System.Net.ServicePointManager]::SecurityProtocol = 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))" && SET "PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin"
```

Then install Make:
```cmd
choco install make
```

**Verify Installation:**
```cmd
make --version
```

---

## Quick Start (Docker - Recommended)

This is the fastest and easiest way to get started.

### Step 1: Clone the Repository

```cmd
REM Navigate to your projects folder
cd C:\Users\YourUsername\Projects

REM Clone the repository
git clone <repository-url> insavein-platform
cd insavein-platform
```

### Step 2: Start Docker Desktop

1. Launch Docker Desktop from Start Menu
2. Wait for Docker to be running (green icon in system tray)
3. Ensure WSL 2 is enabled in Docker Desktop settings

### Step 3: Start All Services

```cmd
REM Build and start all services
docker-compose up -d

REM This will:
REM - Build all Docker images (first time takes 5-10 minutes)
REM - Start PostgreSQL with replication
REM - Start all 8 microservices
REM - Start the frontend application
REM - Run database migrations automatically
```

### Step 4: Check Service Status

```cmd
REM View running containers
docker-compose ps

REM View logs
docker-compose logs -f

REM View specific service logs
docker-compose logs -f auth-service
```

### Step 5: Access the Application

Open your browser and navigate to:
- **Frontend**: http://localhost:3000
- **Auth Service**: http://localhost:8080/health
- **User Service**: http://localhost:8081/health
- **Savings Service**: http://localhost:8082/health
- **Budget Service**: http://localhost:8083/health
- **Goal Service**: http://localhost:8005/health
- **Education Service**: http://localhost:8085/health
- **Notification Service**: http://localhost:8086/health
- **Analytics Service**: http://localhost:8008/health

### Step 6: Create Your First User

1. Go to http://localhost:3000
2. Click "Register"
3. Fill in the registration form
4. Start using the platform!

### Managing Docker Services

```cmd
REM Stop all services
docker-compose down

REM Stop and remove all data (WARNING: Deletes database)
docker-compose down -v

REM Restart a specific service
docker-compose restart auth-service

REM Rebuild after code changes
docker-compose build auth-service
docker-compose up -d auth-service

REM View resource usage
docker stats
```

---

## Manual Setup (Development)

For developers who want to run services individually.

### Step 1: Clone Repository

```cmd
cd C:\Users\YourUsername\Projects
git clone <repository-url> insavein-platform
cd insavein-platform
```

### Step 2: Setup PostgreSQL Database

**Option A: Using Docker (Easier)**
```cmd
REM Start only PostgreSQL
docker-compose up -d postgres-primary

REM Wait for it to be ready
docker-compose ps
```

**Option B: Using Local PostgreSQL**
```cmd
REM Start PostgreSQL service (run as Administrator)
net start postgresql-x64-15

REM Create database and user
psql -U postgres
```

In PostgreSQL shell:
```sql
CREATE DATABASE insavein;
CREATE USER insavein_user WITH PASSWORD 'insavein_password';
GRANT ALL PRIVILEGES ON DATABASE insavein TO insavein_user;
\q
```

### Step 3: Run Database Migrations

```cmd
REM Set database connection string
set DATABASE_URL=postgresql://insavein_user:insavein_password@localhost:5432/insavein?sslmode=disable

REM Navigate to migrations folder
cd migrations

REM Run migrations (Windows batch file)
migrate.bat up

REM Or use migrate directly
migrate -database %DATABASE_URL% -path . up

REM Verify migrations
migrate -database %DATABASE_URL% -path . version
```

Expected output: `10` (current version)

### Step 4: Configure Environment Variables

Create `.env` files for each service:

**auth-service/.env:**
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=insavein_user
DB_PASSWORD=insavein_password
DB_NAME=insavein
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
```

**user-service/.env:**
```env
PORT=8081
DB_HOST=localhost
DB_PORT=5432
DB_USER=insavein_user
DB_PASSWORD=insavein_password
DB_NAME=insavein
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-change-in-production
```

**Repeat for other services** (savings, budget, goal, education, notification, analytics)
- Use ports: 8082, 8083, 8005, 8085, 8086, 8008 respectively
- Same database configuration
- Same JWT_SECRET across all services

### Step 5: Start Backend Services

Open separate Command Prompt windows for each service:

**Terminal 1 - Auth Service:**
```cmd
cd auth-service
go mod download
go run cmd/server/main.go
```

**Terminal 2 - User Service:**
```cmd
cd user-service
go mod download
go run cmd/server/main.go
```

**Terminal 3 - Savings Service:**
```cmd
cd savings-service
go mod download
go run cmd/server/main.go
```

**Terminal 4 - Budget Service:**
```cmd
cd budget-service
go mod download
go run cmd/server/main.go
```

**Terminal 5 - Goal Service:**
```cmd
cd goal-service
go mod download
go run cmd/server/main.go
```

**Terminal 6 - Education Service:**
```cmd
cd education-service
go mod download
go run cmd/server/main.go
```

**Terminal 7 - Notification Service:**
```cmd
cd notification-service
go mod download
go run cmd/server/main.go
```

**Terminal 8 - Analytics Service:**
```cmd
cd analytics-service
go mod download
go run cmd/server/main.go
```

### Step 6: Start Frontend

**Terminal 9 - Frontend:**
```cmd
cd frontend

REM Install dependencies (first time only)
npm install

REM Start development server
npm run dev
```

The frontend will be available at: http://localhost:3000

---

## Verification & Testing

### 1. Check Service Health

```cmd
REM Test each service health endpoint
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
curl http://localhost:8005/health
curl http://localhost:8085/health
curl http://localhost:8086/health
curl http://localhost:8008/health
```

Expected response for each:
```json
{"status":"healthy","service":"<service-name>"}
```

### 2. Test User Registration

```cmd
REM Register a test user
curl -X POST http://localhost:8080/api/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@example.com\",\"password\":\"SecurePass123!\",\"firstName\":\"Test\",\"lastName\":\"User\",\"dateOfBirth\":\"1990-01-01\"}"
```

### 3. Test User Login

```cmd
REM Login with test user
curl -X POST http://localhost:8080/api/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@example.com\",\"password\":\"SecurePass123!\"}"
```

You should receive a JWT token in the response.

### 4. Verify Database

```cmd
REM Connect to database
psql -U insavein_user -d insavein

REM In psql shell:
REM Check tables
\dt

REM Check users
SELECT id, email, first_name, last_name FROM users;

REM Exit
\q
```

### 5. Run Tests

**Backend Tests:**
```cmd
REM Test auth service
cd auth-service
go test ./...

REM Test user service
cd user-service
go test ./...
```

**Frontend Tests:**
```cmd
cd frontend
npm test
```

**Docker Build Tests:**
```cmd
REM Test all Docker builds
docker-build-test.bat
```

---

## Troubleshooting

### Docker Issues

**Problem: Docker Desktop won't start**
- Solution: Enable WSL 2 in Windows Features
  1. Open Command Prompt as Administrator
  2. Run: `wsl --install`
  3. Restart computer
  4. Launch Docker Desktop again

**Problem: "docker-compose: command not found"**
- Solution: Use `docker compose` (without hyphen) for newer Docker versions
  ```cmd
  docker compose up -d
  ```

**Problem: Port already in use**
- Solution: Find and kill the process
  ```cmd
  REM Find process using port 8080
  netstat -ano | findstr "8080"
  
  REM Kill process (replace PID with actual process ID)
  taskkill /PID <PID> /F
  ```

### Database Issues

**Problem: "Failed to connect to database"**
- Check PostgreSQL is running:
  ```cmd
  REM Docker
  docker-compose ps postgres-primary
  
  REM Local service (run as Administrator)
  sc query postgresql-x64-15
  ```
- Verify credentials in `.env` files
- Check firewall isn't blocking port 5432

**Problem: "Migration failed"**
- Check migration version:
  ```cmd
  migrate -database %DATABASE_URL% -path migrations version
  ```
- Force to specific version if needed:
  ```cmd
  migrate -database %DATABASE_URL% -path migrations force 10
  ```
- Re-run migrations:
  ```cmd
  migrate -database %DATABASE_URL% -path migrations up
  ```

### Service Issues

**Problem: "Failed to fetch" in frontend**
- Ensure backend services are running
- Check service health endpoints
- Verify CORS configuration
- Check browser console for errors

**Problem: Service won't start - "port already in use"**
- Find what's using the port:
  ```cmd
  netstat -ano | findstr "8080"
  ```
- Kill the process or change port in `.env`

**Problem: "Cannot find module" in Go**
- Download dependencies:
  ```cmd
  go mod download
  go mod tidy
  ```

**Problem: Frontend build errors**
- Clear node_modules and reinstall:
  ```cmd
  cd frontend
  rmdir /s /q node_modules
  del package-lock.json
  npm install
  ```

### Performance Issues

**Problem: Docker using too much memory**
- Adjust Docker Desktop settings:
  1. Open Docker Desktop
  2. Settings → Resources
  3. Reduce Memory limit to 4GB
  4. Apply & Restart

**Problem: Slow database queries**
- Check database connections:
  ```sql
  SELECT count(*) FROM pg_stat_activity;
  ```
- Restart services to clear connections
- Consider using PgBouncer for connection pooling

### Windows-Specific Issues

**Problem: Line ending issues (CRLF vs LF)**
- Configure Git to handle line endings:
  ```cmd
  git config --global core.autocrlf true
  ```

**Problem: Permission denied on scripts**
- Run Command Prompt as Administrator
- Right-click Command Prompt → Run as administrator

**Problem: Path too long errors**
- Enable long paths in Windows (run as Administrator):
  ```cmd
  reg add HKLM\SYSTEM\CurrentControlSet\Control\FileSystem /v LongPathsEnabled /t REG_DWORD /d 1 /f
  ```
- Restart computer after making this change

**Problem: Environment variables not persisting**
- Set permanently via System Properties:
  1. Right-click "This PC" → Properties
  2. Advanced system settings → Environment Variables
  3. Add variables under "User variables" or "System variables"

---

## Quick Reference Commands

### Docker Commands
```cmd
REM Start all services
docker-compose up -d

REM Stop all services
docker-compose down

REM View logs
docker-compose logs -f

REM Restart service
docker-compose restart <service-name>

REM Rebuild service
docker-compose build <service-name>

REM Remove all containers and volumes
docker-compose down -v
```

### Database Commands
```cmd
REM Connect to database
psql -U insavein_user -d insavein

REM Run migrations
cd migrations
migrate.bat up

REM Check migration version
migrate -database %DATABASE_URL% -path migrations version

REM Rollback migration
migrate.bat down 1
```

### Service Commands
```cmd
REM Start service
cd <service-name>
go run cmd/server/main.go

REM Run tests
go test ./...

REM Build binary
go build -o service.exe cmd/server/main.go
```

### Frontend Commands
```cmd
REM Install dependencies
npm install

REM Start dev server
npm run dev

REM Run tests
npm test

REM Build for production
npm run build
```

---

## Next Steps

After successful setup:

1. **Explore the API**: Check `<service>/API_EXAMPLES.md` files
2. **Read Documentation**: Review `README.md` and service-specific docs
3. **Configure Production**: Update `.env` files with production values
4. **Set up Monitoring**: Configure Prometheus and Grafana (see k8s/)
5. **Deploy to Cloud**: Follow Kubernetes deployment guide (k8s/DEPLOYMENT_GUIDE.md)

---

## Additional Resources

- **Main README**: [README.md](./README.md)
- **Docker Deployment**: [DOCKER_DEPLOYMENT.md](./DOCKER_DEPLOYMENT.md)
- **Database Setup**: [DATABASE_SETUP.md](./DATABASE_SETUP.md)
- **Kubernetes Deployment**: [k8s/DEPLOYMENT_GUIDE.md](./k8s/DEPLOYMENT_GUIDE.md)
- **Security Guide**: [SECURITY_QUICK_REFERENCE.md](./SECURITY_QUICK_REFERENCE.md)

---

## Support

If you encounter issues not covered in this guide:

1. Check service logs: `docker-compose logs <service-name>`
2. Review service README files
3. Check GitHub issues
4. Contact the development team

---

**Happy Coding! 🚀**
