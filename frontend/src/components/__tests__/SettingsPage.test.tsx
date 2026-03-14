import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { SettingsPage } from '../SettingsPage'
import type { UserPreferences } from '../../lib/types/api'

/**
 * Tests for SettingsPage component
 * Validates Requirements 3.3, 3.4
 */

const mockPreferences: UserPreferences = {
  currency: 'USD',
  notifications_enabled: true,
  email_notifications: true,
  push_notifications: false,
  savings_reminders: true,
  reminder_time: '09:00',
  theme: 'light',
}

// Mock the useAuth hook
vi.mock('../../lib/auth/context', () => ({
  useAuth: () => ({
    api: {
      user: {
        getPreferences: vi.fn().mockResolvedValue(mockPreferences),
        updatePreferences: vi.fn().mockResolvedValue(mockPreferences),
        deleteAccount: vi.fn().mockResolvedValue(undefined),
      },
    },
    logout: vi.fn(),
  }),
}))

describe('SettingsPage', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
      },
    })
  })

  it('should display loading state initially', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <SettingsPage />
      </QueryClientProvider>
    )

    expect(screen.getByText('Loading settings...')).toBeTruthy()
  })

  it('should render settings page after loading', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <SettingsPage />
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.queryByText('Loading settings...')).toBeNull()
    })

    // Check that main heading is present
    expect(screen.getByText('Settings')).toBeTruthy()
  })

  it('should display notification settings', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <SettingsPage />
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.queryByText('Loading settings...')).toBeNull()
    })

    expect(screen.getByText(/Enable Notifications/)).toBeTruthy()
    expect(screen.getByText(/Email Notifications/)).toBeTruthy()
  })

  it('should display delete account section', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <SettingsPage />
      </QueryClientProvider>
    )

    await waitFor(() => {
      expect(screen.queryByText('Loading settings...')).toBeNull()
    })

    expect(screen.getByText(/Danger Zone/)).toBeTruthy()
    expect(screen.getByText(/Delete Account/)).toBeTruthy()
  })
})
