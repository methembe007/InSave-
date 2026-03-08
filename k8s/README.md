# InSavein Kubernetes Configuration

This directory contains Kubernetes manifests for deploying the InSavein platform.

## Overview

The InSavein platform uses Kubernetes for orchestration with the following components:
- **Namespace**: Isolated environment for all InSavein resources
- **ConfigMaps**: Non-sensitive configuration (database hosts, service endpoints, rate limits)
- **Secrets**: Sensitive data (database credentials, JWT secrets, API keys)
- **Resource Quotas**: Limits on compute, memory, and storage resources
- **Network Policies**: Security rules for pod-to-pod communication
- **Priority Classes**: Pod scheduling priorities for critical services

## Prerequisites

- Kubernetes cluster (v1.24+)
- kubectl configured to access your cluster
- Sufficient cluster resources (see resource-quota.yaml for requirements)

## Quick Start

### 1. Create Namespace and Base Configuration

```bash
# Create namespace
kubectl apply -f namespace.yaml

# Verify namespace creation
kubectl get namespace insavein
```

### 2. Configure Secrets (IMPORTANT!)

**⚠️ SECURITY WARNING**: The secrets.yaml file contains placeholder values that MUST be replaced before deployment.

Generate secure secrets:

```bash
# Generate JWT secret (64 bytes, base64 encoded)
openssl rand -base64 64

# Generate encryption keys (32 bytes, base64 encoded)
openssl rand -base64 32

# Generate strong passwords
openssl rand -base64 32
```

Edit `secrets.yaml` and replace ALL `CHANGE_ME_*` placeholders with actual values:
- Database credentials
- JWT secret key
- Email service API keys (SendGrid/AWS SES)
- Firebase Cloud Messaging keys
- OpenAI API key (for AI recommendations)
- Encryption keys

**Production Best Practice**: Use external secret management instead of storing secrets in YAML:
- HashiCorp Vault
- AWS Secrets Manager
- Azure Key Vault
- Google Secret Manager

### 3. Apply Configuration

```bash
# Apply priority classes (cluster-wide)
kubectl apply -f priority-class.yaml

# Apply ConfigMap
kubectl apply -f configmap.yaml

# Apply Secrets (after updating with real values!)
kubectl apply -f secrets.yaml

# Apply resource quotas and limits
kubectl apply -f resource-quota.yaml

# Apply network policies
kubectl apply -f network-policy.yaml
```

### 4. Verify Configuration

```bash
# Check all resources in namespace
kubectl get all -n insavein

# Verify ConfigMap
kubectl describe configmap insavein-config -n insavein

# Verify Secrets (values will be hidden)
kubectl describe secret insavein-secrets -n insavein

# Check resource quotas
kubectl describe resourcequota insavein-resource-quota -n insavein

# Check limit ranges
kubectl describe limitrange insavein-limit-range -n insavein

# Verify network policies
kubectl get networkpolicies -n insavein
```

## Configuration Details

### ConfigMap (configmap.yaml)

Contains non-sensitive environment variables:

**Database Configuration**:
- Primary database host and port
- Read replica hosts (2 replicas)
- Connection pooling settings
- PgBouncer configuration

**Service Endpoints**:
- Internal service URLs for all 8 microservices
- Uses Kubernetes DNS (*.insavein.svc.cluster.local)

**Rate Limiting** (Requirement 18.1):
- 100 requests per minute per user
- 1000 requests per minute per IP
- Burst allowance of 20 requests

**JWT Configuration**:
- Access token expiry: 15 minutes
- Refresh token expiry: 7 days
- Algorithm: HMAC-SHA256

**Observability**:
- Metrics port: 9090 (Prometheus)
- Health check port: 8081
- Trace sampling: 10%

### Secrets (secrets.yaml)

Contains sensitive data (Requirement 20.1):

**Database Credentials**:
- Application user credentials
- Replica user credentials
- PostgreSQL superuser credentials
- Replication user credentials

**JWT Secret Key**:
- Used for signing and verifying JWT tokens
- Must be kept secure and rotated every 90 days

**API Keys**:
- Email service (SendGrid/AWS SES)
- Push notifications (Firebase Cloud Messaging)
- AI services (OpenAI for recommendations)

**Encryption Keys**:
- Data encryption key (for data at rest)
- Session encryption key

### Resource Quotas (resource-quota.yaml)

Limits resource usage in the namespace:

**Compute Resources**:
- Total CPU requests: 50 cores
- Total memory requests: 100 GiB
- Total CPU limits: 100 cores
- Total memory limits: 200 GiB

**Storage Resources**:
- Total storage: 500 GiB
- Max PersistentVolumeClaims: 10

**Object Counts**:
- Max pods: 100
- Max services: 20
- Max ConfigMaps: 20
- Max Secrets: 20

**Container Limits** (per container):
- Default CPU: 1 core (limit), 100m (request)
- Default memory: 2 GiB (limit), 256 MiB (request)
- Max CPU: 4 cores
- Max memory: 8 GiB

**Database Limits** (special limits for PostgreSQL):
- Default CPU: 2 cores (limit), 500m (request)
- Default memory: 4 GiB (limit), 2 GiB (request)
- Max CPU: 8 cores
- Max memory: 16 GiB

### Network Policies (network-policy.yaml)

Security rules for pod communication:

1. **Default Deny**: Blocks all ingress and egress by default
2. **Frontend to Services**: Allows frontend to call backend services on port 8080
3. **Services to Database**: Allows backend services to access PostgreSQL (5432) and PgBouncer (6432)
4. **Ingress to Frontend**: Allows ingress controller to route traffic to frontend on port 3000
5. **Prometheus Scraping**: Allows Prometheus to scrape metrics on port 9090
6. **DNS Egress**: Allows all pods to resolve DNS
7. **External APIs**: Allows pods with `needs-external: "true"` label to call external APIs

### Priority Classes (priority-class.yaml)

Pod scheduling priorities:

1. **insavein-critical** (1,000,000): Auth Service, Database
2. **insavein-high** (100,000): User, Savings, Budget services
3. **insavein-medium** (10,000): Goal, Analytics services (default)
4. **insavein-low** (1,000): Education, Notification services

## Security Considerations

### Secrets Management

**Development/Testing**:
- Use Kubernetes Secrets with base64 encoding
- Store secrets.yaml in a secure location (NOT in version control)
- Use `.gitignore` to exclude secrets.yaml

**Production**:
- Use external secret management (Vault, AWS Secrets Manager, etc.)
- Enable secret encryption at rest in etcd
- Rotate secrets every 90 days (Requirement 20.6)
- Use RBAC to restrict secret access

### Network Security

- Network policies enforce zero-trust networking
- All inter-service communication is restricted by default
- Only necessary ports are exposed
- External API access requires explicit label

### TLS/Encryption

- All external traffic uses TLS 1.3 (Requirement 20.4)
- Database connections are encrypted (Requirement 20.5)
- Passwords are never stored in plaintext (Requirement 20.1)

## Monitoring and Observability

### Health Checks

All services expose health check endpoints:
- Liveness probe: `/health/live` on port 8081
- Readiness probe: `/health/ready` on port 8081

### Metrics

All services expose Prometheus metrics:
- Metrics endpoint: `/metrics` on port 9090
- Includes request count, duration, error rate

### Logging

All services use structured JSON logging:
- Log level: info (configurable via LOG_LEVEL)
- Includes trace IDs for distributed tracing
- Retention: 30 days (INFO), 90 days (ERROR)

## Scaling Configuration

### Horizontal Pod Autoscaling

Services will scale based on:
- CPU utilization > 70%
- Memory utilization > 80%

Minimum replicas:
- Auth Service: 3
- Savings Service: 3
- Budget Service: 3
- Other services: 2

### Database Scaling

- Primary database: 1 instance (write operations)
- Read replicas: 2 instances (read-heavy operations)
- Connection pooling via PgBouncer (max 20 connections per service)

## Deployment Order

1. **Infrastructure** (this directory):
   - Namespace
   - Priority classes
   - ConfigMaps
   - Secrets
   - Resource quotas
   - Network policies

2. **Database**:
   - PostgreSQL primary
   - PostgreSQL replicas
   - PgBouncer

3. **Backend Services**:
   - Auth Service (critical)
   - User Service
   - Savings Service
   - Budget Service
   - Goal Service
   - Education Service
   - Notification Service
   - Analytics Service

4. **Frontend**:
   - TanStack Start application

5. **Observability**:
   - Prometheus
   - Grafana
   - OpenTelemetry Collector

## Troubleshooting

### Check Resource Usage

```bash
# View resource quota usage
kubectl describe resourcequota insavein-resource-quota -n insavein

# View pod resource usage
kubectl top pods -n insavein

# View node resource usage
kubectl top nodes
```

### Check Network Policies

```bash
# List all network policies
kubectl get networkpolicies -n insavein

# Describe specific policy
kubectl describe networkpolicy allow-services-to-database -n insavein

# Test connectivity between pods
kubectl run -it --rm debug --image=nicolaka/netshoot -n insavein -- /bin/bash
```

### Check Secrets

```bash
# Verify secret exists
kubectl get secret insavein-secrets -n insavein

# View secret keys (values are hidden)
kubectl describe secret insavein-secrets -n insavein

# Decode secret value (for debugging only!)
kubectl get secret insavein-secrets -n insavein -o jsonpath='{.data.JWT_SECRET_KEY}' | base64 -d
```

### Common Issues

**Pods not starting**:
- Check resource quotas: `kubectl describe resourcequota -n insavein`
- Check limit ranges: `kubectl describe limitrange -n insavein`
- Check pod events: `kubectl describe pod <pod-name> -n insavein`

**Network connectivity issues**:
- Verify network policies allow the connection
- Check service DNS resolution: `nslookup <service-name>.insavein.svc.cluster.local`
- Test connectivity with netshoot pod

**Secret not found**:
- Verify secret exists: `kubectl get secret -n insavein`
- Check secret is in correct namespace
- Verify pod has correct secret reference

## Next Steps

After applying these base configurations:

1. Deploy PostgreSQL database (primary + replicas)
2. Deploy PgBouncer for connection pooling
3. Deploy backend microservices
4. Deploy frontend application
5. Configure ingress controller
6. Set up monitoring and alerting
7. Configure CI/CD pipeline

## References

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/)
- [Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [Secrets Management](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Priority Classes](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/)
