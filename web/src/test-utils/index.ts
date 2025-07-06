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
export { userEvent } from '@testing-library/user-event'

// Common test utilities
export { vi, expect, describe, test, it, beforeEach, afterEach, beforeAll, afterAll } from 'vitest'

// MSW utilities for API mocking
export { server } from '@/mocks/server'
export { http, HttpResponse } from 'msw'

/**
 * Common test data factories
 */
export const TestDataFactory = {
  /**
   * Create a mock spot for testing
   */
  createMockSpot: (overrides: any = {}) => ({
    id: 'test-spot-1',
    name: 'Test Cafe',
    type: 'cafe',
    address: '123 Test St, Tokyo',
    latitude: 35.6762,
    longitude: 139.6503,
    soloFriendly: true,
    soloFriendlyRating: 4.5,
    averageRating: 4.2,
    reviewCount: 42,
    amenities: ['wifi', 'quiet', 'power_outlets'],
    description: 'A test cafe for unit testing',
    photos: ['/images/test-cafe.jpg'],
    openingHours: {
      monday: '08:00-20:00',
      tuesday: '08:00-20:00',
      wednesday: '08:00-20:00',
      thursday: '08:00-20:00',
      friday: '08:00-20:00',
      saturday: '09:00-21:00',
      sunday: '09:00-21:00',
    },
    ...overrides,
  }),

  /**
   * Create a mock user for testing
   */
  createMockUser: (overrides: any = {}) => ({
    id: 'test-user-1',
    email: 'test@example.com',
    name: 'Test User',
    avatar: '/images/test-avatar.jpg',
    preferences: {
      theme: 'light',
      notifications: true,
      language: 'ja',
    },
    createdAt: '2024-01-01T00:00:00Z',
    ...overrides,
  }),

  /**
   * Create a mock review for testing
   */
  createMockReview: (overrides: any = {}) => ({
    id: 'test-review-1',
    spotId: 'test-spot-1',
    userId: 'test-user-1',
    userName: 'Test User',
    userAvatar: '/images/test-avatar.jpg',
    rating: 5,
    soloFriendlyRating: 5,
    comment: 'Great spot for solo work!',
    tags: ['quiet', 'wifi', 'solo-friendly'],
    photos: [],
    helpful: 5,
    notHelpful: 0,
    createdAt: '2024-06-01T10:00:00Z',
    updatedAt: '2024-06-01T10:00:00Z',
    ...overrides,
  }),

  /**
   * Create mock search results for testing
   */
  createMockSearchResults: (count = 3, overrides: any = {}) => ({
    data: Array.from({ length: count }, (_, i) => 
      TestDataFactory.createMockSpot({
        id: `test-spot-${i + 1}`,
        name: `Test Spot ${i + 1}`,
        ...overrides,
      })
    ),
    total: count,
    hasMore: false,
    offset: 0,
    limit: 10,
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
   * Assert that search results are displayed
   */
  expectSearchResults: (count?: number) => {
    const resultsContainer = screen.getByTestId('search-results')
    expect(resultsContainer).toBeInTheDocument()
    
    if (count !== undefined) {
      const spotItems = screen.getAllByTestId('spot-item')
      expect(spotItems).toHaveLength(count)
    }
  },

  /**
   * Assert that authentication is required
   */
  expectAuthenticationRequired: () => {
    expect(screen.getByTestId('login-prompt')).toBeInTheDocument()
  },

  /**
   * Assert that a specific spot is displayed
   */
  expectSpotDetails: (spotName: string) => {
    expect(screen.getByTestId('spot-details')).toBeInTheDocument()
    expect(screen.getByTestId('spot-title')).toHaveTextContent(spotName)
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
   * Simulate user performing a search
   */
  performSearch: async (query: string) => {
    const searchInput = screen.getByPlaceholderText('Search for spots...')
    await userEvent.clear(searchInput)
    await userEvent.type(searchInput, query)
    await userEvent.keyboard('{Enter}')
  },

  /**
   * Simulate user clicking a spot item
   */
  clickSpotItem: async (index = 0) => {
    const spotItems = screen.getAllByTestId('spot-item')
    await userEvent.click(spotItems[index])
  },

  /**
   * Simulate user applying filters
   */
  applyFilters: async (filters: { soloFriendly?: boolean; type?: string }) => {
    const filterButton = screen.getByTestId('filter-button')
    await userEvent.click(filterButton)

    if (filters.soloFriendly) {
      const soloFriendlyFilter = screen.getByTestId('solo-friendly-filter')
      await userEvent.click(soloFriendlyFilter)
    }

    if (filters.type) {
      const typeFilter = screen.getByTestId(`type-filter-${filters.type}`)
      await userEvent.click(typeFilter)
    }

    const applyButton = screen.getByTestId('apply-filters')
    await userEvent.click(applyButton)
  },

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