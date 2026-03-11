// API Request and Response Types

// Auth Types
export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
  date_of_birth: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserSummary
}

export interface TokenResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface UserSummary {
  id: string
  email: string
  first_name: string
  last_name: string
}

// User Types
export interface UserProfile {
  id: string
  email: string
  first_name: string
  last_name: string
  date_of_birth: string
  profile_image_url: string
  created_at: string
  updated_at: string
}

export interface UpdateProfileRequest {
  first_name?: string
  last_name?: string
  date_of_birth?: string
  profile_image_url?: string
}

export interface UserPreferences {
  currency: string
  notifications_enabled: boolean
  email_notifications: boolean
  push_notifications: boolean
  savings_reminders: boolean
  reminder_time: string
  theme: string
}

// Savings Types
export interface SavingsSummary {
  total_saved: number
  current_streak: number
  longest_streak: number
  last_saving_date: string
  monthly_average: number
  this_month_saved: number
}

export interface SavingsTransaction {
  id: string
  user_id: string
  amount: number
  currency: string
  description: string
  category: string
  created_at: string
}

export interface CreateSavingsRequest {
  amount: number
  currency?: string
  description?: string
  category?: string
}

export interface SavingsStreak {
  current_streak: number
  longest_streak: number
  last_save_date: string
}

export interface HistoryParams {
  limit?: number
  offset?: number
  start_date?: string
  end_date?: string
}

// Budget Types
export interface Budget {
  id: string
  user_id: string
  month: string
  total_budget: number
  categories: BudgetCategory[]
  total_spent: number
  remaining_budget: number
  created_at: string
  updated_at: string
}

export interface BudgetCategory {
  id: string
  budget_id: string
  name: string
  allocated_amount: number
  spent_amount: number
  remaining_amount: number
  color: string
}

export interface CreateBudgetRequest {
  month: string
  total_budget: number
  categories: {
    name: string
    allocated_amount: number
    color: string
  }[]
}

export interface UpdateBudgetRequest {
  total_budget?: number
  categories?: {
    id?: string
    name: string
    allocated_amount: number
    color: string
  }[]
}

export interface SpendingRequest {
  budget_id: string
  category_id: string
  amount: number
  description?: string
  merchant?: string
  date: string
}

export interface SpendingTransaction {
  id: string
  user_id: string
  budget_id: string
  category_id: string
  amount: number
  description: string
  merchant: string
  date: string
  created_at: string
}

export interface BudgetAlert {
  category_name: string
  percentage_used: number
  alert_type: 'warning' | 'critical'
  message: string
}

// Goal Types
export interface Goal {
  id: string
  user_id: string
  title: string
  description: string
  target_amount: number
  current_amount: number
  currency: string
  target_date: string
  status: 'active' | 'completed' | 'paused'
  progress_percent: number
  created_at: string
  updated_at: string
}

export interface CreateGoalRequest {
  title: string
  description?: string
  target_amount: number
  currency?: string
  target_date: string
}

export interface UpdateGoalRequest {
  title?: string
  description?: string
  target_amount?: number
  target_date?: string
  status?: 'active' | 'completed' | 'paused'
}

export interface Milestone {
  id: string
  goal_id: string
  title: string
  amount: number
  is_completed: boolean
  completed_at?: string
  order: number
}

// Education Types
export interface Lesson {
  id: string
  title: string
  description: string
  category: string
  duration_minutes: number
  difficulty: 'beginner' | 'intermediate' | 'advanced'
  tags: string[]
  is_completed: boolean
  order: number
}

export interface LessonDetail extends Lesson {
  content: string
  video_url?: string
  resources: Resource[]
  quiz?: QuizQuestion[]
}

export interface Resource {
  title: string
  url: string
  type: string
}

export interface QuizQuestion {
  question: string
  options: string[]
  correct_answer: number
}

export interface EducationProgress {
  total_lessons: number
  completed_lessons: number
  progress_percent: number
  current_streak: number
}

// Analytics Types
export interface SpendingAnalysis {
  period: TimePeriod
  total_spending: number
  category_breakdown: CategorySpending[]
  top_merchants: MerchantSpending[]
  daily_average: number
  comparison_to_previous: number
  trends: SpendingTrend[]
}

export interface TimePeriod {
  start_date: string
  end_date: string
}

export interface CategorySpending {
  category: string
  amount: number
  percentage: number
}

export interface MerchantSpending {
  merchant: string
  amount: number
  transaction_count: number
}

export interface SpendingTrend {
  date: string
  amount: number
}

export interface SavingsPattern {
  pattern_type: 'consistent' | 'irregular' | 'improving'
  average_amount: number
  frequency: string
  best_day_of_week: string
  insights: string[]
}

export interface Recommendation {
  id: string
  type: 'savings' | 'budget' | 'spending'
  priority: 'high' | 'medium' | 'low'
  title: string
  description: string
  action_items: string[]
  potential_savings?: number
}

export interface FinancialHealthScore {
  overall_score: number
  savings_score: number
  budget_score: number
  consistency_score: number
  insights: string[]
  improvement_areas: string[]
}

// Notification Types
export interface Notification {
  id: string
  user_id: string
  type: string
  title: string
  message: string
  is_read: boolean
  created_at: string
}

// Error Types
export interface ApiError {
  error: string
  message: string
  status: number
}
