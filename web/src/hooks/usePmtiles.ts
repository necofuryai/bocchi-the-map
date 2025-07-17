import { useEffect } from 'react';
import maplibregl from 'maplibre-gl';
import { Protocol } from 'pmtiles';

// Protocol reference counting to handle multiple map instances
let protocolRefCount = 0;
let protocolInstance: Protocol | null = null;

export const usePmtiles = () => {
  useEffect(() => {
    try {
      // Register protocol only if it's the first instance
      if (protocolRefCount === 0) {
        console.log("Initializing PMTiles protocol...");
        
        // Initialize with metadata: true for better attribution and inspector support
        // This requires an additional blocking HTTP request but provides better error handling
        protocolInstance = new Protocol({
          metadata: true
        });
        
        // Register the protocol with MapLibre GL
        maplibregl.addProtocol('pmtiles', protocolInstance.tile.bind(protocolInstance));
        
        console.log("PMTiles protocol registered successfully");
      }
      protocolRefCount++;
      
      return () => {
        // Only remove protocol when the last instance unmounts
        protocolRefCount--;
        if (protocolRefCount === 0 && protocolInstance) {
          console.log("Removing PMTiles protocol...");
          try {
            maplibregl.removeProtocol('pmtiles');
            protocolInstance = null;
            console.log("PMTiles protocol removed successfully");
          } catch (error) {
            console.error('Failed to remove PMTiles protocol:', error);
          }
        }
      };
    } catch (error) {
      console.error('Failed to initialize PMTiles protocol:', error);
      // Log additional details for debugging
      console.error('Error details:', {
        message: error instanceof Error ? error.message : 'Unknown error',
        stack: error instanceof Error ? error.stack : undefined
      });
    }
  }, []);
};