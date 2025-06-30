'use client';

import React, { ReactNode } from 'react';
import { useUser } from '@auth0/nextjs-auth0';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';

interface AuthGuardProps {
  children: ReactNode;
  fallback?: ReactNode;
  loadingComponent?: ReactNode;
  errorComponent?: (error: Error) => ReactNode;
}

export function AuthGuard({
  children,
  fallback,
  loadingComponent,
  errorComponent,
}: AuthGuardProps) {
  const { user, isLoading, error } = useUser();

  // Show loading state
  if (isLoading) {
    if (loadingComponent) {
      return <>{loadingComponent}</>;
    }

    return (
      <Card className="p-8 text-center">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/3 mx-auto mb-4"></div>
          <div className="h-4 bg-gray-200 rounded w-1/2 mx-auto"></div>
        </div>
      </Card>
    );
  }

  // Show error state
  if (error) {
    if (errorComponent) {
      return <>{errorComponent(error)}</>;
    }

    return (
      <Card className="p-8 text-center">
        <div className="text-red-600">
          <h2 className="text-xl font-semibold mb-2">Authentication Error</h2>
          <p className="mb-4">{error.message}</p>
          <a href="/auth/login">
            <Button>Try Again</Button>
          </a>
        </div>
      </Card>
    );
  }

  // Show unauthenticated state
  if (!user) {
    if (fallback) {
      return <>{fallback}</>;
    }

    return (
      <Card className="p-8 text-center">
        <h2 className="text-xl font-semibold mb-2">Sign In Required</h2>
        <p className="text-gray-600 mb-4">
          You need to sign in to access this content.
        </p>
        <a href="/auth/login">
          <Button>Sign In</Button>
        </a>
      </Card>
    );
  }

  // User is authenticated, render children
  return <>{children}</>;
}

export default AuthGuard;