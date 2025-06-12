import { useState, useCallback, useMemo } from 'react';

// POI kind types
export type POIKind = 'cafe' | 'atm' | 'restaurant' | 'shop' | 'park' | 'station' | string;

import * as maplibregl from "maplibre-gl";

// MapLibre GL filter expression types
export type FilterExpression = maplibregl.FilterSpecification | null;

// Filter configuration
export interface MapFilter {
  kinds: POIKind[];
  enabled: boolean;
}

// Hook for managing map filter state
export const useMapFilter = (initialKinds: POIKind[] = []) => {
  const [filter, setFilter] = useState<MapFilter>({
    kinds: initialKinds,
    enabled: true,
  });

  // Update filter kinds
  const updateKinds = useCallback((kinds: POIKind[]) => {
    setFilter(prev => ({
      ...prev,
      kinds,
      enabled: kinds.length > 0,
    }));
  }, []);

  // Toggle filter enabled state
  const toggleEnabled = useCallback(() => {
    setFilter(prev => ({
      ...prev,
      enabled: !prev.enabled,
    }));
  }, []);

  // Add kind to filter
  const addKind = useCallback((kind: POIKind) => {
    setFilter(prev => ({
      ...prev,
      kinds: [...prev.kinds, kind],
      enabled: true,
    }));
  }, []);

  // Remove kind from filter
  const removeKind = useCallback((kind: POIKind) => {
    setFilter(prev => {
      const newKinds = prev.kinds.filter(k => k !== kind);
      return {
        kinds: newKinds,
        enabled: newKinds.length > 0,
      };
    });
  }, []);

  // Clear all filters
  const clearFilter = useCallback(() => {
    setFilter({
      kinds: [],
      enabled: false,
    });
  }, []);

  // Generate MapLibre GL filter expression - show only specified POI types with valid names
  const getFilterExpression = useMemo((): FilterExpression => {
    const allowedPOITypes = [
      'cafe',
      'park',
      'library',
      'viewpoint',
      'bench',
      'toilets',
      'charging_station',
      'bicycle_parking',
      'drinking_water'
    ];

    return [
      "all",
      [
        "match",
        ["get", "kind"],
        allowedPOITypes,
        true,
        false,
      ],
      ["has", "name"],
      ["!=", ["get", "name"], ""]
    ];
  }, []);

  return {
    filter,
    updateKinds,
    toggleEnabled,
    addKind,
    removeKind,
    clearFilter,
    getFilterExpression,
  };
};