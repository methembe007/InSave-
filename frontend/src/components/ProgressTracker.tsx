import type { EducationProgress } from '../lib/types/api'

interface ProgressTrackerProps {
  progress?: EducationProgress
}

export function ProgressTracker({ progress }: ProgressTrackerProps) {
  if (!progress) {
    return null
  }

  const { total_lessons, completed_lessons, progress_percent, current_streak } =
    progress

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-xl font-semibold text-gray-900 mb-6">
        Your Progress
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        {/* Total Lessons */}
        <div className="text-center">
          <div className="text-3xl font-bold text-blue-600">
            {total_lessons}
          </div>
          <div className="text-sm text-gray-600 mt-1">Total Lessons</div>
        </div>

        {/* Completed Lessons */}
        <div className="text-center">
          <div className="text-3xl font-bold text-green-600">
            {completed_lessons}
          </div>
          <div className="text-sm text-gray-600 mt-1">Completed</div>
        </div>

        {/* Current Streak */}
        <div className="text-center">
          <div className="text-3xl font-bold text-orange-600">
            {current_streak}
          </div>
          <div className="text-sm text-gray-600 mt-1">Day Streak</div>
        </div>
      </div>

      {/* Progress Bar */}
      <div>
        <div className="flex justify-between items-center mb-2">
          <span className="text-sm font-medium text-gray-700">
            Overall Progress
          </span>
          <span className="text-sm font-semibold text-blue-600">
            {progress_percent.toFixed(1)}%
          </span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
          <div
            className="bg-blue-600 h-full rounded-full transition-all duration-500 ease-out"
            style={{ width: `${progress_percent}%` }}
          />
        </div>
        <div className="mt-2 text-xs text-gray-500 text-center">
          {completed_lessons} of {total_lessons} lessons completed
        </div>
      </div>
    </div>
  )
}
