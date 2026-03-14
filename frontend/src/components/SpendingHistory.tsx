import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'

interface SpendingHistoryProps {
  budgetId: string
}

export function SpendingHistory({ budgetId }: SpendingHistoryProps) {
  const { api } = useAuth()
  const [selectedCategory, setSelectedCategory] = useState<string>('all')

  // Fetch spending history
  const { data: transactions, isLoading } = useQuery({
    queryKey: ['budget', budgetId, 'spending'],
    queryFn: () => api.budget.getSpendingHistory(budgetId),
    enabled: !!budgetId,
  })

  // Get unique categories for filter
  const categories = transactions
    ? Array.from(new Set(transactions.map((t) => t.category_id)))
    : []

  // Filter transactions by category
  const filteredTransactions =
    selectedCategory === 'all'
      ? transactions
      : transactions?.filter((t) => t.category_id === selectedCategory)

  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="animate-pulse space-y-3">
          <div className="h-4 bg-gray-200 rounded w-1/4"></div>
          <div className="h-16 bg-gray-200 rounded"></div>
          <div className="h-16 bg-gray-200 rounded"></div>
          <div className="h-16 bg-gray-200 rounded"></div>
        </div>
      </div>
    )
  }

  if (!transactions || transactions.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="text-center text-gray-500 py-8">
          <svg
            className="mx-auto h-12 w-12 text-gray-400 mb-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
            />
          </svg>
          <p>No spending transactions yet</p>
          <p className="text-sm mt-1">
            Record your first spending to see it here
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      {/* Category Filter */}
      {categories.length > 1 && (
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Filter by Category
          </label>
          <select
            value={selectedCategory}
            onChange={(e) => setSelectedCategory(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="all">All Categories</option>
            {categories.map((categoryId) => (
              <option key={categoryId} value={categoryId}>
                {categoryId}
              </option>
            ))}
          </select>
        </div>
      )}

      {/* Transactions List */}
      <div className="space-y-3 max-h-96 overflow-y-auto">
        {filteredTransactions && filteredTransactions.length > 0 ? (
          filteredTransactions.map((transaction) => {
            const date = new Date(transaction.date).toLocaleDateString('en-US', {
              month: 'short',
              day: 'numeric',
              year: 'numeric',
            })

            return (
              <div
                key={transaction.id}
                className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors"
              >
                <div className="flex justify-between items-start mb-2">
                  <div className="flex-1">
                    <div className="font-medium text-gray-900">
                      {transaction.description || 'No description'}
                    </div>
                    {transaction.merchant && (
                      <div className="text-sm text-gray-600 mt-1">
                        {transaction.merchant}
                      </div>
                    )}
                  </div>
                  <div className="text-right">
                    <div className="font-semibold text-red-600">
                      -${transaction.amount.toFixed(2)}
                    </div>
                    <div className="text-xs text-gray-500 mt-1">{date}</div>
                  </div>
                </div>
              </div>
            )
          })
        ) : (
          <div className="text-center text-gray-500 py-4">
            No transactions for selected category
          </div>
        )}
      </div>
    </div>
  )
}
