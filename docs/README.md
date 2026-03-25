# InSavein Platform Documentation

Complete documentation for the InSavein financial discipline platform.

## Documentation Overview

This directory contains comprehensive documentation for developers, operators, and API consumers.

### For Developers

- **[DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md)** - Complete guide for local development setup
  - Quick start with Docker Compose
  - Manual setup instructions
  - Project structure overview
  - Running tests
  - Code conventions and style guides
  - Adding new services
  - Database management
  - Debugging tips

### For API Consumers

- **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)** - Complete API reference
  - Authentication and authorization
  - All service endpoints with examples
  - Request/response formats
  - Error codes and handling
  - Rate limiting
  - Pagination

- **[openapi.yaml](./openapi.yaml)** - OpenAPI 3.0 specification
  - Machine-readable API specification
  - Import into Postman, Insomnia, or Swagger UI
  - Generate client SDKs

### For Operations

- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Production deployment guide
  - Kubernetes cluster setup (AWS EKS, GKE, AKS)
  - Environment variables and secrets management
  - Database deployment and migrations
  - Backend services deployment
  - Frontend deployment
  - Observability stack setup
  - Deployment process and rollback procedures

- **[OPERATIONS_RUNBOOK.md](./OPERATIONS_RUNBOOK.md)** - Day-to-day operations
  - Monitoring and alerting
  - Common issues and resolutions
  - Backup and restore procedures
  - Scaling procedures
  - Incident response
  - Maintenance tasks

## Quick Links

### Getting Started

1. **New Developer?** Start with [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md#quick-start)
2. **Deploying to Production?** See [DEPLOYMENT.md](./DEPLOYMENT.md#prerequisites)
3. **Integrating with API?** Check [API_DOCUMENTATION.md](./API_DOCUMENTATION.md#authentication)
4. **Troubleshooting Issues?** Refer to [OPERATIONS_RUNBOOK.md](./OPERATIONS_RUNBOOK.md#common-issues-and-resolutions)

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                             │
│                   TanStack Start (React)                     │
│                     Port: 3000                               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend Microservices                     │
├─────────────────────────────────────────────────────────────┤
│  Auth (8080)  │  User (8081)  │  Savings (8082)            │
│  Budget (8083) │  Goal (8005)  │  Education (8085)          │
│  Notification (8086) │  Analytics (8008)                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Data Layer                              │
├─────────────────────────────────────────────────────────────┤
│  PostgreSQL Primary (5432)                                   │
│  PostgreSQL Replica 1 (5433)                                 │
│  PostgreSQL Replica 2 (5434)                                 │
│  PgBouncer (6432) - Connection Pooling                       │
└─────────────────────────────────────────────────────────────┘
```

### Technology Stack

**Frontend**:
- TanStack Start (SSR framework)
- React 18+
- TypeScript
- TanStack Query (data fetching)
- Tailwind CSS

**Backend**:
- Go 1.21+
- PostgreSQL 15+
- JWT authentication
- RESTful APIs

**Infrastructure**:
- Kubernetes
- Docker
- Prometheus & Grafana
- GitHub Actions (CI/CD)

## Service Documentation

Each service has its own README with specific details:

- [auth-service/README.md](../auth-service/README.md) - Authentication service
- [user-service/README.md](../user-service/README.md) - User profile service
- [savings-service/README.md](../savings-service/README.md) - Savings tracking service
- [budget-service/README.md](../budget-service/README.md) - Budget planning service
- [goal-service/README.md](../goal-service/README.md) - Goal management service
- [education-service/README.md](../education-service/README.md) - Education content service
- [notification-service/README.md](../notification-service/README.md) - Notification service
- [analytics-service/README.md](../analytics-service/README.md) - Analytics service
- [frontend/README.md](../frontend/README.md) - Frontend application

## Additional Documentation

### Root Directory

- [README.md](../README.md) - Project overview and quick start
- [DOCKER_DEPLOYMENT.md](../DOCKER_DEPLOYMENT.md) - Docker deployment guide
- [DATABASE_SETUP.md](../DATABASE_SETUP.md) - Database setup and migrations
- [REPLICATION_SETUP.md](../REPLICATION_SETUP.md) - PostgreSQL replication
- [CI_CD_SETUP_GUIDE.md](../CI_CD_SETUP_GUIDE.md) - CI/CD pipeline setup
- [SECURITY_QUICK_REFERENCE.md](../SECURITY_QUICK_REFERENCE.md) - Security guidelines

### Kubernetes

- [k8s/README.md](../k8s/README.md) - Kubernetes configuration overview
- [k8s/DEPLOYMENT_GUIDE.md](../k8s/DEPLOYMENT_GUIDE.md) - K8s deployment steps

### Testing

- [integration-tests/README.md](../integration-tests/README.md) - Integration testing
- [performance-tests/README.md](../performance-tests/README.md) - Performance testing

## API Endpoints Summary

### Authentication (Port 8080)
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout user

### User Service (Port 8081)
- `GET /api/users/profile` - Get user profile
- `PUT /api/users/profile` - Update profile
- `GET /api/users/preferences` - Get preferences
- `PUT /api/users/preferences` - Update preferences

### Savings Service (Port 8082)
- `GET /api/savings/summary` - Get savings summary
- `GET /api/savings/history` - Get transaction history
- `POST /api/savings/transactions` - Create transaction
- `GET /api/savings/streak` - Get savings streak

### Budget Service (Port 8083)
- `GET /api/budgets/current` - Get current budget
- `POST /api/budgets` - Create budget
- `POST /api/budgets/spending` - Record spending
- `GET /api/budgets/alerts` - Get budget alerts

### Goal Service (Port 8005)
- `GET /api/goals` - Get all goals
- `POST /api/goals` - Create goal
- `POST /api/goals/{id}/contribute` - Add contribution
- `GET /api/goals/{id}/milestones` - Get milestones

### Education Service (Port 8085)
- `GET /api/education/lessons` - Get lessons
- `GET /api/education/lessons/{id}` - Get lesson detail
- `POST /api/education/lessons/{id}/complete` - Mark complete

### Notification Service (Port 8086)
- `GET /api/notifications` - Get notifications
- `PUT /api/notifications/{id}/read` - Mark as read

### Analytics Service (Port 8008)
- `GET /api/analytics/spending` - Get spending analysis
- `GET /api/analytics/financial-health` - Get health score
- `GET /api/analytics/recommendations` - Get recommendations

## Environment Variables

### Required for All Services

```bash
PORT=8080                          # Service port
DB_HOST=postgres-primary           # Database host
DB_PORT=5432                       # Database port
DB_USER=insavein_user              # Database user
DB_PASSWORD=<secret>               # Database password
DB_NAME=insavein                   # Database name
DB_SSLMODE=disable                 # SSL mode (require in production)
JWT_SECRET_KEY=<secret>            # JWT signing key
LOG_LEVEL=info                     # Log level (debug, info, error)
```

### Service-Specific

See individual service README files for additional environment variables.

## Common Commands

### Development

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f <service-name>

# Rebuild service
docker-compose build <service-name>
docker-compose up -d <service-name>

# Run tests
cd <service-name>
go test ./...              # Backend
npm test                   # Frontend

# Run migrations
cd migrations && ./migrate.sh
```

### Production

```bash
# Deploy to Kubernetes
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n insavein

# View logs
kubectl logs -l app=<service-name> -n insavein

# Scale service
kubectl scale deployment <service-name> --replicas=5 -n insavein

# Rollback deployment
kubectl rollout undo deployment/<service-name> -n insavein
```

## Troubleshooting

### Service Won't Start

1. Check logs: `docker-compose logs <service-name>`
2. Verify environment variables
3. Check database connectivity
4. Verify port availability

### Database Connection Issues

1. Check database is running: `docker-compose ps postgres-primary`
2. Test connection: `psql -h localhost -p 5432 -U insavein_user -d insavein`
3. Verify credentials in `.env` file
4. Check network connectivity

### API Errors

1. Check service logs
2. Verify authentication token
3. Check request format
4. Review API documentation

For more troubleshooting, see [OPERATIONS_RUNBOOK.md](./OPERATIONS_RUNBOOK.md#common-issues-and-resolutions).

## Performance Benchmarks

**Target SLAs**:
- API Response Time: p95 < 500ms, p99 < 1000ms
- Throughput: 100,000+ requests per minute
- Concurrent Users: 10,000+ (50,000+ peak)
- Uptime: 99.9% (< 43 minutes downtime per month)
- Error Rate: < 0.1%

**Database Performance**:
- Query Latency: p95 < 100ms
- Replication Lag: < 5 seconds
- Connection Pool: 20 connections per service

## Security

**Authentication**: JWT tokens (15-minute access, 7-day refresh)

**Encryption**:
- TLS 1.3 for all external traffic
- Database connections encrypted
- Passwords hashed with bcrypt (cost factor 12)

**Rate Limiting**:
- 100 requests/minute per user
- 1000 requests/minute per IP
- Burst allowance: 20 requests

**Security Headers**:
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- Content-Security-Policy
- Strict-Transport-Security

For more details, see [SECURITY_QUICK_REFERENCE.md](../SECURITY_QUICK_REFERENCE.md).

## Contributing

See [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md#contributing) for contribution guidelines.

## Support

**Documentation Issues**: Open an issue on GitHub

**Technical Support**: support@insavein.com

**Security Issues**: security@insavein.com

## License

MIT License - See LICENSE file for details

---

**Last Updated**: 2026-01-15  
**Documentation Version**: 1.0.0  
**Platform Version**: 1.0.0
