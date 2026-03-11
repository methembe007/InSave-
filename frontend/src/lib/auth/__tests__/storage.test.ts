import { describe, it, expect, beforeEach } from 'vitest'
import { tokenStorage } from '../storage'

/**
 * Tests for token storage utilities
 * Validates Requirements 2.4 (token management)
 */

describe('Token Storage', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear()
  })

  describe('Access Token', () => {
    it('should store and retrieve access token', () => {
      const token = 'test-access-token'
      tokenStorage.setAccessToken(token)
      expect(tokenStorage.getAccessToken()).toBe(token)
    })

    it('should return null when no access token is stored', () => {
      expect(tokenStorage.getAccessToken()).toBeNull()
    })
  })

  describe('Refresh Token', () => {
    it('should store and retrieve refresh token', () => {
      const token = 'test-refresh-token'
      tokenStorage.setRefreshToken(token)
      expect(tokenStorage.getRefreshToken()).toBe(token)
    })

    it('should return null when no refresh token is stored', () => {
      expect(tokenStorage.getRefreshToken()).toBeNull()
    })
  })

  describe('Token Expiry', () => {
    it('should store and retrieve token expiry', () => {
      const expiresIn = 900 // 15 minutes
      tokenStorage.setTokenExpiry(expiresIn)
      const expiry = tokenStorage.getTokenExpiry()
      expect(expiry).toBeGreaterThan(Date.now())
    })

    it('should return null when no expiry is stored', () => {
      expect(tokenStorage.getTokenExpiry()).toBeNull()
    })

    it('should detect expired tokens', () => {
      // Set expiry to 1 second ago
      const pastExpiry = Date.now() - 1000
      localStorage.setItem('insavein_token_expiry', pastExpiry.toString())
      expect(tokenStorage.isTokenExpired()).toBe(true)
    })

    it('should detect valid tokens', () => {
      // Set expiry to 10 minutes from now
      const expiresIn = 600
      tokenStorage.setTokenExpiry(expiresIn)
      expect(tokenStorage.isTokenExpired()).toBe(false)
    })

    it('should consider token expired 1 minute before actual expiry', () => {
      // Set expiry to 30 seconds from now (less than 1 minute buffer)
      const futureExpiry = Date.now() + 30000
      localStorage.setItem('insavein_token_expiry', futureExpiry.toString())
      expect(tokenStorage.isTokenExpired()).toBe(true)
    })
  })

  describe('Clear Tokens', () => {
    it('should clear all tokens', () => {
      tokenStorage.setAccessToken('access-token')
      tokenStorage.setRefreshToken('refresh-token')
      tokenStorage.setTokenExpiry(900)

      tokenStorage.clearTokens()

      expect(tokenStorage.getAccessToken()).toBeNull()
      expect(tokenStorage.getRefreshToken()).toBeNull()
      expect(tokenStorage.getTokenExpiry()).toBeNull()
    })
  })
})
