import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { LessonList } from './LessonList'
import { ProgressTracker } from './ProgressTracker'

export function EducationSection() {
  const { api } = useAuth()

  // Fetch lessons
  const { data: lessons, isLoading: lessonsLoading } = useQuery({
    queryKey: ['education', 'lessons'],
    queryFn: () => api.education.getLessons(),
  })

  // Fetch progress
  const { data: progress, isLoading: progressLoading } = useQuery({
    queryKey: ['education', 'progress'],
    queryFn: () => api.education.getProgress(),
  })

  const isLoading = lessonsLoading || progressLoading

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Financial Education</h1>
          <p className="mt-2 text-gray-600">
            Build your financial knowledge with curated lessons and resources
          </p>
        </div>

        {/* Progress Tracker */}
        {isLoading ? (
          <div className="mb-8 text-center text-gray-500">Loading progress...</div>
        ) : (
          <div className="mb-8">
            <ProgressTracker progress={progress} />
          </div>
        )}

        {/* Lesson List */}
        <LessonList lessons={lessons || []} isLoading={lessonsLoading} />
      </div>
    </div>
  )
}
