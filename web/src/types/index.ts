// Domain entities - MVP simplified version
export interface Review {
  id: string
  spotId: string
  userId: string
  rating: number // 1-5 stars only
  comment?: string
  createdAt: string
  updatedAt: string
}

// MVP simplified user (basic info only)
export interface DomainUser {
  id: string
  email: string
  name?: string
  createdAt: string
  updatedAt: string
}

// Spot entity - MVP simplified version
export interface Spot {
  id: string
  name: string
  category: string
  address: string
  countryCode: string
  createdAt: string
  updatedAt: string
}

// Common response types
export interface PaginationResponse<T = unknown> {
  items: T[]
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
}

// Note: Map-related types removed for MVP simplification
// Geographic search and POI features have been removed

// Supabase Auth related types
export interface SupabaseUser {
  id: string
  email?: string
  user_metadata?: {
    name?: string
    [key: string]: unknown
  }
  created_at?: string
  updated_at?: string
}

export interface AuthUser {
  id: string
  email: string
  name?: string
}

// Auth page component types
export interface AuthPageProps {
  redirectTo?: string
}

export interface LoginPageState {
  isLoading: boolean
  error?: string
}

