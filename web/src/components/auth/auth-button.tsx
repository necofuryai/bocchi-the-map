'use client';

import React from 'react';
import { useUser } from '@auth0/nextjs-auth0';
import { Button } from '@/components/ui/button';

interface AuthButtonProps {
  className?: string;
  variant?: 'default' | 'outline' | 'ghost' | 'link' | 'destructive' | 'secondary';
  size?: 'default' | 'sm' | 'lg' | 'icon';
  showFullText?: boolean;
}

export function AuthButton({ 
  className, 
  variant = 'default', 
  size = 'default',
  showFullText = true 
}: AuthButtonProps) {
  const { user, isLoading } = useUser();

  if (isLoading) {
    return (
      <Button 
        disabled 
        className={className} 
        variant={variant} 
        size={size}
      >
        {showFullText ? '読み込み中...' : '...'}
      </Button>
    );
  }

  if (user) {
    return (
      <a href="/auth/logout">
        <Button 
          className={className} 
          variant={variant} 
          size={size}
        >
          {showFullText ? 'ログアウト' : '出'}
        </Button>
      </a>
    );
  }

  return (
    <a href="/auth/login">
      <Button 
        className={className} 
        variant={variant} 
        size={size}
      >
        {showFullText ? 'ログイン' : '入'}
      </Button>
    </a>
  );
}

export default AuthButton;