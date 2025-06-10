"use client";

import { useEffect, useRef, useState } from "react";
import * as maplibregl from "maplibre-gl";
import { createMapStyle } from "./mapStyle";
import { usePmtiles } from "../hooks/usePmtiles";

// POI関連の色定数
const POI_COLORS = {
  PRIMARY: "#FF4081",
  STROKE: "#FFFFFF",
} as const;

// HTMLエスケープ関数
const escapeHtml = (text: string | number): string => {
  const map: { [key: string]: string } = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  return String(text).replace(/[&<>"']/g, (m) => map[m]);
};

// POI機能をセットアップする関数
const setupPOIFeatures = (map: maplibregl.Map) => {
  // 2️⃣ フォールバックの丸ポチ（filter だけ修正）
  map.addLayer({
    id: "poi-dots",
    type: "circle",
    source: "protomaps",
    "source-layer": "pois",
    minzoom: 13,
    paint: {
      "circle-radius": 5,
      "circle-color": POI_COLORS.PRIMARY,
      "circle-stroke-width": 1,
      "circle-stroke-color": POI_COLORS.STROKE,
    },
  });
  
  // 1️⃣ アイコン：kind をそのままスプライト名に
  map.addLayer(
    {
      id: "poi-icons",
      type: "symbol",
      source: "protomaps",
      "source-layer": "pois",
      minzoom: 13,
      layout: {
        // cafe => cafe-15.png (Maki 規約) + フォールバック
        "icon-image": [
          "case",
          ["has", ["concat", ["get", "kind"], "-15"]], // アイコンが存在するかチェック
          ["concat", ["get", "kind"], "-15"], // 存在する場合は使用
          "circle-15" // デフォルトのフォールバックアイコン
        ],
        "icon-size": 1,
        "icon-allow-overlap": true,
      },
    },
    // "poi-dots" // ← circle レイヤーより上に置く
  );

  // 3️⃣ クリック共通ハンドラ
  const popupHandler = (e: maplibregl.MapLayerMouseEvent) => {
    if (!e.features?.length) return;

    const f = e.features[0];
    const p = f.properties as {
      name?: string;
      kind?: string;
      script?: string;
      min_zoom?: number;
    };

    // 表示用デフォルト
    const name = p.name ?? "名称不明";
    const kind = p.kind ?? "unknown";
    const zoom = p.min_zoom ?? "-";

    // HTMLエスケープを適用
    const escapedName = escapeHtml(name);
    const escapedKind = escapeHtml(kind);
    const escapedZoom = escapeHtml(zoom);

    new maplibregl.Popup({ offset: 8 })
      .setLngLat(
        f.geometry.type === "Point"
          ? (f.geometry.coordinates as [number, number])
          : [e.lngLat.lng, e.lngLat.lat]
      )
      .setHTML(`
        <div style="font-family: system-ui; font-size: 14px; line-height: 1.4; color: black;">
          <strong>${escapedName}</strong><br>
          種類: ${escapedKind}<br>
          最小ズーム: ${escapedZoom}
        </div>
      `)
      .addTo(map);
  };

  ["poi-icons", "poi-dots"].forEach((layer) => {
    map.on("click", layer, popupHandler);
    map.on("mouseenter", layer, () => (map.getCanvas().style.cursor = "pointer"));
    map.on("mouseleave", layer, () => (map.getCanvas().style.cursor = ""));
  });
};

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
        
        if (mapRef.current) {
          setupPOIFeatures(mapRef.current);
        }
        
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