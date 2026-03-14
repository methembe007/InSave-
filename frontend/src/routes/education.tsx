import { createFileRoute } from '@tanstack/react-router'
import { EducationSection } from '../components/EducationSection'
import { ProtectedRoute } from '../components/ProtectedRoute'

export const Route = createFileRoute('/education')({
  component: EducationPage,
})

function EducationPage() {
  return (
    <ProtectedRoute>
      <EducationSection />
    </ProtectedRoute>
  )
}
