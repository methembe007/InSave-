import { createFileRoute } from '@tanstack/react-router'
import { BudgetPlanner } from '../components/BudgetPlanner'
import { ProtectedRoute } from '../components/ProtectedRoute'

export const Route = createFileRoute('/budget')({
  component: BudgetPage,
})

function BudgetPage() {
  return (
    <ProtectedRoute>
      <BudgetPlanner />
    </ProtectedRoute>
  )
}
