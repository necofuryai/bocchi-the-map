import * as maplibregl from "maplibre-gl";
import { escapeHtml } from "@/lib/utils";
import type { POIProperties } from "@/types";

// POI-related color constants
export const POI_COLORS = {
  PRIMARY: "#FF4081",
  STROKE: "#FFFFFF",
} as const;

/**
 * Function to set up POI features
 * @param map MapLibre GL map instance
 * @param filter Optional filter expression for POI layers
 */
export const setupPOIFeatures = (map: maplibregl.Map, filter?: maplibregl.FilterSpecification | null): void => {
  // Variable to track the current popup
  let currentPopup: maplibregl.Popup | null = null;
  
  // Popup handler for POI clicks
  const popupHandler = (e: maplibregl.MapLayerMouseEvent) => {
    if (!e.features?.length) return;

    const f = e.features[0];
    const p = f.properties as POIProperties;

    // Default values for display
    const name = p.name ?? "Unknown name";
    const kind = p.kind ?? "unknown";
    const zoom = p.min_zoom?.toString() ?? "-";

    // Apply HTML escaping
    const escapedName = escapeHtml(name);
    const escapedKind = escapeHtml(kind);
    const escapedZoom = escapeHtml(zoom);

    // Remove existing popup
    if (currentPopup) {
      currentPopup.remove();
    }

    // Create and track new popup
    currentPopup = new maplibregl.Popup({ offset: 8 })
      .setLngLat(
        f.geometry.type === "Point"
          ? (f.geometry.coordinates as [number, number])
          : [e.lngLat.lng, e.lngLat.lat]
      )
      .setHTML(`
        <div style="font-family: system-ui; font-size: 14px; line-height: 1.4; color: black;">
          <strong>${escapedName}</strong><br>
          Type: ${escapedKind}<br>
          Min Zoom: ${escapedZoom}
        </div>
      `)
      .addTo(map);
  };
  
  // Wait for sprites to load
  const setupLayersWithValidatedIcons = () => {
    // Fallback circle dots - check existing layers
    if (!map.getLayer("poi-dots")) {
      const layerConfig: maplibregl.CircleLayerSpecification = {
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
      };
      
      // Apply filter if provided
      if (filter) {
        layerConfig.filter = filter;
      }
      
      map.addLayer(layerConfig);
    }
    
    // POI icon layer - check if icons actually exist
    if (!map.getLayer("poi-icons")) {
      const iconLayerConfig: maplibregl.SymbolLayerSpecification = {
        id: "poi-icons",
        type: "symbol",
        source: "protomaps",
        "source-layer": "pois",
        minzoom: 13,
        layout: {
          // Fallback with icon existence check
          "icon-image": [
            "case",
            ["has", ["concat", ["get", "kind"], "-15"]], // Check if icon exists
            ["concat", ["get", "kind"], "-15"], // Use if exists
            "circle-15" // Fallback if not exists
          ],
          "icon-size": 1,
          "icon-allow-overlap": true,
        },
      };
      
      // Apply filter if provided
      if (filter) {
        iconLayerConfig.filter = filter;
      }
      
      map.addLayer(iconLayerConfig);
    }
    
    // Add event handlers to POI layers
    ["poi-icons", "poi-dots"].forEach((layer) => {
      map.on("click", layer, popupHandler);
      map.on("mouseenter", layer, () => (map.getCanvas().style.cursor = "pointer"));
      map.on("mouseleave", layer, () => (map.getCanvas().style.cursor = ""));
    });
  };

  // Check if sprites are loaded
  if (map.isStyleLoaded()) {
    setupLayersWithValidatedIcons();
  } else {
    map.on("styledata", setupLayersWithValidatedIcons);
  }
};

/**
 * Function to update POI filter on existing layers
 * @param map MapLibre GL map instance
 * @param filter Filter expression to apply (null to remove filter)
 */
export const updatePOIFilter = (map: maplibregl.Map, filter: maplibregl.FilterSpecification | null): void => {
  const layerIds = ["poi-dots", "poi-icons"];
  
  layerIds.forEach((layerId) => {
    if (map.getLayer(layerId)) {
      if (filter) {
        map.setFilter(layerId, filter);
      } else {
        map.setFilter(layerId, null);
      }
    }
  });
};