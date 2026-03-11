import { PiggyBank, Wallet, Target, Flame } from 'lucide-react'
import {
  useSavingsSummary,
  useCurrentBudget,
  useActiveGoals,
} from '../lib/hooks/useDashboardData'

// Requirement 13.3: Quick stats cards
// Display savings this month, budget remaining, goals progress, current streak
export function QuickStatsCards() {
  const { data: savings } = useSavingsSummary()
  const { data: budget } = useCurrentBudget()
  const { data: goals } = useActiveGoals()

  // Requirement 4.5: Savings this month
  const thisMonthSaved = savings?.this_month_saved || 0

  // Requirement 5.2: Budget remaining
  const budgetRemaining = budget?.remaining_budget || 0
  const budgetPercentUsed = budget
    ? ((budget.total_spent / budget.total_budget) * 100).toFixed(0)
    : 0

  // Requirement 6.3: Goals progress
  const activeGoalsCount = goals?.length || 0
  const completedGoalsCount =
    goals?.filter((g) => g.status === 'completed').length || 0
  const avgProgress =
    goals && goals.length > 0
      ? (goals.reduce((sum, g) => sum + g.progress_percent, 0) / goals.length).toFixed(
          0
        )
      : 0

  // Requirement 9.3: Current streak
  const currentStreak = savings?.current_streak || 0

  const stats = [
    {
      name: 'Savings This Month',
      value: `$${thisMonthSaved.toFixed(2)}`,
      icon: PiggyBank,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
      description: 'Keep it up!',
    },
    {
      name: 'Budget Remaining',
      value: `$${budgetRemaining.toFixed(2)}`,
      icon: Wallet,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
      description: budget ? `${budgetPercentUsed}% used` : 'No budget set',
    },
    {
      name: 'Goals Progress',
      value: `${avgProgress}%`,
      icon: Target,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
      description: `${activeGoalsCount} active, ${completedGoalsCount} completed`,
    },
    {
      name: 'Current Streak',
      value: `${currentStreak} days`,
      icon: Flame,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
      description: currentStreak > 0 ? 'On fire! 🔥' : 'Start saving today!',
    },
  ]

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      {stats.map((stat) => {
        const Icon = stat.icon
        return (
          <div
            key={stat.name}
            className="border border-[var(--line)] rounded-xl p-6 hover:shadow-lg transition-shadow"
          >
            <div className="flex items-center justify-between mb-4">
              <div className={`p-3 rounded-lg ${stat.bgColor}`}>
                <Icon className={`w-6 h-6 ${stat.color}`} />
              </div>
            </div>
            <h3 className="text-sm font-semibold text-[var(--sea-ink-soft)] uppercase tracking-wide">
              {stat.name}
            </h3>
            <p className="mt-2 text-2xl font-bold text-[var(--sea-ink)]">
              {stat.value}
            </p>
            <p className="mt-1 text-sm text-[var(--sea-ink-soft)]">
              {stat.description}
            </p>
          </div>
        )
      })}
    </div>
  )
}
