# InSavein Platform

A comprehensive financial discipline application designed to help young people build savings habits, develop financial discipline, and work toward long-term financial independence.

## 🚀 Features

- **User Authentication**: Secure registration and login with JWT tokens
- **Savings Tracking**: Record and track savings transactions with streak counting
- **Budget Planning**: Create monthly budgets with category allocations
- **Spending Tracking**: Monitor spending against budget categories with alerts
- **Financial Goals**: Set and track progress toward long-term financial objectives
- **Education Content**: Access financial literacy lessons and resources
- **Notifications**: Receive email and push notifications for important events
- **Analytics**: Get insights into spending patterns and savings behavior
- **Financial Health Score**: Track overall financial discipline with calculated metrics

## 🏗️ Architecture

### Microservices
- **Auth Service** (Port 8080): Authentication and authorization
- **User Service** (Port 8081): User profile management
- **Savings Service** (Port 8082): Savings transaction tracking
- **Budget Service** (Port 8083): Budget planning and spending tracking
- **Goal Service** (Port 8005): Financial goal management
- **Education Service** (Port 8085): Financial education content delivery
- **Notification Service** (Port 8086): Email and push notifications
- **Analytics Service** (Port 8008): Financial analysis and recommendations

### Frontend
- **TanStack Start Application** (Port 3000): React-based SSR frontend

### Database
- **PostgreSQL Primary** (Port 5432): Write operations
- **PostgreSQL Replica 1** (Port 5433): Read operations (education, notification)
- **PostgreSQL Replica 2** (Port 5434): Read operations (analytics)
- **PgBouncer** (Port 6432): Connection pooling

## 🐳 Quick Start with Docker

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- 8GB RAM
- 20GB disk space

### Start All Services

```bash
# 1. Verify Dockerfiles
bash verify-dockerfiles.sh

# 2. Build and start all services
docker-compose up -d

# 3. Check service health
docker-compose ps

# 4. View logs
docker-compose logs -f
```

### Access the Application

- **Frontend**: http://localhost:3000
- **API Services**: http://localhost:8080-8086, 8008
- **Database**: localhost:5432 (primary), 5433 (replica1), 5434 (replica2)

### Run Tests

```bash
# Linux/Mac
bash docker-build-test.sh

# Windows
docker-build-test.bat
```

### Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove data (⚠️ WARNING: Deletes all data)
docker-compose down -v
```

## 📚 Documentation

- **[DOCKER_DEPLOYMENT.md](./DOCKER_DEPLOYMENT.md)**: Complete Docker deployment guide
- **[DOCKER_QUICK_REFERENCE.md](./DOCKER_QUICK_REFERENCE.md)**: Quick reference for Docker commands
- **[DATABASE_SETUP.md](./DATABASE_SETUP.md)**: Database setup and migration guide
- **[START_SERVICES.md](./START_SERVICES.md)**: Manual service startup guide
- **[REPLICATION_SETUP.md](./REPLICATION_SETUP.md)**: PostgreSQL replication configuration

## 🛠️ Development Setup

### Option 1: Docker (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f <service-name>

# Rebuild after code changes
docker-compose build <service-name>
docker-compose up -d <service-name>
```

### Option 2: Manual Setup

#### Prerequisites
- Go 1.21+
- Node.js 20+
- PostgreSQL 15+

#### Start Database

```bash
# Start PostgreSQL with replication
docker-compose up -d postgres-primary postgres-replica1 postgres-replica2

# Run migrations
cd migrations
./migrate.sh
```

#### Start Backend Services

```bash
# Auth Service
cd auth-service
go run cmd/server/main.go

# User Service
cd user-service
go run cmd/server/main.go

# Repeat for other services...
```

#### Start Frontend

```bash
cd frontend
npm install
npm run dev
```

## 🧪 Testing

### Run All Tests

```bash
# Backend tests
cd <service-name>
go test ./...

# Frontend tests
cd frontend
npm test
```

### Integration Tests

```bash
# With Docker
bash docker-build-test.sh

# Manual
cd user-service
./test-integration.sh
```

## 📦 Project Structure

```
.
├── auth-service/          # Authentication microservice
├── user-service/          # User profile microservice
├── savings-service/       # Savings tracking microservice
├── budget-service/        # Budget planning microservice
├── goal-service/          # Goal management microservice
├── education-service/     # Education content microservice
├── notification-service/  # Notification delivery microservice
├── analytics-service/     # Analytics and insights microservice
├── frontend/              # TanStack Start frontend
├── migrations/            # Database migrations
├── k8s/                   # Kubernetes deployment configs
├── postgres/              # PostgreSQL configuration
├── pgbouncer/             # PgBouncer configuration
├── monitoring/            # Monitoring scripts
├── docker-compose.yml     # Docker Compose configuration
└── README.md              # This file
```

## 🔧 Configuration

### Environment Variables

Each service uses environment variables for configuration. See `.env.example` files in each service directory.

**Common Variables:**
```bash
PORT=<service-port>
DB_HOST=postgres-primary
DB_PORT=5432
DB_USER=insavein_user
DB_PASSWORD=insavein_password
DB_NAME=insavein
DB_SSLMODE=disable
JWT_SECRET=your-secret-key
```

### Docker Compose

The `docker-compose.yml` file includes all services with proper configuration. Modify environment variables as needed for your setup.

## 🚀 Deployment

### Docker Compose (Development)

```bash
docker-compose up -d
```

### Kubernetes (Production)

```bash
cd k8s
kubectl apply -f namespace.yaml
kubectl apply -f secrets.yaml
kubectl apply -f configmap.yaml
kubectl apply -f .
```

See [k8s/README.md](./k8s/README.md) for detailed Kubernetes deployment instructions.

## 🔒 Security

- All services run as non-root users in Docker containers
- JWT-based authentication with token expiration
- Password hashing with bcrypt (cost factor 12)
- Rate limiting on authentication endpoints
- Database connection encryption (configurable)
- CORS policies for API protection
- Input validation and sanitization

## 📊 Monitoring

### Health Checks

All services expose health check endpoints:
- `/health` - Overall service health
- `/health/live` - Liveness probe
- `/health/ready` - Readiness probe

### Logs

```bash
# View all logs
docker-compose logs -f

# View specific service
docker-compose logs -f auth-service

# View last 100 lines
docker-compose logs --tail=100 auth-service
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

### Common Issues

**Docker not starting:**
- Ensure Docker Desktop is running
- Check Docker daemon status: `docker info`

**Database connection failed:**
- Check PostgreSQL is running: `docker-compose ps postgres-primary`
- Verify connection settings in environment variables
- Check logs: `docker-compose logs postgres-primary`

**Service health check failing:**
- Check service logs: `docker-compose logs <service-name>`
- Verify port is not in use: `netstat -an | grep <port>`
- Restart service: `docker-compose restart <service-name>`

### Getting Help

1. Check the documentation in the `docs/` directory
2. Review service-specific README files
3. Check Docker logs: `docker-compose logs -f`
4. Open an issue on GitHub

## 🗺️ Roadmap

- [ ] Mobile application (React Native)
- [ ] AI-powered savings recommendations
- [ ] Social features (savings challenges, leaderboards)
- [ ] Integration with banking APIs
- [ ] Advanced analytics dashboard
- [ ] Multi-currency support
- [ ] Automated savings rules
- [ ] Investment tracking

## 📈 Performance

- **API Response Time**: p95 < 500ms, p99 < 1000ms
- **Throughput**: 100,000+ requests per minute
- **Concurrent Users**: 10,000+ (50,000+ peak)
- **Database**: PostgreSQL with replication and connection pooling
- **Caching**: In-memory caching for analytics (1-hour TTL)

## 🏆 Acknowledgments

- TanStack for the excellent React framework
- Go community for robust backend libraries
- PostgreSQL for reliable database system
- Docker for containerization platform

---

**Built with ❤️ for financial empowerment**


## Git commit for professionals

feat: add login endpoint
fix: resolve password hashing bug
docs: update API documentation


feat(api): add user registration endpoint
fix(db): resolve connection pool exhaustion
refactor(auth): simplify token validation logic
docs: update installation instructions