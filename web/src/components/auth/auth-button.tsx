'use client'

import { useState } from 'react'
import { createClient } from '@/utils/supabase/client'
import { Button } from '@/components/ui/button'

interface AuthButtonProps {
  email: string
  password: string
  onEmailChange?: (email: string) => void
  onPasswordChange?: (password: string) => void
  className?: string
}

export function AuthButton({ 
  email, 
  password, 
  className 
}: AuthButtonProps) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const supabase = createClient()

  const handleSignIn = async () => {
    if (!email || !password) {
      setError('メールアドレスとパスワードを入力してください')
      return
    }

    setLoading(true)
    setError(null)

    try {
      const { error } = await supabase.auth.signInWithPassword({
        email,
        password,
      })

      if (error) {
        setError(error.message)
      } else {
        // Redirect will be handled by middleware
        window.location.href = '/'
      }
    } catch {
      setError('ログインに失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const handleSignUp = async () => {
    if (!email || !password) {
      setError('メールアドレスとパスワードを入力してください')
      return
    }

    setLoading(true)
    setError(null)

    try {
      const { error } = await supabase.auth.signUp({
        email,
        password,
        options: {
          emailRedirectTo: `${window.location.origin}/auth/callback`,
        },
      })

      if (error) {
        setError(error.message)
      } else {
        setError('確認メールを送信しました。メールをご確認ください。')
      }
    } catch {
      setError('新規登録に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={className}>
      {error && (
        <p className="text-sm text-red-600 mb-4">{error}</p>
      )}
      
      <div className="space-y-3">
        <Button 
          onClick={handleSignIn} 
          disabled={loading}
          className="w-full"
        >
          {loading ? 'ログイン中...' : 'ログイン'}
        </Button>
        
        <Button 
          onClick={handleSignUp} 
          variant="outline"
          disabled={loading}
          className="w-full"
        >
          {loading ? '登録中...' : '新規登録'}
        </Button>
      </div>
    </div>
  )
}