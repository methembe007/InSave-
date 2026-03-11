import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { createApiServices, type ApiServices } from '../api'
import { tokenStorage } from './storage'
import type { AuthResponse, LoginRequest, RegisterRequest, UserSummary } from '../types/api'

interface AuthContextType {
  user: UserSummary | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (data: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => Promise<void>
  refreshToken: () => Promise<void>
  api: ApiServices
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserSummary | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const navigate = useNavigate()

  // Token refresh in progress flag
  const [isRefreshing, setIsRefreshing] = useState(false)

  const handleUnauthorized = useCallback(() => {
    tokenStorage.clearTokens()
    setUser(null)
    navigate({ to: '/login' })
  }, [navigate])

  const getToken = useCallback(() => {
    return tokenStorage.getAccessToken()
  }, [])

  // Create API services
  const api = createApiServices(getToken, handleUnauthorized)

  const refreshToken = useCallback(async () => {
    if (isRefreshing) return
    
    const refreshTokenValue = tokenStorage.getRefreshToken()
    if (!refreshTokenValue) {
      handleUnauthorized()
      return
    }

    try {
      setIsRefreshing(true)
      const response = await api.auth.refreshToken(refreshTokenValue)
      
      tokenStorage.setAccessToken(response.access_token)
      tokenStorage.setRefreshToken(response.refresh_token)
      tokenStorage.setTokenExpiry(response.expires_in)
    } catch (error) {
      console.error('Token refresh failed:', error)
      handleUnauthorized()
    } finally {
      setIsRefreshing(false)
    }
  }, [api.auth, handleUnauthorized, isRefreshing])

  const login = useCallback(async (data: LoginRequest) => {
    try {
      const response: AuthResponse = await api.auth.login(data)
      
      tokenStorage.setAccessToken(response.access_token)
      tokenStorage.setRefreshToken(response.refresh_token)
      tokenStorage.setTokenExpiry(response.expires_in)
      
      setUser(response.user)
      navigate({ to: '/dashboard' })
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }, [api.auth, navigate])

  const register = useCallback(async (data: RegisterRequest) => {
    try {
      const response: AuthResponse = await api.auth.register(data)
      
      tokenStorage.setAccessToken(response.access_token)
      tokenStorage.setRefreshToken(response.refresh_token)
      tokenStorage.setTokenExpiry(response.expires_in)
      
      setUser(response.user)
      navigate({ to: '/dashboard' })
    } catch (error) {
      console.error('Registration failed:', error)
      throw error
    }
  }, [api.auth, navigate])

  const logout = useCallback(async () => {
    try {
      await api.auth.logout()
    } catch (error) {
      console.error('Logout failed:', error)
    } finally {
      tokenStorage.clearTokens()
      setUser(null)
      navigate({ to: '/' })
    }
  }, [api.auth, navigate])

  // Check authentication status on mount
  useEffect(() => {
    const checkAuth = async () => {
      const token = tokenStorage.getAccessToken()
      
      if (!token) {
        setIsLoading(false)
        return
      }

      // Check if token is expired
      if (tokenStorage.isTokenExpired()) {
        await refreshToken()
      }

      // Fetch user profile to verify token
      try {
        const profile = await api.user.getProfile()
        setUser({
          id: profile.id,
          email: profile.email,
          first_name: profile.first_name,
          last_name: profile.last_name,
        })
      } catch (error) {
        console.error('Auth check failed:', error)
        tokenStorage.clearTokens()
      } finally {
        setIsLoading(false)
      }
    }

    checkAuth()
  }, [api.user, refreshToken])

  // Set up automatic token refresh
  useEffect(() => {
    const interval = setInterval(() => {
      if (tokenStorage.isTokenExpired() && tokenStorage.getRefreshToken()) {
        refreshToken()
      }
    }, 60000) // Check every minute

    return () => clearInterval(interval)
  }, [refreshToken])

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    register,
    logout,
    refreshToken,
    api,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
