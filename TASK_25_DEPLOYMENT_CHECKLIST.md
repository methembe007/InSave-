# Task 25: Security Implementation - Deployment Checklist

Use this checklist to ensure all security features are properly deployed and configured.

## Pre-Deployment Checklist

### Environment Setup
- [ ] Kubernetes cluster is running (v1.28+)
- [ ] kubectl is configured and connected to cluster
- [ ] NGINX Ingress Controller is installed
- [ ] cert-manager is installed (v1.13.0+)
- [ ] Domain name is configured (api.insavein.com)
- [ ] DNS records point to cluster load balancer
- [ ] Namespace 'insavein' exists

### Verify Prerequisites

```bash
# Check Kubernetes connection
kubectl cluster-info

# Check NGINX Ingress
kubectl get pods -n ingress-nginx

# Check cert-manager
kubectl get pods -n cert-manager

# Check namespace
kubectl get namespace insavein
```

## Deployment Steps

### Step 1: Certificate Management (Subtask 25.5)

- [ ] Review cert-manager issuer configuration
  ```bash
  cat k8s/cert-manager-issuer.yaml
  ```

- [ ] Update email address in issuer (if needed)
  ```yaml
  email: admin@insavein.com  # Change to your email
  ```

- [ ] Apply certificate issuer
  ```bash
  kubectl apply -f k8s/cert-manager-issuer.yaml
  ```

- [ ] Verify ClusterIssuer is ready
  ```bash
  kubectl get clusterissuer
  # Expected: letsencrypt-prod READY=True
  ```

- [ ] Check certificate status
  ```bash
  kubectl get certificate -n insavein
  kubectl describe certificate insavein-tls -n insavein
  ```

- [ ] Wait for certificate to be issued (may take 2-5 minutes)
  ```bash
  kubectl wait --for=condition=ready certificate/insavein-tls -n insavein --timeout=300s
  ```

- [ ] Verify TLS secret exists
  ```bash
  kubectl get secret insavein-tls -n insavein
  ```

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Step 2: Rate Limiting Configuration (Subtask 25.7)

- [ ] Review rate limiting configuration
  ```bash
  cat k8s/rate-limit-config.yaml
  ```

- [ ] Apply rate limiting ConfigMaps
  ```bash
  kubectl apply -f k8s/rate-limit-config.yaml
  ```

- [ ] Verify ConfigMaps are created
  ```bash
  kubectl get configmap -n insavein | grep -E "rate-limit|security"
  ```

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Step 3: Update Ingress (Subtasks 25.6 & 25.7)

- [ ] Review updated ingress configuration
  ```bash
  cat k8s/ingress.yaml
  ```

- [ ] Backup current ingress
  ```bash
  kubectl get ingress insavein-ingress -n insavein -o yaml > ingress-backup.yaml
  ```

- [ ] Apply updated ingress
  ```bash
  kubectl apply -f k8s/ingress.yaml
  ```

- [ ] Verify ingress is updated
  ```bash
  kubectl get ingress insavein-ingress -n insavein
  kubectl describe ingress insavein-ingress -n insavein
  ```

- [ ] Check ingress annotations
  ```bash
  kubectl get ingress insavein-ingress -n insavein -o jsonpath='{.metadata.annotations}' | jq
  ```

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Step 4: Add Validation Middleware (Subtask 25.1)

- [ ] Run middleware application script
  ```bash
  # Linux/Mac
  chmod +x apply-security-middleware.sh
  ./apply-security-middleware.sh
  
  # Windows
  apply-security-middleware.bat
  ```

- [ ] Verify validator is added to all services
  ```bash
  grep -r "github.com/go-playground/validator" */go.mod
  ```

- [ ] Review example implementation
  ```bash
  cat savings-service/internal/handlers/savings_handler_enhanced.go
  ```

- [ ] Update each service handler to use validation
  - [ ] auth-service
  - [ ] user-service
  - [ ] savings-service
  - [ ] budget-service
  - [ ] goal-service
  - [ ] education-service
  - [ ] notification-service
  - [ ] analytics-service

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Step 5: Add Authorization Middleware (Subtask 25.3)

- [ ] Review authorization middleware
  ```bash
  cat shared/middleware/authorization.go
  ```

- [ ] Integrate authorization checks in handlers
  - [ ] Implement ownership checkers
  - [ ] Add RequireOwnership middleware
  - [ ] Add RequireSelfOrAdmin middleware
  - [ ] Add RequireRole middleware (if needed)

- [ ] Test authorization logic
  ```bash
  # Test without token (should return 401)
  curl https://api.insavein.com/api/savings
  
  # Test with invalid token (should return 401)
  curl -H "Authorization: Bearer invalid" https://api.insavein.com/api/savings
  ```

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Step 6: Rebuild and Deploy Services

- [ ] Build Docker images
  ```bash
  docker build -t insavein/auth-service:latest ./auth-service
  docker build -t insavein/user-service:latest ./user-service
  docker build -t insavein/savings-service:latest ./savings-service
  docker build -t insavein/budget-service:latest ./budget-service
  docker build -t insavein/goal-service:latest ./goal-service
  docker build -t insavein/education-service:latest ./education-service
  docker build -t insavein/notification-service:latest ./notification-service
  docker build -t insavein/analytics-service:latest ./analytics-service
  ```

- [ ] Push images to registry
  ```bash
  docker push insavein/auth-service:latest
  docker push insavein/user-service:latest
  # ... repeat for all services
  ```

- [ ] Update Kubernetes deployments
  ```bash
  kubectl rollout restart deployment/auth-service -n insavein
  kubectl rollout restart deployment/user-service -n insavein
  kubectl rollout restart deployment/savings-service -n insavein
  kubectl rollout restart deployment/budget-service -n insavein
  kubectl rollout restart deployment/goal-service -n insavein
  kubectl rollout restart deployment/education-service -n insavein
  kubectl rollout restart deployment/notification-service -n insavein
  kubectl rollout restart deployment/analytics-service -n insavein
  ```

- [ ] Monitor rollout status
  ```bash
  kubectl rollout status deployment/auth-service -n insavein
  kubectl rollout status deployment/user-service -n insavein
  # ... repeat for all services
  ```

- [ ] Verify all pods are running
  ```bash
  kubectl get pods -n insavein
  ```

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

## Verification Checklist

### Test 1: TLS/SSL Configuration

- [ ] Test HTTPS connection
  ```bash
  curl -I https://api.insavein.com/api/auth/health
  ```
  Expected: HTTP 200 or 401

- [ ] Verify certificate details
  ```bash
  openssl s_client -connect api.insavein.com:443 -servername api.insavein.com < /dev/null
  ```
  Expected: Valid certificate from Let's Encrypt

- [ ] Test HTTP to HTTPS redirect
  ```bash
  curl -I http://api.insavein.com/api/auth/health
  ```
  Expected: 301 or 308 redirect to HTTPS

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Test 2: Security Headers

- [ ] Check HSTS header
  ```bash
  curl -I https://api.insavein.com/api/auth/health | grep -i strict-transport
  ```
  Expected: `Strict-Transport-Security: max-age=31536000; includeSubDomains; preload`

- [ ] Check X-Frame-Options
  ```bash
  curl -I https://api.insavein.com/api/auth/health | grep -i x-frame
  ```
  Expected: `X-Frame-Options: DENY`

- [ ] Check X-Content-Type-Options
  ```bash
  curl -I https://api.insavein.com/api/auth/health | grep -i x-content-type
  ```
  Expected: `X-Content-Type-Options: nosniff`

- [ ] Check Content-Security-Policy
  ```bash
  curl -I https://api.insavein.com/api/auth/health | grep -i content-security
  ```
  Expected: CSP header present

- [ ] Check all security headers
  ```bash
  curl -I https://api.insavein.com/api/auth/health
  ```

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Test 3: Input Validation

- [ ] Test with negative amount
  ```bash
  curl -X POST https://api.insavein.com/api/savings \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"amount": -10, "currency": "USD", "category": "test"}'
  ```
  Expected: 400 Bad Request with validation error

- [ ] Test with invalid currency
  ```bash
  curl -X POST https://api.insavein.com/api/savings \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"amount": 10, "currency": "INVALID", "category": "test"}'
  ```
  Expected: 400 Bad Request with validation error

- [ ] Test with missing required field
  ```bash
  curl -X POST https://api.insavein.com/api/savings \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"currency": "USD"}'
  ```
  Expected: 400 Bad Request with validation error

- [ ] Verify error response format
  ```json
  {
    "error": "Bad Request",
    "message": "Validation failed",
    "errors": [
      {
        "field": "Amount",
        "tag": "gt",
        "value": "0",
        "message": "Amount must be greater than 0"
      }
    ]
  }
  ```

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Test 4: Authorization

- [ ] Test without authentication token
  ```bash
  curl https://api.insavein.com/api/savings
  ```
  Expected: 401 Unauthorized

- [ ] Test with invalid token
  ```bash
  curl -H "Authorization: Bearer invalid_token" https://api.insavein.com/api/savings
  ```
  Expected: 401 Unauthorized

- [ ] Test accessing another user's resource
  ```bash
  curl -H "Authorization: Bearer $TOKEN" \
    https://api.insavein.com/api/goals/other-user-goal-id
  ```
  Expected: 403 Forbidden

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Test 5: Rate Limiting

- [ ] Test per-user rate limit (100 req/min)
  ```bash
  for i in {1..110}; do
    curl -H "Authorization: Bearer $TOKEN" \
      https://api.insavein.com/api/savings/summary
  done
  ```
  Expected: 429 Too Many Requests after ~100 requests

- [ ] Verify 429 response format
  ```json
  {
    "error": "Too Many Requests",
    "message": "Rate limit exceeded. Please try again later.",
    "retry_after": 60
  }
  ```

- [ ] Test per-IP rate limit (1000 req/min)
  ```bash
  # Run from different machine or wait 1 minute
  for i in {1..1010}; do
    curl https://api.insavein.com/api/auth/health
  done
  ```
  Expected: 429 Too Many Requests after ~1000 requests

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Test 6: Automated Testing

- [ ] Run comprehensive test suite
  ```bash
  chmod +x test-security-implementation.sh
  export API_URL="https://api.insavein.com"
  export TEST_TOKEN="your-jwt-token-here"
  ./test-security-implementation.sh
  ```

- [ ] Review test results
  - [ ] All security header tests pass
  - [ ] TLS/SSL tests pass
  - [ ] Input validation tests pass
  - [ ] Authorization tests pass
  - [ ] Rate limiting tests pass

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

## Post-Deployment Checklist

### Monitoring Setup

- [ ] Set up certificate expiration alerts
  ```bash
  kubectl apply -f monitoring/certificate-alerts.yaml
  ```

- [ ] Configure rate limit monitoring
  ```bash
  # Check NGINX Ingress logs
  kubectl logs -n ingress-nginx deployment/ingress-nginx-controller | grep "limiting"
  ```

- [ ] Set up security event monitoring
  ```bash
  # Monitor authentication failures
  kubectl logs -n insavein deployment/auth-service | grep "authentication failed"
  
  # Monitor authorization failures
  kubectl logs -n insavein deployment/savings-service | grep "Forbidden"
  ```

- [ ] Create Grafana dashboard for security metrics
  - [ ] Rate limit events
  - [ ] Authentication failures
  - [ ] Authorization failures
  - [ ] Validation errors

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Documentation

- [ ] Update API documentation with validation rules
- [ ] Document security headers for frontend team
- [ ] Create runbook for certificate renewal issues
- [ ] Document rate limiting behavior for users
- [ ] Update deployment documentation

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

### Security Audit

- [ ] Review all validation rules
- [ ] Verify authorization checks on all endpoints
- [ ] Test edge cases and boundary conditions
- [ ] Conduct penetration testing
- [ ] Review security logs for anomalies

**Status:** ⬜ Not Started | ⏳ In Progress | ✅ Complete

---

## Rollback Plan

If issues occur during deployment:

### Rollback Ingress
```bash
kubectl apply -f ingress-backup.yaml
```

### Rollback Service Deployments
```bash
kubectl rollout undo deployment/auth-service -n insavein
kubectl rollout undo deployment/user-service -n insavein
# ... repeat for all services
```

### Remove Certificate Configuration
```bash
kubectl delete -f k8s/cert-manager-issuer.yaml
```

### Remove Rate Limiting
```bash
kubectl delete -f k8s/rate-limit-config.yaml
```

---

## Troubleshooting

### Certificate Not Issuing

**Symptoms:** Certificate status shows "Pending" or "Failed"

**Solutions:**
1. Check cert-manager logs
   ```bash
   kubectl logs -n cert-manager deployment/cert-manager
   ```

2. Check challenge status
   ```bash
   kubectl get challenges -n insavein
   kubectl describe challenge -n insavein
   ```

3. Verify DNS records
   ```bash
   nslookup api.insavein.com
   ```

4. Check firewall rules (ports 80 and 443)

---

### Rate Limiting Not Working

**Symptoms:** No 429 responses after exceeding limits

**Solutions:**
1. Check NGINX Ingress configuration
   ```bash
   kubectl get configmap -n ingress-nginx ingress-nginx-controller -o yaml
   ```

2. Verify ingress annotations
   ```bash
   kubectl get ingress insavein-ingress -n insavein -o yaml
   ```

3. Check NGINX logs
   ```bash
   kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
   ```

4. Restart NGINX Ingress
   ```bash
   kubectl rollout restart deployment/ingress-nginx-controller -n ingress-nginx
   ```

---

### Validation Not Working

**Symptoms:** Invalid data accepted, no validation errors

**Solutions:**
1. Verify validator is initialized
   ```bash
   kubectl logs -n insavein deployment/savings-service | grep validator
   ```

2. Check validation tags on structs
   ```bash
   grep -A 5 "type.*Request struct" savings-service/internal/savings/types.go
   ```

3. Verify handler calls validate.Struct()
   ```bash
   grep "validate.Struct" savings-service/internal/handlers/*.go
   ```

4. Rebuild and redeploy service
   ```bash
   docker build -t insavein/savings-service:latest ./savings-service
   docker push insavein/savings-service:latest
   kubectl rollout restart deployment/savings-service -n insavein
   ```

---

## Sign-Off

### Deployment Team

- [ ] Infrastructure Engineer: _________________ Date: _______
- [ ] Backend Developer: _________________ Date: _______
- [ ] Security Engineer: _________________ Date: _______
- [ ] DevOps Engineer: _________________ Date: _______

### Verification

- [ ] All tests passed
- [ ] Monitoring configured
- [ ] Documentation updated
- [ ] Team trained on new security features
- [ ] Rollback plan tested

### Final Approval

- [ ] Technical Lead: _________________ Date: _______
- [ ] Security Lead: _________________ Date: _______

---

## Notes

Use this section to document any issues, deviations, or special considerations during deployment:

```
[Add notes here]
```

---

## References

- Deployment Guide: `TASK_25_SECURITY_IMPLEMENTATION.md`
- Completion Summary: `TASK_25_COMPLETION_SUMMARY.md`
- Middleware Documentation: `shared/middleware/README.md`
- Test Script: `test-security-implementation.sh`
