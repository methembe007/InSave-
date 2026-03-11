import { ApiClient } from './client'
import type {
  Budget,
  BudgetCategory,
  CreateBudgetRequest,
  UpdateBudgetRequest,
  SpendingRequest,
  SpendingTransaction,
  BudgetAlert,
} from '../types/api'

export class BudgetService {
  private client: ApiClient

  constructor(client: ApiClient) {
    this.client = client
  }

  async getCurrentBudget(): Promise<Budget> {
    return this.client.get<Budget>('/api/budget/current')
  }

  async createBudget(data: CreateBudgetRequest): Promise<Budget> {
    return this.client.post<Budget>('/api/budget', data)
  }

  async updateBudget(id: string, data: UpdateBudgetRequest): Promise<Budget> {
    return this.client.put<Budget>(`/api/budget/${id}`, data)
  }

  async getCategories(): Promise<BudgetCategory[]> {
    return this.client.get<BudgetCategory[]>('/api/budget/categories')
  }

  async recordSpending(data: SpendingRequest): Promise<SpendingTransaction> {
    return this.client.post<SpendingTransaction>('/api/budget/spending', data)
  }

  async getAlerts(): Promise<BudgetAlert[]> {
    return this.client.get<BudgetAlert[]>('/api/budget/alerts')
  }
}
