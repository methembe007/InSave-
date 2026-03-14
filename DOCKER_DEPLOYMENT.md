# Docker Deployment Guide

This guide explains how to build and run the InSavein platform using Docker and Docker Compose.

## Prerequisites

- Docker 20.10 or higher
- Docker Compose 2.0 or higher
- At least 8GB of available RAM
- At least 20GB of available disk space

## Architecture Overview

The InSavein platform consists of:
- **8 Go Microservices**: auth, user, savings, budget, goal, education, notification, analytics
- **1 TanStack Start Frontend**: React-based SSR application
- **PostgreSQL Database**: 1 primary + 2 read replicas with streaming replication
- **PgBouncer**: Connection pooler for database connections

## Quick Start

### 1. Build All Docker Images

```bash
# Build all services
docker-compose build

# Or build specific services
docker-compose build auth-service
docker-compose build frontend
```

### 2. Start All Services

```bash
# Start all services in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# View logs for specific service
docker-compose logs -f auth-service
```

### 3. Verify Services are Running

```bash
# Check service status
docker-compose ps

# Check health of all services
docker-compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
```

### 4. Run Database Migrations

```bash
# Wait for PostgreSQL to be ready (about 30 seconds)
sleep 30

# Run migrations
docker-compose exec postgres-primary psql -U postgres -d insavein -f /migrations/000001_create_users_table.up.sql
# ... run all migration files in order
```

Or use the migration script:
```bash
cd migrations
./migrate.sh
```

### 5. Access the Application

- **Frontend**: http://localhost:3000
- **Auth Service**: http://localhost:8080
- **User Service**: http://localhost:8081
- **Savings Service**: http://localhost:8082
- **Budget Service**: http://localhost:8083
- **Goal Service**: http://localhost:8005
- **Education Service**: http://localhost:8085
- **Notification Service**: http://localhost:8086
- **Analytics Service**: http://localhost:8008
- **pgAdmin** (optional): http://localhost:5050

## Service Details

### Microservices

All Go microservices follow the same pattern:
- **Multi-stage build**: golang:1.21-alpine (builder) → alpine:3.19 (runtime)
- **Non-root user**: Services run as `appuser` (UID 1000)
- **Health checks**: `/health`, `/health/live`, `/health/ready` endpoints
- **Security**: Minimal attack surface with Alpine Linux

#### Port Mapping

| Service | Internal Port | External Port |
|---------|--------------|---------------|
| Auth | 8080 | 8080 |
| User | 8081 | 8081 |
| Savings | 8082 | 8082 |
| Budget | 8083 | 8083 |
| Goal | 8005 | 8005 |
| Education | 8085 | 8085 |
| Notification | 8086 | 8086 |
| Analytics | 8008 | 8008 |

### Frontend

- **Build**: Node.js 20 Alpine
- **Runtime**: Node.js 20 Alpine with production dependencies only
- **Port**: 3000
- **Health check**: HTTP GET on `/`

### Database Architecture

#### PostgreSQL Primary (Write Operations)
- **Port**: 5432
- **Services**: auth, user, savings, budget, goal
- **Replication**: Streaming replication to 2 replicas

#### PostgreSQL Replica 1 (Read Operations)
- **Port**: 5433
- **Services**: education, notification
- **Purpose**: Reduce load on primary for read-heavy operations

#### PostgreSQL Replica 2 (Read Operations)
- **Port**: 5434
- **Services**: analytics
- **Purpose**: Isolate analytics queries from transactional workload

#### PgBouncer (Connection Pooler)
- **Port**: 6432
- **Pool Mode**: Transaction
- **Max Connections**: 1000
- **Default Pool Size**: 20 per service

## Environment Variables

### Common Variables (All Services)

```bash
PORT=<service-port>
DB_HOST=postgres-primary  # or postgres-replica1/postgres-replica2
DB_PORT=5432
DB_USER=insavein_user
DB_PASSWORD=insavein_password
DB_NAME=insavein
DB_SSLMODE=disable  # Use 'require' in production
JWT_SECRET=your-secret-key-change-in-production
```

### Service-Specific Variables

#### Notification Service
```bash
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=notifications@insavein.com
SMTP_PASSWORD=smtp-password
SMTP_FROM=notifications@insavein.com
```

#### Analytics Service
```bash
CACHE_TTL=3600  # Cache time-to-live in seconds
```

#### Frontend
```bash
NODE_ENV=production
VITE_API_BASE_URL=http://localhost:8080
VITE_AUTH_SERVICE_URL=http://auth-service:8080
# ... other service URLs
```

## Docker Compose Commands

### Starting Services

```bash
# Start all services
docker-compose up -d

# Start specific services
docker-compose up -d postgres-primary auth-service frontend

# Start with build
docker-compose up -d --build

# Start and view logs
docker-compose up
```

### Stopping Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes all data)
docker-compose down -v

# Stop specific service
docker-compose stop auth-service
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f auth-service

# Last 100 lines
docker-compose logs --tail=100 auth-service

# Since timestamp
docker-compose logs --since 2024-01-01T00:00:00 auth-service
```

### Scaling Services

```bash
# Scale a service to multiple replicas
docker-compose up -d --scale auth-service=3

# Note: You'll need to configure a load balancer for this to work properly
```

### Rebuilding Services

```bash
# Rebuild all services
docker-compose build

# Rebuild specific service
docker-compose build auth-service

# Rebuild without cache
docker-compose build --no-cache auth-service

# Rebuild and restart
docker-compose up -d --build auth-service
```

## Health Checks

All services include health checks that Docker monitors:

```bash
# Check health status
docker-compose ps

# Inspect health check details
docker inspect --format='{{json .State.Health}}' auth-service | jq
```

Health check endpoints:
- `/health` - Overall service health (includes database connectivity)
- `/health/live` - Liveness probe (service is running)
- `/health/ready` - Readiness probe (service is ready to accept traffic)

## Troubleshooting

### Service Won't Start

```bash
# Check logs
docker-compose logs <service-name>

# Check if port is already in use
netstat -an | grep <port>

# Restart service
docker-compose restart <service-name>
```

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker-compose ps postgres-primary

# Check PostgreSQL logs
docker-compose logs postgres-primary

# Test database connection
docker-compose exec postgres-primary psql -U postgres -d insavein -c "SELECT 1;"

# Check replication status
docker-compose exec postgres-primary psql -U postgres -c "SELECT * FROM pg_stat_replication;"
```

### Service Health Check Failing

```bash
# Check service logs
docker-compose logs <service-name>

# Manually test health endpoint
curl http://localhost:<port>/health

# Restart service
docker-compose restart <service-name>
```

### Out of Memory

```bash
# Check Docker resource usage
docker stats

# Increase Docker memory limit in Docker Desktop settings
# Or add memory limits to docker-compose.yml:
services:
  auth-service:
    mem_limit: 512m
    mem_reservation: 256m
```

### Disk Space Issues

```bash
# Check Docker disk usage
docker system df

# Clean up unused images
docker image prune -a

# Clean up unused volumes
docker volume prune

# Clean up everything (WARNING: removes all unused resources)
docker system prune -a --volumes
```

## Production Considerations

### Security

1. **Change default passwords** in docker-compose.yml
2. **Use secrets management** instead of environment variables
3. **Enable SSL/TLS** for database connections (DB_SSLMODE=require)
4. **Use HTTPS** for all external connections
5. **Rotate JWT secrets** regularly
6. **Scan images** for vulnerabilities: `docker scan <image-name>`

### Performance

1. **Use PgBouncer** for connection pooling (already configured)
2. **Configure resource limits** for each service
3. **Enable database query caching** where appropriate
4. **Use CDN** for frontend static assets
5. **Monitor resource usage** with Prometheus/Grafana

### High Availability

1. **Run multiple replicas** of each service
2. **Use external load balancer** (NGINX, HAProxy, or cloud LB)
3. **Configure automatic failover** for PostgreSQL
4. **Use external storage** for volumes (NFS, cloud storage)
5. **Implement circuit breakers** for service-to-service communication

### Monitoring

1. **Add Prometheus** for metrics collection
2. **Add Grafana** for visualization
3. **Configure alerting** for critical issues
4. **Use distributed tracing** (Jaeger, Zipkin)
5. **Centralized logging** (ELK stack, Loki)

## Development Workflow

### Making Changes

1. **Edit code** in your local directory
2. **Rebuild service**: `docker-compose build <service-name>`
3. **Restart service**: `docker-compose up -d <service-name>`
4. **View logs**: `docker-compose logs -f <service-name>`

### Running Tests

```bash
# Run tests in a service container
docker-compose exec auth-service go test ./...

# Run frontend tests
docker-compose exec frontend npm test
```

### Debugging

```bash
# Execute shell in running container
docker-compose exec auth-service sh

# Run a one-off command
docker-compose run --rm auth-service sh

# Attach to running container
docker attach auth-service
```

## Backup and Restore

### Database Backup

```bash
# Backup primary database
docker-compose exec postgres-primary pg_dump -U postgres insavein > backup.sql

# Backup with compression
docker-compose exec postgres-primary pg_dump -U postgres insavein | gzip > backup.sql.gz
```

### Database Restore

```bash
# Restore from backup
docker-compose exec -T postgres-primary psql -U postgres insavein < backup.sql

# Restore from compressed backup
gunzip -c backup.sql.gz | docker-compose exec -T postgres-primary psql -U postgres insavein
```

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Documentation](https://hub.docker.com/_/postgres)
- [Node.js Docker Best Practices](https://github.com/nodejs/docker-node/blob/main/docs/BestPractices.md)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)

## Support

For issues or questions:
1. Check the logs: `docker-compose logs <service-name>`
2. Review this documentation
3. Check the main README.md for project-specific information
4. Open an issue in the project repository
