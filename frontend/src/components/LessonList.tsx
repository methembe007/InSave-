import { useState, useMemo } from 'react'
import type { Lesson } from '../lib/types/api'

interface LessonListProps {
  lessons: Lesson[]
  isLoading: boolean
}

export function LessonList({ lessons, isLoading }: LessonListProps) {
  const [selectedCategory, setSelectedCategory] = useState<string>('all')

  // Extract unique categories
  const categories = useMemo(() => {
    const cats = new Set(lessons.map((lesson) => lesson.category))
    return ['all', ...Array.from(cats)]
  }, [lessons])

  // Filter lessons by category
  const filteredLessons = useMemo(() => {
    if (selectedCategory === 'all') {
      return lessons
    }
    return lessons.filter((lesson) => lesson.category === selectedCategory)
  }, [lessons, selectedCategory])

  // Sort lessons by order
  const sortedLessons = useMemo(() => {
    return [...filteredLessons].sort((a, b) => a.order - b.order)
  }, [filteredLessons])

  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="text-center text-gray-500">Loading lessons...</div>
      </div>
    )
  }

  if (lessons.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="text-center text-gray-500">No lessons available yet</div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow">
      {/* Category Filter */}
      <div className="border-b border-gray-200 px-6 py-4">
        <div className="flex flex-wrap gap-2">
          {categories.map((category) => (
            <button
              key={category}
              onClick={() => setSelectedCategory(category)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                selectedCategory === category
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {category.charAt(0).toUpperCase() + category.slice(1)}
            </button>
          ))}
        </div>
      </div>

      {/* Lessons List */}
      <div className="divide-y divide-gray-200">
        {sortedLessons.length === 0 ? (
          <div className="px-6 py-8 text-center text-gray-500">
            No lessons found in this category
          </div>
        ) : (
          sortedLessons.map((lesson) => (
            <a
              key={lesson.id}
              href={`/education/lessons/${lesson.id}`}
              className="block px-6 py-4 hover:bg-gray-50 transition-colors"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-lg font-semibold text-gray-900">
                      {lesson.title}
                    </h3>
                    {lesson.is_completed && (
                      <span className="inline-flex items-center text-green-600">
                        <svg
                          className="w-5 h-5"
                          fill="currentColor"
                          viewBox="0 0 20 20"
                        >
                          <path
                            fillRule="evenodd"
                            d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                            clipRule="evenodd"
                          />
                        </svg>
                      </span>
                    )}
                  </div>
                  <p className="text-gray-600 mb-3">{lesson.description}</p>
                  <div className="flex flex-wrap items-center gap-4 text-sm text-gray-500">
                    <span className="flex items-center gap-1">
                      <svg
                        className="w-4 h-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                      {lesson.duration_minutes} min
                    </span>
                    <span
                      className={`px-2 py-1 rounded text-xs font-medium ${
                        lesson.difficulty === 'beginner'
                          ? 'bg-green-100 text-green-800'
                          : lesson.difficulty === 'intermediate'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {lesson.difficulty.charAt(0).toUpperCase() +
                        lesson.difficulty.slice(1)}
                    </span>
                    <span className="px-2 py-1 bg-gray-100 text-gray-700 rounded text-xs font-medium">
                      {lesson.category}
                    </span>
                  </div>
                </div>
                <div className="ml-4">
                  <svg
                    className="w-5 h-5 text-gray-400"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 5l7 7-7 7"
                    />
                  </svg>
                </div>
              </div>
            </a>
          ))
        )}
      </div>
    </div>
  )
}
