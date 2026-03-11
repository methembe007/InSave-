import { TrendingUp, TrendingDown, Minus } from 'lucide-react'
import {
  useSavingsSummary,
  useCurrentBudget,
  useActiveGoals,
  useFinancialHealth,
} from '../lib/hooks/useDashboardData'

// Requirement 13.2: Dashboard summary component
// Fetches data from savings.getSummary, budget.getCurrentBudget, goals.getActiveGoals, analytics.getFinancialHealth
// Displays total saved, current streak, budget status, active goals count
// Shows financial health score with visual indicator
export function DashboardSummary() {
  const { data: savings, isLoading: savingsLoading } = useSavingsSummary()
  const { data: budget, isLoading: budgetLoading } = useCurrentBudget()
  const { data: goals, isLoading: goalsLoading } = useActiveGoals()
  const { data: health, isLoading: healthLoading } = useFinancialHealth()

  const isLoading = savingsLoading || budgetLoading || goalsLoading || healthLoading

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {[...Array(4)].map((_, i) => (
          <div
            key={i}
            className="border border-[var(--line)] rounded-xl p-6 animate-pulse"
          >
            <div className="h-4 bg-[var(--link-bg-hover)] rounded w-24 mb-4" />
            <div className="h-8 bg-[var(--link-bg-hover)] rounded w-32 mb-2" />
            <div className="h-3 bg-[var(--link-bg-hover)] rounded w-20" />
          </div>
        ))}
      </div>
    )
  }

  // Requirement 4.5: Display total saved and current streak
  const totalSaved = savings?.total_saved || 0
  const currentStreak = savings?.current_streak || 0
  const thisMonthSaved = savings?.this_month_saved || 0

  // Requirement 6.3: Display budget status
  const budgetRemaining = budget?.remaining_budget || 0
  const totalBudget = budget?.total_budget || 0

  // Requirement 9.3: Display active goals count
  const activeGoalsCount = goals?.length || 0

  // Requirement 14.1: Display financial health score
  const healthScore = health?.overall_score || 0

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      {/* Total Saved Card */}
      <div className="border border-[var(--line)] rounded-xl p-6 hover:shadow-lg transition-shadow">
        <h3 className="text-sm font-semibold text-[var(--sea-ink-soft)] uppercase tracking-wide">
          Total Saved
        </h3>
        <p className="mt-2 text-3xl font-bold text-[var(--sea-ink)]">
          ${totalSaved.toFixed(2)}
        </p>
        <p className="mt-1 text-sm text-[var(--sea-ink-soft)]">
          This month: ${thisMonthSaved.toFixed(2)}
        </p>
      </div>

      {/* Current Streak Card */}
      <div className="border border-[var(--line)] rounded-xl p-6 hover:shadow-lg transition-shadow">
        <h3 className="text-sm font-semibold text-[var(--sea-ink-soft)] uppercase tracking-wide">
          Current Streak
        </h3>
        <div className="flex items-baseline gap-2 mt-2">
          <p className="text-3xl font-bold text-[var(--sea-ink)]">
            {currentStreak}
          </p>
          <span className="text-xl text-[var(--sea-ink-soft)]">days</span>
          {currentStreak > 0 && <span className="text-2xl">🔥</span>}
        </div>
        <p className="mt-1 text-sm text-[var(--sea-ink-soft)]">
          Longest: {savings?.longest_streak || 0} days
        </p>
      </div>

      {/* Budget Status Card */}
      <div className="border border-[var(--line)] rounded-xl p-6 hover:shadow-lg transition-shadow">
        <h3 className="text-sm font-semibold text-[var(--sea-ink-soft)] uppercase tracking-wide">
          Budget Status
        </h3>
        <p className="mt-2 text-3xl font-bold text-[var(--sea-ink)]">
          ${budgetRemaining.toFixed(2)}
        </p>
        <p className="mt-1 text-sm text-[var(--sea-ink-soft)]">
          {budget ? `of $${totalBudget.toFixed(2)} remaining` : 'No budget set'}
        </p>
      </div>

      {/* Financial Health Card */}
      <div className="border border-[var(--line)] rounded-xl p-6 hover:shadow-lg transition-shadow">
        <h3 className="text-sm font-semibold text-[var(--sea-ink-soft)] uppercase tracking-wide">
          Financial Health
        </h3>
        <div className="flex items-center gap-3 mt-2">
          <p className="text-3xl font-bold text-[var(--sea-ink)]">
            {healthScore}
          </p>
          <HealthIndicator score={healthScore} />
        </div>
        <p className="mt-1 text-sm text-[var(--sea-ink-soft)]">
          {activeGoalsCount} active {activeGoalsCount === 1 ? 'goal' : 'goals'}
        </p>
      </div>
    </div>
  )
}

// Visual indicator for financial health score
function HealthIndicator({ score }: { score: number }) {
  if (score >= 80) {
    return (
      <div className="flex items-center gap-1 text-green-600">
        <TrendingUp className="w-5 h-5" />
        <span className="text-sm font-semibold">Excellent</span>
      </div>
    )
  }
  if (score >= 60) {
    return (
      <div className="flex items-center gap-1 text-blue-600">
        <Minus className="w-5 h-5" />
        <span className="text-sm font-semibold">Good</span>
      </div>
    )
  }
  if (score >= 40) {
    return (
      <div className="flex items-center gap-1 text-yellow-600">
        <TrendingDown className="w-5 h-5" />
        <span className="text-sm font-semibold">Fair</span>
      </div>
    )
  }
  return (
    <div className="flex items-center gap-1 text-red-600">
      <TrendingDown className="w-5 h-5" />
      <span className="text-sm font-semibold">Needs Work</span>
    </div>
  )
}
