import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import type { Goal } from '../lib/types/api'

interface GoalContributionFormProps {
  goal: Goal
  onSuccess: (wasCompleted: boolean) => void
}

export function GoalContributionForm({
  goal,
  onSuccess,
}: GoalContributionFormProps) {
  const { api } = useAuth()
  const queryClient = useQueryClient()
  const [amount, setAmount] = useState<number>(0)
  const [error, setError] = useState<string>('')

  // Add progress mutation with optimistic updates
  const addProgressMutation = useMutation({
    mutationFn: (contributionAmount: number) =>
      api.goals.addProgress(goal.id, contributionAmount),
    onMutate: async (contributionAmount) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: ['goals'] })

      // Snapshot previous value
      const previousGoals = queryClient.getQueryData(['goals', 'active'])

      // Optimistically update
      queryClient.setQueryData(['goals', 'active'], (old: Goal[] | undefined) => {
        if (!old) return old
        return old.map((g) => {
          if (g.id === goal.id) {
            const newCurrentAmount = g.current_amount + contributionAmount
            const newProgressPercent =
              (newCurrentAmount / g.target_amount) * 100
            const newStatus =
              newCurrentAmount >= g.target_amount ? 'completed' : g.status
            return {
              ...g,
              current_amount: newCurrentAmount,
              progress_percent: newProgressPercent,
              status: newStatus,
            }
          }
          return g
        })
      })

      return { previousGoals }
    },
    onError: (_err, _variables, context) => {
      // Rollback on error
      if (context?.previousGoals) {
        queryClient.setQueryData(['goals', 'active'], context.previousGoals)
      }
    },
    onSuccess: (updatedGoal) => {
      // Refetch to ensure consistency
      queryClient.invalidateQueries({ queryKey: ['goals'] })
      queryClient.invalidateQueries({
        queryKey: ['goals', goal.id, 'milestones'],
      })
      
      const wasCompleted = updatedGoal.status === 'completed'
      onSuccess(wasCompleted)
      setAmount(0)
    },
  })

  const validateAmount = (): boolean => {
    if (amount <= 0) {
      setError('Amount must be greater than 0')
      return false
    }
    setError('')
    return true
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateAmount()) {
      return
    }

    try {
      await addProgressMutation.mutateAsync(amount)
    } catch (error) {
      console.error('Failed to add contribution:', error)
      setError('Failed to add contribution. Please try again.')
    }
  }

  const handleAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseFloat(e.target.value) || 0
    setAmount(value)
    if (error) {
      setError('')
    }
  }

  const remainingAmount = goal.target_amount - goal.current_amount

  return (
    <div className="bg-gray-50 rounded-lg p-4">
      <h4 className="text-sm font-medium text-gray-900 mb-3">
        Add Contribution
      </h4>
      <form onSubmit={handleSubmit} className="space-y-3">
        <div>
          <label
            htmlFor="contribution-amount"
            className="block text-sm text-gray-700 mb-1"
          >
            Amount
          </label>
          <div className="relative">
            <span className="absolute left-3 top-2 text-gray-500">$</span>
            <input
              type="number"
              id="contribution-amount"
              value={amount || ''}
              onChange={handleAmountChange}
              step="0.01"
              min="0"
              className={`w-full pl-8 pr-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                error ? 'border-red-500' : 'border-gray-300'
              }`}
              placeholder="0.00"
            />
          </div>
          {error && <p className="mt-1 text-sm text-red-600">{error}</p>}
          <p className="mt-1 text-xs text-gray-500">
            Remaining: ${remainingAmount.toFixed(2)}
          </p>
        </div>

        <button
          type="submit"
          disabled={addProgressMutation.isPending}
          className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
        >
          {addProgressMutation.isPending ? 'Adding...' : 'Add Contribution'}
        </button>
      </form>
    </div>
  )
}
