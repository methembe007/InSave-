# Performance Test Analysis Guide

This guide helps you run load tests and analyze the results to identify bottlenecks and optimization opportunities.

## Running Load Tests

### 1. Prepare the Environment

**Local Testing:**
```bash
# Start all services
docker-compose up -d

# Wait for services to be healthy
docker-compose ps

# Set base URL
export BASE_URL=http://localhost:8080
```

**Staging Environment:**
```bash
export BASE_URL=https://staging.insavein.com
```

### 2. Run Tests

```bash
cd performance-tests

# Run normal load test (recommended first)
./run-tests.sh normal

# Run peak load test
./run-tests.sh peak

# Run stress test (to find breaking point)
./run-tests.sh stress
```

### 3. Monitor During Tests

Open separate terminals to monitor:

**Service Metrics:**
```bash
# Watch Prometheus metrics
kubectl port-forward svc/prometheus 9090:9090
# Open http://localhost:9090

# Watch Grafana dashboards
kubectl port-forward svc/grafana 3000:3000
# Open http://localhost:3000
```

**Database Performance:**
```bash
# Monitor PostgreSQL
kubectl exec -it postgres-0 -- psql -U insavein -c "
SELECT 
  datname,
  numbackends as connections,
  xact_commit as commits,
  xact_rollback as rollbacks,
  blks_read as disk_reads,
  blks_hit as cache_hits
FROM pg_stat_database 
WHERE datname = 'insavein';
"

# Check active queries
kubectl exec -it postgres-0 -- psql -U insavein -c "
SELECT pid, usename, state, query_start, query 
FROM pg_stat_activity 
WHERE state != 'idle' 
ORDER BY query_start;
"
```

**Resource Usage:**
```bash
# Monitor pod resources
kubectl top pods -n insavein

# Watch pod status
watch kubectl get pods -n insavein
```

## Analyzing Results

### Key Performance Indicators (KPIs)

#### 1. Response Time Metrics

**Target Requirements:**
- p95 latency < 500ms
- p99 latency < 1000ms

**Analysis:**
```bash
# Extract response time metrics from k6 output
cat results/normal-load_*.json | jq '.metrics.http_req_duration'
```

**What to Look For:**
- ✅ **Good**: p95 < 500ms, p99 < 1000ms
- ⚠️ **Warning**: p95 500-800ms, p99 1000-2000ms
- ❌ **Critical**: p95 > 800ms, p99 > 2000ms

**Common Causes of Slow Response Times:**
- Database queries without indexes
- N+1 query problems
- Missing caching
- Insufficient connection pool size
- CPU/memory constraints

#### 2. Error Rate

**Target Requirements:**
- Error rate < 0.1% (1 error per 1000 requests)

**Analysis:**
```bash
# Check error rate
cat results/normal-load_*.json | jq '.metrics.http_req_failed'
```

**What to Look For:**
- ✅ **Good**: < 0.1% errors
- ⚠️ **Warning**: 0.1-1% errors
- ❌ **Critical**: > 1% errors

**Common Causes of Errors:**
- Database connection pool exhausted
- Service timeouts
- Rate limiting triggered
- Memory exhaustion
- Deadlocks or lock timeouts

#### 3. Throughput

**Target Requirements:**
- 100,000 requests per minute (1,666 req/sec)

**Analysis:**
```bash
# Check request rate
cat results/peak-load_*.json | jq '.metrics.http_reqs.values.rate'
```

**What to Look For:**
- ✅ **Good**: > 1,666 req/sec sustained
- ⚠️ **Warning**: 1,000-1,666 req/sec
- ❌ **Critical**: < 1,000 req/sec

#### 4. Resource Utilization

**CPU Usage:**
```bash
# Check CPU usage per service
kubectl top pods -n insavein --sort-by=cpu
```

**What to Look For:**
- ✅ **Good**: 40-70% CPU usage
- ⚠️ **Warning**: 70-85% CPU usage
- ❌ **Critical**: > 85% CPU usage (triggers HPA scaling)

**Memory Usage:**
```bash
# Check memory usage per service
kubectl top pods -n insavein --sort-by=memory
```

**What to Look For:**
- ✅ **Good**: < 70% of memory limit
- ⚠️ **Warning**: 70-85% of memory limit
- ❌ **Critical**: > 85% of memory limit (risk of OOM)

### Bottleneck Identification

#### Database Bottlenecks

**Symptoms:**
- High p95/p99 latencies
- Increasing response times under load
- Database CPU at 100%

**Diagnosis:**
```sql
-- Find slow queries
SELECT 
  query,
  calls,
  total_exec_time / 1000 as total_time_sec,
  mean_exec_time / 1000 as mean_time_ms,
  max_exec_time / 1000 as max_time_ms
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 20;

-- Check for missing indexes
SELECT 
  schemaname,
  tablename,
  seq_scan,
  seq_tup_read,
  idx_scan,
  seq_tup_read / seq_scan as avg_seq_tup
FROM pg_stat_user_tables
WHERE seq_scan > 0
ORDER BY seq_tup_read DESC
LIMIT 20;

-- Check connection pool usage
SELECT count(*) as total_connections,
       count(*) FILTER (WHERE state = 'active') as active,
       count(*) FILTER (WHERE state = 'idle') as idle
FROM pg_stat_activity;
```

**Solutions:**
- Add missing indexes (see task 28.3)
- Optimize slow queries
- Increase connection pool size
- Implement caching (see task 28.4)

#### Service Bottlenecks

**Symptoms:**
- Specific service has high CPU/memory
- Errors from one service
- Uneven load distribution

**Diagnosis:**
```bash
# Check service-specific metrics in Prometheus
# Query: rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# Check error rates per service
# Query: rate(http_requests_total{status=~"5.."}[5m])

# Check pod distribution
kubectl get pods -n insavein -o wide
```

**Solutions:**
- Scale up specific service (increase HPA max replicas)
- Optimize service code
- Add caching
- Review resource limits

#### Network Bottlenecks

**Symptoms:**
- High latency but low CPU/memory
- Timeouts
- Connection refused errors

**Diagnosis:**
```bash
# Check ingress metrics
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller

# Check service endpoints
kubectl get endpoints -n insavein

# Test internal service connectivity
kubectl run -it --rm debug --image=nicolaka/netshoot --restart=Never -- curl http://savings-service:8080/health
```

**Solutions:**
- Increase ingress replicas
- Adjust connection timeouts
- Review network policies
- Check DNS resolution

### Performance Optimization Checklist

After identifying bottlenecks, follow this optimization sequence:

#### Phase 1: Quick Wins (Task 28.3 - Database Optimization)
- [ ] Add missing database indexes
- [ ] Fix N+1 queries
- [ ] Enable query result caching
- [ ] Optimize slow queries

#### Phase 2: Caching (Task 28.4)
- [ ] Implement Redis for session caching
- [ ] Cache financial health scores
- [ ] Cache education content
- [ ] Add cache invalidation logic

#### Phase 3: Connection Management (Task 28.5)
- [ ] Configure PgBouncer
- [ ] Set connection pool limits (max 20 per service)
- [ ] Tune connection timeouts
- [ ] Monitor connection usage

#### Phase 4: Scaling
- [ ] Adjust HPA settings based on load patterns
- [ ] Increase resource limits if needed
- [ ] Add more database replicas for read-heavy workloads
- [ ] Consider database sharding for extreme scale

### Reporting Results

Create a performance test report with:

1. **Executive Summary**
   - Test date and environment
   - Pass/fail against requirements
   - Key findings

2. **Metrics Summary**
   - Response times (p50, p95, p99)
   - Error rate
   - Throughput
   - Resource utilization

3. **Bottlenecks Identified**
   - Description of each bottleneck
   - Impact on performance
   - Recommended solutions

4. **Optimization Plan**
   - Prioritized list of optimizations
   - Expected impact
   - Implementation effort

5. **Next Steps**
   - Immediate actions
   - Long-term improvements
   - Re-test schedule

### Example Report Template

```markdown
# Performance Test Report - InSavein Platform

**Date:** 2024-01-15
**Environment:** Staging
**Test Type:** Peak Load Test

## Executive Summary
- ✅ p95 latency: 450ms (target: <500ms)
- ❌ p99 latency: 1,200ms (target: <1000ms)
- ✅ Error rate: 0.05% (target: <0.1%)
- ⚠️ Throughput: 1,400 req/sec (target: 1,666 req/sec)

**Overall Status:** PASS with optimizations needed

## Key Findings

### 1. Analytics Service Slow Queries
- Financial health score calculation takes 800-1500ms
- Missing index on savings_transactions(user_id, created_at)
- No caching implemented

### 2. Database Connection Pool Exhaustion
- Peak connections: 95/100
- Services experiencing connection timeouts
- Need connection pooling with PgBouncer

### 3. Savings Service Memory Pressure
- Memory usage: 850MB/1GB (85%)
- Potential memory leak in streak calculation
- Needs investigation

## Optimization Plan

### High Priority
1. Add database indexes (2 hours)
2. Implement Redis caching for financial health scores (4 hours)
3. Configure PgBouncer (2 hours)

### Medium Priority
4. Optimize streak calculation algorithm (4 hours)
5. Increase savings-service memory limit to 2GB (1 hour)

### Low Priority
6. Add query result caching (8 hours)
7. Implement database query optimization (16 hours)

## Next Steps
1. Implement high-priority optimizations
2. Re-run peak load test
3. Verify p99 latency < 1000ms
4. Schedule stress test for next week
```

## Continuous Performance Monitoring

After initial optimization:

1. **Set up automated performance tests**
   - Run daily in staging
   - Alert on regression

2. **Monitor production metrics**
   - Set up Grafana alerts
   - Track p95/p99 latencies
   - Monitor error rates

3. **Regular performance reviews**
   - Weekly review of metrics
   - Monthly load testing
   - Quarterly capacity planning
