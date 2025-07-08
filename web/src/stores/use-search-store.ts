import { create } from 'zustand'

export interface Spot {
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

export interface SearchFilters {
  soloFriendly?: boolean
  type?: string
  [key: string]: string | boolean | number | undefined
}

interface SearchState {
  spots: Spot[]
  loading: boolean
  error: string | null
  hasMore: boolean
  query: string
  filters: SearchFilters
  total: number
  recentSearches: string[]
  
  setSpots: (spots: Spot[]) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  setHasMore: (hasMore: boolean) => void
  setQuery: (query: string) => void
  setFilters: (filters: SearchFilters) => void
  setTotal: (total: number) => void
  addRecentSearch: (query: string) => void
  clearRecentSearches: () => void
  reset: () => void
}

const initialState = {
  spots: [],
  loading: false,
  error: null,
  hasMore: true,
  query: '',
  filters: {} as SearchFilters,
  total: 0,
  recentSearches: [],
}

export const useSearchStore = create<SearchState>((set, get) => ({
  ...initialState,
  
  setSpots: (spots) => set({ spots }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  setHasMore: (hasMore) => set({ hasMore }),
  setQuery: (query) => set({ query }),
  setFilters: (filters) => set({ filters }),
  setTotal: (total) => set({ total }),
  
  addRecentSearch: (query) => {
    const { recentSearches } = get()
    const filtered = recentSearches.filter(s => s !== query)
    const updated = [query, ...filtered].slice(0, 10) // Keep last 10 searches
    set({ recentSearches: updated })
  },
  
  clearRecentSearches: () => set({ recentSearches: [] }),
  reset: () => set(initialState),
}))