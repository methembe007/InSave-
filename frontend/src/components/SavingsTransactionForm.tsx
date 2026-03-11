import { useState } from 'react'
import { PlusCircle } from 'lucide-react'
import type { CreateSavingsRequest } from '../lib/types/api'

interface SavingsTransactionFormProps {
  onSubmit: (data: CreateSavingsRequest) => Promise<void>
  isSubmitting: boolean
}

const CATEGORIES = [
  'General',
  'Emergency Fund',
  'Vacation',
  'Education',
  'Investment',
  'Other',
]

export function SavingsTransactionForm({
  onSubmit,
  isSubmitting,
}: SavingsTransactionFormProps) {
  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')
  const [category, setCategory] = useState('General')
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    // Client-side validation
    const amountNum = parseFloat(amount)
    if (isNaN(amountNum) || amountNum <= 0) {
      setError('Amount must be greater than 0')
      return
    }

    try {
      await onSubmit({
        amount: amountNum,
        description: description.trim() || undefined,
        category,
        currency: 'USD',
      })

      // Reset form
      setAmount('')
      setDescription('')
      setCategory('General')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create transaction')
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6 border border-gray-200">
      <h2 className="text-xl font-semibold text-gray-900 mb-4 flex items-center gap-2">
        <PlusCircle className="w-5 h-5" />
        Add Savings
      </h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Amount Field */}
        <div>
          <label
            htmlFor="amount"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Amount ($)
          </label>
          <input
            type="number"
            id="amount"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            step="0.01"
            min="0.01"
            required
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0.00"
            disabled={isSubmitting}
          />
        </div>

        {/* Description Field */}
        <div>
          <label
            htmlFor="description"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Description (optional)
          </label>
          <input
            type="text"
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="What did you save for?"
            disabled={isSubmitting}
          />
        </div>

        {/* Category Field */}
        <div>
          <label
            htmlFor="category"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Category
          </label>
          <select
            id="category"
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={isSubmitting}
          >
            {CATEGORIES.map((cat) => (
              <option key={cat} value={cat}>
                {cat}
              </option>
            ))}
          </select>
        </div>

        {/* Error Message */}
        {error && (
          <div className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-md p-3">
            {error}
          </div>
        )}

        {/* Submit Button */}
        <button
          type="submit"
          disabled={isSubmitting}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {isSubmitting ? 'Saving...' : 'Add Savings'}
        </button>
      </form>
    </div>
  )
}
