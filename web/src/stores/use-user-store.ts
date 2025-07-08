import { create } from 'zustand'
import { UserProfile } from '@auth0/nextjs-auth0/client'

interface UserState {
  user: UserProfile | undefined
  error: Error | undefined
  isLoading: boolean
  setUser: (user: UserProfile | undefined) => void
  setError: (error: Error | undefined) => void
  setIsLoading: (isLoading: boolean) => void
}

export const useUserStore = create<UserState>((set) => ({
  user: undefined,
  error: undefined,
  isLoading: true,
  setUser: (user) => set({ user }),
  setError: (error) => set({ error }),
  setIsLoading: (isLoading) => set({ isLoading }),
}))