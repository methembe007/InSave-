import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts'
import { ShoppingCart, TrendingDown, TrendingUp, DollarSign, Store } from 'lucide-react'
import type { TimePeriod } from '../lib/types/api'

interface SpendingAnalysisChartProps {
  period: TimePeriod
}

// Requirement 13.1: Fetch spending analysis data for selected period
// Requirement 13.2: Create pie chart for category breakdown, display stats
export function SpendingAnalysisChart({ period }: SpendingAnalysisChartProps) {
  const { api } = useAuth()
  
  const { data: analysis, isLoading, error } = useQuery({
    queryKey: ['spending-analysis', period],
    queryFn: () => api.analytics.getSpendingAnalysis(period),
  })

  if (isLoading) {
    return (
      <div className="border border-[var(--line)] rounded-xl p-8 animate-pulse">
        <div className="h-96 bg-gray-200 rounded"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="border border-red-200 bg-red-50 rounded-xl p-6">
        <p className="text-red-800">Unable to load spending analysis. Please try again later.</p>
      </div>
    )
  }

  if (!analysis) return null

  const COLORS = [
    '#3b82f6', // blue
    '#10b981', // green
    '#f59e0b', // amber
    '#ef4444', // red
    '#8b5cf6', // violet
    '#ec4899', // pink
    '#06b6d4', // cyan
    '#f97316', // orange
  ]

  const chartData = analysis.category_breakdown.map((cat: { category: string; amount: number; percentage: number }) => ({
    name: cat.category,
    value: cat.amount,
    percentage: cat.percentage,
  }))

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value)
  }

  return (
    <div className="border border-[var(--line)] rounded-xl p-8">
      <div className="flex items-center gap-3 mb-6">
        <ShoppingCart className="w-6 h-6 text-[var(--link)]" />
        <h2 className="text-2xl font-bold text-[var(--sea-ink)]">Spending Analysis</h2>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        {/* Total Spending */}
        <div className="bg-blue-50 rounded-lg p-4 border border-blue-200">
          <div className="flex items-center gap-2 mb-2">
            <DollarSign className="w-5 h-5 text-blue-600" />
            <span className="text-sm font-medium text-blue-900">Total Spending</span>
          </div>
          <p className="text-2xl font-bold text-blue-900">
            {formatCurrency(analysis.total_spending)}
          </p>
        </div>

        {/* Daily Average */}
        <div className="bg-purple-50 rounded-lg p-4 border border-purple-200">
          <div className="flex items-center gap-2 mb-2">
            <DollarSign className="w-5 h-5 text-purple-600" />
            <span className="text-sm font-medium text-purple-900">Daily Average</span>
          </div>
          <p className="text-2xl font-bold text-purple-900">
            {formatCurrency(analysis.daily_average)}
          </p>
        </div>

        {/* Comparison to Previous Period */}
        <div className={`rounded-lg p-4 border ${
          analysis.comparison_to_previous > 0
            ? 'bg-red-50 border-red-200'
            : 'bg-green-50 border-green-200'
        }`}>
          <div className="flex items-center gap-2 mb-2">
            {analysis.comparison_to_previous > 0 ? (
              <TrendingUp className="w-5 h-5 text-red-600" />
            ) : (
              <TrendingDown className="w-5 h-5 text-green-600" />
            )}
            <span className={`text-sm font-medium ${
              analysis.comparison_to_previous > 0 ? 'text-red-900' : 'text-green-900'
            }`}>
              vs Previous Period
            </span>
          </div>
          <p className={`text-2xl font-bold ${
            analysis.comparison_to_previous > 0 ? 'text-red-900' : 'text-green-900'
          }`}>
            {analysis.comparison_to_previous > 0 ? '+' : ''}
            {analysis.comparison_to_previous.toFixed(1)}%
          </p>
        </div>
      </div>

      {/* Category Breakdown Chart and Top Merchants */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Pie Chart */}
        <div>
          <h3 className="text-lg font-semibold text-[var(--sea-ink)] mb-4">
            Category Breakdown
          </h3>
          {chartData.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={chartData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={(entry: any) => `${entry.name}: ${entry.percentage.toFixed(1)}%`}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {chartData.map((_entry: any, index: number) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip
                  formatter={(value: any) => formatCurrency(Number(value))}
                  contentStyle={{
                    backgroundColor: 'white',
                    border: '1px solid #e5e7eb',
                    borderRadius: '8px',
                  }}
                />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <div className="h-[300px] flex items-center justify-center text-[var(--sea-ink-soft)]">
              No spending data for this period
            </div>
          )}
        </div>

        {/* Top Merchants */}
        <div>
          <h3 className="text-lg font-semibold text-[var(--sea-ink)] mb-4 flex items-center gap-2">
            <Store className="w-5 h-5" />
            Top Merchants
          </h3>
          {analysis.top_merchants.length > 0 ? (
            <div className="space-y-3">
              {analysis.top_merchants.slice(0, 5).map((merchant: any, index: number) => (
                <div
                  key={index}
                  className="flex items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-200"
                >
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-[var(--link)] text-white rounded-full flex items-center justify-center font-bold text-sm">
                      {index + 1}
                    </div>
                    <div>
                      <p className="font-medium text-[var(--sea-ink)]">{merchant.merchant}</p>
                      <p className="text-sm text-[var(--sea-ink-soft)]">
                        {merchant.transaction_count} transaction{merchant.transaction_count !== 1 ? 's' : ''}
                      </p>
                    </div>
                  </div>
                  <p className="font-bold text-[var(--sea-ink)]">
                    {formatCurrency(merchant.amount)}
                  </p>
                </div>
              ))}
            </div>
          ) : (
            <div className="h-[300px] flex items-center justify-center text-[var(--sea-ink-soft)]">
              No merchant data available
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
