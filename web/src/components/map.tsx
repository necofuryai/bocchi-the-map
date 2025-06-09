"use client";

import { useEffect, useRef } from "react";
import * as maplibregl from "maplibre-gl";
import { Protocol } from "pmtiles";
import { createMapStyle } from "./mapStyle";

interface MapComponentProps {
  className?: string;
  height?: string;
}

export default function MapComponent({ className = "", height = "480px" }: MapComponentProps) {
  const mapRef = useRef<maplibregl.Map | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (mapRef.current || !containerRef.current) return; // 二重 init 防止 & DOM チェック
    
    // NEXT_PUBLIC_MAP_STYLE_URLが未定義の場合は早期に検知
    if (!process.env.NEXT_PUBLIC_MAP_STYLE_URL) {
      throw new Error("NEXT_PUBLIC_MAP_STYLE_URL が設定されていません");
    }
    
    try {
      const protocol = new Protocol();
      maplibregl.addProtocol("pmtiles", protocol.tile.bind(protocol)); // PMTiles 登録

      const style = createMapStyle(process.env.NEXT_PUBLIC_MAP_STYLE_URL);

      mapRef.current = new maplibregl.Map({
        container: containerRef.current,
        style,
        center: [139.767, 35.681],
        zoom: 15
      });
    } catch (error) {
      console.error("地図の初期化に失敗しました:", error);
    }

    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
      maplibregl.removeProtocol("pmtiles");
    };
  }, []);

  return <div ref={containerRef} className={`w-full ${className}`} style={{ height }} />;
}