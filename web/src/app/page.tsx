"use client";

import MapComponent from '@/components/map';
import { Header } from '@/components/header';
import { Card } from '@/components/ui/card';
import { useMapFilter } from '@/hooks/useMapFilter';

// Default POI kinds to avoid array recreation on every render
const DEFAULT_KINDS = [
  'cafe',
  'park', 
  'library',
  'viewpoint',
  'bench',
  'toilets',
  'charging_station',
  'bicycle_parking',
  'drinking_water'
] as const;

export default function Home() {
  const { filterExpression } = useMapFilter(DEFAULT_KINDS);

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
          </div>
        </main>
    </div>
  );
}
