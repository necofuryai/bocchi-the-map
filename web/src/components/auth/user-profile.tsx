'use client';

import React from 'react';
import { useUser } from '@auth0/nextjs-auth0';
import { Card } from '@/components/ui/card';

interface UserProfileProps {
  className?: string;
}

export function UserProfile({ className }: UserProfileProps) {
  const { user, isLoading, error } = useUser();

  if (isLoading) {
    return (
      <Card className={className}>
        <div className="p-4">
          <div className="animate-pulse">
            <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div className="h-4 bg-gray-200 rounded w-1/2"></div>
          </div>
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className={className}>
        <div className="p-4 text-red-600">
          <p>Error loading user profile: {error.message}</p>
        </div>
      </Card>
    );
  }

  if (!user) {
    return (
      <Card className={className}>
        <div className="p-4">
          <p className="text-gray-500">Please sign in to view your profile.</p>
        </div>
      </Card>
    );
  }

  return (
    <Card className={className}>
      <div className="p-4">
        <div className="flex items-center space-x-4">
          {user.picture && (
            <img
              src={user.picture}
              alt={user.name || 'User'}
              className="h-10 w-10 rounded-full"
            />
          )}
          <div>
            <h3 className="text-lg font-semibold">
              {user.name || 'Anonymous User'}
            </h3>
            {user.email && (
              <p className="text-sm text-gray-600">{user.email}</p>
            )}
          </div>
        </div>
        {user.email_verified !== undefined && (
          <div className="mt-3">
            <span
              className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                user.email_verified
                  ? 'bg-green-100 text-green-800'
                  : 'bg-yellow-100 text-yellow-800'
              }`}
            >
              {user.email_verified ? 'Verified' : 'Unverified'}
            </span>
          </div>
        )}
      </div>
    </Card>
  );
}

export default UserProfile;