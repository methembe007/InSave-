import { ApiClient } from './client'
import type { Notification } from '../types/api'

export class NotificationService {
  private client: ApiClient

  constructor(client: ApiClient) {
    this.client = client
  }

  async getNotifications(): Promise<Notification[]> {
    return this.client.get<Notification[]>('/api/notifications')
  }

  async markAsRead(id: string): Promise<void> {
    return this.client.put<void>(`/api/notifications/${id}/read`)
  }
}
