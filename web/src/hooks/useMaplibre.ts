import { useEffect, useRef, useState } from "react";
import * as maplibregl from "maplibre-gl";
import { createMapStyle } from "@/components/mapStyle";
import { setupPOIFeatures, updatePOIFilter } from "@/components/map/poi-features";
import type { MapError, MapState } from "@/components/map/types";

interface UseMaplibreOptions {
  onClick?: (event: maplibregl.MapMouseEvent) => void;
  onLoad?: (map: maplibregl.Map) => void;
  onError?: (error: MapError) => void;
  defaultCenter?: [number, number];
  defaultZoom?: number;
  poiFilter?: maplibregl.FilterSpecification | null;
}

export const useMaplibre = ({ 
  onClick, 
  onLoad, 
  onError,
  defaultCenter = [139.767, 35.681],
  defaultZoom = 15,
  poiFilter
}: UseMaplibreOptions = {}) => {
  const mapRef = useRef<maplibregl.Map | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [error, setError] = useState<MapError | null>(null);
  const [mapState, setMapState] = useState<MapState>('loading');
  // Optimization pattern to manage onClick callback with ref to avoid useEffect re-execution
  // To prevent main useEffect from re-executing every time onClick changes,
  // store the latest callback in ref and access it within the effect
  const currentOnClickRef = useRef<((event: maplibregl.MapMouseEvent) => void) | undefined>(onClick);


  useEffect(() => {
    if (mapRef.current || !containerRef.current) return;
    
    // Environment variable check
    if (!process.env.NEXT_PUBLIC_MAP_STYLE_URL) {
      const configError: MapError = {
        type: 'configuration',
        message: 'NEXT_PUBLIC_MAP_STYLE_URL is not configured'
      };
      setError(configError);
      setMapState('error');
      onError?.(configError);
      return;
    }
    
    // Define click event handler function here
    const handleClick = (event: maplibregl.MapMouseEvent) => {
      currentOnClickRef.current?.(event);
    };
    
    try {
      const style = createMapStyle(process.env.NEXT_PUBLIC_MAP_STYLE_URL);

      mapRef.current = new maplibregl.Map({
        container: containerRef.current,
        style,
        center: defaultCenter,
        zoom: defaultZoom
      });

      // On map load completion
      mapRef.current.on('load', () => {
        setMapState('loaded');
        setError(null);
        
        if (mapRef.current) {
          setupPOIFeatures(mapRef.current, poiFilter || null);
          onLoad?.(mapRef.current);
        }
      });

      // Error handling
      mapRef.current.on('error', (e: maplibregl.ErrorEvent) => {
        console.error("Map error:", e);
        const loadError: MapError = {
          type: 'loading',
          message: 'Failed to load map',
          originalError: e
        };
        setError(loadError);
        setMapState('error');
        onError?.(loadError);
      });

      // Click event
      if (onClick) {
        mapRef.current.on('click', handleClick);
      }

    } catch (error) {
      console.error("Map initialization failed:", error);
      const initError: MapError = {
        type: 'initialization',
        message: 'Failed to initialize map',
        originalError: error as Error
      };
      setError(initError);
      setMapState('error');
      onError?.(initError);
    }

    return () => {
      if (mapRef.current) {
        if (onClick) {
          // Remove event listener during cleanup
          // Since handleClick function accesses the latest onClick callback through ref,
          // event listener re-registration is unnecessary even when onClick changes
          mapRef.current.off('click', handleClick);
        }
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
    // Note: poiFilter is intentionally excluded from dependencies to prevent map recreation
    // POI filter updates are handled by a separate useEffect below
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [onLoad, onError, onClick, defaultCenter, defaultZoom]);

  // Update ref when onClick handler changes
  useEffect(() => {
    currentOnClickRef.current = onClick;
  }, [onClick]);

  // Update POI filter when poiFilter prop changes
  useEffect(() => {
    if (mapRef.current && mapState === 'loaded') {
      updatePOIFilter(mapRef.current, poiFilter || null);
    }
  }, [poiFilter, mapState]);

  return {
    containerRef,
    mapState,
    error,
    map: mapRef.current
  };
};