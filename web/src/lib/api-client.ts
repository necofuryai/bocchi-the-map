import { refreshAPIToken, clearAPITokens } from './auth'

// API base URL configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

// API response types
interface APIError {
  message: string
  details?: any
}

interface APIResponse<T = any> {
  data?: T
  error?: APIError
  status: number
}

// Generic API client class with automatic authentication
export class APIClient {
  private baseURL: string

  constructor(baseURL?: string) {
    this.baseURL = baseURL || API_BASE_URL
  }

  // Make authenticated request with automatic token refresh
  async request<T = any>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<APIResponse<T>> {
    // Prepare headers
    const headers = new Headers(options.headers)
    headers.set('Content-Type', 'application/json')

    // Authentication handled by HttpOnly cookies

    // Make the request with credentials to include cookies
    let response: Response
    try {
      response = await fetch(`${this.baseURL}${endpoint}`, {
        ...options,
        headers,
        credentials: 'include', // Include cookies for authentication
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

    // Handle 401 Unauthorized - try to refresh token
    if (response.status === 401) {
      const refreshSuccess = await refreshAPIToken()
      if (refreshSuccess) {
        // Retry the request with refreshed HttpOnly cookies
        try {
          response = await fetch(`${this.baseURL}${endpoint}`, {
            ...options,
            headers,
            credentials: 'include', // Include cookies for authentication
          })
        } catch (error) {
          return {
            error: {
              message: 'Network error on retry',
              details: error instanceof Error ? error.message : String(error),
            },
            status: 0,
          }
        }
      } else {
        // Refresh failed, clear tokens
        clearAPITokens()
        return {
          error: {
            message: 'Authentication failed',
            details: 'Token refresh failed',
          },
          status: 401,
        }
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
          data = text as unknown as T
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


  // Convenience methods for different HTTP verbs
  async get<T = any>(endpoint: string): Promise<APIResponse<T>> {
    return this.request<T>(endpoint, { method: 'GET' })
  }

  async post<T = any>(endpoint: string, body?: any): Promise<APIResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  async put<T = any>(endpoint: string, body?: any): Promise<APIResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  async patch<T = any>(endpoint: string, body?: any): Promise<APIResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  async delete<T = any>(endpoint: string): Promise<APIResponse<T>> {
    return this.request<T>(endpoint, { method: 'DELETE' })
  }
}

// Default API client instance
export const apiClient = new APIClient()

// Convenience functions for common API operations
export const api = {
  // User operations
  users: {
    getCurrent: () => apiClient.get('/api/v1/users/me'),
    updatePreferences: (preferences: any) =>
      apiClient.patch('/api/v1/users/me/preferences', { preferences }),
  },

  // Spot operations
  spots: {
    list: (params?: { latitude?: number; longitude?: number; radius?: number; limit?: number }) => {
      const searchParams = new URLSearchParams()
      if (params?.latitude) searchParams.set('latitude', params.latitude.toString())
      if (params?.longitude) searchParams.set('longitude', params.longitude.toString())
      if (params?.radius) searchParams.set('radius', params.radius.toString())
      if (params?.limit) searchParams.set('limit', params.limit.toString())
      
      const query = searchParams.toString()
      return apiClient.get(`/api/v1/spots${query ? `?${query}` : ''}`)
    },
    create: (spot: any) => apiClient.post('/api/v1/spots', spot),
    getById: (id: string) => apiClient.get(`/api/v1/spots/${id}`),
    update: (id: string, spot: any) => apiClient.put(`/api/v1/spots/${id}`, spot),
    delete: (id: string) => apiClient.delete(`/api/v1/spots/${id}`),
  },

  // Review operations
  reviews: {
    list: (params?: { spot_id?: string; user_id?: string; limit?: number }) => {
      const searchParams = new URLSearchParams()
      if (params?.spot_id) searchParams.set('spot_id', params.spot_id)
      if (params?.user_id) searchParams.set('user_id', params.user_id)
      if (params?.limit) searchParams.set('limit', params.limit.toString())
      
      const query = searchParams.toString()
      return apiClient.get(`/api/v1/reviews${query ? `?${query}` : ''}`)
    },
    create: (review: any) => apiClient.post('/api/v1/reviews', review),
    getById: (id: string) => apiClient.get(`/api/v1/reviews/${id}`),
    update: (id: string, review: any) => apiClient.put(`/api/v1/reviews/${id}`, review),
    delete: (id: string) => apiClient.delete(`/api/v1/reviews/${id}`),
  },
}

// Helper function to check if user is authenticated
export function isAuthenticated(): boolean {
  // Authentication is now handled by HttpOnly cookies
  // This would need to be determined by making an API call or checking session state
  return true // This should be implemented based on your authentication strategy
}

// Helper function to handle API errors consistently
export function handleAPIError(error: APIError, fallbackMessage = 'An error occurred'): string {
  if (process.env.NODE_ENV === 'development') {
    console.error('API Error:', error)
  }
  
  return error.message || fallbackMessage
}

export default apiClient