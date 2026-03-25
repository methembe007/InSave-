# Task 25: Security Implementation - Deployment Guide

This document provides comprehensive instructions for deploying the security enhancements implemented in Task 25.

## Overview

Task 25 implements critical security features across the InSavein platform:

1. **Input Validation Middleware** - Validates and sanitizes all user inputs
2. **Authorization Middleware** - Enforces resource ownership and role-based access
3. **TLS/SSL Configuration** - Enables HTTPS with Let's Encrypt certificates
4. **Security Headers** - Adds comprehensive security headers to all responses
5. **Rate Limiting** - Implements per-user and per-IP rate limits

## Prerequisites

Before deploying, ensure you have:

- [ ] Kubernetes cluster running (v1.28+)
- [ ] kubectl configured and connected to cluster
- [ ] NGINX Ingress Controller installed
- [ ] cert-manager installed (for TLS certificates)
- [ ] Domain name configured (api.insavein.com)
- [ ] DNS records pointing to cluster

## Deployment Steps

### Step 1: Install cert-manager (if not already installed)

```bash
# Add cert-manager Helm repository
helm repo add jetstack https://charts.jetstack.io
helm repo update

# Install cert-manager
kubectl create namespace cert-manager
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.13.0 \
  --set installCRDs=true

# Verify installation
kubectl get pods -n cert-manager
```

### Step 2: Deploy Certificate Issuers

```bash
# Apply ClusterIssuer and Certificate resources
kubectl apply -f k8s/cert-manager-issuer.yaml

# Verify ClusterIssuer is ready
kubectl get clusterissuer
kubectl describe clusterissuer letsencrypt-prod

# Check certificate status
kubectl get certificate -n insavein
kubectl describe certificate insavein-tls -n insavein
```

**Expected Output:**
```
NAME              READY   SECRET         AGE
insavein-tls      True    insavein-tls   2m
```

### Step 3: Apply Rate Limiting Configuration

```bash
# Apply rate limiting ConfigMaps
kubectl apply -f k8s/rate-limit-config.yaml

# Verify ConfigMaps
kubectl get configmap -n insavein | grep rate-limit
kubectl get configmap -n insavein | grep security
```

### Step 4: Update Ingress with Security Enhancements

```bash
# Apply updated Ingress configuration
kubectl apply -f k8s/ingress.yaml

# Verify Ingress is updated
kubectl get ingress -n insavein
kubectl describe ingress insavein-ingress -n insavein

# Check for TLS certificate
kubectl get secret insavein-tls -n insavein
```

### Step 5: Update Services with Validation Middleware

For each service, the validation and authorization middleware needs to be integrated. Here's the pattern:

#### 5.1: Add Dependencies

In each service's `go.mod`:

```bash
cd auth-service
go get github.com/go-playground/validator/v10
go mod tidy

cd ../user-service
go get github.com/go-playground/validator/v10
go mod tidy

# Repeat for all services
```

#### 5.2: Update Handlers

Example for savings-service:

```go
// internal/handlers/savings_handler.go
package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/go-playground/validator/v10"
    "github.com/insavein/savings-service/internal/savings"
)

type SavingsHandler struct {
    service  savings.Service
    validate *validator.Validate
}

func NewSavingsHandler(service savings.Service) *SavingsHandler {
    return &SavingsHandler{
        service:  service,
        validate: validator.New(),
    }
}

func (h *SavingsHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(string)
    if !ok {
        respondError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
    
    var req savings.CreateTransactionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    // Validate request
    if err := h.validate.Struct(req); err != nil {
        respondValidationError(w, err)
        return
    }
    
    // Sanitize inputs
    req.Description = sanitizeString(req.Description)
    req.Category = sanitizeString(req.Category)
    
    // Process request
    tx, err := h.service.CreateTransaction(r.Context(), userID, req)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    respondJSON(w, http.StatusCreated, tx)
}
```

#### 5.3: Add Validation Tags to Request Structs

```go
// internal/savings/types.go
type CreateTransactionRequest struct {
    Amount      float64 `json:"amount" validate:"required,gt=0"`
    Currency    string  `json:"currency" validate:"required,len=3"`
    Description string  `json:"description" validate:"max=500"`
    Category    string  `json:"category" validate:"required,max=50"`
}
```

### Step 6: Rebuild and Deploy Services

```bash
# Build Docker images with updated code
docker build -t insavein/auth-service:latest ./auth-service
docker build -t insavein/user-service:latest ./user-service
docker build -t insavein/savings-service:latest ./savings-service
docker build -t insavein/budget-service:latest ./budget-service
docker build -t insavein/goal-service:latest ./goal-service
docker build -t insavein/education-service:latest ./education-service
docker build -t insavein/notification-service:latest ./notification-service
docker build -t insavein/analytics-service:latest ./analytics-service

# Push to registry
docker push insavein/auth-service:latest
docker push insavein/user-service:latest
# ... repeat for all services

# Update Kubernetes deployments
kubectl rollout restart deployment/auth-service -n insavein
kubectl rollout restart deployment/user-service -n insavein
kubectl rollout restart deployment/savings-service -n insavein
kubectl rollout restart deployment/budget-service -n insavein
kubectl rollout restart deployment/goal-service -n insavein
kubectl rollout restart deployment/education-service -n insavein
kubectl rollout restart deployment/notification-service -n insavein
kubectl rollout restart deployment/analytics-service -n insavein

# Monitor rollout status
kubectl rollout status deployment/auth-service -n insavein
```

## Verification

### 1. Verify TLS/SSL Configuration

```bash
# Test HTTPS connection
curl -I https://api.insavein.com/api/auth/health

# Check certificate details
openssl s_client -connect api.insavein.com:443 -servername api.insavein.com

# Verify HSTS header
curl -I https://api.insavein.com/api/auth/health | grep -i strict-transport
```

**Expected Headers:**
```
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'; ...
```

### 2. Verify Rate Limiting

```bash
# Test per-user rate limit (100 req/min)
for i in {1..110}; do
  curl -H "Authorization: Bearer $TOKEN" https://api.insavein.com/api/savings
done

# Should return 429 after 100 requests
```

**Expected Response (after limit):**
```json
{
  "error": "Too Many Requests",
  "message": "Rate limit exceeded. Please try again later.",
  "retry_after": 60
}
```

### 3. Verify Input Validation

```bash
# Test with invalid data
curl -X POST https://api.insavein.com/api/savings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": -10,
    "currency": "INVALID",
    "description": ""
  }'
```

**Expected Response:**
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
    },
    {
      "field": "Currency",
      "tag": "len",
      "value": "3",
      "message": "Currency must be 3 characters"
    }
  ]
}
```

### 4. Verify Authorization

```bash
# Try to access another user's resource
curl -X GET https://api.insavein.com/api/goals/other-user-goal-id \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response:**
```json
{
  "error": "Forbidden",
  "message": "Forbidden: you do not have access to this resource"
}
```

### 5. Verify Security Headers

```bash
# Check all security headers
curl -I https://api.insavein.com/api/auth/health

# Verify CSP
curl -I https://api.insavein.com/api/auth/health | grep -i content-security-policy

# Verify HSTS
curl -I https://api.insavein.com/api/auth/health | grep -i strict-transport-security
```

## Monitoring

### Check Certificate Expiration

```bash
# View certificate details
kubectl get certificate insavein-tls -n insavein -o yaml

# Check renewal status
kubectl describe certificate insavein-tls -n insavein
```

### Monitor Rate Limiting

```bash
# Check NGINX Ingress logs for rate limit events
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller | grep "limiting requests"

# View rate limit metrics (if Prometheus is configured)
kubectl port-forward -n monitoring svc/prometheus 9090:9090
# Open http://localhost:9090 and query: nginx_ingress_controller_requests{status="429"}
```

### Monitor Security Events

```bash
# Check for authentication failures
kubectl logs -n insavein deployment/auth-service | grep "authentication failed"

# Check for authorization failures
kubectl logs -n insavein deployment/savings-service | grep "Forbidden"

# Check for validation errors
kubectl logs -n insavein deployment/budget-service | grep "Validation failed"
```

## Troubleshooting

### Issue: Certificate Not Issuing

```bash
# Check cert-manager logs
kubectl logs -n cert-manager deployment/cert-manager

# Check certificate status
kubectl describe certificate insavein-tls -n insavein

# Check challenge status
kubectl get challenges -n insavein
kubectl describe challenge -n insavein
```

**Common Solutions:**
- Verify DNS records point to cluster
- Check ClusterIssuer configuration
- Ensure NGINX Ingress can handle HTTP-01 challenges
- Check firewall rules allow port 80 and 443

### Issue: Rate Limiting Not Working

```bash
# Check NGINX Ingress configuration
kubectl get configmap -n ingress-nginx ingress-nginx-controller -o yaml

# Check Ingress annotations
kubectl get ingress insavein-ingress -n insavein -o yaml

# Check NGINX logs
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
```

**Common Solutions:**
- Verify rate limit annotations are applied
- Check NGINX Ingress version supports rate limiting
- Ensure ConfigMaps are mounted correctly
- Restart NGINX Ingress controller

### Issue: Validation Not Working

```bash
# Check service logs
kubectl logs -n insavein deployment/savings-service

# Verify validator is initialized
kubectl exec -n insavein deployment/savings-service -- env | grep GO
```

**Common Solutions:**
- Verify go-playground/validator is in go.mod
- Check validation tags on structs
- Ensure validator is initialized in handler
- Rebuild and redeploy service

### Issue: Authorization Failures

```bash
# Check JWT token claims
# Decode token at https://jwt.io

# Check middleware logs
kubectl logs -n insavein deployment/savings-service | grep "user_id"

# Verify context values
kubectl logs -n insavein deployment/savings-service | grep "context"
```

**Common Solutions:**
- Verify JWT token includes user_id claim
- Check auth middleware sets context correctly
- Ensure authorization middleware reads context
- Verify resource ownership logic

## Security Best Practices

### 1. Regular Certificate Renewal

cert-manager automatically renews certificates 30 days before expiration. Monitor renewal:

```bash
# Set up alert for certificate expiration
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: certificate-expiry-alert
  namespace: monitoring
data:
  alert.rules: |
    groups:
    - name: certificates
      rules:
      - alert: CertificateExpiringSoon
        expr: certmanager_certificate_expiration_timestamp_seconds - time() < 604800
        for: 1h
        labels:
          severity: warning
        annotations:
          summary: "Certificate expiring soon"
          description: "Certificate {{ \$labels.name }} expires in less than 7 days"
EOF
```

### 2. Rate Limit Tuning

Adjust rate limits based on traffic patterns:

```yaml
# For higher traffic
nginx.ingress.kubernetes.io/limit-rps: "200"  # 200 req/sec per user
nginx.ingress.kubernetes.io/global-rate-limit: "2000"  # 2000 req/min per IP

# For stricter limits
nginx.ingress.kubernetes.io/limit-rps: "50"  # 50 req/sec per user
nginx.ingress.kubernetes.io/global-rate-limit: "500"  # 500 req/min per IP
```

### 3. Security Header Updates

Keep security headers up to date with best practices:

```yaml
# Update CSP as needed
Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline'; ...

# Add new security headers
Permissions-Policy: geolocation=(), microphone=(), camera=()
Cross-Origin-Embedder-Policy: require-corp
Cross-Origin-Opener-Policy: same-origin
```

### 4. Input Validation Rules

Regularly review and update validation rules:

```go
// Add custom validators
validate.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
    validCurrencies := []string{"USD", "EUR", "GBP"}
    currency := fl.Field().String()
    for _, valid := range validCurrencies {
        if currency == valid {
            return true
        }
    }
    return false
})
```

## Requirements Satisfied

This implementation satisfies the following requirements:

- ✅ **Requirement 17.1**: Return HTTP 400 Bad Request with detailed validation errors
- ✅ **Requirement 17.2**: Include field name in error messages
- ✅ **Requirement 17.3**: Specify valid range in error messages
- ✅ **Requirement 17.5**: Validate all inputs against defined schemas
- ✅ **Requirement 17.6**: Sanitize inputs to prevent SQL injection and XSS
- ✅ **Requirement 3.5**: Only allow users to access their own data
- ✅ **Requirement 15.4**: Return HTTP 403 Forbidden for unauthorized access
- ✅ **Requirement 20.4**: Require TLS 1.3 for all connections
- ✅ **Requirement 20.5**: Include security headers (HSTS, CSP, X-Frame-Options, etc.)
- ✅ **Requirement 18.1**: Rate limit 100 requests per minute per user
- ✅ **Requirement 18.2**: Rate limit 1000 requests per minute per IP
- ✅ **Requirement 18.3**: Include rate limit headers in responses
- ✅ **Requirement 18.5**: Include retry-after header when rate limited

## Next Steps

1. **Monitor Security Metrics**: Set up Grafana dashboards for security events
2. **Penetration Testing**: Conduct security audit and penetration testing
3. **WAF Integration**: Consider adding Web Application Firewall (Cloudflare, AWS WAF)
4. **DDoS Protection**: Implement DDoS protection at edge layer
5. **Security Scanning**: Set up automated security scanning in CI/CD pipeline

## Support

For issues or questions:
- Check logs: `kubectl logs -n insavein deployment/<service-name>`
- Review documentation: `shared/middleware/README.md`
- Contact DevOps team: devops@insavein.com
