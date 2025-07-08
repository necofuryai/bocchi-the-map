'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useUserStore } from '@/stores/use-user-store';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

/**
 * Login page component
 * Handles user authentication via Auth0 and redirects authenticated users
 */
export default function LoginPage() {
  const { user, isLoading } = useUserStore();
  const router = useRouter();

  // Redirect authenticated users to home page
  useEffect(() => {
    if (user && !isLoading) {
      router.push('/');
    }
  }, [user, isLoading, router]);

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

  // If user is already authenticated, don't show login form
  if (user) {
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
            <h2 className="text-xl font-semibold">ログイン</h2>
            <p className="text-sm text-muted-foreground">
              アカウントにサインインして、お気に入りのスポットを見つけましょう
            </p>
          </div>

          <a href="/api/auth/login" className="block">
            <Button className="w-full h-12 text-base">
              サインイン
            </Button>
          </a>

          <div className="text-center">
            <p className="text-xs text-muted-foreground">
              サインインすることで、
              <a href="/terms" className="underline hover:text-primary">利用規約</a>
              および
              <a href="/privacy" className="underline hover:text-primary">プライバシーポリシー</a>
              に同意したことになります
            </p>
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