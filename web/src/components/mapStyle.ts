import * as maplibregl from "maplibre-gl";
import { layers, namedFlavor } from "@protomaps/basemaps";

/**
 * Creates a simple map style configuration using Protomaps basemaps.
 * This is the simplified version that was working during the NextAuth era.
 * 
 * @param mapStyleUrl - The URL to the PMTiles file
 * @returns MapLibre GL style specification
 */
export function createMapStyle(mapStyleUrl: string): maplibregl.StyleSpecification {
  return {
    version: 8,
    glyphs: process.env.NEXT_PUBLIC_MAP_GLYPHS_URL || "https://protomaps.github.io/basemaps-assets/fonts/{fontstack}/{range}.pbf",
    sprite: process.env.NEXT_PUBLIC_MAP_SPRITE_URL || "https://protomaps.github.io/basemaps-assets/sprites/v4/light",
    sources: {
      protomaps: {
        type: "vector" as const,
        url: `pmtiles://${mapStyleUrl}`,
        attribution: "<a href=\"https://github.com/protomaps/basemaps\">Protomaps</a> © <a href=\"https://openstreetmap.org\">OpenStreetMap</a>"
      }
    },
    layers: [
      ...layers("protomaps", namedFlavor("light")),
      // Simple POI layer for basic point of interest display
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