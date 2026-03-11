import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useAuth } from '../lib/auth/context'

export const Route = createFileRoute('/login')({
  component: LoginPage,
})

function LoginPage() {
  const { login } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    try {
      // Requirement 1.6: Login with valid credentials
      await login({ email, password })
      // Success - redirect to dashboard handled by AuthContext
    } catch (err) {
      // Requirement 1.7: Invalid credentials return error
      // Don't reveal whether email or password was incorrect
      setError(err instanceof Error ? err.message : 'Invalid email or password')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-bold display-title">
            Sign in to InSavein
          </h2>
          <p className="mt-2 text-center text-sm text-[var(--sea-ink-soft)]">
            Or{' '}
            <a href="/register" className="font-medium">
              create a new account
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
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="appearance-none relative block w-full px-3 py-2 border border-[var(--line)] rounded-lg placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[var(--lagoon)] focus:border-transparent bg-white/50"
                placeholder="Email address"
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
                autoComplete="current-password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="appearance-none relative block w-full px-3 py-2 border border-[var(--line)] rounded-lg placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[var(--lagoon)] focus:border-transparent bg-white/50"
                placeholder="Password"
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={isLoading}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-[var(--lagoon)] hover:bg-[var(--lagoon-deep)] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-[var(--lagoon)] disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? 'Signing in...' : 'Sign in'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
