
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
  private accessToken: string | null = null

  constructor(baseURL?: string, accessToken?: string) {
    this.baseURL = baseURL || API_BASE_URL
    this.accessToken = accessToken || null
  }

  // Set access token manually (useful for server-side usage)
  setAccessToken(token: string | null): void {
    this.accessToken = token
  }

  // Get Auth0 access token dynamically
  private async getAccessToken(): Promise<string | null> {
    // If token is already set, use it
    if (this.accessToken) {
      return this.accessToken
    }

    // On client side, try to get token from API route
    if (typeof window !== 'undefined') {
      try {
        const response = await fetch('/api/auth/access-token', {
          method: 'GET',
          credentials: 'include',
        })

        if (response.ok) {
          const data = await response.json()
          return data.accessToken || null
        }
      } catch (error) {
        console.warn('Failed to get access token:', error)
      }
    }

    return null
  }

  // Make authenticated request with automatic token refresh
  async request<T = any>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<APIResponse<T>> {
    // Prepare headers
    const headers = new Headers(options.headers)
    headers.set('Content-Type', 'application/json')

    // Try to get Auth0 access token
    const accessToken = await this.getAccessToken()
    if (accessToken) {
      headers.set('Authorization', `Bearer ${accessToken}`)
    }

    // Make the request with credentials to include cookies (fallback auth)
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

    // Handle 401 Unauthorized - authentication required
    if (response.status === 401) {
      // If we have an access token but still get 401, it might be expired
      if (accessToken) {
        // Clear the cached token and potentially redirect to login
        this.accessToken = null
        
        return {
          error: {
            message: 'Authentication expired',
            details: 'Your session has expired. Please log in again.',
          },
          status: 401,
        }
      }
      
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

// Create an authenticated API client with access token (for server-side usage)
export function createAuthenticatedAPIClient(accessToken?: string): APIClient {
  return new APIClient(undefined, accessToken)
}

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
export async function isAuthenticated(): Promise<boolean> {
  try {
    // Check authentication by calling a protected endpoint
    const result = await apiClient.get('/api/v1/users/me')
    return !result.error && result.status === 200
  } catch (error) {
    // Any error means not authenticated
    return false
  }
}

// Helper function to check if user is authenticated with specific token
export async function isAuthenticatedWithToken(accessToken: string): Promise<boolean> {
  try {
    const authenticatedClient = createAuthenticatedAPIClient(accessToken)
    const result = await authenticatedClient.get('/api/v1/users/me')
    return !result.error && result.status === 200
  } catch (error) {
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