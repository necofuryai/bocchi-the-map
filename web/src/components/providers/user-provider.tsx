'use client';

import React, { createContext, useContext, ReactNode } from 'react';
import { UserProfile, Auth0Provider, useUser } from '@auth0/nextjs-auth0';

// Define the user context type
interface UserContextType {
  user: UserProfile | undefined;
  error: Error | undefined;
  isLoading: boolean;
}

// Create the user context
const UserContext = createContext<UserContextType | undefined>(undefined);

// Props for the UserProvider component
interface UserProviderProps {
  children: ReactNode;
}

// Internal component that provides user context value
function UserContextProvider({ children }: UserProviderProps) {
  const { user, error, isLoading } = useUser();

  const contextValue: UserContextType = {
    user: user || undefined,
    error: error || undefined,
    isLoading,
  };

  return (
    <UserContext.Provider value={contextValue}>
      {children}
    </UserContext.Provider>
  );
}

// Main UserProvider component that wraps Auth0Provider and UserContextProvider
export function UserProvider({ children }: UserProviderProps) {
  return (
    <Auth0Provider>
      <UserContextProvider>
        {children}
      </UserContextProvider>
    </Auth0Provider>
  );
}

// Custom hook to use the user context
export function useUserContext(): UserContextType {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error('useUserContext must be used within a UserProvider');
  }
  return context;
}

// Export Auth0 hooks for direct use when needed
export { useUser } from '@auth0/nextjs-auth0';