# InSavein Platform Operations Runbook

Operational procedures for monitoring, troubleshooting, and maintaining the InSavein platform.

## Table of Contents

1. [Monitoring and Alerting](#monitoring-and-alerting)
2. [Common Issues and Resolutions](#common-issues-and-resolutions)
3. [Backup and Restore](#backup-and-restore)
4. [Scaling Procedures](#scaling-procedures)
5. [Incident Response](#incident-response)
6. [Maintenance Tasks](#maintenance-tasks)

---

## Monitoring and Alerting

### Monitoring Stack

**Components**:
- **Prometheus**: Metrics collection and alerting
- **Grafana**: Visualization and dashboards
- **OpenTelemetry**: Distributed tracing
- **ELK/Loki**: Log aggregation

### Access Monitoring Tools

```bash
# Port-forward Prometheus
kubectl port-forward svc/prometheus 9090:9090 -n insavein
# Access: http://localhost:9090

# Port-forward Grafana
kubectl port-forward svc/grafana 3000:3000 -n insavein
# Access: http://localhost:3000
# Default credentials: admin / <from secret>
```

### Key Metrics to Monitor

**Service Health**:
- `up{job="<service-name>"}` - Service availability (1 = up, 0 = down)
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request latency (p50, p95, p99)
- `http_requests_errors_total` - Error count

**System Resources**:
- `container_cpu_usage_seconds_total` - CPU usage
- `container_memory_usage_bytes` - Memory usage
- `container_network_receive_bytes_total` - Network ingress
- `container_network_transmit_bytes_total` - Network egress

**Database**:
- `pg_up` - PostgreSQL availability
- `pg_stat_database_numbackends` - Active connections
- `pg_stat_replication_lag_bytes` - Replication lag
- `pg_stat_database_tup_inserted` - Insert rate
- `pg_stat_database_tup_updated` - Update rate

**Business Metrics**:
- `savings_transactions_total` - Total savings transactions
- `user_registrations_total` - New user registrations
- `budget_alerts_generated_total` - Budget alerts generated
- `goal_completions_total` - Goals completed

### Grafana Dashboards

Pre-configured dashboards available:

1. **Platform Overview** - High-level system health
2. **Service Metrics** - Per-service performance
3. **Database Performance** - PostgreSQL metrics
4. **Business Metrics** - User activity and engagement
5. **Infrastructure** - Node and pod resources

Import dashboards:
```bash
kubectl apply -f k8s/grafana-dashboards.yaml
```

### Alert Rules

**Critical Alerts** (PagerDuty/Slack):
- Service down (no healthy pods)
- Database connection failures
- High error rate (>1%)
- Replication lag (>10s)
- Disk space critical (<10%)

**Warning Alerts** (Slack/Email):
- High CPU usage (>80%)
- High memory usage (>85%)
- Slow response time (p95 >500ms)
- Replication lag (>5s)
- Disk space warning (<20%)

**Alert Configuration**:
```bash
kubectl apply -f k8s/prometheus-alerts.yaml
```

### Alert Response Times

- **Critical**: Respond within 15 minutes
- **Warning**: Respond within 1 hour
- **Info**: Review during business hours

---

## Common Issues and Resolutions

### Issue 1: Service Unavailable (503)

**Symptoms**:
- HTTP 503 errors
- Health checks failing
- No healthy pods

**Diagnosis**:
```bash
# Check pod status
kubectl get pods -l app=<service-name> -n insavein

# Check pod events
kubectl describe pod <pod-name> -n insavein

# Check logs
kubectl logs <pod-name> -n insavein --tail=100

# Check resource usage
kubectl top pod <pod-name> -n insavein
```

**Common Causes & Solutions**:

1. **Database connection failure**:
   ```bash
   # Check database pods
   kubectl get pods -l app=postgres-primary -n insavein
   
   # Test connection from service pod
   kubectl exec -it <service-pod> -n insavein -- \
     psql -h postgres-primary -U insavein_user -d insavein -c "SELECT 1"
   
   # Solution: Restart database or check credentials
   kubectl rollout restart statefulset/postgres-primary -n insavein
   ```

2. **Out of memory (OOMKilled)**:
   ```bash
   # Check if pod was OOMKilled
   kubectl describe pod <pod-name> -n insavein | grep -i oom
   
   # Solution: Increase memory limits
   kubectl edit deployment <service-name> -n insavein
   # Update resources.limits.memory
   ```

3. **Image pull failure**:
   ```bash
   # Check image pull status
   kubectl describe pod <pod-name> -n insavein | grep -i image
   
   # Solution: Verify image exists and credentials are correct
   kubectl get secret -n insavein
   ```

### Issue 2: High Response Time

**Symptoms**:
- p95 latency >500ms
- p99 latency >1000ms
- Slow page loads

**Diagnosis**:
```bash
# Check service metrics in Prometheus
# Query: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Check database query performance
kubectl exec -it postgres-primary-0 -n insavein -- \
  psql -U insavein_user -d insavein -c \
  "SELECT query, mean_exec_time, calls FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;"

# Check pod resource usage
kubectl top pods -n insavein
```

**Solutions**:

1. **Slow database queries**:
   ```sql
   -- Add missing indexes
   CREATE INDEX CONCURRENTLY idx_savings_user_date 
   ON savings_transactions(user_id, created_at DESC);
   
   -- Analyze query plans
   EXPLAIN ANALYZE SELECT * FROM savings_transactions WHERE user_id = '...';
   ```

2. **High CPU usage**:
   ```bash
   # Scale up replicas
   kubectl scale deployment <service-name> --replicas=5 -n insavein
   
   # Or enable HPA
   kubectl autoscale deployment <service-name> \
     --cpu-percent=70 --min=3 --max=10 -n insavein
   ```

3. **Cache misses**:
   ```bash
   # Check Redis cache (if deployed)
   kubectl exec -it redis-0 -n insavein -- redis-cli INFO stats
   
   # Increase cache TTL or warm cache
   ```

### Issue 3: Database Replication Lag

**Symptoms**:
- Replication lag >5 seconds
- Stale data on read replicas
- Alert: "High replication lag"

**Diagnosis**:
```bash
# Check replication status
kubectl exec -it postgres-primary-0 -n insavein -- \
  psql -U postgres -c "SELECT * FROM pg_stat_replication;"

# Check replica lag
kubectl exec -it postgres-replica1-0 -n insavein -- \
  psql -U postgres -c "SELECT now() - pg_last_xact_replay_timestamp() AS lag;"
```

**Solutions**:

1. **High write load**:
   ```bash
   # Reduce write load or scale primary
   # Check write rate
   kubectl exec -it postgres-primary-0 -n insavein -- \
     psql -U postgres -c "SELECT tup_inserted, tup_updated FROM pg_stat_database WHERE datname='insavein';"
   ```

2. **Network issues**:
   ```bash
   # Check network connectivity
   kubectl exec -it postgres-replica1-0 -n insavein -- \
     ping postgres-primary.insavein.svc.cluster.local
   
   # Check network policies
   kubectl get networkpolicies -n insavein
   ```

3. **Replica overloaded**:
   ```bash
   # Check replica resource usage
   kubectl top pod postgres-replica1-0 -n insavein
   
   # Add more replicas or increase resources
   ```

### Issue 4: High Error Rate

**Symptoms**:
- Error rate >0.1%
- 4xx or 5xx HTTP errors
- Alert: "High error rate"

**Diagnosis**:
```bash
# Check error logs
kubectl logs -l app=<service-name> -n insavein | grep -i error

# Check error metrics in Prometheus
# Query: rate(http_requests_errors_total[5m])

# Check specific error types
kubectl logs <pod-name> -n insavein | grep "HTTP 500" | tail -20
```

**Solutions**:

1. **Validation errors (400)**:
   - Review recent code changes
   - Check API client versions
   - Verify request payloads

2. **Authentication errors (401)**:
   ```bash
   # Check JWT secret consistency
   kubectl get secret insavein-secrets -n insavein -o jsonpath='{.data.JWT_SECRET_KEY}' | base64 -d
   
   # Verify token expiration settings
   kubectl get configmap insavein-config -n insavein -o yaml | grep TOKEN
   ```

3. **Internal errors (500)**:
   - Check application logs for stack traces
   - Review recent deployments
   - Consider rollback if errors started after deployment

### Issue 5: Disk Space Full

**Symptoms**:
- Alert: "Disk space critical"
- Database write failures
- Pod evictions

**Diagnosis**:
```bash
# Check PVC usage
kubectl get pvc -n insavein

# Check disk usage in pod
kubectl exec -it postgres-primary-0 -n insavein -- df -h

# Check what's using space
kubectl exec -it postgres-primary-0 -n insavein -- du -sh /var/lib/postgresql/data/*
```

**Solutions**:

1. **Clean up old data**:
   ```sql
   -- Archive old partitions
   -- See migrations/partition-management/archive_old_partitions.sql
   
   -- Vacuum database
   VACUUM FULL;
   ```

2. **Expand PVC**:
   ```bash
   # Edit PVC (if storage class supports expansion)
   kubectl edit pvc postgres-primary-data -n insavein
   # Increase storage size
   ```

3. **Clean up logs**:
   ```bash
   # Rotate logs
   kubectl exec -it <pod-name> -n insavein -- \
     find /var/log -name "*.log" -mtime +7 -delete
   ```

---

## Backup and Restore

### Automated Backups

**Schedule**: Daily at 2 AM UTC

**Backup Types**:
- **Full backup**: Daily
- **Incremental backup**: Every 6 hours
- **WAL archiving**: Continuous

**Backup Location**: S3/GCS/Azure Blob Storage

### Manual Backup

```bash
# Full database backup
kubectl exec -it postgres-primary-0 -n insavein -- \
  pg_dump -U insavein_user -d insavein -F c -f /tmp/backup.dump

# Copy backup from pod
kubectl cp insavein/postgres-primary-0:/tmp/backup.dump ./backup-$(date +%Y%m%d).dump

# Upload to cloud storage
aws s3 cp backup-$(date +%Y%m%d).dump s3://insavein-backups/
```

### Restore from Backup

**⚠️ WARNING**: This will overwrite existing data!

```bash
# 1. Download backup
aws s3 cp s3://insavein-backups/backup-20260115.dump ./backup.dump

# 2. Copy to pod
kubectl cp ./backup.dump insavein/postgres-primary-0:/tmp/backup.dump

# 3. Stop services (prevent writes)
kubectl scale deployment --all --replicas=0 -n insavein

# 4. Restore database
kubectl exec -it postgres-primary-0 -n insavein -- \
  pg_restore -U insavein_user -d insavein -c /tmp/backup.dump

# 5. Verify data
kubectl exec -it postgres-primary-0 -n insavein -- \
  psql -U insavein_user -d insavein -c "SELECT COUNT(*) FROM users;"

# 6. Restart services
kubectl scale deployment --all --replicas=3 -n insavein
```

### Point-in-Time Recovery (PITR)

```bash
# Restore to specific timestamp
# Requires WAL archiving enabled

# 1. Stop PostgreSQL
kubectl scale statefulset postgres-primary --replicas=0 -n insavein

# 2. Restore base backup
# 3. Configure recovery.conf with target time
# 4. Start PostgreSQL
kubectl scale statefulset postgres-primary --replicas=1 -n insavein

# 5. Verify recovery
```

### Backup Verification

**Monthly Test**:
```bash
# 1. Restore to test environment
# 2. Run smoke tests
# 3. Verify data integrity
# 4. Document results
```

---

## Scaling Procedures

### Horizontal Scaling (Add Replicas)

**Manual Scaling**:
```bash
# Scale specific service
kubectl scale deployment auth-service --replicas=5 -n insavein

# Scale all services
kubectl scale deployment --all --replicas=5 -n insavein

# Verify scaling
kubectl get pods -n insavein
```

**Auto-Scaling (HPA)**:
```bash
# Enable HPA
kubectl autoscale deployment auth-service \
  --cpu-percent=70 \
  --min=3 \
  --max=10 \
  -n insavein

# Check HPA status
kubectl get hpa -n insavein

# Describe HPA
kubectl describe hpa auth-service -n insavein
```

### Vertical Scaling (Increase Resources)

```bash
# Edit deployment
kubectl edit deployment auth-service -n insavein

# Update resources:
# resources:
#   requests:
#     cpu: 500m
#     memory: 512Mi
#   limits:
#     cpu: 2000m
#     memory: 2Gi

# Or use kubectl set
kubectl set resources deployment auth-service \
  --requests=cpu=500m,memory=512Mi \
  --limits=cpu=2000m,memory=2Gi \
  -n insavein
```

### Database Scaling

**Add Read Replica**:
```bash
# Scale replica StatefulSet
kubectl scale statefulset postgres-replica --replicas=3 -n insavein

# Verify replication
kubectl exec -it postgres-primary-0 -n insavein -- \
  psql -U postgres -c "SELECT * FROM pg_stat_replication;"
```

**Increase Database Resources**:
```bash
# Edit StatefulSet
kubectl edit statefulset postgres-primary -n insavein

# Update resources and restart
kubectl rollout restart statefulset postgres-primary -n insavein
```

### Scaling Checklist

- [ ] Check current resource usage
- [ ] Determine scaling target (CPU, memory, requests)
- [ ] Plan scaling action (horizontal or vertical)
- [ ] Execute scaling
- [ ] Monitor metrics during scaling
- [ ] Verify health checks pass
- [ ] Test application functionality
- [ ] Document scaling action

---

## Incident Response

### Incident Severity Levels

**SEV1 - Critical**:
- Complete service outage
- Data loss or corruption
- Security breach
- Response time: 15 minutes

**SEV2 - High**:
- Partial service degradation
- High error rate (>5%)
- Performance severely impacted
- Response time: 1 hour

**SEV3 - Medium**:
- Minor service degradation
- Non-critical feature unavailable
- Performance moderately impacted
- Response time: 4 hours

**SEV4 - Low**:
- Cosmetic issues
- Minor bugs
- Documentation errors
- Response time: Next business day

### Incident Response Process

1. **Detect**: Alert triggered or user report
2. **Acknowledge**: On-call engineer acknowledges
3. **Assess**: Determine severity and impact
4. **Communicate**: Notify stakeholders
5. **Mitigate**: Implement temporary fix
6. **Resolve**: Implement permanent fix
7. **Document**: Write post-mortem
8. **Learn**: Implement preventive measures

### Incident Communication

**Status Page**: Update status.insavein.com

**Slack Channels**:
- `#incidents` - Incident coordination
- `#engineering` - Technical discussion
- `#customer-support` - User communication

**Email Templates**:
```
Subject: [INCIDENT] InSavein Service Disruption

We are currently experiencing issues with [service/feature].

Impact: [description]
Status: Investigating/Mitigating/Resolved
ETA: [time]

We will provide updates every 30 minutes.
```

### Post-Incident Review

**Template**:
```markdown
# Post-Incident Review: [Title]

## Incident Summary
- Date: YYYY-MM-DD
- Duration: X hours
- Severity: SEVX
- Impact: [description]

## Timeline
- HH:MM - Alert triggered
- HH:MM - Engineer acknowledged
- HH:MM - Root cause identified
- HH:MM - Mitigation deployed
- HH:MM - Incident resolved

## Root Cause
[Detailed explanation]

## Resolution
[What was done to fix]

## Action Items
- [ ] Improve monitoring
- [ ] Add automated tests
- [ ] Update runbook
- [ ] Implement preventive measures

## Lessons Learned
[What we learned]
```

---

## Maintenance Tasks

### Daily Tasks

- [ ] Check monitoring dashboards
- [ ] Review error logs
- [ ] Check backup status
- [ ] Monitor resource usage
- [ ] Review alerts

### Weekly Tasks

- [ ] Review performance metrics
- [ ] Check database performance
- [ ] Review security logs
- [ ] Update documentation
- [ ] Team sync meeting

### Monthly Tasks

- [ ] Test backup restoration
- [ ] Review and update alerts
- [ ] Capacity planning review
- [ ] Security audit
- [ ] Update dependencies
- [ ] Rotate secrets (every 90 days)

### Quarterly Tasks

- [ ] Disaster recovery drill
- [ ] Performance testing
- [ ] Security penetration testing
- [ ] Infrastructure cost review
- [ ] Team training

### Secret Rotation

```bash
# Generate new secrets
openssl rand -base64 64 > new-jwt-secret.txt

# Update Kubernetes secret
kubectl create secret generic insavein-secrets \
  --from-literal=JWT_SECRET_KEY=$(cat new-jwt-secret.txt) \
  --dry-run=client -o yaml | kubectl apply -f -

# Rolling restart services to pick up new secret
kubectl rollout restart deployment --all -n insavein

# Verify services are healthy
kubectl get pods -n insavein
```

### Database Maintenance

```sql
-- Vacuum and analyze (weekly)
VACUUM ANALYZE;

-- Reindex (monthly)
REINDEX DATABASE insavein;

-- Update statistics
ANALYZE;

-- Check for bloat
SELECT schemaname, tablename, 
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### Log Rotation

```bash
# Configure log rotation in Kubernetes
# Logs are automatically rotated by container runtime

# Manual cleanup if needed
kubectl exec -it <pod-name> -n insavein -- \
  find /var/log -name "*.log" -mtime +30 -delete
```

---

## Emergency Contacts

**On-Call Rotation**: See PagerDuty schedule

**Escalation Path**:
1. On-call engineer
2. Engineering lead
3. CTO
4. CEO (SEV1 only)

**External Contacts**:
- Cloud provider support
- Database vendor support
- Security team

---

## Additional Resources

- [Monitoring Dashboards](http://grafana.insavein.com)
- [Alert Manager](http://prometheus.insavein.com)
- [Status Page](https://status.insavein.com)
- [Internal Wiki](https://wiki.insavein.com)
- [Incident Management Tool](https://pagerduty.com)

---

**Last Updated**: 2026-01-15  
**Version**: 1.0.0
