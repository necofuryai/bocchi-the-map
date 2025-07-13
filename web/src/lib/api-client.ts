import { createClient } from '@/utils/supabase/client'
import type { DomainUser, Review } from '@/types'

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

// Get Supabase access token
const getAccessToken = async (): Promise<string | null> => {
  const supabase = createClient()
  
  try {
    const { data: { session } } = await supabase.auth.getSession()
    return session?.access_token ?? null
  } catch (error) {
    console.warn('Failed to get access token:', error)
    return null
  }
}

// Make authenticated request with automatic token refresh
const request = async <T = Record<string, never>>(
  state: APIClientState,
  endpoint: string,
  options: RequestInit = {}
): Promise<APIResponse<T>> => {
  // Prepare headers
  const headers = new Headers(options.headers)
  headers.set('Content-Type', 'application/json')

  // Try to get Supabase access token
  const accessToken = await getAccessToken()
  if (accessToken) {
    headers.set('Authorization', `Bearer ${accessToken}`)
  }

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
export const createAPIClient = (baseURL?: string) => {
  const state = createInitialState(baseURL)

  return {
    get: <T = Record<string, never>>(endpoint: string): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, { method: 'GET' }),
    post: <T = Record<string, never>>(endpoint: string, body?: object): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, {
        method: 'POST',
        body: body ? JSON.stringify(body) : undefined,
      }),
    put: <T = Record<string, never>>(endpoint: string, body?: object): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, {
        method: 'PUT',
        body: body ? JSON.stringify(body) : undefined,
      }),
    patch: <T = Record<string, never>>(endpoint: string, body?: object): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, {
        method: 'PATCH',
        body: body ? JSON.stringify(body) : undefined,
      }),
    delete: <T = Record<string, never>>(endpoint: string): Promise<APIResponse<T>> => 
      request<T>(state, endpoint, { method: 'DELETE' }),
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

  // Review operations
  reviews: {
    list: (params?: { spot_id?: string; user_id?: string; limit?: number }) => {
      const searchParams = new URLSearchParams()
      if (params?.spot_id) searchParams.set('spot_id', params.spot_id)
      if (params?.user_id) searchParams.set('user_id', params.user_id)
      if (params?.limit) searchParams.set('limit', params.limit.toString())
      
      const query = searchParams.toString()
      return apiClient.get<Review[]>(`/api/v1/reviews${query ? `?${query}` : ''}`)
    },
    create: (review: Omit<Review, 'id' | 'createdAt' | 'updatedAt'>) => 
      apiClient.post<Review>('/api/v1/reviews', review),
    getById: (id: string) => apiClient.get<Review>(`/api/v1/reviews/${id}`),
    update: (id: string, review: Partial<Omit<Review, 'id' | 'createdAt' | 'updatedAt'>>) => 
      apiClient.put<Review>(`/api/v1/reviews/${id}`, review),
    delete: (id: string) => apiClient.delete<Record<string, never>>(`/api/v1/reviews/${id}`),
  },
}

// Helper function to check if user is authenticated
export async function isAuthenticated(): Promise<boolean> {
  const supabase = createClient()
  try {
    const { data: { session } } = await supabase.auth.getSession()
    return !!session?.user
  } catch {
    return false
  }
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