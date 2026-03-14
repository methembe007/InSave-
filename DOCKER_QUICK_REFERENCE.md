# Docker Quick Reference - InSavein Platform

## Essential Commands

### Starting Services

```bash
# Start all services
docker-compose up -d

# Start specific services
docker-compose up -d postgres-primary auth-service

# Start with rebuild
docker-compose up -d --build

# Start and view logs (foreground)
docker-compose up
```

### Stopping Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (⚠️ DELETES DATA)
docker-compose down -v

# Stop specific service
docker-compose stop auth-service
```

### Building Images

```bash
# Build all services
docker-compose build

# Build specific service
docker-compose build auth-service

# Build without cache
docker-compose build --no-cache

# Build and start
docker-compose up -d --build
```

### Viewing Logs

```bash
# All services (follow mode)
docker-compose logs -f

# Specific service
docker-compose logs -f auth-service

# Last 100 lines
docker-compose logs --tail=100 auth-service

# Multiple services
docker-compose logs -f auth-service user-service

# Since timestamp
docker-compose logs --since 2024-01-01T00:00:00
```

### Service Management

```bash
# Check service status
docker-compose ps

# Restart service
docker-compose restart auth-service

# Restart all services
docker-compose restart

# Scale service (requires load balancer)
docker-compose up -d --scale auth-service=3
```

### Debugging

```bash
# Execute command in running container
docker-compose exec auth-service sh

# Run one-off command
docker-compose run --rm auth-service sh

# View container details
docker inspect auth-service

# Check health status
docker inspect --format='{{json .State.Health}}' auth-service | jq
```

### Database Operations

```bash
# Connect to primary database
docker-compose exec postgres-primary psql -U postgres -d insavein

# Run SQL file
docker-compose exec postgres-primary psql -U postgres -d insavein -f /path/to/file.sql

# Backup database
docker-compose exec postgres-primary pg_dump -U postgres insavein > backup.sql

# Restore database
docker-compose exec -T postgres-primary psql -U postgres insavein < backup.sql

# Check replication status
docker-compose exec postgres-primary psql -U postgres -c "SELECT * FROM pg_stat_replication;"
```

### Cleanup

```bash
# Remove stopped containers
docker-compose rm

# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune

# Remove everything (⚠️ NUCLEAR OPTION)
docker system prune -a --volumes

# Check disk usage
docker system df
```

## Service URLs

| Service | URL | Health Check |
|---------|-----|--------------|
| Frontend | http://localhost:3000 | http://localhost:3000/ |
| Auth | http://localhost:8080 | http://localhost:8080/health |
| User | http://localhost:8081 | http://localhost:8081/health |
| Savings | http://localhost:8082 | http://localhost:8082/health |
| Budget | http://localhost:8083 | http://localhost:8083/health |
| Goal | http://localhost:8005 | http://localhost:8005/health |
| Education | http://localhost:8085 | http://localhost:8085/health |
| Notification | http://localhost:8086 | http://localhost:8086/health |
| Analytics | http://localhost:8008 | http://localhost:8008/health |
| PostgreSQL Primary | localhost:5432 | - |
| PostgreSQL Replica 1 | localhost:5433 | - |
| PostgreSQL Replica 2 | localhost:5434 | - |
| PgBouncer | localhost:6432 | - |

## Testing Commands

```bash
# Verify Dockerfiles
bash verify-dockerfiles.sh

# Run full test suite (Linux/Mac)
bash docker-build-test.sh

# Run full test suite (Windows)
docker-build-test.bat

# Test specific service health
curl http://localhost:8080/health

# Test user registration
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","first_name":"Test","last_name":"User","date_of_birth":"1990-01-01"}'

# Test user login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

## Troubleshooting

### Service won't start
```bash
# Check logs
docker-compose logs <service-name>

# Check if port is in use
netstat -an | grep <port>

# Restart service
docker-compose restart <service-name>

# Rebuild and restart
docker-compose up -d --build <service-name>
```

### Database connection issues
```bash
# Check PostgreSQL status
docker-compose ps postgres-primary

# Check PostgreSQL logs
docker-compose logs postgres-primary

# Test connection
docker-compose exec postgres-primary psql -U postgres -c "SELECT 1;"

# Restart database
docker-compose restart postgres-primary
```

### Health check failing
```bash
# Check service logs
docker-compose logs <service-name>

# Manually test health endpoint
curl http://localhost:<port>/health

# Check container health
docker inspect --format='{{json .State.Health}}' <service-name>

# Restart service
docker-compose restart <service-name>
```

### Out of memory
```bash
# Check resource usage
docker stats

# Add memory limits to docker-compose.yml
services:
  auth-service:
    mem_limit: 512m
    mem_reservation: 256m
```

### Disk space issues
```bash
# Check disk usage
docker system df

# Clean up
docker image prune -a
docker volume prune
docker system prune -a --volumes
```

## Environment Variables

### Common Variables
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

### Service-Specific

**Notification Service:**
```bash
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=notifications@insavein.com
SMTP_PASSWORD=smtp-password
```

**Analytics Service:**
```bash
CACHE_TTL=3600
```

**Frontend:**
```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_AUTH_SERVICE_URL=http://auth-service:8080
```

## Development Workflow

### Making Changes

1. Edit code in local directory
2. Rebuild: `docker-compose build <service-name>`
3. Restart: `docker-compose up -d <service-name>`
4. View logs: `docker-compose logs -f <service-name>`

### Running Tests

```bash
# Go service tests
docker-compose exec auth-service go test ./...

# Frontend tests
docker-compose exec frontend npm test
```

### Hot Reload (Development)

For development with hot reload, mount source code as volume:

```yaml
services:
  auth-service:
    volumes:
      - ./auth-service:/app
    command: go run cmd/server/main.go
```

## Production Checklist

- [ ] Change default passwords
- [ ] Use Docker secrets for sensitive data
- [ ] Enable SSL/TLS (DB_SSLMODE=require)
- [ ] Configure resource limits
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure centralized logging
- [ ] Implement automated backups
- [ ] Use container registry
- [ ] Set up CI/CD pipeline
- [ ] Configure auto-scaling
- [ ] Implement rate limiting
- [ ] Set up alerting
- [ ] Document runbooks
- [ ] Test disaster recovery

## Useful Docker Commands

```bash
# List all containers
docker ps -a

# List all images
docker images

# Remove container
docker rm <container-id>

# Remove image
docker rmi <image-id>

# View container logs
docker logs <container-id>

# Execute command in container
docker exec -it <container-id> sh

# Copy files from container
docker cp <container-id>:/path/to/file ./local/path

# Copy files to container
docker cp ./local/file <container-id>:/path/to/destination

# View container resource usage
docker stats

# View container processes
docker top <container-id>

# Inspect container
docker inspect <container-id>

# View container changes
docker diff <container-id>
```

## Network Commands

```bash
# List networks
docker network ls

# Inspect network
docker network inspect insavein-network

# Connect container to network
docker network connect insavein-network <container-id>

# Disconnect container from network
docker network disconnect insavein-network <container-id>
```

## Volume Commands

```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect postgres_primary_data

# Remove volume
docker volume rm postgres_primary_data

# Backup volume
docker run --rm -v postgres_primary_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz /data

# Restore volume
docker run --rm -v postgres_primary_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres_backup.tar.gz -C /
```

## Performance Tips

1. **Use .dockerignore** to exclude unnecessary files
2. **Layer caching**: Order Dockerfile commands from least to most frequently changing
3. **Multi-stage builds**: Separate build and runtime stages
4. **Minimize layers**: Combine RUN commands where possible
5. **Use specific tags**: Avoid `latest` tag in production
6. **Health checks**: Implement proper health checks for all services
7. **Resource limits**: Set CPU and memory limits
8. **Connection pooling**: Use PgBouncer for database connections

## Security Best Practices

1. **Non-root user**: Run containers as non-root
2. **Minimal base images**: Use Alpine Linux
3. **Scan images**: `docker scan <image-name>`
4. **Update regularly**: Keep base images and dependencies updated
5. **Secrets management**: Use Docker secrets or external secret managers
6. **Network isolation**: Use Docker networks
7. **Read-only filesystem**: Mount volumes as read-only where possible
8. **Drop capabilities**: Remove unnecessary Linux capabilities

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [DOCKER_DEPLOYMENT.md](./DOCKER_DEPLOYMENT.md) - Full deployment guide
