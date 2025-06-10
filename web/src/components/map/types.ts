import * as maplibregl from "maplibre-gl";

// エラーの種類を定義
export type MapError = 
  | { type: '設定'; message: string }
  | { type: '初期化'; message: string; originalError?: unknown }
  | { type: '読み込み'; message: string; originalError?: maplibregl.ErrorEvent };

// マップの状態を定義
export type MapState = 'loading' | 'loaded' | 'error';

// MapComponentのProps定義
export interface MapComponentProps {
  className?: string;
  height?: string;
  onClick?: (event: maplibregl.MapMouseEvent) => void;
  onLoad?: (map: maplibregl.Map) => void;
  onError?: (error: MapError) => void;
}