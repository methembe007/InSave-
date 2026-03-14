import { createFileRoute } from '@tanstack/react-router'
import { ProtectedRoute } from '../components/ProtectedRoute'
import { DashboardLayout } from '../components/DashboardLayout'
import { AnalyticsDashboard } from '../components/AnalyticsDashboard'

export const Route = createFileRoute('/analytics')({
  component: AnalyticsPage,
})

// Requirement 13.1: Analytics page layout with spending analysis, patterns, and recommendations
// Requirement 14.1: Display financial health score prominently
function AnalyticsPage() {
  return (
    <ProtectedRoute>
      <DashboardLayout>
        <AnalyticsDashboard />
      </DashboardLayout>
    </ProtectedRoute>
  )
}
