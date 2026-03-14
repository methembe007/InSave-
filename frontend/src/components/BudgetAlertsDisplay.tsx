import type { BudgetAlert } from '../lib/types/api'

interface BudgetAlertsDisplayProps {
  alerts: BudgetAlert[]
  isLoading: boolean
}

export function BudgetAlertsDisplay({
  alerts,
  isLoading,
}: BudgetAlertsDisplayProps) {
  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div className="space-y-3">
            <div className="h-16 bg-gray-200 rounded"></div>
            <div className="h-16 bg-gray-200 rounded"></div>
          </div>
        </div>
      </div>
    )
  }

  if (!alerts || alerts.length === 0) {
    return null
  }

  // Sort alerts: critical first, then by percentage descending
  const sortedAlerts = [...alerts].sort((a, b) => {
    if (a.alert_type === 'critical' && b.alert_type !== 'critical') return -1
    if (a.alert_type !== 'critical' && b.alert_type === 'critical') return 1
    return b.percentage_used - a.percentage_used
  })

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="flex items-center gap-2 mb-4">
        <svg
          className="w-5 h-5 text-orange-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
        <h2 className="text-lg font-semibold text-gray-900">Budget Alerts</h2>
      </div>

      <div className="space-y-3">
        {sortedAlerts.map((alert, index) => {
          const isCritical = alert.alert_type === 'critical'
          const bgColor = isCritical ? 'bg-red-50' : 'bg-yellow-50'
          const borderColor = isCritical
            ? 'border-red-200'
            : 'border-yellow-200'
          const textColor = isCritical ? 'text-red-800' : 'text-yellow-800'
          const iconColor = isCritical ? 'text-red-500' : 'text-yellow-500'

          return (
            <div
              key={`${alert.category_name}-${index}`}
              className={`${bgColor} border ${borderColor} rounded-lg p-4`}
            >
              <div className="flex items-start gap-3">
                <div className={`${iconColor} mt-0.5`}>
                  {isCritical ? (
                    <svg
                      className="w-5 h-5"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                        clipRule="evenodd"
                      />
                    </svg>
                  ) : (
                    <svg
                      className="w-5 h-5"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                        clipRule="evenodd"
                      />
                    </svg>
                  )}
                </div>
                <div className="flex-1">
                  <div className="flex items-center justify-between mb-1">
                    <h3 className={`font-semibold ${textColor}`}>
                      {alert.category_name}
                    </h3>
                    <span
                      className={`text-sm font-medium ${textColor} px-2 py-1 rounded ${
                        isCritical ? 'bg-red-100' : 'bg-yellow-100'
                      }`}
                    >
                      {alert.percentage_used.toFixed(0)}%
                    </span>
                  </div>
                  <p className={`text-sm ${textColor}`}>{alert.message}</p>
                </div>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
