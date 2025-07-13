import * as maplibregl from "maplibre-gl";
import { layers, namedFlavor } from "@protomaps/basemaps";

// Map style configuration constants
const MAP_STYLE_CONFIG = {
  glyphsUrl: process.env.NEXT_PUBLIC_MAP_GLYPHS_URL || "https://protomaps.github.io/basemaps-assets/fonts/{fontstack}/{range}.pbf",
  spriteUrl: process.env.NEXT_PUBLIC_MAP_SPRITE_URL || "https://protomaps.github.io/basemaps-assets/sprites/v4/light",
  textFont: ["Noto Sans Regular"] as string[],
  textSize: 11,
  textColor: "#333",
  textHaloColor: "#fff",
  textHaloWidth: 1,
  textOffset: [0, 0.7] as [number, number]
};

/**
 * Function to generate style configuration for MapLibre GL
 * @param mapStyleUrl Map style URL for PMTiles
 * @returns MapLibre GL StyleSpecification
 */
export function createMapStyle(mapStyleUrl: string): maplibregl.StyleSpecification {
  return {
    version: 8,
    glyphs: MAP_STYLE_CONFIG.glyphsUrl,
    sprite: MAP_STYLE_CONFIG.spriteUrl,
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
          "text-font": MAP_STYLE_CONFIG.textFont,
          "text-size": MAP_STYLE_CONFIG.textSize,
          "text-anchor": "top",
          "text-offset": MAP_STYLE_CONFIG.textOffset,
          "icon-allow-overlap": true
        },
        paint: {
          "text-color": MAP_STYLE_CONFIG.textColor,
          "text-halo-color": MAP_STYLE_CONFIG.textHaloColor,
          "text-halo-width": MAP_STYLE_CONFIG.textHaloWidth
        }
      }
    ]
  };
}