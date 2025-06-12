"use client";

import { useMemo } from 'react';
import MapComponent from '@/components/map';
import { Header } from '@/components/header';
import { Card } from '@/components/ui/card';
import { useMapFilter } from '@/hooks/useMapFilter';

export default function Home() {
  const { getFilterExpression } = useMapFilter([
    'cafe',
    'park', 
    'library',
    'viewpoint',
    'bench',
    'toilets',
    'charging_station',
    'bicycle_parking',
    'drinking_water'
  ]);

  // Memoize the filter expression to avoid unnecessary recalculations
  const filterExpression = useMemo(
    () => getFilterExpression(),
    [getFilterExpression]
  );

  return (
    <div className="min-h-screen bg-background">
      <Header />
      
      <main className="p-6">
          <div className="max-w-7xl mx-auto space-y-6">
            <div className="text-center space-y-2 mb-8">
              <h2 className="text-3xl font-bold tracking-tight">
                一人でも楽しめる場所を見つけよう
              </h2>
              <p className="text-muted-foreground">
                おひとりさまに優しいスポットをマップで探して、レビューを共有しましょう
              </p>
            </div>
            
            <Card className="overflow-hidden">
              <MapComponent 
                className="w-full"
                height="600px"
                poiFilter={filterExpression}
              />
            </Card>
            
            <div className="grid gap-4 md:grid-cols-3">
              <Card className="p-6">
                <h3 className="font-semibold mb-2">🎯 簡単検索</h3>
                <p className="text-sm text-muted-foreground">
                  カテゴリーや場所から、あなたにぴったりのスポットを見つけましょう
                </p>
              </Card>
              <Card className="p-6">
                <h3 className="font-semibold mb-2">💬 リアルな口コミ</h3>
                <p className="text-sm text-muted-foreground">
                  実際に訪れた人のレビューで、一人でも入りやすいお店がわかります
                </p>
              </Card>
              <Card className="p-6">
                <h3 className="font-semibold mb-2">📍 マップ表示</h3>
                <p className="text-sm text-muted-foreground">
                  現在地から近いおひとりさまスポットを地図上で確認できます
                </p>
              </Card>
            </div>
          </div>
        </main>
    </div>
  );
}
