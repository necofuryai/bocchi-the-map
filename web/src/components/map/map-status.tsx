import * as React from "react";
import type { MapError } from "./types";

interface MapErrorDisplayProps {
  error: MapError;
  className?: string;
  height?: string | number;
}

export function MapErrorDisplay({ error, className = "", height }: MapErrorDisplayProps) {
  return (
    <div 
      className={`w-full ${className} flex items-center justify-center bg-gray-100 border border-gray-300 rounded`} 
      style={{ height }}
      role="alert"
    >
      <div className="text-center text-gray-600">
        <p className="text-sm"><span role="img" aria-label="Error">❌</span> {error.message}</p>
        <p className="text-xs mt-1">
          {error.type === 'configuration' ? 'Please check your configuration' : 'Failed to display map'}
        </p>
      </div>
    </div>
  );
}

interface MapLoadingDisplayProps {
  className?: string;
}

export function MapLoadingDisplay({ className = "" }: MapLoadingDisplayProps) {
  return (
    <div 
      className={`absolute inset-0 flex items-center justify-center bg-gray-50 border border-gray-300 rounded ${className}`}
      aria-live="polite"
    >
      <div className="text-center text-gray-600">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-2" />
        <p className="text-sm">Loading map...</p>
      </div>
    </div>
  );
}