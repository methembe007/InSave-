# InSavein Platform - CI/CD Pipeline Overview

## Pipeline Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Developer Workflow                           │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │   Git Push / Pull Request  │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │     GitHub Actions         │
                    └─────────────┬─────────────┘
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        │                         │                         │
┌───────▼────────┐    ┌──────────▼──────────┐    ┌────────▼────────┐
│  Lint (2-3min) │    │  Test (5-7min)      │    │ Security (8min) │
│                │    │                     │    │                 │
│ • golangci-lint│    │ • Go unit tests     │    │ • Trivy scan    │
│ • ESLint       │    │ • Frontend tests    │    │ • Snyk scan     │
│ • TypeScript   │    │ • Integration tests │    │ • GoSec         │
│                │    │ • Coverage report   │    │ • CodeQL        │
└───────┬────────┘    └──────────┬──────────┘    └────────┬────────┘
        │                        │                         │
        └────────────────────────┼─────────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │  Build & Push (10min)   │
                    │                         │
                    │ • Build 9 Docker images │
                    │ • Multi-platform        │
                    │ • Push to GHCR          │
                    │ • Scan images           │
                    └────────────┬────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │   Branch Detection      │
                    └────────────┬────────────┘
                                 │
                ┌────────────────┴────────────────┐
                │                                 │
    ┌───────────▼──────────┐         ┌──────────▼───────────┐
    │  Staging Deployment  │         │ Production Deployment │
    │  (develop branch)    │         │   (main branch)       │
    │                      │         │                       │
    │ • Auto-deploy        │         │ • Manual approval ✋  │
    │ • Zero downtime      │         │ • Backup state        │
    │ • Smoke tests        │         │ • Rolling update      │
    │ • Slack notify       │         │ • Extended tests      │
    │                      │         │ • Health monitoring   │
    │ ✅ staging.insavein  │         │ • Auto-rollback       │
    └──────────────────────┘         │ • Slack + Email       │
                                     │                       │
                                     │ ✅ insavein.com       │
                                     └───────────────────────┘
```

## Workflow Summary

| Workflow | Trigger | Duration | Purpose |
|----------|---------|----------|---------|
| **Lint** | Push/PR | 2-3 min | Code quality checks |
| **Test** | Push/PR | 5-7 min | Unit & integration tests |
| **Security** | Push/PR/Daily | 8 min | Vulnerability scanning |
| **Build** | Push | 10 min | Docker image creation |
| **Deploy Staging** | Push to develop | 5 min | Auto-deploy to staging |
| **Deploy Production** | Push to main | 8 min | Deploy to production |

## Service Matrix

### Backend Services (Go)
1. auth-service
2. user-service
3. savings-service
4. budget-service
5. goal-service
6. education-service
7. notification-service
8. analytics-service

### Frontend
9. frontend (TanStack Start)

## Deployment Flow

### Staging (Automatic)
```
develop branch push
    ↓
Build images (develop-<sha>)
    ↓
Deploy to insavein-staging namespace
    ↓
Run smoke tests
    ↓
Notify team via Slack
```

### Production (Manual Approval)
```
main branch push
    ↓
Build images (main-<sha>)
    ↓
⏸️  Wait for manual approval
    ↓
Backup current state
    ↓
Deploy to insavein-production namespace
    ↓
Run comprehensive tests
    ↓
Monitor health for 2 minutes
    ↓
Notify team via Slack + Email
    ↓
✅ Success or 🔄 Auto-rollback
```

## Security Layers

```
┌─────────────────────────────────────────┐
│         Code Level Security             │
│  • GoSec (static analysis)              │
│  • CodeQL (semantic analysis)           │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Dependency Security                │
│  • Snyk (Go modules + npm)              │
│  • Dependabot alerts                    │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Infrastructure Security            │
│  • Trivy (filesystem + configs)         │
│  • Kubernetes security policies         │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│       Container Security                │
│  • Trivy (image scanning)               │
│  • Multi-stage builds                   │
│  • Non-root users                       │
└─────────────────────────────────────────┘
```

## Key Features

### ✅ Zero Downtime Deployments
- Rolling updates with maxUnavailable=0
- Health checks before traffic routing
- Gradual rollout (25% surge)

### ✅ Automatic Rollback
- Failed health checks trigger rollback
- Previous version restored in < 2 minutes
- Notifications sent to team

### ✅ Comprehensive Testing
- Unit tests with race detection
- Integration tests with PostgreSQL
- Smoke tests for all endpoints
- Critical user flow validation

### ✅ Security First
- Multiple scanning tools
- Daily scheduled scans
- SARIF results in GitHub Security
- Automated dependency updates

### ✅ Multi-Platform Support
- linux/amd64 (Intel/AMD)
- linux/arm64 (ARM/Apple Silicon)

### ✅ Smart Caching
- Go module cache
- npm package cache
- Docker layer cache
- golangci-lint cache

## Monitoring & Alerts

### Notifications
- **Slack**: All deployments
- **Email**: Critical failures only
- **GitHub**: Status checks on PRs

### Metrics Tracked
- Build duration
- Test coverage
- Deployment success rate
- Security vulnerabilities
- Rollback frequency

## Quick Commands

### Check Status
```bash
# View workflow runs
gh run list

# Watch deployment
kubectl get pods -n insavein-production -w
```

### Manual Deploy
```bash
# Trigger staging
gh workflow run deploy-staging.yml

# Trigger production (requires approval)
gh workflow run deploy-production.yml
```

### Rollback
```bash
# Rollback single service
kubectl rollout undo deployment/auth-service -n insavein-production

# Rollback all services
for svc in auth user savings budget goal education notification analytics frontend; do
  kubectl rollout undo deployment/${svc}-service -n insavein-production
done
```

## Documentation

- **Setup Guide**: [CI_CD_SETUP_GUIDE.md](./CI_CD_SETUP_GUIDE.md)
- **Quick Reference**: [CI_CD_QUICK_REFERENCE.md](./CI_CD_QUICK_REFERENCE.md)
- **Workflow Details**: [.github/workflows/README.md](./.github/workflows/README.md)
- **Implementation**: [TASK_26_CI_CD_IMPLEMENTATION.md](./TASK_26_CI_CD_IMPLEMENTATION.md)

## Success Criteria

- ✅ All 6 workflows implemented
- ✅ Linting for Go and TypeScript
- ✅ Automated testing with coverage
- ✅ Multi-layer security scanning
- ✅ Multi-platform Docker builds
- ✅ Staging auto-deployment
- ✅ Production deployment with approval
- ✅ Zero-downtime updates
- ✅ Automatic rollback capability
- ✅ Comprehensive documentation

## Getting Started

1. **Configure Secrets**: See [CI_CD_SETUP_GUIDE.md](./CI_CD_SETUP_GUIDE.md)
2. **Set Up Environments**: Create staging and production
3. **Test Pipeline**: Push to develop branch
4. **Monitor**: Check GitHub Actions tab
5. **Deploy**: Merge to main with approval

---

**Status**: ✅ Fully Implemented  
**Last Updated**: 2024  
**Maintained By**: DevOps Team
