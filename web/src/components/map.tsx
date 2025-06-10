"use client";

import { usePmtiles } from "../hooks/usePmtiles";
import { useMaplibre } from "../hooks/useMaplibre";
import { MapErrorDisplay, MapLoadingDisplay } from "./map/map-status";
import type { MapComponentProps } from "./map/types";

export default function MapComponent({ 
  className = "", 
  height = "480px", 
  onClick, 
  onLoad,
  onError 
}: MapComponentProps) {
  // PMTiles プロトコルを初期化
  usePmtiles();

  // マップの初期化とステート管理
  const { containerRef, mapState, error } = useMaplibre({ 
    onClick, 
    onLoad, 
    onError 
  });

  // エラー状態の表示
  if (mapState === 'error' && error) {
    return <MapErrorDisplay error={error} className={className} height={height} />;
  }

  return (
    <div className="relative">
      <div ref={containerRef} className={`w-full ${className}`} style={{ height }} />
      {mapState === 'loading' && <MapLoadingDisplay />}
    </div>
  );
}