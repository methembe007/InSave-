import { History } from 'lucide-react'
import type { SavingsTransaction } from '../lib/types/api'

interface SavingsHistoryListProps {
  history: SavingsTransaction[]
  isLoading: boolean
}

export function SavingsHistoryList({
  history,
  isLoading,
}: SavingsHistoryListProps) {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  return (
    <div className="bg-white rounded-lg shadow-md border border-gray-200">
      <div className="p-6 border-b border-gray-200">
        <h2 className="text-xl font-semibold text-gray-900 flex items-center gap-2">
          <History className="w-5 h-5" />
          Savings History
        </h2>
      </div>

      <div className="p-6">
        {isLoading ? (
          <div className="text-center text-gray-500 py-8">
            Loading history...
          </div>
        ) : history.length === 0 ? (
          <div className="text-center text-gray-500 py-8">
            No savings transactions yet. Start saving today!
          </div>
        ) : (
          <div className="space-y-3">
            {history.map((transaction) => (
              <div
                key={transaction.id}
                className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
              >
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <div className="flex-1">
                      <p className="font-medium text-gray-900">
                        {transaction.description || 'Savings'}
                      </p>
                      <div className="flex items-center gap-2 mt-1">
                        <span className="text-sm text-gray-500">
                          {transaction.category}
                        </span>
                        <span className="text-gray-300">•</span>
                        <span className="text-sm text-gray-500">
                          {formatDate(transaction.created_at)}
                        </span>
                        <span className="text-gray-300">•</span>
                        <span className="text-sm text-gray-500">
                          {formatTime(transaction.created_at)}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-lg font-semibold text-green-600">
                    +${transaction.amount.toFixed(2)}
                  </p>
                  <p className="text-xs text-gray-500">{transaction.currency}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
