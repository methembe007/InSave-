import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { ProgressTracker } from '../ProgressTracker'
import type { EducationProgress } from '../../lib/types/api'

/**
 * Tests for ProgressTracker component
 * Validates Requirements 11.4, 11.5
 */

describe('ProgressTracker', () => {
  it('should render progress data correctly', () => {
    const progress: EducationProgress = {
      total_lessons: 10,
      completed_lessons: 5,
      progress_percent: 50,
      current_streak: 3,
    }

    render(<ProgressTracker progress={progress} />)

    // Check if all stats are displayed
    expect(screen.getByText('10')).toBeTruthy()
    expect(screen.getByText('Total Lessons')).toBeTruthy()
    expect(screen.getByText('5')).toBeTruthy()
    expect(screen.getByText('Completed')).toBeTruthy()
    expect(screen.getByText('3')).toBeTruthy()
    expect(screen.getByText('Day Streak')).toBeTruthy()
    expect(screen.getByText('50.0%')).toBeTruthy()
  })

  it('should render nothing when progress is undefined', () => {
    const { container } = render(<ProgressTracker progress={undefined} />)
    expect(container.firstChild).toBeNull()
  })

  it('should display correct progress percentage', () => {
    const progress: EducationProgress = {
      total_lessons: 20,
      completed_lessons: 15,
      progress_percent: 75,
      current_streak: 7,
    }

    render(<ProgressTracker progress={progress} />)

    expect(screen.getByText('75.0%')).toBeTruthy()
    expect(screen.getByText('15 of 20 lessons completed')).toBeTruthy()
  })

  it('should handle zero progress', () => {
    const progress: EducationProgress = {
      total_lessons: 10,
      completed_lessons: 0,
      progress_percent: 0,
      current_streak: 0,
    }

    render(<ProgressTracker progress={progress} />)

    expect(screen.getByText('Completed')).toBeTruthy()
    expect(screen.getByText('0.0%')).toBeTruthy()
    expect(screen.getByText('0 of 10 lessons completed')).toBeTruthy()
  })

  it('should handle 100% completion', () => {
    const progress: EducationProgress = {
      total_lessons: 10,
      completed_lessons: 10,
      progress_percent: 100,
      current_streak: 10,
    }

    render(<ProgressTracker progress={progress} />)

    expect(screen.getByText('100.0%')).toBeTruthy()
    expect(screen.getByText('10 of 10 lessons completed')).toBeTruthy()
  })
})
