import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { TrendingUp } from 'lucide-react'
import type { SavingsTransaction } from '../lib/types/api'

interface MonthlySavingsChartProps {
  history: SavingsTransaction[]
}

export function MonthlySavingsChart({ history }: MonthlySavingsChartProps) {
  // Group transactions by month
  const monthlyData = history.reduce((acc, transaction) => {
    const date = new Date(transaction.created_at)
    const monthKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
    const monthLabel = date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
    })

    if (!acc[monthKey]) {
      acc[monthKey] = {
        month: monthLabel,
        amount: 0,
        count: 0,
      }
    }

    acc[monthKey].amount += transaction.amount
    acc[monthKey].count += 1

    return acc
  }, {} as Record<string, { month: string; amount: number; count: number }>)

  // Convert to array and sort by date
  const chartData = Object.entries(monthlyData)
    .sort(([a], [b]) => a.localeCompare(b))
    .slice(-6) // Last 6 months
    .map(([_, data]) => ({
      month: data.month,
      amount: parseFloat(data.amount.toFixed(2)),
      count: data.count,
    }))

  const totalAmount = chartData.reduce((sum, item) => sum + item.amount, 0)

  return (
    <div className="bg-white rounded-lg shadow-md p-6 border border-gray-200">
      <div className="mb-4">
        <h2 className="text-xl font-semibold text-gray-900 flex items-center gap-2">
          <TrendingUp className="w-5 h-5" />
          Monthly Savings
        </h2>
        <p className="text-sm text-gray-600 mt-1">
          Last 6 months • Total: ${totalAmount.toFixed(2)}
        </p>
      </div>

      {chartData.length === 0 ? (
        <div className="text-center text-gray-500 py-12">
          No data to display yet. Start saving to see your progress!
        </div>
      ) : (
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
            <XAxis
              dataKey="month"
              tick={{ fill: '#6b7280', fontSize: 12 }}
              tickLine={{ stroke: '#e5e7eb' }}
            />
            <YAxis
              tick={{ fill: '#6b7280', fontSize: 12 }}
              tickLine={{ stroke: '#e5e7eb' }}
              tickFormatter={(value) => `$${value}`}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: '#fff',
                border: '1px solid #e5e7eb',
                borderRadius: '8px',
                padding: '12px',
              }}
              formatter={(value: any, name: any) => {
                if (name === 'amount' && typeof value === 'number') {
                  return [`$${value.toFixed(2)}`, 'Saved']
                }
                return [value, name]
              }}
              labelStyle={{ fontWeight: 'bold', marginBottom: '8px' }}
            />
            <Bar
              dataKey="amount"
              fill="#3b82f6"
              radius={[8, 8, 0, 0]}
              maxBarSize={60}
            />
          </BarChart>
        </ResponsiveContainer>
      )}
    </div>
  )
}
