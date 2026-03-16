import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/context'
import { User, Mail, Calendar, Image as ImageIcon, Save, X } from 'lucide-react'
import type { UserProfile, UpdateProfileRequest } from '../lib/types/api'

// Requirement 3.1: Display user profile information
// Requirement 3.2: Edit profile form with API integration
export function ProfilePage() {
  const { api } = useAuth()
  const queryClient = useQueryClient()
  const [isEditing, setIsEditing] = useState(false)
  const [formData, setFormData] = useState<UpdateProfileRequest>({})
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  // Fetch user profile
  const { data: profile, isLoading } = useQuery<UserProfile>({
    queryKey: ['user', 'profile'],
    queryFn: () => api.user.getProfile(),
  })

  // Update profile mutation
  const updateProfileMutation = useMutation({
    mutationFn: (data: UpdateProfileRequest) => api.user.updateProfile(data),
    onSuccess: (updatedProfile) => {
      queryClient.setQueryData(['user', 'profile'], updatedProfile)
      setIsEditing(false)
      setSuccess('Profile updated successfully!')
      setError(null)
      setTimeout(() => setSuccess(null), 3000)
    },
    onError: (err: any) => {
      setError(err.message || 'Failed to update profile')
      setSuccess(null)
    },
  })

  const handleEdit = () => {
    if (profile) {
      setFormData({
        first_name: profile.first_name,
        last_name: profile.last_name,
        date_of_birth: profile.date_of_birth,
        profile_image_url: profile.profile_image_url,
      })
    }
    setIsEditing(true)
    setError(null)
    setSuccess(null)
  }

  const handleCancel = () => {
    setIsEditing(false)
    setFormData({})
    setError(null)
    setSuccess(null)
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    // Validate required fields
    if (!formData.first_name || !formData.last_name || !formData.date_of_birth) {
      setError('First name, last name, and date of birth are required')
      return
    }

    updateProfileMutation.mutate(formData)
  }

  const handleInputChange = (field: keyof UpdateProfileRequest, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-[var(--sea-ink-soft)]">Loading profile...</div>
      </div>
    )
  }

  if (!profile) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-red-600">Failed to load profile</div>
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-[var(--sea-ink)]">Profile</h1>
        <p className="mt-2 text-[var(--sea-ink-soft)]">
          Manage your personal information
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

      {/* Profile Card */}
      <div className="border border-[var(--line)] rounded-xl p-6 bg-white">
        {!isEditing ? (
          // View Mode
          <div className="space-y-6">
            {/* Profile Image */}
            <div className="flex items-center gap-6">
              <div className="w-24 h-24 rounded-full bg-[var(--link-bg-hover)] flex items-center justify-center overflow-hidden">
                {profile.profile_image_url ? (
                  <img
                    src={profile.profile_image_url}
                    alt="Profile"
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <User className="w-12 h-12 text-[var(--sea-ink-soft)]" />
                )}
              </div>
              <div>
                <h2 className="text-2xl font-bold text-[var(--sea-ink)]">
                  {profile.first_name} {profile.last_name}
                </h2>
                <p className="text-[var(--sea-ink-soft)]">{profile.email}</p>
              </div>
            </div>

            {/* Profile Information */}
            <div className="space-y-4 pt-6 border-t border-[var(--line)]">
              <div className="flex items-center gap-3">
                <User className="w-5 h-5 text-[var(--sea-ink-soft)]" />
                <div>
                  <p className="text-sm text-[var(--sea-ink-soft)]">Full Name</p>
                  <p className="font-medium text-[var(--sea-ink)]">
                    {profile.first_name} {profile.last_name}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <Mail className="w-5 h-5 text-[var(--sea-ink-soft)]" />
                <div>
                  <p className="text-sm text-[var(--sea-ink-soft)]">Email</p>
                  <p className="font-medium text-[var(--sea-ink)]">{profile.email}</p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <Calendar className="w-5 h-5 text-[var(--sea-ink-soft)]" />
                <div>
                  <p className="text-sm text-[var(--sea-ink-soft)]">Date of Birth</p>
                  <p className="font-medium text-[var(--sea-ink)]">
                    {new Date(profile.date_of_birth).toLocaleDateString()}
                  </p>
                </div>
              </div>

              {profile.profile_image_url && (
                <div className="flex items-center gap-3">
                  <ImageIcon className="w-5 h-5 text-[var(--sea-ink-soft)]" />
                  <div>
                    <p className="text-sm text-[var(--sea-ink-soft)]">Profile Image</p>
                    <p className="font-medium text-[var(--sea-ink)] text-sm truncate max-w-md">
                      {profile.profile_image_url}
                    </p>
                  </div>
                </div>
              )}
            </div>

            {/* Edit Button */}
            <div className="pt-6 border-t border-[var(--line)]">
              <button
                onClick={handleEdit}
                className="px-6 py-2 bg-[var(--sea-ink)] text-white rounded-lg hover:bg-[var(--sea-ink)]/90 transition-colors"
              >
                Edit Profile
              </button>
            </div>
          </div>
        ) : (
          // Edit Mode
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* First Name */}
              <div>
                <label
                  htmlFor="first_name"
                  className="block text-sm font-medium text-[var(--sea-ink)] mb-2"
                >
                  First Name *
                </label>
                <input
                  type="text"
                  id="first_name"
                  value={formData.first_name || ''}
                  onChange={(e) => handleInputChange('first_name', e.target.value)}
                  className="w-full px-4 py-2 border border-[var(--line)] rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--sea-ink)]"
                  required
                />
              </div>

              {/* Last Name */}
              <div>
                <label
                  htmlFor="last_name"
                  className="block text-sm font-medium text-[var(--sea-ink)] mb-2"
                >
                  Last Name *
                </label>
                <input
                  type="text"
                  id="last_name"
                  value={formData.last_name || ''}
                  onChange={(e) => handleInputChange('last_name', e.target.value)}
                  className="w-full px-4 py-2 border border-[var(--line)] rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--sea-ink)]"
                  required
                />
              </div>
            </div>

            {/* Date of Birth */}
            <div>
              <label
                htmlFor="date_of_birth"
                className="block text-sm font-medium text-[var(--sea-ink)] mb-2"
              >
                Date of Birth *
              </label>
              <input
                type="date"
                id="date_of_birth"
                value={formData.date_of_birth || ''}
                onChange={(e) => handleInputChange('date_of_birth', e.target.value)}
                className="w-full px-4 py-2 border border-[var(--line)] rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--sea-ink)]"
                required
              />
            </div>

            {/* Profile Image URL */}
            <div>
              <label
                htmlFor="profile_image_url"
                className="block text-sm font-medium text-[var(--sea-ink)] mb-2"
              >
                Profile Image URL
              </label>
              <input
                type="url"
                id="profile_image_url"
                value={formData.profile_image_url || ''}
                onChange={(e) => handleInputChange('profile_image_url', e.target.value)}
                placeholder="https://example.com/image.jpg"
                className="w-full px-4 py-2 border border-[var(--line)] rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--sea-ink)]"
              />
            </div>

            {/* Action Buttons */}
            <div className="flex gap-3 pt-6 border-t border-[var(--line)]">
              <button
                type="submit"
                disabled={updateProfileMutation.isPending}
                className="flex items-center gap-2 px-6 py-2 bg-[var(--sea-ink)] text-white rounded-lg hover:bg-[var(--sea-ink)]/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Save className="w-4 h-4" />
                {updateProfileMutation.isPending ? 'Saving...' : 'Save Changes'}
              </button>
              <button
                type="button"
                onClick={handleCancel}
                className="flex items-center gap-2 px-6 py-2 border border-[var(--line)] text-[var(--sea-ink)] rounded-lg hover:bg-[var(--link-bg-hover)] transition-colors"
              >
                <X className="w-4 h-4" />
                Cancel
              </button>
            </div>
          </form>
        )}
      </div>

      {/* Account Information */}
      <div className="border border-[var(--line)] rounded-xl p-6 bg-white">
        <h3 className="text-lg font-semibold text-[var(--sea-ink)] mb-4">
          Account Information
        </h3>
        <div className="space-y-3 text-sm">
          <div className="flex justify-between">
            <span className="text-[var(--sea-ink-soft)]">Account Created</span>
            <span className="font-medium text-[var(--sea-ink)]">
              {new Date(profile.created_at).toLocaleDateString()}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-[var(--sea-ink-soft)]">Last Updated</span>
            <span className="font-medium text-[var(--sea-ink)]">
              {new Date(profile.updated_at).toLocaleDateString()}
            </span>
          </div>
        </div>
      </div>
    </div>
  )
}
