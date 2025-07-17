import { useEffect, useRef } from "react";
import * as maplibregl from "maplibre-gl";
import { createMapStyle } from "@/components/mapStyle";
import { useMapStore } from '@/stores/use-map-store';
import type { MapError } from "@/components/map/types";

// Map default configuration
const MAP_DEFAULTS = {
  center: [
    parseFloat(process.env.NEXT_PUBLIC_DEFAULT_MAP_LONGITUDE || '139.767'),
    parseFloat(process.env.NEXT_PUBLIC_DEFAULT_MAP_LATITUDE || '35.681')
  ] as [number, number],
  zoom: parseFloat(process.env.NEXT_PUBLIC_DEFAULT_MAP_ZOOM || '15')
};

interface UseMaplibreOptions {
  onClick?: (event: maplibregl.MapMouseEvent) => void;
  onLoad?: (map: maplibregl.Map) => void;
  onError?: (error: MapError) => void;
  defaultCenter?: [number, number];
  defaultZoom?: number;
}

export const useMaplibre = ({ 
  onClick, 
  onLoad, 
  onError,
  defaultCenter = MAP_DEFAULTS.center,
  defaultZoom = MAP_DEFAULTS.zoom
}: UseMaplibreOptions = {}) => {
  const mapRef = useRef<maplibregl.Map | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const { mapState, error, setMapState, setError } = useMapStore();
  const currentOnClickRef = useRef<((event: maplibregl.MapMouseEvent) => void) | undefined>(onClick);


  useEffect(() => {
    if (mapRef.current || !containerRef.current) return;
    
    // Environment variable check
    if (!process.env.NEXT_PUBLIC_MAP_STYLE_URL) {
      const configError: MapError = {
        type: 'configuration',
        message: 'NEXT_PUBLIC_MAP_STYLE_URL is not configured'
      };
      setError(configError.message);
      setMapState('error');
      onError?.(configError);
      return;
    }
    
    // Define click event handler function here
    const handleClick = (event: maplibregl.MapMouseEvent) => {
      currentOnClickRef.current?.(event);
    };
    
    try {
      console.log("Creating map style with URL:", process.env.NEXT_PUBLIC_MAP_STYLE_URL);
      const style = createMapStyle(process.env.NEXT_PUBLIC_MAP_STYLE_URL);

      console.log("Initializing MapLibre GL map...");
      mapRef.current = new maplibregl.Map({
        container: containerRef.current,
        style,
        center: defaultCenter,
        zoom: defaultZoom
      });

      console.log("Map instance created successfully");

      // On map load completion
      mapRef.current.on('load', () => {
        if (mapRef.current) {
          console.log("Map loaded successfully");
          setMapState('loaded');
          setError(null);
          onLoad?.(mapRef.current);
        }
      });

      // Error handling
      mapRef.current.on('error', (e: maplibregl.ErrorEvent) => {
        console.error("Map error:", e);
        
        // Check if this is a vector tile parsing error
        if (e.error && e.error.message) {
          console.error("Error details:", {
            message: e.error.message,
            error: e.error
          });
          
          // Check for specific vector tile ("vt") errors
          if (e.error.message.includes('vt') || e.error.message.includes('vector tile')) {
            console.error("Vector tile parsing error detected");
            console.error("This might be a PMTiles compatibility issue with MapLibre GL");
          }
        }
        
        const loadError: MapError = {
          type: 'loading',
          message: 'Failed to load map',
          originalError: e
        };
        setError(loadError.message);
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
         originalError: error instanceof maplibregl.ErrorEvent ? error : undefined
       };
      setError(initError.message);
      setMapState('error');
      onError?.(initError);
    }

    return () => {
      if (mapRef.current) {
        if (onClick) {
          mapRef.current.off('click', handleClick);
        }
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [onLoad, onError, onClick, defaultCenter, defaultZoom]);

  useEffect(() => {
    currentOnClickRef.current = onClick;
  }, [onClick]);


  return {
    containerRef,
    mapState,
    error,
    map: mapRef.current
  };
};