import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { SavingsSummaryCards } from './SavingsSummaryCards'
import { SavingsTransactionForm } from './SavingsTransactionForm'
import { SavingsHistoryList } from './SavingsHistoryList'
import { StreakVisualization } from './StreakVisualization'
import { MonthlySavingsChart } from './MonthlySavingsChart'
import type { CreateSavingsRequest } from '../lib/types/api'

export function SavingsTracker() {
  const { api } = useAuth()
  const queryClient = useQueryClient()
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

  // Fetch savings summary
  const { data: summary, isLoading: summaryLoading } = useQuery({
    queryKey: ['savings', 'summary'],
    queryFn: () => api.savings.getSummary(),
  })

  // Fetch savings history
  const { data: history, isLoading: historyLoading } = useQuery({
    queryKey: ['savings', 'history'],
    queryFn: () => api.savings.getHistory({ limit: 50 }),
  })

  // Fetch streak data
  const { data: streak, isLoading: streakLoading } = useQuery({
    queryKey: ['savings', 'streak'],
    queryFn: () => api.savings.getStreak(),
  })

  // Create transaction mutation
  const createTransactionMutation = useMutation({
    mutationFn: (data: CreateSavingsRequest) =>
      api.savings.createTransaction(data),
    onSuccess: () => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ['savings'] })
      setSuccessMessage('Savings transaction created successfully! 🎉')
      setTimeout(() => setSuccessMessage(null), 5000)
    },
  })

  const handleCreateTransaction = async (data: CreateSavingsRequest) => {
    await createTransactionMutation.mutateAsync(data)
  }

  const isLoading = summaryLoading || streakLoading

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Savings Tracker</h1>
          <p className="mt-2 text-gray-600">
            Track your savings journey and build your financial discipline streak
          </p>
        </div>

        {/* Success Message */}
        {successMessage && (
          <div className="mb-6 bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded-lg">
            {successMessage}
          </div>
        )}

        {/* Summary Cards */}
        {isLoading ? (
          <div className="mb-8 text-center text-gray-500">Loading summary...</div>
        ) : (
          <SavingsSummaryCards summary={summary} />
        )}

        {/* Streak Visualization */}
        <div className="mb-8">
          <StreakVisualization
            currentStreak={streak?.current_streak || 0}
            longestStreak={streak?.longest_streak || 0}
            isLoading={streakLoading}
          />
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          {/* Transaction Form */}
          <div>
            <SavingsTransactionForm
              onSubmit={handleCreateTransaction}
              isSubmitting={createTransactionMutation.isPending}
            />
          </div>

          {/* Monthly Chart */}
          <div>
            <MonthlySavingsChart history={history || []} />
          </div>
        </div>

        {/* History List */}
        <SavingsHistoryList
          history={history || []}
          isLoading={historyLoading}
        />
      </div>
    </div>
  )
}
