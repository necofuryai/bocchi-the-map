'use client';

import React, { ReactNode, useEffect } from 'react';
import { Auth0Provider, useUser } from '@auth0/nextjs-auth0';
import { useUserStore } from '@/stores/use-user-store';

// Props for the UserProvider component
interface UserProviderProps {
  children: ReactNode;
}

// Internal component that syncs Auth0 state to Zustand store
function UserStoreSync({ children }: UserProviderProps) {
  const { user, error, isLoading } = useUser();
  const { setUser, setError, setIsLoading } = useUserStore();

  useEffect(() => {
    setUser(user || undefined);
    setError(error || undefined);
    setIsLoading(isLoading);
  }, [user, error, isLoading, setUser, setError, setIsLoading]);

  return <>{children}</>;
}

// Main UserProvider component that wraps Auth0Provider and UserStoreSync
export function UserProvider({ children }: UserProviderProps) {
  return (
    <Auth0Provider>
      <UserStoreSync>
        {children}
      </UserStoreSync>
    </Auth0Provider>
  );
}

// Export Auth0 hooks for direct use when needed
export { useUser } from '@auth0/nextjs-auth0';