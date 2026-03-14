import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { GoalCreationForm } from './GoalCreationForm'
import { GoalCard } from './GoalCard'
import type { CreateGoalRequest } from '../lib/types/api'

export function GoalManager() {
  const { api } = useAuth()
  const queryClient = useQueryClient()
  const [successMessage, setSuccessMessage] = useState<string | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)

  // Fetch active goals
  const { data: goals, isLoading: goalsLoading } = useQuery({
    queryKey: ['goals', 'active'],
    queryFn: () => api.goals.getActiveGoals(),
  })

  // Create goal mutation
  const createGoalMutation = useMutation({
    mutationFn: (data: CreateGoalRequest) => api.goals.createGoal(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] })
      setSuccessMessage('Goal created successfully! 🎯')
      setShowCreateForm(false)
      setTimeout(() => setSuccessMessage(null), 5000)
    },
  })

  const handleCreateGoal = async (data: CreateGoalRequest) => {
    await createGoalMutation.mutateAsync(data)
  }

  const hasGoals = goals && goals.length > 0

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Financial Goals</h1>
            <p className="mt-2 text-gray-600">
              Set and track your long-term financial objectives
            </p>
          </div>
          {!showCreateForm && (
            <button
              onClick={() => setShowCreateForm(true)}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Create New Goal
            </button>
          )}
        </div>

        {/* Success Message */}
        {successMessage && (
          <div className="mb-6 bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded-lg">
            {successMessage}
          </div>
        )}

        {/* Goal Creation Form */}
        {showCreateForm && (
          <div className="mb-8">
            <GoalCreationForm
              onSubmit={handleCreateGoal}
              onCancel={() => setShowCreateForm(false)}
              isSubmitting={createGoalMutation.isPending}
            />
          </div>
        )}

        {/* No Goals State */}
        {!hasGoals && !showCreateForm && !goalsLoading && (
          <div className="bg-white rounded-lg shadow-md p-8 text-center">
            <div className="mb-4">
              <svg
                className="mx-auto h-12 w-12 text-gray-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"
                />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              No Goals Yet
            </h3>
            <p className="text-gray-600 mb-6">
              Create your first financial goal to start working toward your dreams
            </p>
            <button
              onClick={() => setShowCreateForm(true)}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Create Your First Goal
            </button>
          </div>
        )}

        {/* Loading State */}
        {goalsLoading && (
          <div className="text-center text-gray-500">Loading goals...</div>
        )}

        {/* Goals List */}
        {hasGoals && !showCreateForm && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {goals.map((goal) => (
              <GoalCard key={goal.id} goal={goal} />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
