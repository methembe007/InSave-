import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import {
  Settings as SettingsIcon,
  Bell,
  Mail,
  Smartphone,
  Clock,
  DollarSign,
  Palette,
  Save,
  Trash2,
  AlertTriangle,
} from 'lucide-react'
import type { UserPreferences } from '../types/api'

// Requirement 3.3: Display and update user preferences
// Requirement 3.4: Account deletion with confirmation
export function SettingsPage() {
  const { api, logout } = useAuth()
  const queryClient = useQueryClient()
  const [formData, setFormData] = useState<UserPreferences | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [deleteConfirmText, setDeleteConfirmText] = useState('')

  // Fetch user preferences
  const { data: preferences, isLoading } = useQuery<UserPreferences>({
    queryKey: ['user', 'preferences'],
    queryFn: () => api.user.getPreferences(),
  })

  // Initialize formData when preferences load
  if (preferences && !formData) {
    setFormData(preferences)
  }

  // Update preferences mutation
  const updatePreferencesMutation = useMutation({
    mutationFn: (data: UserPreferences) => api.user.updatePreferences(data),
    onSuccess: (updatedPreferences) => {
      queryClient.setQueryData(['user', 'preferences'], updatedPreferences)
      setSuccess('Settings saved successfully!')
      setError(null)
      setTimeout(() => setSuccess(null), 3000)
    },
    onError: (err: any) => {
      setError(err.message || 'Failed to update settings')
      setSuccess(null)
    },
  })

  // Delete account mutation
  const deleteAccountMutation = useMutation({
    mutationFn: () => api.user.deleteAccount(),
    onSuccess: async () => {
      // Logout and redirect after successful deletion
      await logout()
    },
    onError: (err: any) => {
      setError(err.message || 'Failed to delete account')
      setShowDeleteConfirm(false)
      setDeleteConfirmText('')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (formData) {
      updatePreferencesMutation.mutate(formData)
    }
  }

  const handleToggle = (field: keyof UserPreferences) => {
    if (formData) {
      setFormData((prev) => ({
        ...prev!,
        [field]: !prev![field],
      }))
    }
  }

  const handleInputChange = (field: keyof UserPreferences, value: string) => {
    if (formData) {
      setFormData((prev) => ({
        ...prev!,
        [field]: value,
      }))
    }
  }

  const handleDeleteAccount = () => {
    if (deleteConfirmText === 'DELETE') {
      deleteAccountMutation.mutate()
    } else {
      setError('Please type DELETE to confirm')
    }
  }

  if (isLoading || !formData) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-[var(--sea-ink-soft)]">Loading settings...</div>
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-[var(--sea-ink)]">Settings</h1>
        <p className="mt-2 text-[var(--sea-ink-soft)]">
          Manage your preferences and account settings
        </p>
      </div>

      {/* Success/Error Messages */}
      {success && (
        <div className="p-4 bg-green-50 border border-green-200 rounded-lg text-green-800">
          {success}
        </div>
      )}
      {error && (
        <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-800">
          {error}
        </div>
      )}

      {/* Settings Form */}
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* General Settings */}
        <div className="border border-[var(--line)] rounded-xl p-6 bg-white">
          <div className="flex items-center gap-3 mb-6">
            <SettingsIcon className="w-5 h-5 text-[var(--sea-ink)]" />
            <h2 className="text-xl font-semibold text-[var(--sea-ink)]">
              General Settings
            </h2>
          </div>

          <div className="space-y-4">
            {/* Currency */}
            <div>
              <label
                htmlFor="currency"
                className="flex items-center gap-2 text-sm font-medium text-[var(--sea-ink)] mb-2"
              >
                <DollarSign className="w-4 h-4" />
                Currency
              </label>
              <select
                id="currency"
                value={formData.currency}
                onChange={(e) => handleInputChange('currency', e.target.value)}
                className="w-full px-4 py-2 border border-[var(--line)] rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--sea-ink)]"
              >
                <option value="USD">USD - US Dollar</option>
                <option value="EUR">EUR - Euro</option>
                <option value="GBP">GBP - British Pound</option>
                <option value="JPY">JPY - Japanese Yen</option>
                <option value="CAD">CAD - Canadian Dollar</option>
                <option value="AUD">AUD - Australian Dollar</option>
              </select>
            </div>

            {/* Theme */}
            <div>
              <label
                htmlFor="theme"
                className="flex items-center gap-2 text-sm font-medium text-[var(--sea-ink)] mb-2"
              >
                <Palette className="w-4 h-4" />
                Theme
              </label>
              <select
                id="theme"
                value={formData.theme}
                onChange={(e) => handleInputChange('theme', e.target.value)}
                className="w-full px-4 py-2 border border-[var(--line)] rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--sea-ink)]"
              >
                <option value="light">Light</option>
                <option value="dark">Dark</option>
                <option value="auto">Auto (System)</option>
              </select>
            </div>
          </div>
        </div>

        {/* Notification Settings */}
        <div className="border border-[var(--line)] rounded-xl p-6 bg-white">
          <div className="flex items-center gap-3 mb-6">
            <Bell className="w-5 h-5 text-[var(--sea-ink)]" />
            <h2 className="text-xl font-semibold text-[var(--sea-ink)]">
              Notifications
            </h2>
          </div>

          <div className="space-y-4">
            {/* Master Notifications Toggle */}
            <div className="flex items-center justify-between p-4 bg-[var(--link-bg-hover)] rounded-lg">
              <div className="flex items-center gap-3">
                <Bell className="w-5 h-5 text-[var(--sea-ink-soft)]" />
                <div>
                  <p className="font-medium text-[var(--sea-ink)]">
                    Enable Notifications
                  </p>
                  <p className="text-sm text-[var(--sea-ink-soft)]">
                    Master switch for all notifications
                  </p>
                </div>
              </div>
              <button
                type="button"
                onClick={() => handleToggle('notifications_enabled')}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  formData.notifications_enabled
                    ? 'bg-green-600'
                    : 'bg-gray-300'
                }`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    formData.notifications_enabled
                      ? 'translate-x-6'
                      : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {/* Email Notifications */}
            <div className="flex items-center justify-between p-4 border border-[var(--line)] rounded-lg">
              <div className="flex items-center gap-3">
                <Mail className="w-5 h-5 text-[var(--sea-ink-soft)]" />
                <div>
                  <p className="font-medium text-[var(--sea-ink)]">
                    Email Notifications
                  </p>
                  <p className="text-sm text-[var(--sea-ink-soft)]">
                    Receive notifications via email
                  </p>
                </div>
              </div>
              <button
                type="button"
                onClick={() => handleToggle('email_notifications')}
                disabled={!formData.notifications_enabled}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  formData.email_notifications && formData.notifications_enabled
                    ? 'bg-green-600'
                    : 'bg-gray-300'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    formData.email_notifications && formData.notifications_enabled
                      ? 'translate-x-6'
                      : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {/* Push Notifications */}
            <div className="flex items-center justify-between p-4 border border-[var(--line)] rounded-lg">
              <div className="flex items-center gap-3">
                <Smartphone className="w-5 h-5 text-[var(--sea-ink-soft)]" />
                <div>
                  <p className="font-medium text-[var(--sea-ink)]">
                    Push Notifications
                  </p>
                  <p className="text-sm text-[var(--sea-ink-soft)]">
                    Receive push notifications on your device
                  </p>
                </div>
              </div>
              <button
                type="button"
                onClick={() => handleToggle('push_notifications')}
                disabled={!formData.notifications_enabled}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  formData.push_notifications && formData.notifications_enabled
                    ? 'bg-green-600'
                    : 'bg-gray-300'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    formData.push_notifications && formData.notifications_enabled
                      ? 'translate-x-6'
                      : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {/* Savings Reminders */}
            <div className="flex items-center justify-between p-4 border border-[var(--line)] rounded-lg">
              <div className="flex items-center gap-3">
                <Clock className="w-5 h-5 text-[var(--sea-ink-soft)]" />
                <div>
                  <p className="font-medium text-[var(--sea-ink)]">
                    Savings Reminders
                  </p>
                  <p className="text-sm text-[var(--sea-ink-soft)]">
                    Daily reminders to save money
                  </p>
                </div>
              </div>
              <button
                type="button"
                onClick={() => handleToggle('savings_reminders')}
                disabled={!formData.notifications_enabled}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  formData.savings_reminders && formData.notifications_enabled
                    ? 'bg-green-600'
                    : 'bg-gray-300'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    formData.savings_reminders && formData.notifications_enabled
                      ? 'translate-x-6'
                      : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {/* Reminder Time */}
            {formData.savings_reminders && formData.notifications_enabled && (
              <div className="ml-8">
                <label
                  htmlFor="reminder_time"
                  className="block text-sm font-medium text-[var(--sea-ink)] mb-2"
                >
                  Reminder Time
                </label>
                <input
                  type="time"
                  id="reminder_time"
                  value={formData.reminder_time}
                  onChange={(e) => handleInputChange('reminder_time', e.target.value)}
                  className="px-4 py-2 border border-[var(--line)] rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--sea-ink)]"
                />
              </div>
            )}
          </div>
        </div>

        {/* Save Button */}
        <div className="flex justify-end">
          <button
            type="submit"
            disabled={updatePreferencesMutation.isPending}
            className="flex items-center gap-2 px-6 py-2 bg-[var(--sea-ink)] text-white rounded-lg hover:bg-[var(--sea-ink)]/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Save className="w-4 h-4" />
            {updatePreferencesMutation.isPending ? 'Saving...' : 'Save Settings'}
          </button>
        </div>
      </form>

      {/* Danger Zone - Account Deletion */}
      <div className="border border-red-300 rounded-xl p-6 bg-red-50">
        <div className="flex items-center gap-3 mb-4">
          <AlertTriangle className="w-5 h-5 text-red-600" />
          <h2 className="text-xl font-semibold text-red-900">Danger Zone</h2>
        </div>

        <p className="text-sm text-red-800 mb-4">
          Once you delete your account, there is no going back. This will permanently
          delete your profile, savings history, budgets, goals, and all associated data.
        </p>

        {!showDeleteConfirm ? (
          <button
            type="button"
            onClick={() => setShowDeleteConfirm(true)}
            className="flex items-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            <Trash2 className="w-4 h-4" />
            Delete Account
          </button>
        ) : (
          <div className="space-y-4">
            <div>
              <label
                htmlFor="delete_confirm"
                className="block text-sm font-medium text-red-900 mb-2"
              >
                Type <span className="font-bold">DELETE</span> to confirm
              </label>
              <input
                type="text"
                id="delete_confirm"
                value={deleteConfirmText}
                onChange={(e) => setDeleteConfirmText(e.target.value)}
                className="w-full px-4 py-2 border border-red-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-600"
                placeholder="DELETE"
              />
            </div>
            <div className="flex gap-3">
              <button
                type="button"
                onClick={handleDeleteAccount}
                disabled={
                  deleteConfirmText !== 'DELETE' || deleteAccountMutation.isPending
                }
                className="flex items-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Trash2 className="w-4 h-4" />
                {deleteAccountMutation.isPending
                  ? 'Deleting...'
                  : 'Confirm Delete'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowDeleteConfirm(false)
                  setDeleteConfirmText('')
                }}
                className="px-4 py-2 border border-red-300 text-red-900 rounded-lg hover:bg-red-100 transition-colors"
              >
                Cancel
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
