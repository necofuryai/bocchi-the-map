import * as maplibregl from "maplibre-gl";
import { escapeHtml } from "@/lib/utils";
import type { POIProperties } from "@/types";

// POI関連の色定数
export const POI_COLORS = {
  PRIMARY: "#FF4081",
  STROKE: "#FFFFFF",
} as const;

/**
 * POI機能をセットアップする関数
 * @param map MapLibre GL マップインスタンス
 */
export const setupPOIFeatures = (map: maplibregl.Map): void => {
  // フォールバックの丸ポチ
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
  
  // POIアイコンレイヤー
  map.addLayer({
    id: "poi-icons",
    type: "symbol",
    source: "protomaps",
    "source-layer": "pois",
    minzoom: 13,
    layout: {
      // cafe => cafe-15.png (Maki 規約) + フォールバック
      "icon-image": [
        "coalesce",
        ["image", ["concat", ["get", "kind"], "-15"]],
        "circle-15" // デフォルトのフォールバックアイコン
      ],
      "icon-size": 1,
      "icon-allow-overlap": true,
    },
  });

  // POIクリック時のポップアップハンドラ
  const popupHandler = (e: maplibregl.MapLayerMouseEvent) => {
    if (!e.features?.length) return;

    const f = e.features[0];
    const p = f.properties as POIProperties;

    // 表示用デフォルト値
    const name = p.name ?? "名称不明";
    const kind = p.kind ?? "unknown";
    const zoom = p.min_zoom?.toString() ?? "-";

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

  // POIレイヤーにイベントハンドラを追加
  ["poi-icons", "poi-dots"].forEach((layer) => {
    map.on("click", layer, popupHandler);
    map.on("mouseenter", layer, () => (map.getCanvas().style.cursor = "pointer"));
    map.on("mouseleave", layer, () => (map.getCanvas().style.cursor = ""));
  });
};