import { useState } from 'react'
import { FinancialHealthDisplay } from './FinancialHealthDisplay'
import { SpendingAnalysisChart } from './SpendingAnalysisChart'
import { SavingsPatternsDisplay } from './SavingsPatternsDisplay'
import { RecommendationsList } from './RecommendationsList'
import { TrendingUp, Calendar } from 'lucide-react'

// Requirement 13.1: Analytics page layout with spending analysis, patterns, and recommendations
// Requirement 14.1: Display financial health score prominently
export function AnalyticsDashboard() {
  const [selectedPeriod, setSelectedPeriod] = useState<'week' | 'month' | 'quarter'>('month')

  const getDateRange = (period: 'week' | 'month' | 'quarter') => {
    const end = new Date()
    const start = new Date()
    
    switch (period) {
      case 'week':
        start.setDate(end.getDate() - 7)
        break
      case 'month':
        start.setMonth(end.getMonth() - 1)
        break
      case 'quarter':
        start.setMonth(end.getMonth() - 3)
        break
    }
    
    return {
      start_date: start.toISOString().split('T')[0],
      end_date: end.toISOString().split('T')[0],
    }
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[var(--sea-ink)] flex items-center gap-3">
            <TrendingUp className="w-8 h-8 text-[var(--link)]" />
            Analytics Dashboard
          </h1>
          <p className="mt-2 text-[var(--sea-ink-soft)]">
            Insights into your financial health and spending patterns
          </p>
        </div>

        {/* Period Selector */}
        <div className="flex items-center gap-2 border border-[var(--line)] rounded-lg p-1">
          <Calendar className="w-4 h-4 text-[var(--sea-ink-soft)] ml-2" />
          <button
            onClick={() => setSelectedPeriod('week')}
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              selectedPeriod === 'week'
                ? 'bg-[var(--link)] text-white'
                : 'text-[var(--sea-ink-soft)] hover:text-[var(--sea-ink)]'
            }`}
          >
            Week
          </button>
          <button
            onClick={() => setSelectedPeriod('month')}
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              selectedPeriod === 'month'
                ? 'bg-[var(--link)] text-white'
                : 'text-[var(--sea-ink-soft)] hover:text-[var(--sea-ink)]'
            }`}
          >
            Month
          </button>
          <button
            onClick={() => setSelectedPeriod('quarter')}
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              selectedPeriod === 'quarter'
                ? 'bg-[var(--link)] text-white'
                : 'text-[var(--sea-ink-soft)] hover:text-[var(--sea-ink)]'
            }`}
          >
            Quarter
          </button>
        </div>
      </div>

      {/* Financial Health Score - Prominent Display */}
      <FinancialHealthDisplay />

      {/* Spending Analysis */}
      <SpendingAnalysisChart period={getDateRange(selectedPeriod)} />

      {/* Savings Patterns and Recommendations */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <SavingsPatternsDisplay />
        <RecommendationsList />
      </div>
    </div>
  )
}
