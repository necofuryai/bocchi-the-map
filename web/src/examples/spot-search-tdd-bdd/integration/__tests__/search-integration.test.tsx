/**
 * TDD+BDD Integration Test Example
 * 
 * This test demonstrates how components work together
 * to fulfill the user scenarios defined in the E2E tests.
 * 
 * Integration testing bridges the gap between unit tests (TDD)
 * and end-to-end tests (BDD), ensuring components integrate correctly.
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { SearchPage } from '../search-page'

// Mock the search API
vi.mock('@/examples/spot-search-tdd-bdd/hooks/use-spot-search', () => {
  const mockSpots = [
    {
      id: '1',
      name: 'Quiet Coffee House',
      type: 'cafe',
      soloFriendly: true,
      soloFriendlyRating: 4.8,
      averageRating: 4.5,
      reviewCount: 127,
    },
    {
      id: '2',
      name: 'Busy Downtown Cafe',
      type: 'cafe',
      soloFriendly: false,
      soloFriendlyRating: 2.1,
      averageRating: 4.2,
      reviewCount: 89,
    },
  ]

  return {
    useSpotSearch: vi.fn(() => ({
      spots: [],
      loading: false,
      error: null,
      hasMore: true,
      query: '',
      filters: {},
      total: 0,
      search: vi.fn(),
      loadMore: vi.fn(),
      clearSearch: vi.fn(),
      updateFilters: vi.fn(),
    })),
    searchSpots: vi.fn(),
  }
})

import { useSpotSearch } from '@/examples/spot-search-tdd-bdd/hooks/use-spot-search'

describe('Search Integration Tests', () => {
  const mockSearch = vi.fn()
  const mockLoadMore = vi.fn()
  const mockClearSearch = vi.fn()
  const mockUpdateFilters = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Given the SearchPage component is rendered', () => {
    describe('When user performs a search', () => {
      it('Then the search input and results should work together correctly', async () => {
        // Given - Mock the hook to return search results
        const mockSpots = [
          {
            id: '1',
            name: 'Quiet Coffee House',
            type: 'cafe',
            soloFriendly: true,
            soloFriendlyRating: 4.8,
            averageRating: 4.5,
            reviewCount: 127,
          },
          {
            id: '2',
            name: 'Study Library',
            type: 'library',
            soloFriendly: true,
            soloFriendlyRating: 4.9,
            averageRating: 4.7,
            reviewCount: 203,
          },
        ]

        vi.mocked(useSpotSearch).mockReturnValue({
          spots: mockSpots,
          loading: false,
          error: null,
          hasMore: false,
          query: 'quiet',
          filters: {},
          total: 2,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        const user = userEvent.setup()
        render(<SearchPage />)

        // Then - Search results should be displayed
        expect(screen.getByTestId('search-page')).toBeInTheDocument()
        expect(screen.getByTestId('search-results')).toBeInTheDocument()
        
        // And - Both spots should be visible
        expect(screen.getByText('Quiet Coffee House')).toBeInTheDocument()
        expect(screen.getByText('Study Library')).toBeInTheDocument()
        
        // And - Solo-friendly indicators should be shown
        const soloFriendlyBadges = screen.getAllByTestId('solo-friendly-badge')
        expect(soloFriendlyBadges).toHaveLength(2)
      })

      it('Then search input should trigger search function when user types and presses Enter', async () => {
        // Given
        vi.mocked(useSpotSearch).mockReturnValue({
          spots: [],
          loading: false,
          error: null,
          hasMore: true,
          query: '',
          filters: {},
          total: 0,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        const user = userEvent.setup()
        render(<SearchPage />)

        // When - User enters search query and presses Enter
        const searchInput = screen.getByTestId('search-input')
        await user.type(searchInput, 'coffee shops')
        await user.keyboard('{Enter}')

        // Then - Search function should be called
        expect(mockSearch).toHaveBeenCalledWith('coffee shops')
        expect(mockSearch).toHaveBeenCalledTimes(1)
      })
    })

    describe('When filters are applied', () => {
      it('Then updateFilters should be called with correct filter values', async () => {
        // Given
        vi.mocked(useSpotSearch).mockReturnValue({
          spots: [],
          loading: false,
          error: null,
          hasMore: true,
          query: 'cafe',
          filters: {},
          total: 0,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        const user = userEvent.setup()
        render(<SearchPage />)

        // When - User opens filters and applies solo-friendly filter
        const filterToggle = screen.getByTestId('filter-toggle')
        await user.click(filterToggle)

        const soloFriendlyFilter = screen.getByTestId('solo-friendly-filter')
        await user.click(soloFriendlyFilter)

        const applyFilters = screen.getByTestId('apply-filters')
        await user.click(applyFilters)

        // Then - updateFilters should be called with correct values
        expect(mockUpdateFilters).toHaveBeenCalledWith({ soloFriendly: true })
      })

      it('Then filtered results should display appropriate indicators', async () => {
        // Given - Mock filtered results
        const filteredSpots = [
          {
            id: '1',
            name: 'Solo-Friendly Cafe',
            type: 'cafe',
            soloFriendly: true,
            soloFriendlyRating: 4.8,
          },
        ]

        vi.mocked(useSpotSearch).mockReturnValue({
          spots: filteredSpots,
          loading: false,
          error: null,
          hasMore: false,
          query: 'cafe',
          filters: { soloFriendly: true },
          total: 1,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        render(<SearchPage />)

        // Then - Active filter indicator should be shown
        expect(screen.getByTestId('active-filters')).toBeInTheDocument()
        expect(screen.getByTestId('active-filters')).toHaveTextContent('Solo-friendly')
        
        // And - Only solo-friendly spots should be displayed
        expect(screen.getByText('Solo-Friendly Cafe')).toBeInTheDocument()
        expect(screen.getByTestId('solo-friendly-badge')).toBeInTheDocument()
      })
    })

    describe('When loading more results', () => {
      it('Then load more button should trigger loadMore function', async () => {
        // Given
        const initialSpots = [
          { id: '1', name: 'Cafe A', type: 'cafe', soloFriendly: true },
          { id: '2', name: 'Cafe B', type: 'cafe', soloFriendly: false },
        ]

        vi.mocked(useSpotSearch).mockReturnValue({
          spots: initialSpots,
          loading: false,
          error: null,
          hasMore: true,
          query: 'cafe',
          filters: {},
          total: 5,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        const user = userEvent.setup()
        render(<SearchPage />)

        // When - User clicks load more button
        const loadMoreButton = screen.getByTestId('load-more-button')
        await user.click(loadMoreButton)

        // Then - loadMore function should be called
        expect(mockLoadMore).toHaveBeenCalledTimes(1)
      })

      it('Then load more button should not be visible when no more results', () => {
        // Given
        vi.mocked(useSpotSearch).mockReturnValue({
          spots: [{ id: '1', name: 'Cafe A', type: 'cafe' }],
          loading: false,
          error: null,
          hasMore: false,
          query: 'cafe',
          filters: {},
          total: 1,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        render(<SearchPage />)

        // Then - Load more button should not be present
        expect(screen.queryByTestId('load-more-button')).not.toBeInTheDocument()
      })
    })

    describe('When search encounters an error', () => {
      it('Then error message should be displayed with retry option', async () => {
        // Given
        vi.mocked(useSpotSearch).mockReturnValue({
          spots: [],
          loading: false,
          error: 'Network error occurred',
          hasMore: false,
          query: 'cafe',
          filters: {},
          total: 0,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        const user = userEvent.setup()
        render(<SearchPage />)

        // Then - Error message should be displayed
        expect(screen.getByTestId('error-message')).toBeInTheDocument()
        expect(screen.getByTestId('error-message')).toHaveTextContent('Network error occurred')
        
        // And - Retry button should be available
        const retryButton = screen.getByTestId('retry-search')
        expect(retryButton).toBeInTheDocument()

        // When - User clicks retry
        await user.click(retryButton)

        // Then - Search should be called again with the same query
        expect(mockSearch).toHaveBeenCalledWith('cafe')
      })
    })

    describe('When search returns no results', () => {
      it('Then no results message should be displayed with suggestions', () => {
        // Given
        vi.mocked(useSpotSearch).mockReturnValue({
          spots: [],
          loading: false,
          error: null,
          hasMore: false,
          query: 'nonexistent place',
          filters: {},
          total: 0,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        render(<SearchPage />)

        // Then - No results message should be displayed
        expect(screen.getByTestId('no-results-message')).toBeInTheDocument()
        expect(screen.getByTestId('no-results-message')).toHaveTextContent('No spots found')
        
        // And - Search suggestions should be provided
        expect(screen.getByTestId('search-suggestions')).toBeInTheDocument()
      })
    })

    describe('When search is loading', () => {
      it('Then loading state should be displayed across components', () => {
        // Given
        vi.mocked(useSpotSearch).mockReturnValue({
          spots: [],
          loading: true,
          error: null,
          hasMore: true,
          query: 'cafe',
          filters: {},
          total: 0,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        render(<SearchPage />)

        // Then - Loading indicators should be displayed
        expect(screen.getByTestId('search-loading')).toBeInTheDocument()
        expect(screen.getByTestId('results-loading')).toBeInTheDocument()
        
        // And - Search input should be disabled
        const searchInput = screen.getByTestId('search-input')
        expect(searchInput).toBeDisabled()
      })
    })
  })

  describe('Given the user is on mobile device', () => {
    beforeEach(() => {
      // Mock mobile viewport
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value: 375,
      })
    })

    describe('When filters are opened', () => {
      it('Then filter panel should slide up from bottom', async () => {
        // Given
        vi.mocked(useSpotSearch).mockReturnValue({
          spots: [],
          loading: false,
          error: null,
          hasMore: true,
          query: '',
          filters: {},
          total: 0,
          search: mockSearch,
          loadMore: mockLoadMore,
          clearSearch: mockClearSearch,
          updateFilters: mockUpdateFilters,
        })

        const user = userEvent.setup()
        render(<SearchPage />)

        // When - User opens filters on mobile
        const filterToggle = screen.getByTestId('filter-toggle')
        await user.click(filterToggle)

        // Then - Filter panel should have mobile slide-up animation
        const filterPanel = screen.getByTestId('filter-panel')
        expect(filterPanel).toBeInTheDocument()
        expect(filterPanel).toHaveClass('slide-up')
        
        // And - Close button should be available
        expect(screen.getByTestId('close-filters')).toBeInTheDocument()
      })
    })
  })
})

/**
 * Integration Test Summary:
 * 
 * These integration tests verify that:
 * 1. SearchInput component integrates correctly with useSpotSearch hook
 * 2. User interactions flow properly between components
 * 3. State changes propagate correctly through the component tree
 * 4. Error handling works across the integrated system
 * 5. Loading states are coordinated between components
 * 6. Mobile responsiveness works in integrated scenarios
 * 
 * The tests bridge the gap between:
 * - Unit tests (individual component behavior)
 * - E2E tests (full user scenarios)
 * 
 * This ensures that the components work together to fulfill
 * the user stories defined in the BDD E2E tests.
 */