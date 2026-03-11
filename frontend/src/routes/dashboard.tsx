import { createFileRoute } from '@tanstack/react-router'
import { QueryClientProvider } from '@tanstack/react-query'
import { ProtectedRoute } from '../components/ProtectedRoute'
import { DashboardLayout } from '../components/DashboardLayout'
import { DashboardSummary } from '../components/DashboardSummary'
import { QuickStatsCards } from '../components/QuickStatsCards'
import { RecentActivityFeed } from '../components/RecentActivityFeed'
import { queryClient } from '../lib/query/client'
import { useAuth } from '../lib/auth/context'
import { PiggyBank, Wallet, Target, BookOpen } from 'lucide-react'

export const Route = createFileRoute('/dashboard')({
  component: DashboardPage,
})

function DashboardPage() {
  return (
    <ProtectedRoute>
      <QueryClientProvider client={queryClient}>
        <DashboardLayout>
          <DashboardContent />
        </DashboardLayout>
      </QueryClientProvider>
    </ProtectedRoute>
  )
}

// Requirement 13.1: Dashboard layout with navigation
// Requirement 13.2: Dashboard summary component with data fetching
// Requirement 13.3: Quick stats cards
// Requirement 13.4: Recent activity feed
function DashboardContent() {
  const { user } = useAuth()

  return (
    <div className="space-y-8">
      {/* Welcome Header */}
      <div>
        <h1 className="text-3xl font-bold text-[var(--sea-ink)]">
          Welcome back, {user?.first_name}!
        </h1>
        <p className="mt-2 text-[var(--sea-ink-soft)]">
          Here's your financial overview for today
        </p>
      </div>

      {/* Dashboard Summary - Main Stats */}
      <DashboardSummary />

      {/* Quick Stats Cards */}
      <div>
        <h2 className="text-xl font-bold text-[var(--sea-ink)] mb-4">
          Quick Stats
        </h2>
        <QuickStatsCards />
      </div>

      {/* Recent Activity and Quick Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Recent Activity Feed */}
        <RecentActivityFeed />

        {/* Quick Actions */}
        <div className="border border-[var(--line)] rounded-xl p-6">
          <h2 className="text-xl font-bold text-[var(--sea-ink)] mb-6">
            Quick Actions
          </h2>
          <div className="space-y-3">
            <a
              href="/savings"
              className="flex items-center gap-4 p-4 border border-[var(--line)] rounded-lg hover:bg-[var(--link-bg-hover)] transition-colors"
            >
              <div className="p-3 bg-green-50 rounded-lg">
                <PiggyBank className="w-6 h-6 text-green-600" />
              </div>
              <div className="flex-1">
                <h3 className="font-semibold text-[var(--sea-ink)]">
                  Record Savings
                </h3>
                <p className="text-sm text-[var(--sea-ink-soft)]">
                  Track your savings transactions
                </p>
              </div>
            </a>

            <a
              href="/budget"
              className="flex items-center gap-4 p-4 border border-[var(--line)] rounded-lg hover:bg-[var(--link-bg-hover)] transition-colors"
            >
              <div className="p-3 bg-blue-50 rounded-lg">
                <Wallet className="w-6 h-6 text-blue-600" />
              </div>
              <div className="flex-1">
                <h3 className="font-semibold text-[var(--sea-ink)]">
                  Manage Budget
                </h3>
                <p className="text-sm text-[var(--sea-ink-soft)]">
                  Plan and track your spending
                </p>
              </div>
            </a>

            <a
              href="/goals"
              className="flex items-center gap-4 p-4 border border-[var(--line)] rounded-lg hover:bg-[var(--link-bg-hover)] transition-colors"
            >
              <div className="p-3 bg-purple-50 rounded-lg">
                <Target className="w-6 h-6 text-purple-600" />
              </div>
              <div className="flex-1">
                <h3 className="font-semibold text-[var(--sea-ink)]">Set Goals</h3>
                <p className="text-sm text-[var(--sea-ink-soft)]">
                  Create financial goals
                </p>
              </div>
            </a>

            <a
              href="/education"
              className="flex items-center gap-4 p-4 border border-[var(--line)] rounded-lg hover:bg-[var(--link-bg-hover)] transition-colors"
            >
              <div className="p-3 bg-yellow-50 rounded-lg">
                <BookOpen className="w-6 h-6 text-yellow-600" />
              </div>
              <div className="flex-1">
                <h3 className="font-semibold text-[var(--sea-ink)]">
                  Learn & Grow
                </h3>
                <p className="text-sm text-[var(--sea-ink-soft)]">
                  Financial education lessons
                </p>
              </div>
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}
