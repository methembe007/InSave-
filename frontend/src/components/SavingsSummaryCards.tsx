import { DollarSign, TrendingUp, Calendar } from 'lucide-react'
import type { SavingsSummary } from '../lib/types/api'

interface SavingsSummaryCardsProps {
  summary?: SavingsSummary
}

export function SavingsSummaryCards({ summary }: SavingsSummaryCardsProps) {
  if (!summary) {
    return null
  }

  const cards = [
    {
      title: 'Total Saved',
      value: `$${summary.total_saved.toFixed(2)}`,
      icon: DollarSign,
      color: 'bg-blue-500',
    },
    {
      title: 'This Month',
      value: `$${summary.this_month_saved.toFixed(2)}`,
      icon: Calendar,
      color: 'bg-green-500',
    },
    {
      title: 'Monthly Average',
      value: `$${summary.monthly_average.toFixed(2)}`,
      icon: TrendingUp,
      color: 'bg-purple-500',
    },
  ]

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      {cards.map((card) => {
        const Icon = card.icon
        return (
          <div
            key={card.title}
            className="bg-white rounded-lg shadow-md p-6 border border-gray-200"
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">{card.title}</p>
                <p className="mt-2 text-3xl font-bold text-gray-900">
                  {card.value}
                </p>
              </div>
              <div className={`${card.color} p-3 rounded-lg`}>
                <Icon className="w-6 h-6 text-white" />
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
