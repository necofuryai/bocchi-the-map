"use client"

import { useState, useEffect } from "react"
import { signIn, getSession } from "next-auth/react"
import { useRouter, useSearchParams } from "next/navigation"
import { MapPinIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function SignInPage() {
  const [isLoading, setIsLoading] = useState(false)
  const [provider, setProvider] = useState<string | null>(null)
  const router = useRouter()
  const searchParams = useSearchParams()
  const callbackUrl = searchParams.get('callbackUrl') || '/'
  const error = searchParams.get('error')

  useEffect(() => {
    const checkSession = async () => {
      const session = await getSession()
      if (session) {
        router.replace(callbackUrl)
      }
    }
    checkSession()
  }, [router, callbackUrl])

  const handleSignIn = async (providerName: string) => {
    setIsLoading(true)
    setProvider(providerName)
    
    try {
      const result = await signIn(providerName, {
        callbackUrl,
        redirect: false,
      })
      
      if (result?.error) {
        // Sign in failed - reset loading state
        console.error('Sign in error:', result.error)
        setIsLoading(false)
        setProvider(null)
      } else if (result?.url) {
        // Sign in succeeded - navigate to the callback URL
        router.replace(result.url)
      } else {
        // Fallback - reset loading state
        setIsLoading(false)
        setProvider(null)
      }
    } catch (error) {
      console.error('Sign in error:', error)
      setIsLoading(false)
      setProvider(null)
    }
  }

  const getErrorMessage = (error: string) => {
    switch (error) {
      case 'OAuthSignin':
        return 'OAuth認証でエラーが発生しました。'
      case 'OAuthCallback':
        return 'OAuth認証の処理中にエラーが発生しました。'
      case 'OAuthCreateAccount':
        return 'アカウント作成でエラーが発生しました。'
      case 'EmailCreateAccount':
        return 'メールアドレスでのアカウント作成でエラーが発生しました。'
      case 'Callback':
        return '認証処理でエラーが発生しました。'
      case 'OAuthAccountNotLinked':
        return 'このメールアドレスは別の認証方法で既に使用されています。'
      case 'EmailSignin':
        return 'メール認証でエラーが発生しました。'
      case 'CredentialsSignin':
        return '認証情報が正しくありません。'
      case 'SessionRequired':
        return 'この機能を使用するにはログインが必要です。'
      default:
        return '認証でエラーが発生しました。再度お試しください。'
    }
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="flex items-center justify-center space-x-2 mb-2">
            <MapPinIcon className="h-8 w-8 text-primary" />
            <h1 className="text-2xl font-bold">Bocchi The Map</h1>
          </div>
          <CardTitle className="text-xl">ログイン</CardTitle>
          <CardDescription>
            おひとりさまスポットを探すために、アカウントでログインしてください
          </CardDescription>
        </CardHeader>
        
        <CardContent className="space-y-4">
          {error && (
            <div className="p-3 bg-destructive/10 border border-destructive/20 rounded-md">
              <p className="text-sm text-destructive text-center">
                ⚠️ {getErrorMessage(error)}
              </p>
            </div>
          )}
          
          <div className="space-y-3">
            <Button
              variant="outline"
              size="lg"
              onClick={() => handleSignIn('google')}
              disabled={isLoading}
              className="w-full h-12 flex items-center justify-center space-x-2"
            >
              {isLoading && provider === 'google' ? (
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-primary" />
              ) : (
                <>
                  <svg className="h-5 w-5" viewBox="0 0 24 24" aria-hidden="true">
                    <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                    <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                    <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                    <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                  </svg>
                  <span>Googleでログイン</span>
                </>
              )}
            </Button>
            
            <Button
              variant="outline"
              size="lg"
              onClick={() => handleSignIn('twitter')}
              disabled={isLoading}
              className="w-full h-12 flex items-center justify-center space-x-2"
            >
              {isLoading && provider === 'twitter' ? (
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-primary" />
              ) : (
                <>
                  <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
                  </svg>
                  <span>Xでログイン</span>
                </>
              )}
            </Button>
          </div>

          <div className="text-xs text-muted-foreground text-center mt-6">
            ログインすることで、
            <br />
            利用規約およびプライバシーポリシーに同意したものとみなされます
          </div>
        </CardContent>
      </Card>
    </div>
  )
}