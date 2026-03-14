import { createFileRoute } from '@tanstack/react-router'
import { LessonDetail } from '../components/LessonDetail'
import { ProtectedRoute } from '../components/ProtectedRoute'

export const Route = createFileRoute('/education/lessons/$lessonId')({
  component: LessonDetailPage,
})

function LessonDetailPage() {
  const { lessonId } = Route.useParams()

  return (
    <ProtectedRoute>
      <LessonDetail lessonId={lessonId} />
    </ProtectedRoute>
  )
}
