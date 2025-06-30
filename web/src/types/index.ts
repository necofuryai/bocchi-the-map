// Domain entities
export interface Spot {
  id: string
  name: string
  nameI18n?: Record<string, string>
  latitude: number
  longitude: number
  category: string
  address: string
  addressI18n?: Record<string, string>
  countryCode: string
  averageRating: number
  reviewCount: number
  createdAt: string
  updatedAt: string
}

export interface Review {
  id: string
  spotId: string
  userId: string
  rating: number
  comment?: string
  ratingAspects?: Record<string, number>
  createdAt: string
  updatedAt: string
}

export interface DomainUser {
  id: string
  email: string
  displayName: string
  avatarUrl?: string
  preferences: UserPreferences
}

export interface UserPreferences {
  language: 'ja' | 'en'
  darkMode: boolean
  timezone: string
}

// Common response types
export interface PaginationResponse<T = unknown> {
  items: T[]
  totalCount: number
  page: number
  pageSize: number
  totalPages: number
}

// Map-related types
export interface MapPosition {
  latitude: number
  longitude: number
  zoom?: number
  bearing?: number
  pitch?: number
}

export interface POIProperties {
  name?: string
  kind?: string
  script?: string
  min_zoom?: number // snake_case for TileJSON specification compliance
}

// Auth0 related types
export interface Auth0User {
  sub: string
  name?: string
  email?: string
  email_verified?: boolean
  picture?: string
  nickname?: string
  given_name?: string
  family_name?: string
  locale?: string
  updated_at?: string
}

export interface Auth0Session {
  user: Auth0User
  accessToken?: string
  idToken?: string
  refreshToken?: string
}

export interface Auth0Context {
  user?: Auth0User
  isLoading: boolean
  error?: Error
  checkSession: () => Promise<void>
  loginWithRedirect: (options?: any) => Promise<void>
  logout: (options?: any) => void
}

// Auth page component types
export interface AuthPageProps {
  redirectTo?: string
  returnTo?: string
}

export interface LoginPageState {
  isLoading: boolean
  error?: string
}

export interface LogoutPageState {
  isLoggingOut: boolean
  showConfirmation: boolean
}

