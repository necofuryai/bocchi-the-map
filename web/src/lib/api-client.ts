import type { DomainUser, Review, Spot } from '@/types'

// API base URL configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

// API response types
interface APIError {
  message: string
  details?: Record<string, string | number | boolean> | string | Error
}

interface APIResponse<T = Record<string, never>> {
  data?: T
  error?: APIError
  status: number
}

// API client state
interface APIClientState {
  baseURL: string
}

// Create initial state
const createInitialState = (baseURL?: string): APIClientState => ({
  baseURL: baseURL || API_BASE_URL,
})


// Make authenticated request with Clerk token
const request = async <T = Record<string, never>>(
  state: APIClientState,
  endpoint: string,
  options: RequestInit = {}
): Promise<APIResponse<T>> => {
  // Prepare headers
  const headers = new Headers(options.headers)
  headers.set('Content-Type', 'application/json')

  // Note: For authenticated requests, the token should be passed in options.headers
  // This allows for proper token handling from React components using useAuth hook

  // Make the request
  let response: Response
  try {
    response = await fetch(`${state.baseURL}${endpoint}`, {
      ...options,
      headers,
      credentials: 'include',
    })
  } catch (error) {
    return {
      error: {
        message: 'Network error occurred',
        details: error instanceof Error ? error.message : String(error),
      },
      status: 0,
    }
  }

  // Handle 401 Unauthorized - authentication required
  if (response.status === 401) {
    return {
      error: {
        message: 'Authentication required',
        details: 'Please log in to access this resource',
      },
      status: 401,
    }
  }

  // Handle 403 Forbidden - insufficient permissions
  if (response.status === 403) {
    return {
      error: {
        message: 'Access forbidden',
        details: 'You do not have permission to access this resource',
      },
      status: 403,
    }
  }

  // Parse response
  let data: T | undefined
  let error: APIError | undefined

  try {
    if (response.headers.get('content-type')?.includes('application/json')) {
      const jsonData = await response.json()
      if (response.ok) {
        data = jsonData
      } else {
        error = {
          message: jsonData.message || 'API request failed',
          details: jsonData,
        }
      }
    } else {
      const text = await response.text()
      if (response.ok) {
        data = { text } as T
      } else {
        error = {
          message: text || 'API request failed',
          details: { status: response.status, statusText: response.statusText },
        }
      }
    }
  } catch (parseError) {
    error = {
      message: 'Failed to parse response',
      details: parseError instanceof Error ? parseError.message : String(parseError),
    }
  }

  return {
    data,
    error,
    status: response.status,
  }
}

// API client factory function
export const createAPIClient = (baseURL?: string, token?: string) => {
  const state = createInitialState(baseURL)

  const createHeaders = (additionalHeaders?: HeadersInit) => {
    const headers = new Headers(additionalHeaders)
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
    return headers
  }

  return {
    get: <T = Record<string, never>>(endpoint: string, options?: RequestInit): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, { method: 'GET', ...options, headers: createHeaders(options?.headers) }),
    post: <T = Record<string, never>>(endpoint: string, body?: object, options?: RequestInit): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, {
        method: 'POST',
        body: body ? JSON.stringify(body) : undefined,
        ...options,
        headers: createHeaders(options?.headers),
      }),
    put: <T = Record<string, never>>(endpoint: string, body?: object, options?: RequestInit): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, {
        method: 'PUT',
        body: body ? JSON.stringify(body) : undefined,
        ...options,
        headers: createHeaders(options?.headers),
      }),
    patch: <T = Record<string, never>>(endpoint: string, body?: object, options?: RequestInit): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, {
        method: 'PATCH',
        body: body ? JSON.stringify(body) : undefined,
        ...options,
        headers: createHeaders(options?.headers),
      }),
    delete: <T = Record<string, never>>(endpoint: string, options?: RequestInit): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, { method: 'DELETE', ...options, headers: createHeaders(options?.headers) }),
  }
}

// Default API client instance
export const apiClient = createAPIClient()

// Convenience functions for common API operations
export const api = {
  // User operations
  users: {
    getCurrent: () => apiClient.get<DomainUser>('/api/v1/users/me'),
  },

  // Spot operations
  spots: {
    list: (params?: { category?: string; country_code?: string; limit?: number }) => {
      const searchParams = new URLSearchParams()
      if (params?.category) searchParams.set('category', params.category)
      if (params?.country_code) searchParams.set('country_code', params.country_code)
      if (params?.limit) searchParams.set('limit', params.limit.toString())
      
      const query = searchParams.toString()
      return apiClient.get<Spot[]>(`/api/v1/spots${query ? `?${query}` : ''}`)
    },
    create: (spot: Omit<Spot, 'id' | 'createdAt' | 'updatedAt'>) => 
      apiClient.post<Spot>('/api/v1/spots', spot),
    getById: (id: string) => apiClient.get<Spot>(`/api/v1/spots/${id}`),
    update: (id: string, spot: Partial<Omit<Spot, 'id' | 'createdAt' | 'updatedAt'>>) => 
      apiClient.put<Spot>(`/api/v1/spots/${id}`, spot),
    delete: (id: string) => apiClient.delete<Record<string, never>>(`/api/v1/spots/${id}`),
  },

  // Review operations
  reviews: {
    list: (params: { spot_id: string; page?: number; limit?: number }) => {
      const searchParams = new URLSearchParams()
      if (params.page) searchParams.set('page', params.page.toString())
      if (params.limit) searchParams.set('limit', params.limit.toString())
      
      const query = searchParams.toString()
      return apiClient.get<{
        reviews: Review[]
        pagination: {
          page: number
          pageSize: number
          totalCount: number
          totalPages: number
        }
        statistics: {
          averageRating: number
          totalReviews: number
        }
      }>(`/api/v1/spots/${params.spot_id}/reviews${query ? `?${query}` : ''}`)
    },
    create: (review: { spot_id: string; rating: number; comment?: string }) => 
      apiClient.post<Review>('/api/v1/reviews', review),
    getById: (id: string) => apiClient.get<Review>(`/api/v1/reviews/${id}`),
    update: (id: string, review: Partial<Omit<Review, 'id' | 'createdAt' | 'updatedAt'>>) => 
      apiClient.put<Review>(`/api/v1/reviews/${id}`, review),
    delete: (id: string) => apiClient.delete<Record<string, never>>(`/api/v1/reviews/${id}`),
  },
}

// Helper function to check if user is authenticated
export async function isAuthenticated(): Promise<boolean> {
  // With Clerk, authentication state is managed by the middleware
  // We'll need to check the auth state differently
  // For now, return false as we can't directly check from client
  return false
}

// Helper function to handle API errors consistently
export function handleAPIError(error: APIError, fallbackMessage = 'An error occurred'): string {
  if (process.env.NODE_ENV === 'development') {
    console.error('API Error:', error)
  }
  
  return error.message || fallbackMessage
}

// Helper function to handle authentication errors specifically
export function handleAuthError(error: APIError): { shouldRedirectToLogin: boolean; message: string } {
  const isExpired = error.message?.includes('expired')
  const isUnauthorized = error.message?.includes('Authentication required') || 
                         error.message?.includes('Authentication expired')
  
  return {
    shouldRedirectToLogin: isUnauthorized || isExpired,
    message: error.message || 'Authentication error occurred'
  }
}

export default apiClient