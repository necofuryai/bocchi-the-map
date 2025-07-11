/**
 * useSpotSearch Hook Implementation
 * 
 * Custom hook for managing spot search state and operations.
 * Developed using TDD methodology: RED -> GREEN -> REFACTOR
 * 
 * This hook provides:
 * - Search functionality with query and filters
 * - Pagination (load more)
 * - Loading and error states
 * - Race condition handling
 */

import { useState, useCallback, useRef } from 'react'

// Type definitions
interface Spot {
  id: string
  name: string
  type: string
  address?: string
  latitude?: number
  longitude?: number
  soloFriendly?: boolean
  soloFriendlyRating?: number
  averageRating?: number
  reviewCount?: number
  amenities?: string[]
  description?: string
  photos?: string[]
  openingHours?: Record<string, string>
}

interface SearchFilters {
  soloFriendly?: boolean
  type?: string
  [key: string]: any
}

interface SearchResponse {
  data: Spot[]
  total: number
  hasMore: boolean
  offset: number
  limit: number
}

interface UseSpotSearchReturn {
  spots: Spot[]
  loading: boolean
  error: string | null
  hasMore: boolean
  query: string
  filters: SearchFilters
  total: number
  search: (query: string) => Promise<void>
  loadMore: () => Promise<void>
  clearSearch: () => void
  updateFilters: (newFilters: SearchFilters) => Promise<void>
}

// Mock search API function for development
const searchSpots = async (query: string, options: any = {}): Promise<SearchResponse> => {
  // This would be replaced with actual API call
  // For now, return mock data for testing
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({
        data: [],
        total: 0,
        hasMore: false,
        offset: 0,
        limit: 10,
      })
    }, 100)
  })
}

export function useSpotSearch(): UseSpotSearchReturn {
  // State management
  const [spots, setSpots] = useState<Spot[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [hasMore, setHasMore] = useState(true)
  const [query, setQuery] = useState('')
  const [filters, setFilters] = useState<SearchFilters>({})
  const [total, setTotal] = useState(0)

  // Reference for tracking the latest search request
  const searchIdRef = useRef(0)

  // Helper function to extract error message
  const getErrorMessage = (error: any): string => {
    if (error?.response?.data?.error) {
      return error.response.data.error
    }
    if (error?.message) {
      return error.message
    }
    return 'An unexpected error occurred'
  }

  // Search function
  const search = useCallback(async (searchQuery: string) => {
    if (!searchQuery.trim()) {
      return
    }

    // Increment search ID to handle race conditions
    const currentSearchId = ++searchIdRef.current

    setLoading(true)
    setError(null)
    setQuery(searchQuery)

    try {
      const searchOptions = {
        offset: 0,
        limit: 10,
        ...filters,
      }

      const response = await searchSpots(searchQuery, searchOptions)

      // Check if this is still the latest search
      if (currentSearchId === searchIdRef.current) {
        setSpots(response.data)
        setTotal(response.total)
        setHasMore(response.hasMore)
        setLoading(false)
      }
    } catch (err: any) {
      // Only update error if this is still the latest search
      if (currentSearchId === searchIdRef.current) {
        setError(getErrorMessage(err))
        setLoading(false)
        setSpots([])
        setTotal(0)
      }
    }
  }, [filters])

  // Load more function
  const loadMore = useCallback(async () => {
    if (!hasMore || loading || !query) {
      return
    }

    setLoading(true)
    setError(null)

    try {
      const searchOptions = {
        offset: spots.length,
        limit: 10,
        ...filters,
      }

      const response = await searchSpots(query, searchOptions)

      setSpots(prevSpots => [...prevSpots, ...response.data])
      setHasMore(response.hasMore)
      setTotal(response.total)
      setLoading(false)
    } catch (err: any) {
      setError(getErrorMessage(err))
      setLoading(false)
    }
  }, [hasMore, loading, query, spots.length, filters])

  // Clear search function
  const clearSearch = useCallback(() => {
    setSpots([])
    setQuery('')
    setError(null)
    setTotal(0)
    setHasMore(true)
    setLoading(false)
  }, [])

  // Update filters function
  const updateFilters = useCallback(async (newFilters: SearchFilters) => {
    setFilters(newFilters)

    // If there's an active query, re-search with new filters
    if (query) {
      // Reset spots before searching with new filters
      setSpots([])
      setTotal(0)
      setHasMore(true)
      
      const currentSearchId = ++searchIdRef.current
      setLoading(true)
      setError(null)

      try {
        const searchOptions = {
          offset: 0,
          limit: 10,
          ...newFilters,
        }

        const response = await searchSpots(query, searchOptions)

        // Check if this is still the latest search
        if (currentSearchId === searchIdRef.current) {
          setSpots(response.data)
          setTotal(response.total)
          setHasMore(response.hasMore)
          setLoading(false)
        }
      } catch (err: any) {
        if (currentSearchId === searchIdRef.current) {
          setError(getErrorMessage(err))
          setLoading(false)
        }
      }
    }
  }, [query])

  return {
    spots,
    loading,
    error,
    hasMore,
    query,
    filters,
    total,
    search,
    loadMore,
    clearSearch,
    updateFilters,
  }
}

// Export the search function for testing purposes
export { searchSpots }

/**
 * Development Notes:
 * 
 * This hook was built using TDD methodology:
 * 
 * 1. RED Phase:
 *    - Comprehensive test suite written first
 *    - Tests covered all expected functionality and edge cases
 *    - All tests initially failed
 * 
 * 2. GREEN Phase:
 *    - Implemented minimal functionality to pass tests
 *    - Focus on making tests pass rather than perfect code
 *    - Handled all the test scenarios
 * 
 * 3. REFACTOR Phase:
 *    - Improved code organization and readability
 *    - Added proper TypeScript types
 *    - Optimized performance with useCallback
 *    - Added race condition handling
 * 
 * Key Features:
 * - Search with query and filters
 * - Pagination support (load more)
 * - Race condition prevention
 * - Comprehensive error handling
 * - Loading states
 * - TypeScript support
 * 
 * Race Condition Handling:
 * - Uses searchIdRef to track the latest search request
 * - Only updates state if the response is from the latest search
 * - Prevents stale data from overriding newer results
 * 
 * Next Steps:
 * - Integrate with actual search API
 * - Add caching mechanism
 * - Implement debouncing for search queries
 * - Add search history management
 */