import { createFileRoute } from '@tanstack/react-router'
import { QueryClientProvider } from '@tanstack/react-query'
import { ProtectedRoute } from '../components/ProtectedRoute'
import { DashboardLayout } from '../components/DashboardLayout'
import { ProfilePage } from '../components/ProfilePage'
import { queryClient } from '../lib/query/client'

export const Route = createFileRoute('/profile')({
  component: ProfileRoute,
})

function ProfileRoute() {
  return (
    <ProtectedRoute>
      <QueryClientProvider client={queryClient}>
        <DashboardLayout>
          <ProfilePage />
        </DashboardLayout>
      </QueryClientProvider>
    </ProtectedRoute>
  )
}
