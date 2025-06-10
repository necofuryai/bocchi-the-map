import { useCallback, useEffect, useRef, useState } from "react";
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
  // onClickコールバックをrefで管理してuseEffectの再実行を避ける最適化パターン
  // onClickが変更されるたびにメインのuseEffectが再実行されることを防ぐため、
  // 最新のコールバックをrefに保存してエフェクト内でアクセスする
  const currentOnClickRef = useRef<((event: maplibregl.MapMouseEvent) => void) | undefined>(onClick);

  // onLoadとonErrorをuseCallbackでメモ化
  const memoizedOnLoad = useCallback((map: maplibregl.Map) => {
    onLoad?.(map);
  }, [onLoad]);

  const memoizedOnError = useCallback((error: MapError) => {
    onError?.(error);
  }, [onError]);

  useEffect(() => {
    if (mapRef.current || !containerRef.current) return;
    
    // 環境変数チェック
    if (!process.env.NEXT_PUBLIC_MAP_STYLE_URL) {
      const configError: MapError = {
        type: 'configuration',
        message: 'NEXT_PUBLIC_MAP_STYLE_URL が設定されていません'
      };
      setError(configError);
      setMapState('error');
      memoizedOnError(configError);
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
        center: [139.767, 35.681],
        zoom: 15
      });

      // 地図読み込み完了時
      mapRef.current.on('load', () => {
        setMapState('loaded');
        setError(null);
        
        if (mapRef.current) {
          setupPOIFeatures(mapRef.current);
          memoizedOnLoad(mapRef.current);
        }
      });

      // エラーハンドリング
      mapRef.current.on('error', (e: maplibregl.ErrorEvent) => {
        console.error("Map error:", e);
        const loadError: MapError = {
          type: 'loading',
          message: '地図の読み込みに失敗しました',
          originalError: e
        };
        setError(loadError);
        setMapState('error');
        memoizedOnError(loadError);
      });

      // クリックイベント
      if (onClick) {
        mapRef.current.on('click', handleClick);
      }

    } catch (error) {
      console.error("地図の初期化に失敗しました:", error);
      const initError: MapError = {
        type: 'initialization',
        message: '地図の初期化に失敗しました',
        originalError: error as maplibregl.ErrorEvent
      };
      setError(initError);
      setMapState('error');
      memoizedOnError(initError);
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
  }, [memoizedOnLoad, memoizedOnError, onClick]);

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