import { ApiClient } from './client'
import type {
  Goal,
  CreateGoalRequest,
  UpdateGoalRequest,
  Milestone,
} from '../types/api'

export class GoalService {
  private client: ApiClient

  constructor(client: ApiClient) {
    this.client = client
  }

  async getActiveGoals(): Promise<Goal[]> {
    return this.client.get<Goal[]>('/api/goals')
  }

  async getGoal(id: string): Promise<Goal> {
    return this.client.get<Goal>(`/api/goals/${id}`)
  }

  async createGoal(data: CreateGoalRequest): Promise<Goal> {
    return this.client.post<Goal>('/api/goals', data)
  }

  async updateGoal(id: string, data: UpdateGoalRequest): Promise<Goal> {
    return this.client.put<Goal>(`/api/goals/${id}`, data)
  }

  async deleteGoal(id: string): Promise<void> {
    return this.client.delete<void>(`/api/goals/${id}`)
  }

  async getMilestones(goalId: string): Promise<Milestone[]> {
    return this.client.get<Milestone[]>(`/api/goals/${goalId}/milestones`)
  }

  async addProgress(goalId: string, amount: number): Promise<Goal> {
    return this.client.post<Goal>(`/api/goals/${goalId}/progress`, { amount })
  }
}
