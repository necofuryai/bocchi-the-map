/**
 * Test utilities for BDD-style frontend testing
 * 
 * This module exports all testing utilities needed for comprehensive
 * TDD+BDD hybrid testing in the Bocchi The Map project.
 */

// Re-export everything from render utilities
export * from './render-with-providers'

// Re-export BDD helpers
export * from './bdd-helpers'

// Re-export accessibility helpers
export * from './accessibility-helpers'

// Re-export testing library utilities
export * from '@testing-library/react'
export { default as userEvent } from '@testing-library/user-event'

// Import screen for internal use in utility functions
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

// Common test utilities
export { vi, expect, describe, test, it, beforeEach, afterEach, beforeAll, afterAll } from 'vitest'

// MSW utilities for API mocking
export { server } from '@/mocks/server'
export { http, HttpResponse } from 'msw'

// Import domain types for proper typing
import type { Review, DomainUser, Spot } from '@/types'

/**
 * Common test data factories
 */
export const TestDataFactory = {
  /**
   * Create a mock user for testing (MVP simplified)
   */
  createMockUser: (overrides: Partial<DomainUser> = {}) => ({
    id: 'test-user-1',
    email: 'test@example.com',
    name: 'Test User',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    ...overrides,
  }),

  /**
   * Create a mock review for testing (MVP simplified)
   */
  createMockReview: (overrides: Partial<Review> = {}) => ({
    id: 'test-review-1',
    spotId: 'test-spot-1',
    userId: 'test-user-1',
    rating: 5, // 1-5 stars only
    comment: 'Great spot for solo work!',
    createdAt: '2024-06-01T10:00:00Z',
    updatedAt: '2024-06-01T10:00:00Z',
    ...overrides,
  }),

  /**
   * Create a mock spot for testing (MVP)
   */
  createMockSpot: (overrides: Partial<Spot> = {}) => ({
    id: 'test-spot-1',
    name: 'Test Cafe',
    category: 'cafe',
    address: '123 Test Street, Test City',
    countryCode: 'JP',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    ...overrides,
  }),

}

/**
 * Common test assertions for BDD testing
 */
export const BDDAssertions = {
  /**
   * Assert that a loading state is displayed
   */
  expectLoadingState: () => {
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument()
  },

  /**
   * Assert that an error message is displayed
   */
  expectErrorMessage: (message?: string) => {
    const errorElement = screen.getByTestId('error-message')
    expect(errorElement).toBeInTheDocument()
    
    if (message) {
      expect(errorElement).toHaveTextContent(message)
    }
  },


  /**
   * Assert that authentication is required
   */
  expectAuthenticationRequired: () => {
    expect(screen.getByTestId('login-prompt')).toBeInTheDocument()
  },

  /**
   * Assert that pagination controls are working
   */
  expectPagination: (hasMore: boolean) => {
    if (hasMore) {
      expect(screen.getByTestId('load-more-button')).toBeInTheDocument()
    } else {
      expect(screen.queryByTestId('load-more-button')).not.toBeInTheDocument()
    }
  },
}

/**
 * Common user interaction helpers for BDD testing
 */
export const BDDActions = {

  /**
   * Simulate user login
   */
  login: async (email = 'test@example.com', password = 'password123') => {
    const emailInput = screen.getByLabelText(/email/i)
    const passwordInput = screen.getByLabelText(/password/i)
    const loginButton = screen.getByRole('button', { name: /login/i })

    await userEvent.type(emailInput, email)
    await userEvent.type(passwordInput, password)
    await userEvent.click(loginButton)
  },

  /**
   * Simulate user writing a review
   */
  writeReview: async (review: { rating: number; comment: string; tags?: string[] }) => {
    const ratingInput = screen.getByTestId('rating-input')
    const commentInput = screen.getByLabelText(/comment/i)
    const submitButton = screen.getByRole('button', { name: /submit review/i })

    // Set rating (implementation depends on rating component)
    await userEvent.click(ratingInput)
    
    // Type comment
    await userEvent.type(commentInput, review.comment)

    // Add tags if provided
    if (review.tags) {
      for (const tag of review.tags) {
        const tagInput = screen.getByPlaceholderText(/add tag/i)
        await userEvent.type(tagInput, tag)
        await userEvent.keyboard('{Enter}')
      }
    }

    await userEvent.click(submitButton)
  },
}