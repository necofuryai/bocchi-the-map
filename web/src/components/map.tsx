"use client";

import { useEffect, useRef, useState } from "react";
import * as maplibregl from "maplibre-gl";
import { createMapStyle } from "./mapStyle";
import { usePmtiles } from "../hooks/usePmtiles";

// エラーの種類を定義
type MapError = 
  | { type: '設定'; message: string }
  | { type: '初期化'; message: string; originalError?: unknown }
  | { type: '読み込み'; message: string; originalError?: maplibregl.ErrorEvent };

// マップの状態を定義
type MapState = 'loading' | 'loaded' | 'error';

interface MapComponentProps {
  className?: string;
  height?: string;
  onClick?: (event: maplibregl.MapMouseEvent) => void;
  onLoad?: (map: maplibregl.Map) => void;
  onError?: (error: MapError) => void;
}

export default function MapComponent({ 
  className = "", 
  height = "480px", 
  onClick, 
  onLoad,
  onError 
}: MapComponentProps) {
  const mapRef = useRef<maplibregl.Map | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [error, setError] = useState<MapError | null>(null);
  const [mapState, setMapState] = useState<MapState>('loading');

  // PMTiles プロトコルを初期化
  usePmtiles();

  useEffect(() => {
    if (mapRef.current || !containerRef.current) return; // 二重 init 防止 & DOM チェック
    
    // NEXT_PUBLIC_MAP_STYLE_URLが未定義の場合は早期に検知
    if (!process.env.NEXT_PUBLIC_MAP_STYLE_URL) {
      const configError: MapError = {
        type: '設定',
        message: 'NEXT_PUBLIC_MAP_STYLE_URL が設定されていません'
      };
      setError(configError);
      setMapState('error');
      onError?.(configError);
      return;
    }
    
    try {
      const style = createMapStyle(process.env.NEXT_PUBLIC_MAP_STYLE_URL);

      mapRef.current = new maplibregl.Map({
        container: containerRef.current,
        style,
        center: [139.767, 35.681],
        zoom: 15
      });

      mapRef.current.on('load', () => {
        setMapState('loaded');
        setError(null);
        if (onLoad && mapRef.current) {
          onLoad(mapRef.current);
        }
      });

      mapRef.current.on('error', (e: maplibregl.ErrorEvent) => {
        console.error("Map error:", e);
        const loadError: MapError = {
          type: '読み込み',
          message: '地図の読み込みに失敗しました',
          originalError: e
        };
        setError(loadError);
        setMapState('error');
        onError?.(loadError);
      });

      // クリックイベントハンドラーを追加
      if (onClick) {
        mapRef.current.on('click', (e: maplibregl.MapMouseEvent) => {
          onClick(e);
        });
      }

    } catch (error) {
      console.error("地図の初期化に失敗しました:", error);
      const initError: MapError = {
        type: '初期化',
        message: '地図の初期化に失敗しました',
        originalError: error
      };
      setError(initError);
      setMapState('error');
      onError?.(initError);
    }

    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
  }, [onClick, onLoad, onError]);

  if (mapState === 'error' && error) {
    return (
      <div 
        className={`w-full ${className} flex items-center justify-center bg-gray-100 border border-gray-300 rounded`} 
        style={{ height }}
      >
        <div className="text-center text-gray-600">
          <p className="text-sm">❌ {error.message}</p>
          <p className="text-xs mt-1">
            {error.type === '設定' ? '設定を確認してください' : '地図を表示できませんでした'}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="relative">
      <div ref={containerRef} className={`w-full ${className}`} style={{ height }} />
      {mapState === 'loading' && (
        <div 
          className="absolute inset-0 flex items-center justify-center bg-gray-50 border border-gray-300 rounded"
        >
          <div className="text-center text-gray-600">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-2"></div>
            <p className="text-sm">地図を読み込み中...</p>
          </div>
        </div>
      )}
    </div>
  );
}