import { createFileRoute } from '@tanstack/react-router'
import { SavingsTracker } from '../components/SavingsTracker'
import { ProtectedRoute } from '../components/ProtectedRoute'

export const Route = createFileRoute('/savings')({
  component: SavingsPage,
})

function SavingsPage() {
  return (
    <ProtectedRoute>
      <SavingsTracker />
    </ProtectedRoute>
  )
}
