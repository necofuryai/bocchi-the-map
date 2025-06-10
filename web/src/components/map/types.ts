import * as maplibregl from "maplibre-gl";

// エラーの種類を定義
export type MapError = 
  | { type: 'configuration'; message: string }
  | { type: 'initialization'; message: string; originalError?: maplibregl.ErrorEvent }
  | { type: 'loading'; message: string; originalError?: maplibregl.ErrorEvent };

// マップの状態を定義
export type MapState = 'loading' | 'loaded' | 'error';

// MapComponentのProps定義
export interface MapComponentProps {
  className?: string;
  height?: number | string;
  onClick?: (event: maplibregl.MapMouseEvent) => void;
  onLoad?: (map: maplibregl.Map) => void;
  onError?: (error: MapError) => void;
}