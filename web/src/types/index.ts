// Domain entities
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
  updated_at?: string
}

export interface Auth0Session {
  user: Auth0User
  accessToken?: string
  idToken?: string
  refreshToken?: string
}

// Auth0 configuration types
interface Auth0RedirectOptions {
  appState?: Record<string, string | number | boolean>
  fragment?: string
  redirectUri?: string
  screen_hint?: 'signup' | 'login'
  prompt?: 'none' | 'login' | 'consent' | 'select_account'
  max_age?: number
  login_hint?: string
  acr_values?: string
  scope?: string
  audience?: string
  connection?: string
  [key: string]: string | number | boolean | Record<string, string | number | boolean> | undefined
}

interface Auth0LogoutOptions {
  logoutParams?: {
    returnTo?: string
    client_id?: string
    federated?: boolean
  }
  returnTo?: string
}

export interface Auth0Context {
  user?: Auth0User
  isLoading: boolean
  error?: Error
  checkSession: () => Promise<void>
  loginWithRedirect: (options?: Auth0RedirectOptions) => Promise<void>
  logout: (options?: Auth0LogoutOptions) => void
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

