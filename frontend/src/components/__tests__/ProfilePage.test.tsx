import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ProfilePage } from '../ProfilePage'
import type { UserProfile } from '../../lib/types/api'

/**
 * Tests for ProfilePage component
 * Validates Requirements 3.1, 3.2
 */

// Mock the useAuth hook
vi.mock('../../lib/auth/context', () => ({
  useAuth: () => ({
    api: {
      user: {
        getProfile: vi.fn().mockResolvedValue({
          id: '123',
          email: 'test@example.com',
          first_name: 'John',
          last_name: 'Doe',
          date_of_birth: '1990-01-01',
          profile_image_url: 'https://example.com/image.jpg',
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        }),
        updateProfile: vi.fn(),
      },
    },
  }),
}))

describe('ProfilePage', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
      },
    })
  })

  it('should render profile information correctly', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <ProfilePage />
      </QueryClientProvider>
    )

    // Wait for profile to load
    await waitFor(() => {
      expect(screen.getAllByText(/John/)).toBeTruthy()
    })

    expect(screen.getAllByText(/test@example.com/)).toBeTruthy()
    expect(screen.getByText('Full Name')).toBeTruthy()
    expect(screen.getByText('Email')).toBeTruthy()
    expect(screen.getByText('Date of Birth')).toBeTruthy()
  })

  it('should display loading state initially', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <ProfilePage />
      </QueryClientProvider>
    )

    expect(screen.getByText('Loading profile...')).toBeTruthy()
  })

  it('should display edit button in view mode', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <ProfilePage />
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.getByText('Edit Profile')).toBeTruthy()
    })
  })

  it('should display account information section', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <ProfilePage />
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.getByText('Account Information')).toBeTruthy()
    })

    expect(screen.getByText('Account Created')).toBeTruthy()
    expect(screen.getByText('Last Updated')).toBeTruthy()
  })
})
