import MapComponent from '@/components/map';

export default function Home() {
  return (
    <div className="min-h-screen bg-background">
      <header className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 py-4">
          <h1 className="text-2xl font-bold text-foreground">
            Bocchi The Map
          </h1>
          <p className="text-sm text-muted-foreground">
            おひとりさま向けスポットレビューアプリ
          </p>
        </div>
      </header>
      
      <main className="container mx-auto px-4 py-6">
        <div className="rounded-lg border bg-card">
          <MapComponent 
            className="rounded-lg overflow-hidden"
            height="600px"
          />
        </div>
        
        <div className="mt-6 text-center text-sm text-muted-foreground">
          <p>一人でも楽しめる場所を見つけて、レビューを共有しよう</p>
        </div>
      </main>
    </div>
  );
}
