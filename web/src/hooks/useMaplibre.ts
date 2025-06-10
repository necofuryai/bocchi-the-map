import { useEffect, useRef, useState } from "react";
import * as maplibregl from "maplibre-gl";
import { createMapStyle } from "@/components/mapStyle";
import { setupPOIFeatures } from "@/components/map/poi-features";
import type { MapError, MapState } from "@/components/map/types";

interface UseMaplibreOptions {
  onClick?: (event: maplibregl.MapMouseEvent) => void;
  onLoad?: (map: maplibregl.Map) => void;
  onError?: (error: MapError) => void;
}

export const useMaplibre = ({ onClick, onLoad, onError }: UseMaplibreOptions = {}) => {
  const mapRef = useRef<maplibregl.Map | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [error, setError] = useState<MapError | null>(null);
  const [mapState, setMapState] = useState<MapState>('loading');

  useEffect(() => {
    if (mapRef.current || !containerRef.current) return;
    
    // 環境変数チェック
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

      // 地図読み込み完了時
      mapRef.current.on('load', () => {
        setMapState('loaded');
        setError(null);
        
        if (mapRef.current) {
          setupPOIFeatures(mapRef.current);
          onLoad?.(mapRef.current);
        }
      });

      // エラーハンドリング
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

      // クリックイベント
      if (onClick) {
        mapRef.current.on('click', onClick);
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

  return {
    containerRef,
    mapState,
    error,
    map: mapRef.current
  };
};