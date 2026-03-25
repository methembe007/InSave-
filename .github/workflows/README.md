# CI/CD Pipeline Documentation

This directory contains GitHub Actions workflows for the InSavein Platform CI/CD pipeline.

## Workflows Overview

### 1. Lint (`lint.yml`)
**Triggers:** Push and PR to `main` and `develop` branches

**Purpose:** Ensures code quality and consistency across all services

**Jobs:**
- `lint-go`: Runs golangci-lint on all 8 Go microservices
- `lint-frontend`: Runs ESLint and TypeScript type checking on the frontend
- `lint-shared`: Runs golangci-lint on shared middleware

**Features:**
- Matrix strategy for parallel linting of all services
- Caching for faster execution
- Strict linting rules with zero warnings policy

### 2. Test (`test.yml`)
**Triggers:** Push and PR to `main` and `develop` branches

**Purpose:** Runs unit tests and integration tests with coverage reporting

**Jobs:**
- `test-go`: Runs Go tests for all microservices with PostgreSQL service
- `test-frontend`: Runs Vitest tests for the frontend
- `test-integration`: Runs integration tests across services

**Features:**
- PostgreSQL service container for database tests
- Code coverage reporting to Codecov
- Race condition detection with `-race` flag
- Database migrations before tests

### 3. Security Scanning (`security.yml`)
**Triggers:** 
- Push and PR to `main` and `develop` branches
- Scheduled daily at 2 AM UTC

**Purpose:** Identifies security vulnerabilities in code and dependencies

**Jobs:**
- `trivy-filesystem`: Scans filesystem for vulnerabilities
- `trivy-config`: Scans configuration files (Kubernetes, Docker)
- `snyk-go`: Scans Go dependencies for known vulnerabilities
- `snyk-frontend`: Scans npm dependencies for known vulnerabilities
- `gosec`: Static security analysis for Go code
- `codeql`: Advanced semantic code analysis

**Features:**
- SARIF format results uploaded to GitHub Security tab
- Multiple security tools for comprehensive coverage
- Scheduled daily scans for continuous monitoring

### 4. Build and Push (`build-push.yml`)
**Triggers:** 
- Push to `main` and `develop` branches
- Push of version tags (`v*`)
- Pull requests

**Purpose:** Builds Docker images and pushes to GitHub Container Registry

**Jobs:**
- `build-go-services`: Builds all 8 Go microservice images
- `build-frontend`: Builds the frontend image
- `verify-images`: Verifies built images are pullable
- `scan-images`: Scans built images for vulnerabilities

**Features:**
- Multi-platform builds (linux/amd64, linux/arm64)
- Smart tagging strategy (branch, SHA, semver, latest)
- Build caching for faster builds
- Automatic vulnerability scanning of built images
- Only pushes on non-PR events

**Image Tags:**
- `develop-<sha>`: Develop branch builds
- `main-<sha>`: Main branch builds
- `v1.2.3`: Semantic version tags
- `latest`: Latest main branch build

### 5. Deploy to Staging (`deploy-staging.yml`)
**Triggers:** 
- Push to `develop` branch
- Manual workflow dispatch

**Purpose:** Deploys the application to the staging environment

**Environment:** `staging` (https://staging.insavein.com)

**Jobs:**
- `deploy-staging`: Deploys all services to Kubernetes staging namespace
- `notify-deployment`: Sends Slack notification about deployment status

**Features:**
- Automatic deployment on develop branch push
- Database migration execution
- Rolling updates with zero downtime
- Health check verification for all services
- Comprehensive smoke tests
- Slack notifications

**Deployment Steps:**
1. Configure kubectl with staging cluster
2. Create/update ConfigMaps and Secrets
3. Deploy PostgreSQL (if needed)
4. Run database migrations
5. Deploy all 8 backend services
6. Deploy frontend
7. Apply ingress configuration
8. Wait for rollout completion
9. Run smoke tests
10. Send notification

### 6. Deploy to Production (`deploy-production.yml`)
**Triggers:** 
- Push to `main` branch
- Push of version tags (`v*`)
- Manual workflow dispatch (requires approval)

**Purpose:** Deploys the application to the production environment

**Environment:** `production` (https://insavein.com)

**Jobs:**
- `deploy-production`: Deploys all services to Kubernetes production namespace
- `notify-deployment`: Sends Slack and email notifications

**Features:**
- Manual approval required (GitHub environment protection)
- Backup of current deployment state
- Rolling updates with maxUnavailable=0 (zero downtime)
- Comprehensive smoke tests including critical user flows
- 2-minute health monitoring after deployment
- Automatic rollback on failure
- Slack and email notifications
- Enhanced error handling

**Deployment Steps:**
1. Configure kubectl with production cluster
2. Backup current deployment state
3. Create/update ConfigMaps and Secrets
4. Deploy PostgreSQL (if needed)
5. Run database migrations
6. Deploy all services with rolling update strategy
7. Wait for rollout completion (10min timeout per service)
8. Run comprehensive smoke tests
9. Monitor deployment health for 2 minutes
10. Send notifications
11. Automatic rollback if any step fails

## Required Secrets

### GitHub Secrets
Configure these in your repository settings:

#### Staging Environment
- `KUBE_CONFIG_STAGING`: Base64-encoded kubeconfig for staging cluster
- `JWT_SECRET_STAGING`: JWT signing secret for staging
- `DATABASE_URL_STAGING`: PostgreSQL connection string for staging

#### Production Environment
- `KUBE_CONFIG_PRODUCTION`: Base64-encoded kubeconfig for production cluster
- `JWT_SECRET_PRODUCTION`: JWT signing secret for production
- `DATABASE_URL_PRODUCTION`: PostgreSQL connection string for production

#### Shared Secrets
- `CODECOV_TOKEN`: Token for uploading coverage reports
- `SNYK_TOKEN`: Token for Snyk security scanning
- `SMTP_USERNAME`: SMTP username for email notifications
- `SMTP_PASSWORD`: SMTP password for email notifications
- `SLACK_WEBHOOK_URL`: Slack webhook URL for notifications
- `ALERT_EMAIL`: Email address for critical alerts

#### Automatic Secrets
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions

## Environment Protection Rules

### Staging Environment
- No approval required
- Automatic deployment on develop branch push

### Production Environment
- **Manual approval required** from designated reviewers
- Deployment only from main branch or version tags
- Additional security checks

## Workflow Dependencies

```
┌─────────────┐
│   Lint      │
└─────────────┘
       ↓
┌─────────────┐
│   Test      │
└─────────────┘
       ↓
┌─────────────┐
│  Security   │
└─────────────┘
       ↓
┌─────────────┐
│ Build/Push  │
└─────────────┘
       ↓
┌─────────────┐     ┌──────────────┐
│   Staging   │ →   │  Production  │
└─────────────┘     └──────────────┘
```

## Best Practices

### Branch Strategy
- `develop`: Development branch, auto-deploys to staging
- `main`: Production branch, auto-deploys to production (with approval)
- Feature branches: Create PRs to develop
- Release tags: `v1.2.3` format for versioned releases

### Deployment Strategy
1. Merge feature to `develop` → Auto-deploy to staging
2. Test in staging environment
3. Merge `develop` to `main` → Requires approval → Deploy to production
4. Tag release with semantic version

### Rollback Procedure
If production deployment fails:
1. Automatic rollback is triggered
2. Previous version is restored
3. Notifications are sent
4. Manual investigation required

For manual rollback:
```bash
# Rollback specific service
kubectl rollout undo deployment/auth-service -n insavein-production

# Rollback to specific revision
kubectl rollout undo deployment/auth-service --to-revision=2 -n insavein-production
```

### Monitoring Deployments
```bash
# Watch deployment status
kubectl rollout status deployment/auth-service -n insavein-production

# View deployment history
kubectl rollout history deployment/auth-service -n insavein-production

# View pod logs
kubectl logs -f deployment/auth-service -n insavein-production
```

## Troubleshooting

### Build Failures
- Check Go version compatibility (1.25+)
- Verify all dependencies are available
- Check Docker build context

### Test Failures
- Ensure PostgreSQL service is healthy
- Check database migration status
- Verify environment variables

### Deployment Failures
- Verify kubeconfig is valid and not expired
- Check cluster connectivity
- Ensure namespace exists
- Verify secrets are configured
- Check pod logs for errors

### Security Scan Failures
- Review SARIF results in GitHub Security tab
- Update vulnerable dependencies
- Fix identified security issues

## Performance Optimization

### Caching Strategy
- Go module cache: Speeds up dependency downloads
- golangci-lint cache: Speeds up linting
- Docker layer cache: Speeds up image builds
- npm cache: Speeds up frontend builds

### Parallel Execution
- Matrix strategy for parallel service builds
- Concurrent linting of all services
- Parallel test execution

## Maintenance

### Regular Tasks
- Update action versions quarterly
- Review and update security scanning rules
- Rotate secrets every 90 days
- Review and optimize workflow performance
- Update Go and Node.js versions

### Monitoring
- Check GitHub Actions usage and billing
- Monitor workflow execution times
- Review security scan results weekly
- Track deployment success rates

## Support

For issues or questions about the CI/CD pipeline:
1. Check workflow logs in GitHub Actions tab
2. Review this documentation
3. Contact DevOps team
4. Create an issue in the repository
