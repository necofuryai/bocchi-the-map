"use client";

import React from "react";
import { usePmtiles } from "../hooks/usePmtiles";
import { useMaplibre } from "../hooks/useMaplibre";
import { MapErrorDisplay, MapLoadingDisplay } from "./map/map-status";
import type { MapComponentProps } from "./map/types";

export default function MapComponent({ 
  className = "", 
  height = "480px", 
  onClick, 
  onLoad,
  onError,
  poiFilter
}: MapComponentProps) {
  // Initialize PMTiles protocol
  usePmtiles();

  // Initialize map and manage state
  const { containerRef, mapState, error } = useMaplibre({ 
    onClick, 
    onLoad, 
    onError,
    poiFilter
  });

  // Display error state
  if (mapState === 'error' && error) {
    return <MapErrorDisplay error={error} className={className} height={height} />;
  }

  return (
    <div className="relative">
      <div ref={containerRef} className={`w-full ${className}`} style={{ height }} />
      {mapState === 'loading' && <MapLoadingDisplay />}
    </div>
  );
}