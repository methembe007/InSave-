import { createFileRoute } from '@tanstack/react-router'
import { GoalManager } from '../components/GoalManager'

export const Route = createFileRoute('/goals')({
  component: GoalsPage,
})

function GoalsPage() {
  return <GoalManager />
}
