/**
 * TDD Unit Test Example: useSpotSearch Hook
 * 
 * This demonstrates TDD for custom React hooks:
 * RED -> GREEN -> REFACTOR
 * 
 * The hook manages search state and API calls,
 * supporting the search functionality required by the E2E tests.
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, act, waitFor } from '@testing-library/react'
import { useSpotSearch } from '../use-spot-search'

// Mock the API service
vi.mock('@/services/spot-api', () => ({
  searchSpots: vi.fn(),
}))

import { searchSpots } from '@/services/spot-api'

describe('useSpotSearch Hook', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('Given the useSpotSearch hook is initialized', () => {
    describe('When the hook is first called', () => {
      it('Then it should return initial state correctly', () => {
        // Given & When
        const { result } = renderHook(() => useSpotSearch())
        
        // Then
        expect(result.current.spots).toEqual([])
        expect(result.current.loading).toBe(false)
        expect(result.current.error).toBeNull()
        expect(result.current.hasMore).toBe(true)
        expect(result.current.query).toBe('')
        expect(result.current.filters).toEqual({})
        expect(result.current.total).toBe(0)
        expect(typeof result.current.search).toBe('function')
        expect(typeof result.current.loadMore).toBe('function')
        expect(typeof result.current.clearSearch).toBe('function')
        expect(typeof result.current.updateFilters).toBe('function')
      })
    })

    describe('When search is called with a query', () => {
      it('Then it should set loading state and fetch spots successfully', async () => {
        // Given
        const mockSpots = [
          { id: '1', name: 'Quiet Cafe', type: 'cafe', soloFriendly: true },
          { id: '2', name: 'Study Library', type: 'library', soloFriendly: true },
        ]
        const mockResponse = {
          data: mockSpots,
          total: 2,
          hasMore: false,
          offset: 0,
          limit: 10,
        }
        
        vi.mocked(searchSpots).mockResolvedValueOnce(mockResponse)
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When
        await act(async () => {
          await result.current.search('coffee')
        })
        
        // Then
        expect(result.current.spots).toEqual(mockSpots)
        expect(result.current.loading).toBe(false)
        expect(result.current.error).toBeNull()
        expect(result.current.hasMore).toBe(false)
        expect(result.current.query).toBe('coffee')
        expect(result.current.total).toBe(2)
        
        expect(searchSpots).toHaveBeenCalledWith('coffee', {
          offset: 0,
          limit: 10,
        })
      })

      it('Then it should handle loading state correctly during search', async () => {
        // Given
        vi.mocked(searchSpots).mockImplementation(
          () => new Promise(resolve => setTimeout(() => resolve({
            data: [],
            total: 0,
            hasMore: false,
            offset: 0,
            limit: 10,
          }), 100))
        )
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When
        act(() => {
          result.current.search('coffee')
        })
        
        // Then - Should be loading immediately
        expect(result.current.loading).toBe(true)
        expect(result.current.error).toBeNull()
        
        // Wait for search to complete
        await waitFor(() => {
          expect(result.current.loading).toBe(false)
        })
      })

      it('Then it should handle empty search results', async () => {
        // Given
        const mockResponse = {
          data: [],
          total: 0,
          hasMore: false,
          offset: 0,
          limit: 10,
        }
        
        vi.mocked(searchSpots).mockResolvedValueOnce(mockResponse)
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When
        await act(async () => {
          await result.current.search('nonexistent')
        })
        
        // Then
        expect(result.current.spots).toEqual([])
        expect(result.current.total).toBe(0)
        expect(result.current.hasMore).toBe(false)
        expect(result.current.loading).toBe(false)
        expect(result.current.error).toBeNull()
      })
    })

    describe('When search fails', () => {
      it('Then it should set error state and clear loading', async () => {
        // Given
        const error = new Error('Network error')
        vi.mocked(searchSpots).mockRejectedValueOnce(error)
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When
        await act(async () => {
          await result.current.search('coffee')
        })
        
        // Then
        expect(result.current.spots).toEqual([])
        expect(result.current.loading).toBe(false)
        expect(result.current.error).toBe('Network error')
        expect(result.current.hasMore).toBe(true)
        expect(result.current.query).toBe('coffee')
      })

      it('Then it should handle API errors with custom messages', async () => {
        // Given
        const apiError = {
          response: {
            status: 429,
            data: { error: 'Rate limit exceeded' }
          }
        }
        vi.mocked(searchSpots).mockRejectedValueOnce(apiError)
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When
        await act(async () => {
          await result.current.search('coffee')
        })
        
        // Then
        expect(result.current.error).toBe('Rate limit exceeded')
      })
    })

    describe('When loadMore is called', () => {
      it('Then it should append new spots to existing results', async () => {
        // Given
        const initialSpots = [
          { id: '1', name: 'Cafe A', type: 'cafe' }
        ]
        const additionalSpots = [
          { id: '2', name: 'Cafe B', type: 'cafe' }
        ]
        
        const firstResponse = {
          data: initialSpots,
          total: 2,
          hasMore: true,
          offset: 0,
          limit: 1,
        }
        
        const secondResponse = {
          data: additionalSpots,
          total: 2,
          hasMore: false,
          offset: 1,
          limit: 1,
        }
        
        vi.mocked(searchSpots)
          .mockResolvedValueOnce(firstResponse)
          .mockResolvedValueOnce(secondResponse)
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When - Initial search
        await act(async () => {
          await result.current.search('coffee')
        })
        
        // Then - Should have first set of results
        expect(result.current.spots).toEqual(initialSpots)
        expect(result.current.hasMore).toBe(true)
        
        // When - Load more
        await act(async () => {
          await result.current.loadMore()
        })
        
        // Then - Should have combined results
        expect(result.current.spots).toEqual([...initialSpots, ...additionalSpots])
        expect(result.current.hasMore).toBe(false)
        expect(result.current.total).toBe(2)
      })

      it('Then it should not load more if no more results available', async () => {
        // Given
        const mockResponse = {
          data: [{ id: '1', name: 'Cafe A', type: 'cafe' }],
          total: 1,
          hasMore: false,
          offset: 0,
          limit: 10,
        }
        
        vi.mocked(searchSpots).mockResolvedValueOnce(mockResponse)
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When
        await act(async () => {
          await result.current.search('coffee')
        })
        
        const searchCallCount = vi.mocked(searchSpots).mock.calls.length
        
        await act(async () => {
          await result.current.loadMore()
        })
        
        // Then - Should not make another API call
        expect(vi.mocked(searchSpots).mock.calls.length).toBe(searchCallCount)
        expect(result.current.hasMore).toBe(false)
      })

      it('Then it should handle loadMore errors gracefully', async () => {
        // Given
        const initialResponse = {
          data: [{ id: '1', name: 'Cafe A', type: 'cafe' }],
          total: 2,
          hasMore: true,
          offset: 0,
          limit: 1,
        }
        
        vi.mocked(searchSpots)
          .mockResolvedValueOnce(initialResponse)
          .mockRejectedValueOnce(new Error('Load more failed'))
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When
        await act(async () => {
          await result.current.search('coffee')
        })
        
        await act(async () => {
          await result.current.loadMore()
        })
        
        // Then - Should keep existing spots and set error
        expect(result.current.spots).toHaveLength(1)
        expect(result.current.error).toBe('Load more failed')
        expect(result.current.loading).toBe(false)
      })
    })

    describe('When clearSearch is called', () => {
      it('Then it should reset all search state', async () => {
        // Given
        const mockResponse = {
          data: [{ id: '1', name: 'Cafe A', type: 'cafe' }],
          total: 1,
          hasMore: false,
          offset: 0,
          limit: 10,
        }
        
        vi.mocked(searchSpots).mockResolvedValueOnce(mockResponse)
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When - Perform search first
        await act(async () => {
          await result.current.search('coffee')
        })
        
        // Then - Verify search results exist
        expect(result.current.spots).toHaveLength(1)
        expect(result.current.query).toBe('coffee')
        
        // When - Clear search
        act(() => {
          result.current.clearSearch()
        })
        
        // Then - Should reset to initial state
        expect(result.current.spots).toEqual([])
        expect(result.current.query).toBe('')
        expect(result.current.error).toBeNull()
        expect(result.current.total).toBe(0)
        expect(result.current.hasMore).toBe(true)
      })
    })

    describe('When updateFilters is called', () => {
      it('Then it should update filters and re-search if query exists', async () => {
        // Given
        const mockResponse = {
          data: [{ id: '1', name: 'Solo Cafe', type: 'cafe', soloFriendly: true }],
          total: 1,
          hasMore: false,
          offset: 0,
          limit: 10,
        }
        
        vi.mocked(searchSpots).mockResolvedValue(mockResponse)
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When - Initial search
        await act(async () => {
          await result.current.search('coffee')
        })
        
        // Clear previous calls for assertion
        vi.mocked(searchSpots).mockClear()
        
        // When - Update filters
        await act(async () => {
          await result.current.updateFilters({ soloFriendly: true, type: 'cafe' })
        })
        
        // Then - Should update filters and re-search
        expect(result.current.filters).toEqual({ soloFriendly: true, type: 'cafe' })
        expect(searchSpots).toHaveBeenCalledWith('coffee', {
          offset: 0,
          limit: 10,
          soloFriendly: true,
          type: 'cafe',
        })
      })

      it('Then it should update filters without searching if no query', async () => {
        // Given
        const { result } = renderHook(() => useSpotSearch())
        
        // When
        act(() => {
          result.current.updateFilters({ soloFriendly: true })
        })
        
        // Then
        expect(result.current.filters).toEqual({ soloFriendly: true })
        expect(searchSpots).not.toHaveBeenCalled()
      })
    })

    describe('When multiple searches are performed quickly', () => {
      it('Then it should handle race conditions correctly', async () => {
        // Given
        let resolveFirst: (value: any) => void
        let resolveSecond: (value: any) => void
        
        const firstPromise = new Promise(resolve => { resolveFirst = resolve })
        const secondPromise = new Promise(resolve => { resolveSecond = resolve })
        
        vi.mocked(searchSpots)
          .mockReturnValueOnce(firstPromise)
          .mockReturnValueOnce(secondPromise)
        
        const { result } = renderHook(() => useSpotSearch())
        
        // When - Start two searches quickly
        act(() => {
          result.current.search('first')
        })
        
        act(() => {
          result.current.search('second')
        })
        
        // Resolve second search first (simulating race condition)
        const secondResults = {
          data: [{ id: '2', name: 'Second Result', type: 'cafe' }],
          total: 1,
          hasMore: false,
          offset: 0,
          limit: 10,
        }
        
        await act(async () => {
          resolveSecond!(secondResults)
        })
        
        // Then resolve first search (should be ignored)
        const firstResults = {
          data: [{ id: '1', name: 'First Result', type: 'cafe' }],
          total: 1,
          hasMore: false,
          offset: 0,
          limit: 10,
        }
        
        await act(async () => {
          resolveFirst!(firstResults)
        })
        
        // Then - Should only show results from the latest search
        expect(result.current.spots).toEqual(secondResults.data)
        expect(result.current.query).toBe('second')
      })
    })
  })
})