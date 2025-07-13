'use client'

import { useState } from 'react'
import Link from 'next/link'
import { MapPinIcon } from '@heroicons/react/24/outline'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { AuthButton } from '@/components/auth/auth-button'

export default function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center space-y-4">
          <Link href="/" className="inline-flex items-center space-x-2">
            <MapPinIcon className="h-8 w-8 text-primary" />
            <h1 className="text-2xl font-bold">Bocchi The Map</h1>
          </Link>
          <p className="text-muted-foreground">
            おひとりさまスポットを見つけよう
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="text-center">ログイン / 新規登録</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="email" className="text-sm font-medium">
                メールアドレス
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="your@email.com"
                className="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring"
                required
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="text-sm font-medium">
                パスワード
              </label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="パスワードを入力"
                className="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring"
                required
              />
            </div>

            <AuthButton
              email={email}
              password={password}
            />

            <div className="text-center">
              <Link href="/">
                <Button variant="link" size="sm">
                  ホームに戻る
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>

        <div className="text-center text-sm text-muted-foreground">
          新規登録の場合、確認メールが送信されます
        </div>
      </div>
    </div>
  )
}