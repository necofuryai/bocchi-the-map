import NextAuth from "next-auth"
import Google from "next-auth/providers/google"
import Twitter from "next-auth/providers/twitter"

export const { handlers, auth, signIn, signOut } = NextAuth({
  providers: [
    Google({
      clientId: process.env.GOOGLE_CLIENT_ID || (() => { throw new Error("GOOGLE_CLIENT_ID is required") })(),
      clientSecret: process.env.GOOGLE_CLIENT_SECRET || (() => { throw new Error("GOOGLE_CLIENT_SECRET is required") })(),
    }),
    Twitter({
      clientId: process.env.TWITTER_CLIENT_ID || (() => { throw new Error("TWITTER_CLIENT_ID is required") })(),
      clientSecret: process.env.TWITTER_CLIENT_SECRET || (() => { throw new Error("TWITTER_CLIENT_SECRET is required") })(),
      version: "2.0",
    }),
  ],
  callbacks: {
    async signIn({ user, account, profile: _profile }) {
      if (account?.provider === 'google') {
        try {
          // Check if user.email is null/undefined before making API call
          if (!user.email) {
            console.warn('User email is null/undefined, skipping API call and continuing sign-in')
            return true
          }
          
          const apiUrl = process.env.API_URL || 'http://localhost:8080'
          const userData = {
            email: user.email,
            display_name: user.name,
            avatar_url: user.image,
            auth_provider: 'google',
            auth_provider_id: account.providerAccountId,
          }
          
          if (process.env.NODE_ENV === 'development') {
            console.log('Creating/updating user:', userData)
          }
          
          const abortController = new AbortController()
          const timeoutId = setTimeout(() => abortController.abort(), 15000)
          
          const response = await fetch(`${apiUrl}/api/users`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData),
            signal: abortController.signal,
          })
          
          clearTimeout(timeoutId)
          
          if (!response.ok) {
            const errorText = await response.text()
            console.error('Failed to create/update user:', {
              status: response.status,
              statusText: response.statusText,
              error: errorText,
              user: user.email,
            })
            // Allow sign-in to continue even if user creation fails
            // The user will be created on next successful API call
          } else {
            if (process.env.NODE_ENV === 'development') {
              console.log('User created/updated successfully:', user.email)
            }
          }
        } catch (error) {
          console.error('Error creating/updating user:', {
            error: error instanceof Error ? error.message : String(error),
            stack: error instanceof Error ? error.stack : undefined,
            user: user.email,
            provider: account.provider,
          })
          // Allow sign-in to continue even if user creation fails
        }
      }
      return true
    },
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