import { useCallback, useEffect, useRef, useState } from "react";
import * as maplibregl from "maplibre-gl";
import { createMapStyle } from "@/components/mapStyle";
import { setupPOIFeatures } from "@/components/map/poi-features";
import type { MapError, MapState } from "@/components/map/types";

interface UseMaplibreOptions {
  onClick?: (event: maplibregl.MapMouseEvent) => void;
  onLoad?: (map: maplibregl.Map) => void;
  onError?: (error: MapError) => void;
  defaultCenter?: [number, number];
  defaultZoom?: number;
}

export const useMaplibre = ({ 
  onClick, 
  onLoad, 
  onError,
  defaultCenter = [139.767, 35.681],
  defaultZoom = 15
}: UseMaplibreOptions = {}) => {
  const mapRef = useRef<maplibregl.Map | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [error, setError] = useState<MapError | null>(null);
  const [mapState, setMapState] = useState<MapState>('loading');
  // onClickコールバックをrefで管理してuseEffectの再実行を避ける最適化パターン
  // onClickが変更されるたびにメインのuseEffectが再実行されることを防ぐため、
  // 最新のコールバックをrefに保存してエフェクト内でアクセスする
  const currentOnClickRef = useRef<((event: maplibregl.MapMouseEvent) => void) | undefined>(onClick);


  useEffect(() => {
    if (mapRef.current || !containerRef.current) return;
    
    // 環境変数チェック
    if (!process.env.NEXT_PUBLIC_MAP_STYLE_URL) {
      const configError: MapError = {
        type: 'configuration',
        message: 'NEXT_PUBLIC_MAP_STYLE_URL is not configured'
      };
      setError(configError);
      setMapState('error');
      onError?.(configError);
      return;
    }
    
    // クリックイベント処理関数をここで定義
    const handleClick = (event: maplibregl.MapMouseEvent) => {
      currentOnClickRef.current?.(event);
    };
    
    try {
      const style = createMapStyle(process.env.NEXT_PUBLIC_MAP_STYLE_URL);

      mapRef.current = new maplibregl.Map({
        container: containerRef.current,
        style,
        center: defaultCenter,
        zoom: defaultZoom
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
          type: 'loading',
          message: 'Failed to load map',
          originalError: e
        };
        setError(loadError);
        setMapState('error');
        onError?.(loadError);
      });

      // クリックイベント
      if (onClick) {
        mapRef.current.on('click', handleClick);
      }

    } catch (error) {
      console.error("Map initialization failed:", error);
      const initError: MapError = {
        type: 'initialization',
        message: 'Failed to initialize map',
        ...(error instanceof Error && { originalError: error })
      };
      setError(initError);
      setMapState('error');
      onError?.(initError);
    }

    return () => {
      if (mapRef.current) {
        if (onClick) {
          // クリーンアップ時にイベントリスナーを削除
          // handleClick関数はrefを通じて最新のonClickコールバックにアクセスするため、
          // onClickが変更されてもイベントリスナーの再登録は不要
          mapRef.current.off('click', handleClick);
        }
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
  }, [onLoad, onError, onClick, defaultCenter, defaultZoom]);

  // onClickハンドラーが変更されたときにrefを更新
  useEffect(() => {
    currentOnClickRef.current = onClick;
  }, [onClick]);

  return {
    containerRef,
    mapState,
    error,
    map: mapRef.current
  };
};