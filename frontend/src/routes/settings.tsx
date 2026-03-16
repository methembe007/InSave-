import { createFileRoute } from '@tanstack/react-router'
import { ProtectedRoute } from '../components/ProtectedRoute'
import { DashboardLayout } from '../components/DashboardLayout'
import { SettingsPage } from '../components/SettingsPage'

export const Route = createFileRoute('/settings')({
  component: SettingsRoute,
})

function SettingsRoute() {
  return (
    <ProtectedRoute>
      <DashboardLayout>
        <SettingsPage />
      </DashboardLayout>
    </ProtectedRoute>
  )
}
