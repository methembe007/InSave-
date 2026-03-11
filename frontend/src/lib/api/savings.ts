import { ApiClient } from './client'
import type {
  SavingsSummary,
  SavingsTransaction,
  CreateSavingsRequest,
  SavingsStreak,
  HistoryParams,
} from '../types/api'

export class SavingsService {
  private client: ApiClient

  constructor(client: ApiClient) {
    this.client = client
  }

  async getSummary(): Promise<SavingsSummary> {
    return this.client.get<SavingsSummary>('/api/savings/summary')
  }

  async getHistory(params?: HistoryParams): Promise<SavingsTransaction[]> {
    const queryParams = new URLSearchParams()
    if (params?.limit) queryParams.append('limit', params.limit.toString())
    if (params?.offset) queryParams.append('offset', params.offset.toString())
    if (params?.start_date) queryParams.append('start_date', params.start_date)
    if (params?.end_date) queryParams.append('end_date', params.end_date)

    const query = queryParams.toString()
    return this.client.get<SavingsTransaction[]>(
      `/api/savings/history${query ? `?${query}` : ''}`
    )
  }

  async createTransaction(
    data: CreateSavingsRequest
  ): Promise<SavingsTransaction> {
    return this.client.post<SavingsTransaction>('/api/savings/transactions', data)
  }

  async getStreak(): Promise<SavingsStreak> {
    return this.client.get<SavingsStreak>('/api/savings/streak')
  }
}
