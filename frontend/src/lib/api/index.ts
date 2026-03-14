import { ApiClient } from './client'
import { AuthService } from './auth'
import { UserService } from './user'
import { SavingsService } from './savings'
import { BudgetService } from './budget'
import { GoalService } from './goals'
import { EducationService } from './education'
import { AnalyticsService } from './analytics'
import { NotificationService } from './notifications'

// Service URLs from environment variables
const AUTH_SERVICE_URL = import.meta.env.VITE_AUTH_SERVICE_URL || 'http://localhost:8080'
const USER_SERVICE_URL = import.meta.env.VITE_USER_SERVICE_URL || 'http://localhost:8081'
const SAVINGS_SERVICE_URL = import.meta.env.VITE_SAVINGS_SERVICE_URL || 'http://localhost:8082'
const BUDGET_SERVICE_URL = import.meta.env.VITE_BUDGET_SERVICE_URL || 'http://localhost:8083'
const GOAL_SERVICE_URL = import.meta.env.VITE_GOAL_SERVICE_URL || 'http://localhost:8005'
const EDUCATION_SERVICE_URL = import.meta.env.VITE_EDUCATION_SERVICE_URL || 'http://localhost:8085'
const NOTIFICATION_SERVICE_URL = import.meta.env.VITE_NOTIFICATION_SERVICE_URL || 'http://localhost:8086'
const ANALYTICS_SERVICE_URL = import.meta.env.VITE_ANALYTICS_SERVICE_URL || 'http://localhost:8008'

export function createApiServices(
  getToken: () => string | null,
  onUnauthorized: () => void
) {
  // Create API clients for each service
  const authClient = new ApiClient(AUTH_SERVICE_URL, getToken, onUnauthorized)
  const userClient = new ApiClient(USER_SERVICE_URL, getToken, onUnauthorized)
  const savingsClient = new ApiClient(SAVINGS_SERVICE_URL, getToken, onUnauthorized)
  const budgetClient = new ApiClient(BUDGET_SERVICE_URL, getToken, onUnauthorized)
  const goalClient = new ApiClient(GOAL_SERVICE_URL, getToken, onUnauthorized)
  const educationClient = new ApiClient(EDUCATION_SERVICE_URL, getToken, onUnauthorized)
  const notificationClient = new ApiClient(NOTIFICATION_SERVICE_URL, getToken, onUnauthorized)
  const analyticsClient = new ApiClient(ANALYTICS_SERVICE_URL, getToken, onUnauthorized)

  return {
    auth: new AuthService(authClient),
    user: new UserService(userClient),
    savings: new SavingsService(savingsClient),
    budget: new BudgetService(budgetClient),
    goals: new GoalService(goalClient),
    education: new EducationService(educationClient),
    notifications: new NotificationService(notificationClient),
    analytics: new AnalyticsService(analyticsClient),
  }
}

export type ApiServices = ReturnType<typeof createApiServices>
