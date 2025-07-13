import Link from 'next/link'
import { MapPinIcon } from '@heroicons/react/24/outline'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export default function AuthCodeErrorPage() {
  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center space-y-4">
          <Link href="/" className="inline-flex items-center space-x-2">
            <MapPinIcon className="h-8 w-8 text-primary" />
            <h1 className="text-2xl font-bold">Bocchi The Map</h1>
          </Link>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="text-center text-red-600">認証エラー</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4 text-center">
            <p className="text-muted-foreground">
              認証プロセスでエラーが発生しました。
            </p>
            
            <div className="space-y-2">
              <Link href="/auth/login">
                <Button className="w-full">
                  再度ログインを試す
                </Button>
              </Link>
              
              <Link href="/">
                <Button variant="outline" className="w-full">
                  ホームに戻る
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}