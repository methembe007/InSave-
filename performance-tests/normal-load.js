import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration for normal load
export const options = {
  stages: [
    { duration: '2m', target: 100 },   // Ramp up to 100 users
    { duration: '5m', target: 1000 },  // Ramp up to 1000 users
    { duration: '10m', target: 1000 }, // Stay at 1000 users
    { duration: '2m', target: 0 },     // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95% < 500ms, 99% < 1000ms
    http_req_failed: ['rate<0.001'],                 // Error rate < 0.1%
    errors: ['rate<0.001'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Test data
const users = [];
let authToken = '';

export function setup() {
  // Register test users
  const testUsers = [];
  for (let i = 0; i < 10; i++) {
    const email = `loadtest${Date.now()}_${i}@example.com`;
    const password = 'TestPassword123!';
    
    const registerRes = http.post(`${BASE_URL}/api/auth/register`, JSON.stringify({
      email: email,
      password: password,
      first_name: 'Load',
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
  // Select a random user
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
  
  // Simulate realistic user behavior with weighted scenarios
  const scenario = Math.random();
  
  if (scenario < 0.3) {
    // 30% - View dashboard
    const dashboardRes = http.get(`${BASE_URL}/api/user/profile`, params);
    check(dashboardRes, {
      'dashboard status 200': (r) => r.status === 200,
      'dashboard response time < 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);
    
    const savingsRes = http.get(`${BASE_URL}/api/savings/summary`, params);
    check(savingsRes, { 'savings summary ok': (r) => r.status === 200 }) || errorRate.add(1);
    
    const budgetRes = http.get(`${BASE_URL}/api/budget/current`, params);
    check(budgetRes, { 'budget ok': (r) => r.status === 200 }) || errorRate.add(1);
    
  } else if (scenario < 0.5) {
    // 20% - Create savings transaction
    const savingsRes = http.post(`${BASE_URL}/api/savings/transactions`, JSON.stringify({
      amount: Math.floor(Math.random() * 100) + 1,
      currency: 'USD',
      description: 'Load test savings',
      category: 'general',
    }), params);
    
    check(savingsRes, {
      'savings created': (r) => r.status === 200 || r.status === 201,
      'savings response time < 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);
    
  } else if (scenario < 0.7) {
    // 20% - Record spending
    const spendingRes = http.post(`${BASE_URL}/api/budget/spending`, JSON.stringify({
      amount: Math.floor(Math.random() * 50) + 1,
      category_id: 'test-category',
      description: 'Load test spending',
      merchant: 'Test Store',
      date: new Date().toISOString().split('T')[0],
    }), params);
    
    check(spendingRes, {
      'spending recorded': (r) => r.status === 200 || r.status === 201,
      'spending response time < 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);
    
  } else if (scenario < 0.85) {
    // 15% - View analytics
    const analyticsRes = http.get(`${BASE_URL}/api/analytics/health`, params);
    check(analyticsRes, {
      'analytics ok': (r) => r.status === 200,
      'analytics response time < 1000ms': (r) => r.timings.duration < 1000,
    }) || errorRate.add(1);
    
  } else {
    // 15% - View goals
    const goalsRes = http.get(`${BASE_URL}/api/goals`, params);
    check(goalsRes, {
      'goals ok': (r) => r.status === 200,
      'goals response time < 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);
  }
  
  // Think time between requests (simulate real user behavior)
  sleep(Math.random() * 3 + 1); // 1-4 seconds
}

export function teardown(data) {
  // Cleanup test users if needed
  console.log('Load test completed');
}
