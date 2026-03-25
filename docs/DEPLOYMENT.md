# InSavein Platform Deployment Guide

Complete guide for deploying the InSavein platform to production using Kubernetes.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Kubernetes Cluster Setup](#kubernetes-cluster-setup)
3. [Environment Variables and Secrets](#environment-variables-and-secrets)
4. [Database Deployment](#database-deployment)
5. [Backend Services Deployment](#backend-services-deployment)
6. [Frontend Deployment](#frontend-deployment)
7. [Observability Stack](#observability-stack)
8. [Deployment Process](#deployment-process)
9. [Rollback Procedures](#rollback-procedures)
10. [Health Checks and Monitoring](#health-checks-and-monitoring)
11. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Tools

- **kubectl** v1.24+
- **helm** v3.0+ (optional, for observability stack)
- **docker** v20.10+ (for building images)
- **openssl** (for generating secrets)

### Infrastructure Requirements

**Kubernetes Cluster**:
- Kubernetes v1.24 or higher
- Minimum 3 worker nodes
- Total resources:
  - CPU: 50 cores (requests), 100 cores (limits)
  - Memory: 100 GiB (requests), 200 GiB (limits)
  - Storage: 500 GiB

**External Services** (optional):
- Container registry (Docker Hub, ECR, GCR, ACR)
- Load balancer (cloud provider or NGINX)
- DNS provider
- TLS certificate provider (Let's Encrypt, cloud provider)

---

## Kubernetes Cluster Setup

### Option 1: Managed Kubernetes (Recommended for Production)

#### AWS EKS

```bash
# Install eksctl
curl --silent --location "https://github.com/weksctl-io/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin

# Create cluster
eksctl create cluster \
  --name insavein-prod \
  --region us-east-1 \
  --nodegroup-name standard-workers \
  --node-type t3.xlarge \
  --nodes 3 \
  --nodes-min 3 \
  --nodes-max 10 \
  --managed

# Configure kubectl
aws eks update-kubeconfig --region us-east-1 --name insavein-prod
```

#### Google GKE

```bash
# Create cluster
gcloud container clusters create insavein-prod \
  --zone us-central1-a \
  --num-nodes 3 \
  --machine-type n1-standard-4 \
  --enable-autoscaling \
  --min-nodes 3 \
  --max-nodes 10

# Get credentials
gcloud container clusters get-credentials insavein-prod --zone us-central1-a
```

#### Azure AKS

```bash
# Create resource group
az group create --name insavein-rg --location eastus

# Create cluster
az aks create \
  --resource-group insavein-rg \
  --name insavein-prod \
  --node-count 3 \
  --node-vm-size Standard_D4s_v3 \
  --enable-cluster-autoscaler \
  --min-count 3 \
  --max-count 10 \
  --generate-ssh-keys

# Get credentials
az aks get-credentials --resource-group insavein-rg --name insavein-prod
```

### Option 2: Self-Hosted Kubernetes

For self-hosted clusters, use tools like:
- **kubeadm**: Official Kubernetes cluster setup tool
- **k3s**: Lightweight Kubernetes distribution
- **RKE**: Rancher Kubernetes Engine

Refer to official documentation for setup instructions.

### Verify Cluster

```bash
# Check cluster info
kubectl cluster-info

# Check nodes
kubectl get nodes

# Check available resources
kubectl top nodes
```

---

## Environment Variables and Secrets

### Step 1: Generate Secrets

```bash
cd k8s

# Generate all secrets
./generate-secrets.sh

# Or manually generate individual secrets
openssl rand -base64 64  # JWT secret
openssl rand -base64 32  # Encryption keys
openssl rand -base64 32  # Passwords
```

### Step 2: Create Kubernetes Secrets

**Database Credentials**:
```bash
kubectl create secret generic postgres-credentials \
  --from-literal=POSTGRES_USER=insavein_user \
  --from-literal=POSTGRES_PASSWORD=$(openssl rand -base64 32) \
  --from-literal=POSTGRES_DB=insavein \
  --from-literal=REPLICATION_USER=replicator \
  --from-literal=REPLICATION_PASSWORD=$(openssl rand -base64 32) \
  -n insavein
```

**Application Secrets**:
```bash
kubectl create secret generic insavein-secrets \
  --from-literal=JWT_SECRET_KEY=$(openssl rand -base64 64) \
  --from-literal=DATA_ENCRYPTION_KEY=$(openssl rand -base64 32) \
  --from-literal=SESSION_ENCRYPTION_KEY=$(openssl rand -base64 32) \
  --from-literal=EMAIL_API_KEY=your-sendgrid-api-key \
  --from-literal=FCM_SERVER_KEY=your-fcm-server-key \
  --from-literal=OPENAI_API_KEY=your-openai-api-key \
  -n insavein
```

### Step 3: Create ConfigMap

```bash
kubectl apply -f k8s/configmap.yaml
```

### Environment Variables Reference

**Required for all services**:
- `PORT`: Service port (8080-8086, 8008)
- `DB_HOST`: Database host
- `DB_PORT`: Database port (5432)
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password (from secret)
- `DB_NAME`: Database name
- `DB_SSLMODE`: SSL mode (require for production)
- `JWT_SECRET_KEY`: JWT signing key (from secret)
- `LOG_LEVEL`: Logging level (info, debug, error)

**Service-specific**:
- **Auth Service**: `BCRYPT_COST`, `ACCESS_TOKEN_EXPIRY`, `REFRESH_TOKEN_EXPIRY`
- **Notification Service**: `EMAIL_API_KEY`, `FCM_SERVER_KEY`
- **Analytics Service**: `OPENAI_API_KEY`, `CACHE_TTL`

---

## Database Deployment

### Step 1: Deploy PostgreSQL Primary

```bash
kubectl apply -f k8s/postgres-statefulset.yaml
```

**Verify deployment**:
```bash
kubectl get statefulset postgres-primary -n insavein
kubectl get pods -l app=postgres-primary -n insavein
kubectl logs postgres-primary-0 -n insavein
```

### Step 2: Run Database Migrations

```bash
# Port-forward to database
kubectl port-forward svc/postgres-primary 5432:5432 -n insavein

# Run migrations (in another terminal)
cd migrations
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=insavein_user
export DB_PASSWORD=<password-from-secret>
export DB_NAME=insavein
./migrate.sh

# Or use Docker
docker run --rm --network host \
  -v $(pwd)/migrations:/migrations \
  -e DB_HOST=localhost \
  -e DB_PORT=5432 \
  -e DB_USER=insavein_user \
  -e DB_PASSWORD=<password> \
  -e DB_NAME=insavein \
  migrate/migrate -path=/migrations -database "postgres://insavein_user:<password>@localhost:5432/insavein?sslmode=disable" up
```

### Step 3: Deploy Read Replicas

```bash
# Deploy replicas
kubectl apply -f k8s/postgres-replica-statefulset.yaml

# Verify replication
kubectl exec -it postgres-primary-0 -n insavein -- psql -U insavein_user -d insavein -c "SELECT * FROM pg_stat_replication;"
```

### Step 4: Deploy PgBouncer (Connection Pooling)

```bash
kubectl apply -f k8s/pgbouncer-deployment.yaml

# Verify
kubectl get pods -l app=pgbouncer -n insavein
```

---

## Backend Services Deployment

### Deployment Order

Deploy services in this order to respect dependencies:

1. **Auth Service** (critical - required by all services)
2. **User Service** (high priority)
3. **Savings Service** (high priority)
4. **Budget Service** (high priority)
5. **Goal Service** (medium priority)
6. **Education Service** (low priority)
7. **Notification Service** (low priority)
8. **Analytics Service** (medium priority)

### Deploy All Services

```bash
# Deploy in order
kubectl apply -f k8s/auth-service-deployment.yaml
kubectl apply -f k8s/user-service-deployment.yaml
kubectl apply -f k8s/savings-service-deployment.yaml
kubectl apply -f k8s/budget-service-deployment.yaml
kubectl apply -f k8s/goal-service-deployment.yaml
kubectl apply -f k8s/education-service-deployment.yaml
kubectl apply -f k8s/notification-service-deployment.yaml
kubectl apply -f k8s/analytics-service-deployment.yaml

# Or deploy all at once
kubectl apply -f k8s/ -l tier=backend
```

### Verify Deployments

```bash
# Check all deployments
kubectl get deployments -n insavein

# Check pods
kubectl get pods -n insavein

# Check services
kubectl get services -n insavein

# Check logs for a specific service
kubectl logs -l app=auth-service -n insavein --tail=50

# Check health
kubectl exec -it <pod-name> -n insavein -- curl http://localhost:8081/health
```

---

## Frontend Deployment

### Build Frontend Image

```bash
cd frontend

# Build Docker image
docker build -t insavein/frontend:latest .

# Tag for registry
docker tag insavein/frontend:latest <registry>/insavein/frontend:v1.0.0

# Push to registry
docker push <registry>/insavein/frontend:v1.0.0
```

### Deploy Frontend

```bash
# Update image in deployment manifest
# Edit k8s/frontend-deployment.yaml and set image

# Deploy
kubectl apply -f k8s/frontend-deployment.yaml

# Verify
kubectl get pods -l app=frontend -n insavein
kubectl logs -l app=frontend -n insavein
```

### Configure Ingress

```bash
# Deploy ingress controller (if not already installed)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Deploy ingress rules
kubectl apply -f k8s/ingress.yaml

# Get ingress IP
kubectl get ingress -n insavein
```

### Configure DNS

Point your domain to the ingress IP:
```
A    insavein.com          -> <ingress-ip>
A    api.insavein.com      -> <ingress-ip>
A    www.insavein.com      -> <ingress-ip>
```

### Configure TLS

```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Create ClusterIssuer for Let's Encrypt
kubectl apply -f k8s/cert-manager-issuer.yaml

# Certificate will be automatically provisioned
kubectl get certificate -n insavein
```

---

## Observability Stack

### Deploy Prometheus

```bash
kubectl apply -f k8s/prometheus-deployment.yaml

# Verify
kubectl get pods -l app=prometheus -n insavein
```

### Deploy Grafana

```bash
kubectl apply -f k8s/grafana-deployment.yaml

# Get Grafana password
kubectl get secret grafana-admin -n insavein -o jsonpath='{.data.password}' | base64 -d

# Port-forward to access
kubectl port-forward svc/grafana 3000:3000 -n insavein
# Access at http://localhost:3000
```

### Import Dashboards

```bash
# Dashboards are pre-configured in grafana-dashboards.yaml
kubectl apply -f k8s/grafana-dashboards.yaml
```

### Configure Alerts

```bash
# Alert rules are defined in prometheus-alerts.yaml
kubectl apply -f k8s/prometheus-alerts.yaml

# Verify alerts
kubectl exec -it <prometheus-pod> -n insavein -- promtool check rules /etc/prometheus/alerts.yml
```

---

## Deployment Process

### Rolling Deployment Strategy

All services use rolling updates with:
- **Max Surge**: 25% (1 extra pod during update)
- **Max Unavailable**: 0 (no downtime)

### Deploy New Version

```bash
# Update image tag in deployment
kubectl set image deployment/auth-service \
  auth-service=<registry>/auth-service:v1.1.0 \
  -n insavein

# Or edit deployment
kubectl edit deployment auth-service -n insavein

# Watch rollout
kubectl rollout status deployment/auth-service -n insavein

# Check rollout history
kubectl rollout history deployment/auth-service -n insavein
```

### Automated Deployment (CI/CD)

See `.github/workflows/deploy-production.yml` for automated deployment pipeline.

**Deployment steps**:
1. Run tests
2. Build Docker images
3. Push to registry
4. Update Kubernetes manifests
5. Apply to cluster
6. Verify health checks
7. Run smoke tests
8. Rollback on failure

---

## Rollback Procedures

### Quick Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/auth-service -n insavein

# Rollback to specific revision
kubectl rollout undo deployment/auth-service --to-revision=3 -n insavein

# Check rollback status
kubectl rollout status deployment/auth-service -n insavein
```

### Manual Rollback

```bash
# Scale down new version
kubectl scale deployment auth-service --replicas=0 -n insavein

# Deploy previous version
kubectl apply -f k8s/auth-service-deployment-v1.0.0.yaml

# Verify
kubectl get pods -l app=auth-service -n insavein
```

### Database Rollback

```bash
# Rollback migrations
cd migrations
./migrate.sh down 1  # Rollback 1 migration

# Or rollback to specific version
./migrate.sh goto 5  # Rollback to version 5
```

### Rollback Checklist

- [ ] Identify failing service/version
- [ ] Check error logs and metrics
- [ ] Execute rollback command
- [ ] Verify health checks pass
- [ ] Check application functionality
- [ ] Monitor error rates
- [ ] Notify team of rollback
- [ ] Document incident
- [ ] Plan fix for next deployment

---

## Health Checks and Monitoring

### Health Check Endpoints

All services expose:
- `/health` - Overall health
- `/health/live` - Liveness probe
- `/health/ready` - Readiness probe

### Kubernetes Probes

**Liveness Probe**:
```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8081
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

**Readiness Probe**:
```yaml
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8081
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

### Monitoring Metrics

**Key Metrics**:
- Request rate (requests/second)
- Error rate (errors/total requests)
- Response time (p50, p95, p99)
- CPU usage (%)
- Memory usage (%)
- Database connections
- Queue depth

**Access Metrics**:
```bash
# Port-forward Prometheus
kubectl port-forward svc/prometheus 9090:9090 -n insavein

# Access at http://localhost:9090
```

### Alerts

**Critical Alerts**:
- Service down (no healthy pods)
- High error rate (>1%)
- High response time (p95 >1s)
- Database connection failures
- Out of memory errors

**Warning Alerts**:
- High CPU usage (>80%)
- High memory usage (>85%)
- Slow response time (p95 >500ms)
- Replication lag (>5s)

---

## Troubleshooting

### Common Issues

#### Pods Not Starting

```bash
# Check pod status
kubectl describe pod <pod-name> -n insavein

# Check events
kubectl get events -n insavein --sort-by='.lastTimestamp'

# Check logs
kubectl logs <pod-name> -n insavein --previous
```

**Common causes**:
- Image pull errors (check registry credentials)
- Resource limits exceeded (check resource quotas)
- Missing secrets/configmaps
- Failed health checks

#### Database Connection Errors

```bash
# Check database pods
kubectl get pods -l app=postgres-primary -n insavein

# Check database logs
kubectl logs postgres-primary-0 -n insavein

# Test connection from service pod
kubectl exec -it <service-pod> -n insavein -- \
  psql -h postgres-primary -U insavein_user -d insavein
```

#### High Memory Usage

```bash
# Check memory usage
kubectl top pods -n insavein

# Check memory limits
kubectl describe pod <pod-name> -n insavein | grep -A 5 Limits

# Increase memory limits if needed
kubectl edit deployment <service-name> -n insavein
```

#### Service Unavailable

```bash
# Check service endpoints
kubectl get endpoints <service-name> -n insavein

# Check if pods are ready
kubectl get pods -l app=<service-name> -n insavein

# Check service logs
kubectl logs -l app=<service-name> -n insavein --tail=100
```

### Debug Commands

```bash
# Get all resources
kubectl get all -n insavein

# Describe resource
kubectl describe <resource-type> <resource-name> -n insavein

# Get logs
kubectl logs <pod-name> -n insavein -f

# Execute command in pod
kubectl exec -it <pod-name> -n insavein -- /bin/sh

# Port-forward for local access
kubectl port-forward <pod-name> 8080:8080 -n insavein

# Check resource usage
kubectl top pods -n insavein
kubectl top nodes
```

---

## Security Checklist

- [ ] All secrets stored in Kubernetes Secrets (not in code)
- [ ] TLS enabled for all external traffic
- [ ] Database connections encrypted
- [ ] Network policies configured
- [ ] RBAC policies configured
- [ ] Pod security policies enabled
- [ ] Container images scanned for vulnerabilities
- [ ] Resource limits set for all pods
- [ ] Secrets rotation schedule configured
- [ ] Audit logging enabled

---

## Post-Deployment Checklist

- [ ] All pods running and healthy
- [ ] Health checks passing
- [ ] Database migrations completed
- [ ] Ingress configured and accessible
- [ ] TLS certificates provisioned
- [ ] DNS records updated
- [ ] Monitoring dashboards configured
- [ ] Alerts configured and tested
- [ ] Backup jobs scheduled
- [ ] Documentation updated
- [ ] Team notified of deployment

---

## Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [k8s/README.md](../k8s/README.md) - Kubernetes configuration details
- [k8s/DEPLOYMENT_GUIDE.md](../k8s/DEPLOYMENT_GUIDE.md) - Base configuration deployment
- [CI_CD_SETUP_GUIDE.md](../CI_CD_SETUP_GUIDE.md) - CI/CD pipeline setup

---

**Last Updated**: 2026-01-15  
**Version**: 1.0.0
