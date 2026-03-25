import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const responseTime = new Trend('response_time');
const failedRequests = new Counter('failed_requests');

// Stress test configuration - gradually increase load to find breaking point
export const options = {
  stages: [
    { duration: '2m', target: 5000 },   // Ramp to 5k users
    { duration: '3m', target: 10000 },  // Ramp to 10k users
    { duration: '3m', target: 20000 },  // Ramp to 20k users
    { duration: '3m', target: 30000 },  // Ramp to 30k users
    { duration: '3m', target: 40000 },  // Ramp to 40k users
    { duration: '3m', target: 50000 },  // Ramp to 50k users (stress point)
    { duration: '5m', target: 50000 },  // Hold at 50k
    { duration: '3m', target: 0 },      // Ramp down
  ],
  thresholds: {
    // More relaxed thresholds for stress test - we expect some degradation
    http_req_duration: ['p(95)<2000', 'p(99)<5000'],
    http_req_failed: ['rate<0.05'], // Allow up to 5% error rate
    errors: ['rate<0.05'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
  // Register many test users for stress test
  const testUsers = [];
  const batchSize = 100;
  
  for (let i = 0; i < batchSize; i++) {
    const email = `stresstest${Date.now()}_${i}@example.com`;
    const password = 'TestPassword123!';
    
    const registerRes = http.post(`${BASE_URL}/api/auth/register`, JSON.stringify({
      email: email,
      password: password,
      first_name: 'Stress',
      last_name: `Test${i}`,
      date_of_birth: '1990-01-01',
    }), {
      headers: { 'Content-Type': 'application/json' },
      timeout: '10s',
    });
    
    if (registerRes.status === 200 || registerRes.status === 201) {
      const body = JSON.parse(registerRes.body);
      testUsers.push({
        email: email,
        password: password,
        token: body.access_token,
      });
    }
    
    // Small delay to avoid overwhelming during setup
    if (i % 10 === 0) {
      sleep(0.1);
    }
  }
  
  console.log(`Setup complete: ${testUsers.length} users registered`);
  return { users: testUsers };
}

export default function(data) {
  const user = data.users[Math.floor(Math.random() * data.users.length)];
  
  if (!user) {
    failedRequests.add(1);
    return;
  }
  
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${user.token}`,
    },
    timeout: '10s', // Longer timeout for stress conditions
  };
  
  // Aggressive request patterns
  const scenario = Math.random();
  
  try {
    if (scenario < 0.5) {
      // 50% - Read operations
      const res = http.get(`${BASE_URL}/api/savings/summary`, params);
      responseTime.add(res.timings.duration);
      
      const success = check(res, {
        'status ok': (r) => r.status === 200,
        'response received': (r) => r.body && r.body.length > 0,
      });
      
      if (!success) {
        errorRate.add(1);
        failedRequests.add(1);
      }
      
    } else if (scenario < 0.75) {
      // 25% - Write operations
      const res = http.post(`${BASE_URL}/api/savings/transactions`, JSON.stringify({
        amount: Math.floor(Math.random() * 100) + 1,
        currency: 'USD',
        description: 'Stress test transaction',
        category: 'general',
      }), params);
      
      responseTime.add(res.timings.duration);
      
      const success = check(res, {
        'write ok': (r) => r.status === 200 || r.status === 201,
      });
      
      if (!success) {
        errorRate.add(1);
        failedRequests.add(1);
      }
      
    } else {
      // 25% - Complex queries (analytics)
      const res = http.get(`${BASE_URL}/api/analytics/health`, params);
      responseTime.add(res.timings.duration);
      
      const success = check(res, {
        'analytics ok': (r) => r.status === 200 || r.status === 503, // 503 acceptable under stress
      });
      
      if (!success) {
        errorRate.add(1);
        failedRequests.add(1);
      }
    }
  } catch (e) {
    errorRate.add(1);
    failedRequests.add(1);
    console.error(`Request failed: ${e.message}`);
  }
  
  // Minimal think time for stress test
  sleep(Math.random() * 1 + 0.1); // 0.1-1.1 seconds
}

export function teardown(data) {
  console.log('Stress test completed');
  console.log(`Total users: ${data.users.length}`);
}

export function handleSummary(data) {
  return {
    'stress-test-summary.json': JSON.stringify(data),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  };
}

function textSummary(data, options) {
  const indent = options.indent || '';
  const enableColors = options.enableColors || false;
  
  let summary = '\n' + indent + '=== Stress Test Summary ===\n\n';
  
  // Request metrics
  summary += indent + `Total Requests: ${data.metrics.http_reqs.values.count}\n`;
  summary += indent + `Failed Requests: ${data.metrics.http_req_failed.values.passes}\n`;
  summary += indent + `Request Rate: ${data.metrics.http_reqs.values.rate.toFixed(2)} req/s\n\n`;
  
  // Response time metrics
  summary += indent + 'Response Times:\n';
  summary += indent + `  p50: ${data.metrics.http_req_duration.values['p(50)'].toFixed(2)}ms\n`;
  summary += indent + `  p95: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms\n`;
  summary += indent + `  p99: ${data.metrics.http_req_duration.values['p(99)'].toFixed(2)}ms\n`;
  summary += indent + `  max: ${data.metrics.http_req_duration.values.max.toFixed(2)}ms\n\n`;
  
  // Error rate
  const errorRate = (data.metrics.errors.values.rate * 100).toFixed(2);
  summary += indent + `Error Rate: ${errorRate}%\n`;
  
  return summary;
}
