import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { NotificationsDropdown } from '../NotificationsDropdown'
import type { Notification } from '../../lib/types/api'

// Mock the auth context
const mockApi = {
  notifications: {
    getNotifications: vi.fn(),
    markAsRead: vi.fn(),
  },
}

vi.mock('../../lib/auth/context', () => ({
  useAuth: () => ({
    api: mockApi,
  }),
}))

describe('NotificationsDropdown', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders bell icon', () => {
    render(<NotificationsDropdown />)
    
    const bellButton = screen.getByRole('button', { name: /notifications/i })
    expect(bellButton).toBeTruthy()
  })

  it('fetches and displays notifications when opened', async () => {
    const mockNotifications: Notification[] = [
      {
        id: '1',
        user_id: 'user1',
        type: 'savings',
        title: 'Savings Goal Reached',
        message: 'Congratulations! You reached your savings goal.',
        is_read: false,
        created_at: new Date().toISOString(),
      },
      {
        id: '2',
        user_id: 'user1',
        type: 'budget',
        title: 'Budget Alert',
        message: 'You have exceeded 80% of your budget.',
        is_read: true,
        created_at: new Date(Date.now() - 86400000).toISOString(), // 1 day ago
      },
    ]

    mockApi.notifications.getNotifications.mockResolvedValue(mockNotifications)

    render(<NotificationsDropdown />)

    // Open dropdown
    const bellButton = screen.getByRole('button', { name: /notifications/i })
    fireEvent.click(bellButton)

    // Wait for notifications to load
    await waitFor(() => {
      expect(screen.getByText('Savings Goal Reached')).toBeTruthy()
      expect(screen.getByText('Budget Alert')).toBeTruthy()
    })

    expect(mockApi.notifications.getNotifications).toHaveBeenCalledTimes(1)
  })

  it('displays unread count badge', async () => {
    const mockNotifications: Notification[] = [
      {
        id: '1',
        user_id: 'user1',
        type: 'savings',
        title: 'Notification 1',
        message: 'Message 1',
        is_read: false,
        created_at: new Date().toISOString(),
      },
      {
        id: '2',
        user_id: 'user1',
        type: 'budget',
        title: 'Notification 2',
        message: 'Message 2',
        is_read: false,
        created_at: new Date().toISOString(),
      },
      {
        id: '3',
        user_id: 'user1',
        type: 'goal',
        title: 'Notification 3',
        message: 'Message 3',
        is_read: true,
        created_at: new Date().toISOString(),
      },
    ]

    mockApi.notifications.getNotifications.mockResolvedValue(mockNotifications)

    render(<NotificationsDropdown />)

    // Open dropdown to trigger fetch
    const bellButton = screen.getByRole('button', { name: /notifications/i })
    fireEvent.click(bellButton)

    // Wait for notifications to load
    await waitFor(() => {
      expect(screen.getByText('Notification 1')).toBeTruthy()
    })

    // Check unread count badge shows 2
    expect(screen.getByText('2')).toBeTruthy()
  })

  it('marks notification as read when clicked', async () => {
    const mockNotifications: Notification[] = [
      {
        id: '1',
        user_id: 'user1',
        type: 'savings',
        title: 'Unread Notification',
        message: 'This is unread',
        is_read: false,
        created_at: new Date().toISOString(),
      },
    ]

    mockApi.notifications.getNotifications.mockResolvedValue(mockNotifications)
    mockApi.notifications.markAsRead.mockResolvedValue(undefined)

    render(<NotificationsDropdown />)

    // Open dropdown
    const bellButton = screen.getByRole('button', { name: /notifications/i })
    fireEvent.click(bellButton)

    // Wait for notification to appear
    await waitFor(() => {
      expect(screen.getByText('Unread Notification')).toBeTruthy()
    })

    // Click on the notification
    const notification = screen.getByText('Unread Notification')
    fireEvent.click(notification)

    // Verify markAsRead was called
    expect(mockApi.notifications.markAsRead).toHaveBeenCalledWith('1')
  })

  it('displays empty state when no notifications', async () => {
    mockApi.notifications.getNotifications.mockResolvedValue([])

    render(<NotificationsDropdown />)

    // Open dropdown
    const bellButton = screen.getByRole('button', { name: /notifications/i })
    fireEvent.click(bellButton)

    // Wait for empty state
    await waitFor(() => {
      expect(screen.getByText('No notifications yet')).toBeTruthy()
    })
  })

  it('displays error message when fetch fails', async () => {
    mockApi.notifications.getNotifications.mockRejectedValue(
      new Error('Failed to fetch notifications')
    )

    render(<NotificationsDropdown />)

    // Open dropdown
    const bellButton = screen.getByRole('button', { name: /notifications/i })
    fireEvent.click(bellButton)

    // Wait for error message
    await waitFor(() => {
      expect(screen.getByText('Failed to fetch notifications')).toBeTruthy()
    })
  })

  it('orders notifications by date descending', async () => {
    const mockNotifications: Notification[] = [
      {
        id: '1',
        user_id: 'user1',
        type: 'savings',
        title: 'Oldest',
        message: 'Message 1',
        is_read: false,
        created_at: new Date('2024-01-01').toISOString(),
      },
      {
        id: '2',
        user_id: 'user1',
        type: 'budget',
        title: 'Newest',
        message: 'Message 2',
        is_read: false,
        created_at: new Date('2024-01-03').toISOString(),
      },
      {
        id: '3',
        user_id: 'user1',
        type: 'goal',
        title: 'Middle',
        message: 'Message 3',
        is_read: true,
        created_at: new Date('2024-01-02').toISOString(),
      },
    ]

    mockApi.notifications.getNotifications.mockResolvedValue(mockNotifications)

    render(<NotificationsDropdown />)

    // Open dropdown
    const bellButton = screen.getByRole('button', { name: /notifications/i })
    fireEvent.click(bellButton)

    // Wait for notifications to load
    await waitFor(() => {
      expect(screen.getByText('Newest')).toBeTruthy()
    })

    // Get all notification titles
    const titles = screen.getAllByRole('heading', { level: 4 })
    
    // Verify order: Newest, Middle, Oldest
    expect(titles[0].textContent).toBe('Newest')
    expect(titles[1].textContent).toBe('Middle')
    expect(titles[2].textContent).toBe('Oldest')
  })
})
