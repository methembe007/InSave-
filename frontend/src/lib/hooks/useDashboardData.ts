import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../auth/context'
import type {
  SavingsSummary,
  Budget,
  Goal,
  FinancialHealthScore,
  SavingsTransaction,
  SpendingTransaction,
} from '../types/api'

// Requirement 4.5: Fetch savings summary
export function useSavingsSummary() {
  const { api } = useAuth()
  return useQuery<SavingsSummary>({
    queryKey: ['savings', 'summary'],
    queryFn: () => api.savings.getSummary(),
  })
}

// Requirement 6.3: Fetch current budget
export function useCurrentBudget() {
  const { api } = useAuth()
  return useQuery<Budget | null>({
    queryKey: ['budget', 'current'],
    queryFn: async () => {
      try {
        return await api.budget.getCurrentBudget()
      } catch (error: any) {
        // If no budget exists, return null instead of throwing
        if (error.status === 404) {
          return null
        }
        throw error
      }
    },
  })
}

// Requirement 9.3: Fetch active goals
export function useActiveGoals() {
  const { api } = useAuth()
  return useQuery<Goal[]>({
    queryKey: ['goals', 'active'],
    queryFn: () => api.goals.getActiveGoals(),
  })
}

// Requirement 14.1: Fetch financial health score
export function useFinancialHealth() {
  const { api } = useAuth()
  return useQuery<FinancialHealthScore | null>({
    queryKey: ['analytics', 'health'],
    queryFn: async () => {
      try {
        return await api.analytics.getFinancialHealth()
      } catch (error: any) {
        // If insufficient data, return null
        if (error.status === 400) {
          return null
        }
        throw error
      }
    },
  })
}

// Requirement 4.4: Fetch recent savings transactions
export function useRecentSavings(limit: number = 5) {
  const { api } = useAuth()
  return useQuery<SavingsTransaction[]>({
    queryKey: ['savings', 'history', { limit }],
    queryFn: () => api.savings.getHistory({ limit }),
  })
}

// Requirement 7.1: Fetch recent spending transactions
export function useRecentSpending(limit: number = 5) {
  return useQuery<SpendingTransaction[]>({
    queryKey: ['spending', 'history', { limit }],
    queryFn: async () => {
      // Note: This would need to be implemented in the budget service
      // For now, return empty array as placeholder
      return []
    },
  })
}
