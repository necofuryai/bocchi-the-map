"use client";

import { usePmtiles } from "../hooks/usePmtiles";
import { useMaplibre } from "../hooks/useMaplibre";
import { MapErrorDisplay, MapLoadingDisplay } from "./map/map-status";
import type { MapComponentProps } from "./map/types";

// Default map height - can be overridden via environment variable or props
const DEFAULT_MAP_HEIGHT = process.env.NEXT_PUBLIC_DEFAULT_MAP_HEIGHT || "480px";

export default function MapComponent({ 
  className = "", 
  height = DEFAULT_MAP_HEIGHT, 
  onClick, 
  onLoad,
  onError
}: MapComponentProps) {
  // Initialize PMTiles protocol
  usePmtiles();

  // Initialize map and manage state
  const { containerRef, mapState, error } = useMaplibre({ 
    onClick, 
    onLoad, 
    onError
  });

  // Display error state
  if (mapState === 'error' && error) {
    return <MapErrorDisplay error={{ type: 'initialization', message: error }} className={className} height={height} />;
  }

  return (
    <div className="relative">
      <div ref={containerRef} className={`w-full ${className}`} style={{ height }} data-testid="map-container" />
      {mapState === 'loading' && <MapLoadingDisplay />}
    </div>
  );
}