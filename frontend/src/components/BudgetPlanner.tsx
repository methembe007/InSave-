import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { BudgetOverview } from './BudgetOverview'
import { BudgetCreationForm } from './BudgetCreationForm'
import { BudgetCategoryCards } from './BudgetCategoryCards'
import { SpendingTransactionForm } from './SpendingTransactionForm'
import { BudgetAlertsDisplay } from './BudgetAlertsDisplay'
import { SpendingHistory } from './SpendingHistory'
import type { CreateBudgetRequest, SpendingRequest } from '../lib/types/api'

export function BudgetPlanner() {
  const { api } = useAuth()
  const queryClient = useQueryClient()
  const [successMessage, setSuccessMessage] = useState<string | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)

  // Fetch current budget
  const { data: budget, isLoading: budgetLoading, error: budgetError } = useQuery({
    queryKey: ['budget', 'current'],
    queryFn: () => api.budget.getCurrentBudget(),
    retry: false,
  })

  // Fetch budget alerts
  const { data: alerts, isLoading: alertsLoading } = useQuery({
    queryKey: ['budget', 'alerts'],
    queryFn: () => api.budget.getAlerts(),
    enabled: !!budget,
  })

  // Create budget mutation
  const createBudgetMutation = useMutation({
    mutationFn: (data: CreateBudgetRequest) => api.budget.createBudget(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['budget'] })
      setSuccessMessage('Budget created successfully! 🎉')
      setShowCreateForm(false)
      setTimeout(() => setSuccessMessage(null), 5000)
    },
  })

  // Record spending mutation
  const recordSpendingMutation = useMutation({
    mutationFn: (data: SpendingRequest) => api.budget.recordSpending(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['budget'] })
      setSuccessMessage('Spending recorded successfully! 💰')
      setTimeout(() => setSuccessMessage(null), 5000)
    },
  })

  const handleCreateBudget = async (data: CreateBudgetRequest) => {
    await createBudgetMutation.mutateAsync(data)
  }

  const handleRecordSpending = async (data: SpendingRequest) => {
    await recordSpendingMutation.mutateAsync(data)
  }

  const noBudget = budgetError || !budget

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Budget Planner</h1>
            <p className="mt-2 text-gray-600">
              Plan your spending and track your budget progress
            </p>
          </div>
          {!noBudget && !showCreateForm && (
            <button
              onClick={() => setShowCreateForm(true)}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Create New Budget
            </button>
          )}
        </div>

        {/* Success Message */}
        {successMessage && (
          <div className="mb-6 bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded-lg">
            {successMessage}
          </div>
        )}

        {/* Budget Alerts */}
        {alerts && alerts.length > 0 && (
          <div className="mb-8">
            <BudgetAlertsDisplay alerts={alerts} isLoading={alertsLoading} />
          </div>
        )}

        {/* No Budget State */}
        {noBudget && !showCreateForm && (
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
                  d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
                />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              No Budget Found
            </h3>
            <p className="text-gray-600 mb-6">
              Create your first budget to start tracking your spending
            </p>
            <button
              onClick={() => setShowCreateForm(true)}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Create Budget
            </button>
          </div>
        )}

        {/* Budget Creation Form */}
        {showCreateForm && (
          <div className="mb-8">
            <BudgetCreationForm
              onSubmit={handleCreateBudget}
              onCancel={() => setShowCreateForm(false)}
              isSubmitting={createBudgetMutation.isPending}
            />
          </div>
        )}

        {/* Budget Overview and Management */}
        {budget && !showCreateForm && (
          <>
            {/* Budget Overview */}
            <div className="mb-8">
              <BudgetOverview budget={budget} isLoading={budgetLoading} />
            </div>

            {/* Category Cards */}
            <div className="mb-8">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">
                Budget Categories
              </h2>
              <BudgetCategoryCards categories={budget.categories} />
            </div>

            {/* Main Content Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
              {/* Spending Form */}
              <div>
                <h2 className="text-xl font-semibold text-gray-900 mb-4">
                  Record Spending
                </h2>
                <SpendingTransactionForm
                  budget={budget}
                  onSubmit={handleRecordSpending}
                  isSubmitting={recordSpendingMutation.isPending}
                />
              </div>

              {/* Spending History */}
              <div>
                <h2 className="text-xl font-semibold text-gray-900 mb-4">
                  Recent Spending
                </h2>
                <SpendingHistory budgetId={budget.id} />
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
