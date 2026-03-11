import { ArrowUpCircle, ArrowDownCircle, Calendar } from 'lucide-react'
import { useRecentSavings } from '../lib/hooks/useDashboardData'

// Requirement 13.4: Recent activity feed
// Fetch recent savings and spending transactions
// Display in chronological order
// Show transaction type, amount, description, date
export function RecentActivityFeed() {
  const { data: savings, isLoading } = useRecentSavings(10)

  if (isLoading) {
    return (
      <div className="border border-[var(--line)] rounded-xl p-6">
        <h2 className="text-xl font-bold text-[var(--sea-ink)] mb-4">
          Recent Activity
        </h2>
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="flex items-center gap-4 animate-pulse">
              <div className="w-10 h-10 bg-[var(--link-bg-hover)] rounded-full" />
              <div className="flex-1">
                <div className="h-4 bg-[var(--link-bg-hover)] rounded w-32 mb-2" />
                <div className="h-3 bg-[var(--link-bg-hover)] rounded w-24" />
              </div>
              <div className="h-5 bg-[var(--link-bg-hover)] rounded w-16" />
            </div>
          ))}
        </div>
      </div>
    )
  }

  // Combine and sort transactions by date
  const transactions = [
    ...(savings || []).map((t) => ({
      id: t.id,
      type: 'savings' as const,
      amount: t.amount,
      description: t.description || 'Savings deposit',
      category: t.category,
      date: new Date(t.created_at),
    })),
    // TODO: Add spending transactions when available
  ].sort((a, b) => b.date.getTime() - a.date.getTime())

  return (
    <div className="border border-[var(--line)] rounded-xl p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-bold text-[var(--sea-ink)]">Recent Activity</h2>
        <a
          href="/savings"
          className="text-sm text-[var(--sea-ink-soft)] hover:text-[var(--sea-ink)] transition-colors"
        >
          View all →
        </a>
      </div>

      {transactions.length === 0 ? (
        <div className="text-center py-8">
          <Calendar className="w-12 h-12 text-[var(--sea-ink-soft)] mx-auto mb-3" />
          <p className="text-[var(--sea-ink-soft)]">No recent activity</p>
          <p className="text-sm text-[var(--sea-ink-soft)] mt-1">
            Start by recording your first savings transaction
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {transactions.map((transaction) => (
            <TransactionItem key={transaction.id} transaction={transaction} />
          ))}
        </div>
      )}
    </div>
  )
}

interface Transaction {
  id: string
  type: 'savings' | 'spending'
  amount: number
  description: string
  category: string
  date: Date
}

function TransactionItem({ transaction }: { transaction: Transaction }) {
  const isSavings = transaction.type === 'savings'

  return (
    <div className="flex items-center gap-4 p-3 rounded-lg hover:bg-[var(--link-bg-hover)] transition-colors">
      {/* Icon */}
      <div
        className={`
        p-2 rounded-full
        ${isSavings ? 'bg-green-50' : 'bg-red-50'}
      `}
      >
        {isSavings ? (
          <ArrowUpCircle className="w-6 h-6 text-green-600" />
        ) : (
          <ArrowDownCircle className="w-6 h-6 text-red-600" />
        )}
      </div>

      {/* Details */}
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-[var(--sea-ink)] truncate">
          {transaction.description}
        </p>
        <div className="flex items-center gap-2 mt-1">
          {transaction.category && (
            <span className="text-xs px-2 py-0.5 bg-[var(--link-bg-hover)] rounded-full text-[var(--sea-ink-soft)]">
              {transaction.category}
            </span>
          )}
          <span className="text-xs text-[var(--sea-ink-soft)]">
            {formatDate(transaction.date)}
          </span>
        </div>
      </div>

      {/* Amount */}
      <div className="text-right">
        <p
          className={`
          text-lg font-bold
          ${isSavings ? 'text-green-600' : 'text-red-600'}
        `}
        >
          {isSavings ? '+' : '-'}${transaction.amount.toFixed(2)}
        </p>
      </div>
    </div>
  )
}

function formatDate(date: Date): string {
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins}m ago`
  if (diffHours < 24) return `${diffHours}h ago`
  if (diffDays < 7) return `${diffDays}d ago`

  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  })
}
