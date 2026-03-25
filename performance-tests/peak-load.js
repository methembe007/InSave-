import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const responseTime = new Trend('response_time');

// Test configuration for peak load
export const options = {
  stages: [
    { duration: '3m', target: 1000 },   // Ramp up to 1000 users
    { duration: '5m', target: 5000 },   // Ramp up to 5000 users
    { duration: '5m', target: 10000 },  // Ramp up to 10000 users (peak)
    { duration: '10m', target: 10000 }, // Stay at peak for 10 minutes
    { duration: '3m', target: 0 },      // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // Performance targets
    http_req_failed: ['rate<0.001'],                 // Error rate < 0.1%
    errors: ['rate<0.001'],
    'http_reqs': ['rate>1666'],                      // ~100k req/min = 1666 req/sec
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
  // Register more test users for peak load
  const testUsers = [];
  for (let i = 0; i < 50; i++) {
    const email = `peaktest${Date.now()}_${i}@example.com`;
    const password = 'TestPassword123!';
    
    const registerRes = http.post(`${BASE_URL}/api/auth/register`, JSON.stringify({
      email: email,
      password: password,
      first_name: 'Peak',
      last_name: `Test${i}`,
      date_of_birth: '1990-01-01',
    }), {
      headers: { 'Content-Type': 'application/json' },
    });
    
    if (registerRes.status === 200 || registerRes.status === 201) {
      const body = JSON.parse(registerRes.body);
      testUsers.push({
        email: email,
        password: password,
        token: body.access_token,
      });
    }
  }
  
  return { users: testUsers };
}

export default function(data) {
  const user = data.users[Math.floor(Math.random() * data.users.length)];
  
  if (!user) {
    return;
  }
  
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${user.token}`,
    },
  };
  
  // More aggressive request patterns for peak load
  const scenario = Math.random();
  
  if (scenario < 0.4) {
    // 40% - Read operations (dashboard, summaries)
    const requests = [
      http.get(`${BASE_URL}/api/user/profile`, params),
      http.get(`${BASE_URL}/api/savings/summary`, params),
      http.get(`${BASE_URL}/api/budget/current`, params),
    ];
    
    requests.forEach(res => {
      responseTime.add(res.timings.duration);
      check(res, {
        'read status ok': (r) => r.status === 200,
        'read response time < 500ms': (r) => r.timings.duration < 500,
      }) || errorRate.add(1);
    });
    
  } else if (scenario < 0.6) {
    // 20% - Write operations (transactions)
    const savingsRes = http.post(`${BASE_URL}/api/savings/transactions`, JSON.stringify({
      amount: Math.floor(Math.random() * 100) + 1,
      currency: 'USD',
      description: 'Peak load test',
      category: 'general',
    }), params);
    
    responseTime.add(savingsRes.timings.duration);
    check(savingsRes, {
      'write status ok': (r) => r.status === 200 || r.status === 201,
      'write response time < 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);
    
  } else if (scenario < 0.8) {
    // 20% - Analytics (more expensive queries)
    const analyticsRes = http.get(`${BASE_URL}/api/analytics/health`, params);
    
    responseTime.add(analyticsRes.timings.duration);
    check(analyticsRes, {
      'analytics status ok': (r) => r.status === 200,
      'analytics response time < 1000ms': (r) => r.timings.duration < 1000,
    }) || errorRate.add(1);
    
  } else {
    // 20% - Mixed operations
    const historyRes = http.get(`${BASE_URL}/api/savings/history?limit=20`, params);
    
    responseTime.add(historyRes.timings.duration);
    check(historyRes, {
      'history status ok': (r) => r.status === 200,
      'history response time < 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);
  }
  
  // Shorter think time for peak load
  sleep(Math.random() * 2 + 0.5); // 0.5-2.5 seconds
}

export function teardown(data) {
  console.log('Peak load test completed');
}
