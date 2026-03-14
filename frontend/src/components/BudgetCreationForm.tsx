import { useState } from 'react'
import type { CreateBudgetRequest } from '../lib/types/api'

interface BudgetCreationFormProps {
  onSubmit: (data: CreateBudgetRequest) => Promise<void>
  onCancel: () => void
  isSubmitting: boolean
}

interface CategoryInput {
  id: string
  name: string
  allocated_amount: string
  color: string
}

const DEFAULT_COLORS = [
  '#3B82F6', // blue
  '#10B981', // green
  '#F59E0B', // amber
  '#EF4444', // red
  '#8B5CF6', // purple
  '#EC4899', // pink
  '#14B8A6', // teal
  '#F97316', // orange
]

export function BudgetCreationForm({
  onSubmit,
  onCancel,
  isSubmitting,
}: BudgetCreationFormProps) {
  const [totalBudget, setTotalBudget] = useState('')
  const [month, setMonth] = useState(() => {
    const now = new Date()
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
  })
  const [categories, setCategories] = useState<CategoryInput[]>([
    {
      id: '1',
      name: 'Food & Dining',
      allocated_amount: '',
      color: DEFAULT_COLORS[0],
    },
    {
      id: '2',
      name: 'Transportation',
      allocated_amount: '',
      color: DEFAULT_COLORS[1],
    },
    {
      id: '3',
      name: 'Entertainment',
      allocated_amount: '',
      color: DEFAULT_COLORS[2],
    },
  ])
  const [error, setError] = useState<string | null>(null)

  const addCategory = () => {
    const newId = String(Date.now())
    const colorIndex = categories.length % DEFAULT_COLORS.length
    setCategories([
      ...categories,
      {
        id: newId,
        name: '',
        allocated_amount: '',
        color: DEFAULT_COLORS[colorIndex],
      },
    ])
  }

  const removeCategory = (id: string) => {
    if (categories.length > 1) {
      setCategories(categories.filter((cat) => cat.id !== id))
    }
  }

  const updateCategory = (
    id: string,
    field: keyof CategoryInput,
    value: string
  ) => {
    setCategories(
      categories.map((cat) => (cat.id === id ? { ...cat, [field]: value } : cat))
    )
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    // Validation
    const budgetAmount = parseFloat(totalBudget)
    if (isNaN(budgetAmount) || budgetAmount <= 0) {
      setError('Please enter a valid budget amount')
      return
    }

    const validCategories = categories.filter(
      (cat) => cat.name.trim() && cat.allocated_amount.trim()
    )

    if (validCategories.length === 0) {
      setError('Please add at least one category with a name and amount')
      return
    }

    const totalAllocated = validCategories.reduce(
      (sum, cat) => sum + parseFloat(cat.allocated_amount),
      0
    )

    if (totalAllocated > budgetAmount) {
      setError(
        `Total allocated (${totalAllocated.toFixed(2)}) exceeds budget (${budgetAmount.toFixed(2)})`
      )
      return
    }

    try {
      await onSubmit({
        month: `${month}-01`,
        total_budget: budgetAmount,
        categories: validCategories.map((cat) => ({
          name: cat.name,
          allocated_amount: parseFloat(cat.allocated_amount),
          color: cat.color,
        })),
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create budget')
    }
  }

  const totalAllocated = categories.reduce((sum, cat) => {
    const amount = parseFloat(cat.allocated_amount)
    return sum + (isNaN(amount) ? 0 : amount)
  }, 0)

  const budgetAmount = parseFloat(totalBudget)
  const remaining = isNaN(budgetAmount) ? 0 : budgetAmount - totalAllocated

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-xl font-semibold text-gray-900 mb-6">
        Create Monthly Budget
      </h2>

      <form onSubmit={handleSubmit}>
        {/* Error Message */}
        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg">
            {error}
          </div>
        )}

        {/* Month and Total Budget */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Month
            </label>
            <input
              type="month"
              value={month}
              onChange={(e) => setMonth(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Total Budget
            </label>
            <input
              type="number"
              step="0.01"
              min="0"
              value={totalBudget}
              onChange={(e) => setTotalBudget(e.target.value)}
              placeholder="0.00"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              required
            />
          </div>
        </div>

        {/* Categories */}
        <div className="mb-6">
          <div className="flex justify-between items-center mb-4">
            <label className="block text-sm font-medium text-gray-700">
              Budget Categories
            </label>
            <button
              type="button"
              onClick={addCategory}
              className="text-sm text-blue-600 hover:text-blue-700 font-medium"
            >
              + Add Category
            </button>
          </div>

          <div className="space-y-3">
            {categories.map((category, index) => (
              <div
                key={category.id}
                className="flex gap-3 items-start bg-gray-50 p-3 rounded-lg"
              >
                <div className="flex-1">
                  <input
                    type="text"
                    value={category.name}
                    onChange={(e) =>
                      updateCategory(category.id, 'name', e.target.value)
                    }
                    placeholder="Category name"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
                <div className="w-32">
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={category.allocated_amount}
                    onChange={(e) =>
                      updateCategory(
                        category.id,
                        'allocated_amount',
                        e.target.value
                      )
                    }
                    placeholder="Amount"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
                <div className="w-16">
                  <input
                    type="color"
                    value={category.color}
                    onChange={(e) =>
                      updateCategory(category.id, 'color', e.target.value)
                    }
                    className="w-full h-10 border border-gray-300 rounded-lg cursor-pointer"
                  />
                </div>
                {categories.length > 1 && (
                  <button
                    type="button"
                    onClick={() => removeCategory(category.id)}
                    className="text-red-600 hover:text-red-700 p-2"
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
                        d="M6 18L18 6M6 6l12 12"
                      />
                    </svg>
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Budget Summary */}
        {!isNaN(budgetAmount) && budgetAmount > 0 && (
          <div className="mb-6 bg-blue-50 rounded-lg p-4">
            <div className="flex justify-between text-sm mb-2">
              <span className="text-gray-700">Total Budget:</span>
              <span className="font-semibold text-gray-900">
                ${budgetAmount.toFixed(2)}
              </span>
            </div>
            <div className="flex justify-between text-sm mb-2">
              <span className="text-gray-700">Total Allocated:</span>
              <span className="font-semibold text-gray-900">
                ${totalAllocated.toFixed(2)}
              </span>
            </div>
            <div className="flex justify-between text-sm pt-2 border-t border-blue-200">
              <span className="text-gray-700">Remaining:</span>
              <span
                className={`font-semibold ${
                  remaining < 0 ? 'text-red-600' : 'text-green-600'
                }`}
              >
                ${remaining.toFixed(2)}
              </span>
            </div>
          </div>
        )}

        {/* Actions */}
        <div className="flex gap-3 justify-end">
          <button
            type="button"
            onClick={onCancel}
            disabled={isSubmitting}
            className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isSubmitting}
            className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            {isSubmitting ? 'Creating...' : 'Create Budget'}
          </button>
        </div>
      </form>
    </div>
  )
}
