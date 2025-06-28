import NextAuth from "next-auth"
import Google from "next-auth/providers/google"
import Twitter from "next-auth/providers/twitter"

// Auth provider constants
const AUTH_PROVIDER = {
  GOOGLE: 'google' as const,
  TWITTER: 'twitter' as const,
} as const

// Type augmentation for Auth.js v5
declare module "next-auth" {
  interface User {
    provider?: string
    providerAccountId?: string
  }
  
  interface Session {
    user: {
      id?: string
      email?: string
      name?: string
      image?: string
      provider?: string
      providerAccountId?: string
    }
  }
}

export const { handlers, auth, signIn, signOut } = NextAuth({
  providers: [
    Google({
      clientId: process.env.GOOGLE_CLIENT_ID || (() => { throw new Error("GOOGLE_CLIENT_ID is required") })(),
      clientSecret: process.env.GOOGLE_CLIENT_SECRET || (() => { throw new Error("GOOGLE_CLIENT_SECRET is required") })(),
    }),
    Twitter({
      clientId: process.env.TWITTER_CLIENT_ID || (() => { throw new Error("TWITTER_CLIENT_ID is required") })(),
      clientSecret: process.env.TWITTER_CLIENT_SECRET || (() => { throw new Error("TWITTER_CLIENT_SECRET is required") })(),
    }),
  ],
  callbacks: {
    // Handle sign-in callback with user creation/update
    async signIn({ user, account }) {
      if (account && (account.provider === AUTH_PROVIDER.GOOGLE || account.provider === AUTH_PROVIDER.TWITTER)) {
        try {
          // Check if user.email is null/undefined before making API call
          if (!user.email) {
            console.warn('User email is null/undefined, skipping API call and continuing sign-in')
            return true
          }
          
          const apiUrl = process.env.API_URL || 'http://localhost:8080'
          const userData = {
            email: user.email,
            name: user.name,
            image: user.image,
            provider: account.provider,
            provider_id: account.providerAccountId,
          }
          
          if (process.env.NODE_ENV === 'development') {
            console.log('Creating/updating user:', userData)
          }
          
          const abortController = new AbortController()
          const timeoutId = setTimeout(() => abortController.abort(), 15000)
          
          let response: Response
          try {
            response = await fetch(`${apiUrl}/api/v1/users`, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify(userData),
              signal: abortController.signal,
            })
          } finally {
            clearTimeout(timeoutId)
          }
          
          if (!response.ok) {
            const errorText = await response.text()
            if (process.env.NODE_ENV !== 'production') {
              console.error(JSON.stringify({
                timestamp: new Date().toISOString(),
                level: 'error',
                message: 'Failed to create/update user',
                details: {
                  status: response.status,
                  statusText: response.statusText,
                  error: errorText,
                  user: user.email,
                }
              }))
              // Log specific error for debugging
              if (response.status >= 500) {
                console.error(JSON.stringify({
                  timestamp: new Date().toISOString(),
                  level: 'error',
                  message: 'Server error - user creation will be retried on next login',
                  statusCode: response.status
                }))
              } else if (response.status === 400) {
                console.error(JSON.stringify({
                  timestamp: new Date().toISOString(),
                  level: 'error',
                  message: 'Invalid request data - check OAuth provider configuration',
                  statusCode: response.status
                }))
              }
            }
            // Allow sign-in to continue even if user creation fails
            // The user will be created on next successful API call
          } else {
            if (process.env.NODE_ENV === 'development') {
              console.log('User created/updated successfully:', user.email)
            }
            
            // After successful user creation, generate API token
            await generateAPIToken(userData, apiUrl)
          }
        } catch (error) {
          if (process.env.NODE_ENV !== 'production') {
            console.error('Error creating/updating user:', {
              message: error instanceof Error ? error.message : String(error),
              stack: error instanceof Error ? error.stack : undefined,
              name: error instanceof Error ? error.name : 'Unknown',
              cause: error instanceof Error ? error.cause : undefined,
              user: user.email,
              provider: account.provider,
              // Network error details if available (safely typed)
              networkCode: (error && typeof error === 'object' && 'code' in error) ? (error as Record<string, unknown>).code : undefined,
              httpStatus: (error && typeof error === 'object' && 'status' in error) ? (error as Record<string, unknown>).status : undefined,
              responseData: (error && typeof error === 'object' && 'response' in error) ? (error as Record<string, unknown>).response : undefined,
              timestamp: new Date().toISOString(),
            })
          }
          // Allow sign-in to continue even if user creation fails
        }
      }
      return true
    },
    // Customize session object with additional user data
    async session({ session, token }) {
      if (session?.user) {
        const userId = token?.uid ?? token?.sub
        if (typeof userId === 'string') {
          session.user.id = userId
        }
        if (typeof token?.provider === 'string') {
          session.user.provider = token.provider
        }
        if (typeof token?.providerAccountId === 'string') {
          session.user.providerAccountId = token.providerAccountId
        }
      }
      return session
    },
    // Customize JWT token with additional data
    async jwt({ token, user, account }) {
      if (user) {
        token.uid = user.id
      }
      if (account) {
        token.provider = account.provider
        token.providerAccountId = account.providerAccountId
      }
      return token
    },
  },
  pages: {
    signIn: "/auth/signin",
    error: "/auth/error",
  },
  session: {
    strategy: "jwt",
  },
})

// Generate API token for authenticated user
async function generateAPIToken(userData: {
  email: string
  name?: string | null
  image?: string | null
  provider: string
  provider_id: string
}, apiUrl: string): Promise<void> {
  try {
    const abortController = new AbortController()
    const timeoutId = setTimeout(() => abortController.abort(), 15000)
    
    let response: Response
    try {
      response = await fetch(`${apiUrl}/api/v1/auth/token`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: userData.email,
          provider: userData.provider,
          provider_id: userData.provider_id,
        }),
        credentials: 'include',
        signal: abortController.signal,
      })
    } finally {
      clearTimeout(timeoutId)
    }
    
    if (response.ok) {
      if (process.env.NODE_ENV === 'development') {
        console.log('API token authentication successful')
      }
    } else {
      const errorText = await response.text()
      if (process.env.NODE_ENV !== 'production') {
        console.error('Failed to generate API token:', {
          status: response.status,
          statusText: response.statusText,
          error: errorText,
          email: userData.email,
        })
      }
    }
  } catch (error) {
    if (process.env.NODE_ENV !== 'production') {
      console.error('Error generating API token:', {
        message: error instanceof Error ? error.message : String(error),
        email: userData.email,
      })
    }
  }
}

// Clear API tokens (logout)
export function clearAPITokens(): void {
  // Token clearing now handled by server-side logout endpoint
  // HttpOnly cookies will be cleared by the server
  if (process.env.NODE_ENV === 'development') {
    console.log('Token clearing handled by server-side logout')
  }
}

// Refresh API token using HttpOnly cookies
export async function refreshAPIToken(): Promise<boolean> {
  try {
    const apiUrl = process.env.API_URL || 'http://localhost:8080'
    
    const response = await fetch(`${apiUrl}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    })
    
    if (response.ok) {
      if (process.env.NODE_ENV === 'development') {
        console.log('API token refresh successful')
      }
      return true
    } else {
      // If refresh fails, clear tokens
      clearAPITokens()
      return false
    }
  } catch (error) {
    if (process.env.NODE_ENV !== 'production') {
      console.error('Error refreshing API token:', error)
    }
    clearAPITokens()
    return false
  }
}