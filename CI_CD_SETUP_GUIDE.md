# CI/CD Pipeline Setup Guide

This guide walks you through setting up the complete CI/CD pipeline for the InSavein Platform.

## Prerequisites

Before setting up the CI/CD pipeline, ensure you have:

1. **GitHub Repository**: Repository with admin access
2. **Kubernetes Clusters**: Staging and production clusters
3. **Container Registry**: GitHub Container Registry (GHCR) access
4. **External Services**:
   - Codecov account for coverage reporting
   - Snyk account for security scanning
   - Slack workspace for notifications
   - SMTP server for email notifications

## Step 1: Configure GitHub Secrets

Navigate to your repository settings → Secrets and variables → Actions → New repository secret

### Required Secrets

#### Kubernetes Configuration

**KUBE_CONFIG_STAGING**
```bash
# Get your kubeconfig and encode it
cat ~/.kube/config-staging | base64 -w 0
```

**KUBE_CONFIG_PRODUCTION**
```bash
# Get your kubeconfig and encode it
cat ~/.kube/config-production | base64 -w 0
```

#### Application Secrets

**JWT_SECRET_STAGING**
```bash
# Generate a secure random secret
openssl rand -base64 32
```

**JWT_SECRET_PRODUCTION**
```bash
# Generate a different secure random secret
openssl rand -base64 32
```

**DATABASE_URL_STAGING**
```
postgres://username:password@staging-db-host:5432/insavein?sslmode=require
```

**DATABASE_URL_PRODUCTION**
```
postgres://username:password@production-db-host:5432/insavein?sslmode=require
```

#### Third-Party Service Tokens

**CODECOV_TOKEN**
1. Go to https://codecov.io
2. Add your repository
3. Copy the upload token

**SNYK_TOKEN**
1. Go to https://snyk.io
2. Navigate to Account Settings → API Token
3. Copy your token

**SLACK_WEBHOOK_URL**
1. Go to https://api.slack.com/apps
2. Create a new app or select existing
3. Enable Incoming Webhooks
4. Create a webhook for your channel
5. Copy the webhook URL

**SMTP_USERNAME** and **SMTP_PASSWORD**
- Your SMTP server credentials
- For Gmail: Use App Password (not regular password)

**ALERT_EMAIL**
- Email address to receive critical alerts
- Example: devops@insavein.com

## Step 2: Configure GitHub Environments

### Create Staging Environment

1. Go to Settings → Environments → New environment
2. Name: `staging`
3. Environment URL: `https://staging.insavein.com`
4. No protection rules needed (auto-deploy)

### Create Production Environment

1. Go to Settings → Environments → New environment
2. Name: `production`
3. Environment URL: `https://insavein.com`
4. Configure protection rules:
   - ✅ Required reviewers (add team members)
   - ✅ Wait timer: 0 minutes (or set delay if needed)
   - ✅ Deployment branches: Only `main` branch

## Step 3: Set Up Kubernetes Clusters

### Staging Cluster Setup

```bash
# Create namespace
kubectl create namespace insavein-staging

# Create secrets
kubectl create secret generic app-secrets \
  --from-literal=jwt-secret="YOUR_JWT_SECRET_STAGING" \
  --from-literal=database-url="YOUR_DATABASE_URL_STAGING" \
  --from-literal=smtp-password="YOUR_SMTP_PASSWORD" \
  --namespace=insavein-staging

# Apply resource quotas and limits
kubectl apply -f k8s/resource-quota.yaml -n insavein-staging
kubectl apply -f k8s/network-policy.yaml -n insavein-staging
```

### Production Cluster Setup

```bash
# Create namespace
kubectl create namespace insavein-production

# Create secrets
kubectl create secret generic app-secrets \
  --from-literal=jwt-secret="YOUR_JWT_SECRET_PRODUCTION" \
  --from-literal=database-url="YOUR_DATABASE_URL_PRODUCTION" \
  --from-literal=smtp-password="YOUR_SMTP_PASSWORD" \
  --namespace=insavein-production

# Apply resource quotas and limits
kubectl apply -f k8s/resource-quota.yaml -n insavein-production
kubectl apply -f k8s/network-policy.yaml -n insavein-production

# Set up monitoring
kubectl apply -f k8s/prometheus-deployment.yaml -n insavein-production
kubectl apply -f k8s/grafana-deployment.yaml -n insavein-production
```

## Step 4: Configure Branch Protection Rules

### Main Branch Protection

1. Go to Settings → Branches → Add rule
2. Branch name pattern: `main`
3. Configure rules:
   - ✅ Require a pull request before merging
   - ✅ Require approvals: 2
   - ✅ Dismiss stale pull request approvals
   - ✅ Require status checks to pass before merging
     - Select: `Lint`, `Test`, `Security Scanning`
   - ✅ Require branches to be up to date before merging
   - ✅ Require conversation resolution before merging
   - ✅ Include administrators

### Develop Branch Protection

1. Go to Settings → Branches → Add rule
2. Branch name pattern: `develop`
3. Configure rules:
   - ✅ Require a pull request before merging
   - ✅ Require approvals: 1
   - ✅ Require status checks to pass before merging
     - Select: `Lint`, `Test`
   - ✅ Require branches to be up to date before merging

## Step 5: Install Required Dependencies

### Frontend Dependencies

```bash
cd frontend
npm install --save-dev \
  eslint \
  @typescript-eslint/parser \
  @typescript-eslint/eslint-plugin \
  eslint-plugin-react \
  eslint-plugin-react-hooks \
  eslint-plugin-jsx-a11y
```

### Go Services

Each service should have a `go.mod` file. Ensure all dependencies are up to date:

```bash
# For each service
cd auth-service
go mod tidy
go mod verify
```

## Step 6: Test the Pipeline

### Test Linting

```bash
# Create a feature branch
git checkout -b feature/test-pipeline

# Make a small change
echo "# Test" >> README.md

# Commit and push
git add .
git commit -m "test: CI/CD pipeline"
git push origin feature/test-pipeline

# Create a pull request
# The lint workflow should run automatically
```

### Test Building

```bash
# Merge to develop branch
git checkout develop
git merge feature/test-pipeline
git push origin develop

# The build-push workflow should run
# Check GitHub Actions tab for progress
```

### Test Staging Deployment

```bash
# After build completes, staging deployment should start automatically
# Monitor in GitHub Actions tab

# Verify deployment
kubectl get pods -n insavein-staging
kubectl get deployments -n insavein-staging

# Test staging URL
curl https://staging.insavein.com/api/auth/health
```

### Test Production Deployment

```bash
# Create a pull request from develop to main
git checkout main
git pull origin main
git merge develop
git push origin main

# Approve the deployment in GitHub
# Go to Actions → Deploy to Production → Review deployments → Approve

# Verify deployment
kubectl get pods -n insavein-production
kubectl get deployments -n insavein-production

# Test production URL
curl https://insavein.com/api/auth/health
```

## Step 7: Configure Monitoring and Alerts

### Prometheus Alerts

```bash
# Apply Prometheus alert rules
kubectl apply -f k8s/prometheus-alerts.yaml -n insavein-production
```

### Grafana Dashboards

```bash
# Apply Grafana dashboards
kubectl apply -f k8s/grafana-dashboards.yaml -n insavein-production
```

### Access Grafana

```bash
# Port forward to access Grafana
kubectl port-forward -n insavein-production svc/grafana 3000:3000

# Open http://localhost:3000
# Default credentials: admin/admin (change immediately)
```

## Step 8: Verify Security Scanning

### Check Security Tab

1. Go to Security → Code scanning alerts
2. Verify Trivy, Snyk, and CodeQL results appear
3. Review and triage any findings

### Check Dependabot

1. Go to Security → Dependabot alerts
2. Enable Dependabot security updates
3. Configure auto-merge for patch updates

## Troubleshooting

### Workflow Fails to Start

**Issue**: Workflow doesn't trigger on push

**Solution**:
- Check branch protection rules
- Verify workflow file syntax: `yamllint .github/workflows/*.yml`
- Check GitHub Actions is enabled: Settings → Actions → General

### Authentication Errors

**Issue**: Cannot push to GHCR or access Kubernetes

**Solution**:
```bash
# Verify GITHUB_TOKEN has packages:write permission
# Check Settings → Actions → General → Workflow permissions

# Verify kubeconfig is valid
echo "$KUBE_CONFIG_STAGING" | base64 -d > /tmp/kubeconfig
kubectl --kubeconfig=/tmp/kubeconfig get nodes
```

### Build Failures

**Issue**: Docker build fails

**Solution**:
```bash
# Test build locally
docker build -t test-image ./auth-service

# Check Dockerfile syntax
docker build --no-cache -t test-image ./auth-service
```

### Deployment Failures

**Issue**: Kubernetes deployment fails

**Solution**:
```bash
# Check pod status
kubectl get pods -n insavein-staging
kubectl describe pod <pod-name> -n insavein-staging
kubectl logs <pod-name> -n insavein-staging

# Check deployment events
kubectl get events -n insavein-staging --sort-by='.lastTimestamp'

# Verify secrets exist
kubectl get secrets -n insavein-staging
```

### Test Failures

**Issue**: Tests fail in CI but pass locally

**Solution**:
- Check PostgreSQL service is healthy
- Verify environment variables are set
- Check for race conditions: `go test -race ./...`
- Review test logs in GitHub Actions

## Maintenance

### Weekly Tasks

- Review security scan results
- Check deployment success rates
- Monitor workflow execution times
- Review and merge Dependabot PRs

### Monthly Tasks

- Update action versions
- Review and optimize caching strategy
- Audit access permissions
- Test disaster recovery procedures

### Quarterly Tasks

- Rotate all secrets and tokens
- Review and update security policies
- Performance audit of CI/CD pipeline
- Update documentation

## Best Practices

### Commit Messages

Use conventional commits:
```
feat: add new feature
fix: fix bug
docs: update documentation
test: add tests
chore: update dependencies
ci: update CI/CD configuration
```

### Pull Request Workflow

1. Create feature branch from `develop`
2. Make changes and commit
3. Push and create PR to `develop`
4. Wait for CI checks to pass
5. Request review
6. Merge to `develop` (auto-deploys to staging)
7. Test in staging
8. Create PR from `develop` to `main`
9. Get approval and merge (deploys to production)

### Rollback Procedure

If production deployment fails:
```bash
# Automatic rollback is triggered by the workflow

# For manual rollback:
kubectl rollout undo deployment/auth-service -n insavein-production

# Rollback all services:
for service in auth-service user-service savings-service budget-service goal-service education-service notification-service analytics-service frontend; do
  kubectl rollout undo deployment/$service -n insavein-production
done
```

## Support

For issues or questions:
1. Check workflow logs in GitHub Actions
2. Review this guide and README.md
3. Check Kubernetes pod logs
4. Contact DevOps team
5. Create an issue in the repository

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Codecov Documentation](https://docs.codecov.com/)
- [Snyk Documentation](https://docs.snyk.io/)
