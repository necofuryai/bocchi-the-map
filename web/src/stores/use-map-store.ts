import { create } from 'zustand'
import { Spot } from './use-search-store'

interface MapFilter {
  kinds?: string[]
  enabled?: boolean
  [key: string]: string | boolean | number | string[] | undefined
}

interface MapState {
  mapState: 'loading' | 'loaded' | 'error'
  error: string | null
  filters: MapFilter
  selectedSpot: Spot | null
  
  setMapState: (state: 'loading' | 'loaded' | 'error') => void
  setError: (error: string | null) => void
  setFilters: (filters: MapFilter) => void
  setSelectedSpot: (spot: Spot | null) => void
  reset: () => void
}

const initialState = {
  mapState: 'loading' as const,
  error: null,
  filters: {},
  selectedSpot: null,
}

export const useMapStore = create<MapState>((set) => ({
  ...initialState,
  
  setMapState: (mapState) => set({ mapState }),
  setError: (error) => set({ error }),
  setFilters: (filters) => set({ filters }),
  
  
  setSelectedSpot: (selectedSpot) => set({ selectedSpot }),
  reset: () => set(initialState),
}))