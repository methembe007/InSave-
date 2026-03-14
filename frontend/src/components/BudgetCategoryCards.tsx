import type { BudgetCategory } from '../lib/types/api'

interface BudgetCategoryCardsProps {
  categories: BudgetCategory[]
}

export function BudgetCategoryCards({ categories }: BudgetCategoryCardsProps) {
  if (!categories || categories.length === 0) {
    return (
      <div className="text-center text-gray-500 py-8">
        No categories found
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {categories.map((category) => {
        const percentageUsed =
          category.allocated_amount > 0
            ? (category.spent_amount / category.allocated_amount) * 100
            : 0

        const getStatusColor = () => {
          if (percentageUsed >= 100) return 'red'
          if (percentageUsed >= 80) return 'yellow'
          return 'green'
        }

        const statusColor = getStatusColor()

        const colorClasses = {
          red: {
            bg: 'bg-red-50',
            border: 'border-red-200',
            text: 'text-red-900',
            progress: 'bg-red-500',
            badge: 'bg-red-100 text-red-800',
          },
          yellow: {
            bg: 'bg-yellow-50',
            border: 'border-yellow-200',
            text: 'text-yellow-900',
            progress: 'bg-yellow-500',
            badge: 'bg-yellow-100 text-yellow-800',
          },
          green: {
            bg: 'bg-green-50',
            border: 'border-green-200',
            text: 'text-green-900',
            progress: 'bg-green-500',
            badge: 'bg-green-100 text-green-800',
          },
        }

        const colors = colorClasses[statusColor]

        return (
          <div
            key={category.id}
            className={`${colors.bg} border ${colors.border} rounded-lg p-4 transition-all hover:shadow-md`}
          >
            {/* Category Header */}
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                <div
                  className="w-4 h-4 rounded-full"
                  style={{ backgroundColor: category.color }}
                ></div>
                <h3 className="font-semibold text-gray-900">
                  {category.name}
                </h3>
              </div>
              <span
                className={`text-xs font-medium px-2 py-1 rounded-full ${colors.badge}`}
              >
                {percentageUsed.toFixed(0)}%
              </span>
            </div>

            {/* Amounts */}
            <div className="space-y-2 mb-3">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Allocated:</span>
                <span className="font-medium text-gray-900">
                  ${category.allocated_amount.toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Spent:</span>
                <span className={`font-medium ${colors.text}`}>
                  ${category.spent_amount.toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Remaining:</span>
                <span
                  className={`font-medium ${
                    category.remaining_amount < 0
                      ? 'text-red-600'
                      : 'text-green-600'
                  }`}
                >
                  ${category.remaining_amount.toFixed(2)}
                </span>
              </div>
            </div>

            {/* Progress Bar */}
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className={`h-2 rounded-full transition-all ${colors.progress}`}
                style={{ width: `${Math.min(percentageUsed, 100)}%` }}
              ></div>
            </div>

            {/* Warning Message */}
            {percentageUsed >= 100 && (
              <div className="mt-3 text-xs text-red-600 font-medium">
                ⚠️ Budget exceeded!
              </div>
            )}
            {percentageUsed >= 80 && percentageUsed < 100 && (
              <div className="mt-3 text-xs text-yellow-600 font-medium">
                ⚠️ Approaching limit
              </div>
            )}
          </div>
        )
      })}
    </div>
  )
}
