import type { Budget } from '../lib/types/api'

interface BudgetOverviewProps {
  budget: Budget | undefined
  isLoading: boolean
}

export function BudgetOverview({ budget, isLoading }: BudgetOverviewProps) {
  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="h-20 bg-gray-200 rounded"></div>
            <div className="h-20 bg-gray-200 rounded"></div>
            <div className="h-20 bg-gray-200 rounded"></div>
          </div>
        </div>
      </div>
    )
  }

  if (!budget) {
    return null
  }

  const percentageUsed = (budget.total_spent / budget.total_budget) * 100
  const month = new Date(budget.month).toLocaleDateString('en-US', {
    month: 'long',
    year: 'numeric',
  })

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="mb-4">
        <h2 className="text-lg font-semibold text-gray-900">
          Budget for {month}
        </h2>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Total Budget */}
        <div className="bg-blue-50 rounded-lg p-4">
          <div className="text-sm text-blue-600 font-medium mb-1">
            Total Budget
          </div>
          <div className="text-2xl font-bold text-blue-900">
            ${budget.total_budget.toFixed(2)}
          </div>
        </div>

        {/* Total Spent */}
        <div className="bg-orange-50 rounded-lg p-4">
          <div className="text-sm text-orange-600 font-medium mb-1">
            Total Spent
          </div>
          <div className="text-2xl font-bold text-orange-900">
            ${budget.total_spent.toFixed(2)}
          </div>
          <div className="text-xs text-orange-600 mt-1">
            {percentageUsed.toFixed(1)}% used
          </div>
        </div>

        {/* Remaining Budget */}
        <div className="bg-green-50 rounded-lg p-4">
          <div className="text-sm text-green-600 font-medium mb-1">
            Remaining
          </div>
          <div className="text-2xl font-bold text-green-900">
            ${budget.remaining_budget.toFixed(2)}
          </div>
          <div className="text-xs text-green-600 mt-1">
            {(100 - percentageUsed).toFixed(1)}% left
          </div>
        </div>
      </div>

      {/* Progress Bar */}
      <div className="mt-6">
        <div className="flex justify-between text-sm text-gray-600 mb-2">
          <span>Budget Progress</span>
          <span>{percentageUsed.toFixed(1)}%</span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-3">
          <div
            className={`h-3 rounded-full transition-all ${
              percentageUsed >= 100
                ? 'bg-red-500'
                : percentageUsed >= 80
                ? 'bg-yellow-500'
                : 'bg-green-500'
            }`}
            style={{ width: `${Math.min(percentageUsed, 100)}%` }}
          ></div>
        </div>
      </div>
    </div>
  )
}
