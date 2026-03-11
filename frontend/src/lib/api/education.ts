import { ApiClient } from './client'
import type { Lesson, LessonDetail, EducationProgress } from '../types/api'

export class EducationService {
  private client: ApiClient

  constructor(client: ApiClient) {
    this.client = client
  }

  async getLessons(): Promise<Lesson[]> {
    return this.client.get<Lesson[]>('/api/education/lessons')
  }

  async getLesson(id: string): Promise<LessonDetail> {
    return this.client.get<LessonDetail>(`/api/education/lessons/${id}`)
  }

  async markLessonComplete(id: string): Promise<void> {
    return this.client.post<void>(`/api/education/lessons/${id}/complete`)
  }

  async getProgress(): Promise<EducationProgress> {
    return this.client.get<EducationProgress>('/api/education/progress')
  }
}
