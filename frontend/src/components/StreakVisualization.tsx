import { Flame, Trophy } from 'lucide-react'

interface StreakVisualizationProps {
  currentStreak: number
  longestStreak: number
  isLoading: boolean
}

export function StreakVisualization({
  currentStreak,
  longestStreak,
  isLoading,
}: StreakVisualizationProps) {
  const getMotivationalMessage = (streak: number): string => {
    if (streak === 0) return "Start your savings journey today!"
    if (streak === 1) return "Great start! Keep it going!"
    if (streak < 7) return "You're building momentum!"
    if (streak < 30) return "Amazing consistency! Keep it up!"
    if (streak < 100) return "You're on fire! Incredible discipline!"
    return "Legendary streak! You're a savings champion!"
  }

  const getStreakColor = (streak: number): string => {
    if (streak === 0) return 'text-gray-400'
    if (streak < 7) return 'text-orange-500'
    if (streak < 30) return 'text-orange-600'
    return 'text-red-600'
  }

  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6 border border-gray-200">
        <div className="text-center text-gray-500">Loading streak...</div>
      </div>
    )
  }

  return (
    <div className="bg-gradient-to-r from-orange-50 to-red-50 rounded-lg shadow-md p-8 border border-orange-200">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Current Streak */}
        <div className="text-center">
          <div className="flex justify-center mb-4">
            <div className="relative">
              <Flame
                className={`w-20 h-20 ${getStreakColor(currentStreak)}`}
                fill="currentColor"
              />
              {currentStreak > 0 && (
                <div className="absolute inset-0 flex items-center justify-center">
                  <span className="text-2xl font-bold text-white drop-shadow-lg">
                    {currentStreak}
                  </span>
                </div>
              )}
            </div>
          </div>
          <h3 className="text-2xl font-bold text-gray-900 mb-2">
            Current Streak
          </h3>
          <p className="text-4xl font-extrabold text-orange-600 mb-2">
            {currentStreak} {currentStreak === 1 ? 'day' : 'days'}
          </p>
          <p className="text-sm text-gray-600 italic">
            {getMotivationalMessage(currentStreak)}
          </p>
        </div>

        {/* Longest Streak */}
        <div className="text-center">
          <div className="flex justify-center mb-4">
            <div className="relative">
              <Trophy className="w-20 h-20 text-yellow-500" fill="currentColor" />
              {longestStreak > 0 && (
                <div className="absolute inset-0 flex items-center justify-center">
                  <span className="text-2xl font-bold text-white drop-shadow-lg">
                    {longestStreak}
                  </span>
                </div>
              )}
            </div>
          </div>
          <h3 className="text-2xl font-bold text-gray-900 mb-2">
            Longest Streak
          </h3>
          <p className="text-4xl font-extrabold text-yellow-600 mb-2">
            {longestStreak} {longestStreak === 1 ? 'day' : 'days'}
          </p>
          <p className="text-sm text-gray-600 italic">
            {longestStreak === currentStreak && longestStreak > 0
              ? "You're at your personal best!"
              : longestStreak > 0
              ? 'Your best achievement so far!'
              : 'Start saving to set your record!'}
          </p>
        </div>
      </div>

      {/* Progress Bar */}
      {longestStreak > 0 && (
        <div className="mt-6">
          <div className="flex justify-between text-sm text-gray-600 mb-2">
            <span>Progress to longest streak</span>
            <span>
              {Math.round((currentStreak / longestStreak) * 100)}%
            </span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-3">
            <div
              className="bg-gradient-to-r from-orange-500 to-red-500 h-3 rounded-full transition-all duration-500"
              style={{
                width: `${Math.min((currentStreak / longestStreak) * 100, 100)}%`,
              }}
            />
          </div>
        </div>
      )}
    </div>
  )
}
