"use client";

import { useEffect, useRef } from "react";
import * as maplibregl from "maplibre-gl";
import { Protocol } from "pmtiles";
import { layers, namedFlavor } from "@protomaps/basemaps";
import "maplibre-gl/dist/maplibre-gl.css";

interface MapComponentProps {
  className?: string;
  height?: string;
}

export default function MapComponent({ className = "", height = "480px" }: MapComponentProps) {
  const mapRef = useRef<maplibregl.Map | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (mapRef.current || !containerRef.current) return; // 二重 init 防止 & DOM チェック
    
    try {
      const protocol = new Protocol();
      maplibregl.addProtocol("pmtiles", protocol.tile); // PMTiles 登録

      const style: maplibregl.StyleSpecification = {
        version: 8,
        glyphs: "https://protomaps.github.io/basemaps-assets/fonts/{fontstack}/{range}.pbf",
        sprite: "https://protomaps.github.io/basemaps-assets/sprites/v4/light",
        sources: {
          protomaps: {
            type: "vector" as const,
            url: `pmtiles://${process.env.NEXT_PUBLIC_MAP_STYLE_URL}`,
            attribution: "<a href=\"https://github.com/protomaps/basemaps\">Protomaps</a> © <a href=\"https://openstreetmap.org\">OpenStreetMap</a>"
          }
        },
        layers: [
          ...layers("protomaps", namedFlavor("light")),
          {
            id: "poi-cafe-atm",
            type: "symbol",
            source: "protomaps",
            "source-layer": "pois",
            layout: {
              "text-field": ["get", "name"],
              "text-font": ["Noto Sans Regular"],
              "text-size": 11,
              "text-anchor": "top",
              "text-offset": [0, 0.7],
              "icon-allow-overlap": true
            },
            paint: {
              "text-color": "#333",
              "text-halo-color": "#fff",
              "text-halo-width": 1
            }
          }
        ]
      };

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