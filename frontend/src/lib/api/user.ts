import { ApiClient } from './client'
import type {
  UserProfile,
  UpdateProfileRequest,
  UserPreferences,
} from '../types/api'

export class UserService {
  private client: ApiClient

  constructor(client: ApiClient) {
    this.client = client
  }

  async getProfile(): Promise<UserProfile> {
    return this.client.get<UserProfile>('/api/user/profile')
  }

  async updateProfile(data: UpdateProfileRequest): Promise<UserProfile> {
    return this.client.put<UserProfile>('/api/user/profile', data)
  }

  async getPreferences(): Promise<UserPreferences> {
    return this.client.get<UserPreferences>('/api/user/preferences')
  }

  async updatePreferences(data: UserPreferences): Promise<UserPreferences> {
    return this.client.put<UserPreferences>('/api/user/preferences', data)
  }

  async deleteAccount(): Promise<void> {
    return this.client.delete<void>('/api/user/account')
  }
}
