import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import React from 'react'

// Mock the router
const mockNavigate = vi.fn()
vi.mock('@tanstack/react-router', () => ({
  createFileRoute: () => ({
    component: null,
  }),
  useNavigate: () => mockNavigate,
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}))

// Mock the auth context
const mockRegister = vi.fn()
const mockUseAuth = vi.fn(() => ({
  register: mockRegister,
  isLoading: false,
}))

vi.mock('../../lib/auth/context', () => ({
  useAuth: () => mockUseAuth(),
}))

// Register component for testing
function RegisterPage() {
  const { register, isLoading } = mockUseAuth()
  const [formData, setFormData] = React.useState({
    email: '',
    password: '',
    confirmPassword: '',
    first_name: '',
    last_name: '',
    date_of_birth: '',
  })
  const [error, setError] = React.useState('')

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    // Validation
    if (!formData.email || !formData.password || !formData.first_name || !formData.last_name || !formData.date_of_birth) {
      setError('All fields are required')
      return
    }

    if (formData.password.length < 8) {
      setError('Password must be at least 8 characters')
      return
    }

    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match')
      return
    }

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(formData.email)) {
      setError('Invalid email format')
      return
    }

    try {
      await register({
        email: formData.email,
        password: formData.password,
        first_name: formData.first_name,
        last_name: formData.last_name,
        date_of_birth: formData.date_of_birth,
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    }
  }

  return (
    <div>
      <h1>Register</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          name="first_name"
          placeholder="First Name"
          value={formData.first_name}
          onChange={handleChange}
          data-testid="first-name-input"
        />
        <input
          type="text"
          name="last_name"
          placeholder="Last Name"
          value={formData.last_name}
          onChange={handleChange}
          data-testid="last-name-input"
        />
        <input
          type="email"
          name="email"
          placeholder="Email"
          value={formData.email}
          onChange={handleChange}
          data-testid="email-input"
        />
        <input
          type="date"
          name="date_of_birth"
          value={formData.date_of_birth}
          onChange={handleChange}
          data-testid="dob-input"
        />
        <input
          type="password"
          name="password"
          placeholder="Password"
          value={formData.password}
          onChange={handleChange}
          data-testid="password-input"
        />
        <input
          type="password"
          name="confirmPassword"
          placeholder="Confirm Password"
          value={formData.confirmPassword}
          onChange={handleChange}
          data-testid="confirm-password-input"
        />
        <button type="submit" disabled={isLoading} data-testid="submit-button">
          {isLoading ? 'Loading...' : 'Register'}
        </button>
      </form>
      {error && <div data-testid="error-message">{error}</div>}
    </div>
  )
}

describe('Register Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render registration form', () => {
    render(<RegisterPage />)

    expect(screen.getByText('Register')).toBeInTheDocument()
    expect(screen.getByTestId('first-name-input')).toBeInTheDocument()
    expect(screen.getByTestId('last-name-input')).toBeInTheDocument()
    expect(screen.getByTestId('email-input')).toBeInTheDocument()
    expect(screen.getByTestId('dob-input')).toBeInTheDocument()
    expect(screen.getByTestId('password-input')).toBeInTheDocument()
    expect(screen.getByTestId('confirm-password-input')).toBeInTheDocument()
    expect(screen.getByTestId('submit-button')).toBeInTheDocument()
  })

  it('should validate required fields', async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)

    const submitButton = screen.getByTestId('submit-button')
    await user.click(submitButton)

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent('All fields are required')
    })

    expect(mockRegister).not.toHaveBeenCalled()
  })

  it('should validate password length', async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)

    await user.type(screen.getByTestId('first-name-input'), 'John')
    await user.type(screen.getByTestId('last-name-input'), 'Doe')
    await user.type(screen.getByTestId('email-input'), 'john@example.com')
    await user.type(screen.getByTestId('dob-input'), '1990-01-01')
    await user.type(screen.getByTestId('password-input'), 'short')
    await user.type(screen.getByTestId('confirm-password-input'), 'short')

    await user.click(screen.getByTestId('submit-button'))

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent(
        'Password must be at least 8 characters'
      )
    })

    expect(mockRegister).not.toHaveBeenCalled()
  })

  it('should validate password match', async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)

    await user.type(screen.getByTestId('first-name-input'), 'John')
    await user.type(screen.getByTestId('last-name-input'), 'Doe')
    await user.type(screen.getByTestId('email-input'), 'john@example.com')
    await user.type(screen.getByTestId('dob-input'), '1990-01-01')
    await user.type(screen.getByTestId('password-input'), 'password123')
    await user.type(screen.getByTestId('confirm-password-input'), 'password456')

    await user.click(screen.getByTestId('submit-button'))

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent('Passwords do not match')
    })

    expect(mockRegister).not.toHaveBeenCalled()
  })

  it('should validate email format', async () => {
    const user = userEvent.setup()
    render(<RegisterPage />)

    await user.type(screen.getByTestId('first-name-input'), 'John')
    await user.type(screen.getByTestId('last-name-input'), 'Doe')
    await user.type(screen.getByTestId('email-input'), 'invalid-email')
    await user.type(screen.getByTestId('dob-input'), '1990-01-01')
    await user.type(screen.getByTestId('password-input'), 'password123')
    await user.type(screen.getByTestId('confirm-password-input'), 'password123')

    await user.click(screen.getByTestId('submit-button'))

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent('Invalid email format')
    })

    expect(mockRegister).not.toHaveBeenCalled()
  })

  it('should submit valid registration data', async () => {
    mockRegister.mockResolvedValue(undefined)

    const user = userEvent.setup()
    render(<RegisterPage />)

    await user.type(screen.getByTestId('first-name-input'), 'John')
    await user.type(screen.getByTestId('last-name-input'), 'Doe')
    await user.type(screen.getByTestId('email-input'), 'john@example.com')
    await user.type(screen.getByTestId('dob-input'), '1990-01-01')
    await user.type(screen.getByTestId('password-input'), 'password123')
    await user.type(screen.getByTestId('confirm-password-input'), 'password123')

    await user.click(screen.getByTestId('submit-button'))

    await waitFor(() => {
      expect(mockRegister).toHaveBeenCalledWith({
        email: 'john@example.com',
        password: 'password123',
        first_name: 'John',
        last_name: 'Doe',
        date_of_birth: '1990-01-01',
      })
    })
  })

  it('should handle registration failure', async () => {
    mockRegister.mockRejectedValue(new Error('Email already exists'))

    const user = userEvent.setup()
    render(<RegisterPage />)

    await user.type(screen.getByTestId('first-name-input'), 'John')
    await user.type(screen.getByTestId('last-name-input'), 'Doe')
    await user.type(screen.getByTestId('email-input'), 'existing@example.com')
    await user.type(screen.getByTestId('dob-input'), '1990-01-01')
    await user.type(screen.getByTestId('password-input'), 'password123')
    await user.type(screen.getByTestId('confirm-password-input'), 'password123')

    await user.click(screen.getByTestId('submit-button'))

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent('Email already exists')
    })
  })

  it('should disable submit button while loading', () => {
    mockUseAuth.mockReturnValue({
      register: mockRegister,
      isLoading: true,
    })

    render(<RegisterPage />)

    const submitButton = screen.getByTestId('submit-button')
    expect(submitButton).toBeDisabled()
    expect(submitButton).toHaveTextContent('Loading...')
  })

  it('should accept valid email formats', async () => {
    mockRegister.mockResolvedValue(undefined)

    const validEmails = [
      'test@example.com',
      'user.name@example.com',
      'user+tag@example.co.uk',
      'user_name@example-domain.com',
    ]

    for (const email of validEmails) {
      vi.clearAllMocks()
      const user = userEvent.setup()
      render(<RegisterPage />)

      await user.type(screen.getByTestId('first-name-input'), 'John')
      await user.type(screen.getByTestId('last-name-input'), 'Doe')
      await user.type(screen.getByTestId('email-input'), email)
      await user.type(screen.getByTestId('dob-input'), '1990-01-01')
      await user.type(screen.getByTestId('password-input'), 'password123')
      await user.type(screen.getByTestId('confirm-password-input'), 'password123')

      await user.click(screen.getByTestId('submit-button'))

      await waitFor(() => {
        expect(mockRegister).toHaveBeenCalled()
      })
    }
  })

  it('should reject invalid email formats', async () => {
    const invalidEmails = [
      'notanemail',
      '@example.com',
      'user@',
      'user @example.com',
      'user@example',
    ]

    for (const email of invalidEmails) {
      vi.clearAllMocks()
      const user = userEvent.setup()
      render(<RegisterPage />)

      await user.type(screen.getByTestId('first-name-input'), 'John')
      await user.type(screen.getByTestId('last-name-input'), 'Doe')
      await user.type(screen.getByTestId('email-input'), email)
      await user.type(screen.getByTestId('dob-input'), '1990-01-01')
      await user.type(screen.getByTestId('password-input'), 'password123')
      await user.type(screen.getByTestId('confirm-password-input'), 'password123')

      await user.click(screen.getByTestId('submit-button'))

      await waitFor(() => {
        expect(screen.getByTestId('error-message')).toHaveTextContent('Invalid email format')
      })

      expect(mockRegister).not.toHaveBeenCalled()
    }
  })
})
