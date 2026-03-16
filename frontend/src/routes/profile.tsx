import { createFileRoute } from '@tanstack/react-router'
import { ProtectedRoute } from '../components/ProtectedRoute'
import { DashboardLayout } from '../components/DashboardLayout'
import { ProfilePage } from '../components/ProfilePage'

export const Route = createFileRoute('/profile')({
  component: ProfileRoute,
})

function ProfileRoute() {
  return (
    <ProtectedRoute>
      <DashboardLayout>
        <ProfilePage />
      </DashboardLayout>
    </ProtectedRoute>
  )
}
