import type { ApiError } from '../types/api'

export class ApiClient {
  private baseUrl: string
  private getToken: () => string | null
  private onUnauthorized: () => void

  constructor(
    baseUrl: string,
    getToken: () => string | null,
    onUnauthorized: () => void
  ) {
    this.baseUrl = baseUrl
    this.getToken = getToken
    this.onUnauthorized = onUnauthorized
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
    retries = 3
  ): Promise<T> {
    const token = this.getToken()
    const url = `${this.baseUrl}${endpoint}`

    console.log('[ApiClient] Making request to:', endpoint, 'with token:', !!token)

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }

    // Add custom headers from options
    if (options.headers) {
      const optionsHeaders = options.headers as Record<string, string>
      Object.assign(headers, optionsHeaders)
    }

    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const config: RequestInit = {
      ...options,
      headers,
    }

    try {
      const response = await fetch(url, config)

      // Handle 401 Unauthorized
      if (response.status === 401) {
        console.error('[ApiClient] 401 Unauthorized for:', endpoint)
        // Only trigger unauthorized handler if this is an auth endpoint
        // For other endpoints, just throw the error without clearing tokens
        if (endpoint.includes('/api/auth/') || endpoint.includes('/api/user/profile')) {
          this.onUnauthorized()
        }
        throw new Error('Unauthorized')
      }

      // Handle other error responses
      if (!response.ok) {
        const errorData: ApiError = await response.json().catch(() => ({
          error: 'Unknown Error',
          message: response.statusText,
          status: response.status,
        }))

        throw new Error(errorData.message || `HTTP ${response.status}`)
      }

      // Handle 204 No Content
      if (response.status === 204) {
        return {} as T
      }

      return await response.json()
    } catch (error) {
      // Retry logic for network errors
      if (retries > 0 && error instanceof TypeError) {
        await new Promise((resolve) => setTimeout(resolve, 1000))
        return this.request<T>(endpoint, options, retries - 1)
      }

      // Enhance error message for network failures
      if (error instanceof TypeError) {
        throw new Error(`Network error: Unable to connect to ${this.baseUrl}. Please ensure the service is running.`)
      }

      throw error
    }
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' })
  }

  async post<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async put<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' })
  }
}
