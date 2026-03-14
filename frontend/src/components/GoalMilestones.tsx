import type { Milestone } from '../lib/types/api'

interface GoalMilestonesProps {
  goalId: string
  milestones: Milestone[]
  isLoading: boolean
}

export function GoalMilestones({
  milestones,
  isLoading,
}: GoalMilestonesProps) {
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount)
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

  if (isLoading) {
    return (
      <div className="bg-gray-50 rounded-lg p-4">
        <p className="text-sm text-gray-500 text-center">Loading milestones...</p>
      </div>
    )
  }

  if (milestones.length === 0) {
    return (
      <div className="bg-gray-50 rounded-lg p-4">
        <p className="text-sm text-gray-500 text-center">
          No milestones set for this goal
        </p>
      </div>
    )
  }

  // Sort milestones by order
  const sortedMilestones = [...milestones].sort((a, b) => a.order - b.order)

  return (
    <div className="bg-gray-50 rounded-lg p-4">
      <h4 className="text-sm font-medium text-gray-900 mb-3">Milestones</h4>
      <div className="space-y-3">
        {sortedMilestones.map((milestone) => (
          <div
            key={milestone.id}
            className={`flex items-start gap-3 p-3 rounded-lg ${
              milestone.is_completed
                ? 'bg-green-50 border border-green-200'
                : 'bg-white border border-gray-200'
            }`}
          >
            {/* Checkmark Icon */}
            <div className="flex-shrink-0 mt-0.5">
              {milestone.is_completed ? (
                <svg
                  className="w-5 h-5 text-green-600"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                    clipRule="evenodd"
                  />
                </svg>
              ) : (
                <svg
                  className="w-5 h-5 text-gray-400"
                  fill="none"
                  viewBox="0 0 20 20"
                  stroke="currentColor"
                >
                  <circle
                    cx="10"
                    cy="10"
                    r="8"
                    strokeWidth="2"
                    className="stroke-current"
                  />
                </svg>
              )}
            </div>

            {/* Milestone Details */}
            <div className="flex-1 min-w-0">
              <div className="flex items-start justify-between gap-2">
                <h5
                  className={`text-sm font-medium ${
                    milestone.is_completed
                      ? 'text-green-900'
                      : 'text-gray-900'
                  }`}
                >
                  {milestone.title}
                </h5>
                <span
                  className={`text-sm font-medium whitespace-nowrap ${
                    milestone.is_completed
                      ? 'text-green-700'
                      : 'text-gray-700'
                  }`}
                >
                  {formatCurrency(milestone.amount)}
                </span>
              </div>
              {milestone.is_completed && milestone.completed_at && (
                <p className="text-xs text-green-600 mt-1">
                  Completed on {formatDate(milestone.completed_at)}
                </p>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
