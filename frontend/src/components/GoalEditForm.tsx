import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import type { Goal, UpdateGoalRequest } from '../lib/types/api'

interface GoalEditFormProps {
  goal: Goal
  onCancel: () => void
  onSuccess: () => void
}

export function GoalEditForm({ goal, onCancel, onSuccess }: GoalEditFormProps) {
  const { api } = useAuth()
  const queryClient = useQueryClient()
  const [formData, setFormData] = useState<UpdateGoalRequest>({
    title: goal.title,
    description: goal.description,
    target_amount: goal.target_amount,
    target_date: goal.target_date.split('T')[0], // Extract date part
    status: goal.status,
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  // Update goal mutation
  const updateGoalMutation = useMutation({
    mutationFn: (data: UpdateGoalRequest) =>
      api.goals.updateGoal(goal.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] })
      onSuccess()
    },
  })

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {}

    if (formData.title && !formData.title.trim()) {
      newErrors.title = 'Title cannot be empty'
    }

    if (formData.target_amount && formData.target_amount <= 0) {
      newErrors.target_amount = 'Target amount must be greater than 0'
    }

    if (formData.target_date) {
      const targetDate = new Date(formData.target_date)
      const today = new Date()
      today.setHours(0, 0, 0, 0)

      if (targetDate <= today) {
        newErrors.target_date = 'Target date must be in the future'
      }
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    try {
      await updateGoalMutation.mutateAsync(formData)
    } catch (error) {
      console.error('Failed to update goal:', error)
    }
  }

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'target_amount' ? parseFloat(value) || 0 : value,
    }))
    // Clear error for this field
    if (errors[name]) {
      setErrors((prev) => {
        const newErrors = { ...prev }
        delete newErrors[name]
        return newErrors
      })
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-xl font-semibold text-gray-900 mb-4">Edit Goal</h3>
      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Title */}
        <div>
          <label
            htmlFor="edit-title"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Goal Title
          </label>
          <input
            type="text"
            id="edit-title"
            name="title"
            value={formData.title}
            onChange={handleChange}
            className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
              errors.title ? 'border-red-500' : 'border-gray-300'
            }`}
          />
          {errors.title && (
            <p className="mt-1 text-sm text-red-600">{errors.title}</p>
          )}
        </div>

        {/* Description */}
        <div>
          <label
            htmlFor="edit-description"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Description
          </label>
          <textarea
            id="edit-description"
            name="description"
            value={formData.description}
            onChange={handleChange}
            rows={3}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Target Amount */}
        <div>
          <label
            htmlFor="edit-target-amount"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Target Amount
          </label>
          <div className="relative">
            <span className="absolute left-3 top-2 text-gray-500">$</span>
            <input
              type="number"
              id="edit-target-amount"
              name="target_amount"
              value={formData.target_amount || ''}
              onChange={handleChange}
              step="0.01"
              min="0"
              className={`w-full pl-8 pr-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                errors.target_amount ? 'border-red-500' : 'border-gray-300'
              }`}
            />
          </div>
          {errors.target_amount && (
            <p className="mt-1 text-sm text-red-600">{errors.target_amount}</p>
          )}
        </div>

        {/* Target Date */}
        <div>
          <label
            htmlFor="edit-target-date"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Target Date
          </label>
          <input
            type="date"
            id="edit-target-date"
            name="target_date"
            value={formData.target_date}
            onChange={handleChange}
            className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
              errors.target_date ? 'border-red-500' : 'border-gray-300'
            }`}
          />
          {errors.target_date && (
            <p className="mt-1 text-sm text-red-600">{errors.target_date}</p>
          )}
        </div>

        {/* Status */}
        <div>
          <label
            htmlFor="edit-status"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Status
          </label>
          <select
            id="edit-status"
            name="status"
            value={formData.status}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="active">Active</option>
            <option value="paused">Paused</option>
            <option value="completed">Completed</option>
          </select>
        </div>

        {/* Buttons */}
        <div className="flex gap-3 pt-4">
          <button
            type="submit"
            disabled={updateGoalMutation.isPending}
            className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {updateGoalMutation.isPending ? 'Saving...' : 'Save Changes'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            disabled={updateGoalMutation.isPending}
            className="flex-1 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 disabled:bg-gray-100 disabled:cursor-not-allowed transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}
