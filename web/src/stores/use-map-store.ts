import { create } from 'zustand'
import type { Spot, Review } from '@/types'

interface POIPopupState {
  isOpen: boolean
  spot: Spot | null
  reviews: Review[]
  averageRating: number | null
  position: [number, number] | null
  isLoadingReviews: boolean
  reviewsError: string | null
}

interface MapState {
  mapState: 'loading' | 'loaded' | 'error'
  error: string | null
  popup: POIPopupState
  
  setMapState: (state: 'loading' | 'loaded' | 'error') => void
  setError: (error: string | null) => void
  
  // Popup actions
  openPopup: (spot: Spot, position: [number, number]) => void
  closePopup: () => void
  setPopupReviews: (reviews: Review[]) => void
  setPopupReviewsLoading: (loading: boolean) => void
  setPopupReviewsError: (error: string | null) => void
  
  reset: () => void
}

const initialPopupState: POIPopupState = {
  isOpen: false,
  spot: null,
  reviews: [],
  averageRating: null,
  position: null,
  isLoadingReviews: false,
  reviewsError: null,
}

const initialState = {
  mapState: 'loading' as const,
  error: null,
  popup: initialPopupState,
}

// Helper function to calculate average rating
const calculateAverageRating = (reviews: Review[]): number | null => {
  if (reviews.length === 0) return null
  const sum = reviews.reduce((acc, review) => acc + review.rating, 0)
  return Math.round((sum / reviews.length) * 10) / 10
}

export const useMapStore = create<MapState>((set) => ({
  ...initialState,
  
  setMapState: (mapState) => set({ mapState }),
  setError: (error) => set({ error }),
  
  // Popup actions
  openPopup: (spot, position) => set((state) => ({
    popup: {
      ...state.popup,
      isOpen: true,
      spot,
      position,
      reviews: [],
      averageRating: null,
      isLoadingReviews: false,
      reviewsError: null,
    }
  })),
  
  closePopup: () => set((state) => ({
    popup: {
      ...state.popup,
      isOpen: false,
      spot: null,
      reviews: [],
      averageRating: null,
      position: null,
      isLoadingReviews: false,
      reviewsError: null,
    }
  })),
  
  setPopupReviews: (reviews) => set((state) => ({
    popup: {
      ...state.popup,
      reviews,
      averageRating: calculateAverageRating(reviews),
      isLoadingReviews: false,
      reviewsError: null,
    }
  })),
  
  setPopupReviewsLoading: (isLoadingReviews) => set((state) => ({
    popup: {
      ...state.popup,
      isLoadingReviews,
      reviewsError: null,
    }
  })),
  
  setPopupReviewsError: (reviewsError) => set((state) => ({
    popup: {
      ...state.popup,
      reviewsError,
      isLoadingReviews: false,
    }
  })),
  
  reset: () => set(initialState),
}))