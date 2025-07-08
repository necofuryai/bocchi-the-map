import { useCallback, useMemo } from 'react';
import { useMapStore } from '@/stores/use-map-store';

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
  const { filters, setFilters } = useMapStore();
  
  // Get current filter state from store
  const filter: MapFilter = {
    kinds: (filters.kinds as POIKind[]) || [...initialKinds],
    enabled: filters.enabled ?? initialKinds.length > 0,
  };

  // Update filter kinds
  const updateKinds = useCallback((kinds: readonly POIKind[]) => {
    const uniqueKinds = Array.from(new Set(kinds));
    const currentKinds = filter.kinds;
    
    if (
      uniqueKinds.length === currentKinds.length &&
      uniqueKinds.every(k => currentKinds.includes(k))
    ) {
      return; // No change, avoid unnecessary state update
    }
    
    setFilters({
      ...filters,
      kinds: uniqueKinds,
      enabled: uniqueKinds.length > 0,
    });
  }, [filter.kinds, filters, setFilters]);

  // Toggle filter enabled state
  const toggleEnabled = useCallback(() => {
    setFilters({
      ...filters,
      enabled: !filter.enabled,
    });
  }, [filters, filter.enabled, setFilters]);

  // Add kind to filter
  const addKind = useCallback((kind: POIKind) => {
    if (filter.kinds.includes(kind)) {
      return;
    }
    
    setFilters({
      ...filters,
      kinds: [...filter.kinds, kind],
      enabled: true,
    });
  }, [filter.kinds, filters, setFilters]);

  // Remove kind from filter
  const removeKind = useCallback((kind: POIKind) => {
    const newKinds = filter.kinds.filter(k => k !== kind);
    
    setFilters({
      ...filters,
      kinds: newKinds,
      enabled: newKinds.length > 0,
    });
  }, [filter.kinds, filters, setFilters]);

  // Clear all filters
  const clearFilter = useCallback(() => {
    setFilters({
      ...filters,
      kinds: [],
      enabled: false,
    });
  }, [filters, setFilters]);

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