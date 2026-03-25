# CI/CD Pipeline Quick Reference

Quick reference guide for common CI/CD operations on the InSavein Platform.

## Workflow Triggers

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| **Lint** | Push/PR to main/develop | Code quality checks |
| **Test** | Push/PR to main/develop | Run tests with coverage |
| **Security** | Push/PR + Daily 2AM | Security vulnerability scanning |
| **Build/Push** | Push to main/develop, Tags | Build and push Docker images |
| **Deploy Staging** | Push to develop | Auto-deploy to staging |
| **Deploy Production** | Push to main, Tags | Deploy to production (requires approval) |

## Common Commands

### Check Workflow Status

```bash
# View recent workflow runs
gh run list

# View specific workflow
gh run view <run-id>

# Watch workflow in real-time
gh run watch
```

### Manual Workflow Trigger

```bash
# Trigger staging deployment
gh workflow run deploy-staging.yml

# Trigger production deployment
gh workflow run deploy-production.yml
```

### Check Deployment Status

```bash
# Staging
kubectl get pods -n insavein-staging
kubectl get deployments -n insavein-staging

# Production
kubectl get pods -n insavein-production
kubectl get deployments -n insavein-production
```

### View Logs

```bash
# View service logs (staging)
kubectl logs -f deployment/auth-service -n insavein-staging

# View service logs (production)
kubectl logs -f deployment/auth-service -n insavein-production

# View logs from all pods
kubectl logs -l app=auth-service -n insavein-production --tail=100
```

### Rollback Deployment

```bash
# Rollback single service
kubectl rollout undo deployment/auth-service -n insavein-production

# Rollback to specific revision
kubectl rollout undo deployment/auth-service --to-revision=2 -n insavein-production

# View rollout history
kubectl rollout history deployment/auth-service -n insavein-production
```

### Scale Services

```bash
# Scale up
kubectl scale deployment/auth-service --replicas=5 -n insavein-production

# Scale down
kubectl scale deployment/auth-service --replicas=3 -n insavein-production
```

## Image Tags

| Tag Format | Example | Description |
|------------|---------|-------------|
| `develop-<sha>` | `develop-abc1234` | Develop branch builds |
| `main-<sha>` | `main-def5678` | Main branch builds |
| `v*` | `v1.2.3` | Semantic version tags |
| `latest` | `latest` | Latest main branch |

## Environment URLs

| Environment | URL | Namespace |
|-------------|-----|-----------|
| **Staging** | https://staging.insavein.com | `insavein-staging` |
| **Production** | https://insavein.com | `insavein-production` |

## Health Check Endpoints

```bash
# Auth Service
curl https://insavein.com/api/auth/health

# User Service
curl https://insavein.com/api/user/health

# Savings Service
curl https://insavein.com/api/savings/health

# Budget Service
curl https://insavein.com/api/budget/health

# Goal Service
curl https://insavein.com/api/goal/health

# Education Service
curl https://insavein.com/api/education/health

# Notification Service
curl https://insavein.com/api/notification/health

# Analytics Service
curl https://insavein.com/api/analytics/health

# Frontend
curl https://insavein.com
```

## Secrets Management

### View Secrets (GitHub CLI)

```bash
# List all secrets
gh secret list

# View secret (shows when it was updated, not the value)
gh secret view JWT_SECRET_PRODUCTION
```

### Update Secrets

```bash
# Update secret
gh secret set JWT_SECRET_STAGING

# Update from file
gh secret set KUBE_CONFIG_STAGING < kubeconfig-staging.txt

# Update with base64 encoding
cat kubeconfig-staging.yaml | base64 -w 0 | gh secret set KUBE_CONFIG_STAGING
```

### Rotate Secrets

```bash
# Generate new JWT secret
openssl rand -base64 32

# Update in GitHub
gh secret set JWT_SECRET_PRODUCTION

# Update in Kubernetes
kubectl create secret generic app-secrets \
  --from-literal=jwt-secret="NEW_SECRET" \
  --namespace=insavein-production \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart deployments to pick up new secret
kubectl rollout restart deployment/auth-service -n insavein-production
```

## Troubleshooting Quick Fixes

### Workflow Stuck

```bash
# Cancel running workflow
gh run cancel <run-id>

# Re-run failed workflow
gh run rerun <run-id>
```

### Pod CrashLoopBackOff

```bash
# Check pod status
kubectl describe pod <pod-name> -n insavein-production

# Check logs
kubectl logs <pod-name> -n insavein-production --previous

# Check events
kubectl get events -n insavein-production --sort-by='.lastTimestamp' | grep <pod-name>
```

### Image Pull Errors

```bash
# Verify image exists
docker pull ghcr.io/your-org/insavein-auth-service:main-abc1234

# Recreate image pull secret
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<github-username> \
  --docker-password=<github-token> \
  --namespace=insavein-production \
  --dry-run=client -o yaml | kubectl apply -f -
```

### Database Connection Issues

```bash
# Test database connectivity from pod
kubectl run -it --rm debug --image=postgres:16-alpine --restart=Never -n insavein-production -- \
  psql "postgres://username:password@postgres:5432/insavein"

# Check database service
kubectl get svc postgres -n insavein-production
kubectl describe svc postgres -n insavein-production
```

## Monitoring

### Prometheus Queries

```promql
# Request rate
rate(http_requests_total[5m])

# Error rate
rate(http_requests_total{status=~"5.."}[5m])

# Request duration (p95)
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Pod CPU usage
rate(container_cpu_usage_seconds_total{namespace="insavein-production"}[5m])

# Pod memory usage
container_memory_usage_bytes{namespace="insavein-production"}
```

### Access Monitoring Tools

```bash
# Port forward Prometheus
kubectl port-forward -n insavein-production svc/prometheus 9090:9090

# Port forward Grafana
kubectl port-forward -n insavein-production svc/grafana 3000:3000

# Access in browser
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3000
```

## Performance Optimization

### Check Workflow Duration

```bash
# View workflow timing
gh run view <run-id> --log

# View workflow metrics
gh api repos/:owner/:repo/actions/workflows/:workflow_id/timing
```

### Clear Caches

```bash
# Clear GitHub Actions cache
gh cache list
gh cache delete <cache-id>

# Clear all caches for a branch
gh cache delete --all --branch develop
```

## Security

### View Security Alerts

```bash
# View Dependabot alerts
gh api repos/:owner/:repo/dependabot/alerts

# View code scanning alerts
gh api repos/:owner/:repo/code-scanning/alerts
```

### Scan Local Images

```bash
# Scan with Trivy
trivy image ghcr.io/your-org/insavein-auth-service:latest

# Scan with Snyk
snyk container test ghcr.io/your-org/insavein-auth-service:latest
```

## Deployment Checklist

### Before Deploying to Production

- [ ] All tests passing in staging
- [ ] Security scans completed with no critical issues
- [ ] Database migrations tested in staging
- [ ] Performance testing completed
- [ ] Rollback plan documented
- [ ] Team notified of deployment
- [ ] Monitoring dashboards ready
- [ ] On-call engineer available

### After Deploying to Production

- [ ] Verify all health checks passing
- [ ] Check error rates in monitoring
- [ ] Verify critical user flows working
- [ ] Monitor for 15 minutes
- [ ] Update deployment documentation
- [ ] Notify team of successful deployment

## Emergency Procedures

### Complete System Rollback

```bash
#!/bin/bash
# Rollback all services to previous version

NAMESPACE="insavein-production"
SERVICES=(
  "auth-service"
  "user-service"
  "savings-service"
  "budget-service"
  "goal-service"
  "education-service"
  "notification-service"
  "analytics-service"
  "frontend"
)

for service in "${SERVICES[@]}"; do
  echo "Rolling back $service..."
  kubectl rollout undo deployment/$service -n $NAMESPACE
done

echo "Waiting for rollback to complete..."
for service in "${SERVICES[@]}"; do
  kubectl rollout status deployment/$service -n $NAMESPACE
done

echo "Rollback complete!"
```

### Emergency Scale Down

```bash
#!/bin/bash
# Scale down all services to minimum replicas

NAMESPACE="insavein-production"
SERVICES=(
  "auth-service"
  "user-service"
  "savings-service"
  "budget-service"
  "goal-service"
  "education-service"
  "notification-service"
  "analytics-service"
  "frontend"
)

for service in "${SERVICES[@]}"; do
  echo "Scaling down $service..."
  kubectl scale deployment/$service --replicas=1 -n $NAMESPACE
done
```

### Emergency Maintenance Mode

```bash
# Enable maintenance mode (redirect all traffic to maintenance page)
kubectl apply -f k8s/maintenance-mode.yaml -n insavein-production

# Disable maintenance mode
kubectl delete -f k8s/maintenance-mode.yaml -n insavein-production
```

## Useful Aliases

Add to your `.bashrc` or `.zshrc`:

```bash
# Kubernetes aliases
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgd='kubectl get deployments'
alias kl='kubectl logs -f'
alias kd='kubectl describe'

# InSavein specific
alias kp-staging='kubectl -n insavein-staging'
alias kp-prod='kubectl -n insavein-production'
alias kp-logs='kubectl logs -f -n insavein-production'
alias kp-pods='kubectl get pods -n insavein-production'

# GitHub CLI aliases
alias ghw='gh workflow'
alias ghr='gh run'
alias ghrl='gh run list'
alias ghrv='gh run view'
```

## Contact Information

| Role | Contact | Availability |
|------|---------|--------------|
| **DevOps Lead** | devops@insavein.com | 24/7 |
| **On-Call Engineer** | oncall@insavein.com | 24/7 |
| **Security Team** | security@insavein.com | Business hours |
| **Slack Channel** | #devops-alerts | 24/7 |

## Additional Resources

- [Full CI/CD Setup Guide](./CI_CD_SETUP_GUIDE.md)
- [Workflow Documentation](./.github/workflows/README.md)
- [Kubernetes Deployment Guide](./k8s/DEPLOYMENT_GUIDE.md)
- [Security Implementation](./TASK_25_SECURITY_IMPLEMENTATION.md)
- [Observability Guide](./TASK_24_OBSERVABILITY_IMPLEMENTATION.md)
