import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { Lightbulb, CheckCircle2, AlertCircle, Info } from 'lucide-react'

// Requirement 13.4: Fetch and display AI-assisted recommendations
export function RecommendationsList() {
  const { api } = useAuth()
  
  const { data: recommendations, isLoading, error } = useQuery({
    queryKey: ['recommendations'],
    queryFn: () => api.analytics.getRecommendations(),
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
        <p className="text-red-800">Unable to load recommendations.</p>
      </div>
    )
  }

  if (!recommendations || recommendations.length === 0) {
    return (
      <div className="border border-[var(--line)] rounded-xl p-6">
        <h2 className="text-xl font-bold text-[var(--sea-ink)] mb-4">
          Recommendations
        </h2>
        <p className="text-[var(--sea-ink-soft)]">
          No recommendations at this time. Keep tracking your finances!
        </p>
      </div>
    )
  }

  const getPriorityIcon = (priority: string) => {
    switch (priority) {
      case 'high':
        return <AlertCircle className="w-5 h-5 text-red-600" />
      case 'medium':
        return <Info className="w-5 h-5 text-yellow-600" />
      case 'low':
        return <Lightbulb className="w-5 h-5 text-blue-600" />
      default:
        return <Info className="w-5 h-5 text-gray-600" />
    }
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high':
        return 'bg-red-50 border-red-200'
      case 'medium':
        return 'bg-yellow-50 border-yellow-200'
      case 'low':
        return 'bg-blue-50 border-blue-200'
      default:
        return 'bg-gray-50 border-gray-200'
    }
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value)
  }

  // Sort by priority: high > medium > low
  const sortedRecommendations = [...recommendations].sort((a, b) => {
    const priorityOrder = { high: 0, medium: 1, low: 2 }
    return priorityOrder[a.priority as keyof typeof priorityOrder] - 
           priorityOrder[b.priority as keyof typeof priorityOrder]
  })

  return (
    <div className="border border-[var(--line)] rounded-xl p-6">
      <div className="flex items-center gap-3 mb-6">
        <Lightbulb className="w-6 h-6 text-[var(--link)]" />
        <h2 className="text-xl font-bold text-[var(--sea-ink)]">
          AI Recommendations
        </h2>
      </div>

      <div className="space-y-4">
        {sortedRecommendations.map((rec) => (
          <div
            key={rec.id}
            className={`rounded-lg p-4 border ${getPriorityColor(rec.priority)}`}
          >
            {/* Header with Priority */}
            <div className="flex items-start justify-between mb-3">
              <div className="flex items-center gap-2">
                {getPriorityIcon(rec.priority)}
                <h3 className="font-bold text-[var(--sea-ink)]">{rec.title}</h3>
              </div>
              <span className="text-xs font-semibold uppercase tracking-wide px-2 py-1 rounded">
                {rec.priority}
              </span>
            </div>

            {/* Description */}
            <p className="text-sm text-[var(--sea-ink)] mb-3">
              {rec.description}
            </p>

            {/* Potential Savings */}
            {rec.potential_savings && rec.potential_savings > 0 && (
              <div className="mb-3 p-2 bg-white rounded border border-current opacity-50">
                <p className="text-xs font-medium">Potential Savings</p>
                <p className="text-lg font-bold">
                  {formatCurrency(rec.potential_savings)}
                </p>
              </div>
            )}

            {/* Action Items Checklist */}
            {rec.action_items.length > 0 && (
              <div className="mt-3 pt-3 border-t border-current opacity-30">
                <p className="text-xs font-semibold mb-2">Action Items:</p>
                <ul className="space-y-1">
                  {rec.action_items.map((item: string, idx: number) => (
                    <li key={idx} className="flex items-start gap-2 text-sm">
                      <CheckCircle2 className="w-4 h-4 mt-0.5 flex-shrink-0" />
                      <span>{item}</span>
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
