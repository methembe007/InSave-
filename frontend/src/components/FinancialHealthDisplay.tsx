import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { Activity, TrendingUp, Wallet, Target, AlertCircle } from 'lucide-react'

// Requirement 14.1: Display overall score with visual gauge/meter
// Requirement 14.2: Show component scores (savings, budget, consistency)
// Requirement 14.3: Display insights and improvement areas
export function FinancialHealthDisplay() {
  const { api } = useAuth()
  
  const { data: healthScore, isLoading, error } = useQuery({
    queryKey: ['financial-health'],
    queryFn: () => api.analytics.getFinancialHealth(),
    staleTime: 1000 * 60 * 60, // 1 hour cache as per requirement 13.6
  })

  if (isLoading) {
    return (
      <div className="border border-[var(--line)] rounded-xl p-8 animate-pulse">
        <div className="h-48 bg-gray-200 rounded"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="border border-red-200 bg-red-50 rounded-xl p-6">
        <div className="flex items-center gap-3 text-red-800">
          <AlertCircle className="w-5 h-5" />
          <p>
            {error instanceof Error && error.message.includes('insufficient data')
              ? 'Not enough data yet. Keep saving for 30 days to see your financial health score!'
              : 'Unable to load financial health score. Please try again later.'}
          </p>
        </div>
      </div>
    )
  }

  if (!healthScore) return null

  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-600'
    if (score >= 60) return 'text-yellow-600'
    return 'text-red-600'
  }

  const getScoreBgColor = (score: number) => {
    if (score >= 80) return 'bg-green-50'
    if (score >= 60) return 'bg-yellow-50'
    return 'bg-red-50'
  }

  const getScoreLabel = (score: number) => {
    if (score >= 80) return 'Excellent'
    if (score >= 60) return 'Good'
    if (score >= 40) return 'Fair'
    return 'Needs Improvement'
  }

  return (
    <div className="border border-[var(--line)] rounded-xl p-8 bg-gradient-to-br from-blue-50 to-purple-50">
      <div className="flex items-center gap-3 mb-6">
        <Activity className="w-6 h-6 text-[var(--link)]" />
        <h2 className="text-2xl font-bold text-[var(--sea-ink)]">
          Financial Health Score
        </h2>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Overall Score - Prominent Display */}
        <div className="lg:col-span-1 flex flex-col items-center justify-center">
          <div className="relative w-48 h-48">
            {/* Circular Progress */}
            <svg className="w-48 h-48 transform -rotate-90">
              <circle
                cx="96"
                cy="96"
                r="88"
                stroke="currentColor"
                strokeWidth="12"
                fill="none"
                className="text-gray-200"
              />
              <circle
                cx="96"
                cy="96"
                r="88"
                stroke="currentColor"
                strokeWidth="12"
                fill="none"
                strokeDasharray={`${(healthScore.overall_score / 100) * 553} 553`}
                className={getScoreColor(healthScore.overall_score)}
                strokeLinecap="round"
              />
            </svg>
            <div className="absolute inset-0 flex flex-col items-center justify-center">
              <span className={`text-5xl font-bold ${getScoreColor(healthScore.overall_score)}`}>
                {healthScore.overall_score}
              </span>
              <span className="text-sm text-[var(--sea-ink-soft)] mt-1">
                {getScoreLabel(healthScore.overall_score)}
              </span>
            </div>
          </div>
        </div>

        {/* Component Scores */}
        <div className="lg:col-span-2 space-y-4">
          <h3 className="text-lg font-semibold text-[var(--sea-ink)] mb-4">
            Score Breakdown
          </h3>

          {/* Savings Score (40% weight) */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <TrendingUp className="w-5 h-5 text-green-600" />
                <span className="font-medium text-[var(--sea-ink)]">
                  Savings Score
                </span>
                <span className="text-xs text-[var(--sea-ink-soft)]">(40% weight)</span>
              </div>
              <span className={`font-bold ${getScoreColor(healthScore.savings_score)}`}>
                {healthScore.savings_score}
              </span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-3">
              <div
                className={`h-3 rounded-full transition-all ${getScoreBgColor(healthScore.savings_score)}`}
                style={{ width: `${healthScore.savings_score}%` }}
              />
            </div>
          </div>

          {/* Budget Score (30% weight) */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Wallet className="w-5 h-5 text-blue-600" />
                <span className="font-medium text-[var(--sea-ink)]">
                  Budget Score
                </span>
                <span className="text-xs text-[var(--sea-ink-soft)]">(30% weight)</span>
              </div>
              <span className={`font-bold ${getScoreColor(healthScore.budget_score)}`}>
                {healthScore.budget_score}
              </span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-3">
              <div
                className={`h-3 rounded-full transition-all ${getScoreBgColor(healthScore.budget_score)}`}
                style={{ width: `${healthScore.budget_score}%` }}
              />
            </div>
          </div>

          {/* Consistency Score (30% weight) */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Target className="w-5 h-5 text-purple-600" />
                <span className="font-medium text-[var(--sea-ink)]">
                  Consistency Score
                </span>
                <span className="text-xs text-[var(--sea-ink-soft)]">(30% weight)</span>
              </div>
              <span className={`font-bold ${getScoreColor(healthScore.consistency_score)}`}>
                {healthScore.consistency_score}
              </span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-3">
              <div
                className={`h-3 rounded-full transition-all ${getScoreBgColor(healthScore.consistency_score)}`}
                style={{ width: `${healthScore.consistency_score}%` }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Insights and Improvement Areas */}
      <div className="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Insights */}
        {healthScore.insights.length > 0 && (
          <div className="bg-white rounded-lg p-6 border border-green-200">
            <h3 className="text-lg font-semibold text-[var(--sea-ink)] mb-3 flex items-center gap-2">
              <TrendingUp className="w-5 h-5 text-green-600" />
              Positive Insights
            </h3>
            <ul className="space-y-2">
              {healthScore.insights.map((insight: string, index: number) => (
                <li key={index} className="flex items-start gap-2 text-sm text-[var(--sea-ink)]">
                  <span className="text-green-600 mt-0.5">✓</span>
                  <span>{insight}</span>
                </li>
              ))}
            </ul>
          </div>
        )}

        {/* Improvement Areas */}
        {healthScore.improvement_areas.length > 0 && (
          <div className="bg-white rounded-lg p-6 border border-yellow-200">
            <h3 className="text-lg font-semibold text-[var(--sea-ink)] mb-3 flex items-center gap-2">
              <AlertCircle className="w-5 h-5 text-yellow-600" />
              Areas to Improve
            </h3>
            <ul className="space-y-2">
              {healthScore.improvement_areas.map((area: string, index: number) => (
                <li key={index} className="flex items-start gap-2 text-sm text-[var(--sea-ink)]">
                  <span className="text-yellow-600 mt-0.5">→</span>
                  <span>{area}</span>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  )
}
