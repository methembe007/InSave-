import { ApiClient } from './client'
import type {
  SpendingAnalysis,
  SavingsPattern,
  Recommendation,
  FinancialHealthScore,
  TimePeriod,
} from '../types/api'

export class AnalyticsService {
  private client: ApiClient

  constructor(client: ApiClient) {
    this.client = client
  }

  async getSpendingAnalysis(period: TimePeriod): Promise<SpendingAnalysis> {
    const params = new URLSearchParams({
      start_date: period.start_date,
      end_date: period.end_date,
    })
    return this.client.get<SpendingAnalysis>(
      `/api/analytics/spending?${params.toString()}`
    )
  }

  async getSavingsPatterns(): Promise<SavingsPattern[]> {
    return this.client.get<SavingsPattern[]>('/api/analytics/patterns')
  }

  async getRecommendations(): Promise<Recommendation[]> {
    return this.client.get<Recommendation[]>('/api/analytics/recommendations')
  }

  async getFinancialHealth(): Promise<FinancialHealthScore> {
    return this.client.get<FinancialHealthScore>('/api/analytics/health')
  }
}
