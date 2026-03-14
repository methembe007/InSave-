import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { TrendingUp, Calendar, DollarSign, Lightbulb } from 'lucide-react'

// Requirement 13.3: Fetch and display savings patterns
export function SavingsPatternsDisplay() {
  const { api } = useAuth()
  
  const { data: patterns, isLoading, error } = useQuery({
    queryKey: ['savings-patterns'],
    queryFn: () => api.analytics.getSavingsPatterns(),
  })

  if (isLoading) {
    return (
      <div className="border border-[var(--line)] rounded-xl p-6 animate-pulse">
        <div className="h-64 bg-gray-200 rounded"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="border border-red-200 bg-red-50 rounded-xl p-6">
        <p className="text-red-800">Unable to load savings patterns.</p>
      </div>
    )
  }

  if (!patterns || patterns.length === 0) {
    return (
      <div className="border border-[var(--line)] rounded-xl p-6">
        <h2 className="text-xl font-bold text-[var(--sea-ink)] mb-4">Savings Patterns</h2>
        <p className="text-[var(--sea-ink-soft)]">
          Keep saving to see your patterns!
        </p>
      </div>
    )
  }

  const getPatternColor = (type: string) => {
    switch (type) {
      case 'consistent':
        return 'bg-green-50 border-green-200 text-green-800'
      case 'improving':
        return 'bg-blue-50 border-blue-200 text-blue-800'
      case 'irregular':
        return 'bg-yellow-50 border-yellow-200 text-yellow-800'
      default:
        return 'bg-gray-50 border-gray-200 text-gray-800'
    }
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value)
  }

  return (
    <div className="border border-[var(--line)] rounded-xl p-6">
      <div className="flex items-center gap-3 mb-6">
        <TrendingUp className="w-6 h-6 text-[var(--link)]" />
        <h2 className="text-xl font-bold text-[var(--sea-ink)]">Savings Patterns</h2>
      </div>

      <div className="space-y-4">
        {patterns.map((pattern: any, index: number) => (
          <div
            key={index}
            className={`rounded-lg p-4 border ${getPatternColor(pattern.pattern_type)}`}
          >
            {/* Pattern Type Badge */}
            <div className="flex items-center justify-between mb-3">
              <span className="text-sm font-semibold uppercase tracking-wide">
                {pattern.pattern_type}
              </span>
              <Calendar className="w-4 h-4" />
            </div>

            {/* Pattern Stats */}
            <div className="grid grid-cols-2 gap-4 mb-3">
              <div>
                <p className="text-xs opacity-75 mb-1">Average Amount</p>
                <p className="font-bold flex items-center gap-1">
                  <DollarSign className="w-4 h-4" />
                  {formatCurrency(pattern.average_amount)}
                </p>
              </div>
              <div>
                <p className="text-xs opacity-75 mb-1">Frequency</p>
                <p className="font-bold">{pattern.frequency}</p>
              </div>
            </div>

            {/* Best Day */}
            <div className="mb-3">
              <p className="text-xs opacity-75 mb-1">Best Day to Save</p>
              <p className="font-bold">{pattern.best_day_of_week}</p>
            </div>

            {/* Insights */}
            {pattern.insights.length > 0 && (
              <div className="mt-4 pt-4 border-t border-current opacity-50">
                <div className="flex items-start gap-2">
                  <Lightbulb className="w-4 h-4 mt-0.5 flex-shrink-0" />
                  <div className="space-y-1">
                    {pattern.insights.map((insight: string, idx: number) => (
                      <p key={idx} className="text-sm">
                        {insight}
                      </p>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
