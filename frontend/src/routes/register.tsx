import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { useAuth } from '../lib/auth/context'

export const Route = createFileRoute('/register')({
  component: RegisterPage,
})

function RegisterPage() {
  const { register, isAuthenticated } = useAuth()
  const navigate = useNavigate()
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    date_of_birth: '',
  })
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  // Navigate to dashboard when authenticated
  useEffect(() => {
    if (isAuthenticated) {
      navigate({ to: '/dashboard' })
    }
  }, [isAuthenticated, navigate])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    })
  }

  const validateEmail = (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    // Client-side validation - email format (Requirement 17.1)
    if (!validateEmail(formData.email)) {
      setError('Please enter a valid email address')
      return
    }

    // Client-side validation - password length (Requirement 1.4)
    if (formData.password.length < 8) {
      setError('Password must be at least 8 characters long')
      return
    }

    // Required field validation (Requirement 17.2)
    if (!formData.first_name || !formData.last_name || !formData.date_of_birth) {
      setError('All fields are required')
      return
    }

    setIsLoading(true)

    try {
      await register(formData)
      // Navigation handled by useEffect above (Requirement 1.1)
    } catch (err) {
      // Handle error states with detailed messages (Requirement 17.1)
      if (err instanceof Error) {
        if (err.message.includes('duplicate') || err.message.includes('already exists')) {
          setError('An account with this email already exists')
        } else {
          setError(err.message)
        }
      } else {
        setError('Registration failed. Please try again.')
      }
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-bold display-title">
            Create your account
          </h2>
          <p className="mt-2 text-center text-sm text-[var(--sea-ink-soft)]">
            Or{' '}
            <a href="/login" className="font-medium">
              sign in to existing account
            </a>
          </p>
        </div>
        <form className="mt-8 space-y-6 island-shell rounded-2xl p-8" onSubmit={handleSubmit}>
          {error && (
            <div className="rounded-md bg-red-50 p-4 border border-red-200">
              <p className="text-sm text-red-800">{error}</p>
            </div>
          )}
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="first_name" className="block text-sm font-medium mb-2">
                  First Name
                </label>
                <input
                  id="first_name"
                  name="first_name"
                  type="text"
                  required
                  value={formData.first_name}
                  onChange={handleChange}
                  className="appearance-none relative block w-full px-3 py-2 border border-[var(--line)] rounded-lg placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[var(--lagoon)] focus:border-transparent bg-white/50"
                  placeholder="First name"
                />
              </div>
              <div>
                <label htmlFor="last_name" className="block text-sm font-medium mb-2">
                  Last Name
                </label>
                <input
                  id="last_name"
                  name="last_name"
                  type="text"
                  required
                  value={formData.last_name}
                  onChange={handleChange}
                  className="appearance-none relative block w-full px-3 py-2 border border-[var(--line)] rounded-lg placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[var(--lagoon)] focus:border-transparent bg-white/50"
                  placeholder="Last name"
                />
              </div>
            </div>
            <div>
              <label htmlFor="email" className="block text-sm font-medium mb-2">
                Email address
              </label>
              <input
                id="email"
                name="email"
                type="email"
                autoComplete="email"
                required
                value={formData.email}
                onChange={handleChange}
                className="appearance-none relative block w-full px-3 py-2 border border-[var(--line)] rounded-lg placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[var(--lagoon)] focus:border-transparent bg-white/50"
                placeholder="Email address"
              />
            </div>
            <div>
              <label htmlFor="date_of_birth" className="block text-sm font-medium mb-2">
                Date of Birth
              </label>
              <input
                id="date_of_birth"
                name="date_of_birth"
                type="date"
                required
                value={formData.date_of_birth}
                onChange={handleChange}
                className="appearance-none relative block w-full px-3 py-2 border border-[var(--line)] rounded-lg placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[var(--lagoon)] focus:border-transparent bg-white/50"
              />
            </div>
            <div>
              <label htmlFor="password" className="block text-sm font-medium mb-2">
                Password
              </label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="new-password"
                required
                value={formData.password}
                onChange={handleChange}
                className="appearance-none relative block w-full px-3 py-2 border border-[var(--line)] rounded-lg placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[var(--lagoon)] focus:border-transparent bg-white/50"
                placeholder="Password (min. 8 characters)"
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={isLoading}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-[var(--lagoon)] hover:bg-[var(--lagoon-deep)] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-[var(--lagoon)] disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? 'Creating account...' : 'Create account'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
