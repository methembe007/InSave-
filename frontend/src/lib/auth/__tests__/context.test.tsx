import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { AuthProvider, useAuth } from '../context'
import { tokenStorage } from '../storage'
import type { AuthResponse, LoginRequest, RegisterRequest } from '../../types/api'

// Mock the router
vi.mock('@tanstack/react-router', () => ({
  useNavigate: () => vi.fn(),
}))

// Mock the API services
vi.mock('../../api', () => ({
  createApiServices: () => ({
    auth: {
      login: vi.fn(),
      register: vi.fn(),
      logout: vi.fn(),
      refreshToken: vi.fn(),
    },
    user: {
      getProfile: vi.fn(),
    },
  }),
}))

// Test component that uses auth
function TestComponent() {
  const { user, isAuthenticated, isLoading, login, register, logout } = useAuth()

  return (
    <div>
      <div data-testid="loading">{isLoading ? 'Loading' : 'Ready'}</div>
      <div data-testid="authenticated">{isAuthenticated ? 'Yes' : 'No'}</div>
      <div data-testid="user">{user ? `${user.first_name} ${user.last_name}` : 'None'}</div>
      <button onClick={() => login({ email: 'test@example.com', password: 'password123' })}>
        Login
      </button>
      <button onClick={() => register({
        email: 'new@example.com',
        password: 'password123',
        first_name: 'John',
        last_name: 'Doe',
        date_of_birth: '1990-01-01',
      })}>
        Register
      </button>
      <button onClick={() => logout()}>Logout</button>
    </div>
  )
}

describe('AuthContext', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    tokenStorage.clearTokens()
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('should provide initial unauthenticated state', async () => {
    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('Ready')
    })

    expect(screen.getByTestId('authenticated')).toHaveTextContent('No')
    expect(screen.getByTestId('user')).toHaveTextContent('None')
  })

  it('should handle login successfully', async () => {
    const mockAuthResponse: AuthResponse = {
      access_token: 'mock-access-token',
      refresh_token: 'mock-refresh-token',
      expires_in: 900,
      user: {
        id: '123',
        email: 'test@example.com',
        first_name: 'Test',
        last_name: 'User',
      },
    }

    const { createApiServices } = await import('../../api')
    const mockApi = createApiServices(() => null, () => {})
    vi.mocked(mockApi.auth.login).mockResolvedValue(mockAuthResponse)

    const user = userEvent.setup()
    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('Ready')
    })

    const loginButton = screen.getByText('Login')
    await user.click(loginButton)

    await waitFor(() => {
      expect(tokenStorage.getAccessToken()).toBe('mock-access-token')
      expect(tokenStorage.getRefreshToken()).toBe('mock-refresh-token')
    })
  })

  it('should handle registration successfully', async () => {
    const mockAuthResponse: AuthResponse = {
      access_token: 'mock-access-token',
      refresh_token: 'mock-refresh-token',
      expires_in: 900,
      user: {
        id: '456',
        email: 'new@example.com',
        first_name: 'John',
        last_name: 'Doe',
      },
    }

    const { createApiServices } = await import('../../api')
    const mockApi = createApiServices(() => null, () => {})
    vi.mocked(mockApi.auth.register).mockResolvedValue(mockAuthResponse)

    const user = userEvent.setup()
    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('Ready')
    })

    const registerButton = screen.getByText('Register')
    await user.click(registerButton)

    await waitFor(() => {
      expect(tokenStorage.getAccessToken()).toBe('mock-access-token')
      expect(tokenStorage.getRefreshToken()).toBe('mock-refresh-token')
    })
  })

  it('should handle logout successfully', async () => {
    // Set up initial authenticated state
    tokenStorage.setAccessToken('existing-token')
    tokenStorage.setRefreshToken('existing-refresh')

    const { createApiServices } = await import('../../api')
    const mockApi = createApiServices(() => null, () => {})
    vi.mocked(mockApi.auth.logout).mockResolvedValue(undefined)

    const user = userEvent.setup()
    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('Ready')
    })

    const logoutButton = screen.getByText('Logout')
    await user.click(logoutButton)

    await waitFor(() => {
      expect(tokenStorage.getAccessToken()).toBeNull()
      expect(tokenStorage.getRefreshToken()).toBeNull()
    })
  })

  it('should handle login failure', async () => {
    const { createApiServices } = await import('../../api')
    const mockApi = createApiServices(() => null, () => {})
    vi.mocked(mockApi.auth.login).mockRejectedValue(new Error('Invalid credentials'))

    const user = userEvent.setup()
    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('Ready')
    })

    const loginButton = screen.getByText('Login')
    
    // Should throw error
    await expect(async () => {
      await user.click(loginButton)
    }).rejects.toThrow()

    // Should remain unauthenticated
    expect(screen.getByTestId('authenticated')).toHaveTextContent('No')
  })

  it('should restore session from stored tokens', async () => {
    // Set up stored tokens
    tokenStorage.setAccessToken('stored-token')
    tokenStorage.setRefreshToken('stored-refresh')
    tokenStorage.setTokenExpiry(900)

    const mockProfile = {
      id: '789',
      email: 'stored@example.com',
      first_name: 'Stored',
      last_name: 'User',
      date_of_birth: '1990-01-01',
      profile_image_url: '',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    const { createApiServices } = await import('../../api')
    const mockApi = createApiServices(() => 'stored-token', () => {})
    vi.mocked(mockApi.user.getProfile).mockResolvedValue(mockProfile)

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('Ready')
    })

    await waitFor(() => {
      expect(screen.getByTestId('authenticated')).toHaveTextContent('Yes')
      expect(screen.getByTestId('user')).toHaveTextContent('Stored User')
    })
  })

  it('should clear tokens on auth check failure', async () => {
    // Set up stored tokens
    tokenStorage.setAccessToken('invalid-token')
    tokenStorage.setRefreshToken('invalid-refresh')

    const { createApiServices } = await import('../../api')
    const mockApi = createApiServices(() => 'invalid-token', () => {})
    vi.mocked(mockApi.user.getProfile).mockRejectedValue(new Error('Unauthorized'))

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('Ready')
    })

    await waitFor(() => {
      expect(tokenStorage.getAccessToken()).toBeNull()
      expect(tokenStorage.getRefreshToken()).toBeNull()
      expect(screen.getByTestId('authenticated')).toHaveTextContent('No')
    })
  })
})
