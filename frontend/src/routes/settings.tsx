import { createFileRoute } from '@tanstack/react-router'
import { QueryClientProvider } from '@tanstack/react-query'
import { ProtectedRoute } from '../components/ProtectedRoute'
import { DashboardLayout } from '../components/DashboardLayout'
import { SettingsPage } from '../components/SettingsPage'
import { queryClient } from '../lib/query/client'

export const Route = createFileRoute('/settings')({
  component: SettingsRoute,
})

function SettingsRoute() {
  return (
    <ProtectedRoute>
      <QueryClientProvider client={queryClient}>
        <DashboardLayout>
          <SettingsPage />
        </DashboardLayout>
      </QueryClientProvider>
    </ProtectedRoute>
  )
}
