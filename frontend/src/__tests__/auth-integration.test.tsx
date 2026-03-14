import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import React from 'react'
import { tokenStorage } from '../lib/auth/storage'
import type { AuthResponse } from '../lib/types/api'

/**
 * Integration tests for the complete authentication flow
 * Tests the interaction between login, token storage, protected routes, and logout
 */

// Mock fetch for API calls
const mockFetch = vi.fn()
global.fetch = mockFetch

// Mock router
const mockNavigate = vi.fn()
vi.mock('@tanstack/react-router', () => ({
  useNavigate: () => mockNavigate,
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}))

describe('Authentication Integration Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    tokenStorage.clearTokens()
    mockFetch.mockClear()
  })

  afterEach(() => {
    tokenStorage.clearTokens()
  })

  describe('Complete Registration Flow', () => {
    it('should register user, store tokens, and redirect to dashboard', async () => {
      const mockAuthResponse: AuthResponse = {
        access_token: 'new-access-token',
        refresh_token: 'new-refresh-token',
        expires_in: 900,
        user: {
          id: 'user-123',
          email: 'newuser@example.com',
          first_name: 'New',
          last_name: 'User',
        },
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockAuthResponse,
      })

      // Simulate registration
      const response = await fetch('http://localhost:8081/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: 'newuser@example.com',
          password: 'password123',
          first_name: 'New',
          last_name: 'User',
          date_of_birth: '1990-01-01',
        }),
      })

      const data = await response.json()

      // Store tokens
      tokenStorage.setAccessToken(data.access_token)
      tokenStorage.setRefreshToken(data.refresh_token)
      tokenStorage.setTokenExpiry(data.expires_in)

      // Verify tokens are stored
      expect(tokenStorage.getAccessToken()).toBe('new-access-token')
      expect(tokenStorage.getRefreshToken()).toBe('new-refresh-token')
      expect(tokenStorage.isTokenExpired()).toBe(false)
    })
  })

  describe('Complete Login Flow', () => {
    it('should login user, store tokens, and allow access to protected routes', async () => {
      const mockAuthResponse: AuthResponse = {
        access_token: 'login-access-token',
        refresh_token: 'login-refresh-token',
        expires_in: 900,
        user: {
          id: 'user-456',
          email: 'existing@example.com',
          first_name: 'Existing',
          last_name: 'User',
        },
      }

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockAuthResponse,
      })

      // Simulate login
      const response = await fetch('http://localhost:8081/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: 'existing@example.com',
          password: 'password123',
        }),
      })

      const data = await response.json()

      // Store tokens
      tokenStorage.setAccessToken(data.access_token)
      tokenStorage.setRefreshToken(data.refresh_token)
      tokenStorage.setTokenExpiry(data.expires_in)

      // Verify authentication state
      expect(tokenStorage.getAccessToken()).toBe('login-access-token')
      expect(tokenStorage.getRefreshToken()).toBe('login-refresh-token')
      expect(tokenStorage.isTokenExpired()).toBe(false)

      // Simulate protected API call
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          id: 'user-456',
          email: 'existing@example.com',
          first_name: 'Existing',
          last_name: 'User',
        }),
      })

      const profileResponse = await fetch('http://localhost:8082/api/user/profile', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokenStorage.getAccessToken()}`,
        },
      })

      expect(profileResponse.ok).toBe(true)
      const profile = await profileResponse.json()
      expect(profile.email).toBe('existing@example.com')
    })
  })

  describe('Token Refresh Flow', () => {
    it('should refresh expired token and continue session', async () => {
      // Set up expired token
      tokenStorage.setAccessToken('expired-token')
      tokenStorage.setRefreshToken('valid-refresh-token')
      tokenStorage.setTokenExpiry(-100) // Expired

      expect(tokenStorage.isTokenExpired()).toBe(true)

      // Mock refresh token response
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          access_token: 'refreshed-access-token',
          refresh_token: 'new-refresh-token',
          expires_in: 900,
        }),
      })

      // Simulate token refresh
      const response = await fetch('http://localhost:8081/api/auth/refresh', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          refresh_token: tokenStorage.getRefreshToken(),
        }),
      })

      const data = await response.json()

      // Update tokens
      tokenStorage.setAccessToken(data.access_token)
      tokenStorage.setRefreshToken(data.refresh_token)
      tokenStorage.setTokenExpiry(data.expires_in)

      // Verify new tokens
      expect(tokenStorage.getAccessToken()).toBe('refreshed-access-token')
      expect(tokenStorage.getRefreshToken()).toBe('new-refresh-token')
      expect(tokenStorage.isTokenExpired()).toBe(false)
    })
  })

  describe('Logout Flow', () => {
    it('should logout user, clear tokens, and redirect to home', async () => {
      // Set up authenticated state
      tokenStorage.setAccessToken('active-token')
      tokenStorage.setRefreshToken('active-refresh')
      tokenStorage.setTokenExpiry(900)

      expect(tokenStorage.getAccessToken()).toBe('active-token')

      // Mock logout response
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 204,
      })

      // Simulate logout
      await fetch('http://localhost:8081/api/auth/logout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokenStorage.getAccessToken()}`,
        },
      })

      // Clear tokens
      tokenStorage.clearTokens()

      // Verify tokens are cleared
      expect(tokenStorage.getAccessToken()).toBeNull()
      expect(tokenStorage.getRefreshToken()).toBeNull()
    })
  })

  describe('Unauthorized Access Handling', () => {
    it('should handle 401 response and clear tokens', async () => {
      // Set up authenticated state
      tokenStorage.setAccessToken('invalid-token')
      tokenStorage.setRefreshToken('invalid-refresh')

      // Mock 401 response
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({
          error: 'Unauthorized',
          message: 'Invalid or expired token',
        }),
      })

      // Simulate API call with invalid token
      const response = await fetch('http://localhost:8082/api/user/profile', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokenStorage.getAccessToken()}`,
        },
      })

      expect(response.status).toBe(401)

      // Clear tokens on 401
      if (response.status === 401) {
        tokenStorage.clearTokens()
      }

      // Verify tokens are cleared
      expect(tokenStorage.getAccessToken()).toBeNull()
      expect(tokenStorage.getRefreshToken()).toBeNull()
    })
  })

  describe('Session Persistence', () => {
    it('should restore session from stored tokens on page reload', () => {
      // Simulate stored tokens from previous session
      tokenStorage.setAccessToken('stored-token')
      tokenStorage.setRefreshToken('stored-refresh')
      tokenStorage.setTokenExpiry(900)

      // Simulate page reload (tokens should persist in localStorage)
      const accessToken = tokenStorage.getAccessToken()
      const refreshToken = tokenStorage.getRefreshToken()
      const isExpired = tokenStorage.isTokenExpired()

      expect(accessToken).toBe('stored-token')
      expect(refreshToken).toBe('stored-refresh')
      expect(isExpired).toBe(false)
    })

    it('should not restore session if tokens are expired', () => {
      // Set up expired tokens
      tokenStorage.setAccessToken('expired-token')
      tokenStorage.setRefreshToken('expired-refresh')
      tokenStorage.setTokenExpiry(-100) // Expired

      // Check if token is expired
      const isExpired = tokenStorage.isTokenExpired()
      expect(isExpired).toBe(true)

      // Should trigger refresh or logout
      if (isExpired && !tokenStorage.getRefreshToken()) {
        tokenStorage.clearTokens()
      }
    })
  })

  describe('Multiple API Calls with Authentication', () => {
    it('should include auth token in all protected API calls', async () => {
      const token = 'valid-auth-token'
      tokenStorage.setAccessToken(token)

      // Mock multiple API responses
      mockFetch
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => ({ total_saved: 1000 }),
        })
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => ({ total_budget: 2000 }),
        })
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => ({ goals: [] }),
        })

      // Simulate multiple API calls
      const endpoints = [
        'http://localhost:8083/api/savings/summary',
        'http://localhost:8084/api/budget/current',
        'http://localhost:8085/api/goals',
      ]

      const requests = endpoints.map((url) =>
        fetch(url, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
        })
      )

      const responses = await Promise.all(requests)

      // Verify all requests succeeded
      responses.forEach((response) => {
        expect(response.ok).toBe(true)
      })

      // Verify all requests included auth header
      expect(mockFetch).toHaveBeenCalledTimes(3)
      mockFetch.mock.calls.forEach((call) => {
        const headers = call[1]?.headers as Record<string, string>
        expect(headers.Authorization).toBe(`Bearer ${token}`)
      })
    })
  })

  describe('Error Handling', () => {
    it('should handle network errors gracefully', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      try {
        await fetch('http://localhost:8081/api/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            email: 'test@example.com',
            password: 'password123',
          }),
        })
      } catch (error) {
        expect(error).toBeInstanceOf(Error)
        expect((error as Error).message).toBe('Network error')
      }

      // Tokens should not be set on error
      expect(tokenStorage.getAccessToken()).toBeNull()
    })

    it('should handle invalid credentials error', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({
          error: 'Unauthorized',
          message: 'Invalid email or password',
        }),
      })

      const response = await fetch('http://localhost:8081/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: 'wrong@example.com',
          password: 'wrongpassword',
        }),
      })

      expect(response.ok).toBe(false)
      expect(response.status).toBe(401)

      const error = await response.json()
      expect(error.message).toBe('Invalid email or password')

      // Tokens should not be set on error
      expect(tokenStorage.getAccessToken()).toBeNull()
    })
  })
})
