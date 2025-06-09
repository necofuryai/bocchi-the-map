import * as maplibregl from "maplibre-gl";
import { layers, namedFlavor } from "@protomaps/basemaps";

/**
 * MapLibre GL用のスタイル設定を生成する関数
 * @param mapStyleUrl PMTiles用のマップスタイルURL
 * @returns MapLibre GL StyleSpecification
 */
export function createMapStyle(mapStyleUrl: string): maplibregl.StyleSpecification {
  return {
    version: 8,
    glyphs: "https://protomaps.github.io/basemaps-assets/fonts/{fontstack}/{range}.pbf",
    sprite: "https://protomaps.github.io/basemaps-assets/sprites/v4/light",
    sources: {
      protomaps: {
        type: "vector" as const,
        url: `pmtiles://${mapStyleUrl}`,
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
}