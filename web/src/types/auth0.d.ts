/**
 * Auth0 TypeScript type declarations for Bocchi The Map
 * 
 * This file extends the Auth0 types with custom properties
 * and ensures proper TypeScript support for Auth0 integration.
 */

declare module '@auth0/nextjs-auth0' {
  interface UserProfile {
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

  interface Session {
    user: UserProfile
    accessToken?: string
    idToken?: string
    refreshToken?: string
  }
}

// Extend global environment variables
declare global {
  namespace NodeJS {
    interface ProcessEnv {
      AUTH0_SECRET: string
      AUTH0_BASE_URL: string
      AUTH0_DOMAIN: string
      AUTH0_CLIENT_ID: string
      AUTH0_CLIENT_SECRET: string
      AUTH0_AUDIENCE?: string
      AUTH0_SCOPE?: string
    }
  }
}

export {};