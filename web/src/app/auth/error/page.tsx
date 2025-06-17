"use client"

import { useEffect, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { MapPinIcon, AlertTriangleIcon, RefreshCwIcon, HomeIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function AuthErrorPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const error = searchParams.get('error')
  const [errorDetails, setErrorDetails] = useState<{
    title: string
    description: string
    suggestion: string
  }>()

  useEffect(() => {
    const getErrorDetails = (error: string) => {
      switch (error) {
        case 'Configuration':
          return {
            title: '設定エラー',
            description: 'アプリの認証設定に問題があります。',
            suggestion: 'しばらく時間をおいてから再度お試しください。問題が続く場合は、サポートにご連絡ください。'
          }
        case 'AccessDenied':
          return {
            title: 'アクセス拒否',
            description: '認証プロバイダーからアクセスが拒否されました。',
            suggestion: '異なる認証方法をお試しいただくか、時間をおいて再度お試しください。'
          }
        case 'Verification':
          return {
            title: '認証エラー',
            description: 'メール認証の確認に失敗しました。',
            suggestion: '認証リンクの有効期限が切れている可能性があります。再度ログインをお試しください。'
          }
        case 'OAuthSignin':
          return {
            title: 'OAuth認証エラー',
            description: 'OAuth認証プロバイダーとの接続に問題があります。',
            suggestion: 'ネットワーク接続を確認して、再度お試しください。'
          }
        case 'OAuthCallback':
          return {
            title: 'OAuth認証処理エラー',
            description: 'OAuth認証の処理中にエラーが発生しました。',
            suggestion: 'ポップアップブロッカーが有効になっていないか確認して、再度お試しください。'
          }
        case 'OAuthCreateAccount':
          return {
            title: 'アカウント作成エラー',
            description: 'OAuth認証を使用したアカウント作成に失敗しました。',
            suggestion: 'アカウント情報に問題がある可能性があります。異なる認証方法をお試しください。'
          }
        case 'EmailCreateAccount':
          return {
            title: 'メールアカウント作成エラー',
            description: 'メールアドレスを使用したアカウント作成に失敗しました。',
            suggestion: 'メールアドレスが正しいか確認して、再度お試しください。'
          }
        case 'Callback':
          return {
            title: '認証コールバックエラー',
            description: '認証プロセスの最終段階でエラーが発生しました。',
            suggestion: 'Cookieが有効になっているか確認して、再度お試しください。'
          }
        case 'OAuthAccountNotLinked':
          return {
            title: 'アカウント連携エラー',
            description: 'このメールアドレスは既に別の認証方法で使用されています。',
            suggestion: '以前に使用した認証方法でログインするか、サポートにお問い合わせください。'
          }
        case 'EmailSignin':
          return {
            title: 'メール認証エラー',
            description: 'メール認証に失敗しました。',
            suggestion: 'メールアドレスが正しいか確認して、再度お試しください。'
          }
        case 'CredentialsSignin':
          return {
            title: 'ログイン情報エラー',
            description: 'ユーザー名またはパスワードが正しくありません。',
            suggestion: 'ログイン情報を確認して、再度お試しください。'
          }
        case 'SessionRequired':
          return {
            title: 'セッション必須',
            description: 'この機能を使用するにはログインが必要です。',
            suggestion: 'ログインしてから再度お試しください。'
          }
        case 'Default':
        default:
          return {
            title: '認証エラー',
            description: '予期しないエラーが発生しました。',
            suggestion: 'しばらく時間をおいてから再度お試しください。問題が続く場合は、サポートにご連絡ください。'
          }
      }
    }

    setErrorDetails(getErrorDetails(error || 'Default'))
  }, [error])

  const handleRetry = () => {
    router.push('/auth/signin')
  }

  const handleGoHome = () => {
    router.push('/')
  }

  if (!errorDetails) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="flex items-center justify-center space-x-2 mb-2">
            <MapPinIcon className="h-8 w-8 text-primary" />
            <h1 className="text-2xl font-bold">Bocchi The Map</h1>
          </div>
          <div className="flex items-center justify-center mb-2">
            <AlertTriangleIcon className="h-8 w-8 text-destructive" />
          </div>
          <CardTitle className="text-xl text-destructive">
            {errorDetails.title}
          </CardTitle>
          <CardDescription>
            {errorDetails.description}
          </CardDescription>
        </CardHeader>
        
        <CardContent className="space-y-4">
          <div className="p-4 bg-muted rounded-lg">
            <p className="text-sm text-muted-foreground text-center">
              💡 {errorDetails.suggestion}
            </p>
          </div>
          
          {error && (
            <div className="p-3 bg-destructive/10 border border-destructive/20 rounded-md">
              <p className="text-xs text-destructive/80 text-center">
                エラーコード: {error}
              </p>
            </div>
          )}
          
          <div className="space-y-3">
            <Button
              onClick={handleRetry}
              className="w-full"
              size="lg"
            >
              <RefreshCwIcon className="h-4 w-4 mr-2" />
              再度ログインを試す
            </Button>
            
            <Button
              onClick={handleGoHome}
              variant="outline"
              className="w-full"
              size="lg"
            >
              <HomeIcon className="h-4 w-4 mr-2" />
              ホームに戻る
            </Button>
          </div>

          <div className="text-xs text-muted-foreground text-center mt-6">
            問題が解決しない場合は、
            <br />
            サポートまでお問い合わせください
          </div>
        </CardContent>
      </Card>
    </div>
  )
}