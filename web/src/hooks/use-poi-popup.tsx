import React, { useEffect, useRef, useCallback } from "react"
import { createRoot } from "react-dom/client"
import * as maplibregl from "maplibre-gl"
import { useMapStore } from "@/stores/use-map-store"
import { POIPopup } from "@/components/poi-popup"
import type { Spot } from "@/types"

interface UsePOIPopupOptions {
  map: maplibregl.Map | null
}

export const usePOIPopup = ({ map }: UsePOIPopupOptions) => {
  const popupRef = useRef<maplibregl.Popup | null>(null)
  const containerRef = useRef<HTMLDivElement | null>(null)
  const rootRef = useRef<ReturnType<typeof createRoot> | null>(null)
  const { popup, openPopup, closePopup } = useMapStore()

  // Create popup instance
  useEffect(() => {
    if (!map) return

    // Create popup container
    const container = document.createElement("div")
    container.className = "poi-popup-container"
    containerRef.current = container

    // Create React root
    rootRef.current = createRoot(container)

    // Create MapLibre popup
    popupRef.current = new maplibregl.Popup({
      closeButton: false,
      closeOnClick: false,
      closeOnMove: false,
      maxWidth: "380px",
      className: "poi-popup",
      anchor: "bottom",
      offset: [0, -10] as [number, number],
    })

    // Handle popup close
    popupRef.current.on("close", () => {
      closePopup()
    })

    return () => {
      if (rootRef.current) {
        rootRef.current.unmount()
        rootRef.current = null
      }
      if (popupRef.current) {
        popupRef.current.remove()
        popupRef.current = null
      }
      if (containerRef.current) {
        containerRef.current.remove()
        containerRef.current = null
      }
    }
  }, [map, closePopup])

  // Handle popup state changes
  useEffect(() => {
    if (!popupRef.current || !containerRef.current || !rootRef.current) return

    if (popup.isOpen && popup.spot && popup.position) {
      // Render React component
      rootRef.current.render(<POIPopup />)
      
      // Set popup content and position
      popupRef.current
        .setDOMContent(containerRef.current)
        .setLngLat(popup.position)
        .addTo(map!)
    } else {
      // Close popup
      if (popupRef.current.isOpen()) {
        popupRef.current.remove()
      }
    }
  }, [popup.isOpen, popup.spot, popup.position, map])

  // Handle POI click events
  const handlePOIClick = useCallback((event: maplibregl.MapMouseEvent) => {
    if (!map) return

    // Query for POI features at click point
    const features = map.queryRenderedFeatures(event.point, {
      layers: ["poi-icons-priority", "poi-icons"], // Handle both priority and regular POI layers
    })

    if (features.length > 0) {
      const feature = features[0]
      const { properties, geometry } = feature

      // Extract POI data from feature properties
      if (properties && geometry.type === "Point") {
        const spot: Spot = {
          id: properties.id || `poi-${Date.now()}`, // Generate ID if not available
          name: properties.name || "Unknown Place",
          category: properties.class || "Unknown",
          address: properties.address || "Address not available",
          countryCode: properties.country_code || "XX",
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        }

        const position: [number, number] = [
          geometry.coordinates[0],
          geometry.coordinates[1],
        ]

        // Open popup with POI data
        openPopup(spot, position)
      }
    }
  }, [map, openPopup])

  // Handle clicking outside POI to close popup
  const handleMapClick = useCallback((event: maplibregl.MapMouseEvent) => {
    if (!map) return

    // Check if click is on a POI
    const features = map.queryRenderedFeatures(event.point, {
      layers: ["poi-icons-priority", "poi-icons"],
    })

    if (features.length === 0 && popup.isOpen) {
      // Click outside POI, close popup
      closePopup()
    } else if (features.length > 0) {
      // Click on POI, handle POI click
      handlePOIClick(event)
    }
  }, [map, popup.isOpen, closePopup, handlePOIClick])

  // Setup event listeners
  useEffect(() => {
    if (!map) return

    // Add click event listener
    map.on("click", handleMapClick)

    // Add cursor styling for POI layers
    const setCursorPointer = () => {
      map.getCanvas().style.cursor = "pointer"
    }

    const setCursorDefault = () => {
      map.getCanvas().style.cursor = ""
    }

    map.on("mouseenter", "poi-icons", setCursorPointer)
    map.on("mouseleave", "poi-icons", setCursorDefault)
    map.on("mouseenter", "poi-icons-priority", setCursorPointer)
    map.on("mouseleave", "poi-icons-priority", setCursorDefault)

    return () => {
      map.off("click", handleMapClick)
      map.off("mouseenter", "poi-icons", setCursorPointer)
      map.off("mouseleave", "poi-icons", setCursorDefault)
      map.off("mouseenter", "poi-icons-priority", setCursorPointer)
      map.off("mouseleave", "poi-icons-priority", setCursorDefault)
    }
  }, [map, handleMapClick])

  return {
    popup: popupRef.current,
    isOpen: popup.isOpen,
    handlePOIClick,
    closePopup,
  }
}