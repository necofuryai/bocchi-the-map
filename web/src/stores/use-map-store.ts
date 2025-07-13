import { create } from 'zustand'


interface MapState {
  mapState: 'loading' | 'loaded' | 'error'
  error: string | null
  
  setMapState: (state: 'loading' | 'loaded' | 'error') => void
  setError: (error: string | null) => void
  reset: () => void
}

const initialState = {
  mapState: 'loading' as const,
  error: null,
}

export const useMapStore = create<MapState>((set) => ({
  ...initialState,
  
  setMapState: (mapState) => set({ mapState }),
  setError: (error) => set({ error }),
  reset: () => set(initialState),
}))