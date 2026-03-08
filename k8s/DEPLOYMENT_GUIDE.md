# InSavein Kubernetes Deployment Guide

This guide provides step-by-step instructions for deploying the InSavein platform base configuration to Kubernetes.

## Prerequisites

Before you begin, ensure you have:

1. **Kubernetes Cluster** (v1.24 or higher)
   - Managed cluster (EKS, GKE, AKS) or self-hosted
   - Sufficient resources (see resource-quota.yaml)
   - kubectl configured and authenticated

2. **Required Tools**:
   - `kubectl` (v1.24+)
   - `openssl` (for generating secrets)
   - `make` (optional, for using Makefile)

3. **Access Requirements**:
   - Cluster admin permissions (for creating namespace, priority classes)
   - Ability to create secrets and configmaps

## Deployment Steps

### Step 1: Verify Cluster Access

```bash
# Check kubectl is configured
kubectl cluster-info

# Check available resources
kubectl get nodes
kubectl top nodes

# Verify you have admin permissions
kubectl auth can-i create namespace
```

### Step 2: Generate Secrets

**Option A: Using the script (Linux/Mac)**

```bash
cd k8s
chmod +x generate-secrets.sh
./generate-secrets.sh
```

The script will:
- Generate secure random values for all secrets
- Optionally create a `secrets-generated.yaml` file
- Display all generated values

**Option B: Using Make (Linux/Mac)**

```bash
cd k8s
make generate-secrets
```

**Option C: Manual generation**

```bash
# JWT Secret (64 bytes)
openssl rand -base64 64

# Encryption keys (32 bytes)
openssl rand -base64 32

# Passwords
openssl rand -base64 32
```

### Step 3: Update Secrets File

1. Copy the generated values from Step 2
2. Edit `secrets.yaml` or use `secrets-generated.yaml`
3. Replace ALL `CHANGE_ME_*` placeholders:

   **Required replacements**:
   - `DB_PASSWORD` - Database user password
   - `DB_REPLICA_PASSWORD` - Replica user password
   - `JWT_SECRET_KEY` - JWT signing key
   - `DATA_ENCRYPTION_KEY` - Data encryption key
   - `SESSION_ENCRYPTION_KEY` - Session encryption key
   - `POSTGRES_PASSWORD` - PostgreSQL superuser password
   - `REPLICATION_PASSWORD` - Replication user password

   **External service credentials** (add when available):
   - `EMAIL_API_KEY` - SendGrid or AWS SES API key
   - `EMAIL_API_SECRET` - Email service secret
   - `FCM_SERVER_KEY` - Firebase Cloud Messaging key
   - `FCM_PROJECT_ID` - Firebase project ID
   - `OPENAI_API_KEY` - OpenAI API key (for AI recommendations)

4. **IMPORTANT**: Do NOT commit the updated secrets file to version control!

### Step 4: Review ConfigMap

Review `configmap.yaml` and update if needed:

- **Database hosts**: Update if using external database
- **Service URLs**: Update if using different service names
- **Rate limits**: Adjust based on requirements
- **CORS origins**: Update with your actual domain names

### Step 5: Deploy Base Configuration

**Option A: Using Make (recommended)**

```bash
cd k8s

# Check for placeholder values in secrets
make check-secrets

# Apply all configurations
make apply

# Verify deployment
make verify
```

**Option B: Manual deployment**

```bash
cd k8s

# 1. Create namespace
kubectl apply -f namespace.yaml

# 2. Create priority classes (cluster-wide)
kubectl apply -f priority-class.yaml

# 3. Create ConfigMap
kubectl apply -f configmap.yaml

# 4. Create Secrets
kubectl apply -f secrets.yaml
# OR if you generated secrets-generated.yaml:
kubectl apply -f secrets-generated.yaml

# 5. Apply resource quotas and limits
kubectl apply -f resource-quota.yaml

# 6. Apply network policies
kubectl apply -f network-policy.yaml
```

### Step 6: Verify Deployment

```bash
# Check namespace
kubectl get namespace insavein

# Check all resources
kubectl get all -n insavein

# Verify ConfigMap
kubectl describe configmap insavein-config -n insavein

# Verify Secrets (values will be hidden)
kubectl describe secret insavein-secrets -n insavein
kubectl describe secret postgres-credentials -n insavein

# Check resource quotas
kubectl describe resourcequota insavein-resource-quota -n insavein

# Check limit ranges
kubectl describe limitrange insavein-limit-range -n insavein
kubectl describe limitrange insavein-database-limits -n insavein

# Verify network policies
kubectl get networkpolicies -n insavein

# Check priority classes
kubectl get priorityclass | grep insavein
```

Expected output:
- Namespace: `insavein` (Active)
- ConfigMaps: 1 (insavein-config)
- Secrets: 2 (insavein-secrets, postgres-credentials)
- ResourceQuotas: 1 (insavein-resource-quota)
- LimitRanges: 2 (insavein-limit-range, insavein-database-limits)
- NetworkPolicies: 7 (various policies)
- PriorityClasses: 4 (critical, high, medium, low)

## Configuration Validation

### Test ConfigMap Access

```bash
# Create a test pod that uses the ConfigMap
kubectl run test-config \
  --image=busybox \
  --restart=Never \
  -n insavein \
  --rm -it \
  --env="DB_HOST=$(kubectl get configmap insavein-config -n insavein -o jsonpath='{.data.DB_HOST}')" \
  -- sh -c 'echo "DB_HOST=$DB_HOST"'
```

### Test Secret Access

```bash
# Verify secret exists and has correct keys
kubectl get secret insavein-secrets -n insavein -o jsonpath='{.data}' | jq 'keys'

# Check if JWT secret is set (without revealing value)
kubectl get secret insavein-secrets -n insavein -o jsonpath='{.data.JWT_SECRET_KEY}' | wc -c
# Should output a number > 0
```

### Test Network Policies

```bash
# Create a test pod
kubectl run test-network \
  --image=nicolaka/netshoot \
  -n insavein \
  --labels="tier=backend" \
  -- sleep 3600

# Test DNS resolution
kubectl exec -it test-network -n insavein -- nslookup kubernetes.default

# Test connectivity (should work with proper labels)
kubectl exec -it test-network -n insavein -- curl -v telnet://postgres-primary.insavein.svc.cluster.local:5432

# Clean up
kubectl delete pod test-network -n insavein
```

### Test Resource Limits

```bash
# Try to create a pod that exceeds limits (should fail)
kubectl run test-limits \
  --image=nginx \
  -n insavein \
  --requests='cpu=100,memory=20Gi' \
  --limits='cpu=100,memory=20Gi'

# Should fail with error about exceeding resource quota
```

## Troubleshooting

### Issue: Secrets contain placeholder values

**Error**: `make check-secrets` fails with placeholder warning

**Solution**:
1. Run `make generate-secrets` to generate secure values
2. Update `secrets.yaml` with generated values
3. Ensure ALL `CHANGE_ME_*` placeholders are replaced

### Issue: Resource quota exceeded

**Error**: "exceeded quota: insavein-resource-quota"

**Solution**:
1. Check current usage: `kubectl describe resourcequota -n insavein`
2. Delete unused resources or increase quota in `resource-quota.yaml`
3. Reapply: `kubectl apply -f resource-quota.yaml`

### Issue: Network policy blocking traffic

**Error**: Pods cannot communicate with each other

**Solution**:
1. Check network policies: `kubectl get networkpolicies -n insavein`
2. Verify pod labels match policy selectors
3. Test with: `kubectl describe networkpolicy <policy-name> -n insavein`
4. Temporarily disable policies for testing:
   ```bash
   kubectl delete networkpolicy insavein-default-deny -n insavein
   ```

### Issue: Priority class not found

**Error**: "priorityclass.scheduling.k8s.io not found"

**Solution**:
1. Priority classes are cluster-wide resources
2. Apply them first: `kubectl apply -f priority-class.yaml`
3. Verify: `kubectl get priorityclass | grep insavein`

### Issue: Cannot create namespace

**Error**: "namespaces is forbidden: User cannot create resource"

**Solution**:
- You need cluster admin permissions
- Contact your cluster administrator
- Or use an existing namespace (update all manifests)

## Security Best Practices

### 1. Secrets Management

**Development/Testing**:
- Use Kubernetes Secrets with base64 encoding
- Store secrets files outside version control
- Use `.gitignore` to exclude secrets files

**Production**:
- Use external secret management:
  - HashiCorp Vault
  - AWS Secrets Manager
  - Azure Key Vault
  - Google Secret Manager
- Enable encryption at rest for etcd
- Rotate secrets every 90 days
- Use RBAC to restrict secret access

### 2. Network Security

- Keep network policies enabled
- Use zero-trust networking model
- Only expose necessary ports
- Regularly audit network policies

### 3. Access Control

- Use RBAC for fine-grained permissions
- Create service accounts for applications
- Avoid using cluster admin for applications
- Regularly audit RBAC policies

### 4. Monitoring

- Monitor resource usage
- Set up alerts for quota violations
- Track secret access (audit logs)
- Monitor network policy violations

## Next Steps

After successfully deploying the base configuration:

1. **Deploy Database**:
   - PostgreSQL primary instance
   - PostgreSQL read replicas (2)
   - PgBouncer for connection pooling

2. **Deploy Backend Services**:
   - Auth Service (priority: critical)
   - User Service (priority: high)
   - Savings Service (priority: high)
   - Budget Service (priority: high)
   - Goal Service (priority: medium)
   - Education Service (priority: low)
   - Notification Service (priority: low)
   - Analytics Service (priority: medium)

3. **Deploy Frontend**:
   - TanStack Start application
   - Configure ingress controller

4. **Set Up Observability**:
   - Prometheus for metrics
   - Grafana for dashboards
   - OpenTelemetry for tracing
   - ELK/Loki for logging

5. **Configure CI/CD**:
   - GitHub Actions / GitLab CI
   - Automated testing
   - Rolling deployments
   - Automated rollbacks

## Cleanup

To remove all resources:

**Using Make**:
```bash
cd k8s
make delete
```

**Manual cleanup**:
```bash
# Delete in reverse order
kubectl delete -f network-policy.yaml
kubectl delete -f resource-quota.yaml
kubectl delete -f secrets.yaml
kubectl delete -f configmap.yaml
kubectl delete -f priority-class.yaml
kubectl delete -f namespace.yaml
```

**⚠️ WARNING**: This will delete ALL resources in the insavein namespace!

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review Kubernetes documentation
3. Check cluster logs: `kubectl logs -n insavein <pod-name>`
4. Contact your DevOps team

## References

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/)
- [Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [Secrets Management](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Priority and Preemption](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/)
