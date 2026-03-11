import { ApiClient } from './client'
import type {
  RegisterRequest,
  LoginRequest,
  AuthResponse,
  TokenResponse,
} from '../types/api'

export class AuthService {
  private client: ApiClient

  constructor(client: ApiClient) {
    this.client = client
  }

  async register(data: RegisterRequest): Promise<AuthResponse> {
    return this.client.post<AuthResponse>('/api/auth/register', data)
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    return this.client.post<AuthResponse>('/api/auth/login', data)
  }

  async logout(): Promise<void> {
    return this.client.post<void>('/api/auth/logout')
  }

  async refreshToken(refreshToken: string): Promise<TokenResponse> {
    return this.client.post<TokenResponse>('/api/auth/refresh', {
      refresh_token: refreshToken,
    })
  }
}
