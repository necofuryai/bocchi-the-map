import { useState, useCallback, useMemo } from 'react';

// POI kind types
export type POIKind =
  | 'cafe'
  | 'park'
  | 'library'
  | 'viewpoint'
  | 'bench'
  | 'toilets'
  | 'charging_station'
  | 'bicycle_parking'
  | 'drinking_water';

import * as maplibregl from "maplibre-gl";

// MapLibre GL filter expression types
export type FilterExpression = maplibregl.FilterSpecification | null;

// Filter configuration
export interface MapFilter {
  kinds: POIKind[];
  enabled: boolean;
}

// Hook for managing map filter state
export const useMapFilter = (initialKinds: readonly POIKind[] = []) => {
  const [filter, setFilter] = useState<MapFilter>({
    kinds: [...initialKinds],
    enabled: initialKinds.length > 0,
  });

  // Update filter kinds
  const updateKinds = useCallback((kinds: readonly POIKind[]) => {
    setFilter(prev => {
      const uniqueKinds = Array.from(new Set(kinds));
      if (
        uniqueKinds.length === prev.kinds.length &&
        uniqueKinds.every(k => prev.kinds.includes(k))
      ) {
        return prev; // No change, avoid unnecessary state update
      }
      return {
        ...prev,
        kinds: uniqueKinds,
        enabled: uniqueKinds.length > 0,
      };
    });
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
    setFilter(prev => {
      if (prev.kinds.includes(kind)) {
        return prev;
      }
      return {
        ...prev,
        kinds: [...prev.kinds, kind],
        enabled: true,
      };
    });
  }, []);

  // Remove kind from filter
  const removeKind = useCallback((kind: POIKind) => {
    setFilter(prev => {
      const newKinds = prev.kinds.filter(k => k !== kind);
      return {
        ...prev,
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
  const filterExpression = useMemo((): FilterExpression => {
    if (!filter.enabled || filter.kinds.length === 0) {
      return null;
    }

    return [
      "all",
      [
        "match",
        ["get", "kind"],
        filter.kinds,
        true,
        false,
      ],
      ["has", "name"],
      ["!=", ["get", "name"], ""]
    ];
  }, [filter.enabled, filter.kinds]);

  return {
    filter,
    updateKinds,
    toggleEnabled,
    addKind,
    removeKind,
    clearFilter,
    filterExpression,
  };
};