"use client";

import { useEffect, useRef, useState } from "react";
import * as maplibregl from "maplibre-gl";
import { Protocol } from "pmtiles";
import { createMapStyle } from "./mapStyle";

interface MapComponentProps {
  className?: string;
  height?: string;
}

// Protocol reference counting to handle multiple map instances
let protocolRefCount = 0;
let protocolInstance: Protocol | null = null;

export default function MapComponent({ className = "", height = "480px" }: MapComponentProps) {
  const mapRef = useRef<maplibregl.Map | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (mapRef.current || !containerRef.current) return; // 二重 init 防止 & DOM チェック
    
    // NEXT_PUBLIC_MAP_STYLE_URLが未定義の場合は早期に検知
    if (!process.env.NEXT_PUBLIC_MAP_STYLE_URL) {
      const errorMsg = "NEXT_PUBLIC_MAP_STYLE_URL が設定されていません";
      setError(errorMsg);
      setIsLoading(false);
      return;
    }
    
    try {
      // Register protocol only if it's the first instance
      if (protocolRefCount === 0) {
        protocolInstance = new Protocol();
        maplibregl.addProtocol("pmtiles", protocolInstance.tile.bind(protocolInstance));
      }
      protocolRefCount++;

      const style = createMapStyle(process.env.NEXT_PUBLIC_MAP_STYLE_URL);

      mapRef.current = new maplibregl.Map({
        container: containerRef.current,
        style,
        center: [139.767, 35.681],
        zoom: 15
      });

      mapRef.current.on('load', () => {
        setIsLoading(false);
      });

      mapRef.current.on('error', (e) => {
        console.error("Map error:", e);
        setError("地図の読み込みに失敗しました");
        setIsLoading(false);
      });

    } catch (error) {
      console.error("地図の初期化に失敗しました:", error);
      setError("地図の初期化に失敗しました");
      setIsLoading(false);
    }

    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
      
      // Only remove protocol when the last instance unmounts
      protocolRefCount--;
      if (protocolRefCount === 0 && protocolInstance) {
        maplibregl.removeProtocol("pmtiles");
        protocolInstance = null;
      }
    };
  }, []);

  if (error) {
    return (
      <div 
        className={`w-full ${className} flex items-center justify-center bg-gray-100 border border-gray-300 rounded`} 
        style={{ height }}
      >
        <div className="text-center text-gray-600">
          <p className="text-sm">❌ {error}</p>
          <p className="text-xs mt-1">地図を表示できませんでした</p>
        </div>
      </div>
    );
  }

  return (
    <div className="relative">
      <div ref={containerRef} className={`w-full ${className}`} style={{ height }} />
      {isLoading && (
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