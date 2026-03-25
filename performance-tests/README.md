# Performance Testing with k6

This directory contains k6 load test scripts for the InSavein platform.

## Prerequisites

1. Install k6: https://k6.io/docs/getting-started/installation/
2. Ensure the platform is running (locally or in staging environment)

## Test Scripts

### 1. Normal Load Test (`normal-load.js`)
Simulates typical production load with realistic user behavior.

- **Users**: Ramps up to 1,000 concurrent users
- **Duration**: 19 minutes total
- **Request Rate**: ~10 req/sec per user
- **Scenarios**: Dashboard views (30%), savings transactions (20%), spending records (20%), analytics (15%), goals (15%)
- **Performance Targets**: p95 < 500ms, p99 < 1000ms, error rate < 0.1%

**Run:**
```bash
k6 run normal-load.js
```

### 2. Peak Load Test (`peak-load.js`)
Tests system behavior under peak traffic conditions.

- **Users**: Ramps up to 10,000 concurrent users
- **Duration**: 26 minutes total
- **Request Rate**: ~50 req/sec per user
- **Target**: 100,000 requests per minute
- **Performance Targets**: p95 < 500ms, p99 < 1000ms, error rate < 0.1%

**Run:**
```bash
k6 run peak-load.js
```

### 3. Stress Test (`stress-test.js`)
Gradually increases load to find the system's breaking point.

- **Users**: Ramps up to 50,000 concurrent users
- **Duration**: 25 minutes total
- **Purpose**: Identify bottlenecks and failure modes
- **Relaxed Targets**: p95 < 2000ms, p99 < 5000ms, error rate < 5%

**Run:**
```bash
k6 run stress-test.js
```

## Configuration

Set the base URL using environment variable:

```bash
# Local testing
k6 run -e BASE_URL=http://localhost:8080 normal-load.js

# Staging environment
k6 run -e BASE_URL=https://staging.insavein.com peak-load.js
```

## Analyzing Results

### Key Metrics to Monitor

1. **Response Time**
   - p95 latency should be < 500ms
   - p99 latency should be < 1000ms
   - Watch for degradation as load increases

2. **Error Rate**
   - Should remain < 0.1% under normal/peak load
   - May increase to 5% under stress test

3. **Throughput**
   - Target: 100,000 requests per minute (1,666 req/sec)
   - Monitor actual vs. target throughput

4. **Resource Utilization**
   - Monitor CPU, memory, database connections
   - Check for bottlenecks in specific services

### k6 Cloud (Optional)

For advanced analytics, use k6 Cloud:

```bash
k6 cloud normal-load.js
```

## Integration with CI/CD

Add performance tests to your CI/CD pipeline:

```yaml
# .github/workflows/performance-test.yml
name: Performance Tests
on:
  schedule:
    - cron: '0 2 * * *'  # Run daily at 2 AM
  workflow_dispatch:

jobs:
  performance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install k6
        run: |
          sudo gpg -k
          sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6
      - name: Run normal load test
        run: k6 run -e BASE_URL=${{ secrets.STAGING_URL }} performance-tests/normal-load.js
```

## Troubleshooting

### High Error Rates
- Check service logs for errors
- Verify database connection pool settings
- Check for rate limiting issues

### Slow Response Times
- Review database query performance
- Check for N+1 query problems
- Verify caching is working
- Monitor database replication lag

### Connection Timeouts
- Increase connection pool size
- Check network latency
- Verify load balancer configuration

## Next Steps After Testing

1. **Analyze Results**: Review metrics and identify bottlenecks
2. **Optimize**: Implement caching, query optimization, connection pooling
3. **Re-test**: Run tests again to verify improvements
4. **Monitor**: Set up continuous monitoring in production
