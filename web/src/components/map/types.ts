import * as maplibregl from "maplibre-gl";

// Define error types
export type MapError = 
  | { type: 'configuration'; message: string }
  | { type: 'initialization'; message: string; originalError?: maplibregl.ErrorEvent }
  | { type: 'loading'; message: string; originalError?: maplibregl.ErrorEvent };

// Define map state
export type MapState = 'loading' | 'loaded' | 'error';

// MapComponent props definition
export interface MapComponentProps {
  className?: string;
  height?: number | string;
  onClick?: (event: maplibregl.MapMouseEvent) => void;
  onLoad?: (map: maplibregl.Map) => void;
  onError?: (error: MapError) => void;
  poiFilter?: maplibregl.FilterSpecification | null;
}