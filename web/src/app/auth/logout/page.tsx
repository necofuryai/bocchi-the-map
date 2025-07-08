'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useUserStore } from '@/stores/use-user-store';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

/**
 * Logout page component
 * Handles user logout confirmation and Auth0 logout process
 */
export default function LogoutPage() {
  const { user, isLoading } = useUserStore();
  const router = useRouter();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  // Redirect non-authenticated users to home page
  useEffect(() => {
    if (!user && !isLoading) {
      router.push('/');
    }
  }, [user, isLoading, router]);

  const handleLogout = async () => {
    setIsLoggingOut(true);
    // Navigate to Auth0 logout endpoint
    window.location.href = '/api/auth/logout';
  };

  const handleCancel = () => {
    router.push('/');
  };

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card className="w-full max-w-md p-8">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">認証状態を確認中...</p>
          </div>
        </Card>
      </div>
    );
  }

  // If user is not authenticated, don't show logout form
  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <Card className="w-full max-w-md p-8 space-y-6">
        <div className="text-center space-y-4">
          <h1 className="text-2xl font-bold tracking-tight">
            Bocchi The Map
          </h1>
          <p className="text-muted-foreground">
            おひとりさま向けスポットレビューアプリ
          </p>
        </div>

        <div className="space-y-4">
          <div className="text-center space-y-2">
            <h2 className="text-xl font-semibold">ログアウト</h2>
            <div className="space-y-3">
              {user.picture && (
                <img
                  src={user.picture}
                  alt="プロフィール画像"
                  className="w-16 h-16 rounded-full mx-auto border-2 border-border"
                />
              )}
              <div>
                <p className="font-medium">{user.name}</p>
                <p className="text-sm text-muted-foreground">{user.email}</p>
              </div>
            </div>
            <p className="text-sm text-muted-foreground pt-2">
              本当にログアウトしますか？
            </p>
          </div>

          <div className="space-y-3">
            <Button 
              onClick={handleLogout}
              disabled={isLoggingOut}
              variant="destructive"
              className="w-full h-12 text-base"
            >
              {isLoggingOut ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  ログアウト中...
                </>
              ) : (
                'ログアウト'
              )}
            </Button>

            <Button 
              onClick={handleCancel}
              variant="outline"
              disabled={isLoggingOut}
              className="w-full h-12 text-base"
            >
              キャンセル
            </Button>
          </div>
        </div>

        <div className="text-center pt-4 border-t">
          <a 
            href="/" 
            className="text-sm text-muted-foreground hover:text-primary transition-colors"
          >
            ← ホームに戻る
          </a>
        </div>
      </Card>
    </div>
  );
}