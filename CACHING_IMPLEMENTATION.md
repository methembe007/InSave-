# Redis Caching Implementation Guide

## Overview

Task 28.4 implements a comprehensive Redis caching strategy to improve performance and reduce database load.

## Caching Strategy

### Cache TTL Configuration

| Data Type | TTL | Rationale |
|-----------|-----|-----------|
| Sessions | 15 minutes | Matches JWT access token expiry |
| Financial Health Scores | 1 hour | Expensive calculation, changes slowly |
| Education Content | 24 hours | Static content, rarely changes |
| User Profiles | 30 minutes | Moderate update frequency |
| Budgets | 15 minutes | Updated frequently during spending |
| Savings Summary | 5 minutes | Updated with each transaction |
| Goals | 10 minutes | Moderate update frequency |

### Cache Invalidation Strategy

**Write-Through Pattern:**
- Update database first
- Invalidate cache on success
- Next read will populate cache

**Cache-Aside Pattern:**
- Check cache first
- On miss, query database
- Store result in cache

## Implementation

### 1. Install Redis in Kubernetes

```bash
# Apply Redis deployment
kubectl apply -f k8s/redis-deployment.yaml

# Verify Redis is running
kubectl get pods -n insavein | grep redis

# Test Redis connection
kubectl exec -it redis-0 -n insavein -- redis-cli -a <password> ping
```

### 2. Add Redis Client to Services

**Add dependency to go.mod:**
```go
require (
    github.com/redis/go-redis/v9 v9.3.0
)
```

**Initialize cache in service:**
```go
import (
    "github.com/insavein/shared/cache"
)

func main() {
    // Initialize Redis cache
    redisCache, err := cache.NewRedisCache(
        os.Getenv("REDIS_ADDR"),
        os.Getenv("REDIS_PASSWORD"),
        0, // database number
    )
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    defer redisCache.Close()

    // Create cache manager
    cacheManager := cache.NewCacheManager(redisCache, nil)

    // Pass to service
    service := NewService(db, cacheManager)
}
```

### 3. Update Service Code

#### Analytics Service - Financial Health Score Caching

**Before (no caching):**
```go
func (s *Service) GetFinancialHealth(ctx context.Context, userID string) (*FinancialHealthScore, error) {
    // Expensive calculation
    score, err := s.calculateFinancialHealth(ctx, userID)
    if err != nil {
        return nil, err
    }
    return score, nil
}
```

**After (with caching):**
```go
func (s *Service) GetFinancialHealth(ctx context.Context, userID string) (*FinancialHealthScore, error) {
    // Try cache first
    var score FinancialHealthScore
    err := s.cacheManager.GetFinancialHealth(ctx, userID)
    if err == nil {
        return &score, nil
    }
    
    // Cache miss - calculate
    calculatedScore, err := s.calculateFinancialHealth(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Store in cache (fire and forget)
    go s.cacheManager.SetFinancialHealth(context.Background(), userID, calculatedScore)
    
    return calculatedScore, nil
}

// Invalidate cache when data changes
func (s *Service) RecordTransaction(ctx context.Context, userID string, tx *Transaction) error {
    if err := s.repo.CreateTransaction(ctx, tx); err != nil {
        return err
    }
    
    // Invalidate financial health cache
    go s.cacheManager.InvalidateFinancialHealth(context.Background(), userID)
    
    return nil
}
```

#### Education Service - Content Caching

```go
func (s *Service) GetLesson(ctx context.Context, lessonID string) (*LessonDetail, error) {
    // Try cache first
    var lesson LessonDetail
    cached, err := s.cacheManager.GetLesson(ctx, lessonID)
    if err == nil {
        return cached.(*LessonDetail), nil
    }
    
    // Cache miss - query database
    lesson, err := s.repo.GetLesson(ctx, lessonID)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    go s.cacheManager.SetLesson(context.Background(), lessonID, lesson)
    
    return lesson, nil
}

func (s *Service) GetLessons(ctx context.Context, category string) ([]Lesson, error) {
    // Try cache first
    cached, err := s.cacheManager.GetLessonList(ctx, category)
    if err == nil {
        return cached.([]Lesson), nil
    }
    
    // Cache miss - query database
    lessons, err := s.repo.GetLessons(ctx, category)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    go s.cacheManager.SetLessonList(context.Background(), category, lessons)
    
    return lessons, nil
}
```

#### Auth Service - Session Caching

```go
func (s *Service) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
    // Parse token to get session ID
    claims, err := s.parseToken(token)
    if err != nil {
        return nil, err
    }
    
    // Check if session is cached (not revoked)
    sessionID := claims.SessionID
    exists, err := s.cacheManager.GetSession(ctx, sessionID)
    if err != nil && !cache.IsCacheMiss(err) {
        // Cache error - fall back to database
        return s.validateTokenFromDB(ctx, token)
    }
    
    if exists != nil {
        // Session found in cache - valid
        return claims, nil
    }
    
    // Session not in cache - check if revoked
    revoked, err := s.repo.IsTokenRevoked(ctx, sessionID)
    if err != nil {
        return nil, err
    }
    
    if revoked {
        return nil, ErrTokenRevoked
    }
    
    // Valid session - cache it
    go s.cacheManager.SetSession(context.Background(), sessionID, map[string]interface{}{
        "user_id": claims.UserID,
        "email":   claims.Email,
    })
    
    return claims, nil
}

func (s *Service) Logout(ctx context.Context, sessionID string) error {
    // Revoke in database
    if err := s.repo.RevokeToken(ctx, sessionID); err != nil {
        return err
    }
    
    // Remove from cache
    return s.cacheManager.DeleteSession(ctx, sessionID)
}
```

#### Savings Service - Summary Caching

```go
func (s *Service) GetSummary(ctx context.Context, userID string) (*SavingsSummary, error) {
    // Try cache first
    cached, err := s.cacheManager.GetSavingsSummary(ctx, userID)
    if err == nil {
        return cached.(*SavingsSummary), nil
    }
    
    // Cache miss - calculate from database
    summary, err := s.repo.GetSummary(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    go s.cacheManager.SetSavingsSummary(context.Background(), userID, summary)
    
    return summary, nil
}

func (s *Service) CreateTransaction(ctx context.Context, userID string, tx *SavingsTransaction) error {
    if err := s.repo.CreateTransaction(ctx, tx); err != nil {
        return err
    }
    
    // Invalidate summary cache
    go s.cacheManager.InvalidateSavingsSummary(context.Background(), userID)
    
    return nil
}
```

### 4. Environment Configuration

**Update service deployments:**
```yaml
env:
  - name: REDIS_ADDR
    value: "redis:6379"
  - name: REDIS_PASSWORD
    valueFrom:
      secretKeyRef:
        name: redis-secret
        key: password
  - name: REDIS_DB
    value: "0"
```

**Service-specific database numbers:**
- Auth Service: DB 0
- Analytics Service: DB 1
- Education Service: DB 2
- Savings Service: DB 3
- Budget Service: DB 4
- Goal Service: DB 5
- User Service: DB 6

## Monitoring

### Redis Metrics

**Add to Prometheus scrape config:**
```yaml
- job_name: 'redis'
  static_configs:
    - targets: ['redis:6379']
```

**Key metrics to monitor:**
- `redis_connected_clients`: Number of connected clients
- `redis_used_memory_bytes`: Memory usage
- `redis_keyspace_hits_total`: Cache hits
- `redis_keyspace_misses_total`: Cache misses
- `redis_evicted_keys_total`: Evicted keys

**Cache hit ratio:**
```promql
rate(redis_keyspace_hits_total[5m]) / 
(rate(redis_keyspace_hits_total[5m]) + rate(redis_keyspace_misses_total[5m]))
```

### Grafana Dashboard

Create dashboard with panels for:
1. Cache hit ratio (target: > 80%)
2. Memory usage (target: < 80% of max)
3. Connected clients
4. Operations per second
5. Eviction rate
6. Latency (p95, p99)

## Testing

### Test Cache Functionality

```bash
# Connect to Redis
kubectl exec -it redis-0 -n insavein -- redis-cli -a <password>

# Check keys
KEYS *

# Get a cached value
GET financial_health:user:123:score

# Check TTL
TTL financial_health:user:123:score

# Monitor commands in real-time
MONITOR

# Get cache statistics
INFO stats
```

### Load Test with Caching

```bash
# Run load test
cd performance-tests
./run-tests.sh peak

# Monitor cache hit ratio during test
kubectl exec -it redis-0 -n insavein -- redis-cli -a <password> INFO stats | grep keyspace
```

## Cache Invalidation Patterns

### Pattern 1: Invalidate on Write

```go
func (s *Service) UpdateBudget(ctx context.Context, userID string, budget *Budget) error {
    // Update database
    if err := s.repo.UpdateBudget(ctx, budget); err != nil {
        return err
    }
    
    // Invalidate cache
    month := budget.Month.Format("2006-01")
    return s.cacheManager.InvalidateBudget(ctx, userID, month)
}
```

### Pattern 2: Invalidate Related Caches

```go
func (s *Service) RecordSpending(ctx context.Context, userID string, spending *Spending) error {
    // Update database
    if err := s.repo.RecordSpending(ctx, spending); err != nil {
        return err
    }
    
    // Invalidate multiple related caches
    go func() {
        ctx := context.Background()
        s.cacheManager.InvalidateBudget(ctx, userID, spending.Month)
        s.cacheManager.InvalidateFinancialHealth(ctx, userID)
    }()
    
    return nil
}
```

### Pattern 3: Invalidate All User Data

```go
func (s *Service) DeleteAccount(ctx context.Context, userID string) error {
    // Delete from database
    if err := s.repo.DeleteUser(ctx, userID); err != nil {
        return err
    }
    
    // Invalidate all user caches
    return s.cacheManager.InvalidateUserCache(ctx, userID)
}
```

## Performance Impact

### Expected Improvements

| Metric | Before Caching | After Caching | Improvement |
|--------|----------------|---------------|-------------|
| Financial Health API (p95) | 800ms | 50ms | 94% |
| Education Content API (p95) | 150ms | 20ms | 87% |
| Savings Summary API (p95) | 200ms | 30ms | 85% |
| Database Load | 100% | 40% | 60% reduction |
| Cache Hit Ratio | N/A | 85% | N/A |

### Cost Savings

- Reduced database CPU usage: 60%
- Reduced database I/O: 70%
- Reduced database connection usage: 50%
- Improved user experience: 3x faster response times

## Troubleshooting

### High Cache Miss Rate

**Symptoms:**
- Cache hit ratio < 50%
- High database load

**Solutions:**
- Increase TTL for stable data
- Pre-warm cache for common queries
- Review cache key generation

### Memory Pressure

**Symptoms:**
- Redis memory usage > 90%
- Frequent evictions

**Solutions:**
- Increase Redis memory limit
- Reduce TTL for less important data
- Review maxmemory-policy (allkeys-lru)

### Cache Stampede

**Symptoms:**
- Multiple requests for same data simultaneously
- Database spikes when cache expires

**Solutions:**
- Implement cache locking
- Use probabilistic early expiration
- Stagger TTL values

## Best Practices

1. **Always handle cache errors gracefully** - Fall back to database
2. **Use async cache operations** - Don't block on cache writes
3. **Monitor cache hit ratio** - Target > 80%
4. **Set appropriate TTLs** - Balance freshness vs. performance
5. **Invalidate on writes** - Keep cache consistent
6. **Use namespaced keys** - Avoid key collisions
7. **Test cache failures** - Ensure system works without cache

## Next Steps

1. Deploy Redis to Kubernetes
2. Update all services with caching code
3. Run load tests to verify improvements
4. Monitor cache metrics in Grafana
5. Tune TTL values based on usage patterns
