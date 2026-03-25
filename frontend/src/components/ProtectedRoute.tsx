import { useEffect, type ReactNode } from 'react'
import { useNavigate, useLocation } from '@tanstack/react-router'
import { useAuth } from '../lib/auth/context'

interface ProtectedRouteProps {
  children: ReactNode
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  useEffect(() => {
    // Only redirect if not loading and not authenticated
    // Also check we're not already on the login page to avoid loops
    if (!isLoading && !isAuthenticated && location.pathname !== '/login') {
      navigate({ to: '/login' })
    }
  }, [isAuthenticated, isLoading, navigate, location.pathname])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite]" />
          <p className="mt-4 text-lg">Loading...</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return null
  }

  return <>{children}</>
}
