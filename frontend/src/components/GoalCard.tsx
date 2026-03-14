import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { GoalContributionForm } from './GoalContributionForm'
import { GoalMilestones } from './GoalMilestones'
import { GoalEditForm } from './GoalEditForm'
import type { Goal } from '../lib/types/api'

interface GoalCardProps {
  goal: Goal
}

export function GoalCard({ goal }: GoalCardProps) {
  const { api } = useAuth()
  const queryClient = useQueryClient()
  const [showContributionForm, setShowContributionForm] = useState(false)
  const [showMilestones, setShowMilestones] = useState(false)
  const [showEditForm, setShowEditForm] = useState(false)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [showCelebration, setShowCelebration] = useState(false)

  // Fetch milestones
  const { data: milestones } = useQuery({
    queryKey: ['goals', goal.id, 'milestones'],
    queryFn: () => api.goals.getMilestones(goal.id),
    enabled: showMilestones,
  })

  // Delete goal mutation
  const deleteGoalMutation = useMutation({
    mutationFn: () => api.goals.deleteGoal(goal.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] })
      setShowDeleteConfirm(false)
    },
  })

  const handleDelete = async () => {
    await deleteGoalMutation.mutateAsync()
  }

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: goal.currency,
    }).format(amount)
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800'
      case 'active':
        return 'bg-blue-100 text-blue-800'
      case 'paused':
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getProgressColor = (percent: number) => {
    if (percent >= 100) return 'bg-green-500'
    if (percent >= 75) return 'bg-blue-500'
    if (percent >= 50) return 'bg-yellow-500'
    return 'bg-orange-500'
  }

  const handleContributionSuccess = (wasCompleted: boolean) => {
    setShowContributionForm(false)
    if (wasCompleted) {
      setShowCelebration(true)
      setTimeout(() => setShowCelebration(false), 5000)
    }
  }

  if (showEditForm) {
    return (
      <GoalEditForm
        goal={goal}
        onCancel={() => setShowEditForm(false)}
        onSuccess={() => setShowEditForm(false)}
      />
    )
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6 relative overflow-hidden">
      {/* Celebration Animation */}
      {showCelebration && (
        <div className="absolute inset-0 bg-green-500 bg-opacity-90 flex items-center justify-center z-10 animate-pulse">
          <div className="text-center text-white">
            <div className="text-6xl mb-4">🎉</div>
            <div className="text-2xl font-bold">Goal Completed!</div>
            <div className="text-lg">Congratulations!</div>
          </div>
        </div>
      )}

      {/* Status Badge */}
      <div className="flex justify-between items-start mb-4">
        <span
          className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(
            goal.status
          )}`}
        >
          {goal.status.charAt(0).toUpperCase() + goal.status.slice(1)}
        </span>
        <div className="flex gap-2">
          <button
            onClick={() => setShowEditForm(true)}
            className="text-gray-400 hover:text-blue-600 transition-colors"
            title="Edit goal"
          >
            <svg
              className="w-5 h-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
              />
            </svg>
          </button>
          <button
            onClick={() => setShowDeleteConfirm(true)}
            className="text-gray-400 hover:text-red-600 transition-colors"
            title="Delete goal"
          >
            <svg
              className="w-5 h-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              />
            </svg>
          </button>
        </div>
      </div>

      {/* Goal Title and Description */}
      <h3 className="text-xl font-semibold text-gray-900 mb-2">
        {goal.title}
      </h3>
      {goal.description && (
        <p className="text-gray-600 text-sm mb-4">{goal.description}</p>
      )}

      {/* Amounts */}
      <div className="mb-4">
        <div className="flex justify-between text-sm mb-1">
          <span className="text-gray-600">Progress</span>
          <span className="font-medium text-gray-900">
            {goal.progress_percent.toFixed(1)}%
          </span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-3 mb-2">
          <div
            className={`h-3 rounded-full transition-all duration-500 ${getProgressColor(
              goal.progress_percent
            )}`}
            style={{ width: `${Math.min(goal.progress_percent, 100)}%` }}
          />
        </div>
        <div className="flex justify-between text-sm">
          <span className="text-gray-600">
            {formatCurrency(goal.current_amount)}
          </span>
          <span className="font-medium text-gray-900">
            {formatCurrency(goal.target_amount)}
          </span>
        </div>
      </div>

      {/* Target Date */}
      <div className="mb-4 text-sm text-gray-600">
        <span className="font-medium">Target Date:</span> {formatDate(goal.target_date)}
      </div>

      {/* Action Buttons */}
      <div className="space-y-2">
        {goal.status === 'active' && (
          <button
            onClick={() => setShowContributionForm(!showContributionForm)}
            className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            {showContributionForm ? 'Cancel' : 'Add Contribution'}
          </button>
        )}
        <button
          onClick={() => setShowMilestones(!showMilestones)}
          className="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
        >
          {showMilestones ? 'Hide Milestones' : 'View Milestones'}
        </button>
      </div>

      {/* Contribution Form */}
      {showContributionForm && (
        <div className="mt-4">
          <GoalContributionForm
            goal={goal}
            onSuccess={handleContributionSuccess}
          />
        </div>
      )}

      {/* Milestones */}
      {showMilestones && (
        <div className="mt-4">
          <GoalMilestones
            goalId={goal.id}
            milestones={milestones || []}
            isLoading={!milestones}
          />
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      {showDeleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md mx-4">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              Delete Goal?
            </h3>
            <p className="text-gray-600 mb-6">
              Are you sure you want to delete "{goal.title}"? This action cannot
              be undone and will remove all associated milestones.
            </p>
            <div className="flex gap-3">
              <button
                onClick={handleDelete}
                disabled={deleteGoalMutation.isPending}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:bg-gray-400 transition-colors"
              >
                {deleteGoalMutation.isPending ? 'Deleting...' : 'Delete'}
              </button>
              <button
                onClick={() => setShowDeleteConfirm(false)}
                disabled={deleteGoalMutation.isPending}
                className="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
